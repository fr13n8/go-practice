package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"

	"github.com/fr13n8/go-practice/internal/config"
	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/fr13n8/go-practice/pkg/redis"
	"github.com/fr13n8/go-practice/pkg/utils"
)

func main() {
	run()
}

func run() {
	// es, err := elasticsearch.NewDefaultClient()
	// if err != nil {
	// 	log.Fatalf("Error creating the client: %s", err)
	// }
	// log.Println(elasticsearch.Version)
	// res, err := es.Info()
	// if err != nil {
	// 	log.Fatalf("Error getting response: %s", err)
	// }
	// defer res.Body.Close()
	// log.Println(res)
	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	rdb, err := redis.NewRedis(&cfg.Redis)
	if err != nil {
		panic(err)
	}
	defer rdb.Close()

	pubsub := rdb.PSubscribe("task.*")
	_, err = pubsub.Receive()
	if err != nil {
		panic(err)
	}
	defer pubsub.Close()
	ch := pubsub.Channel()

	for msg := range ch {
		var task domain.Task
		b := bytes.NewReader([]byte(msg.Payload))
		if err := gob.NewDecoder(b).Decode(&task); err != nil {
			log.Println(err)
			continue
		}

		log.Println("Received message", task)
	}
}
