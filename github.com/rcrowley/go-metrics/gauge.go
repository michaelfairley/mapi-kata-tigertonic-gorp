package metrics

import "sync/atomic"

// Gauges hold an int64 value that can be set arbitrarily.
//
// This is an interface so as to encourage other structs to implement
// the Gauge API as appropriate.
type Gauge interface {
	Update(int64)
	Value() int64
}

// Get an existing or create and register a new Gauge.
func GetOrRegisterGauge(name string, r Registry) Gauge {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewGauge()).(Gauge)
}

// Create a new Gauge.
func NewGauge() Gauge {
	if UseNilMetrics {
		return NilGauge{}
	}
	return &StandardGauge{0}
}

// Create and register a new Gauge.
func NewRegisteredGauge(name string, r Registry) Gauge {
	c := NewGauge()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// No-op Gauge.
type NilGauge struct{}

// No-op.
func (g NilGauge) Update(v int64) {}

// No-op.
func (g NilGauge) Value() int64 { return 0 }

// The standard implementation of a Gauge uses the sync/atomic package
// to manage a single int64 value.
type StandardGauge struct {
	value int64
}

// Update the gauge's value.
func (g *StandardGauge) Update(v int64) {
	atomic.StoreInt64(&g.value, v)
}

// Return the gauge's current value.
func (g *StandardGauge) Value() int64 {
	return atomic.LoadInt64(&g.value)
}
