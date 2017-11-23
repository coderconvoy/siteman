package usr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/coderconvoy/lazyf"
	"github.com/coderconvoy/passmaker/pswd"
	"github.com/pkg/errors"
)

const (
	CAN_EDIT = iota
	CAN_READ
	NO_READ
)

type Usr struct {
	Username string
	Password pswd.Password
	Root     string
	Paths    map[string]int
}

func Load(lz lazyf.LZ) (Usr, error) {
	pstr, err := lz.PString("pass", "pswd")
	if err != nil {
		return Usr{}, errors.Errorf("User Must have a Password")
	}

	pass, err := pswd.Parse(pstr)
	if err != nil {
		return Usr{}, errors.Errorf("Password Interpretation fail")
	}

	root, err := lz.PString("root")
	if err != nil {
		return Usr{}, errors.Errorf("User Musr have a root folder")
	}
	res := Usr{
		Username: lz.Name,
		Password: pass,
		Root:     root,
		Paths:    make(map[string]int),
	}

	for _, v := range lz.PStringAr("path") {
		sp := strings.Split(v, ":")
		if len(sp) == 1 {
			res.Paths[v] = CAN_EDIT
			continue
		}
		switch sp[0] {
		case "read":
			res.Paths[sp[1]] = CAN_READ
		case "no":
			res.Paths[sp[1]] = NO_READ
		default:
			res.Paths[sp[1]] = CAN_EDIT
		}
	}
	return res, nil
}

func absPath(uu []Usr, fpath string) []Usr {
	res := []Usr{}
	for _, v := range uu {
		u := v
		s := lazyf.EnvReplace(u.Root)
		if len(s) == 0 {
			continue
		}
		if s[0] != '/' {
			u.Root = path.Join(fpath, s)
		}
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

func LoadUsers(fname string) ([]Usr, error) {
	ulist, _, err := lazyf.GetConfig(fname)
	if err != nil {
		return []Usr{}, err
	}

	if len(ulist) == 0 {
		return []Usr{}, errors.Errorf("No users in userfile")
	}
	res := []Usr{}
	for _, v := range ulist {
		u, err := Load(v)
		if err != nil {
			//TODO, add logger
			fmt.Println(err)
			continue
		}
		res = append(res, u)
	}
	fpath := filepath.Dir(fname)
	return absPath(res, fpath), nil
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

func (u Usr) GlobalPermission(fname string) int {
	if !strings.HasPrefix(fname, u.Root) {
		return NO_READ
	}
	tr := strings.TrimPrefix(fname, u.Root)
	return u.Permission(tr)

}

func (u Usr) Permission(fname string) int {
	longest := ""
	res := CAN_EDIT
	for k, v := range u.Paths {
		if len(k) <= len(longest) {
			continue
		}
		if !strings.HasPrefix(fname, k) {
			continue
		}
		longest = k
		res = v
	}
	return res
}
