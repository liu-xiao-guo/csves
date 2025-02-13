package models

// FieldConfig represents a field configuration
type FieldConfig struct {
	Name     string `json:"name"`     // Field name
	Required bool   `json:"required"` // Whether the field is required
	CSVName  string `json:"csv_name"` // Original CSV header name
}

// Document represents a dynamic document with arbitrary fields
type Document struct {
	Fields map[string]string
}

// NewDocument creates a new document with initialized fields
func NewDocument() *Document {
	return &Document{
		Fields: make(map[string]string),
	}
}

// SetField sets a field value
func (d *Document) SetField(name, value string) {
	d.Fields[name] = value
}

// GetField gets a field value
func (d *Document) GetField(name string) string {
	return d.Fields[name]
}
