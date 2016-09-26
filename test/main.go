package main

import (
	"fmt"
	"github.com/jackmanlabs/codegen/generate_sql"
	"github.com/jackmanlabs/codegen/structfinder"
	"github.com/jackmanlabs/errors"
	"log"
)

func main() {
	filename := "/home/jackman/gopath/src/github.com/jackmanlabs/v/types/world.go"
	structFinder, err := structfinder.NewStructFinderFromFile(filename)
	if err != nil {
		log.Print(errors.Stack(err))
	}

	defs := structFinder.FindStructs()

	for _, def := range defs {
		s := generate_sql.SelectPlural(def)
		fmt.Println("==================================================================")
		fmt.Println(s)
		fmt.Println("==================================================================")
	}
}
