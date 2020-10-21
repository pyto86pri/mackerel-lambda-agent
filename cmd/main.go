package main

import (
	"context"

	"github.com/pyto86pri/mackerel-lambda-agent/cmd/config"
	"github.com/pyto86pri/mackerel-lambda-agent/cmd/metrics"

	mackerel "github.com/mackerelio/mackerel-client-go"
	"github.com/pyto86pri/mackerel-lambda-agent/cmd/agent"
	"github.com/pyto86pri/mackerel-lambda-agent/cmd/app"
	"github.com/pyto86pri/mackerel-lambda-agent/cmd/extensions"
	log "github.com/sirupsen/logrus"
)

var version, revision string

func main() {
	log.Printf("Starting mackerel-lambda-agent (version:%s, revision:%s)", version, revision)
	conf, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config", err)
	}
	ec, err := extensions.NewClient()
	if err != nil {
		log.Fatal("Failed to initialize extensions API client", err)
	}
	mc := mackerel.NewClient(conf.APIKey)
	if conf.APIBase != nil {
		mc.BaseURL = conf.APIBase
	}
	agent := agent.New()
	bucket := metrics.NewBucket()
	app := &app.App{
		Version:          version,
		Revision:         revision,
		Config:           conf,
		MackerelClient:   mc,
		ExtensionsClient: ec,
		Agent:            agent,
		Bucket:           bucket,
	}

	ctx := context.Background()
	app.Run(ctx)
}
