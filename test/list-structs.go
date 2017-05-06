package main

import (
	"encoding/json"
	"github.com/jackmanlabs/codegen/extractor"
	"github.com/jackmanlabs/errors"
	"log"
	"os"
)

func main() {

	pkgPath := "github.com/jackmanlabs/codegen/extractor"

	pkgs, err := extractor.PackageStructs(pkgPath)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(pkgs)
}