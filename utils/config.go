package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/nbkit/mdf/gmap"
	"github.com/nbkit/mdf/log"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type AppConfig struct {
	Token   string `mapstructure:"token" json:"token"`
	Code    string `mapstructure:"code" json:"code"`
	Active  string `mapstructure:"active" json:"active"`
	Name    string `mapstructure:"name" json:"name"`
	Port    string `mapstructure:"port" json:"port"`
	Locale  string `mapstructure:"locale" json:"locale"`
	Mode    string `mapstructure:"mode" json:"mode"`
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
	Mode      string `mapstructure:"mode" json:"mode"`
}

func (s *DbConfig) fill() {
	if s.Driver == "" {
		s.Driver = ORM_DRIVER_MYSQL
	}
	if s.Host == "" {
		s.Host = "localhost"
	}
	if s.Port == "" {
		if s.Driver == ORM_DRIVER_MYSQL {
			s.Port = "3306"
		}
		if s.Driver == ORM_DRIVER_MSSQL {
			s.Port = "1433"
		}
	}
	if s.Charset == "" {
		s.Charset = "utf8mb4"
	}
	if s.Collation == "" {
		s.Collation = "utf8mb4_general_ci"
	}
}

func (s DbConfig) GetDsnString(useDB bool) string {
	s.fill()
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//mssql:   =>  sqlserver://username:password@localhost:1433?database=dbname
	//mysql => user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	str := ""
	// 创建连接
	if s.Driver == ORM_DRIVER_MSSQL {
		var buf bytes.Buffer
		buf.WriteString("sqlserver://")
		buf.WriteString(s.Username)
		if s.Password != "" {
			buf.WriteByte(':')
			buf.WriteString(s.Password)
		}
		buf.WriteByte('@')
		if s.Host != "" {
			buf.WriteString(s.Host)
			if s.Port != "" {
				buf.WriteByte(':')
				buf.WriteString(s.Port)
			} else {
				buf.WriteString(":1433")
			}
		}
		if s.Database != "" && useDB {
			buf.WriteString("?database=")
			buf.WriteString(s.Database)
		} else {
			buf.WriteString("?database=master")
		}
		str = buf.String()
		return str
	}
	{
		config := mysql.Config{
			User:                 s.Username,
			Passwd:               s.Password,
			Net:                  "tcp",
			Addr:                 s.Host,
			AllowNativePasswords: true,
			ParseTime:            true,
			Loc:                  time.Local,
		}
		if useDB {
			config.DBName = s.Database
		}
		if s.Port != "" {
			config.Addr = fmt.Sprintf("%s:%s", s.Host, s.Port)
		}
		str = config.FormatDSN()
	}
	return str
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
	Auth AuthConfig             `mapstructure:"auth" json:"auth"`
	Data map[string]interface{} `json:"data"`
}

type EnvConfig struct {
	App  AppConfig  `mapstructure:"app" json:"app"`
	Db   DbConfig   `mapstructure:"db" json:"db"`
	Auth AuthConfig `mapstructure:"auth" json:"auth"`
	data map[string]interface{}
}

func (c EnvConfig) MarshalJSON() ([]byte, error) {
	jsonMap := jsonConfig{}
	jsonMap.App = c.App
	jsonMap.Db = c.Db
	jsonMap.Auth = c.Auth
	jsonMap.Data = c.data
	return json.Marshal(jsonMap)
}

func (c *EnvConfig) UnmarshalJSON(b []byte) error {
	jsonMap := jsonConfig{}
	json.Unmarshal(b, &jsonMap)
	c.App = jsonMap.App
	c.Db = jsonMap.Db
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

func (s *EnvConfig) decrypt(value string) string {
	if value != "" && strings.Index(value, "ECN(") == 0 {
		value = value[4 : len(value)-1]
		value, _ = AesCFBDecrypt(value, Config.App.Token)
	}
	return value
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
	//v.SetConfigType(envType)
	v.SetConfigName(name)
	v.AddConfigPath(JoinCurrentPath(envPath))
	//v.SetConfigFile(envFile)

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
	if Config.App.Mode == "" {
		Config.App.Mode = "release"
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

	Config.Db.Password = Config.decrypt(Config.Db.Password)
	Config.Db.Host = Config.decrypt(Config.Db.Host)
	Config.Db.Username = Config.decrypt(Config.Db.Username)
	Config.Db.Database = Config.decrypt(Config.Db.Database)

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
