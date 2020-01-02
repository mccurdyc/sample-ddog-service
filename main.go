package main

import (
	"fmt"
	"log"
	"os"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/gin-gonic/gin"
)

func main() {
	// https://godoc.org/gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer
	tracer.Start(
		tracer.WithAnalytics(true),
		tracer.WithServiceName("go-sample"),
		tracer.WithAgentAddr(fmt.Sprintf("%s:%s", os.Getenv("DD_AGENT_HOST"), os.Getenv("DD_AGENT_PORT"))),
		tracer.WithRuntimeMetrics(),
	)
	defer tracer.Stop()

	// Create a gin.Engine
	r := gin.New()

	// Use the tracer middleware with your desired service name.
	r.Use(gintrace.Middleware("go-sample"))

	// Continue using the router as normal.
	r.GET("/hello", func(c *gin.Context) {
		// https://docs.datadoghq.com/tracing/advanced/connect_logs_and_traces/?tab=go#manual-trace-id-injection
		span := tracer.StartSpan("web.request", tracer.ResourceName("/hello"))
		defer span.Finish()

		// Retrieve Trace ID and Span ID
		traceID := span.Context().TraceID()
		spanID := span.Context().SpanID()

		// Append them to log messages as fields:
		log.Printf("my log message dd.trace_id=%d dd.span_id=%d", traceID, spanID)
		c.String(404, "Hello World!")
	})

	r.Run(":8080")
}
