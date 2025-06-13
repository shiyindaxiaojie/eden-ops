package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

// Config 系统配置结构
type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Log       LogConfig       `mapstructure:"log"`
	Tencent   TencentConfig   `mapstructure:"tencent"`
	IPLocator IPLocatorConfig `mapstructure:"iplocator"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	GinMode string `mapstructure:"gin_mode"` // gin运行模式: debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string        `mapstructure:"host"`
	Port            int           `mapstructure:"port"`
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	DBName          string        `mapstructure:"dbname"`
	Charset         string        `mapstructure:"charset"`
	MaxIdleConns    int           `mapstructure:"maxIdleConns"`
	MaxOpenConns    int           `mapstructure:"maxOpenConns"`
	ConnMaxLifetime time.Duration `mapstructure:"connMaxLifetime"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire"`
	Issuer string `mapstructure:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level       string `mapstructure:"level"`
	Format      string `mapstructure:"format"`
	Output      string `mapstructure:"output"`
	File        string `mapstructure:"file"`
	SQLEnabled  bool   `mapstructure:"sql_enabled"`  // 是否启用SQL日志
	SQLDetailed bool   `mapstructure:"sql_detailed"` // 是否显示详细SQL日志
}

// TencentConfig 腾讯云配置
type TencentConfig struct {
	SecretID  string `mapstructure:"secret_id" env:"TENCENT_SECRET_ID"`
	SecretKey string `mapstructure:"secret_key" env:"TENCENT_SECRET_KEY"`
	Region    string `mapstructure:"region" env:"TENCENT_REGION"`
	K8sConfig string `mapstructure:"k8s_config"` // Kubernetes配置文件路径
}

// IPLocatorConfig IP定位配置
type IPLocatorConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	BaseURL   string `mapstructure:"base_url"`
}

// LoadFromEnv 从环境变量加载配置
func (c *TencentConfig) LoadFromEnv() {
	if id, ok := os.LookupEnv("TENCENT_SECRET_ID"); ok {
		c.SecretID = id
	}
	if key, ok := os.LookupEnv("TENCENT_SECRET_KEY"); ok {
		c.SecretKey = key
	}
	if region, ok := os.LookupEnv("TENCENT_REGION"); ok {
		c.Region = region
	}
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	log.Printf("配置加载成功: %s", configPath)
	return &config, nil
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.DBName, c.Charset)
}

// LoadConfigFromYAML 从文件加载配置
func LoadConfigFromYAML(filename string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// 解析配置
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() (*Config, error) {
	cfg := &Config{}

	// 服务器配置
	cfg.Server.Host = os.Getenv("SERVER_HOST")
	port, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	cfg.Server.Port = port
	cfg.Server.Mode = os.Getenv("SERVER_MODE")

	// 数据库配置
	cfg.Database.Host = os.Getenv("DB_HOST")
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	cfg.Database.Port = dbPort
	cfg.Database.Username = os.Getenv("DB_USERNAME")
	cfg.Database.Password = os.Getenv("DB_PASSWORD")
	cfg.Database.DBName = os.Getenv("DB_NAME")
	cfg.Database.Charset = os.Getenv("DB_CHARSET")
	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	cfg.Database.MaxIdleConns = maxIdleConns
	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	cfg.Database.MaxOpenConns = maxOpenConns
	connMaxLifetime, _ := time.ParseDuration(os.Getenv("DB_CONN_MAX_LIFETIME"))
	cfg.Database.ConnMaxLifetime = connMaxLifetime

	// JWT配置
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	expire, _ := strconv.Atoi(os.Getenv("JWT_EXPIRE"))
	cfg.JWT.Expire = expire
	cfg.JWT.Issuer = "eden-ops"

	// 日志配置
	cfg.Log.Level = os.Getenv("LOG_LEVEL")
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info" // 默认使用info级别
	}
	cfg.Log.Format = os.Getenv("LOG_FORMAT")
	if cfg.Log.Format == "" {
		cfg.Log.Format = "text" // 默认使用text格式
	}
	cfg.Log.Output = os.Getenv("LOG_OUTPUT")
	if cfg.Log.Output == "" {
		cfg.Log.Output = "console" // 默认输出到控制台
	}
	cfg.Log.File = os.Getenv("LOG_FILE")
	if cfg.Log.File == "" {
		cfg.Log.File = "logs/eden-ops.log" // 默认日志文件路径
	}

	return cfg, nil
}
