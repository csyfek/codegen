package main

//go:generate ./html2go.sh

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/", &handlerGenerateSql{})

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func writeError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	w.Write([]byte(err.Error()))
}
