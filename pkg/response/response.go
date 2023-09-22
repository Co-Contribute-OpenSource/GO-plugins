package response

import (
	"fmt"
	"net/http"

	"github.com/AccelByte/go-restful-plugins/v4/pkg/logger/event"
	"github.com/emicklei/go-restful/v3"
	"github.com/pkg/errors"
)

// Error is response sent when an error occurs
// Use event ID for error code, register your event ID at:
// https://docs.google.com/spreadsheets/d/1tUB0BSNLyPgeWEtnNzVQkl6Shud-_ErJja2RjIyt1B0/edit?usp=sharing
type Error struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
	ErrorLogMsg  string `json:"-"`
}

const (
	levelInfo  = 3
	levelWarn  = 4
	levelError = 5

	unableToWriteResponse = 20000
)

// Write sends response with specified values
func Write(request *restful.Request, response *restful.Response, httpStatusCode int, serviceType int, eventID int,
	message string, entity interface{}) {
	err := response.WriteHeaderAndJson(httpStatusCode, entity, restful.MIME_JSON)
	if err != nil {
		WriteErrorWithEventID(request, response, http.StatusInternalServerError, serviceType, eventID, errors.WithStack(err),
			&Error{
				ErrorCode:    unableToWriteResponse,
				ErrorMessage: "unable to write response",
				ErrorLogMsg:  fmt.Sprintf("unable to write response: %+v, body: %+v, error: %v", response, entity, err),
			})

		return
	}

	event.Info(request, eventID, serviceType, levelInfo, fmt.Sprintf("%s, response: %+v", message, entity))
}

// WriteError sends error message
func WriteError(request *restful.Request, response *restful.Response, httpStatusCode int, serviceType int,
	eventErr error, errorResponse *Error) {
	WriteErrorWithEventID(request, response, httpStatusCode, serviceType, errorResponse.ErrorCode, eventErr, errorResponse)
}

// WriteErrorWithEventID sends error message with Event ID
func WriteErrorWithEventID(request *restful.Request, response *restful.Response, httpStatusCode int,
	serviceType int, eventID int, eventErr error, errorResponse *Error) {
	err := response.WriteHeaderAndJson(httpStatusCode, errorResponse, restful.MIME_JSON)
	if err != nil {
		err = errors.Wrap(err, "unable to write error response")
		event.Error(request, unableToWriteResponse, serviceType, levelError,
			fmt.Sprintf("%v: %+v: %v", err, errorResponse, eventErr))
		fmt.Printf("%+v\n", err)

		return
	}

	if httpStatusCode >= 500 {
		event.Error(request, eventID, serviceType, levelError,
			fmt.Sprintf("error: %+v: %v", errorResponse, eventErr))
		fmt.Printf("%+v\n", eventErr)

		return
	}

	event.Warn(request, eventID, serviceType, levelWarn,
		fmt.Sprintf("error: %+v: %v", errorResponse, eventErr))
	fmt.Printf("%+v\n", eventErr)
}
