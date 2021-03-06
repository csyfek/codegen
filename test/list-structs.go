package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/jackmanlabs/errors"
)

func main() {

	pkgPath := "github.com/jackmanlabs/codegen/types"

	pkgs, err := pkger.PackageStructs(pkgPath)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "\t")
	enc.Encode(pkgs)
}
