package main

import (
	"github.com/jackmanlabs/codegen/mysql"
	"github.com/jackmanlabs/codegen/structfinder"
	"github.com/jackmanlabs/errors"
	"html/template"
	"log"
	"net/http"
)

type handlerGenerateSql struct{}

func (this *handlerGenerateSql) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	var data GenerateSqlData = GenerateSqlData{
		Structs: make([]string, 0),
	}

	if r.Method == "POST" {
		data.Input = r.FormValue("Input")
		data.Struct = r.FormValue("Struct")

		structFinder, err := structfinder.NewStructFinderFromBytes([]byte(data.Input))
		if err != nil {
			log.Print(errors.Stack(err))
		}

		structDatum := structFinder.FindStructs()

		// Set up the struct selection drop-down menu.
		structSelected := false
		var selectedStruct structfinder.StructDefinition
		for _, structData := range structDatum {
			selectOption := SelectOption{
				Name:     structData.Name,
				Selected: data.Struct == structData.Name,
			}

			if selectOption.Selected {
				structSelected = true
				selectedStruct = structData
			}
			data.Structs = append(data.Structs, selectOption)
		}

		if !structSelected && len(data.Structs) > 0 {
			data.Structs[0].Selected = true
		}

		// Make sure that we have a valid struct selected for generation.
		if !structSelected && len(structDatum) > 0 {
			selectedStruct = structDatum[0]
		} else if !structSelected {
			goto PostProcessing
		}

		data.SelectOne = mysql.SelectOne(selectedStruct)
		data.SelectMany = mysql.SelectMany(selectedStruct)
	}

PostProcessing:
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
	Input      string
	Schema     string
	SelectOne  string
	SelectMany string
	Insert     string
	Update     string
	Delete     string
	Errors     string
	Structs    []SelectOption
	Struct     string
}

type SelectOption struct {
	Name     string
	Selected bool
}

var generateSqlHtml string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Generate SQL From Golang</title>
    <style>
        textarea,div {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
<form method="POST">
    <div style="clear:both; width:33%; float:left;">
        <label>Input:</label>
        <br/>
        <textarea style="width:100%;"
                  rows="20"
                  name="Input">{{.Input}}</textarea>
        <br/>
    </div>
    <div style="width:33%; float:left;">
        <label>Struct:</label><br/>
        <select name="Struct"> {{range .Structs}}
            <option {{if
                    .Selected}}
                    selected="selected"
                    {{end}}
                    value="{{.Name}}">{{.Name}}
            </option>
            {{end}} </select>
        <br/>
        <input type="submit"
               value="Submit">
    </div>
</form>
<div style="width:33%; float:left;">
    <label>Errors:</label>
    <br/>
    <textarea style="width:100%;"
              rows="20"
              name="Errors">{{.Errors}}</textarea>
</div>
<div style="clear:both; width:33%; float:left;">
    <label>Create:</label>
    <br/>
    <textarea style="width:100%;"
              rows="20"
              name="Create">{{.Create}}</textarea>
</div>
<div style="width:33%; float:left;">
    <label>SelectOne:</label>
    <br/>
    <textarea style="width:100%;"
              rows="20"
              name="SelectOne">{{.SelectOne}}</textarea>
</div>
<div style="width:33%; float:left;">
    <label>SelectMany:</label>
    <br/>
    <textarea style="width:100%;"
              rows="20"
              name="SelectMany">{{.SelectMany}}</textarea>
</div>
<div style="width:33%; float:left;">
    <label>Insert:</label>
    <br/>
    <textarea style="width:100%;"
              rows="20"
              name="Insert">{{.Insert}}</textarea>
</div>
<div style="width:33%; float:left;">
    <label>Update:</label>
    <br/>
    <textarea style="width:100%;"
              rows="20"
              name="Update">{{.Update}}</textarea>
</div>
<div style="width:33%; float:left;">
    <label>Delete:</label>
    <br/>
    <textarea style="width:100%;"
              rows="20"
              name="Delete">{{.Delete}}</textarea>
</div>
</body>
</html>
`
