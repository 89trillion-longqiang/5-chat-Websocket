package main

import (
	"testing"

	"chat/module"
)

func Test_NewHub(t *testing.T){

	ret := module.NewHub()

	if ret == nil {
		t.Errorf("Tset_NewHub error ")
	}else {
		t.Log("Tset_NewHub pass ")
	}
}