package conf

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"time"

	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

var confPath = os.Getenv("GOUTILS_CONF_PATH")

type Log struct {
	lumberjack.Logger `yaml:",inline"`
	OutputEncoder     string `yaml:"outputEncoder"`
	Stdout            bool   `yaml:"stdout"`
	FileOut           bool   `yaml:"fileOut"`
	HttpOut           bool   `yaml:"httpOut"`
	HttpUrl           string `yaml:"httpUrl"`     // 请求地址
	HttpDebug         bool   `yaml:"httpDebug"`   // 是否打印请求输出成功和失败的情况
	HttpTimeout       int    `yaml:"httpTimeout"` // httpclient超时
}

type Database struct {
	Key      string                 `yaml:"key"`
	Addrs    []string               `yaml:"addrs,flow"`
	User     string                 `yaml:"user"`
	Password string                 `yaml:"password"`
	DBName   string                 `yaml:"dbName"`
	EXT      map[string]interface{} `yaml:",flow"`
}

type DBSsh struct {
	On      bool   `yaml:"on"`
	Address string `yaml:"address"` //ip:port
	User    string `yaml:"user"`    //ssh user
	PriKey  string `yaml:"priKey"`  //base64
	KeyPass string `yaml:"keyPass"` //key password
}

type Mongo struct {
	Database   `yaml:",inline"`
	AuthSource string `yaml:"authSource"`

	// 连接管理
	Direct         bool          `yaml:"direct"`
	ConnectTimeout time.Duration `yaml:"connectTimeout"` //连接超时. 默认为10秒
	Keepalive      time.Duration `yaml:"keepalive"`      //DialInfo.DialServer实现
	WriteTimeout   time.Duration `yaml:"writeTimeout"`   //写超时, 默认为ConnectTimeout
	ReadTimeout    time.Duration `yaml:"readTimeout"`    //读超时, 默认为ConnectTimeout
	Compressors    []string      `yaml:"compressors,flow"`

	// 连接池管理
	MinPoolSize       int `yaml:"minPoolSize"`       //对应DialInfo.MinPoolSize
	MaxPoolSize       int `yaml:"maxPoolSize"`       //对应DialInfo.PoolLimit
	MaxPoolWaitTimeMS int `yaml:"maxPoolWaitTimeMS"` //对应DialInfo.PoolTimeout获取连接超时, 默认为0永不超时
	MaxPoolIdleTimeMS int `yaml:"maxPoolIdleTimeMS"` //对应DialInfo.MaxIdleTimeMS

	//读写偏好
	Mode string     `yaml:"mode"`
	Safe *MongoSafe `yaml:"safe"`

	//ssh tune(experimental)
	Ssh *DBSsh `yaml:"ssh"`
}

type MongoSafe struct {
	W        int    `yaml:"w"`        // Min # of servers to ack before success
	WMode    string `yaml:"wMode"`    // Write mode for MongoDB 2.0+ (e.g. "majority")
	RMode    string `yaml:"rMode"`    // Read mode for MonogDB 3.2+ ("majority", "local", "linearizable")
	WTimeout int    `yaml:"wTimeout"` // Milliseconds to wait for W before timing out
	FSync    bool   `yaml:"fSync"`    // Sync via the journal if present, or via data files sync otherwise
	J        bool   `yaml:"j"`        // Sync via the journal if present
}

type Elasticsearch struct {
	Database    `yaml:",inline"`
	Type        string `yaml:"type"`
	MaxPoolSize int    `yaml:"maxPoolSize"`
}

type Redis struct {
	Key              string   `yaml:"key"`
	MasterName       string   `yaml:"masterName"`
	Addrs            []string `yaml:"addrs,flow"`
	Db               int      `yaml:"db"`
	PoolSize         int      `yaml:"poolSize"`
	Password         string   `yaml:"password"`
	SentinelPassword string   `yaml:"sentinelPassword"`
	DialTimeout      string   `yaml:"dialTimeout"`
	ReadTimeout      string   `yaml:"readTimeout"`
	WriteTimeout     string   `yaml:"writeTimeout"`
	IdleTimeout      string   `yaml:"idleTimeout"`
	ReadOnly         bool     `yaml:"readOnly"`
	RouteByLatency   bool     `yaml:"routeByLatency"`
	RouteRandomly    bool     `yaml:"routeRandomly"`
}

