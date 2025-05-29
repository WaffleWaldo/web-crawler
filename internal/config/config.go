package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure
type Config struct {
	Crawler CrawlerConfig `yaml:"crawler"`
	Storage StorageConfig `yaml:"storage"`
	HTTP    HTTPConfig    `yaml:"http"`
	Filters FiltersConfig `yaml:"filters"`
}

// CrawlerConfig holds crawler-specific settings
type CrawlerConfig struct {
	Workers   int `yaml:"workers"`
	RateLimit int `yaml:"rate_limit"` // Milliseconds between requests
	Timeout   int `yaml:"timeout"`
	MaxDepth  int `yaml:"max_depth"`
	MaxPages  int `yaml:"max_pages"`
}

// GetRateLimit returns the rate limit as a time.Duration
func (c *CrawlerConfig) GetRateLimit() time.Duration {
	return time.Duration(c.RateLimit) * time.Millisecond
}

// StorageConfig holds storage-related settings
type StorageConfig struct {
	MongoDB MongoDBConfig `yaml:"mongodb"`
}

// MongoDBConfig holds MongoDB-specific settings
type MongoDBConfig struct {
	Database    string `yaml:"database"`
	Collection  string `yaml:"collection"`
	Timeout     int    `yaml:"timeout"`
	MaxPoolSize uint64 `yaml:"max_pool_size"`
}

// HTTPConfig holds HTTP client settings
type HTTPConfig struct {
	UserAgent      string      `yaml:"user_agent"`
	FollowRedirect bool        `yaml:"follow_redirects"`
	MaxRedirects   int         `yaml:"max_redirects"`
	Timeout        int         `yaml:"timeout"`
	Retry          RetryConfig `yaml:"retry"`
}

// RetryConfig holds retry settings
type RetryConfig struct {
	MaxAttempts  int `yaml:"max_attempts"`
	InitialDelay int `yaml:"initial_delay"`
	MaxDelay     int `yaml:"max_delay"`
}

// FiltersConfig holds URL filtering settings
type FiltersConfig struct {
	AllowedDomains     []string `yaml:"allowed_domains"`
	ExcludedPaths      []string `yaml:"excluded_paths"`
	AllowedSchemes     []string `yaml:"allowed_schemes"`
	ExcludedExtensions []string `yaml:"excluded_extensions"`
}

// LoadConfig loads configuration from a YAML file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Crawler: CrawlerConfig{
			Workers:   5,
			RateLimit: 500, // 500ms between requests
			Timeout:   30,
			MaxDepth:  10,
			MaxPages:  1000,
		},
		Storage: StorageConfig{
			MongoDB: MongoDBConfig{
				Database:    "webcrawler",
				Collection:  "webpages",
				Timeout:     10,
				MaxPoolSize: 100,
			},
		},
		HTTP: HTTPConfig{
			UserAgent:      "GoWebCrawler/1.0",
			FollowRedirect: true,
			MaxRedirects:   10,
			Timeout:        30,
			Retry: RetryConfig{
				MaxAttempts:  3,
				InitialDelay: 1,
				MaxDelay:     5,
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
	}
}
