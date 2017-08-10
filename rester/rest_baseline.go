package rester

import (
	"bytes"
	"fmt"
)

func Main(registers []string) string {

	b := bytes.NewBuffer(nil)

	fmt.Fprint(b, `
package main

import (
	"github.com/jackmanlabs/errors"
	"net/http"
	"strings"
	"log"
	"encoding/xml"
	"encoding/json"
	"github.com/gorilla/mux"
)

func main(){

	var r *mux.Router = mux.NewRouter()


`)

	for _, register := range registers {
		fmt.Fprintln(b, register)
	}

	fmt.Fprint(b, `

	log.Fatal(http.ListenAndServe(":8080", r))
}

func deserialize(r *http.Request, target interface{}) error {
	contentType := r.Header.Get("Content-GoType")

	var err error

	if strings.Contains(contentType, "xml") {
		err = xml.NewDecoder(r.Body).Decode(target)
	} else {
		// Assume JSON as default.
		err = json.NewDecoder(r.Body).Decode(target)
	}

	return errors.Stack(err)
}


func serialize(w http.ResponseWriter, r *http.Request, source interface{}) error {
	contentType := r.Header.Get("Accept")

	// Include some fallback behavior.
	if !(strings.Contains(contentType, "xml") || strings.Contains(contentType, "json")) {
		contentType = r.Header.Get("Content-GoType")
	}

	var err error

	if strings.Contains(contentType, "xml") {
		err = xml.NewEncoder(w).Encode(source)
	} else {
		// Assume JSON as default.
		err = json.NewEncoder(w).Encode(source)
	}

	return errors.Stack(err)
}

// This is a basic filter for error reporting.
// Update it to your liking.
type ErrFilter func(w http.ResponseWriter, r *http.Request) error

func (this ErrFilter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := this(w, r)
	if err != nil {
		w.WriteHeader(500)
		log.Print(err)
	}
}


`)

	return b.String()
}
