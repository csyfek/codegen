package main

import (
	"github.com/jackmanlabs/codegen/structfinder"
	"github.com/jackmanlabs/errors"
	"html/template"
	"log"
	"net/http"
	"database/sql"
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

		data.SelectSingular = generateSelectSingular(selectedStruct)
		data.SelectPlural = generateSelectPlural(selectedStruct)
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
	Input          string
	SelectSingular string
	SelectPlural   string
	Insert         string
	Update         string
	Delete         string
	Errors         string
	Structs        []SelectOption
	Struct         string
}

type SelectOption struct {
	Name     string
	Selected bool
}

func generateSelectSingular(structfinder.StructDefinition) string {

	var ps_StoreDynamic *sql.Stmt

	// This method assumes that the caller specifies the ID, a UUID.
	func StoreDynamic(noun string, id string, data []byte) error {
		db, err := db()
		if err != nil {
		return errors.Stack(err)
		}

		if ps_StoreDynamic == nil {
		q := `INSERT INTO dynamic (id, noun, data) VALUES (?, ?, ?);`

		ps_StoreDynamic, err = db.Prepare(q)
		if err != nil {
		return errors.Stack(err)
		}
		}

		args := []interface{}{
		id,
		noun,
		data,
		}

		_, err = ps_StoreDynamic.Exec(args...)
		if err != nil {
		return errors.Stack(err)
		}

		return nil
	}

	return ""
}
func generateSelectPlural(structfinder.StructDefinition) string {

	return ""
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
        <textarea cols="80"
                  rows="40"
                  name="Input">{{.Input}}</textarea>
        <br/>
    </div>
    <div style="float:left;">
        <select name="Struct"> {{range .Structs}}
            <option {{if
                    .Selected}}
                    selected="selected"
                    {{end}}
                    value="{{.Name}}">{{.Name}}
            </option>
            {{end}} </select>
        <input type="submit"
               value="Submit">
    </div>
</form>
<div style="float:left;">
    <label>Errors:</label>
    <br/>
    <textarea cols="80"
              rows="40"
              name="Errors">{{.Errors}}</textarea>
</div>
<div style="clear:both; float:left;">
    <label>SelectSingular:</label>
    <br/>
    <textarea cols="80"
              rows="40"
              name="SelectSingular">{{.SelectSingular}}</textarea>
</div>
<div style="float:left;">
    <label>SelectPlural:</label>
    <br/>
    <textarea cols="80"
              rows="40"
              name="SelectPlural">{{.SelectPlural}}</textarea>
</div>
<div style="float:left;">
    <label>Insert:</label>
    <br/>
    <textarea cols="80"
              rows="40"
              name="Insert">{{.Insert}}</textarea>
</div>
<div style="float:left;">
    <label>Update:</label>
    <br/>
    <textarea cols="80"
              rows="40"
              name="Update">{{.Update}}</textarea>
</div>
<div style="float:left;">
    <label>Delete:</label>
    <br/>
    <textarea cols="80"
              rows="40"
              name="Delete">{{.Delete}}</textarea>
</div>
</body>
</html>
`
