package rester

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
)

func GetCollection(def *codegen.Model) (string, string) {

	resourceName := resource(def.Name)
	models := codegen.Plural(def.Name)

	register := fmt.Sprintf(`r.Path("/%s").Methods("GET").Handler(ErrFilter(handleGet%s))`, resourceName, models)

	b := bytes.NewBuffer(nil) // The handler.

	fmt.Fprintf(b, "func handleGet%s(w http.ResponseWriter, r *http.Request) error {\n\n", models)

	fmt.Fprint(b, "var vars map[string]string = mux.Vars(r)\n")

	fmt.Fprint(b, "// FILTER DETECTION\n\n")

	fmt.Fprintf(b, "var filter filters.%s\n\n", def.Name)
	for _, member := range def.Members {
		if member.IsNumeric() {

			fmt.Fprintf(b, "{ // %s\n\n", member.GoName)

			// Min
			fmt.Fprintf(b, "if min_, ok := vars[\"%s_min\"]; ok {\n", member.JsonName)
			fmt.Fprint(b, "\tif min, err := strconv.ParseFloat(min_, 64); err == nil {\n")
			fmt.Fprintf(b, "\tfilter.%s_Min = &%s(min)", member.GoName, member.GoType)
			fmt.Fprint(b, "\t}\n")
			fmt.Fprint(b, "}\n\n")

			// Max
			fmt.Fprintf(b, "if max_, ok := vars[\"%s_max\"]; ok {\n", member.JsonName)
			fmt.Fprint(b, "\tif max, err := strconv.ParseFloat(max_, 64); err == nil {\n")
			fmt.Fprintf(b, "\tfilter.%s_Max =  &%s(max)", member.GoName, member.GoType)
			fmt.Fprint(b, "\t}\n")
			fmt.Fprint(b, "}\n\n")

			// Exact
			fmt.Fprintf(b, "if v_, ok := vars[\"%s\"]; ok {\n", member.JsonName)
			fmt.Fprint(b, "\tif v, err := strconv.ParseFloat(v_, 64); err == nil {\n")
			fmt.Fprintf(b, "\tfilter.%s = &%s(v)", member.GoName, member.GoType)
			fmt.Fprint(b, "\t}\n")
			fmt.Fprint(b, "}\n\n")

			fmt.Fprint(b, "}\n\n")

		} else if member.GoType == "bool" {

		} else if member.GoType == "time.Time" {

		} else if member.GoType == "time.Duration" {

		} else if member.GoType == "string" {
			fmt.Fprintf(b, "if v, ok := vars[\"%s\"]; ok {\n", member.JsonName)
			fmt.Fprintf(b, "\tfilter.%s = &v", member.GoName)
			fmt.Fprint(b, "}\n\n")
		}
	}

	fmt.Fprintf(b, "z, err := control.Get%s(filter)", models)
	fmt.Fprint(b, `
	if err != nil {
		return errors.Stack(err)
	}

	err = serialize(w,r,z)
	if err != nil{
		return errors.Stack(err)
	}

	return nil
}
	`)

	return register, b.String()
}
