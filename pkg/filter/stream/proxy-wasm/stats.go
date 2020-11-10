package proxywasm

import (
	"sync"

	"github.com/rcrowley/go-metrics"
	metrics2 "mosn.io/mosn/pkg/metrics"
	"mosn.io/mosn/pkg/types"
)

type WasmMetric interface {
	// ID return the id of metric
	ID() uint32
	// Name return the name of metric
	Name() string
	// Value return the value of metric
	Value() int64
	// Record set the value of metric
	Record(int64)
	// Increment add the value to the metric
	Increment(int64) int64
}

type counterWrapper struct {
	id   uint32
	name string
	metrics.Counter
}

func (c *counterWrapper) ID() uint32 {
	return c.id
}

func (c *counterWrapper) Name() string {
	return c.name
}

func (c *counterWrapper) Value() int64 {
	return c.Count()
}

func (c *counterWrapper) Record(value int64) {
	diff := value - c.Count()
	c.Inc(diff)
}

func (c *counterWrapper) Increment(value int64) int64 {
	c.Inc(value)
	return c.Count()
}

type gaugeWrapper struct {
	id   uint32
	name string
	metrics.Gauge
}

func (g *gaugeWrapper) ID() uint32 {
	return g.id
}

func (g *gaugeWrapper) Name() string {
	return g.name
}

func (g *gaugeWrapper) Record(value int64) {
	g.Update(value)
}

func (g *gaugeWrapper) Increment(value int64) int64 {
	newValue := g.Value() + value
	g.Record(newValue)
	return newValue
}

type histogramWrapper struct {
	id   uint32
	name string
	metrics.Histogram
}

func (h *histogramWrapper) ID() uint32 {
	return h.id
}

func (h *histogramWrapper) Name() string {
	return h.name
}

func (h *histogramWrapper) Value() int64 {
	return 0
}

func (h *histogramWrapper) Record(value int64) {
	return
}

func (h *histogramWrapper) Increment(value int64) int64 {
	return 0
}

type metricsManager struct {
	metrics            types.Metrics
	metricIDGenerator  uint32
	metricsMap         map[uint32]WasmMetric
	metricsNameToIDMap map[string]uint32
	mut                sync.RWMutex
}

func newWasmMetricsManager(wasmExtentionName string) *metricsManager {
	return &metricsManager{
		metrics:            metrics2.NewWasmStats(wasmExtentionName),
		metricIDGenerator:  0,
		metricsMap:         make(map[uint32]WasmMetric),
		metricsNameToIDMap: make(map[string]uint32),
	}
}

func (mm *metricsManager) newCounter(name string) uint32 {
	mm.mut.Lock()
	defer mm.mut.Unlock()
	// metric already exists
	if id, ok := mm.metricsNameToIDMap[name]; ok {
		return id
	}
	// otherwise, create new one
	mm.metricIDGenerator++
	counter := &counterWrapper{
		id:      mm.metricIDGenerator,
		name:    name,
		Counter: mm.metrics.Counter(name),
	}
	mm.metricsMap[counter.id] = counter
	mm.metricsNameToIDMap[name] = counter.id
	return counter.id
}

func (mm *metricsManager) newGauge(name string) uint32 {
	mm.mut.Lock()
	defer mm.mut.Unlock()
	// metric already exists
	if id, ok := mm.metricsNameToIDMap[name]; ok {
		return id
	}
	// otherwise, create new one
	mm.metricIDGenerator++
	gauge := &gaugeWrapper{
		id:    mm.metricIDGenerator,
		name:  name,
		Gauge: mm.metrics.Gauge(name),
	}
	mm.metricsMap[gauge.id] = gauge
	return gauge.id
}

func (mm *metricsManager) newHistogram(name string) uint32 {
	mm.mut.Lock()
	defer mm.mut.Unlock()
	// metric already exists
	if id, ok := mm.metricsNameToIDMap[name]; ok {
		return id
	}
	// otherwise, create new one
	mm.metricIDGenerator++
	histogram := &histogramWrapper{
		id:        mm.metricIDGenerator,
		name:      name,
		Histogram: mm.metrics.Histogram(name),
	}
	mm.metricsMap[histogram.id] = histogram
	return histogram.id
}

func (mm *metricsManager) incrementMetric(metricId uint32, offset int64) {
	if mm.metricsMap == nil {
		return
	}
	mm.mut.RLock()
	metric, ok := mm.metricsMap[metricId]
	mm.mut.RUnlock()
	if !ok {
		return
	}
	metric.Increment(offset)
}

func (mm *metricsManager) recordMetric(metricId uint32, value int64) {
	if mm.metricsMap == nil {
		return
	}
	mm.mut.RLock()
	metric, ok := mm.metricsMap[metricId]
	mm.mut.RUnlock()
	if !ok {
		return
	}
	metric.Record(value)
}

func (mm *metricsManager) getMetric(metricId uint32) int64 {
	if mm.metricsMap == nil {
		return 0
	}
	mm.mut.RLock()
	metric, ok := mm.metricsMap[metricId]
	mm.mut.RUnlock()
	if !ok {
		return 0
	}
	return metric.Value()
}
