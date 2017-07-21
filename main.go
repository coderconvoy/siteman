package main

import (
	"fmt"
	"log"
	"net/http"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Login")

}

func Handle(w http.ResponseWriter, r *http.Request) {

	w.Write(IndexPage().Bytes())
}

func main() {

	http.HandleFunc("/login", HandleLogin)
	http.HandleFunc("/", Handle)

	fmt.Println("Starting Server")

	log.Fatal(http.ListenAndServe(":8090", nil))

}
