package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/coderconvoy/dbase"
	"github.com/coderconvoy/gojs"
	"github.com/coderconvoy/lazyf"
	"github.com/coderconvoy/siteman/usr"
)

func LoginHandler(uu []usr.Usr, sc *dbase.SessionControl) MuxFunc {
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
		sc.Login(w, uu[found])
		//Add to sessioncontrol
		//point to home
		http.Redirect(w, r, "/home", 303)

	}
}

func LogoutHandler(uu []usr.Usr, sc *dbase.SessionControl) MuxFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		sc.Logout(w, r)
		Handle(w, r)
	}
}

func Handle(w http.ResponseWriter, r *http.Request) {

	w.Write(IndexPage().Bytes())
}

func main() {
	//Underscores are for flag pointers here
	_usrf := lazyf.FlagString("usrf", "{HOME}/.config/users", "userfile", "Set Userdata file")
	_insec := lazyf.FlagBool("i", "insec", "Run Insecure")
	_port := lazyf.FlagString("p", "8081", "port", "Port to run through")
	_pubkey := lazyf.FlagString("pub", "", "pubkey", "Location of public key")
	_privkey := lazyf.FlagString("priv", "", "privkey", "Location of private key")

	lazyf.FlagLoad("c", "{HOME}/.config/init")

	if *_usrf == "" {
		log.Fatal("No Users Cannot run without users")
	}

	users, err := usr.LoadUsers(*_usrf)
	if err != nil {
		fmt.Println("No USers err")
		log.Fatal(err)
	}

	for _, v := range users {
		fmt.Printf("USR,%s,%s\n", v.Username, v.Root)
	}

	sesh := dbase.NewSessionControl(time.Minute * 15)

	gojs.Single.AddFuncs(Asset, AssetDir)

	http.HandleFunc("/save", NewHandler(users, sesh, FileSaver))
	http.HandleFunc("/newfile", NewHandler(users, sesh, FileCreator))
	http.HandleFunc("/delete", NewHandler(users, sesh, FileDeleter))
	http.HandleFunc("/move", NewHandler(users, sesh, FileMover))
	http.HandleFunc("/upload", NewHandler(users, sesh, FileUploader))
	http.HandleFunc("/mkdir", NewHandler(users, sesh, Mkdir))
	http.HandleFunc("/ass/", gojs.AssetHandler("/ass/", gojs.Single))
	http.HandleFunc("/home", NewHandler(users, sesh, HomeView))
	http.HandleFunc("/login", LoginHandler(users, sesh))
	http.HandleFunc("/logout", LogoutHandler(users, sesh))
	http.HandleFunc("/usr/", NewHandler(users, sesh, FileGetter))
	http.HandleFunc("/view/", NewHandler(users, sesh, FileGetter))
	http.HandleFunc("/tabusr/", NewHandler(users, sesh, FileGetter))
	http.HandleFunc("/", Handle)

	fmt.Println("Starting Server")

	if *_insec {
		//Run without https
		log.Fatal(http.ListenAndServe("localhost:"+*_port, nil))
	}

	log.Fatal(http.ListenAndServeTLS(":"+*_port, *_pubkey, *_privkey, nil))
}
