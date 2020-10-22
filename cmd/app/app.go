package app

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"

	"github.com/pyto86pri/mackerel-lambda-agent/cmd/agent"
	"github.com/pyto86pri/mackerel-lambda-agent/cmd/config"
	"github.com/pyto86pri/mackerel-lambda-agent/cmd/extensions"
	"github.com/pyto86pri/mackerel-lambda-agent/cmd/libs"
	"github.com/pyto86pri/mackerel-lambda-agent/cmd/metrics"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/mackerelio/mackerel-client-go"
	log "github.com/sirupsen/logrus"
)

// MetaData ...
type MetaData struct {
	FunctionArn string `json:"functionArn"`
	Version     string `json:"version"`
	MemorySize  string `json:"memorySize"`
}

// App ...
type App struct {
	Version          string
	Revision         string
	Config           *config.Config
	MackerelClient   *mackerel.Client
	ExtensionsClient *extensions.Client
	Agent            *agent.Agent
	Bucket           *metrics.ValuesBucket

	hostID string
}

func getAccountID() (string, error) {
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)
	id, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return "", err
	}
	return *id.Account, nil
}

func getEnvironmentID() (string, error) {
	// XXX: maybe this is not safe
	content, err := ioutil.ReadFile("/proc/sys/kernel/random/boot_id")
	if err != nil {
		return "", err
	}
	bootID := strings.TrimRight(string(content), "\r\n")
	return bootID, nil
}

func (app *App) init() (err error) {
	err = app.MackerelClient.CreateGraphDefs([]*mackerel.GraphDefsParam{
		metrics.CPUGraphDefs,
		metrics.FilesystemGraphDefs,
		metrics.DiskGraphDefs,
		metrics.LoadavgGraphDefs,
		metrics.MemoryGraphDefs,
		metrics.NetworkGraphDefs,
	})
	if err != nil {
		return
	}
	accountID, err := getAccountID()
	if err != nil {
		return
	}
	// AWS Lambda execution environment id
	environmentID, err := getEnvironmentID()
	if err != nil {
		return
	}
	region := os.Getenv("AWS_REGION")
	functionName := os.Getenv("AWS_LAMBDA_FUNCTION_NAME")
	functionArn := fmt.Sprintf("arn:aws:lambda:%s:%s:function:%s", region, accountID, functionName)
	param := &mackerel.CreateHostParam{
		Name: environmentID, // or function name + environment id
		Meta: mackerel.HostMeta{
			AgentName:     "mackerel-lambda-agent",
			AgentVersion:  app.Version,
			AgentRevision: app.Revision,
			Cloud: &mackerel.Cloud{
				Provider: "lambda",
				MetaData: &MetaData{
					FunctionArn: functionArn,
					Version:     os.Getenv("AWS_LAMBDA_FUNCTION_VERSION"),
					MemorySize:  os.Getenv("AWS_LAMBDA_FUNCTION_MEMORY_SIZE"),
				},
			},
		},
	}
	if app.Config.DisplayName != "" {
		param.DisplayName = app.Config.DisplayName
	}
	param.RoleFullnames = app.Config.Roles
	// Everytime create new host and retire on shutdown
	app.hostID, err = app.MackerelClient.CreateHost(param)
	if err != nil {
		return
	}
	return

}

func (app *App) collectMetrics(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			go app.Agent.Collect(app.Bucket)
		case <-ctx.Done():
			return
		}
	}
}

func (app *App) sendMetrics(now int64, values *metrics.Values) (err error) {
	if len(*values) == 0 {
		return
	}
	var metricValues []*mackerel.HostMetricValue
	for name, value := range *values {
		metricValues = append(metricValues, &mackerel.HostMetricValue{
			HostID: app.hostID,
			MetricValue: &mackerel.MetricValue{
				Name:  name,
				Time:  now,
				Value: value,
			},
		})
	}
	err = app.MackerelClient.PostHostMetricValues(metricValues)
	return
}

func (app *App) flushOnce() {
	now := time.Now().Unix()
	// Skip if not one hour passed after last flushing
	if app.Bucket.LastFlushedAt()+60 > now {
		return
	}
	// TODO: Switch reduce fn depending on metrics
	values := libs.MapReduce(app.Bucket.Flush(), math.Max)
	err := app.sendMetrics(now, values)
	if err != nil {
		log.Error("Failed to post metrics", err)
	}
}

func (app *App) flush(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			app.flushOnce()
		case <-ctx.Done():
			app.flushOnce()
			return
		}
	}
}

func (app *App) loop() {
	for {
		event, err := app.ExtensionsClient.Next()
		if err != nil {
			continue
		}
		switch event.EventType {
		case "INVOKE":
			app.flushOnce()
		case "SHUTDOWN":
			return
		default:
			log.Warning("Unknown event type", event)
		}
	}
}

// Run run application
func (app *App) Run(ctx context.Context) {
	_, err := app.ExtensionsClient.Register()
	if err != nil {
		log.Fatal("Failed to register", err)
	}
	err = app.init()
	if err != nil {
		app.ExtensionsClient.InitError(&extensions.ErrorRequest{})
		log.Fatal("Failed to initialize", err)
	}
	ctx, cancel := context.WithCancel(ctx)
	go app.collectMetrics(ctx, 1*time.Second)
	go app.flush(ctx, 60*time.Second)
	app.loop()
	cancel()
	err = app.MackerelClient.RetireHost(app.hostID)
	if err != nil {
		log.Error("Failed to retire", err)
	}
	// TODO: wait for workers
	os.Exit(0)
}
