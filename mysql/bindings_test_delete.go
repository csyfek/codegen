package mysql

var templateTestDelete string = `
func TestDelete{{.model}}(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}

	id := ""

	err = ds.Delete{{.model}}(id)
	if err != nil {
		t.Fatal(errs.Stack(err))
	}
}
`

var templateTestDeleteTx string = `
func TestDelete{{.model}}Tx(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}
	
	tx, err := ds.Begin()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}
	
	id := ""

	err = ds.Delete{{.model}}Tx(tx, id)
	if err != nil {
		tx.Rollback()
		t.Fatal(errs.Stack(err))
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(errs.Stack(err))
	}
}`
