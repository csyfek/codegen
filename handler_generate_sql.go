package main

import (
	"fmt"
	"github.com/jackmanlabs/errors"
	"go/ast"
	"go/parser"
	"go/scanner"
	"go/token"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

type handlerGenerateSql struct{}

func (this *handlerGenerateSql) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var data GenerateSqlData

	if r.Method == "POST" {
		data.Input = r.FormValue("Input")

		filename := "/home/jackman/gopath/src/v/codegen/handler_generate_sql.go"

		structFinder, err := NewStructFinderFromFile(filename)
		if err != nil {
			log.Print(errors.Stack(err))
		}

		err = structFinder.FindStructs()
		if err != nil {
			log.Print(errors.Stack(err))
		}

		//data.SelectSingular = buf.String()

		//test, err := json.MarshalIndent(t, "", "\t")
		//data.SelectSingular = string(test)
	}

	t, err := template.New("generateSql").Parse(generateSqlHtml)
	if err != nil {
		writeError(w, errors.Stack(err))
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Print(errors.Stack(err))
	}
}

type GenerateSqlData struct {
	Input          string
	SelectSingular string
	SelectPlural   string
	Insert         string
	Update         string
	Delete         string
	Errors         string
}

var generateSqlHtml string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Generate SQL From Golang</title>
</head>
<body>
<form method="POST">
    <div style="float:left;">
        <label>Input:</label>
        <br/>
        <textarea cols="80" rows="40" name="Input">{{.Input}}</textarea>
        <br/>
    </div>
    <div style="float:left;">
        <input type="submit" value="Submit">
    </div>
</form>
<div style="float:left;">
    <label>Errors:</label>
    <br/>
    <textarea cols="80" rows="40" name="Errors">{{.Errors}}</textarea>
</div>
<div style="clear:both; float:left;">
    <label>SelectSingular:</label>
    <br/>
    <textarea cols="80" rows="40" name="SelectSingular">{{.SelectSingular}}</textarea>
</div>
<div style="float:left;">
    <label>SelectPlural:</label>
    <br/>
    <textarea cols="80" rows="40" name="SelectPlural">{{.SelectPlural}}</textarea>
</div>
<div style="float:left;">
    <label>Insert:</label>
    <br/>
    <textarea cols="80" rows="40" name="Insert">{{.Insert}}</textarea>
</div>
<div style="float:left;">
    <label>Update:</label>
    <br/>
    <textarea cols="80" rows="40" name="Update">{{.Update}}</textarea>
</div>
<div style="float:left;">
    <label>Delete:</label>
    <br/>
    <textarea cols="80" rows="40" name="Delete">{{.Delete}}</textarea>
</div>
</body>
</html>
`

func NewStructFinderFromFile(filename string) (*StructFinder, error) {

	var (
		this *StructFinder = new(StructFinder)
		err  error
	)

	this.Data, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Stack(err)
	}

	this.FSet = token.NewFileSet()
	this.File, err = parser.ParseFile(this.FSet, filename, nil, 0)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return this, nil
}

type StructFinder struct {
	File *ast.File
	Data []byte
	FSet *token.FileSet
}

func (this *StructFinder) FindStructs() error {
	for _, dec := range this.File.Decls {
		switch dec.(type) {
		case *ast.GenDecl:
			genDecl := dec.(*ast.GenDecl)
			if genDecl.Tok == token.TYPE {
				var s scanner.Scanner
				s.Init(this.FSet.File(genDecl.Pos()), this.Data, nil /* no error handler */, scanner.ScanComments)

				var (
					pos token.Pos
					tok token.Token
					lit string
				)

				// fast forward scanner
				for pos < genDecl.TokPos {
					pos, tok, lit = s.Scan()
				}

				// This should yield the struct name
				pos, tok, lit = s.Scan()
				fmt.Printf("%s\t%s\t%q\n", this.FSet.Position(pos), tok, lit)

				// This should yield token.STRUCT
				pos, tok, lit = s.Scan()
				fmt.Printf("%s\t%s\t%q\n", this.FSet.Position(pos), tok, lit)
				if tok != token.STRUCT{
					continue
				}

			}

			//case *ast.FuncDecl:
			//case *ast.BadDecl:
			//default:
		}
		//fmt.Printf("%s\n", raw[dec.Pos()-1:dec.End()])
	}

	return nil
}

type StructDefinition struct {
	Package string
	Members []StructMemberDefinition
}

type StructMemberDefinition struct {
	Name string
	Type string
}
