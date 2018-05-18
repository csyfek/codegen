package mysql


var templateTestUpdateOne string = `
func TestUpdate{{.model}}(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}

	x := new({{.modelPackageName}}.{{.model}})

	err = ds.Update{{.model}}(x)
	if err != nil {
		t.Fatal(errors.Stack(err))
	}
}
`

var templateTestUpdateOneTx string = `
func TestUpdate{{.model}}Tx(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}
	
	tx, err := ds.Begin()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}
	
	x := new({{.modelPackageName}}.{{.model}})

	err = ds.Update{{.model}}Tx(tx, x)
	if err != nil {
		t.Fatal(errors.Stack(err))
	}

	err = tx.Rollback()
	if err != nil {
		t.Fatal(errors.Stack(err))
	}
}`

