package util

import "strings"

// GetDomain is used to get domain (scheme+host) from the specified URI
func GetDomain(uri string) string {
	afterDoubleSlashIndex := strings.Index(uri, "//")
	if afterDoubleSlashIndex == -1 {
		afterDoubleSlashIndex = 0
	} else {
		afterDoubleSlashIndex += 2
	}

	pathIndex := afterDoubleSlashIndex + strings.Index(uri[afterDoubleSlashIndex:], "/")
	if pathIndex > afterDoubleSlashIndex {
		return uri[:pathIndex]
	}
	return uri
}
