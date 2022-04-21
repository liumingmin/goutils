package es7

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

	"github.com/olivere/elastic/v7"
)

type Client struct {
	*elastic.Client
}

var esClients = make(map[string]*Client)

func InitClients() {
	dbs := conf.Conf.Databases
	if dbs == nil {
		fmt.Fprintf(os.Stderr, "No database configuration")
		return
	}

	for _, database := range dbs {
		if database.Type == db.ES7 {
			client, err := initClient(database)
			if err != nil {
				continue
			}

			esClients[database.Key] = client
		}
	}
}

func initClient(dbconf *conf.Database) (ret *Client, err error) {
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
	return &Client{client}, err
}

// get ES 封装客户端
//client := GetEsClient(keyName)
//if client == nil {
//log.Error(ctx, "ES client is nil")
//return errors.New("ES client is nil")
//}

func GetEsClient(key string) *Client {
	client, _ := esClients[key]
	return client
}

func (t *Client) FindByModel(ctx context.Context, model elasticsearch.QueryModel) error {
	searchService := t.Search().Index(model.IndexName)

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

	return t.find(ctx, searchService, model.Results, model.Total)
}

func (t *Client) FindBySource(ctx context.Context, model elasticsearch.SourceModel) error {
	searchService := t.Search().Index(model.IndexName)

	searchService = searchService.Source(model.Source)

	return t.find(ctx, searchService, model.Results, model.Total)
}

func (t *Client) AggregateBySource(ctx context.Context, model elasticsearch.AggregateModel, result ...interface{}) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("agg error: %v", err)
	})

	searchService := t.Search().Index(model.IndexName)

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
				log.Error(ctx, "RawData marshalJSON err: %v", err1)
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

func (t *Client) find(ctx context.Context, searchService *elastic.SearchService, results interface{}, total *int64) (err error) {
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

func (t *Client) Insert(ctx context.Context, esIndexName, id string, data interface{}) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("Insert error: %v", err)
	})

	_, err = t.Index().Index(esIndexName).Id(id).BodyJson(data).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *Client) BatchInsert(ctx context.Context, esIndexName string, ids []string, items []interface{}) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("BatchInsert error: %v", err)
	})

	bulkService := t.Bulk()

	for i, id := range ids {
		bulkService.Add(elastic.NewBulkIndexRequest().Index(esIndexName).
			Id(id).Doc(items[i]))
	}

	_, err = bulkService.Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *Client) UpdateById(ctx context.Context, esIndexName, id string, updateM map[string]interface{}) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("Update error: %v", err)
	})

	var b strings.Builder
	for field, _ := range updateM {
		fmt.Fprintf(&b, "ctx._source.%s=params.%s;", field, field)
	}
	_, err = t.Update().Index(esIndexName).Id(id).
		Script(elastic.NewScriptInline(b.String()).Lang("painless").Params(updateM)).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *Client) DeleteById(ctx context.Context, esIndexName, id string) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("Delete error: %v", err)
	})

	_, err = t.Delete().Index(esIndexName).Id(id).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (t *Client) CreateIndexByMapping(ctx context.Context, esIndexName, esMapping string) (err error) {
	defer log.Recover(ctx, func(e interface{}) string {
		err = fmt.Errorf("%v", e)
		return fmt.Sprintf("CreateIndexIfNotExist error: %v", err)
	})

	exist, _ := t.IndexExists(esIndexName).Do(ctx)
	if exist {
		return nil
	}

	_, err = t.CreateIndex(esIndexName).Body(esMapping).Do(ctx)
	if err != nil {
		log.Error(ctx, "CreateIndexIfNotExist failed, err: %v", err)
	}
	return err
}

func (t *Client) CreateIndexByModel(ctx context.Context, esIndexName string, model *MappingModel) (err error) {
	esMapping, err := json.Marshal(model)
	if err != nil {
		log.Error(ctx, "CreateIndexByModelIfNotExist failed, err: %v", err)
		return err
	}

	mappingBody := string(esMapping)
	log.Debug(ctx, "CreateIndexByModel mapping is: %v", mappingBody)

	return t.CreateIndexByMapping(ctx, esIndexName, mappingBody)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type MappingModel struct {
	Mapping  `json:"mappings"`
	Settings `json:"settings"`
}

type Mapping struct {
	Dynamic    bool                                      `json:"dynamic"` // false
	Properties map[string]*elasticsearch.MappingProperty `json:"properties"`
}

type Settings struct {
	IndexMappingIgnoreMalformed bool  `json:"index.mapping.ignore_malformed,omitempty"` // true
	NumberOfReplicas            int64 `json:"number_of_replicas"`                       // 1
	NumberOfShards              int64 `json:"number_of_shards"`                         // 3
}
