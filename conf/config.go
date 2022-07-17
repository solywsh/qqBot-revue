package conf

import (
	"github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"sync"
)

var (
	configOnce sync.Once
)

// Config 配置相关
type Config struct {
	ListenPort            string                `yaml:"listenPort"`            // 监听端口
	AdminUserOrderHeader  string                `yaml:"adminUserOrderHeader"`  // 管理员命令头 adminUserOrderHeader
	ListenGroup           []string              `yaml:"listenGroup"`           // 监听群列表
	ForwardAuthentication ForwardAuthentication `yaml:"forwardAuthentication"` // 正向鉴权 forward authentication
	ReverseAuthentication ReverseAuthentication `yaml:"reverseAuthentication"` // 反向鉴权 reverse authentication
	Revue                 Revue                 `yaml:"revue"`                 // revue相关
	UrlHeader             string                `yaml:"urlHeader"`             // url
	SelfId                string                `yaml:"selfId"`                // 机器人的qq
	AdminUser             []string              `yaml:"adminUser"`             // 管理员列表
	Database              Database              `yaml:"Database"`              // 数据库相关
}

// ForwardAuthentication 正向鉴权相关
type ForwardAuthentication struct {
	Enable bool   `yaml:"enable"`
	Token  string `yaml:"token"`
}

// ReverseAuthentication 反向鉴权相关
type ReverseAuthentication struct {
	Enable bool   `yaml:"enable"`
	Secret string `yaml:"secret"`
}

// Revue 相关
type Revue struct {
	Enable bool `yaml:"enable"`
}

// Database 数据库相关
type Database struct {
	Sqlite Sqlite `yaml:"sqlite"`
	Mysql  Mysql  `yaml:"mysql"`
	Mongo  Mongo  `yaml:"mongo"`
}

// Sqlite 数据库
type Sqlite struct {
	Enable bool   `yaml:"enable"`
	Path   string `yaml:"path"`
}

// Mysql 数据库
type Mysql struct {
	Dbname   string `yaml:"dbname"`
	Charset  string `yaml:"charset"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Enable   bool   `yaml:"enable"`
	Address  string `yaml:"address"`
}

// Mongo 数据库相关服务配置
type Mongo struct {
	HImgDB HImgDB `yaml:"hImgDB"`
}

// HImgDB 涩图数据库
type HImgDB struct {
	Enable bool   `yaml:"enable"`
	Url    string `yaml:"url"`
}

func NewConf(yamlPath string) (conf *Config, err error) {
	configOnce.Do(func() {
		conf = &Config{}
		yamlFile, err := ioutil.ReadFile(yamlPath)
		if err != nil {
			log.Println(err)
			return
		}
		err = yaml.Unmarshal(yamlFile, conf)
		if err != nil {
			log.Println(err)
			return
		}
	})
	return conf, err
}
