package metrics

import (
	"github.com/mackerelio/go-osstat/memory"
	mackerel "github.com/mackerelio/mackerel-client-go"
)

// MemoryGenerator ...
type MemoryGenerator struct{}

// Generate ...
func (g *MemoryGenerator) Generate() (Values, error) {
	memory, err := memory.Get()
	if err != nil {
		return nil, err
	}

	return Values{
		"custom.aws.lambda.extensions.memory.total":  float64(memory.Total),
		"custom.aws.lambda.extensions.memory.used":   float64(memory.Used),
		"custom.aws.lambda.extensions.memory.cached": float64(memory.Cached),
		"custom.aws.lambda.extensions.memory.free":   float64(memory.Free),
	}, nil
}

// MemoryGraphDefs ...
var MemoryGraphDefs = &mackerel.GraphDefsParam{
	Name:        "custom.aws.lambda.extensions.memory",
	DisplayName: "Memory",
	Unit:        "bytes",
	Metrics: []*mackerel.GraphDefsMetric{
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.memory.total",
			DisplayName: "Total",
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.memory.used",
			DisplayName: "Used",
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.memory.cached",
			DisplayName: "Cached",
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.memory.free",
			DisplayName: "Free",
		},
	},
}
