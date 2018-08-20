package mysql

var templateTestSelectOne string = `
func TestSelect{{.model}}(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}

	id := ""

	_, err = ds.Select{{.model}}(id)
	if err != nil {
		t.Fatal(errs.Stack(err))
	}
}
`

var templateTestSelectOneTx string = `
func TestSelect{{.model}}Tx(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}

	tx, err := ds.Begin()
	if err != nil{
		t.Fatal(errs.Stack(err))
	}
	
	id := ""

	_, err = ds.Select{{.model}}Tx(tx, id)
	if err != nil {
		tx.Rollback()
		t.Fatal(errs.Stack(err))
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(errs.Stack(err))
	}
}`
