package elastic

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
)

func NewElasticClient() {
	certPath, err := filepath.Abs("pkg/elastic/certs/http_ca.crt")
	if err != nil {
		log.Fatal(err)
	}
	cert, err := os.ReadFile(certPath)
	if err != nil {
		log.Fatal(err)
	}

	es, err := elasticsearch.NewClient(
		elasticsearch.Config{
			Addresses: []string{"https://localhost:9200"},
			Username:  "elastic",
			Password:  "4XTdk*2iyN*n2jYYqgg1",
			// ServiceToken: "eyJ2ZXIiOiI4LjMuMyIsImFkciI6WyIxNzIuMjAuMC4yOjkyMDAiXSwiZmdyIjoiNzU0MGY5NDdmMzVkNjQxMWFjOTJhNzM3ZTFhMzJhYzBmZDJhMDc2YjNmMjM4ODQ3M2RhNjI1M2JlOTFmMzFjYSIsImtleSI6InljajZzSUlCR0hiVWkzSTdWR1VDOmdaSGE1QllhUXBxUDJJYkpnTHJTWVEifQ==",
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   10,
				ResponseHeaderTimeout: time.Second,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
			CACert: cert,
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
	log.Println(res)
}
