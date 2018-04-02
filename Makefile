.PHONY: clean charts all

all: structs2bindings structs2schema

clean:
	rm -f structs2bindings
	rm -f structs2schema

structs2bindings:
	go build github.com/jackmanlabs/codegen/cmd/structs2bindings

structs2schema:
	go build github.com/jackmanlabs/codegen/cmd/structs2schema
