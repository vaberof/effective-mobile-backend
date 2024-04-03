package apicar

import (
	"github.com/go-resty/resty/v2"
	"github.com/vaberof/effective-mobile-backend/pkg/logging/logs"
	"log/slog"
	"time"
)

type HttpClientConfig struct {
	Host    string        `yaml:"host"`
	Timeout time.Duration `yaml:"timeout"`
}

type HttpClient struct {
	client *resty.Client
	host   string
	logger *slog.Logger
}

func NewHttpClient(config *HttpClientConfig, logsBuilder *logs.Logs) *HttpClient {
	return &HttpClient{
		host:   config.Host,
		client: resty.New().SetTimeout(config.Timeout),
		logger: logsBuilder.WithName("http-client.api.car"),
	}
}
