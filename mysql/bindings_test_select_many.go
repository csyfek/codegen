package mysql


var templateTestSelectMany string = `
func TestSelect{{.models}}(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}

	_, err = ds.Select{{.models}}()
	if err != nil {
		t.Fatal(errors.Stack(err))
	}
}
`

var templateTestSelectManyTx string = `
func TestSelect{{.models}}Tx(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}
	
	tx, err := ds.Begin()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}

	_, err = ds.Select{{.models}}Tx(tx)
	if err != nil {
		t.Fatal(errors.Stack(err))
	}

	err = tx.Rollback()
	if err != nil {
		t.Fatal(errors.Stack(err))
	}
}`

