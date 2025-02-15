# CSVES (CSV to Elasticsearch)

A flexible tool for importing CSV data into Elasticsearch with automatic field detection and mapping.

## Warning

This tool (probably) won't work for all csv files out of the box. 

Running test mode first is recommended.

## Features

- 🔍 Automatic CSV delimiter detection
- 📄 Dynamic field mapping
- 🧹 Automatic whitespace and control character cleaning
- 🎯 Field selection and filtering
- ⚙️ Configurable through command-line flags or environment variables
- 🧪 Test mode for data verification
- 📝 Custom field mapping through JSON configuration

## Installation

### Prerequisites

- Go 1.23 or higher
- Elasticsearch 8.x
- Access to an Elasticsearch instance

### Build Steps

1. Clone the repository:
```bash
git clone https://github.com/githubesson/csves
cd csves
```

2. Build the binary:
```bash
go build -o csves cmd/csves/main.go
```

## Usage

### Basic Usage

```bash
# Using .env settings
./csves

# Test mode (no Elasticsearch connection) without .env settings
./csves -csv="data.csv" -test

# Select specific fields without .env settings
./csves -csv="data.csv" -select="email,phone,address"
```

### Command Line Flags

| Flag | Description | Default | Required |
|------|-------------|---------|----------|
| `-csv` | Path to CSV file | - | Yes |
| `-es-url` | Elasticsearch URL | http://localhost:9200 | No |
| `-index` | Elasticsearch index name | csv_data | No |
| `-fields` | Path to field configuration file | - | No |
| `-select` | Comma-separated list of fields to include | - | No |
| `-delimiter` | CSV delimiter character | auto-detect | No |
| `-test` | Run in test mode | false | No |

### Environment Variables

You can also configure the tool using environment variables in a `.env` file:

```env
ELASTICSEARCH_URL=https://localhost:9200
INDEX_NAME=csv_test
CSV_FILE_PATH=./example.csv
USER_NAME=elastic
PASSWORD="y9NWnPq0++V=WxMXxSmr"
FIELD_CONFIG_PATH=fields.json
ELASTICSEARCH_CERT_PATH=/Users/liuxg/elastic/elasticsearch-8.17.1/config/certs/http_ca.crt
```

### Field Configuration

Create a `fields.json` file to specify field mappings and requirements:

```json
[
    {
        "name": "User Id",
        "required": true,
        "csv_name": "userid"
    },
    {
        "name": "Email",
        "required": true,
        "csv_name": "email"
    }
]
```

- `name`: Field name in Elasticsearch
- `required`: Whether the field must exist in CSV
- `csv_name`: Column header name in CSV file

## Examples

### 1. Basic Import
```bash
./csves -csv="users.csv"
```

### 2. Custom Elasticsearch Configuration
```bash
./csves -csv="users.csv" -es-url="http://elasticsearch:9200" -index="users_v1"
```

### 3. Field Selection
```bash
./csves -csv="users.csv" -select="email,phone" -test
```

### 4. Custom Field Mapping
```bash
./csves -csv="users.csv" -fields="fields.json"
```

### 5. Specific Delimiter
```bash
./csves -csv="users.csv" -delimiter=";"
```

## Data Cleaning

The tool automatically:
- Removes leading and trailing whitespace
- Removes control characters
- Normalizes internal spaces
- Skips empty fields
- Handles multi-line values

## Error Handling

- Validates required fields
- Reports parsing errors
- Shows bulk indexing failures
- Provides detailed error messages

## Development

### Project Structure
```
csves/
├── cmd/
│   └── csves/
│       └── main.go           # Entry point
├── pkg/
│   ├── config/              # Configuration handling
│   ├── csv/                 # CSV processing
│   ├── elasticsearch/       # ES operations
│   └── models/              # Data models
├── go.mod                   # Go modules file
├── go.sum                   # Dependencies checksum
└── README.md               # This file
```

## License

This project is licensed under the MIT License - see the LICENSE file for details. 