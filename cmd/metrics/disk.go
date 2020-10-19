package metrics

import (
	"time"

	"github.com/mackerelio/go-osstat/disk"
	mackerel "github.com/mackerelio/mackerel-client-go"
)

// DiskGenerator ...
type DiskGenerator struct {
	Interval time.Duration
}

func calcTotalReadsWrites(ds []disk.Stats) (reads, writes uint64) {
	for _, d := range ds {
		reads += d.ReadsCompleted
	}
	for _, d := range ds {
		writes += d.WritesCompleted
	}
	return
}

// Generate ...
func (g *DiskGenerator) Generate() (Values, error) {
	prev, err := disk.Get()
	if err != nil {
		return nil, err
	}
	prevReads, prevWrites := calcTotalReadsWrites(prev)

	time.Sleep(g.Interval)

	curr, err := disk.Get()
	if err != nil {
		return nil, err
	}
	currReads, currWrites := calcTotalReadsWrites(curr)

	return Values{
		"custom.aws.lambda.extensions.disk.reads":  float64(currReads-prevReads) / g.Interval.Seconds(),
		"custom.aws.lambda.extensions.disk.writes": float64(currWrites-prevWrites) / g.Interval.Seconds(),
	}, nil
}

// DiskGraphDefs ...
var DiskGraphDefs = &mackerel.GraphDefsParam{
	Name:        "custom.aws.lambda.extensions.disk",
	DisplayName: "Disk I/O",
	Unit:        "iops",
	Metrics: []*mackerel.GraphDefsMetric{
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.disk.reads",
			DisplayName: "Reads",
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.disk.writes",
			DisplayName: "Writes",
		},
	},
}
