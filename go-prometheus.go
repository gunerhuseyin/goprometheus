package goprometheus

import (
	metricCollector "github.com/afex/hystrix-go/hystrix/metric_collector"
	"github.com/gin-gonic/gin"
	p "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

type IHystrixPrometheus interface {
	Middleware(name string) metricCollector.MetricCollector
}

type IGinPrometheus interface {
	Middleware(c *gin.Context)
	GetEngine()     *gin.Engine
}

type IGoPrometheus interface {
	AddSummaryVector(name, desc string, labels []string)
	AddCounterVector(name, desc string, labels []string)
	AddGaugeVector(name, desc string, labels []string)
	AddHistogramVector(name, desc string, labels []string)
}

type GoPrometheus struct {
	Vectors *Vectors
	gin     IGinPrometheus
	hystrix IHystrixPrometheus
}

type Vectors struct {
	SummaryVectors   map[string]*SummaryVector
	CounterVectors   map[string]*CounterVector
	GaugeVectors     map[string]*GaugeVector
	HistogramVectors map[string]*HistogramVector
}

type SummaryVector struct {
	Vector *p.SummaryVec
	Labels []string
}

type CounterVector struct {
	Vector *p.CounterVec
	Labels []string
}

type GaugeVector struct {
	Vector *p.GaugeVec
	Labels []string
}

type HistogramVector struct {
	Vector *p.HistogramVec
	Labels []string
}

func New() *GoPrometheus {

	vectors := &Vectors{
		SummaryVectors:   make(map[string]*SummaryVector),
		CounterVectors:   make(map[string]*CounterVector),
		GaugeVectors:     make(map[string]*GaugeVector),
		HistogramVectors: make(map[string]*HistogramVector),
	}

	return &GoPrometheus{
		Vectors: vectors,
	}
}

func (gp *GoPrometheus) UseGin(gin IGinPrometheus) {
	gp.gin = gin
}

func (gp *GoPrometheus) UseHystrix(hystrix IHystrixPrometheus) {
	gp.hystrix = hystrix
}

func (gp *GoPrometheus) AddSummaryVector(name, desc string, labels []string) {
	if _, isPresent := gp.Vectors.SummaryVectors[name]; isPresent {
		return
	}

	gp.Vectors.SummaryVectors[name] = &SummaryVector{
		Vector: promauto.NewSummaryVec(p.SummaryOpts{
			Name: name,
			Help: desc,
		}, labels),
		Labels: labels,
	}
}

func (gp *GoPrometheus) AddCounterVector(name, desc string, labels []string) {
	if _, isPresent := gp.Vectors.CounterVectors[name]; isPresent {
		return
	}

	gp.Vectors.CounterVectors[name] = &CounterVector{
		Vector: promauto.NewCounterVec(p.CounterOpts{
			Name: name,
			Help: desc,
		}, labels),
		Labels: labels,
	}
}

func (gp *GoPrometheus) AddGaugeVector(name, desc string, labels []string) {
	if _, isPresent := gp.Vectors.GaugeVectors[name]; isPresent {
		return
	}

	gp.Vectors.GaugeVectors[name] = &GaugeVector{
		Vector: promauto.NewGaugeVec(p.GaugeOpts{
			Name: name,
			Help: desc,
		}, labels),
		Labels: labels,
	}
}

func (gp *GoPrometheus) AddHistogramVector(name, desc string, labels []string) {
	if _, isPresent := gp.Vectors.HistogramVectors[name]; isPresent {
		return
	}

	gp.Vectors.HistogramVectors[name] = &HistogramVector{
		Vector: promauto.NewHistogramVec(p.HistogramOpts{
			Name: name,
			Help: desc,
		}, labels),
		Labels: labels,
	}
}

func (v *SummaryVector) AddMetric(value float64, labelValues ...string) {
	v.Vector.WithLabelValues(labelValues...).Observe(value)
}

func (v *CounterVector) AddMetric(value float64, labelValues ...string) {
	v.Vector.WithLabelValues(labelValues...).Add(value)
}

func (v *GaugeVector) AddMetric(value float64, labelValues ...string) {
	v.Vector.WithLabelValues(labelValues...).Add(value)
}

func (v *HistogramVector) AddMetric(value float64, labelValues ...string) {
	v.Vector.WithLabelValues(labelValues...).Observe(value)
}

func (gp *GoPrometheus) Run() {

	http.Handle("/metrics", promhttp.Handler())

	if gp.hystrix != nil {
		metricCollector.Registry.Register(gp.hystrix.Middleware)
	}
	if gp.gin != nil {
		gp.gin.GetEngine().Use(gp.gin.Middleware)
	}
}
