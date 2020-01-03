package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var ids = map[int]uuid.UUID{
	0: uuid.New(),
	1: uuid.New(),
	2: uuid.New(),
	3: uuid.New(),
	4: uuid.New(),
	5: uuid.New(),
}

func main() {
	f, err := os.OpenFile("/var/log.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		logrus.Fatalf("failed to open log file")
	}
	logrus.SetOutput(f)

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
	gin.DefaultWriter = ioutil.Discard

	// Use the tracer middleware with your desired service name.
	r.Use(gintrace.Middleware("go-sample"))

	// https://docs.datadoghq.com/logs/log_collection/go/#configure-your-logger

	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Continue using the router as normal.
	r.GET("/hello", func(c *gin.Context) {
		// https://docs.datadoghq.com/tracing/advanced/connect_logs_and_traces/?tab=go#manual-trace-id-injection
		span := tracer.StartSpan("web.request", tracer.ResourceName("/hello"))
		defer span.Finish()

		// Retrieve Trace ID and Span ID
		traceID := span.Context().TraceID()
		spanID := span.Context().SpanID()

		rnd := rand.Intn(5)
		span.SetTag("sku_uuid", ids[rnd])

		// https://docs.datadoghq.com/tracing/advanced/connect_logs_and_traces/?tab=go
		// for now, you manually have to add these fields in Go for them to
		// be tied into DD log collection
		logrus.WithFields(logrus.Fields{"sku_uuid": ids[rnd], "dd.trace_id": traceID, "dd.span_id": spanID}).Infof("my log message")

		ctx := tracer.ContextWithSpan(context.Background(), span)
		wait(ctx, rnd, time.Second*2)

		c.String(200, "Hello World!")
	})

	r.Run(":8080")
}

func wait(ctx context.Context, fanout int, dur time.Duration) {
	span, ctx := tracer.StartSpanFromContext(ctx, "wait", tracer.ResourceName("wait"))
	defer span.Finish()

	for i := 0; i < fanout; i++ {
		innerSpan := tracer.StartSpan("innerwait", tracer.ResourceName("innerwait"), tracer.ChildOf(span.Context()))
		time.Sleep(dur)
		innerSpan.Finish()
	}
}
