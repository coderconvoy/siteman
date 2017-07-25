package main

import (
	"net/http"

	"github.com/coderconvoy/dbase"
)

type MuxFunc func(http.ResponseWriter, *http.Request)
type UHandleFunc func(Usr, http.ResponseWriter, *http.Request)

func NewHandler(u []User, sc *dbase.SessionControl, f UHandleFunc) MuxFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		l, ok := sc.GetLogin(w, r)
		if ok != dbase.OK {
		}
	}
}

func HandleView(u Usr, w http.ResponseWriter, r *http.Request) {

}
