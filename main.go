package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/coderconvoy/siteman/usr"
)

func HandleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Login")

}

func Handle(w http.ResponseWriter, r *http.Request) {

	w.Write(IndexPage().Bytes())
}

func main() {
	usrn := flag.Bool("usr", false, "Create or Edit a User")
	usrf := flag.String("usrf", "usrdata.json", "Set Userdata file")

	flag.Parse()

	if *usrn {
		usr.RunUserFunc(*usrf)
		return
	}

	users, err := usr.LoadUsers(*usrf)
	if err != nil {
		fmt.Println(err)
		return
	}

	_ = users

	http.HandleFunc("/login", HandleLogin)
	http.HandleFunc("/", Handle)

	fmt.Println("Starting Server")

	log.Fatal(http.ListenAndServe(":8090", nil))

}
