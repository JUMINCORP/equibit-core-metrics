package micrometrics

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	mutex            = &sync.Mutex{}
	formattedMetrics = make([]string, 0)
)

type PrometheusExporter struct {
	address string
}

func NewPrometheusExporter(address string) *PrometheusExporter {
	p := new(PrometheusExporter)
	p.address = address
	return p
}

func httpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		for _, metric := range formattedMetrics {
			w.Write([]byte(metric))
		}
	}
}

func (p *PrometheusExporter) setup() {
	http.Handle("/metrics", httpHandler())
	log.Fatal(http.ListenAndServe(p.address, nil))
}

func formatMetric(m Metric) string {
	var sb strings.Builder

	sb.WriteString(m.Name)

	sortedLabels := make([]string, len(m.Labels))
	i := 0
	for k := range m.Labels {
		sortedLabels[i] = k
		i++
	}
	sort.Strings(sortedLabels)

	sb.WriteRune('{')
	for i := range sortedLabels {
		if i != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(sortedLabels[i])
		sb.WriteString("=\"")
		sb.WriteString(m.Labels[sortedLabels[i]])
		sb.WriteRune('"')
	}
	sb.WriteString("} ")

	sb.WriteString(strconv.FormatFloat(m.Value, 'f', -1, 64))

	sb.WriteRune('\n')

	return sb.String()
}

func (p *PrometheusExporter) export(metrics []Metric) error {
	var err error

	mutex.Lock()
	defer mutex.Unlock()

	formattedMetrics = make([]string, len(metrics))

	for i, m := range metrics {
		//formattedLabels := formatLabels(m.labels)
		formattedMetrics[i] = formatMetric(m)
	}
	return err
}
