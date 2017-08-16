package main

import (
	"net/http"

	"github.com/coderconvoy/dbase"
	"github.com/coderconvoy/htmq"
	"github.com/coderconvoy/siteman/usr"
)

type MuxFunc func(http.ResponseWriter, *http.Request)
type UHandleFunc func(usr.Usr, http.ResponseWriter, *http.Request)

func NewHandler(u []usr.Usr, sc *dbase.SessionControl, f UHandleFunc) MuxFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		l, ok := sc.GetLogin(w, r)

		if ok != dbase.OK {
			http.Redirect(w, r, "/", 303)
			return
		}
		us, dok := l.Data.(usr.Usr)
		if dok {
			f(us, w, r)
		}
	}
}

func HomeView(u usr.Usr, w http.ResponseWriter, r *http.Request) {
	p, b := htmq.NewPage("Home", "/ass/css/main.css", "https://ajax.googleapis.com/ajax/libs/jquery/3.2.1/jquery.min.js,/ass/js/fold.js")
	fv, err := FileView(u.Root, "", 6)
	if err != nil {
		b.AddChildren(htmq.NewText("Cannot read home directory: " + err.Error()))
	}

	fv = RootWrap(fv)
	//Copy Paste Area
	cpdiv := htmq.NewParent("div", []*htmq.Tag{
		htmq.NewTextTag("p", "", "id", "loc-p"),
		htmq.QBut("Select", "selectFile(this)"),
		htmq.QBut("Move Here", "moveHere(this)", "class", "with_select hidden"), htmq.NewText("<br>"),
		htmq.QBut("Rename", "rename(this)"),
	}, "id", "copydiv")
	//File View
	tdiv := htmq.NewParent("div", []*htmq.Tag{
		htmq.NewTag("textarea", "id", "filebox"),
		htmq.NewText("<br>"),
		htmq.QBut("Delete", "deleteFile(this)"),
		htmq.QBut("Save", "saveFile(this)"),
	}, "id", "filediv")

	//Folder View
	foldiv := htmq.NewParent("div", []*htmq.Tag{

		htmq.QInput("text", "filename", "id", "foltext"),
		htmq.QBut("Add File", "addFile(this)"),
	}, "id", "foldiv", "style", "display:none;")

	b.AddChildren(cpdiv, fv, tdiv, foldiv)

	b.AddChildren(htmq.QScript("foldStart();"))

	w.Write(p.Bytes())
}
