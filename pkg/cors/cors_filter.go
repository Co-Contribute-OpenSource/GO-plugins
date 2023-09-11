package cors

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
)

// CrossOriginResourceSharing is used to create a Container Filter that implements CORS.
// Cross-origin resource sharing (CORS) is a mechanism that allows JavaScript on a web page
// to make XMLHttpRequests to another domain, not the domain the JavaScript originated from.
//
// http://en.wikipedia.org/wiki/Cross-origin_resource_sharing
// http://enable-cors.org/server.html
// https://web.dev/cross-origin-resource-sharing
type CrossOriginResourceSharing struct {
	
}