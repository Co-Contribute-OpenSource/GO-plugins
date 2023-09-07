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
func Trace