package metrics

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Songmu/timeout"
	mackerel "github.com/mackerelio/mackerel-client-go"
	log "github.com/sirupsen/logrus"
)

// FilesystemGenerator ...
type FilesystemGenerator struct {
	Interval time.Duration
}

// Generate ...
func (g *FilesystemGenerator) Generate() (Values, error) {
	filesystems, err := Get()
	if err != nil {
		return nil, err
	}

	values := make(Values)
	for _, filesystem := range filesystems {
		// AWS Lambda only allows to use "/tmp"
		if filesystem.Mounted == "/tmp" {
			base := filepath.Base(filesystem.Mounted)
			values["custom.aws.lambda.extensions.filesystem."+base+".used"] = float64(filesystem.Used) * 1024
			values["custom.aws.lambda.extensions.filesystem."+base+".available"] = float64(filesystem.Available) * 1024
			values["custom.aws.lambda.extensions.filesystem."+base+".total"] = float64(filesystem.Used+filesystem.Available) * 1024
		}
	}

	return values, nil
}

// FilesystemGraphDefs ...
var FilesystemGraphDefs = &mackerel.GraphDefsParam{
	Name:        "custom.aws.lambda.extensions.filesystem.#",
	DisplayName: "Filesystem",
	Unit:        "bytes",
	Metrics: []*mackerel.GraphDefsMetric{
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.filesystem.#.used",
			DisplayName: "Used",
			IsStacked:   true,
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.filesystem.#.available",
			DisplayName: "Available",
			IsStacked:   true,
		},
		&mackerel.GraphDefsMetric{
			Name:        "custom.aws.lambda.extensions.filesystem.#.total",
			DisplayName: "Total",
		},
	},
}

// Stat ...
type Stat struct {
	Name      string
	Blocks    uint64
	Used      uint64
	Available uint64
	Capacity  uint8
	Mounted   string
}

// Get ...
func Get() ([]*Stat, error) {
	cmd := exec.Command("df")
	to := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  15 * time.Second,
		KillAfter: 5 * time.Second,
	}
	exitStatus, stdout, stderr, err := to.Run()
	if err != nil {
		return nil, fmt.Errorf("Failed to invoke 'df' command: %q", err)
	}
	if exitStatus.Code != 0 {
		return nil, fmt.Errorf("'df' command exited with a non-zero status: %d: %q", exitStatus.Code, stderr)
	}
	stats, err := parseLines(stdout)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

func parseLines(out string) ([]*Stat, error) {
	s := bufio.NewScanner(strings.NewReader(out))
	var filesystems []*Stat
	// Skip headers
	if s.Scan() {
		for s.Scan() {
			line := s.Text()
			stat, err := parseLine(line)
			if err != nil {
				log.Warningf("Failed to parse df line: %q", err)
				continue
			}
			filesystems = append(filesystems, stat)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return filesystems, nil
}

func parseLine(line string) (*Stat, error) {
	matches := regexp.MustCompile(`^(.+?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(\d+)%\s+(.+)$`).FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("Failed to parse line: [%s]", line)
	}
	name := matches[1]
	blocks, _ := strconv.ParseUint(matches[2], 0, 64)
	used, _ := strconv.ParseUint(matches[3], 0, 64)
	available, _ := strconv.ParseUint(matches[4], 0, 64)
	capacity, _ := strconv.ParseUint(matches[5], 0, 8)
	mounted := matches[6]

	return &Stat{
		Name:      name,
		Blocks:    blocks,
		Used:      used,
		Available: available,
		Capacity:  uint8(capacity),
		Mounted:   mounted,
	}, nil
}
