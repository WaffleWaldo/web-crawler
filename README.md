# Web Crawler

A concurrent web crawler written in Go with MongoDB integration for storing crawled pages.

## Features

- Concurrent crawling with configurable number of workers
- Rate limiting to avoid overwhelming target servers
- MongoDB integration for storing crawled pages
- HTML content parsing and link extraction
- URL normalization and deduplication
- Graceful shutdown handling

## Project Structure

```
.
├── cmd/
│   └── crawler/         # Main application entry point
├── internal/
│   ├── crawler/         # Core crawler implementation
│   ├── queue/           # URL queue management
│   └── storage/         # Storage interfaces and implementations
├── pkg/
│   └── utils/           # Shared utility functions
└── configs/             # Configuration files
```

## Prerequisites

- Go 1.21 or later
- MongoDB Atlas account (or local MongoDB instance)

## Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod download
   ```

## Usage

Run the crawler with a seed URL:

```bash
go run cmd/crawler/main.go -seed "https://example.com"
```

To use MongoDB for storing crawled pages:

```bash
go run cmd/crawler/main.go -mongo "mongodb+srv://<username>:<password>@<cluster>.mongodb.net" -seed "https://example.com"
```

## Configuration

The crawler can be configured using command-line flags:

- `-seed`: The starting URL for crawling (required)
- `-mongo`: MongoDB connection string (optional)

## License

MIT 