package metrics

type Metrics struct {
	Requests  uint64
	Hits      uint64
	Misses    uint64
	Evictions uint64
}

func (m *Metrics) RecordHit() {
	m.Requests++
	m.Hits++
}

func (m *Metrics) RecordMiss() {
	m.Requests++
	m.Misses++
}

func (m Metrics) HitRate() float64 {
	if m.Requests == 0 {
		return 0
	}

	return float64(m.Hits) / float64(m.Requests)
}
