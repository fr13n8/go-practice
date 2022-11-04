package config

import "time"

type Config struct {
	HTTP          HttpConfig
	ElasticConfig ElasticConfig
	Database      string
}

type HttpConfig struct {
	ServerName         string
	Host               string
	Port               string
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	MaxRequestBodySize int
}

func NewConfig() *Config {
	return &Config{
		HTTP: HttpConfig{
			ServerName:         "My Server",
			Host:               "",
			Port:               "3000",
			ReadTimeout:        time.Second * 5,
			WriteTimeout:       time.Second * 5,
			MaxRequestBodySize: 1048576,
		},
		ElasticConfig: ElasticConfig{
			Host:     "localhost",
			Port:     "9200",
			User:     "elastic",
			Password: "gvETj_R1HXb8hp0blG39",
		},
		Database: "database.db",
	}
}

type ElasticConfig struct {
	Host     string
	Port     string
	Password string
	User     string
}
