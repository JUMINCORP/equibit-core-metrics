package micrometrics

type Metric struct {
	Labels map[string]string
	Name   string
	Value  float64
}

type Exporter interface {
	Setup()
	Export([]Metric) error
}
