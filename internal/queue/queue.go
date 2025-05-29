package queue

import (
	"sync"
)

// URLQueue manages the queue of URLs to be crawled
type URLQueue struct {
	urls  []string
	mutex sync.Mutex
}

// NewURLQueue creates a new URL queue
func NewURLQueue() *URLQueue {
	return &URLQueue{
		urls: make([]string, 0),
	}
}

// Push adds a URL to the queue
func (q *URLQueue) Push(url string) {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	q.urls = append(q.urls, url)
}

// Pop removes and returns a URL from the queue
func (q *URLQueue) Pop() (string, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.urls) == 0 {
		return "", false
	}

	url := q.urls[0]
	q.urls = q.urls[1:]
	return url, true
}

// Size returns the current size of the queue
func (q *URLQueue) Size() int {
	q.mutex.Lock()
	defer q.mutex.Unlock()
	return len(q.urls)
}