type KafkaProducer struct {
	Key           string   `yaml:"key"`
	Address       []string `yaml:"address"`
	Async         bool     `yaml:"async"`
	ReturnSuccess bool     `yaml:"returnSuccess"`
	ReturnError   bool     `yaml:"returnError"`
	//username and password for SASL/PLAIN  or SASL/SCRAM authentication
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type KafkaConsumer struct {
	Key          string        `yaml:"key"`
	Address      []string      `yaml:"address"` // kafka地址
	Group        string        `yaml:"group"`   // groupId
	Offset       int64         `yaml:"offset"`
	Ack          int           `yaml:"ack"`          // ack类型
	DialTimeout  time.Duration `yaml:"dialTimeout"`  // How long to wait for the initial connection.
	ReadTimeout  time.Duration `yaml:"readTimeout"`  // How long to wait for a response.
	WriteTimeout time.Duration `yaml:"writeTimeout"` // How long to wait for a transmit.
	KeepAlive    time.Duration `yaml:"keepAlive"`
	//username and password for SASL/PLAIN  or SASL/SCRAM authentication
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// Get value of given key of Database section.
func (d *Database) Ext(key string, defaultVal ...interface{}) interface{} {
	if v, exist := d.EXT[key]; exist {
		return v
	}

	if len(defaultVal) > 0 {
		return defaultVal[0]
	}

	return ""
}

func (d *Database) ExtString(key string, defaultVal ...interface{}) string {
	keyValue := d.Ext(key, defaultVal...)

	return fmt.Sprint(keyValue)
}

func (d *Database) ExtInt(key string, defaultVal ...interface{}) int {
	return d.Ext(key, defaultVal...).(int)
}

func (d *Database) ExtDuration(key string, defaultVal ...interface{}) time.Duration {
	keyValue := d.Ext(key, defaultVal...)

	str := fmt.Sprint(keyValue)
	t, _ := time.ParseDuration(str)
	return t
}

func (d *Database) ExtBool(key string, defaultVal ...interface{}) bool {
	return d.Ext(key, defaultVal...).(bool)
}

type Config struct {
	LogLevel        string                 `yaml:"logLevel"`
	Logs            []*Log                 `yaml:"logs,flow"`
	Mongos          []*Mongo               `yaml:"mongos,flow"`
	Elasticsearches []*Elasticsearch       `yaml:"elasticsearches,flow"`
	Redises         []*Redis               `yaml:"redises,flow"`
	KafkaProducers  []*KafkaProducer       `yaml:"kafkaProducers,flow"`
	KafkaConsumers  []*KafkaConsumer       `yaml:"kafkaConsumers,flow"`
	EXT             map[string]interface{} `yaml:"ext,flow"`
}

// Ext will return the value of the EXT config, the keys is a string
// separated by DOT(.). If you provide a default value, this method
// will return the it while the key cannot be found. otherwise it
// will raise a panic!
func (c *Config) Ext(keys string, defaultVal ...interface{}) interface{} {
	r, e := c.ExtSep(keys, ".")
	if e != nil || r == nil {
		if len(defaultVal) > 0 {
			return defaultVal[0]
		} else {
			panic(e)
		}
	} else {
		return r
	}
}

func (c *Config) ExtString(keys string, defaultVal ...interface{}) string {
	return fmt.Sprintf("%v", c.Ext(keys, defaultVal...))
}

func (c *Config) ExtInt(keys string, defaultVal ...interface{}) int {
	return c.Ext(keys, defaultVal...).(int)
}

func (c *Config) ExtInt8(keys string, defaultVal ...interface{}) int8 {
	return int8(c.Ext(keys, defaultVal...).(int))
}

func (c *Config) ExtInt16(keys string, defaultVal ...interface{}) int16 {
	return int16(c.Ext(keys, defaultVal...).(int))
}

func (c *Config) ExtInt32(keys string, defaultVal ...interface{}) int32 {
	return int32(c.Ext(keys, defaultVal...).(int))
}

func (c *Config) ExtInt64(keys string, defaultVal ...interface{}) int64 {
	return int64(c.Ext(keys, defaultVal...).(int))
}

func (c *Config) ExtBool(keys string, defaultVal ...interface{}) bool {
	return c.Ext(keys, defaultVal...).(bool)
}

func (c *Config) ExtFloat64(keys string, defaultVal ...interface{}) float64 {
	return c.Ext(keys, defaultVal...).(float64)
}

func (c *Config) ExtFloat32(keys string, defaultVal ...interface{}) float32 {
	return c.Ext(keys, defaultVal...).(float32)
}

func (c *Config) ExtDuration(keys string, defaultVal ...interface{}) time.Duration {
	str := fmt.Sprintf("%v", c.Ext(keys, defaultVal...))
	t, _ := time.ParseDuration(str)
	return t
}

// Ext will return the value of the EXT config, the keys is separated
// by the given sep string.
func (c *Config) ExtSep(keys, sep string) (interface{}, error) {
	ks := strings.Split(keys, sep)
	var result interface{}
	var isFinal, success bool
	result = c.EXT
	for _, k := range ks {
		result, isFinal, success = find(result, k)
		if !success {
			return "", fmt.Errorf("no such key: %v", k)
		} else if isFinal {
			break
		}
	}

	if success {
		return result, nil
	} else {
		return "", fmt.Errorf("not found")
	}
}

func find(v interface{}, key interface{}) (result interface{}, isFinal, success bool) {
	switch m := v.(type) {
	case map[string]interface{}:
		result, success = m[key.(string)]
		isFinal = reflect.TypeOf(result) != nil && reflect.TypeOf(result).Kind() != reflect.Map
	case map[interface{}]interface{}:
		result, success = m[key]
		isFinal = reflect.TypeOf(result) != nil && reflect.TypeOf(result).Kind() != reflect.Map
	}
	return
}

var (
	Conf = Config{}
)

func LoadConf(path string) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		return
	}
	if err = yaml.Unmarshal(c, &Conf); err != nil {
		fmt.Fprint(os.Stderr, err)
	}
}

