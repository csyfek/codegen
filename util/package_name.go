package util

import "strings"

func PackageName(packagePath string) string {
	chunks := strings.Split(packagePath, "/")
	return chunks[len(chunks)-1]
}
