package metrics

import (
	"sync"
	"time"
)

// Values ...
type Values map[string]float64

// Generator ...
type Generator interface {
	Generate() (Values, error)
}

// ValuesBucket ...
type ValuesBucket struct {
	mu sync.Mutex
	v  []*Values
	f  int64
}

// NewBucket ...
func NewBucket() *ValuesBucket {
	return &ValuesBucket{v: []*Values{}}
}

// Push ...
func (b *ValuesBucket) Push(values *Values) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.v = append(b.v, values)
}

// Flush ...
func (b *ValuesBucket) Flush() []*Values {
	b.mu.Lock()
	defer b.mu.Unlock()
	vs := b.v
	b.v = []*Values{}
	b.f = time.Now().Unix()
	return vs
}

// LastFlushedAt ...
func (b *ValuesBucket) LastFlushedAt() int64 {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.f
}
