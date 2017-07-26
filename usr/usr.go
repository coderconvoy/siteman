package usr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/coderconvoy/dbase"
	"github.com/pkg/errors"
)

type Usr struct {
	Username string
	Password dbase.Password
	Root     string
	Edit     string
}

func absPath(uu []Usr, fpath string) []Usr {
	res := []Usr{}
	for _, v := range uu {
		u := v
		u.Root = path.Join(fpath, u.Root)
		res = append(res, u)
	}
	return res
}

func (u Usr) ConvertPath(p string) (string, error) {
	fp := path.Join(u.Root, p)
	if !strings.HasPrefix(fp, u.Root) {
		return u.Root, errors.New("Cannot reach outside Root folder")
	}
	return fp, nil
}

func RunUserFunc(fname string) {
	ask := NewAsker()
	askEdit := func(ex int, s string) bool {
		if ex < 0 {
			return true
		}
		fmt.Println(s)
		return ask.AskBool("Edit?", false)
	}

	fmt.Println("Running User Setup with ", fname)
	uu, err := ReadUsers(fname)
	if err != nil {
		fmt.Println("Could not read file", err)

		if !ask.AskBool("continue", false) {
			fmt.Println("Exiting")
			return
		}
	}
	fmt.Println("Continuing")
	name := ask.Ask("Enter a Username", "")

	if name == "" {
		fmt.Println("Not valid")
		return
	}

	existing := -1
	var res Usr

	for k, v := range uu {
		if v.Username == name {
			if !ask.AskBool("Edit Existing User?", false) {
				return
			}
			existing = k
			res = v
			break
		}
	}
	res.Username = name

	if askEdit(existing, "Root Folder = "+res.Root) {
		folder := ask.Ask("Enter a Root Folder", "")
		if folder == "" {
			fmt.Println("Not Valid")
			return
		}
		res.Root = folder
	}

	//Editable folders
	if askEdit(existing, "Editable Paths = "+res.Edit) {
		editable := ask.Ask("Enter editable roots, relative to Root folder", "")
		if editable == "" {
			fmt.Println("Not Valid")
			return
		}
		res.Edit = editable
	}

	if askEdit(existing, "Edit Password?") {
		for {
			pass1 := ask.Ask("Enter Password", "")
			if len(pass1) < 8 {
				fmt.Println("Too short")
				continue
			}
			pass2 := ask.Ask("Confirm Password", "")

			if pass2 != pass1 {
				fmt.Println("No match")
				continue
			}
			res.Password, err = dbase.NewPassword(pass1)
			if err != nil {
				fmt.Println("Password problems:", err)
			}
			break
		}
	}

	//Put it all together
	if existing >= 0 {
		uu[existing] = res
	} else {
		uu = append(uu, res)
	}

	err = WriteUsers(fname, uu)
	if err != nil {
		fmt.Println(err)
	}
}

func LoadUsers(fname string) ([]Usr, error) {
	uu, err := ReadUsers(fname)
	fpath := filepath.Dir(fname)
	return absPath(uu, fpath), err
}

func ReadUsers(fname string) ([]Usr, error) {
	d, err := ioutil.ReadFile(fname)
	if err != nil {
		return []Usr{}, errors.Wrap(err, "Could not load User file")
	}

	var res []Usr
	err = json.Unmarshal(d, &res)
	if err != nil {
		return []Usr{}, errors.Wrap(err, "Could not read User file")
	}

	return res, nil
}

func WriteUsers(fname string, u []Usr) error {
	d, err := json.Marshal(u)
	if err != nil {
		return errors.Wrap(err, "Could not marshal Data")
	}
	err = ioutil.WriteFile(fname, d, 0777)
	if err != nil {
		return errors.Wrap(err, "Could not save")
	}
	return nil
}

func (u Usr) ChPass(npass string) Usr {
	u.Password, _ = dbase.NewPassword(npass)
	return u
}
