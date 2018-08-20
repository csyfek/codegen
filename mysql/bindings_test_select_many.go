package mysql

var templateTestSelectMany string = `
func TestSelect{{.models}}(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}

	_, err = ds.Select{{.models}}()
	if err != nil {
		t.Fatal(errs.Stack(err))
	}
}
`

var templateTestSelectManyTx string = `
func TestSelect{{.models}}Tx(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}
	
	tx, err := ds.Begin()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}

	_, err = ds.Select{{.models}}Tx(tx)
	if err != nil {
		tx.Rollback()
		t.Fatal(errs.Stack(err))
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(errs.Stack(err))
	}
}`
