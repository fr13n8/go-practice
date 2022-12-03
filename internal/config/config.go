package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	HTTP          HttpConfig
	GRPC          GrpcConfig
	ElasticConfig ElasticConfig
	Database      DatabaseConfig
	Redis         RedisConfig
	Jaeger        JaegerConfig
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
	// return &Config{
	// 	HTTP: HttpConfig{
	// 		ServerName:         "Http Server",
	// 		Host:               "0.0.0.0",
	// 		Port:               "80",
	// 		ReadTimeout:        time.Second * 5,
	// 		WriteTimeout:       time.Second * 5,
	// 		MaxRequestBodySize: 1048576,
	// 	},
	// 	GRPC: GrpcConfig{
	// 		ServerName: "Grpc Server",
	// 		Host:       "0.0.0.0",
	// 		Port:       "8080",
	// 	},
	// 	ElasticConfig: ElasticConfig{
	// 		Host:     "localhost",
	// 		Port:     "9200",
	// 		User:     "elastic",
	// 		Password: "gvETj_R1HXb8hp0blG39",
	// 	},
	// 	Database: DatabaseConfig{
	// 		Host:     "postgres",
	// 		Port:     "5432",
	// 		User:     "postgres",
	// 		Password: "postgres",
	// 		Name:     "go-practice",
	// 	},
	// 	Redis: RedisConfig{
	// 		Host:     "redis",
	// 		Port:     "6379",
	// 		Password: "",
	// 	},
	// }
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
