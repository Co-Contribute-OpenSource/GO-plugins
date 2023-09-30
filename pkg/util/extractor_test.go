package util

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AccelByte/go-jose/jwt"
	"github.com/AccelByte/go-restful-plugins/v4/pkg/logger/event"
	"github.com/AccelByte/iam-go-sdk/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/stretchr/testify/assert"
)

// nolint: dupl,funlen // most part of the test is identical
func TestExtractDefaultWithJWT(t *testing.T) {
	t.Parallel()

	ws := new(restful.WebService)

	var UserID, Namespace, traceID, sessionID string

	var ClientIDs []string

	ws.Filter(event.Log("test", "iam", ExtractDefault))
	ws.Route(
		ws.GET("/namespace/{namespace}/user/{id}").
			Param(restful.PathParameter("namespace", "namespace")).
			Param(restful.PathParameter("id", "user ID")).
			To(func(request *restful.Request, response *restful.Response) {
				request.SetAttribute("JWTClaims", &iam.JWTClaims{
					Namespace: "testNamespace",
					ClientID:  "testClientID",
					Claims: jwt.Claims{
						Subject: "testUserID",
					},
				})

				UserID, ClientIDs, Namespace, traceID, sessionID = ExtractDefault(request)
			}))

	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest(http.MethodGet, "/namespace/abc/user/def", nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	req.Header.Set(traceIDKey, "testTraceID")
	req.Header.Set(sessionIDKey, "testSesssionID")

	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	assert.Equal(t, "testUserID", UserID)
	assert.Equal(t, []string{"testClientID"}, ClientIDs)
	assert.Equal(t, "testNamespace", Namespace)
	assert.Equal(t, "testTraceID", traceID)
	assert.Equal(t, "testSesssionID", sessionID)
}

func TestExtractDefaultWithoutJWT(t *testing.T) {
	t.Parallel()

	ws := new(restful.WebService)

	var UserID, Namespace, traceID, sessionID string

	var ClientIDs []string

	ws.Filter(event.Log("test", "iam", ExtractDefault))
	ws.Route(
		ws.GET("/namespace/{namespace}/user/{id}").
			Param(restful.PathParameter("namespace", "namespace")).
			Param(restful.PathParameter("id", "user ID")).
			To(func(request *restful.Request, response *restful.Response) {
				UserID, ClientIDs, Namespace, traceID, sessionID = ExtractDefault(request)
			}))

	container := restful.NewContainer()
	container.Add(ws)

	req := httptest.NewRequest(http.MethodGet, "/namespace/abc/user/def", nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	req.Header.Set(traceIDKey, "testTraceID")
	req.Header.Set(sessionIDKey, "testSesssionID")

	resp := httptest.NewRecorder()
	container.ServeHTTP(resp, req)

	assert.Equal(t, "", UserID)
	assert.Equal(t, []string{}, ClientIDs)
	assert.Equal(t, "", Namespace)
	assert.Equal(t, "testTraceID", traceID)
	assert.Equal(t, "testSesssionID", sessionID)
}
