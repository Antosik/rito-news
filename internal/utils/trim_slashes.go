package utils

import "strings"

func TrimSlashes(path string) string {
	return strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/")
}
