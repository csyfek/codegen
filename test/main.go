package main

import (
	"github.com/jackmanlabs/bucket/jlog"
	"github.com/jackmanlabs/codegen/structfinder"
	"github.com/jackmanlabs/errors"
	"log"
)

func main() {
	filename := "/home/jackman/gopath/src/github.com/jackmanlabs/codegen/structfinder/structfinder.go"
	structFinder, err := structfinder.NewStructFinderFromFile(filename)
	if err != nil {
		log.Print(errors.Stack(err))
	}

	defs := structFinder.FindStructs()
	jlog.Log(defs)
}
