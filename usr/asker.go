package usr

import (
	"bufio"
	"fmt"
	"os"
)

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

func (a Asker) AskBool(s string, def bool) bool {
	r := a.Ask(s, "")
	switch r {
	case "y", "Y", "yes", "YES", "Yes", "t", "true":
		return true
	case "n", "N", "no", "NO", "No", "f", "false":
		return false
	}
	return def
}
