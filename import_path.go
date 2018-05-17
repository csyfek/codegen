package codegen

import (
	"github.com/jackmanlabs/errors"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func ImportPath(path string) string {

	var err error

	path, err = filepath.Abs(path)
	if err != nil {
		log.Print(errors.Stack(err))
		return ""
	}

	gopath := os.Getenv("GOPATH")
	gopath, err = filepath.Abs(gopath)
	if err != nil {
		log.Print(errors.Stack(err))
		return ""
	}

	if !strings.HasPrefix(path, gopath) {
		return ""
	}

	pkgpath := strings.TrimPrefix(path, gopath+"/src/")

	return pkgpath
}
