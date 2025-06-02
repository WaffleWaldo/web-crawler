package queue

import (
	"sync/atomic"
	"time"
)

// Priority levels for URL crawling
const (
	PriorityHigh = iota
	PriorityNormal
	PriorityLow
)

// URLItem represents a URL with priority and metadata
type URLItem struct {
	URL      string
	Priority int
	Host     string
	Depth    int
	QueuedAt time.Time // For performance tracking
}

// URLQueue is a high-performance priority queue using channels with enhanced buffering
type URLQueue struct {
	highPriority   chan URLItem
	normalPriority chan URLItem
	lowPriority    chan URLItem
	size           int64
	closed         int64

	// Performance counters
	totalQueued   int64
	totalDequeued int64
	highCount     int64
	normalCount   int64
	lowCount      int64
}

// NewURLQueue creates a new high-performance URL queue with enhanced buffering
func NewURLQueue() *URLQueue {
	return &URLQueue{
		// Significantly increased buffer sizes for better performance
		highPriority:   make(chan URLItem, 2000),  // Increased from 1000
		normalPriority: make(chan URLItem, 20000), // Increased from 10000
		lowPriority:    make(chan URLItem, 10000), // Increased from 5000
	}
}

// Push adds a URL to the appropriate priority queue
func (q *URLQueue) Push(url string) {
	q.PushWithPriority(url, PriorityNormal, "", 0)
}

// PushWithPriority adds a URL with specific priority and metadata (enhanced)
func (q *URLQueue) PushWithPriority(url string, priority int, host string, depth int) {
	if atomic.LoadInt64(&q.closed) == 1 {
		return
	}

	item := URLItem{
		URL:      url,
		Priority: priority,
		Host:     host,
		Depth:    depth,
		QueuedAt: time.Now(),
	}

	// Enhanced non-blocking push with improved fallback strategy
	switch priority {
	case PriorityHigh:
		select {
		case q.highPriority <- item:
			atomic.AddInt64(&q.size, 1)
			atomic.AddInt64(&q.totalQueued, 1)
			atomic.AddInt64(&q.highCount, 1)
		default:
			// High priority queue full, try urgent fallback to normal
			select {
			case q.normalPriority <- item:
				atomic.AddInt64(&q.size, 1)
				atomic.AddInt64(&q.totalQueued, 1)
				atomic.AddInt64(&q.normalCount, 1)
			default:
				// Both full, drop to prevent blocking (rare case)
			}
		}
	case PriorityLow:
		select {
		case q.lowPriority <- item:
			atomic.AddInt64(&q.size, 1)
			atomic.AddInt64(&q.totalQueued, 1)
			atomic.AddInt64(&q.lowCount, 1)
		default:
			// Low priority queue full, just drop (acceptable for low priority)
		}
	default: // PriorityNormal
		select {
		case q.normalPriority <- item:
			atomic.AddInt64(&q.size, 1)
			atomic.AddInt64(&q.totalQueued, 1)
			atomic.AddInt64(&q.normalCount, 1)
		default:
			// Normal queue full, try low priority as fallback
			select {
			case q.lowPriority <- item:
				atomic.AddInt64(&q.size, 1)
				atomic.AddInt64(&q.totalQueued, 1)
				atomic.AddInt64(&q.lowCount, 1)
			default:
				// Both full, drop to prevent blocking
			}
		}
	}
}

// Pop removes and returns the highest priority URL available (enhanced)
// Returns empty URLItem and false if no URLs are available
func (q *URLQueue) Pop() (URLItem, bool) {
	// Try high priority first (with higher probability)
	select {
	case item := <-q.highPriority:
		atomic.AddInt64(&q.size, -1)
		atomic.AddInt64(&q.totalDequeued, 1)
		return item, true
	default:
	}

	// Try normal priority with higher frequency than low
	select {
	case item := <-q.normalPriority:
		atomic.AddInt64(&q.size, -1)
		atomic.AddInt64(&q.totalDequeued, 1)
		return item, true
	default:
		// Only check low priority if normal is empty
		select {
		case item := <-q.lowPriority:
			atomic.AddInt64(&q.size, -1)
			atomic.AddInt64(&q.totalDequeued, 1)
			return item, true
		default:
		}
	}

	return URLItem{}, false
}

// PopBlocking waits for a URL to become available with enhanced priority ordering
func (q *URLQueue) PopBlocking() (URLItem, bool) {
	if atomic.LoadInt64(&q.closed) == 1 {
		return URLItem{}, false
	}

	// Enhanced priority selection with weighted approach
	select {
	case item := <-q.highPriority:
		atomic.AddInt64(&q.size, -1)
		atomic.AddInt64(&q.totalDequeued, 1)
		return item, true
	default:
		// Weighted selection between normal and low priority
		select {
		case item := <-q.highPriority:
			atomic.AddInt64(&q.size, -1)
			atomic.AddInt64(&q.totalDequeued, 1)
			return item, true
		case item := <-q.normalPriority:
			atomic.AddInt64(&q.size, -1)
			atomic.AddInt64(&q.totalDequeued, 1)
			return item, true
		case item := <-q.lowPriority:
			atomic.AddInt64(&q.size, -1)
			atomic.AddInt64(&q.totalDequeued, 1)
			return item, true
		}
	}
}

// PopBatch returns multiple items at once for batch processing (new method)
func (q *URLQueue) PopBatch(maxItems int) []URLItem {
	items := make([]URLItem, 0, maxItems)

	for i := 0; i < maxItems; i++ {
		item, ok := q.Pop()
		if !ok {
			break
		}
		items = append(items, item)
	}

	return items
}

// Size returns the current approximate size of all queues
func (q *URLQueue) Size() int {
	return int(atomic.LoadInt64(&q.size))
}

// GetStats returns detailed queue statistics for monitoring
func (q *URLQueue) GetStats() map[string]int64 {
	return map[string]int64{
		"size":          atomic.LoadInt64(&q.size),
		"totalQueued":   atomic.LoadInt64(&q.totalQueued),
		"totalDequeued": atomic.LoadInt64(&q.totalDequeued),
		"highCount":     atomic.LoadInt64(&q.highCount),
		"normalCount":   atomic.LoadInt64(&q.normalCount),
		"lowCount":      atomic.LoadInt64(&q.lowCount),
		"highBuffer":    int64(len(q.highPriority)),
		"normalBuffer":  int64(len(q.normalPriority)),
		"lowBuffer":     int64(len(q.lowPriority)),
	}
}

// IsFull checks if any queue is approaching capacity
func (q *URLQueue) IsFull() bool {
	highFull := len(q.highPriority) > cap(q.highPriority)*9/10 // 90% full
	normalFull := len(q.normalPriority) > cap(q.normalPriority)*9/10
	lowFull := len(q.lowPriority) > cap(q.lowPriority)*9/10

	return highFull || normalFull || lowFull
}

// Close closes the queue and prevents new items from being added
func (q *URLQueue) Close() {
	atomic.StoreInt64(&q.closed, 1)
	close(q.highPriority)
	close(q.normalPriority)
	close(q.lowPriority)
}
