package main

import (
	"flag"
	"fmt"
	"github.com/jackmanlabs/codegen/pg"
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

	stat, err := os.Stat(*path)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}
	if stat.IsDir() {
		log.Println("You must specify a Go file to analyze for this tool to work.")
		log.Println("Not a directory. Don't specify a directory.")
		flag.Usage()
		os.Exit(1)
	}

	structFinder, err := structfinder.NewStructFinderFromFile(*path)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	defs := structFinder.FindStructs()

	for _, def := range defs {
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(pg.Create(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(pg.SelectSingular(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(pg.SelectPlural(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(pg.Update(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(pg.Insert(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
		fmt.Println(pg.Delete(def))
		fmt.Println("/*============================================================================*/")
		fmt.Println()
	}
}
