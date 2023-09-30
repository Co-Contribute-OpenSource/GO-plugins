package util

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/sirupsen/logrus"
)

// MinifyJSON is used to compact the json bytes (e.g. removing whitespaces and new lines character)
func MinifyJSON(jsonBytes []byte) string {
	buffer := new(bytes.Buffer)
	if err := json.Compact(buffer, jsonBytes); err != nil {
		logrus.Warnf("Fail to compact json: %v", err)

		// the uncompleted json will throw an error, i.e. {"foo":"bar"
		// so instead, switch with the simple parse (e.g. removing new lines)
		jsonString := string(jsonBytes)
		jsonString = strings.ReplaceAll(jsonString, "\n", "")
		jsonString = strings.ReplaceAll(jsonString, "\r", "")
		return jsonString
	}
	return buffer.String()
}
