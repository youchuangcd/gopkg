package initialize

import (
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"strings"
)

// TracerProvider returns an OpenTelemetry TracerProvider configured to use
// the Jaeger exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
//
// TracerProvider returns an OpenTelemetry TracerProvider configured to use
//
//	@Description:
//	@param endpoint jaeger 接收数据地址
//	@param serviceName
//	@param env
//	@param simplerRatio 采样率 1位百分之百 0.01=百分之一
//	@return *tracesdk.TracerProvider
//	@return error
func TracerProvider(endpoint, serviceName, env string, simplerRatioServiceMap map[string]float64) (*tracesdk.TracerProvider, error) {
	var endpointOption jaeger.EndpointOption
	if strings.HasPrefix(endpoint, "http") {
		// HTTP.
		endpointOption = jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint))
	} else {
		// UDP.
		endpointOption = jaeger.WithAgentEndpoint(jaeger.WithAgentHost(endpoint))
	}
	// Create the Jaeger exporter
	exp, err := jaeger.New(endpointOption)
	if err != nil {
		return nil, err
	}
	var simplerRatio float64
	if ratio, ok := simplerRatioServiceMap[serviceName]; ok {
		simplerRatio = ratio
	} else if defaultRatio, ok2 := simplerRatioServiceMap["default"]; ok2 {
		simplerRatio = defaultRatio
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(simplerRatio)),
		//tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			attribute.String("environment", env),
		)),
	)
	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader))))
	return tp, nil
}
