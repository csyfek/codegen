# codegen

A collection of code generation tools for creating SQL schemas, SQL-Go bindings,
HTTP-REST interfaces, and model libraries from Go packages or existing SQL
databases.

This tool is intended to generate code that is able to be compiled and run
immediately. If it doesn't, please submit an issue or a pull request. It is
unlikely, however, that the generated code will be useful without modification.
The concept is to take care of boilerplate for you so that you only need to
worry about the stuff that matters to you.

This suite has evolved significantly from its humble beginnings. Those who have
worked with me have probably seen code generated by this tool or one of its
ancestors. 

The first generation of this tool was typically a small web app with basic HTML
(no JavaScript) interfaces for various kinds of transformations.

In subsequent iterations, it became a command-line tool for creating SQL
bindings and schemas from Go struct definitions.

As my career involved more reverse engineering and language migrations, this
suite has expanded to include the generation of entire REST interfaces and model
libraries from existing databases.

Please note that this suite has not traditionally been consumed by people other
than myself and my coworkers. While I try to keep the documentation up, I do not
intend to ever promise interface stability. Instead, my intention is to preserve
functionality as easily as possible.

When in doubt, just execute the various binaries to see a list of command-line
options:

* `structs2bindings`
* `structs2interface`
* `structs2schema`

## Installation

You can build and install the binaries with standard Go tools. For your
convenience, and mine, there are some make targets:

```
$ make clean
rm -rf structs2bindings structs2schema structs2interface
$ make install
go build github.com/jackmanlabs/codegen/cmd/structs2bindings
go build github.com/jackmanlabs/codegen/cmd/structs2schema
go build github.com/jackmanlabs/codegen/cmd/structs2interface
go install github.com/jackmanlabs/codegen/cmd/structs2bindings
go install github.com/jackmanlabs/codegen/cmd/structs2schema
go install github.com/jackmanlabs/codegen/cmd/structs2interface
```

Make sure `$GOPATH/bin` is in your `$PATH`.

## SQL Bindings Generation: `structs2bindings`

This tool is for generating Go-SQL bindings from Go struct definitions.

```
$ structs2bindings 
Usage of structs2bindings:
  -driver string
        The SQL driver relevant to your request; one of 'sqlite', 'mysql', 'pg', or 'mssql'. (default "mysql")
  -dst string
        The desired output path of the bindings.
  -pkg string
        The package that you want to use for source material.
2018/06/14 16:48:20 The 'pkg' argument is required.
```

Because the default behavior is create bindings for *all* the structs it finds,
a single dump of source would become quite unwieldy. Therefore, the `dst` flag
must be a folder. The tool will create a basic database initializer in that
folder (`bindings.go`) and a number of other files corresponding to the structs
it discovers in the package path you specify.

The `driver` flag specifies the driver for which bindings will be built. As of
this writing `mysql` and `sqlite` are either done or mostly done. The drivers 
`pg` and `mssql` are currently being refactored.

The `pkg` flag is the Go package path that contains the structs for which you
want to create bindings. Any Go package path (import path) is acceptable as long
as it's in your `$GOPATH`. Please note that this is not a path on your
filesystem; see https://blog.golang.org/organizing-go-code if you have 
questions.

In addition to the bindings themselves, a series of test files will also be
created. The intent of these is to validate the SQL and Go syntax of the
bindings. These can be discarded safely, as they're mostly to ensure the
generator did its job properly.

I also use the following package for error reporting:
`github.com/jackmanlabs/errors`. This package has the Stack() method to assemble
a call stack in addition to the error to facilitate debugging. If you don't like
it, it should easy enough to replace it with your own error package or error
messages.

### Example

```
$ structs2bindings 
    \ -driver sqlite
    \ -dst /home/jackman/gopath/src/github.com/jackmanlabs/foo/ds/sqlite
    \ -pkg github.com/jackmanlabs/foo
```

This example creates bindings in the folder
`/home/jackman/gopath/src/github.com/jackmanlabs/foo/ds/sqlite` for the structs
found in the Go package `github.com/jackmanlabs/foo`. The bindings written will
be for Sqlite.

## Data Source Interface Generation: `structs2interface`

This tool is for generating data source interfaces from Go struct definitions.
The generated interfaces should be consistent with the bindings created using
`structs2bindings`.

```
$ structs2interface
Usage of structs2interface:
  -dst string
        The desired output path of the bindings.
  -pkg string
        The package that you want to use for source material.
2018/06/14 17:14:44 The 'pkg' argument is required.
```

