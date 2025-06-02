package utils

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ContentSaver handles saving crawled page content to files
type ContentSaver struct {
	baseDir     string
	enabled     bool
	maxFileSize int64 // Maximum file size to save (in bytes)
}

// NewContentSaver creates a new content saver
func NewContentSaver(baseDir string, enabled bool, maxFileSize int64) *ContentSaver {
	return &ContentSaver{
		baseDir:     baseDir,
		enabled:     enabled,
		maxFileSize: maxFileSize,
	}
}

// SavePageContent saves page content to a file with metadata
func (cs *ContentSaver) SavePageContent(pageURL, title, content, contentType string, statusCode int, crawledAt time.Time) error {
	if !cs.enabled {
		return nil
	}

	// Skip if content is too large
	if cs.maxFileSize > 0 && int64(len(content)) > cs.maxFileSize {
		return nil
	}

	// Create safe filename from URL
	filename := cs.createSafeFilename(pageURL)

	// Create domain-based directory structure
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		return err
	}

	domainDir := filepath.Join(cs.baseDir, cs.sanitizeDomain(parsedURL.Host))
	if err := os.MkdirAll(domainDir, 0755); err != nil {
		return err
	}

	// Create full file path
	filePath := filepath.Join(domainDir, filename+".html")

	// Create metadata header
	metadata := cs.createMetadataHeader(pageURL, title, contentType, statusCode, crawledAt, len(content))

	// Combine metadata and content
	fullContent := metadata + "\n\n" + content

	// Write to file
	return os.WriteFile(filePath, []byte(fullContent), 0644)
}

// createSafeFilename creates a safe filename from URL
func (cs *ContentSaver) createSafeFilename(pageURL string) string {
	parsedURL, err := url.Parse(pageURL)
	if err != nil {
		// Use MD5 hash as fallback
		hash := md5.Sum([]byte(pageURL))
		return fmt.Sprintf("page_%x", hash)
	}

	// Use path for filename
	path := parsedURL.Path
	if path == "" || path == "/" {
		path = "index"
	}

	// Remove leading slash and replace special characters
	path = strings.TrimPrefix(path, "/")
	path = strings.ReplaceAll(path, "/", "_")
	path = strings.ReplaceAll(path, "?", "_")
	path = strings.ReplaceAll(path, "&", "_")
	path = strings.ReplaceAll(path, "=", "_")
	path = strings.ReplaceAll(path, "#", "_")
	path = strings.ReplaceAll(path, "%", "_")
	path = strings.ReplaceAll(path, " ", "_")

	// Limit filename length
	if len(path) > 100 {
		// Use first 80 chars + hash of full path
		hash := md5.Sum([]byte(parsedURL.Path))
		path = path[:80] + fmt.Sprintf("_%x", hash)[:16]
	}

	// Add query parameters if short enough
	if parsedURL.RawQuery != "" && len(path) < 80 {
		query := strings.ReplaceAll(parsedURL.RawQuery, "=", "_")
		query = strings.ReplaceAll(query, "&", "_")
		if len(path)+len(query) < 100 {
			path += "_" + query
		}
	}

	return path
}

// sanitizeDomain creates a safe directory name from domain
func (cs *ContentSaver) sanitizeDomain(domain string) string {
	// Remove www prefix
	domain = strings.TrimPrefix(domain, "www.")

	// Replace dots with underscores for directory safety
	domain = strings.ReplaceAll(domain, ".", "_")
	domain = strings.ReplaceAll(domain, ":", "_")

	return domain
}

// createMetadataHeader creates HTML comment with page metadata
func (cs *ContentSaver) createMetadataHeader(pageURL, title, contentType string, statusCode int, crawledAt time.Time, contentSize int) string {
	return fmt.Sprintf(`<!--
CRAWLED PAGE METADATA
=====================
URL: %s
Title: %s
Content-Type: %s
Status Code: %d
Crawled At: %s
Content Size: %d bytes
Crawler: Ultra-High-Performance Go Web Crawler
=====================
-->`, pageURL, title, contentType, statusCode, crawledAt.Format(time.RFC3339), contentSize)
}

// GetSavedFiles returns a list of all saved files
func (cs *ContentSaver) GetSavedFiles() ([]string, error) {
	if !cs.enabled {
		return nil, nil
	}

	var files []string
	err := filepath.Walk(cs.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}

// GetStats returns statistics about saved content
func (cs *ContentSaver) GetStats() (map[string]interface{}, error) {
	if !cs.enabled {
		return map[string]interface{}{"enabled": false}, nil
	}

	files, err := cs.GetSavedFiles()
	if err != nil {
		return nil, err
	}

	var totalSize int64
	domainCount := make(map[string]int)

	for _, file := range files {
		// Get file size
		if info, err := os.Stat(file); err == nil {
			totalSize += info.Size()
		}

		// Count by domain
		dir := filepath.Dir(file)
		domain := filepath.Base(dir)
		domainCount[domain]++
	}

	return map[string]interface{}{
		"enabled":      true,
		"total_files":  len(files),
		"total_size":   totalSize,
		"domains":      len(domainCount),
		"domain_count": domainCount,
		"base_dir":     cs.baseDir,
	}, nil
}
