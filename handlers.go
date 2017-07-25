package main

import "net/http"

type MuxFunc func(http.ResponseWriter, *http.Request)
type UHandleFunc func(Usr, http.ResponseWriter, *http.Request)

func NewHandler(u []User, f UHandleFunc) MuxFunc {

	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func HandleView(u Usr, w http.ResponseWriter, r *http.Request) {

}
