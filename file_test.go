package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestAnalayzeFile4Python(t *testing.T) {
	tmpfile, err := ioutil.TempFile("", "tmp.py")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	t.Log(tmpfile.Name())
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3
	"""
	pass
`))

	language := NewLanguage("Python", "#", "\"\"\"", "\"\"\"")
	clocFile := analyzeFile(tmpfile.Name(), language)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 4 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 3 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
}

func TestAnalayzeFile4PythonInvalid(t *testing.T) {
	t.SkipNow()
	tmpfile, err := ioutil.TempFile("", "tmp.py")
	if err != nil {
		t.Logf("ioutil.TempFile() error. err=[%v]", err)
		return
	}
	t.Log(tmpfile.Name())
	defer os.Remove(tmpfile.Name())

	tmpfile.Write([]byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3"""
	pass
`))

	language := NewLanguage("Python", "#", "\"\"\"", "\"\"\"")
	clocFile := analyzeFile(tmpfile.Name(), language)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 3 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 3 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
}
