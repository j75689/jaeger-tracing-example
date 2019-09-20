package main

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/opengintracing"
	"github.com/gin-gonic/gin"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/zipkin"
)

func main() {
	// Configure tracing
	propagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	trace, closer := jaeger.NewTracer(
		"openintracing example v1",
		jaeger.NewConstSampler(true),
		jaeger.NewNullReporter(),
		jaeger.TracerOptions.Injector(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.Extractor(opentracing.HTTPHeaders, propagator),
		jaeger.TracerOptions.ZipkinSharedRPCSpan(true),
	)
	defer closer.Close()
	opentracing.SetGlobalTracer(trace)

	r := gin.Default()
	r.POST("/headers", opengintracing.NewSpan(trace, "get headers"), handlePrintHeaders)

	fmt.Println(r.Run(":8080"))
}

func handlePrintHeaders(c *gin.Context) {
	for k, v := range c.Request.Header {
		c.String(http.StatusOK, fmt.Sprintf("%s: %s\n", k, v))
	}
}
