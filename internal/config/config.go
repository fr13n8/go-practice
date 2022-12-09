package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP     HttpConfig
	GRPC     GrpcConfig
	Elastic  ElasticConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Jaeger   JaegerConfig
}

type HttpConfig struct {
	ServerName         string
	Host               string
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxRequestBodySize int
}

type ElasticConfig struct {
	Host     string
	Port     string
	Password string
	User     string
}

type GrpcConfig struct {
	ServerName        string
	Host              string
	Port              string
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Password string
	User     string
	Name     string
}

type JaegerConfig struct {
	Host        string
	Port        string
	ServiceName string
	LogSpans    bool
}

func NewConfig(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
