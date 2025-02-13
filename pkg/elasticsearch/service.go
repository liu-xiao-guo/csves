package elasticsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"csves/pkg/config"
	"csves/pkg/models"

	"github.com/elastic/go-elasticsearch/v8"
)

// Service handles Elasticsearch operations
type Service struct {
	client *elasticsearch.Client
	config *config.Config
}

// NewService creates a new Elasticsearch service
func NewService(cfg *config.Config) (*Service, error) {
	esCfg := elasticsearch.Config{
		Addresses: []string{cfg.ElasticsearchURL},
	}

	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, fmt.Errorf("error creating Elasticsearch client: %w", err)
	}

	return &Service{
		client: client,
		config: cfg,
	}, nil
}

// Setup ensures the index exists with proper mapping
func (s *Service) Setup() error {
	// Check if index exists
	res, err := s.client.Indices.Exists([]string{s.config.IndexName})
	if err != nil {
		return fmt.Errorf("error checking index existence: %w", err)
	}

	if res.StatusCode == 404 {
		// Create index with basic mapping for source_csv
		mapping := `{
			"mappings": {
				"properties": {
					"source_csv": {
						"type": "keyword"
					}
				},
				"dynamic": true
			}
		}`

		res, err = s.client.Indices.Create(
			s.config.IndexName,
			s.client.Indices.Create.WithBody(strings.NewReader(mapping)),
		)
		if err != nil {
			return fmt.Errorf("error creating index: %w", err)
		}
		if res.IsError() {
			return fmt.Errorf("error creating index: %s", res.String())
		}
	}

	return nil
}

// BulkIndex indexes multiple documents in bulk
func (s *Service) BulkIndex(documents []models.Document) error {
	var buf bytes.Buffer

	for _, doc := range documents {
		meta := map[string]interface{}{
			"index": map[string]interface{}{
				"_index": s.config.IndexName,
			},
		}

		if err := json.NewEncoder(&buf).Encode(meta); err != nil {
			return fmt.Errorf("error encoding metadata: %w", err)
		}

		if err := json.NewEncoder(&buf).Encode(doc.Fields); err != nil {
			return fmt.Errorf("error encoding document: %w", err)
		}
	}

	res, err := s.client.Bulk(bytes.NewReader(buf.Bytes()))
	if err != nil {
		return fmt.Errorf("error bulk indexing: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("bulk indexing failed: %s", res.String())
	}

	var bulkResponse map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&bulkResponse); err != nil {
		return fmt.Errorf("error parsing bulk response: %w", err)
	}

	if bulkResponse["errors"].(bool) {
		fmt.Println("Some documents failed to index")
		items := bulkResponse["items"].([]interface{})
		for _, item := range items {
			indexResp := item.(map[string]interface{})["index"].(map[string]interface{})
			if indexResp["error"] != nil {
				fmt.Printf("Error: %v\n", indexResp["error"])
			}
		}
		return fmt.Errorf("some documents failed to index")
	}

	return nil
}
