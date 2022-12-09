package main

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/fr13n8/go-practice/internal/config"
	"github.com/fr13n8/go-practice/internal/domain"
	"github.com/fr13n8/go-practice/pkg/elastic"
	"github.com/fr13n8/go-practice/pkg/redis"
	"github.com/fr13n8/go-practice/pkg/utils"
)

func main() {
	run()
}

func run() {
	configPath := utils.GetConfigPath(os.Getenv("config"))
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	es, err := elastic.NewElasticClient(&cfg.Elastic)
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
		event := msg.Channel[5:]
		log.Println("Event:", event)
		b := bytes.NewReader([]byte(msg.Payload))
		if err := gob.NewDecoder(b).Decode(&task); err != nil {
			log.Println(err)
			continue
		}

		switch event {
		case "created":
			log.Println("Created task", task)
			req := esapi.IndexRequest{
				Index:      "tasks",
				DocumentID: task.ID,
				Body:       bytes.NewReader([]byte(task.Name)),
				Refresh:    "true",
			}
			res, err := req.Do(context.Background(), es)
			if err != nil {
				log.Println(err)
				continue
			}
			defer res.Body.Close()
			if res.IsError() {
				log.Printf("[%s] Error indexing document ID=%s", res.Status(), task.ID)
			} else {
				// Deserialize the response into a map.
				var r map[string]interface{}
				if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
					log.Printf("Error parsing the response body: %s", err)
				} else {
					// Print the response status and indexed document version.
					log.Printf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
				}
			}

			log.Println("Indexed task", task)
		case "updated":
			log.Println("Updated task", task)
		case "deleted":
			log.Println("Deleted task", task)
		}
	}
}
