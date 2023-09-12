package jaeger

import (
	"context"
	"fmt"

	"github.com/AccelByte/go-restful-plugins/v4/pkg/trace"
	"github.com/emicklei/go-restful/v3"
)

type contextKeyType string

const (
	spanContextKey = contextKeyType("span")
)

func Filter() restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		traceID := req.HeaderParameter(trace.TraceIDKey)

		span, ctx := StartSpan(req, "Request "+req.Request.Method+" "+req.Request.URL.Path)

		span.SetTag(trace.TraceIDKey, traceID)
		defer Finish(span)

		ctx = context.WithValue(ctx, spanContextKey, span)

		req.Request = req.Request.WithContext(ctx)

		chain.ProcessFilter(req, resp)

		AddLog(span, "Response status code", fmt.Sprintf("%v", resp.StatusCode()))
	}
}