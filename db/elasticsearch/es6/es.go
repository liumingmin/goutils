package es6

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/liumingmin/goutils/conf"
	"github.com/liumingmin/goutils/db"
	"github.com/liumingmin/goutils/db/elasticsearch"
	"github.com/liumingmin/goutils/log"

	"github.com/olivere/elastic"
)

var esClients = make(map[string]*elastic.Client)

func InitClients() {
	dbs := conf.Conf.Databases
	if dbs == nil {
		fmt.Fprintf(os.Stderr, "No database configuration")
		return
	}

	for _, database := range dbs {
		if database.Type == db.ES6 {
			client, err := initClient(database)
			if err != nil {
				continue
			}

			esClients[database.Key] = client
		}
	}
}

func initClient(dbconf *conf.Database) (ret *elastic.Client, err error) {
	defer log.Recover(context.Background(), func(e interface{}) string {
		err = e.(error)
		return fmt.Sprintf("initClient failed. error: %v", err)
	})

	//初始化
	var options []elastic.ClientOptionFunc
	hosts := strings.Split(dbconf.Host, ",")
	options = append(options, elastic.SetURL(hosts...))

	if dbconf.User != "" {
		options = append(options, elastic.SetBasicAuth(dbconf.User, dbconf.Password))
	}

	sniffEnabled := dbconf.ExtBool("sniffEnabled", true)
	options = append(options, elastic.SetSniff(sniffEnabled))
	if sniffEnabled {
		options = append(options, elastic.SetSnifferInterval(dbconf.ExtDuration("sniffInterval", "15m")))
	}

	healthCheck := dbconf.ExtBool("healthCheck", true)
	options = append(options, elastic.SetHealthcheck(healthCheck))
	if healthCheck {
		options = append(options, elastic.SetHealthcheckInterval(dbconf.ExtDuration("healthCheckInterval", "60s")))
	}

	// http.DefaultClient.set
	maxPerHost := dbconf.MaxPoolSize / len(hosts)
	if maxPerHost > 2 {
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          dbconf.MaxPoolSize,
				MaxIdleConnsPerHost:   maxPerHost,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}
		options = append(options, elastic.SetHttpClient(httpClient))
	}

	client, err := elastic.NewClient(options...)
	if err != nil {
		log.Error(context.Background(), "elastic.NewClient err: %v", err)
	}
	return client, err
}

func GetEsClient(ctx context.Context, key string) *elastic.Client {
	client, _ := esClients[key]
	return client
}

func FindByModel(ctx context.Context, model elasticsearch.QueryModel) error {
	client := GetEsClient(ctx, model.KeyName) // get ES 客户端
	if client == nil {
		log.Error(ctx, "ES client is nil")
		return errors.New("ES client is nil")
	}

	searchService := client.Search().Index(model.IndexName)

	if model.TypeName != "" {
		searchService = searchService.Type(model.TypeName)
	}

	searchService = searchService.Query(model.Query.(elastic.Query))
	if model.Cursor > 0 {
		searchService = searchService.From(model.Cursor)
	}

	if model.Size > 0 {
		searchService = searchService.Size(model.Size)
	}

	for _, field := range model.Sort {
		bSort := true
		switch field[0] {
		case '+':
			field = field[1:]
		case '-':
			bSort = false
			field = field[1:]
		}
		searchService = searchService.Sort(field, bSort)
	}

	return find(ctx, searchService, model.Results, model.Total)
}

func FindBySource(ctx context.Context, model elasticsearch.SourceModel) error {
	client := GetEsClient(ctx, model.KeyName) // get ES 客户端
	if client == nil {
		log.Error(ctx, "ES client is nil")
		return errors.New("ES client is nil")
	}

	searchService := client.Search().Index(model.IndexName)

	if model.TypeName != "" {
		searchService = searchService.Type(model.TypeName)
	}
	searchService = searchService.Source(model.Source)

	return find(ctx, searchService, model.Results, model.Total)
}

