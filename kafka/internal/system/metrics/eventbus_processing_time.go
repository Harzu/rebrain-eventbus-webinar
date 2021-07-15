package metrics

import (
	"math"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const labelMetricName = "metric_name"

type eventbusMessageProcessingTime struct {
	mu     *sync.Mutex
	metric *prometheus.GaugeVec
	values map[string]timeRecord
}

func newEventbusMessageProcessingTime() *eventbusMessageProcessingTime {
	return &eventbusMessageProcessingTime{
		mu: &sync.Mutex{},
		metric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "eventbus_message_processing_time",
			Help: "eventbus message processing time",
		}, []string{labelMessageType, labelMetricName}),
		values: map[string]timeRecord{},
	}
}

func (m *eventbusMessageProcessingTime) Register(registry *prometheus.Registry) {
	registry.MustRegister(m.metric)
}

func (m *eventbusMessageProcessingTime) SetToPrometheus() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metric.Reset()
	for key, v := range m.values {
		m.metric.With(prometheus.Labels{labelMessageType: key, labelMetricName: "max"}).Set(v.Max)
		m.metric.With(prometheus.Labels{labelMessageType: key, labelMetricName: "sum"}).Set(v.Sum)
		m.metric.With(prometheus.Labels{labelMessageType: key, labelMetricName: "amount"}).Set(float64(v.Amount))
	}
	m.values = map[string]timeRecord{}
}

func (m *eventbusMessageProcessingTime) Add(messageType string, duration float64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var record timeRecord
	if existRecord, ok := m.values[messageType]; ok {
		record = existRecord
	}

	record.Add(duration)
	m.values[messageType] = record
}

type timeRecord struct {
	Max    float64
	Sum    float64
	Amount int
}

func (d *timeRecord) Add(duration float64) {
	d.Max = math.Max(d.Max, duration)
	d.Sum += duration
	d.Amount++
}
