package log

import (
	"github.com/emicklei/go-restful/v3"
)

const (
	MaskedQueryParamsAttribute    = "MaskedQueryParams"
	MaskedRequestFieldsAttribute  = "MaskedRequestFields"
	MaskedResponseFieldsAttribute = "MaskedResponseFields"
	UserIDAttribute               = "LogUserId"
	ClientIDAttribute             = "LogClientId"
	NamespaceAttribute            = "LogNamespace"
)

// Option contains attribute options for log functionality
type Option struct {
	// Query param that need to masked in url, separated with comma
	MaskedQueryParams string
	// Field that need to masked in request body, separated with comma
	MaskedRequestFields string
	// Field that need to masked in response body, separated with comma
	MaskedResponseFields string
}

// Attributes filter is used to define the log attributes for the endpoint
func Attribute(opt Option) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain){
		if option.MaskedQueryParams != "" {
			req.SetAttribute(MaskedQueryParamsAttribute, option.MaskedQueryParams)
		}
		if option.MaskedRequestFields != "" {
			req.SetAttribute(MaskedRequestFieldsAttribute, option.MaskedRequestFields)
		}
		if option.MaskedResponseFields != "" {
			req.SetAttribute(MaskedResponseFieldsAttribute, option.MaskedResponseFields)
		}
		chain.ProcessFilter(req, resp)
	}
}