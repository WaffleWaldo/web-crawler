# Go Web Crawler

A high-performance, concurrent web crawler written in Go with MongoDB integration for storing crawled pages. The crawler is highly configurable and can be used with or without MongoDB storage.

## Features

- **Concurrent Crawling**: Multiple worker goroutines for efficient crawling
- **Rate Limiting**: Configurable rate limiting per worker to avoid overwhelming target servers
- **MongoDB Integration**: Optional storage of crawled pages in MongoDB
- **Configurable**: YAML-based configuration for easy customization
- **URL Filtering**: Configurable domain, path, and file type filtering
- **Graceful Shutdown**: Proper handling of shutdown signals
- **Performance Benchmarking**: Real-time performance monitoring with graph generation
- **HTTP Features**:
  - Custom User-Agent
  - Configurable redirects
  - Timeout handling
  - Retry mechanism

## Project Structure

```
.
├── cmd/
│   └── crawler/         # Main application entry point
├── configs/             # Configuration files
│   └── default.yaml     # Default configuration
├── internal/            # Internal packages
│   ├── benchmark/      # Performance benchmarking
│   ├── config/         # Configuration handling
│   ├── crawler/        # Core crawler implementation
│   ├── queue/          # URL queue management
│   └── storage/        # Storage interfaces and MongoDB implementation
├── pkg/                # Public packages
│   └── utils/          # Shared utilities
└── scripts/            # Helper scripts
    └── run.sh          # Crawler execution script
```

## Prerequisites

- Go 1.21 or later
- MongoDB (optional)
  - Local instance or
  - MongoDB Atlas account

## Installation

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd web-crawler
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

## Configuration

The crawler is configured via YAML files in the `configs/` directory. The default configuration is in `configs/default.yaml`.

### Main Configuration Sections:

1. **Crawler Settings**:
   ```yaml
   crawler:
     workers: 5              # Number of concurrent workers
     rate_limit: 500         # Milliseconds between requests
     timeout: 30             # Request timeout in seconds
     max_depth: 10           # Maximum crawl depth
     max_pages: 1000         # Maximum pages to crawl (0 for unlimited)
   ```

2. **MongoDB Settings**:
   ```yaml
   storage:
     mongodb:
       database: "webcrawler"
       collection: "webpages"
       timeout: 10           # Connection timeout
       max_pool_size: 100    # Connection pool size
   ```

3. **HTTP Settings**:
   ```yaml
   http:
     user_agent: "GoWebCrawler/1.0"
     follow_redirects: true
     max_redirects: 10
     timeout: 30
     retry:
       max_attempts: 3
       initial_delay: 1
       max_delay: 5
   ```

4. **URL Filtering**:
   ```yaml
   filters:
     allowed_domains: []     # Empty means same domain as seed
     excluded_paths:         # Paths to skip
       - "/wp-admin"
       - "/login"
     allowed_schemes:        # URL schemes to allow
       - "http"
       - "https"
     excluded_extensions:    # File extensions to skip
       - ".pdf"
       - ".jpg"
   ```

5. **Benchmark Settings**:
   ```yaml
   benchmark:
     enabled: true          # Enable performance benchmarking
     interval: 1s           # Metric collection interval
     output_dir: "benchmarks" # Directory for graph output
   ```

## Usage

### Using the Run Script

The `scripts/run.sh` script provides an easy way to run the crawler:

1. Basic usage (with default configuration):
   ```bash
   ./scripts/run.sh https://example.com
   ```

2. With MongoDB:
   ```bash
   ./scripts/run.sh -m "mongodb+srv://user:pass@cluster.mongodb.net" https://example.com
   ```

3. With custom configuration:
   ```bash
   ./scripts/run.sh -c "configs/custom.yaml" https://example.com
   ```

### Script Options

- `domain`: The starting URL for crawling
- `-m, --mongo`: MongoDB connection string (optional)
- `-c, --config`: Path to configuration file (default: configs/default.yaml)
- `-h, --help`: Show help message

## MongoDB Integration

When MongoDB integration is enabled:

1. Each crawled page is stored with:
   - URL
   - Title
   - Content
   - Extracted links
   - Crawl timestamp
   - HTTP status code
   - Content type

2. The data structure in MongoDB:
   ```json
   {
     "url": "https://example.com",
     "title": "Page Title",
     "content": "Page Content",
     "links": ["https://example.com/page1", "..."],
     "crawled_at": "2024-05-29T19:28:16Z",
     "status_code": 200,
     "content_type": "text/html"
   }
   ```

## Performance Benchmarking

The crawler includes real-time performance monitoring that generates graphs showing:

1. **Pages Crawled vs Time**: Shows the rate of page crawling over time
2. **Crawled/Queued Ratio vs Time**: Shows the efficiency of the crawler in processing the URL queue

Graphs are automatically generated in the configured output directory (default: `benchmarks/`) when the crawler stops. Each graph is saved with a timestamp in the filename for easy tracking of different crawl sessions.

### Reading the Graphs

1. **Pages Crawled vs Time**:
   - X-axis: Time in seconds since start
   - Y-axis: Total number of pages crawled
   - Helps identify crawling speed and any slowdowns

2. **Crawled/Queued Ratio vs Time**:
   - X-axis: Time in seconds since start
   - Y-axis: Ratio of crawled pages to queued URLs
   - Ratio > 1: Crawler is keeping up with new URLs
   - Ratio < 1: Queue is growing faster than crawling

### Benchmark Configuration

Adjust benchmark settings in your config file:
```yaml
benchmark:
  enabled: true          # Enable/disable benchmarking
  interval: 1s           # How often to collect metrics
  output_dir: "benchmarks" # Where to save the graphs
```

## Development

### Adding New Features

1. **Storage Backends**: Implement the `storage.Archiver` interface in `internal/storage/`
2. **URL Filters**: Add new filters in `internal/crawler/crawler.go`
3. **Configuration**: Extend the config structs in `internal/config/config.go`
4. **Benchmarks**: Add new metrics in `internal/benchmark/types.go`

### Running Tests

```bash
go test ./...
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT License - see LICENSE file for details 