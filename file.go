package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/coderconvoy/htmq"
	"github.com/coderconvoy/siteman/usr"
)

func RootWrap(ul *htmq.Tag) *htmq.Tag {

	return htmq.NewParent("ul", []*htmq.Tag{
		htmq.NewTextTag("li", "/", "onclick", "fold(this)", "class", "treefolder"),
		ul,
	}, "id", "treetop")
}

func FileView(root, lpath string, md int) (*htmq.Tag, error) {
	cpath := path.Join(root, lpath)
	if !strings.HasPrefix(cpath, root) {
		return nil, errors.New("Tried to escape the root")
	}
	dir, err := ioutil.ReadDir(cpath)
	if err != nil {
		return nil, err
	}
	chids := []*htmq.Tag{}
	for _, v := range dir {
		if v.IsDir() && md > 0 {
			dp := path.Join(lpath, v.Name())
			chids = append(chids, htmq.NewTextTag("li", v.Name(), "onclick", "fold(this)", "class", "treefolder"))
			inner, e2 := FileView(root, dp, md-1)
			inner.AddAttrs("style", "display:none;")
			if e2 != nil {
				err = e2
			}
			chids = append(chids, inner)
			continue
		}
		chids = append(chids, htmq.NewTextTag("li", v.Name(), "onclick", "showFile(this)", "class", "treefile"))
	}
	return htmq.NewParent("ul", chids), err
}

func FileGetter(u usr.Usr, w http.ResponseWriter, r *http.Request) {
	p := strings.TrimPrefix(r.URL.Path, "/usr/")
	p2, err := u.ConvertPath(p)
	if err != nil {
		fmt.Fprintln(w, "Could not find file ", p)
		return
	}
	cc, err := ioutil.ReadFile(p2)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(cc)
}

func FileSaver(u usr.Usr, w http.ResponseWriter, r *http.Request) {
	p := strings.TrimSpace(r.FormValue("fname"))
	if p == "" {
		http.Error(w, "No Filename given", 400)
		return
	}
	p2, err := u.ConvertPath(p)
	if err != nil {
		http.Error(w, "Could not write file: "+err.Error(), 400)
		return
	}
	ioutil.WriteFile(p2, []byte(r.FormValue("fcontents")), 0777)
	return
}

func FileDeleter(u usr.Usr, w http.ResponseWriter, r *http.Request) {
	p := strings.TrimSpace(r.FormValue("fname"))
	if p == "" {
		http.Error(w, "No Filename given", 400)
		return
	}
	p2, err := u.ConvertPath(p)
	if err != nil {
		http.Error(w, "Could not Delete File: "+err.Error(), 400)
		return
	}
	err = os.Remove(p2)
	if err != nil {
		http.Error(w, "Could not Delete File: "+err.Error(), 400)
	}
	return
}

func FileMover(u usr.Usr, w http.ResponseWriter, r *http.Request) {
	fpath := strings.TrimSpace(r.FormValue("fname"))
	if fpath == "" {
		http.Error(w, "No -From- Filename Given", 400)
		return
	}
	tpath := strings.TrimSpace(r.FormValue("tname"))
	if tpath == "" {
		http.Error(w, "No -To- Filename Given", 400)
		return
	}

	sfrom, err := u.ConvertPath(fpath)
	if err != nil {
		http.Error(w, "Could not Move File: "+err.Error(), 400)
		return
	}

	sto, err := u.ConvertPath(tpath)
	if err != nil {
		http.Error(w, "Could not Move File: "+err.Error(), 400)
		return
	}

	err = os.Rename(sfrom, sto)
	if err != nil {
		s := strings.Replace(err.Error(), u.Root, "/", -1)
		http.Error(w, "Could not Move File: "+s, 400)
		return
	}

}

func FileUploader(u usr.Usr, w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.Error(w, "Upload Expected Post", 400)
		return
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("uploadfile")

	if err != nil {
		http.Error(w, "Upload Error: "+err.Error(), 400)
		return
	}
	defer file.Close()

	uppath, err := u.ConvertPath("uploads")

	if err != nil {
		http.Error(w, "Upload Error: "+err.Error(), 400)
		return
	}

	err = os.MkdirAll(uppath, 0777)
	if err != nil {
		http.Error(w, "Upload Error: Could not make uploads Directory: "+err.Error(), 400)
		return
	}

	spath := path.Join(uppath, handler.Filename)
	if !strings.HasPrefix(spath, uppath) {
		http.Error(w, "Path tried to escape parent folder", 400)
		return
	}

	out, err := os.OpenFile(spath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		http.Error(w, "Could not open file for writing", 400)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "File could not write correctly", 400)
	}

	http.Redirect(w, r, "/home", 303)

}
