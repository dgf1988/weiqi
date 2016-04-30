package weiqi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	c_configFilename = "config.json"

	c_defaultConfigJson = `
{
  "BasePath":"d:/Project/src/github.com/dgf1988/weiqi/",
  "SiteTitle":"围棋163",
  "SiteDesc":"围棋综合网站",
  "SiteKeywords": ["围棋", "棋谱", "棋手", "新闻", "文章", "资料"],

  "SiteAuthorName":"DGF",
  "SiteAuthorUrl":"http://www.dingguofeng.com",
  "SiteAuthorEmail":"dgf1988@qq.com",
  "SiteICP":"闽ICP备14014166号-2",

  "DbDriver":"mysql",
  "DbUsername":"root",
  "DbPassword":"guofeng001",
  "DbHost":"localhost",
  "DbPort":3306,
  "DbName":"weiqi_2",
  "DbCharset":"utf8",

  "UploadPath":"D:/upload/"
}
	`
)

type Config struct {
	BasePath string

	SiteTitle    string
	SiteDesc     string
	SiteKeywords []string

	SiteAuthorName  string
	SiteAuthorUrl   string
	SiteAuthorEmail string
	SiteICP         string

	DbDriver   string
	DbUsername string
	DbPassword string
	DbHost     string
	DbPort     int
	DbName     string
	DbCharset  string

	UploadPath string
}

var (
	config *Config
)

func init() {
	c, err := loadConfig(c_configFilename)
	if err != nil || c == nil {
		panic(err.Error())
	}
	config = c
}

func defConfig() *Config {
	var c = &Config{}
	err := json.Unmarshal([]byte(c_defaultConfigJson), c)
	if err != nil {
		panic(err.Error())
	}
	return c
}

func loadConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	jsbyte, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	var c = &Config{}
	err = json.Unmarshal(jsbyte, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//"root:guofeng001@tcp(localhost:3306)/weiqi2?charset=utf8&parseTime=true"
func (c Config) DbConnectString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true",
		c.DbUsername, c.DbPassword, c.DbHost, c.DbPort, c.DbName, c.DbCharset)
}

func (c *Config) ToJson() string {
	jsbyte, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(jsbyte)
}

func (c *Config) Save(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(c.ToJson())
	return err
}
