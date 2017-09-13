package main

import (
	"flag"
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
	usrn := flag.Bool("usr", false, "Create or Edit a User")
	usrf := flag.String("usrf", "", "Set Userdata file")
	insec := flag.Bool("i", false, "Run insecure")
	noconf := flag.Bool("noconf", false, "Use Default Configuration")
	confloc := flag.String("config", "", "Config File Location")

	flag.Parse()

	conf, err := getConfig(*confloc, *noconf)
	if err != nil {
		fmt.Println("Config Error:", err)
		return
	}
	//Testing map
	fmt.Println("Config")
	for k, v := range conf.Deets {
		fmt.Println("\t", k, ":", v)
	}

	usrloc := *usrf
	if *usrf == "" {
		loc := conf.PStringD("usrdata.json", "userfile")
		fmt.Println("userfile == ", loc)
		usrloc = lazyf.EnvReplace(loc)
	}
	fmt.Println("Userfile at :" + usrloc)

	if *usrn {
		usr.RunUserFunc(usrloc)
		return
	}

	users, err := usr.LoadUsers(usrloc)
	if err != nil {
		fmt.Println(err)
		return
	}

	sesh := dbase.NewSessionControl(time.Minute * 15)

	gojs.Single.AddFuncs(Asset, AssetDir)

	http.HandleFunc("/save", NewHandler(users, sesh, FileSaver))
	http.HandleFunc("/delete", NewHandler(users, sesh, FileDeleter))
	http.HandleFunc("/move", NewHandler(users, sesh, FileMover))
	http.HandleFunc("/upload", NewHandler(users, sesh, FileUploader))
	http.HandleFunc("/mkdir", NewHandler(users, sesh, Mkdir))
	http.HandleFunc("/ass/", gojs.AssetHandler("/ass/", gojs.Single))
	http.HandleFunc("/home", NewHandler(users, sesh, HomeView))
	http.HandleFunc("/login", LoginHandler(users, sesh))
	http.HandleFunc("/logout", LogoutHandler(users, sesh))
	http.HandleFunc("/usr/", NewHandler(users, sesh, FileGetter))
	http.HandleFunc("/tabusr/", NewHandler(users, sesh, FileGetter))
	http.HandleFunc("/", Handle)

	fmt.Println("Starting Server")

	if *insec || conf.PBoolD(false, "insec", "insecure") {
		//Run without https
		insPort := conf.PStringD("8090", "ins-port")
		log.Fatal(http.ListenAndServe(":"+insPort, nil))
	}

	pubkey := conf.PStringD("data/server.pub", "pubkey")
	privkey := conf.PStringD("data/server.key", "privkey")
	secPort := conf.PStringD("8091", "port")

	log.Fatal(http.ListenAndServeTLS(":"+secPort, pubkey, privkey, nil))

}