func AggregateBySource(ctx context.Context, model elasticsearch.AggregateModel, result ...interface{}) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("agg error: %v", err)
	})

	client := GetEsClient(ctx, model.KeyName) // get ES 客户端
	if client == nil {
		log.Error(ctx, "ES client is nil")
		return errors.New("ES client is nil")
	}

	searchService := client.Search().Index(model.IndexName)

	if model.TypeName != "" {
		searchService = searchService.Type(model.TypeName)
	}
	searchService = searchService.Source(model.Source)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		log.Error(ctx, "SearchService do err: %v", err)
		return err
	}

	for i, aggKey := range model.AggKeys {
		if i >= len(result) {
			return nil
		}

		if rawMsg, ok := searchResult.Aggregations[aggKey]; ok && rawMsg != nil {
			rawData, err1 := rawMsg.MarshalJSON()
			if err1 != nil {
				log.Error(ctx, "RawMsg marshalJSON err: %v", err1)
				continue
			}

			err1 = json.Unmarshal(rawData, result[i])
			if err1 != nil {
				log.Error(ctx, "RawData Unmarshal err: %v", err1)
			}
		}
	}

	return nil
}

func find(ctx context.Context, searchService *elastic.SearchService, results interface{}, total *int64) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("Find error: %v", err)
	})

	destType := reflect.TypeOf(results)
	if destType.Kind() != reflect.Ptr || destType.Elem().Kind() != reflect.Slice {
		return errors.New("dest type must be a slice address")
	}

	destElemType := destType.Elem().Elem()
	var isSliceElemPtr = false
	if destElemType.Kind() == reflect.Ptr {
		destElemType = destElemType.Elem()
		isSliceElemPtr = true
	}

	ptrDestValue := reflect.ValueOf(results)
	destValue := reflect.Indirect(ptrDestValue) //.Elem()
	destValue = destValue.Slice(0, destValue.Cap())

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		return err
	}

	for _, item := range searchResult.Each(destElemType) {
		if isSliceElemPtr {
			destElemValuePtr := reflect.New(destElemType)
			destElemValuePtr.Elem().Set(reflect.ValueOf(item))
			destValue = reflect.Append(destValue, destElemValuePtr)
		} else {
			destValue = reflect.Append(destValue, reflect.ValueOf(item))
		}
	}

	ptrDestValue.Elem().Set(destValue.Slice(0, destValue.Len()))

	if total != nil {
		if searchResult.Hits != nil {
			*(total) = searchResult.TotalHits()
		}
	}

	return nil
}

func Insert(ctx context.Context, esKeyName, esIndexName, esTypeName, id string, data interface{}) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("Insert error: %v", err)
	})

	client := GetEsClient(ctx, esKeyName) // get ES 客户端
	if client == nil {
		log.Error(ctx, "ES client is nil")
		return errors.New("ES client is nil")
	}

	_, err = client.Index().Index(esIndexName).Type(esTypeName).Id(id).BodyJson(data).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func BatchInsert(ctx context.Context, esKeyName, esIndexName, esTypeName string, ids []string, items []interface{}) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("BatchInsert error: %v", err)
	})

	client := GetEsClient(ctx, esKeyName) // get ES 客户端
	if client == nil {
		log.Error(ctx, "ES client is nil")
		return errors.New("ES client is nil")
	}

	bulkService := client.Bulk()

	for i, id := range ids {
		bulkService.Add(elastic.NewBulkIndexRequest().Index(esIndexName).Type(esTypeName).
			Id(id).Doc(items[i]))
	}

	_, err = bulkService.Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func UpdateById(ctx context.Context, esKeyName, esIndexName, esTypeName, id string, updateM map[string]interface{}) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("Update error: %v", err)
	})

	client := GetEsClient(ctx, esKeyName) // get ES 客户端
	if client == nil {
		log.Error(ctx, "ES client is nil")
		return errors.New("ES client is nil")
	}

	var b strings.Builder
	for field, _ := range updateM {
		fmt.Fprintf(&b, "ctx._source.%s=params.%s;", field, field)
	}
	_, err = client.Update().Index(esIndexName).Type(esTypeName).Id(id).
		Script(elastic.NewScriptInline(b.String()).Lang("painless").Params(updateM)).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func DeleteById(ctx context.Context, esKeyName, esIndexName, esTypeName, id string) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("Delete error: %v", err)
	})

	client := GetEsClient(ctx, esKeyName) // get ES 客户端
	if client == nil {
		log.Error(ctx, "ES client is nil")
		return errors.New("ES client is nil")
	}

	_, err = client.Delete().Index(esIndexName).Type(esTypeName).Id(id).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
