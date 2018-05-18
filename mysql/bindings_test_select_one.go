package mysql

var templateTestSelectOne string = `
func TestSelect{{.model}}(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}

	id := ""

	_, err = ds.Select{{.model}}(id)
	if err != nil {
		t.Fatal(errors.Stack(err))
	}
}
`

var templateTestSelectOneTx string = `
func TestSelect{{.model}}Tx(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}
	
	tx, err := ds.Begin()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}
	
	id := ""

	_, err = ds.Select{{.model}}Tx(tx, id)
	if err != nil {
		t.Fatal(errors.Stack(err))
	}

	err = tx.Rollback()
	if err != nil {
		t.Fatal(errors.Stack(err))
	}
}`

