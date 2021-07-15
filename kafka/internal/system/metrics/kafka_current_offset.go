package metrics

import (
	"math"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type messageCurrentOffsetNumber struct {
	mu     *sync.Mutex
	metric *prometheus.GaugeVec
	values map[string]int
}

func newMessageCurrentOffsetNumber() *messageCurrentOffsetNumber {
	return &messageCurrentOffsetNumber{
		mu: &sync.Mutex{},
		metric: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "kafka_topic_current_offset",
			Help: "Kafka topic current cursor offset",
		}, []string{labelMessageType}),
		values: map[string]int{},
	}
}

func (m *messageCurrentOffsetNumber) Register(registry *prometheus.Registry) {
	registry.MustRegister(m.metric)
}

func (m messageCurrentOffsetNumber) SetToPrometheus() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.metric.Reset()
	for key, v := range m.values {
		m.metric.With(prometheus.Labels{
			labelMessageType: key,
		}).Set(float64(v))
	}
	m.values = map[string]int{}
}

func (m messageCurrentOffsetNumber) Set(topicName string, offsetCurrent int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.values[topicName] = int(math.Max(float64(m.values[topicName]), float64(offsetCurrent)))
}
