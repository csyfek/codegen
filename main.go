package main

import (
	"flag"
	"fmt"
	"github.com/jackmanlabs/codegen/generate_mysql"
	"github.com/jackmanlabs/codegen/structfinder"
	"github.com/jackmanlabs/errors"
	"log"
	"os"
)

func main() {
	var path *string = flag.String("file", "", "The file that you want to analyze.")
	flag.Parse()

	if *path == "" {
		log.Println("You must specify a Go file to analyze for this tool to work.")
		flag.Usage()
		os.Exit(1)
	}

	structFinder, err := structfinder.NewStructFinderFromFile(*path)
	if err != nil {
		log.Print(errors.Stack(err))
	}

	defs := structFinder.FindStructs()

	for _, def := range defs {
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(generate_mysql.Create(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(generate_mysql.SelectSingular(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(generate_mysql.SelectPlural(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(generate_mysql.Update(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(generate_mysql.Insert(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(generate_mysql.Delete(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
	}
}
