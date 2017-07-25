package usr

import (
	"testing"

	"github.com/coderconvoy/dbase"
)

func Test_User(t *testing.T) {
	tdat := []Usr{
		{"Alice", dbase.Password{}, "/lo", "/pop"},
	}
	t2 := absPath(tdat, "/poop")
	if t2[0].Root != "/poop/lo" {
		t.Errorf("Expected /poop/lo, got %s", t2[0].Root)
	}

}
