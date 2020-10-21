package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
)

type rawConfig struct {
	APIBase     string   `yaml:"apibase"`
	APIKey      string   `yaml:"apikey"`
	DisplayName string   `yaml:"display_name"`
	Roles       []string `yaml:"roles"`
}

// Config ...
type Config struct {
	APIBase     *url.URL
	APIKey      string
	DisplayName string
	Roles       []string
}

// TODO: fetch config from S3, EFS, ...
func fetch() ([]byte, error) {
	return nil, nil
}

func parseConfig(data []byte) (*rawConfig, error) {
	var conf struct {
		rawConfig `yaml:",inline"`
	}
	err := yaml.Unmarshal(data, &conf)
	if err != nil {
		return nil, err
	}
	return &conf.rawConfig, nil
}

func parseRoles(value string) []string {
	var roles []string
	for _, v := range strings.Split(value, ",") {
		roles = append(roles, strings.Trim(v, " "))
	}
	return roles
}

// Load ...
func Load() (*Config, error) {
	var conf Config
	// TODO: fetch raw config
	rawConf := &rawConfig{}

	var urlString string
	if rawConf.APIBase == "" {
		urlString = os.Getenv("MACKEREL_API_BASE")
	} else {
		urlString = rawConf.APIBase
	}
	if urlString != "" {
		url, err := url.Parse(urlString)
		if err != nil {
			return nil, fmt.Errorf("Invalid api base")
		}
		conf.APIBase = url
	}

	var apiKey string
	if rawConf.APIKey == "" {
		apiKey = os.Getenv("MACKEREL_API_KEY")
	} else {
		apiKey = rawConf.APIKey
	}
	if apiKey == "" {
		return nil, fmt.Errorf("Please set mackerel api key")
	}
	conf.APIKey = apiKey

	if v, ok := os.LookupEnv("MACKEREL_ROLES"); len(rawConf.Roles) == 0 && ok {
		conf.Roles = parseRoles(v)
	} else {
		conf.Roles = rawConf.Roles
	}

	return &conf, nil
}
