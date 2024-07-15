package gocloc

import (
	"bytes"
	"os"
	"testing"
)

func TestAnalayzeFile4Python(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tmp.py")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3
	"""
	pass
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
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
	if clocFile.Lang != "Python" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4PythonInvalid(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tmp.py")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3"""
	pass
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
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
	if clocFile.Lang != "Python" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4PythonNoShebang(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tmp.py")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`a = '''hello
	world
	'''

	b = 1
	"""hello
	commen
	"""

	print a, b
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 2 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 3 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 5 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Python" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4Go(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tmp.go")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`package main

func main() {
	var n string /*
		comment
		comment
	*/
}
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 3 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 4 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Go" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4GoWithOnelineBlockComment(t *testing.T) {
	t.SkipNow()
	tmpfile, err := os.CreateTemp("", "tmp.go")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`package main

func main() {
	st := "/*"
	a := 1
	en := "*/"
	/* comment */
}
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 1 { // cloc->3, tokei->1, gocloc->4
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 6 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Go" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4GoWithCommentInnerBlockComment(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tmp.go")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`package main

func main() {
	// comment /*
	a := 1
	b := 2
}
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 1 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 5 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Go" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalyzeFile4GoWithNoComment(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tmp.go")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`package main

	func main() {
		a := "/*                */"
		b := "//                  "
	}
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 0 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 5 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Go" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalyzeFile4ATSWithDoubleMultilineComments(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tmp.java")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`/* com */
(* co *)

vo (*
com *)

vo /*
jife */

vo /* ff */
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("ATS", []string{"//"}, [][]string{{"(*", "*)"}, {"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 3 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 4 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 3 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "ATS" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalyzeFile4JavaWithCommentInCodeLine(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tmp.java")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`public class Sample {
		public static void main(String args[]){
		int a; /* A takes care of counts */
		int b;
		int c;
		String d; /*Just adding comments */
		bool e; /*
		comment*/
		bool f; /*
		comment1
		comment2
		*/
		/*End of Main*/
		}
		}
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Java", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 0 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 5 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 10 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Java" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalyzeFile4Makefile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "Makefile.am")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`# This is a simple Makefile with comments
	.PHONY: test build

	build:
		mkdir -p bin
		GO111MODULE=on go build -o ./bin/gocloc cmd/gocloc/main.go

	# Another comment
	update-package:
		GO111MODULE=on go get -u github.com/hhatto/gocloc

	run-example:
		GO111MODULE=on go run examples/languages.go
		GO111MODULE=on go run examples/files.go

	test:
		GO111MODULE=on go test -v
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Makefile", []string{"#"}, [][]string{{"", ""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 4 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 2 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 11 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Makefile" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeReader(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3
	"""
	pass
`))

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeReader("test.py", language, buf, clocOpts)

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 4 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 3 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Python" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeReader_OnCallbacks(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`foo
		"""bar

`))

	var lines int
	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocOpts.OnCode = func(line string) {
		if line != "foo" {
			t.Errorf("invalid logic. code_line=%v", line)
		}
		lines++
	}

	clocOpts.OnBlank = func(line string) {
		if line != "" {
			t.Errorf("invalid logic. blank_line=%v", line)
		}
		lines++
	}

	clocOpts.OnComment = func(line string) {
		if line != "\"\"\"bar" {
			t.Errorf("invalid logic. comment_line=%v", line)
		}
		lines++
	}

	AnalyzeReader("test.py", language, buf, clocOpts)

	if lines != 3 {
		t.Errorf("invalid logic. lines=%v", lines)
	}
}

func TestAnalyzeFile4Imba(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test.imba")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`###
This color is my favorite
I need several lines to really
emphasize this fact.
###
const color = "blue"

# this is line comment
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Imba", []string{"#"}, [][]string{{"###", "###"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 1 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 6 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 1 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Imba" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}

func TestAnalayzeFile4Just(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "tmp.go")
	if err != nil {
		t.Logf("os.CreateTemp() error. err=[%v]", err)
		return
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(`polyglot: python js perl sh ruby nu

python:
  #!/usr/bin/env python3
  print('Hello from python!')

js:
  #!/usr/bin/env node
  console.log('Greetings from JavaScript!')  # with comment

# this is comment
`),
	); err != nil {
		t.Fatalf("tmpfile.Write() error. err=[%v]", err)
	}

	language := NewLanguage("Just", []string{"#"}, [][]string{{"", ""}}).
		WithRegexLineComments([]string{`^#[^!].*`})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile.Name(), language, clocOpts)
	tmpfile.Close()

	if clocFile.Blanks != 3 {
		t.Errorf("invalid logic. blanks=%v", clocFile.Blanks)
	}
	if clocFile.Comments != 1 {
		t.Errorf("invalid logic. comments=%v", clocFile.Comments)
	}
	if clocFile.Code != 7 {
		t.Errorf("invalid logic. code=%v", clocFile.Code)
	}
	if clocFile.Lang != "Just" {
		t.Errorf("invalid logic. lang=%v", clocFile.Lang)
	}
}
