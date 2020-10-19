package agent

import (
	"sync"
	"time"

	"github.com/pyto86pri/mackerel-agent-lambda/cmd/metrics"

	log "github.com/sirupsen/logrus"
)

// Agent ...
type Agent struct {
	gs []metrics.Generator
}

// New construct a new Agent
func New() *Agent {
	return &Agent{
		gs: []metrics.Generator{
			&metrics.CPUGenerator{Interval: 1 * time.Second},
			&metrics.DiskGenerator{Interval: 1 * time.Second},
			&metrics.MemoryGenerator{},
			&metrics.NetworkGenerator{},
			&metrics.LoadavgGenerator{},
		},
	}
}

// Collect collect metrics and push them to bucket
func (agent *Agent) Collect(bucket *metrics.ValuesBucket) {
	var wg sync.WaitGroup
	for _, g := range agent.gs {
		wg.Add(1)
		go func(g metrics.Generator) {
			defer wg.Done()

			values, err := g.Generate()
			if err != nil {
				log.Errorf("Failed to generate value")
				return
			}
			bucket.Push(&values)
		}(g)
	}
	wg.Wait()
}
