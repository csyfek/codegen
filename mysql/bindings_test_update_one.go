package mysql

var templateTestUpdateOne string = `
func TestUpdate{{.model}}(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}

	x := new({{.modelPackageName}}.{{.model}})

	err = ds.Update{{.model}}(x)
	if err != nil {
		t.Fatal(errs.Stack(err))
	}
}
`

var templateTestUpdateOneTx string = `
func TestUpdate{{.model}}Tx(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}
	
	tx, err := ds.Begin()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}
	
	x := new({{.modelPackageName}}.{{.model}})

	err = ds.Update{{.model}}Tx(tx, x)
	if err != nil {
		tx.Rollback()
		t.Fatal(errs.Stack(err))
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(errs.Stack(err))
	}
}`
