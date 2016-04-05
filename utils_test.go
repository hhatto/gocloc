package main

import "testing"

func TestCheckMD5SumIgnore(t *testing.T) {
	fileCache = make(map[string]struct{})

	if checkMD5Sum("./utils_test.go") {
		t.Errorf("invalid sequence")
	}
	if !checkMD5Sum("./utils_test.go") {
		t.Errorf("invalid sequence")
	}
}
