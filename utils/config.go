package utils

import (
	"encoding/json"
	"github.com/nbkit/mdf/gmap"
	"github.com/nbkit/mdf/log"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"strings"
)

type AppConfig struct {
	Token   string `mapstructure:"token" json:"token"`
	Code    string `mapstructure:"code" json:"code"`
	Active  string `mapstructure:"active" json:"active"`
	Name    string `mapstructure:"name" json:"name"`
	Port    string `mapstructure:"port" json:"port"`
	Locale  string `mapstructure:"locale" json:"locale"`
	Debug   bool   `mapstructure:"debug" json:"debug"`
	Storage string `mapstructure:"storage" json:"storage"`
	//注册中心
	Registry string `mapstructure:"registry" json:"registry"`
	//服务地址，带端口号
	Address       string `mapstructure:"address" json:"address"`
	PublicAddress string `mapstructure:"public_address" json:"public_address"`
}
type DbConfig struct {
	Driver    string `mapstructure:"driver" json:"driver"`
	Host      string `mapstructure:"host" json:"host"`
	Port      string `mapstructure:"port" json:"port"`
	Database  string `mapstructure:"database" json:"database"`
	Username  string `mapstructure:"username" json:"username"`
	Password  string `mapstructure:"password" json:"password"`
	Charset   string `mapstructure:"charset" json:"charset"`
	Collation string `mapstructure:"collation" json:"collation"`
}
type LogConfig struct {
	Level string `mapstructure:"level" json:"level"`
	Path  string `mapstructure:"path" json:"path"`
	Stack bool   `mapstructure:"stack" json:"stack"`
}
type AuthConfig struct {
	//权限中心地址
	Address string `mapstructure:"address" json:"address"`
	//权限中心编码
	Code string `mapstructure:"code" json:"code"`
}

const AppConfigName = "app"

var Config *EnvConfig
var _ENVMaps = gmap.New()

type jsonConfig struct {
	App  AppConfig              `mapstructure:"app" json:"app"`
	Db   DbConfig               `mapstructure:"db" json:"db"`
	Log  LogConfig              `mapstructure:"log" json:"log"`
	Auth AuthConfig             `mapstructure:"auth" json:"auth"`
	Data map[string]interface{} `json:"data"`
}

type EnvConfig struct {
	App  AppConfig  `mapstructure:"app" json:"app"`
	Db   DbConfig   `mapstructure:"db" json:"db"`
	Log  LogConfig  `mapstructure:"log" json:"log"`
	Auth AuthConfig `mapstructure:"auth" json:"auth"`
	data map[string]interface{}
}

func (c EnvConfig) MarshalJSON() ([]byte, error) {
	jsonMap := jsonConfig{}
	jsonMap.App = c.App
	jsonMap.Db = c.Db
	jsonMap.Log = c.Log
	jsonMap.Auth = c.Auth
	jsonMap.Data = c.data
	return json.Marshal(jsonMap)
}

func (c *EnvConfig) UnmarshalJSON(b []byte) error {
	jsonMap := jsonConfig{}
	json.Unmarshal(b, &jsonMap)
	c.App = jsonMap.App
	c.Db = jsonMap.Db
	c.Log = jsonMap.Log
	c.Auth = jsonMap.Auth
	c.data = jsonMap.Data
	return nil
}
func (s *EnvConfig) GetValue(name string, envNames ...string) string {
	return s.getViper(envNames...).GetString(name)
}
func (s *EnvConfig) UnmarshalValue(name string, rawVal interface{}, envNames ...string) error {
	return s.getViper(envNames...).UnmarshalKey(name, rawVal)
}
func (s *EnvConfig) Unmarshal(rawVal interface{}, envNames ...string) error {
	return s.getViper(envNames...).Unmarshal(rawVal)
}
func (s *EnvConfig) GetBool(name string, envNames ...string) bool {
	return s.getViper(envNames...).GetBool(name)
}
func (s *EnvConfig) GetObject(name string, envNames ...string) interface{} {
	return s.getViper(envNames...).Get(name)
}
func (s *EnvConfig) SetValue(name string, value interface{}, envNames ...string) *EnvConfig {
	s.getViper(envNames...).Set(name, value)
	return s
}
func (s *EnvConfig) getViper(envNames ...string) *viper.Viper {
	if len(envNames) > 0 {
		return getConfigViper(envNames[0])
	} else {
		return getConfigViper(AppConfigName)
	}
}
func getConfigViper(name string) *viper.Viper {
	if name == "" {
		name = AppConfigName
	}
	name = strings.ToLower(name)
	if v, ok := _ENVMaps.Get(name); ok && v != nil {
		if vv, ok := v.(*viper.Viper); ok && vv != nil {
			return vv
		}
	}

	envType := "yaml"
	envPath := "env"

	envFile := JoinCurrentPath(path.Join(envPath, name+"."+envType))
	if !PathExists(envFile) {
		paths := []string{path.Join(envPath, name+".prod."+envType), path.Join(envPath, name+".dev."+envType)}
		for _, p := range paths {
			//不存在时，自动由dev创建
			if !PathExists(envFile) {
				devFile := JoinCurrentPath(p)
				if PathExists(devFile) {
					if s, err := os.Open(devFile); err == nil {
						defer s.Close()
						if newEnv, err := os.Create(envFile); err == nil {
							defer newEnv.Close()
							io.Copy(newEnv, s)
							break
						}
					}
				}
			}
		}
	}
	v := viper.New()
	v.SetConfigType(envType)
	v.SetConfigName(name)
	v.AddConfigPath(JoinCurrentPath(envPath))
	v.SetConfigFile(envFile)

	if err := v.ReadInConfig(); err != nil {
		log.ErrorF("Fatal error when reading %s config file:%s", name, err)
	}
	_ENVMaps.Set(name, v)
	return v
}
func newInitConfig() {
	Config = &EnvConfig{}
	vp := getConfigViper(AppConfigName)
	if err := vp.Unmarshal(&Config); err != nil {
		log.ErrorF("Fatal error when reading %s config file:%s", AppConfigName, err)
	}
	if Config.App.Port == "" {
		Config.App.Port = "8080"
	}
	if Config.App.Locale == "" {
		Config.App.Locale = "zh-cn"
	}
	if Config.App.Token == "" {
		Config.App.Token = "01e8f6a984101b20b24e4d172ec741be"
	}
	if Config.App.Storage == "" {
		Config.App.Storage = "./storage"
	}
	if Config.Db.Driver == "" {
		Config.Db.Driver = "mysql"
	}
	if Config.Db.Host == "" {
		Config.Db.Host = "localhost"
	}
	if Config.Db.Port == "" {
		if Config.Db.Driver == "mysql" {
			Config.Db.Port = "3306"
		}
		if Config.Db.Driver == "mssql" {
			Config.Db.Port = "1433"
		}
	}
	if Config.Db.Charset == "" {
		Config.Db.Charset = "utf8mb4"
	}
	if Config.Db.Collation == "" {
		Config.Db.Collation = "utf8mb4_general_ci"
	}
	if Config.Auth.Code == "" {
		Config.Auth.Code = Config.App.Code
	}
	kvs := make(map[string]interface{})
	if err := vp.Unmarshal(&kvs); err != nil {
		log.ErrorF("Fatal error when reading %s config file:%s", AppConfigName, err)
	}
	if len(kvs) > 0 {
		for k, v := range kvs {
			Config.SetValue(k, v)
		}
	}
}
