package benchmark

import (
	"sync"
	"time"
)

// Metric represents a single data point in time
type Metric struct {
	Timestamp   time.Time
	PagesCount  int
	QueuedCount int
}

// Recorder handles the collection and storage of benchmark metrics
type Recorder struct {
	metrics []Metric
	mu      sync.RWMutex
	start   time.Time
}

// New creates a new benchmark recorder
func New() *Recorder {
	return &Recorder{
		metrics: make([]Metric, 0),
		start:   time.Now(),
	}
}

// Record adds a new metric point
func (r *Recorder) Record(pagesCount, queuedCount int) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.metrics = append(r.metrics, Metric{
		Timestamp:   time.Now(),
		PagesCount:  pagesCount,
		QueuedCount: queuedCount,
	})
}

// GetMetrics returns a copy of all recorded metrics
func (r *Recorder) GetMetrics() []Metric {
	r.mu.RLock()
	defer r.mu.RUnlock()

	metrics := make([]Metric, len(r.metrics))
	copy(metrics, r.metrics)
	return metrics
}

// ElapsedSeconds returns the number of seconds since the recorder started
func (r *Recorder) ElapsedSeconds() float64 {
	return time.Since(r.start).Seconds()
}
