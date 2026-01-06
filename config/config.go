package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/gogf/gf/util/gconv"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

const (
	ServerName     = "app"
	ServerFullName = "app"
)

var (
	etcdKey         = fmt.Sprintf("/configs/%s/system", ServerFullName)
	etcdAddr        string
	localConfigPath string
	GlobalConfig    Config
)

func init() {
	flag.StringVar(&localConfigPath, "c", ServerName+"_local.yml", "default config path")
	flag.StringVar(&etcdAddr, "r", os.Getenv("ETCD_ADDR"), "default consul address")
}

type Config struct {
	Server   Server   `yaml:"server"`
	Mysql    Mysql    `yaml:"mysql"`
	Postgres Postgres `yaml:"postgres"`
	Redis    Redis    `yaml:"redis"`
}

type Server struct {
	HttpPort        int    `yaml:"http_port"`
	Env             string `yaml:"env"`
	EnablePprof     bool   `yaml:"enable_pprof"`
	LogLevel        string `yaml:"log_level"`
	ShutdownTimeout int    `yaml:"shutdown_timeout"` // 优雅关闭超时时间(秒)
}

type Mysql struct {
	Dialect  string `yaml:"dialect"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Charset  string `yaml:"charset"`
	ShowSql  bool   `yaml:"show_sql"`
	MaxOpen  int    `yaml:"max_open"`
	MaxIdle  int    `yaml:"max_idle"`
}

type Postgres struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	User          string `yaml:"user"`
	Password      string `yaml:"password"`
	Database      string `yaml:"database"`
	SSLMode       string `yaml:"ssl_mode"`
	MaxOpen       int    `yaml:"max_open"`
	MaxIdle       int    `yaml:"max_idle"`
	EnablePostGIS bool   `yaml:"enable_postgis"`
}

func (m *Mysql) GetDsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=Local",
		m.User, m.Password, m.Host, m.Port, m.Database, m.Charset)
}

func (p *Postgres) GetDsn() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Database, p.SSLMode)
}

type Redis struct {
	Addr    string `yaml:"addr"`
	PWD     string `yaml:"password"`
	DBIndex int    `yaml:"db_index"`
	MaxIdle int    `yaml:"max_idle"`
	MaxOpen int    `yaml:"max_open"`
}

func InitConfig() *Config {
	var (
		err      error
		tempConf = &Config{}
		vipConf  = viper.New()
	)

	flag.Parse()

	// etcd地址存在，优先使用etcd的配置
	if etcdAddr != "" {
		tempConf, err = getFromRemoteAndWatchUpdate(vipConf)
		if err != nil {
			panic(err)
		}
		applyEnvOverrides(tempConf)
		return tempConf
	}

	// 从本地获取
	tempConf, err = getFromLocal()
	if err != nil {
		panic(err)
	}

	// 应用环境变量覆盖
	applyEnvOverrides(tempConf)
	return tempConf
}

// applyEnvOverrides 环境变量覆盖配置
func applyEnvOverrides(conf *Config) {
	// Server
	if v := os.Getenv("APP_HTTP_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			conf.Server.HttpPort = port
		}
	}
	if v := os.Getenv("APP_ENV"); v != "" {
		conf.Server.Env = v
	}
	if v := os.Getenv("APP_LOG_LEVEL"); v != "" {
		conf.Server.LogLevel = v
	}

	// MySQL
	if v := os.Getenv("MYSQL_HOST"); v != "" {
		conf.Mysql.Host = v
	}
	if v := os.Getenv("MYSQL_PORT"); v != "" {
		if port, err := strconv.Atoi(v); err == nil {
			conf.Mysql.Port = port
		}
	}
	if v := os.Getenv("MYSQL_USER"); v != "" {
		conf.Mysql.User = v
	}
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" {
		conf.Mysql.Password = v
	}
	if v := os.Getenv("MYSQL_DATABASE"); v != "" {
		conf.Mysql.Database = v
	}

	// Redis
	if v := os.Getenv("REDIS_ADDR"); v != "" {
		conf.Redis.Addr = v
	}
	if v := os.Getenv("REDIS_PASSWORD"); v != "" {
		conf.Redis.PWD = v
	}

	// 设置默认值
	if conf.Server.ShutdownTimeout == 0 {
		conf.Server.ShutdownTimeout = 10
	}
}

func getFromRemoteAndWatchUpdate(v *viper.Viper) (*Config, error) {
	tempConf := Config{}
	if err := v.AddRemoteProvider("etcd3", etcdAddr, etcdKey); err != nil {
		return nil, err
	}
	if err := v.ReadRemoteConfig(); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(&tempConf); err != nil {
		return nil, err
	}

	go func() {
		for {
			time.Sleep(time.Minute * 1)
			if err := v.WatchRemoteConfig(); err == nil {
				_ = v.Unmarshal(&GlobalConfig)
				fmt.Println(">>> etcd config hot-reloaded: ", gconv.String(GlobalConfig))
			}
		}
	}()
	return &tempConf, nil
}

func getFromLocal() (*Config, error) {
	tempConf := Config{}
	if _, err := os.Stat(localConfigPath); err == nil {
		content, err := os.ReadFile(localConfigPath)
		if err != nil {
			return nil, err
		}
		err = yaml.Unmarshal(content, &tempConf)
		if err != nil {
			return nil, err
		}
		return &tempConf, nil
	}
	return nil, fmt.Errorf("local config file not found, file_name: %s", localConfigPath)
}
