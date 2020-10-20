package main

import (
	"context"
	"os"

	"github.com/pyto86pri/mackerel-agent-lambda/cmd/metrics"

	mackerel "github.com/mackerelio/mackerel-client-go"
	"github.com/pyto86pri/mackerel-agent-lambda/cmd/agent"
	"github.com/pyto86pri/mackerel-agent-lambda/cmd/app"
	"github.com/pyto86pri/mackerel-agent-lambda/cmd/extensions"
	log "github.com/sirupsen/logrus"
)

var version, revision string

func main() {
	log.Printf("Starting mackerel-agent-lambda (version:%s, revision:%s)", version, revision)
	ec, err := extensions.NewClient()
	if err != nil {
		log.Fatal("Failed to initialize extensions API client")
	}
	apiKey := os.Getenv("MACKEREL_API_KEY")
	if apiKey == "" {
		log.Fatal("Please set mackerel api key in environment variables")
	}
	mc := mackerel.NewClient(apiKey)
	agent := agent.New()
	bucket := metrics.NewBucket()
	app := &app.App{
		Version:          version,
		Revision:         revision,
		MackerelClient:   mc,
		ExtensionsClient: ec,
		Agent:            agent,
		Bucket:           bucket,
	}

	ctx := context.Background()
	app.Run(ctx)
}
