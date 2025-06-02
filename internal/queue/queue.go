package queue

import (
	"sync/atomic"
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
}

// URLQueue is a high-performance priority queue using channels
type URLQueue struct {
	highPriority   chan URLItem
	normalPriority chan URLItem
	lowPriority    chan URLItem
	size           int64
	closed         int64
}

// NewURLQueue creates a new high-performance URL queue
func NewURLQueue() *URLQueue {
	return &URLQueue{
		highPriority:   make(chan URLItem, 1000),  // High priority buffer
		normalPriority: make(chan URLItem, 10000), // Normal priority buffer
		lowPriority:    make(chan URLItem, 5000),  // Low priority buffer
	}
}

// Push adds a URL to the appropriate priority queue
func (q *URLQueue) Push(url string) {
	q.PushWithPriority(url, PriorityNormal, "", 0)
}

// PushWithPriority adds a URL with specific priority and metadata
func (q *URLQueue) PushWithPriority(url string, priority int, host string, depth int) {
	if atomic.LoadInt64(&q.closed) == 1 {
		return
	}

	item := URLItem{
		URL:      url,
		Priority: priority,
		Host:     host,
		Depth:    depth,
	}

	// Non-blocking push with priority queuing
	switch priority {
	case PriorityHigh:
		select {
		case q.highPriority <- item:
			atomic.AddInt64(&q.size, 1)
		default:
			// High priority queue full, fallback to normal
			select {
			case q.normalPriority <- item:
				atomic.AddInt64(&q.size, 1)
			default:
				// Skip if all queues are full to prevent blocking
			}
		}
	case PriorityLow:
		select {
		case q.lowPriority <- item:
			atomic.AddInt64(&q.size, 1)
		default:
			// Skip if low priority queue is full
		}
	default: // PriorityNormal
		select {
		case q.normalPriority <- item:
			atomic.AddInt64(&q.size, 1)
		default:
			// Normal queue full, try low priority
			select {
			case q.lowPriority <- item:
				atomic.AddInt64(&q.size, 1)
			default:
				// Skip if all queues are full
			}
		}
	}
}

// Pop removes and returns the highest priority URL available
// Returns empty string and false if no URLs are available
func (q *URLQueue) Pop() (URLItem, bool) {
	// Try high priority first
	select {
	case item := <-q.highPriority:
		atomic.AddInt64(&q.size, -1)
		return item, true
	default:
	}

	// Then normal priority
	select {
	case item := <-q.normalPriority:
		atomic.AddInt64(&q.size, -1)
		return item, true
	default:
	}

	// Finally low priority
	select {
	case item := <-q.lowPriority:
		atomic.AddInt64(&q.size, -1)
		return item, true
	default:
	}

	return URLItem{}, false
}

// PopBlocking waits for a URL to become available with priority ordering
func (q *URLQueue) PopBlocking() (URLItem, bool) {
	if atomic.LoadInt64(&q.closed) == 1 {
		return URLItem{}, false
	}

	// Use select with priority ordering
	select {
	case item := <-q.highPriority:
		atomic.AddInt64(&q.size, -1)
		return item, true
	default:
		select {
		case item := <-q.highPriority:
			atomic.AddInt64(&q.size, -1)
			return item, true
		case item := <-q.normalPriority:
			atomic.AddInt64(&q.size, -1)
			return item, true
		case item := <-q.lowPriority:
			atomic.AddInt64(&q.size, -1)
			return item, true
		}
	}
}

// Size returns the current approximate size of all queues
func (q *URLQueue) Size() int {
	return int(atomic.LoadInt64(&q.size))
}

// Close closes the queue and prevents new items from being added
func (q *URLQueue) Close() {
	atomic.StoreInt64(&q.closed, 1)
	close(q.highPriority)
	close(q.normalPriority)
	close(q.lowPriority)
}
