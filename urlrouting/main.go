package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	r := httprouter.New()
	r.GET("/", HomeHandler)

	// Posts collection
	r.GET("/posts", PostsIndexHandler)
	r.POST("/posts", PostsCreateHandler)

	// Posts singular
	r.GET("/posts/:id", PostsShowHandler)
	r.PUT("/posts/:id", PostUpdateHandler)
	r.GET("/posts/:id/edit", PostsEditHandler)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", r)
}

func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Home")
}

func PostsIndexHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Posts Index")
}

func PostsCreateHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Posts Create")
}

func PostsShowHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	fmt.Fprintln(rw, "Showing Post", id)
}

func PostUpdateHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Posts Update")
}

func PostsEditHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	fmt.Fprintln(rw, "Post Edit", id)
}