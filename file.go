package main

import (
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"

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
		return nil, errors.Errorf("Tried to escape the root :%s:%s ", root, cpath)
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
	var maxSize int64 = 10000 //
	p := strings.TrimPrefix(r.URL.Path, "/usr/")

	if strings.HasPrefix(r.URL.Path, "/view/") {
		maxSize *= 1000
		p = strings.TrimPrefix(r.URL.Path, "/view/")
	}

	p2, err := u.ConvertPath(p)
	if err != nil {
		http.Error(w, "Could not access file by that name"+err.Error(), 400)
		return
	}
	f, err := os.Open(p2)
	if err != nil {
		http.Error(w, "Could not read file: "+err.Error(), 400)
		return
	}
	defer f.Close()

	finfo, err := f.Stat()
	if err != nil {
		http.Error(w, "Could not get file Size: "+err.Error(), 400)
		return
	}
	if finfo.Size() > maxSize {
		http.Error(w, "big file", 400)
		return
	}

	w.Header().Set("Content-Type", mime.TypeByExtension(path.Ext(p)))
	io.Copy(w, f)
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

	WriteEdits(w, NewEdit("say", "File Saved : "+p), NewEdit("unchange"))
}

func FileCreator(u usr.Usr, w http.ResponseWriter, r *http.Request) {
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
	WriteEdits(w, NewEdit("new", p), NewEdit("say", "File Created : "+p))
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
	//Delete properly if already in trash folder:
	if strings.HasPrefix(p2, path.Join(u.Root, "/trash/")) {
		err = os.RemoveAll(p2)
		if err != nil {
			http.Error(w, "Could not Delete File: "+err.Error(), 400)
			return
		}
		WriteEdits(w, NewEdit("rm", p), NewEdit("say", "Permanently Deleted: "+p))
		return
	}

	tpath := path.Join(u.Root, "/trash/")
	err = os.MkdirAll(tpath, 0777)
	if err != nil {
		http.Error(w, "Could not Create Trash directory: "+err.Error(), 400)
		return
	}

	npath := path.Join(tpath, path.Base(p))

	err = os.Rename(p2, npath)
	if err != nil {
		http.Error(w, "Could not move to trash"+err.Error(), 400)
		return
	}

	bn := path.Base(p)
	WriteEdits(w, NewEdit("mkdir", "trash"), NewEdit("mv", p, path.Join("/trash/", bn)), NewEdit("say", p+" moved to trash"))
}

func Mkdir(u usr.Usr, w http.ResponseWriter, r *http.Request) {
	p := strings.TrimSpace(r.FormValue("fname"))

	path, err := u.ConvertPath(p)
	if err != nil {
		http.Error(w, "Could not create directory: "+err.Error(), 400)
		return
	}

	err = os.MkdirAll(path, 0777)
	if err != nil {
		http.Error(w, "Could not Create directory: "+err.Error(), 400)
		return
	}

	WriteEdits(w, NewEdit("mkdir", p))

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

	WriteEdits(w, NewEdit("mv", fpath, tpath), NewEdit("say", "Moved"))

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

	fup := r.FormValue("fup-location")
	if fup == "" {
		http.Error(w, "No Folder Location provided", 400)
		return
	}

	uppath, err := u.ConvertPath(fup)
	if err != nil {
		http.Error(w, "Upload Error: "+err.Error(), 400)
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
