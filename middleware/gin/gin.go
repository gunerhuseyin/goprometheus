package ginmiddleware

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gunerhuseyin/goprometheus"
	"os"
	"strconv"
	"time"
)

type GinPrometheus struct {
	goprometheus *goprometheus.GoPrometheus
	config       *Config
	values       []string
}

type Config struct {
	labels              []string
	hostname            string
	Engine              *gin.Engine
	MetricName          string
	MetricDescription   string
	IgnorePaths         map[string]bool
	TimeFormat          string
	DurationType        string
	EnableTime          bool
	EnableRequestHeader bool
	EnableRequestBody   bool
}

func Default(gop *goprometheus.GoPrometheus) *GinPrometheus {

	c := &Config{
		hostname:            hostname(),
		Engine:              gin.Default(),
		MetricName:          "gin_requests_duration",
		MetricDescription:   "The duration of requests",
		DurationType:        "ms",
		EnableTime:          false,
		EnableRequestHeader: false,
		EnableRequestBody:   false,
		labels:              []string{"node", "code", "method", "handler", "host", "url", "time"},
		IgnorePaths: map[string]bool{
			"/metrics": true,
			"/ping":    true,
		},
	}

	gp := GinPrometheus{
		config:       c,
		goprometheus: gop,
	}

	gop.AddSummaryVector(c.MetricName, c.MetricDescription, c.labels)

	return &gp
}

func New(gop *goprometheus.GoPrometheus, c *Config) *GinPrometheus {

	c.hostname = hostname()
	c.labels = []string{"node", "code", "method", "handler", "host", "url"}

	if c.EnableRequestBody {
		c.labels = append(c.labels, "body")
	}

	if c.EnableRequestHeader {
		c.labels = append(c.labels, "header")
	}

	if c.EnableTime {
		c.labels = append(c.labels, "time")
	}

	gp := GinPrometheus{
		config:       c,
		goprometheus: gop,
	}

	gop.AddSummaryVector(c.MetricName, c.MetricDescription, c.labels)

	return &gp
}

func (gp *GinPrometheus) Middleware(c *gin.Context) {

	defer gp.trace(c, time.Now())

	c.Next()
}

func (gp *GinPrometheus) GetEngine() *gin.Engine {
	return gp.config.Engine
}

func (gp *GinPrometheus) trace(c *gin.Context, start time.Time) {

	status := strconv.Itoa(c.Writer.Status())
	url := c.Request.URL.String()

	if value, isPresent := gp.config.IgnorePaths[url]; isPresent {
		if value {
			return
		}
	}

	gp.values = []string{hostname(), status, c.Request.Method, c.HandlerName(), c.Request.Host, url}

	if gp.config.EnableRequestBody {
		out, err := json.Marshal(c.Request.Body)
		if err != nil {
			gp.values = append(gp.values, "")
		} else {
			gp.values = append(gp.values, string(out))
		}
	}

	if gp.config.EnableRequestHeader {
		out, err := json.Marshal(c.Request.Header)
		if err != nil {
			gp.values = append(gp.values, "")
		} else {
			gp.values = append(gp.values, string(out))
		}
	}

	if gp.config.EnableTime {
		gp.values = append(gp.values, time.Now().UTC().Format(gp.config.TimeFormat))
	}

	elapsed := gp.since(start)

	gp.goprometheus.Vectors.SummaryVectors[gp.config.MetricName].AddMetric(elapsed, gp.values...)
}

func (gp *GinPrometheus) since(start time.Time) float64 {

	var elapsed float64

	switch gp.config.DurationType {
	case "m":
		elapsed = float64(time.Since(start).Minutes())
	case "s":
		elapsed = float64(time.Since(start).Seconds())
	case "ns":
		elapsed = float64(time.Since(start).Nanoseconds())
	default:
		elapsed = float64(time.Since(start).Nanoseconds()) / 1000000
	}

	return elapsed

}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Print(name)

	return name
}
