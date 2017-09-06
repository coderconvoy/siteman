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
		htmq.QBut("Select", "selectFile(this)"),
		htmq.QBut("Move Here", "moveHere(this)", "class", "with_select hidden"),
		htmq.QBut("Rename", "rename(this)"),
		htmq.NewTextTag("p", "", "id", "loc-p"),
	}, "id", "copydiv")
	//File View
	tdiv := htmq.NewParent("div", []*htmq.Tag{
		htmq.NewTag("textarea", "id", "filebox"),
		htmq.NewTag("img", "id", "fileimg", "class", "hidden"),
		htmq.NewText("<br>"),
		htmq.NewParent("div", []*htmq.Tag{
			htmq.QBut("Delete", "deleteFile(this)"),
			htmq.QBut("Save", "saveFile(this)"),
		}, "class", "bottomrow"),
	}, "id", "filediv")

	//Folder View
	foldiv := htmq.NewParent("div", []*htmq.Tag{
		htmq.QBut("Add Folder Here", "addFolder(this)"),
		htmq.QBut("Add File Here", "addFile(this)"),
		htmq.QText("<br>"),
		htmq.QUpload("/upload", []*htmq.Tag{htmq.QInput("text", "fup-location", "id", "fup-location", "--hidden")}),
		htmq.QText("<br>"),
		htmq.QBut("Delete Folder", "deleteFolder(this)"),
	}, "id", "foldiv", "style", "display:none;")

	b.AddChildren(fv, htmq.NewParent("div", []*htmq.Tag{htmq.NewTag("div", "id", "messbar", "class", "hidden"), cpdiv, tdiv, foldiv}, "id", "rightdiv"))

	b.AddChildren(htmq.QScript("foldStart();"))

	w.Write(p.Bytes())
}
