package elastic

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/fr13n8/go-practice/internal/config"
)

func NewElasticClient(cfg *config.ElasticConfig) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{
				"http://" + cfg.Host + ":" + cfg.Port,
			},
		},
	)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Fatalf("Error: %s", res.String())
	}

	return es, nil
}
