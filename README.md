# GoPrometheus

GoPrometheus, GoGin and GoHystrix metrics exporter for Prometheus

## Installation



```bash
go get github.com/gunerhuseyin/goprometheus
```

## Usage

Import packages

```golang
import(
	"github.com/gunerhuseyin/goprometheus"
	ginmiddleware "github.com/gunerhuseyin/goprometheus/middleware/gin"
	hystrixmiddleware "github.com/gunerhuseyin/goprometheus/middleware/hystrix"
)
```
Create new GoPrometheus and GoGin

```
    r := gin.Default()
    gpm := goprometheus.New()
```

Configure gin middleware

```go
	gpGinConfig := &ginmiddleware.Config{
		Engine:              r,  //Gin engine
		TimeFormat:          "2006-01-02T15:04:05",
		MetricName:          "gin_requests_duration",
		DurationType:        "ms",
		EnableRequestHeader: false,
		EnableRequestBody:   false,
		IgnorePaths: map[string]bool{
			"/metrics": true,
			"/ping":    true,
		},
	}

    gpGin := ginmiddleware.New(gpm, gpGinConfig)
```

Configure hystrix middleware

```go
	gpHystrixConfig := &hystrixmiddleware.Config{
		Prefix: "hystrix_circuit_breaker_",
	}

	gpHystrix := hystrixmiddleware.New(gpm, gpHystrixConfig)
```


Use and run GoPrometheus
```
	gpm.UseGin(gpGin)
	gpm.UseHystrix(gpHystrix)
	gpm.Run()
```
Start http server
```
	http.Handle("/", r)	
    _ = http.ListenAndServe(":8080", nil)
```
## Result

Metrics export to `:8080/metrics` 

Sample gin metrics
```
gin_requests_duration_sum{code="200",handler="main.main.func3",host="localhost:8080",method="GET",node="gopoc-deployment-5597c7fdc4-npwzz",time="2019-07-02T12:16:06",url="/test2"} 55.95169900000001
gin_requests_duration_count{code="200",handler="main.main.func3",host="localhost:8080",method="GET",node="gopoc-deployment-5597c7fdc4-npwzz",time="2019-07-02T12:16:06",url="/test2"} 34
```
34 requests in one second, and a total of 55,95169900000001 ms. Average response time 1,6456382059 ms

Sample hystrix metrics
```
hystrix_circuit_breaker_attempts{name="test"} 100
hystrix_circuit_breaker_concurrency_in_use{name="test"} 0
hystrix_circuit_breaker_context_deadline_exceeded{name="test"} 0
hystrix_circuit_breaker_errors{name="test"} 100
hystrix_circuit_breaker_failures{name="test"} 50
hystrix_circuit_breaker_fallback_failures{name="test"} 100
hystrix_circuit_breaker_fallback_successes{name="test"} 0
hystrix_circuit_breaker_rejects{name="test"} 0
hystrix_circuit_breaker_run_duration{name="test"} 186636
hystrix_circuit_breaker_short_circuits{name="test"} 50
hystrix_circuit_breaker_successes{name="test"} 0
hystrix_circuit_breaker_timeouts{name="test"} 0
hystrix_circuit_breaker_total_duration{name="test"} 209263
```

## License
[MIT](https://choosealicense.com/licenses/mit/)