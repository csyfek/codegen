.PHONY: clean charts all

all: structs2bindings

clean:
	rm -f structs2bindings

structs2bindings:
	go build github.com/jackmanlabs/codegen/cmd/structs2bindings
