package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type promMetric interface {
	// SetToPrometheus is expected to bulk-update specific metrics.
	SetToPrometheus()
}

type Client struct {
	registry *prometheus.Registry

	EventbusMessagesProcessingTime *eventbusMessageProcessingTime
	EventbusMessageCount           *eventbusMessageCount
}

func New() *Client {
	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(prometheus.NewGoCollector())

	client := &Client{
		registry: registry,

		EventbusMessagesProcessingTime: newEventbusMessageProcessingTime(),
		EventbusMessageCount:           newEventbusMessageCount(),
	}

	client.EventbusMessagesProcessingTime.Register(registry)
	client.EventbusMessageCount.Register(registry)

	return client
}

func (c *Client) Metrics() []promMetric {
	return []promMetric{
		c.EventbusMessageCount,
		c.EventbusMessagesProcessingTime,
	}
}

func (c *Client) Handler() http.Handler {
	return http.HandlerFunc(func(rsp http.ResponseWriter, req *http.Request) {
		for _, m := range c.Metrics() {
			m.SetToPrometheus()
		}
		promhttp.HandlerFor(c.registry, promhttp.HandlerOpts{}).ServeHTTP(rsp, req)
	})
}
