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
	config       Config
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
	EnableRequestHeader bool
	EnableRequestBody   bool
}

func Default(gop *goprometheus.GoPrometheus) *GinPrometheus {

	c := Config{
		hostname:            hostname(),
		Engine:              gin.Default(),
		MetricName:          "gin_requests_duration",
		MetricDescription:   "The duration of requests",
		TimeFormat:          "2006-01-02T15:04:05.9",
		DurationType:        "ms",
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

func New(gop *goprometheus.GoPrometheus, c Config) *GinPrometheus {

	c.hostname = hostname()
	c.labels = []string{"node", "code", "method", "handler", "host", "url", "time"}

	if c.EnableRequestBody {
		c.labels = append(c.labels, "body")
	}

	if c.EnableRequestHeader {
		c.labels = append(c.labels, "header")
	}

	gp := GinPrometheus{
		config:       c,
		goprometheus: gop,
	}

	gop.AddSummaryVector(c.MetricName, c.MetricDescription, c.labels)

	return &gp
}

func (gp *GinPrometheus) Middleware(c *gin.Context) {

	status := strconv.Itoa(c.Writer.Status())
	url := c.Request.URL.String()

	if value, isPresent := gp.config.IgnorePaths[url]; isPresent {
		if value {
			return
		}
	}

	gp.values = []string{hostname(), status, c.Request.Method, c.HandlerName(), c.Request.Host, url, time.Now().UTC().Format(gp.config.TimeFormat)}

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

	defer gp.timeTrack(time.Now())

	c.Next()
}

func (gp *GinPrometheus) GetEngine() *gin.Engine {
	return gp.config.Engine
}

func (gp *GinPrometheus) timeTrack(start time.Time) {

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

	gp.goprometheus.Vectors.SummaryVectors[gp.config.MetricName].AddMetric(elapsed, gp.values...)

}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	fmt.Print(name)

	return name
}
