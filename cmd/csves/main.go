package main

import (
	"log"

	"csves/pkg/config"
	"csves/pkg/csv"
	"csves/pkg/elasticsearch"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	// Initialize CSV service
	csvService, err := csv.NewService(cfg)
	if err != nil {
		log.Fatalf("Error initializing CSV service: %s", err)
	}

	// Process header
	headerMap, err := csvService.ProcessHeader()
	if err != nil {
		log.Fatalf("Error processing CSV header: %s", err)
	}

	// Process records
	documents, err := csvService.ProcessRecords(headerMap)
	if err != nil {
		log.Fatalf("Error processing CSV records: %s", err)
	}

	if cfg.TestMode {
		csvService.PrintDocuments(documents, true)
		return
	}

	// Initialize Elasticsearch service
	esService, err := elasticsearch.NewService(cfg)
	if err != nil {
		log.Fatalf("Error initializing Elasticsearch service: %s", err)
	}

	// Setup Elasticsearch
	if err := esService.Setup(); err != nil {
		log.Fatalf("Error setting up Elasticsearch: %s", err)
	}

	// Print sample for verification
	csvService.PrintDocuments(documents, false)

	// Bulk index documents
	if err := esService.BulkIndex(documents); err != nil {
		log.Fatalf("Error bulk indexing documents: %s", err)
	}

	log.Println("All documents indexed successfully")
}