Because the default behavior is to create (sub-)interfaces for *all* the structs
it finds, a single dump of source could become unwieldy. Therefore, the `dst`
flag must be a folder. The tool will create a master interface in that
folder (`ds.go`) and a number of other files corresponding to the structs
it discovers in the package path you specify.

The `pkg` flag is the Go package path that contains the structs for which you
want to create interfaces. Any Go package path (import path) is acceptable as 
long as it's in your `$GOPATH`. Please note that this is not a path on your
filesystem; see https://blog.golang.org/organizing-go-code if you have 
questions.

### Example

```
$ structs2interface 
    \ -pkg github.com/jackmanlabs/foo
    \ -dst /home/jackman/gopath/src/github.com/jackmanlabs/foo/ds
```

This example creates interface definitions in the folder
`/home/jackman/gopath/src/github.com/jackmanlabs/foo/ds` for the structs
found in the Go package `github.com/jackmanlabs/foo`. The interface written will
be driver agnostic.

## SQL Schema Generation: `structs2schema`

This tool is for generating an SQL schema from Go struct definitions. The
resulting schema should be compatible with the SQL in the bindings created by
`structs2bindings`.

```
$ structs2schema
Usage of structs2schema:
  -driver string
        The SQL driver relevant to your request; one of 'sqlite', 'mysql', 'pg', or 'mssql'. (default "mysql")
  -pkg string
        The package that you want to use for source material.
2018/06/14 17:22:13 The 'pkg' argument is required.
```

Unlike the other tools, this tool dumps the resulting schema to `stdout`. Simply
use pipes to get it where you want it.

The `driver` flag specifies the driver for which the schema will be built. As of
this writing `mysql` and `sqlite` are either done or mostly done. The drivers 
`pg` and `mssql` are currently being refactored.

The `pkg` flag is the Go package path that contains the structs for which you
want to create bindings. Any Go package path (import path) is acceptable as long
as it's in your `$GOPATH`. Please note that this is not a path on your
filesystem; see https://blog.golang.org/organizing-go-code if you have 
questions.

This tool currently guesses foreign key constraints based on common conventions.
The schema will likely need to be rearranged, however, so that it can be
executed repeatably. Eventually, more intelligent behavior will be implemented.

### Example

```
$ structs2bindings 
    \ -driver sqlite
    \ -dst /home/jackman/gopath/src/github.com/jackmanlabs/foo/ds/sqlite
    \ -pkg github.com/jackmanlabs/foo
```

This example creates bindings in the folder
`/home/jackman/gopath/src/github.com/jackmanlabs/foo/ds/sqlite` for the structs
found in the Go package `github.com/jackmanlabs/foo`. The bindings written will
be for Sqlite.

## Struct Generation from a Database: `db2structs`

This tool is for creating structs (data access objects, DAOs) for already
existing tables. The resulting structs can be then used with the other tools to
create a complete application around an existing database.

The tool requires verification and testing. I haven't used this except in
reverse-engineering gigs, which are usually a little more *improvisational* in
nature. Therefore, formalizing this tool is less of a priority than the others. 

An additional impediment is that databases vary wildly in sophistication and
consistency. Thus, per GIGO, it's quite a task to try and normalize
table and column names so that the structs created with this tool don't trigger
nonsensical output when processed with the other tools.

## REST Interface Generation: `structs2rest`

This tool is for creating an HTTP REST interface based on structs.

Again, this tool hasn't been in huge demand, and needs verification and testing.

Historically, this tool creates a very flat API, which does only basic
type validation and relies on the database for relationship constraint 
validation. This kind of API has been criticized for being too transparent with
regards to the database layer. There's no reason you can't rearrange the API to
be more hierarchical or otherwise build upon the boilerplate.

As a bonus, however, the resulting code has the potential to self-document with 
Swagger/OpenAPI.

There have also been attempts to do validations at the REST layer based on
tags on the source structs.

## Technical Debt

I fully acknowledge that I've built some technical debt into this product. The
need to get something working right away has pushed me to write sloppy code that
is just good enough. If you decide you want to clean this up and make it more
extensible, there may be an invitation to a pizza dinner in your future.

Here are some things that I know need improvement and I'll be working on
eventually:

* Order the `DROP` and `CREATE` statements in the schema intelligently based on foreign keys.
* Optionally switch between primary keys that are auto-incrementing integers or UUIDs.
* Implement a dummy database generator that can be used for unit testing.
* The `db2structs` tool needs to be verified and tested.
* The `structs2rest` tool needs to be verified and tested.