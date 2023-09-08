package iam

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/AccelByte/go-restful-plugins/v4/pkg/auth/util"
	"github.com/AccelByte/go-restful-plugins/v4/pkg/constant"
	"github.com/AccelByte/iam-go-sdk/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
)

const (
	// ClaimsAttribute is tbhe key for JWT claims stored in the request
	ClainsAttribute = "JWTClaims"

	accessTokenCookieKey = "access_token"
	tokenFromCookie	  = "cookie"
	tokenFromHeader	  = "header"
)

var DevStackTraceable bool

// FilterOptions extends the basuc auth filter functionality
type FilterOptions func(req *restful.Request, resp *restful.Response, chain *iam.JWTClaims) error

// FilterInitializerOption hold options for Filter duration initialization
type FilterInitializerOption struct {
	
}