package csv

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"csves/pkg/config"
	"csves/pkg/models"
)

// cleanString removes leading/trailing spaces and control characters
func cleanString(s string) string {
	// First trim spaces and control characters
	s = strings.TrimFunc(s, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsControl(r)
	})

	// Then normalize internal spaces (replace multiple spaces with single space)
	return strings.Join(strings.Fields(s), " ")
}

// Service handles CSV operations
type Service struct {
	config *config.Config
	reader *csv.Reader
	file   *os.File
}

// DetectDelimiter tries to detect the CSV delimiter by checking common ones
func DetectDelimiter(filePath string) (rune, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return ',', fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Read first line to analyze
	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return ',', fmt.Errorf("empty file")
	}
	firstLine := scanner.Text()

	// Common delimiters to check
	delimiters := []rune{',', ';', '\t', '|'}
	maxCount := 0
	bestDelimiter := ','

	for _, d := range delimiters {
		count := strings.Count(firstLine, string(d))
		if count > maxCount {
			maxCount = count
			bestDelimiter = d
		}
	}

	return bestDelimiter, nil
}

// DetectFields detects available fields from CSV header
func DetectFields(header []string) []models.FieldConfig {
	var fields []models.FieldConfig
	for _, h := range header {
		name := cleanString(h)
		if name == "" {
			continue // Skip empty field names
		}
		fields = append(fields, models.FieldConfig{
			Name:     name,
			CSVName:  name,
			Required: false,
		})
	}
	return fields
}

// NewService creates a new CSV service
func NewService(cfg *config.Config) (*Service, error) {
	file, err := os.Open(cfg.CSVFilePath)
	if err != nil {
		return nil, fmt.Errorf("error opening CSV file: %w", err)
	}

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true // Enable built-in CSV trimming

	if cfg.DelimiterChar == 0 {
		// Auto-detect delimiter if not specified
		delimiter, err := DetectDelimiter(cfg.CSVFilePath)
		if err != nil {
			return nil, fmt.Errorf("error detecting delimiter: %w", err)
		}
		reader.Comma = delimiter
		fmt.Printf("Detected delimiter: '%c'\n", delimiter)
	} else {
		reader.Comma = cfg.DelimiterChar
	}

	return &Service{
		config: cfg,
		reader: reader,
		file:   file,
	}, nil
}

// Close closes the CSV file
func (s *Service) Close() error {
	return s.file.Close()
}

// ProcessHeader processes the CSV header and maps column names to indices
func (s *Service) ProcessHeader() (map[string]int, error) {
	header, err := s.reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}

	// Clean header fields
	for i, h := range header {
		header[i] = cleanString(h)
	}

	// Detect fields if not provided
	if len(s.config.Fields) == 0 {
		s.config.Fields = DetectFields(header)
		fmt.Println("Detected fields:", s.config.Fields)
	}

	headerMap := make(map[string]int)
	for i, column := range header {
		if column != "" { // Skip empty column names
			headerMap[strings.ToLower(column)] = i
		}
	}

	// Validate required fields if any
	for _, field := range s.config.Fields {
		if field.Required {
			if _, exists := headerMap[strings.ToLower(field.CSVName)]; !exists {
				return nil, fmt.Errorf("required field '%s' not found in CSV header", field.CSVName)
			}
		}
	}

	fmt.Println("CSV Header mapping:", headerMap)
	return headerMap, nil
}

// ProcessRecords processes CSV records and returns documents
func (s *Service) ProcessRecords(headerMap map[string]int) ([]models.Document, error) {
	var documents []models.Document

	// Get the base filename without path
	sourceFile := filepath.Base(s.config.CSVFilePath)

	for {
		record, err := s.reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Warning: Error reading CSV record: %s", err)
			continue
		}

		// Clean all record values
		for i, v := range record {
			record[i] = cleanString(v)
		}

		doc := models.NewDocument()
		// Set fields based on configuration
		for _, field := range s.config.Fields {
			if idx, exists := headerMap[strings.ToLower(field.CSVName)]; exists {
				value := record[idx]
				if value != "" { // Only set non-empty values
					doc.SetField(field.Name, value)
				}
			}
		}

		// Add source CSV filename
		doc.SetField("source_csv", sourceFile)

		documents = append(documents, *doc)
	}

	return documents, nil
}

// PrintDocuments prints documents in a readable format
func (s *Service) PrintDocuments(documents []models.Document, printAll bool) {
	if printAll {
		fmt.Println("Test Mode - Printing all processed records:")
		for i, doc := range documents {
			fmt.Printf("Record %d:\n", i+1)
			for _, field := range s.config.Fields {
				value := doc.GetField(field.Name)
				if value != "" { // Only print non-empty values
					fmt.Printf("  %s: %s\n", field.Name, value)
				}
			}
			// Always print source CSV
			fmt.Printf("  source_csv: %s\n", doc.GetField("source_csv"))
			fmt.Println()
		}
	} else {
		fmt.Println("Sample of processed records:")
		for i := 0; i < 2 && i < len(documents); i++ {
			fmt.Printf("Fields: %+v\n", documents[i].Fields)
		}
	}
	fmt.Printf("Total records processed: %d\n", len(documents))
}
