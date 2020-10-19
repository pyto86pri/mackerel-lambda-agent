package metrics

import (
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	mackerel "github.com/mackerelio/mackerel-client-go"
)

// CPUGenerator ...
type CPUGenerator struct {
	Interval time.Duration
}

// Generate ...
func (g *CPUGenerator) Generate() (Values, error) {
	prev, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	time.Sleep(g.Interval)

	curr, err := cpu.Get()
	if err != nil {
		return nil, err
	}

	total := float64(curr.Total - prev.Total)

	return Values{
		"custom.lambda.extensions.cpu.user":   float64((curr.User - prev.User)) * 100.0 / total,
		"custom.lambda.extensions.cpu.nice":   float64(curr.Nice-prev.Nice) * 100.0 / total,
		"custom.lambda.extensions.cpu.system": float64(curr.System-prev.System) * 100.0 / total,
		"custom.lambda.extensions.cpu.idle":   float64(curr.Idle-prev.Idle) * 100.0 / total,
	}, nil
}

// CPUGraphDefs ...
var CPUGraphDefs = &mackerel.GraphDefsParam{
	Name:        "custom.lambda.extensions.cpu",
	DisplayName: "CPU",
	Unit:        "percentage",
	Metrics: []*mackerel.GraphDefsMetric{
		&mackerel.GraphDefsMetric{
			Name:        "user",
			DisplayName: "User",
		},
		&mackerel.GraphDefsMetric{
			Name:        "nice",
			DisplayName: "Nice",
		},
		&mackerel.GraphDefsMetric{
			Name:        "system",
			DisplayName: "System",
		},
		&mackerel.GraphDefsMetric{
			Name:        "idle",
			DisplayName: "Idle",
		},
	},
}
