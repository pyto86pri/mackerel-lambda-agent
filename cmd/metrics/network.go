package metrics

import (
	"github.com/mackerelio/go-osstat/network"
	mackerel "github.com/mackerelio/mackerel-client-go"
)

// NetworkGenerator ...
type NetworkGenerator struct{}

// Generate ...
func (g *NetworkGenerator) Generate() (Values, error) {
	networks, err := network.Get()
	if err != nil {
		return nil, err
	}

	var in uint64
	for _, network := range networks {
		in += network.RxBytes
	}
	var out uint64
	for _, network := range networks {
		out += network.TxBytes
	}

	return Values{
		"custom.aws.lambda.extensions.network.in":  float64(in),
		"custom.aws.lambda.extensions.network.out": float64(out),
	}, nil
}

// NetworkGraphDefs ...
var NetworkGraphDefs = &mackerel.GraphDefsParam{
	Name:        "custom.aws.lambda.extensions.network",
	DisplayName: "Network I/O",
	Unit:        "bytes",
	Metrics: []*mackerel.GraphDefsMetric{
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.network.in",
			DisplayName: "In",
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.network.out",
			DisplayName: "Out",
		},
	},
}
