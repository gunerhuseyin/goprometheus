package hystrixmiddleware

import (
	"github.com/afex/hystrix-go/hystrix/metric_collector"
	gp "github.com/gunerhuseyin/goprometheus"
)

type Config struct {
	Prefix string
}

type HystrixPrometheus struct {
	config                  *Config
	goprometheus            *gp.GoPrometheus
	name                    string
	circuitState            string
	attempts                string
	errors                  string
	successes               string
	failures                string
	rejects                 string
	shortCircuits           string
	timeouts                string
	fallbackSuccesses       string
	fallbackFailures        string
	contextCanceled         string
	contextDeadlineExceeded string
	totalDuration           string
	runDuration             string
	concurrencyInUse        string
}

const (
	CircuitState            = "circuit_state"
	Attempts                = "attempts"
	Errors                  = "errors"
	Successes               = "successes"
	Failures                = "failures"
	Rejects                 = "rejects"
	ShortCircuits           = "short_circuits"
	Timeouts                = "timeouts"
	FallbackSuccesses       = "fallback_successes"
	FallbackFailures        = "fallback_failures"
	ContextCanceled         = "context_canceled"
	ContextDeadlineExceeded = "context_deadline_exceeded"
	TotalDuration           = "total_duration"
	RunDuration             = "run_duration"
	ConcurrencyInUse        = "concurrency_in_use"
)

func New(gp *gp.GoPrometheus, config *Config) *HystrixPrometheus {

	gp.AddCounterVector(config.Prefix+Attempts, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+Errors, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+Successes, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+Failures, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+Rejects, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+ShortCircuits, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+Timeouts, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+FallbackSuccesses, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+FallbackFailures, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+ContextCanceled, "circuit open", []string{"name"})
	gp.AddCounterVector(config.Prefix+ContextDeadlineExceeded, "circuit open", []string{"name"})
	gp.AddGaugeVector(config.Prefix+TotalDuration, "circuit open", []string{"name"})
	gp.AddGaugeVector(config.Prefix+RunDuration, "circuit open", []string{"name"})
	gp.AddGaugeVector(config.Prefix+ConcurrencyInUse, "circuit open", []string{"name"})

	return &HystrixPrometheus{
		goprometheus: gp,
		config:       config,
	}
}

func Default(gp *gp.GoPrometheus) *HystrixPrometheus {

	config := &Config{
		Prefix: "hystrix_circuit_breaker_",
	}

	gp.AddCounterVector(config.Prefix+Attempts, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+Errors, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+Successes, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+Failures, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+Rejects, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+ShortCircuits, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+Timeouts, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+FallbackSuccesses, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+FallbackFailures, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+ContextCanceled, "", []string{"name"})
	gp.AddCounterVector(config.Prefix+ContextDeadlineExceeded, "", []string{"name"})
	gp.AddGaugeVector(config.Prefix+TotalDuration, "", []string{"name"})
	gp.AddGaugeVector(config.Prefix+RunDuration, "", []string{"name"})
	gp.AddGaugeVector(config.Prefix+ConcurrencyInUse, "", []string{"name"})

	return &HystrixPrometheus{
		goprometheus: gp,
		config:       config,
	}
}

func (hp *HystrixPrometheus) Middleware(name string) metricCollector.MetricCollector {
	return &HystrixPrometheus{
		goprometheus:            hp.goprometheus,
		name:                    name,
		circuitState:            hp.config.Prefix + CircuitState,
		attempts:                hp.config.Prefix + Attempts,
		errors:                  hp.config.Prefix + Errors,
		successes:               hp.config.Prefix + Successes,
		failures:                hp.config.Prefix + Failures,
		rejects:                 hp.config.Prefix + Rejects,
		shortCircuits:           hp.config.Prefix + ShortCircuits,
		timeouts:                hp.config.Prefix + Timeouts,
		fallbackSuccesses:       hp.config.Prefix + FallbackSuccesses,
		fallbackFailures:        hp.config.Prefix + FallbackFailures,
		contextCanceled:         hp.config.Prefix + ContextCanceled,
		contextDeadlineExceeded: hp.config.Prefix + ContextDeadlineExceeded,
		totalDuration:           hp.config.Prefix + TotalDuration,
		runDuration:             hp.config.Prefix + RunDuration,
		concurrencyInUse:        hp.config.Prefix + ConcurrencyInUse,
	}
}

func (hp *HystrixPrometheus) Update(r metricCollector.MetricResult) {

	hp.goprometheus.Vectors.CounterVectors[hp.attempts].AddMetric(r.Attempts, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.errors].AddMetric(r.Errors, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.successes].AddMetric(r.Successes, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.failures].AddMetric(r.Failures, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.rejects].AddMetric(r.Rejects, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.shortCircuits].AddMetric(r.ShortCircuits, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.timeouts].AddMetric(r.Timeouts, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.fallbackSuccesses].AddMetric(r.FallbackSuccesses, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.fallbackFailures].AddMetric(r.FallbackFailures, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.fallbackSuccesses].AddMetric(r.ContextCanceled, hp.name)
	hp.goprometheus.Vectors.CounterVectors[hp.contextDeadlineExceeded].AddMetric(r.ContextDeadlineExceeded, hp.name)
	hp.goprometheus.Vectors.GaugeVectors[hp.totalDuration].AddMetric(float64(r.TotalDuration.Nanoseconds()/1000), hp.name)
	hp.goprometheus.Vectors.GaugeVectors[hp.runDuration].AddMetric(float64(r.RunDuration.Nanoseconds()/1000), hp.name)
	hp.goprometheus.Vectors.GaugeVectors[hp.concurrencyInUse].AddMetric(float64(100*r.ConcurrencyInUse), hp.name)

}

func (hpm *HystrixPrometheus) Reset() {}
