package main

import (
	"fmt"
	"net/http"

	"github.com/coderconvoy/dbase"
	"github.com/coderconvoy/siteman/usr"
)

type MuxFunc func(http.ResponseWriter, *http.Request)
type UHandleFunc func(usr.Usr, http.ResponseWriter, *http.Request)

func NewHandler(u []usr.Usr, sc *dbase.SessionControl, f UHandleFunc) MuxFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		l, ok := sc.GetLogin(w, r)

		if ok != dbase.OK {

		}
		us, dok := l.Data.(usr.Usr)
		if dok {
			f(us, w, r)
		}
	}
}

func HandleView(u usr.Usr, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home is where the heart is")
}
