package main

import (
	"net/http"
	"fmt"
	"log"

	"github.com/julienschmidt/httprouter"
)

func index(rw http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(rw, "Welcome.\n")
}

func hello(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	p := ps.ByName("name")
	fmt.Fprintf(rw, "Hello, %s.\n", p)
}

func main() {
	r := httprouter.New()
	r.GET("/", index)
	r.GET("/hello/:name", hello)

	log.Fatal(http.ListenAndServe(":8080", r))
}