package datadog

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/mocktracer"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// nolint:paralleltest
func TestChildSpan(t *testing.T) {
	assert := assert.New(t)

	mt := mocktracer.Start()
	defer mt.Stop()

	ws := new(restful.WebService)
	ws.Filter(Trace)
	ws.Route(ws.GET("/user/{id}").To(func(request *restful.Request, response *restful.Response) {
		_, ok := tracer.SpanFromContext(request.Request.Context())
		assert.True(ok)
	}))

	container := restful.NewContainer()
	container.Add(ws)

	r := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	container.ServeHTTP(w, r)
}

// nolint:paralleltest
func TestTrace200(t *testing.T) {
	assert := assert.New(t)

	mt := mocktracer.Start()
	defer mt.Stop()

	ws := new(restful.WebService)
	ws.Filter(Trace)
	ws.Route(ws.GET("/user/{id}").Param(restful.PathParameter("id", "user ID")).
		To(func(request *restful.Request, response *restful.Response) {
			_, ok := tracer.SpanFromContext(request.Request.Context())
			assert.True(ok)
			id := request.PathParameter("id")
			ignoreErr(response.Write([]byte(id)))
		}))

	container := restful.NewContainer()
	container.Add(ws)

	r := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	container.ServeHTTP(w, r)
	response := w.Result()
	assert.Equal(response.StatusCode, 200)

	spans := mt.FinishedSpans()

	assert.Len(spans, 1)

	if len(spans) < 1 {
		t.Fatalf("no spans")
	}

	span := spans[0]
	assert.Equal("http.request", span.OperationName())
	assert.Equal(ext.SpanTypeWeb, span.Tag(ext.SpanType))
	assert.Contains(span.Tag(ext.ResourceName), "/user/{id}")
	assert.Equal("200", span.Tag(ext.HTTPCode))
	assert.Equal("GET", span.Tag(ext.HTTPMethod))
	assert.Equal("/user/123", span.Tag(ext.HTTPURL))
}

// nolint:paralleltest
func TestError(t *testing.T) {
	assert := assert.New(t)

	mt := mocktracer.Start()
	defer mt.Stop()

	wantErr := errors.New("oh no")

	ws := new(restful.WebService)
	ws.Filter(Trace)
	ws.Route(ws.GET("/err").To(func(request *restful.Request, response *restful.Response) {
		ignoreErr(response.WriteError(500, wantErr))
	}))

	container := restful.NewContainer()
	container.Add(ws)

	r := httptest.NewRequest("GET", "/err", nil)
	w := httptest.NewRecorder()

	container.ServeHTTP(w, r)
	response := w.Result()
	assert.Equal(response.StatusCode, 500)

	spans := mt.FinishedSpans()
	assert.Len(spans, 1)

	if len(spans) < 1 {
		t.Fatalf("no spans")
	}

	span := spans[0]
	assert.Equal("http.request", span.OperationName())
	assert.Equal("500", span.Tag(ext.HTTPCode))
	assert.Equal(wantErr.Error(), span.Tag(ext.Error).(error).Error())
}

// nolint:paralleltest
func TestGetSpanNotInstrumented(t *testing.T) {
	assert := assert.New(t)

	ws := new(restful.WebService)
	ws.Route(ws.GET("/ping").To(func(request *restful.Request, response *restful.Response) {
		// Assert we don't have a span on the context.
		_, ok := tracer.SpanFromContext(request.Request.Context())
		assert.False(ok)
		ignoreErr(response.Write([]byte("ok")))
	}))

	container := restful.NewContainer()
	container.Add(ws)

	r := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	container.ServeHTTP(w, r)
	response := w.Result()
	assert.Equal(response.StatusCode, 200)
}

// nolint:paralleltest
func TestPropagation(t *testing.T) {
	assert := assert.New(t)

	mt := mocktracer.Start()
	defer mt.Stop()

	r := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	pspan := tracer.StartSpan("test")
	ignoreErr(tracer.Inject(pspan.Context(), tracer.HTTPHeadersCarrier(r.Header)))

	ws := new(restful.WebService)
	ws.Filter(Trace)
	ws.Route(ws.GET("/user/{id}").To(func(request *restful.Request, response *restful.Response) {
		span, ok := tracer.SpanFromContext(request.Request.Context())
		assert.True(ok)
		assert.Equal(span.(mocktracer.Span).ParentID(), pspan.(mocktracer.Span).SpanID())
	}))

	container := restful.NewContainer()
	container.Add(ws)

	container.ServeHTTP(w, r)
}

// nolint:paralleltest
func TestInject(t *testing.T) {
	assert := assert.New(t)

	mt := mocktracer.Start()
	defer mt.Stop()

	ws := new(restful.WebService)
	ws.Filter(Trace)
	ws.Route(ws.GET("/user/{id}").To(func(request *restful.Request, response *restful.Response) {
		outReq := httptest.NewRequest("GET", "/example", nil)
		ignoreErr(Inject(outReq, request))
		_, err := tracer.Extract(tracer.HTTPHeadersCarrier(outReq.Header))
		assert.Nil(err)
	}))

	container := restful.NewContainer()
	container.Add(ws)

	r := httptest.NewRequest("GET", "/user/123", nil)
	w := httptest.NewRecorder()

	container.ServeHTTP(w, r)
}

func ignoreErr(_ ...interface{}) {}
