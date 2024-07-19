package config

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func OtelInitialSetting() {
	// OTLP Exporter設定
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint("localhost:31000"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		fmt.Printf("Failed to create exporter: %v", err)
		return
	}

	// Tracerの設定
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,                       // SchemaURL is the schema URL used to generate the trace ID. Must be set to an absolute URL.
			semconv.ServiceNameKey.String("Portal"), // ServiceNameKey is the key used to identify the service name in a Resource.
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}
