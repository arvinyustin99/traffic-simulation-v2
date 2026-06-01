package model

type Metrics struct {
	AvgSpeed float64
}

func NewMetrics() *Metrics {
	return &Metrics{}
}

func (s *Simulation) CollectMetrics() {
	s.Metrics.Update(s)
}

func (m *Metrics) Update(sim *Simulation) {

}
