package metrics

import (
	"github.com/mackerelio/go-osstat/loadavg"
	mackerel "github.com/mackerelio/mackerel-client-go"
)

// LoadavgGenerator ...
type LoadavgGenerator struct{}

// Generate ...
func (g *LoadavgGenerator) Generate() (Values, error) {
	loadavg, err := loadavg.Get()
	if err != nil {
		return nil, err
	}

	return Values{
		"custom.lambda.extensions.loadavg.loadavg1":  loadavg.Loadavg1,
		"custom.lambda.extensions.loadavg.loadavg5":  loadavg.Loadavg5,
		"custom.lambda.extensions.loadavg.loadavg15": loadavg.Loadavg15,
	}, nil
}

// LoadavgGraphDefs ...
var LoadavgGraphDefs = &mackerel.GraphDefsParam{
	Name:        "custom.lambda.extensions.loadavg",
	DisplayName: "Loadavg",
	Metrics: []*mackerel.GraphDefsMetric{
		&mackerel.GraphDefsMetric{
			Name:        "loadavg1",
			DisplayName: "Loadavg1",
		},
		&mackerel.GraphDefsMetric{
			Name:        "loadavg5",
			DisplayName: "Loadavg5",
		},
		&mackerel.GraphDefsMetric{
			Name:        "loadavg15",
			DisplayName: "Loadavg15",
		},
	},
}
