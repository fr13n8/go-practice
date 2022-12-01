package config

import "time"

type Config struct {
	HTTP          HttpConfig
	GRPC          GrpcConfig
	ElasticConfig ElasticConfig
	Database      DatabaseConfig
	Redis         RedisConfig
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
	ServerName string
	Host       string
	Port       string
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

func NewConfig() *Config {
	return &Config{
		HTTP: HttpConfig{
			ServerName:         "Http Server",
			Host:               "0.0.0.0",
			Port:               "80",
			ReadTimeout:        time.Second * 5,
			WriteTimeout:       time.Second * 5,
			MaxRequestBodySize: 1048576,
		},
		GRPC: GrpcConfig{
			ServerName: "Grpc Server",
			Host:       "0.0.0.0",
			Port:       "8080",
		},
		ElasticConfig: ElasticConfig{
			Host:     "localhost",
			Port:     "9200",
			User:     "elastic",
			Password: "gvETj_R1HXb8hp0blG39",
		},
		Database: DatabaseConfig{
			Host:     "postgres",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			Name:     "go-practice",
		},
		Redis: RedisConfig{
			Host:     "redis",
			Port:     "6379",
			Password: "",
		},
	}
}