func Ext(keys string, defVal ...interface{}) interface{} {
	return Conf.Ext(keys, defVal...)
}

// ExtString is a shorcut for (c *Config) ExtString()
func ExtString(keys string, defVal ...interface{}) string {
	return Conf.ExtString(keys, defVal...)
}

// ExtBool is a shorcut for (c *Config) ExtBool()
func ExtBool(keys string, defVal ...interface{}) bool {
	return Conf.ExtBool(keys, defVal...)
}

// ExtInt is a shorcut for (c *Config) ExtInt()
func ExtInt(keys string, defVal ...interface{}) int {
	return Conf.ExtInt(keys, defVal...)
}

// ExtInt8 is a shorcut for (c *Config) ExtInt8()
func ExtInt8(keys string, defVal ...interface{}) int8 {
	return Conf.ExtInt8(keys, defVal...)
}

// ExtInt16 is a shorcut for (c *Config) ExtInt16()
func ExtInt16(keys string, defVal ...interface{}) int16 {
	return Conf.ExtInt16(keys, defVal...)
}

// ExtInt32 is a shorcut for (c *Config) ExtInt32()
func ExtInt32(keys string, defVal ...interface{}) int32 {
	return Conf.ExtInt32(keys, defVal...)
}

// ExtInt64 is a shorcut for (c *Config) ExtInt64()
func ExtInt64(keys string, defVal ...interface{}) int64 {
	return Conf.ExtInt64(keys, defVal...)
}

// ExtFloat32 is a shorcut for (c *Config) ExtFloat32()
func ExtFloat32(keys string, defVal ...interface{}) float32 {
	return Conf.ExtFloat32(keys, defVal...)
}

// ExtFloat64 is a shorcut for (c *Config) ExtFloat64()
func ExtFloat64(keys string, defVal ...interface{}) float64 {
	return Conf.ExtFloat64(keys, defVal...)
}

// ExtDuration is a shorcut for (c *Config) ExtDuration()
func ExtDuration(keys string, defVal ...interface{}) time.Duration {
	return Conf.ExtDuration(keys, defVal...)
}

func init() {
	if confPath == "" {
		confPath = "conf.yml"
	}
	LoadConf(confPath)
}
