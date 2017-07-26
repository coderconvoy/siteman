package main

import (
	"fmt"
	"testing"
)

func Test_Files(t *testing.T) {
	fv, err := FileView("test_data", "", 3)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(fv)
}
