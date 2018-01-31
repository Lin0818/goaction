package main

import (
    "net/http"
)

func main() {
	// func ListenAndServe(addr string, handler Handler) error
	// func FileServer(root FileSystem) Handler
	// func (d Dir) Open(name string) (File, error)
    http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))
}
