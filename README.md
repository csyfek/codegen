# codegen
Multi-purpose code generator for JSON, YAML, SQL, and Go.

Those who have worked with me have probably seen code generated by this tool or
one of its ancestors. Its ancestors were typically a small web-app with a basic
HTML (no JavaScript) interface for various kinds of transformations. My intent
with this tool is to make it much more capable by inspecting entire packages
during the data gathering stage.

## SQL Bindings Generation

At the time of writing, this project is used primarily for generating SQL
bindings. The project root is an application that generates Golang bindings and
SQL table definitions based on the structs that are defined in the package you
provide.

The application requires the use of the `pkg` flag and either the `go` or `sql`
flag with no arguments. The `pkg` argument should be the package path of the
package that you want to use a basis for the output. Using the `sql` flag causes
the command to generate SQL table definitions to stdout. Similarly, use of the
`go` command causes the application to generate the Golang bindings to stdout.

The reason I chose to only generate one language at a time was so that you could
easily pipe the output to language-specific files that you can manipulate with
syntax highlighting right away.

Regarding the output, the code assumes that you have a db() method somewhere in
your code that returns a database pointer and an error. I use something that
looks like this:

```go
var _db *sql.DB
func db() (*sql.DB, error) {

    if _db != nil{
        return _db, nil
    }

    database := "foo"
    host := "bar"
    port := 1234
    password := "username"
    username := "password"

    connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
    _db, err := sql.Open("postgres", connString)
    if err != nil {
        return nil, errors.Stack(err)
    }

    return _db, nil
}
```

In addition to expecting this sort of function available, We also use the
following package for error reporting: `github.com/jackmanlabs/errors`. This
package has the Stack() method to assemble a call stack in addition to the
error to facilitate debugging. If you don't like it, it should easy enough to
replace it with your own error package or error messages.

## Technical Debt

I fully acknowledge that I've built some technical debt into this product. The
need to get something working right away has pushed me to write sloppy code that
is just good enough. If you decide you want to clean this up and make it more
extensible, there may be an invitation to pizza dinner in your future.