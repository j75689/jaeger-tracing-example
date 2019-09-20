package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

var version string

func main() {
	// Configure tracing
	propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	sender, _ := jaeger.NewUDPTransport(os.Getenv("JAEGER_AGENT_HOST"), 0)
	trace, closer := jaeger.NewTracer(
		"openintracing example "+version,
		jaeger.NewConstSampler(true),
		jaeger.NewRemoteReporter(sender, jaeger.ReporterOptions.BufferFlushInterval(1*time.Second)),
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)
	defer closer.Close()
	opentracing.SetGlobalTracer(trace)

	r := gin.Default()
	r.POST("/headers", opengintracing.NewSpan(trace, "get headers"), handlePrintHeaders)
	r.GET("/version", opengintracing.NewSpan(trace, "version"), handleVersion)
	fmt.Println(r.Run(":8080"))
}

func handlePrintHeaders(c *gin.Context) {
	for k, v := range c.Request.Header {
		c.String(http.StatusOK, fmt.Sprintf("%s: %s\n", k, v))
	}
}

func handleVersion(c *gin.Context) {
	c.String(http.StatusOK, fmt.Sprintf("Version: %s", version))
}
