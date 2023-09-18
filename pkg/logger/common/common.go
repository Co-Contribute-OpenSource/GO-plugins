package common

import (
	"net/http"
	"time"

	"github.com/AccelByte/go-restful-plugins/v4/pkg/logger/log"
	publicsourceip "github.com/AccelByte/public-source-ip"
	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
)

const (
	commonLogFormat = `%s - %s [%s] "%s %s %s" %d %d %d`
)

// Log is a filter that will log incoming request into the Common Log format
func Log(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if log.FullAccessLogEnabled {
		// Notes: If FullAccessLogEnabled is true, show full access log to avoid breaking changes for existing implementation
		log.AccessLog(req, resp, chain)
	} else {
		commonLogFilter(req, resp, chain)
	}
}

// simpleAccessLogFilter will print the access log in simple common log format
func commonLogFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	start := time.Now()
	username := "-"

	if req.Request.URL.User != nil {
		if name := req.Request.URL.User.Username(); name != "" {
			username = name
		}
	}

	chain.ProcessFilter(req, resp)

	duration := time.Since(start)
	logrus.Infof(commonLogFormat,
		publicsourceip.PublicIP(&http.Request{Header: req.Request.Header}),
		username,
		time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		req.Request.Method,
		req.Request.URL.RequestURI(),
		req.Request.Proto,
		resp.StatusCode(),
		resp.ContentLength(),
		duration.Milliseconds(),
	)
}