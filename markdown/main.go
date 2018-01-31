package main

import (
    "net/http"
    "fmt"
    "os"
    "github.com/russross/blackfriday"
)

func main() {
	port := os.GetEnv(PORT)
	if port == "" {
		port = "8080"
	}


	// func Handle(pattern string, handler Handler)
	// func HandleFunc(pattern string, handler func(ResponseWriter, *Request))
	// func ListenAndServe(addr string, handler Handler) error
	http.HandleFunc("/markdown", GenerateMarkDown)
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":"+port, nil)
}

// type ResponseWriter interface
// -- Write([]byte) (int, error)
// -- WriteHeader(int)
// type Request struct
func GenerateMarkDown(rw http.ResponseWriter, r *http.Request) {
	// func MarkdownCommon(input []byte) []byte
	// func (r *Request) FormValue(key string) string
    markdown := blackfriday.MarkdownCommon([]byte(r.FormValue("body")))
    rw.Write(markdown)
    rw.WriteHeader(http.StatusOK)
    fmt.Println(r.Header)
}
