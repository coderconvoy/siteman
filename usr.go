package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/coderconvoy/dbase"
	"github.com/pkg/errors"
)

type Usr struct {
	Username string
	Password dbase.Password
	Folder   string
}

type Asker struct {
	s *bufio.Scanner
}

func NewAsker() Asker {
	return Asker{bufio.NewScanner(os.Stdin)}
}

func (a Asker) Ask(s, def string) string {
	fmt.Printf("%s\n>>", s)
	b := a.s.Scan()
	if !b {
		return def
	}

	return a.s.Text()
}

func RunUserFunc(fname string) {
	fmt.Println("Running User Setup with ", fname)
	uu, err := ReadUsers(fname)
	ask := NewAsker()
	if err != nil {
		fmt.Println("Could not read file", err)
		s := ask.Ask("continue", "n")
		if strings.ToLower(s) != "y" {
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

	for _, v := range uu {

		fmt.Println(v.Username)
	}
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
