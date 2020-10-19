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
		"custom.aws.lambda.extensions.cpu.user":   float64((curr.User - prev.User)) * 100.0 / total,
		"custom.aws.lambda.extensions.cpu.nice":   float64(curr.Nice-prev.Nice) * 100.0 / total,
		"custom.aws.lambda.extensions.cpu.system": float64(curr.System-prev.System) * 100.0 / total,
		"custom.aws.lambda.extensions.cpu.idle":   float64(curr.Idle-prev.Idle) * 100.0 / total,
	}, nil
}

// CPUGraphDefs ...
var CPUGraphDefs = &mackerel.GraphDefsParam{
	Name:        "custom.aws.lambda.extensions.cpu",
	DisplayName: "CPU",
	Unit:        "percentage",
	Metrics: []*mackerel.GraphDefsMetric{
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.cpu.user",
			DisplayName: "User",
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.cpu.nice",
			DisplayName: "Nice",
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.cpu.system",
			DisplayName: "System",
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.cpu.idle",
			DisplayName: "Idle",
		},
	},
}
