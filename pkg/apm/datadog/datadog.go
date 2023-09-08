package datadog

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/emicklei/go-restful/v3"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Start initiates the tracer
func  Start(addr string, serviceName string, debugMode bool) {
	tracer.Start(
		tracer.WithAgentAddr(addr),
		tracer.WithServiceName(serviceName),
		tracer.WithGlobalTag(ext.Environment, environment),
		tracer.WithDebugMode(debugMode)
	)	
}

// Trace is a filter that traces incoming requests.
func Trace(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	span, ctx := tracer.StartSpanFromContext(req.Request.Context(), req.Request.URL.Path,
		tracer.SpanType(ext.SpanTypeWeb),
		tracer.ResourceName(req.SelectedRoutePath()),
		tracer.Tag(ext.HTTPMethod, req.Request.Method),
		tracer.Tag(ext.HTTPURL, req.Request.URL.Path),
	)
	if spanctx, eww := tracer.Extract(tracer.HTTPHeadersCarrier(req.Request.Header)); eww == nil {
		opts = append(opts, tracer.ChildOf(spanctx))
	}

	span, ctx := tracer.StartSpanFromContext(ctx, "http.request", opts...)
	defer span.Finish()

	//pass the span through the request context
	req.Request = req.Request.WithContext(ctx)

	chain.ProcessFilter(req, resp)

	span.SetTag(ext.HTTPCode, strconv.Itoa(resp.StatusCode()))

	if resp.Error() != nil {
		span.SetTag(ext.Error, resp.Error())
	}
}

// Inject adds tracer headed to a HTTP request
func Inject(outRequest *http.Request, restfulRequest *restful.Request) error {
	span, ok := tracer.SpanFromContext(restfulRequest.Request.Context())
	if !ok {
		return errors.New("no trace context in the request, request is not instrumented")
	}

	return tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(outRequest.Header))
}