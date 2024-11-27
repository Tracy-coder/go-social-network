// from: https://github.com/chenghonour/formulago

package configs

import (
	"embed"
	"os"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/jinzhu/configor"
)

//go:embed *.yaml
var configFiles embed.FS

// GlobalConfig .
var globalConfig Config
var isInit = false

func InitConfig() {
	// log print embed config file
	DirEntry, err := configFiles.ReadDir(".")
	if err != nil {
		hlog.Fatal("read config embed dir error: ", err)
	}
	for _, v := range DirEntry {
		hlog.Info("embed config file: ", v.Name())
	}
	// load config
	globalConfig, _ = load()
}

func Data() Config {
	if !isInit {
		InitConfig()
		isInit = true
	}
	return globalConfig
}

func ReLoad() {
	globalConfig, _ = load()
}

// load from config.yaml (embed)
func load() (config Config, err error) {
	// Obtain whether it is prod environment from the environment variable
	IsProd := os.Getenv("IS_PROD")
	// If it is prod environment, load the prod configuration file
	if IsProd == "true" {
		hlog.Info("load prod config")
		err = configor.New(&configor.Config{ErrorOnUnmatchedKeys: true, FS: configFiles}).
			Load(&config, "config.yaml")
		if err != nil {
			hlog.Fatal(err)
		}
		return
	}
	// If it is not prod environment, load the dev configuration file
	hlog.Info("load dev config")
	err = configor.New(&configor.Config{ErrorOnUnmatchedKeys: true, FS: configFiles}).
		Load(&config, "config_dev.yaml")
	if err != nil {
		hlog.Fatal(err)
	}
	return
}

// Config is the configuration of the project.
type Config struct {
	Name     string     `yaml:"Name"`
	IsDemo   bool       `yaml:"IsDemo"`
	IsProd   bool       `yaml:"IsProd"`
	Host     string     `yaml:"Host"`
	Port     int        `yaml:"Port"`
	Timeout  int        `yaml:"Timeout"`
	Captcha  Captcha    `yaml:"Captcha"`
	Redis    Redis      `yaml:"Redis"`
	Database Database   `yaml:"Database"`
	Casbin   CasbinConf `yaml:"Casbin"`
	Kafka    KafkaConf  `yaml:"Kafka"`
}

// Captcha is the configuration of the captcha.
type Captcha struct {
	KeyLong   int `yaml:"KeyLong"`
	ImgWidth  int `yaml:"ImgWidth"`
	ImgHeight int `yaml:"ImgHeight"`
}

// Redis is the configuration of the redis.
type Redis struct {
	Enable   bool   `yaml:"Enable"`
	Host     string `yaml:"Host"`
	Port     int    `yaml:"Port"`
	Type     string `yaml:"Type"`
	Password string `yaml:"Password"`
	DB       int    `yaml:"DB"`
}

// Database is the configuration of the database.
type Database struct {
	Type        string      `yaml:"Type"`
	Host        string      `yaml:"Host"`
	Port        int         `yaml:"Port"`
	DBName      string      `yaml:"DBName"`
	Username    interface{} `yaml:"Username"`
	Password    interface{} `yaml:"Password"`
	MaxOpenConn int         `yaml:"MaxOpenConn"`
	SSLMode     bool        `yaml:"SSLMode"`
	CacheTime   int         `yaml:"CacheTime"`
}

// CasbinConf is the configuration of the casbin.
type CasbinConf struct {
	ModelText string `yaml:"ModelText"`
}

type KafkaConf struct {
	Brokers []string `yaml:"Brokers"`
	Topic   string   `yaml:"Topic"`
	GroupID string   `yaml:"GroupID"`
}
