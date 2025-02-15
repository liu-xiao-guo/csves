package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"csves/pkg/models"

	"github.com/joho/godotenv"
)

// Config holds all configuration parameters
type Config struct {
	ElasticsearchURL string
	IndexName        string
	CSVFilePath      string
	DelimiterChar    rune
	HeaderMap        map[string]int
	TestMode         bool
	Fields           []models.FieldConfig
	FieldConfigPath  string
	SelectedFields   []string
	UserName		 string
	Password		 string
	CertPath		 string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		HeaderMap: make(map[string]int),
	}

	flag.StringVar(&config.ElasticsearchURL, "es-url", getEnvOrDefault("ELASTICSEARCH_URL", "http://localhost:9200"), "Elasticsearch URL")
	flag.StringVar(&config.IndexName, "index", getEnvOrDefault("INDEX_NAME", "csv_data"), "Elasticsearch index name")
	flag.StringVar(&config.CSVFilePath, "csv", getEnvOrDefault("CSV_FILE_PATH", ""), "Path to CSV file")
	flag.StringVar(&config.FieldConfigPath, "fields", getEnvOrDefault("FIELD_CONFIG_PATH", ""), "Path to field configuration JSON file")
	selectedFields := flag.String("select", "", "Comma-separated list of fields to include (empty for all fields)")
	delimiter := flag.String("delimiter", "", "CSV delimiter character (auto-detect if not specified)")
	flag.BoolVar(&config.TestMode, "test", false, "Test mode - only parse and print documents without connecting to Elasticsearch")
	flag.StringVar(&config.UserName, "username", getEnvOrDefault("USER_NAME", "elastic"), "User name")
	flag.StringVar(&config.Password, "password", getEnvOrDefault("PASSWORD", "123456"), "Password")
	flag.StringVar(&config.CertPath, "certpath", getEnvOrDefault("ELASTICSEARCH_CERT_PATH", "./http_ca.crt"), "Password")

	fmt.Println("es-url:", config.ElasticsearchURL)
	fmt.Println("index:", config.IndexName)
	fmt.Println("csv:", config.CSVFilePath)
	fmt.Println("fields:", config.FieldConfigPath)
	fmt.Println("fields:", config.FieldConfigPath)
	fmt.Println("test:", config.TestMode)
	fmt.Println("username:", config.UserName)
	fmt.Println("password:", config.Password)
	fmt.Println("certpath:", config.CertPath)

	flag.Parse()

	if config.CSVFilePath == "" {
		return nil, fmt.Errorf("CSV file path is required")
	}

	if *selectedFields != "" {
		config.SelectedFields = strings.Split(*selectedFields, ",")
		for i, field := range config.SelectedFields {
			config.SelectedFields[i] = strings.TrimSpace(field)
		}
	}

	if *delimiter != "" {
		config.DelimiterChar = rune((*delimiter)[0])
	}

	if config.FieldConfigPath != "" {
		fields, err := loadFieldConfig(config.FieldConfigPath)
		if err != nil {
			return nil, fmt.Errorf("error loading field configuration: %w", err)
		}
		config.Fields = fields

		if len(config.SelectedFields) > 0 {
			config.Fields = filterFields(config.Fields, config.SelectedFields)
		}
	}

	return config, nil
}

func filterFields(fields []models.FieldConfig, selectedFields []string) []models.FieldConfig {
	if len(selectedFields) == 0 {
		return fields
	}

	var filtered []models.FieldConfig
	for _, field := range fields {
		for _, selected := range selectedFields {
			if strings.EqualFold(field.Name, selected) {
				filtered = append(filtered, field)
				break
			}
		}
	}
	return filtered
}

func loadFieldConfig(path string) ([]models.FieldConfig, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading field config file: %w", err)
	}

	var fields []models.FieldConfig
	if err := json.Unmarshal(file, &fields); err != nil {
		return nil, fmt.Errorf("error parsing field config: %w", err)
	}

	return fields, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
