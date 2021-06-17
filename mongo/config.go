package mongo

import (
	"fmt"
	"strings"
	"time"

	"github.com/demdxx/gocast"
	"github.com/liumingmin/goutils/conf"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Safe struct {
	W        int    // Min # of servers to ack before success
	WMode    string // Write mode for MongoDB 2.0+ (e.g. "majority")
	RMode    string // Read mode for MonogDB 3.2+ ("majority", "local", "linearizable")
	WTimeout int    // Milliseconds to wait for W before timing out
	FSync    bool   // Sync via the journal if present, or via data files sync otherwise
	J        bool   // Sync via the journal if present
}

type Config struct {
	// 连接URL, 格式为[mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	Address     []string
	Database    string
	Username    string
	Password    string
	Source      string
	Safe        *Safe
	Mode        readpref.Mode
	Compressors []string

	// 连接管理
	Direct         bool
	ConnectTimeout time.Duration //连接超时. 默认为10秒
	Keepalive      time.Duration //DialInfo.DialServer实现
	WriteTimeout   time.Duration //写超时, 默认为ConnectTimeout
	ReadTimeout    time.Duration //读超时, 默认为ConnectTimeout

	// 连接池管理
	MinPoolSize       int //对应DialInfo.MinPoolSize
	MaxPoolSize       int //对应DialInfo.PoolLimit
	MaxPoolWaitTimeMS int //对应DialInfo.PoolTimeout获取连接超时, 默认为0永不超时
	MaxPoolIdleTimeMS int //对应DialInfo.MaxIdleTimeMS
}

func newConfig(dbconf *conf.Database) *Config {
	safeConf, _ := dbconf.Ext("safe", map[interface{}]interface{}{}).(map[interface{}]interface{})
	if safeConf == nil {
		safeConf = map[interface{}]interface{}{}
	}
	safe := getSafe(safeConf)

	modeConf, _ := dbconf.Ext("mode", "primary").(string)
	mode := getMode(modeConf)

	direct, _ := dbconf.Ext("direct", false).(bool)

	return &Config{
		Address:           strings.Split(dbconf.Host, ","),
		Database:          dbconf.Name,
		Username:          dbconf.User,
		Password:          dbconf.Password,
		Source:            dbconf.ExtString("authSource", dbconf.Name),
		Safe:              safe,
		Mode:              mode,
		Direct:            direct,
		Keepalive:         dbconf.ExtDuration("keepalive", ""),
		ConnectTimeout:    dbconf.ExtDuration("connectTimeout", "10s"),
		ReadTimeout:       dbconf.ExtDuration("readTimeout", ""),
		WriteTimeout:      dbconf.ExtDuration("writeTimeout", ""),
		MinPoolSize:       dbconf.ExtInt("minPoolSize", 0),
		MaxPoolSize:       dbconf.ExtInt("maxPoolSize", 10),
		MaxPoolWaitTimeMS: dbconf.ExtInt("maxPoolWaitTimeMS", 0),
		MaxPoolIdleTimeMS: dbconf.ExtInt("maxPoolIdleTimeMS", 0),
	}
}

func getMode(val interface{}) readpref.Mode {
	switch val := val.(type) {
	case string:
		ret, err := readpref.ModeFromString(val)
		if err != nil {
			return readpref.PrimaryMode
		}
		return ret
	case int:
		return readpref.Mode(val)
	case int64:
		return readpref.Mode(val)
	}
	panic("unsupport mode type: " + fmt.Sprint(val))
}

func getSafe(val map[interface{}]interface{}) *Safe {
	safe := &Safe{
		W: 1,
	}
	for k, v := range val {
		switch k {
		case "W", "w":
			safe.W = gocast.ToInt(v)
		case "WMode", "wmode":
			safe.WMode = gocast.ToString(v)
		case "RMode", "rmode":
			safe.RMode = gocast.ToString(v)
		case "WTimeout", "wtimeout":
			safe.WTimeout = gocast.ToInt(v)
		case "FSync", "fsync":
			safe.FSync = gocast.ToBool(v)
		case "J", "j":
			safe.J = gocast.ToBool(v)
		}
	}
	return safe
}

func InitClient(dbconf *conf.Database) (ret *Client, err error) {
	return newClient(newConfig(dbconf))
}
