package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coderconvoy/dbase"
	"github.com/coderconvoy/siteman/usr"
)

func LoginHandler(sc *dbase.SessionControl, uu []usr.Usr) MuxFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//Check for match
		found := -1
		for k, v := range uu {
			if v.Username == r.FormValue("username") && v.Password.Check(r.FormValue("password")) {
				found = k
				break
			}
		}
		if found == -1 {
			//TODO Add Message somewhere
			http.Redirect(w, r, "/", 303)
			return
		}
		//Add to sessioncontrol
		//point to home

		fmt.Fprintln(w, "Login")
	}
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

	sesh := dbase.NewSessionControl(time.Minute * 15)
	_ = users

	http.HandleFunc("/login", HandleLogin)
	http.HandleFunc("/", Handle)

	fmt.Println("Starting Server")

	log.Fatal(http.ListenAndServe(":8090", nil))

}
