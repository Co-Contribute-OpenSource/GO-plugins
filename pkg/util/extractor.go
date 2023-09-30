package util

import (
	"github.com/AccelByte/go-restful-plugins/v4/pkg/auth/iam"
	"github.com/emicklei/go-restful/v3"
)

const (
	traceIDKey   = "X-Ab-TraceID"
	sessionIDKey = "X-Ab-SessionID"
)

// ExtractDefault is default function for extracting attribute for filter event logger
func ExtractDefault(req *restful.Request) (userID string, clientID []string,
	namespace string, traceID string, sessionID string) {
	traceID = req.HeaderParameter(traceIDKey)
	sessionID = req.HeaderParameter(sessionIDKey)

	claims := iam.RetrieveJWTClaims(req)
	if claims != nil {
		return claims.Subject, []string{claims.ClientID}, claims.Namespace, traceID, sessionID
	}

	return "", []string{}, "", traceID, sessionID
}
