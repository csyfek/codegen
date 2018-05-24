.PHONY: clean charts all

BINARIES=structs2bindings structs2schema structs2interface

all: $(BINARIES)

clean:
	rm -rf $(BINARIES)

structs2bindings:
	go build github.com/jackmanlabs/codegen/cmd/structs2bindings

structs2schema:
	go build github.com/jackmanlabs/codegen/cmd/structs2schema

structs2interface:
	go build github.com/jackmanlabs/codegen/cmd/structs2interface

install: $(BINARIES)
	go install github.com/jackmanlabs/codegen/cmd/structs2interface
	go install github.com/jackmanlabs/codegen/cmd/structs2schema
	go install github.com/jackmanlabs/codegen/cmd/structs2bindings
