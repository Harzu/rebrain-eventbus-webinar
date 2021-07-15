package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelType        = "type"
	labelMessageType = "message_type"
)

const (
	typeClaimed = "claimed"
	typeSkipped = "skipped"
	typeSuccess = "success"
	typeFailed  = "failed"
)

type eventbusMessageCount struct {
	mu     *sync.Mutex
	metric *prometheus.GaugeVec
	values map[messageCountData]int
}

func newEventbusMessageCount() *eventbusMessageCount {
	return &eventbusMessageCount{
		mu: &sync.Mutex{},
		metric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "eventbus_message_count",
			Help: "count of messages",
		}, []string{labelType, labelMessageType}),
		values: map[messageCountData]int{},
	}
}

func (m *eventbusMessageCount) Register(registry *prometheus.Registry) {
	registry.MustRegister(m.metric)
}

func (m *eventbusMessageCount) SetToPrometheus() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metric.Reset()
	for data, value := range m.values {
		m.metric.With(prometheus.Labels{
			labelType:        data.metricType,
			labelMessageType: data.messageType,
		}).Set(float64(value))
	}
	m.values = map[messageCountData]int{}
}

func (m *eventbusMessageCount) add(data messageCountData) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.values[data]++
}

func (m *eventbusMessageCount) AddClaimed(messageType string) {
	m.add(messageCountData{
		metricType:  typeClaimed,
		messageType: messageType,
	})
}

func (m *eventbusMessageCount) AddSkipped(messageType string) {
	m.add(messageCountData{
		metricType:  typeSkipped,
		messageType: messageType,
	})
}

func (m *eventbusMessageCount) AddSuccess(messageType string) {
	m.add(messageCountData{
		metricType:  typeSuccess,
		messageType: messageType,
	})
}

func (m *eventbusMessageCount) AddFailed(messageType string) {
	m.add(messageCountData{
		metricType:  typeFailed,
		messageType: messageType,
	})
}

type messageCountData struct {
	metricType  string
	messageType string
}
