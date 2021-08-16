package manager

import "github.com/elastic/go-elasticsearch/v8"

type ElasticManager struct {
	client *elasticsearch.Client
}

func NewElasticManager(addr string) (*ElasticManager, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{addr},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return &ElasticManager{es}, nil
}
