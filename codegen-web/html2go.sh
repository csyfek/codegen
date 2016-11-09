#!/usr/bin/env bash

# OUTPUT FILE
echo 'package main' > html.go
echo '// This file created using go generate.' >> html.go

# GENERATE SQL
echo 'var HtmlGenerateSql string = `' >> html.go
cat generate_sql.html >> html.go
echo '`' >> html.go

# FORMAT OUTPUT FILE
gofmt -w html.go