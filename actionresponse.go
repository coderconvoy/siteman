package main

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type Edit struct {
	Op     string
	Params string
}

func NewEdit(op string, params ...string) Edit {

	ps := ""
	for k, v := range params {
		if k > 0 {
			ps += ","
		}
		ps += v
	}
	return Edit{op, ps}
}

func WriteEdits(w http.ResponseWriter, edits ...Edit) error {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(edits)
	if err != nil {
		return errors.Errorf("Could marshal json")
	}
	_, err = w.Write(data)
	return err
}
