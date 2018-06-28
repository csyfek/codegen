package mysql

var templateTestInsertOne string = `
func TestInsert{{.model}}(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}

	x := new({{.modelPackageName}}.{{.model}})

	err = ds.Insert{{.model}}(x)
	if err != nil {
		t.Fatal(errors.Stack(err))
	}
}
`

var templateTestInsertOneTx string = `
func TestInsert{{.model}}Tx(t *testing.T) {

	ds, err := New()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}
	
	tx, err := ds.Begin()
	if err != nil{
		t.Fatal(errors.Stack(err))
	}
	
	x := new({{.modelPackageName}}.{{.model}})

	err = ds.Insert{{.model}}Tx(tx, x)
	if err != nil {
		tx.Rollback()
		t.Fatal(errors.Stack(err))
	}

	err = tx.Commit()
	if err != nil {
		t.Fatal(errors.Stack(err))
	}
}`

