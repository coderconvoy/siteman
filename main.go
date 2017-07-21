package main

import (
	"flag"
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
	usr := flag.Bool("usr", false, "Create or Edit a User")

	flag.Parse()

	if *usr {
		RunUserFunc("test_data/out/ulist")
		return
	}

	fmt.Println("Starting Server")

	log.Fatal(http.ListenAndServe(":8090", nil))

}
