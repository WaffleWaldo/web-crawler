package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Crawler   CrawlerConfig   `yaml:"crawler"`
	Storage   StorageConfig   `yaml:"storage"`
	HTTP      HTTPConfig      `yaml:"http"`
	Filters   FiltersConfig   `yaml:"filters"`
	Benchmark BenchmarkConfig `yaml:"benchmark"`
}

// CrawlerConfig holds crawler-specific settings
type CrawlerConfig struct {
	Workers   int           `yaml:"workers"`
	RateLimit time.Duration `yaml:"rate_limit"`
	Timeout   time.Duration `yaml:"timeout"`
	MaxDepth  int           `yaml:"max_depth"`
	MaxPages  int           `yaml:"max_pages"`
}

// GetRateLimit returns the rate limit as a time.Duration
func (c *CrawlerConfig) GetRateLimit() time.Duration {
	return c.RateLimit
}

// StorageConfig holds storage-related settings
type StorageConfig struct {
	MongoDB MongoDBConfig `yaml:"mongodb"`
}

// MongoDBConfig holds MongoDB-specific settings
type MongoDBConfig struct {
	Database    string        `yaml:"database"`
	Collection  string        `yaml:"collection"`
	Timeout     time.Duration `yaml:"timeout"`
	MaxPoolSize uint64        `yaml:"max_pool_size"`
	MinPoolSize uint64        `yaml:"min_pool_size"`
	MaxIdleTime time.Duration `yaml:"max_idle_time"`
}

// HTTPConfig holds HTTP client settings
type HTTPConfig struct {
	UserAgent      string        `yaml:"user_agent"`
	FollowRedirect bool          `yaml:"follow_redirects"`
	MaxRedirects   int           `yaml:"max_redirects"`
	Timeout        time.Duration `yaml:"timeout"`
	Retry          RetryConfig   `yaml:"retry"`
}

// RetryConfig holds retry settings
type RetryConfig struct {
	MaxAttempts  int           `yaml:"max_attempts"`
	InitialDelay time.Duration `yaml:"initial_delay"`
	MaxDelay     time.Duration `yaml:"max_delay"`
}

// FiltersConfig holds URL filtering settings
type FiltersConfig struct {
	AllowedDomains     []string `yaml:"allowed_domains"`
	ExcludedPaths      []string `yaml:"excluded_paths"`
	AllowedSchemes     []string `yaml:"allowed_schemes"`
	ExcludedExtensions []string `yaml:"excluded_extensions"`
}

// BenchmarkConfig holds benchmark settings
type BenchmarkConfig struct {
	Enabled   bool          `yaml:"enabled"`
	Interval  time.Duration `yaml:"interval"`
	OutputDir string        `yaml:"output_dir"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Crawler: CrawlerConfig{
			Workers:   5,
			RateLimit: 500 * time.Millisecond, // 500ms between requests
			Timeout:   30 * time.Second,
			MaxDepth:  10,
			MaxPages:  1000,
		},
		Storage: StorageConfig{
			MongoDB: MongoDBConfig{
				Database:    "webcrawler",
				Collection:  "webpages",
				Timeout:     30 * time.Second,
				MaxPoolSize: 50,
				MinPoolSize: 10,
				MaxIdleTime: 5 * time.Minute,
			},
		},
		HTTP: HTTPConfig{
			UserAgent:      "GoWebCrawler/1.0",
			FollowRedirect: true,
			MaxRedirects:   10,
			Timeout:        30 * time.Second,
			Retry: RetryConfig{
				MaxAttempts:  3,
				InitialDelay: 1 * time.Second,
				MaxDelay:     5 * time.Second,
			},
		},
		Filters: FiltersConfig{
			AllowedDomains: []string{},
			ExcludedPaths: []string{
				"/wp-admin",
				"/wp-login",
				"/wp-content",
				"/admin",
				"/login",
			},
			AllowedSchemes: []string{
				"http",
				"https",
			},
			ExcludedExtensions: []string{
				".pdf", ".jpg", ".jpeg", ".png", ".gif",
				".zip", ".tar", ".gz", ".rar", ".exe",
				".doc", ".docx", ".xls", ".xlsx",
				".ppt", ".pptx",
			},
		},
		Benchmark: BenchmarkConfig{
			Enabled:   false,
			Interval:  0,
			OutputDir: "",
		},
	}
}
