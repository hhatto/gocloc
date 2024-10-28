package gocloc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeFile4Python(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "tmp.py")

	if err := os.WriteFile(tmpfile, []byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3
	"""
	pass
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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

func TestAnalyzeFile4PythonInvalid(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "tmp.py")

	if err := os.WriteFile(tmpfile, []byte(`#!/bin/python

class A:
	"""comment1
	comment2
	comment3"""
	pass
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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

func TestAnalyzeFile4PythonNoShebang(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "tmp.py")

	if err := os.WriteFile(tmpfile, []byte(`#!/bin/python
	world
	'''

	b = 1
	"""hello
	commen
	"""

	print a, b
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Python", []string{"#"}, [][]string{{"\"\"\"", "\"\"\""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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

func TestAnalyzeFile4Go(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "tmp.go")

	if err := os.WriteFile(tmpfile, []byte(`package main

func main() {
	var n string /*
		comment
		comment
	*/
}
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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

func TestAnalyzeFile4GoWithOnelineBlockComment(t *testing.T) {
	t.SkipNow()
	tmpfile := filepath.Join(t.TempDir(), "tmp.go")

	if err := os.WriteFile(tmpfile, []byte(`package main

func main() {
	st := "/*"
	a := 1
	en := "*/"
	/* comment */
}
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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

func TestAnalyzeFile4GoWithCommentInnerBlockComment(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "tmp.go")

	if err := os.WriteFile(tmpfile, []byte(`package main

func main() {
	// comment /*
	a := 1
	b := 2
}
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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
	tmpfile := filepath.Join(t.TempDir(), "tmp.go")

	if err := os.WriteFile(tmpfile, []byte(`package main

	func main() {
		a := "/*                */"
		b := "//                  "
	}
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Go", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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
	tmpfile := filepath.Join(t.TempDir(), "tmp.java")

	if err := os.WriteFile(tmpfile, []byte(`/* com */
(* co *)

vo (*
com *)

vo /*
jife */

vo /* ff */
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("ATS", []string{"//"}, [][]string{{"(*", "*)"}, {"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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
	tmpfile := filepath.Join(t.TempDir(), "tmp.java")

	if err := os.WriteFile(tmpfile, []byte(`public class Sample {
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
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Java", []string{"//"}, [][]string{{"/*", "*/"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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
	tmpfile := filepath.Join(t.TempDir(), "Makefile.am")

	if err := os.WriteFile(tmpfile, []byte(`# This is a simple Makefile with comments
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
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Makefile", []string{"#"}, [][]string{{"", ""}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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

func TestAnalyzeReader(t *testing.T) {
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

func TestAnalyzeReader_OnCallbacks(t *testing.T) {
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
	tmpfile := filepath.Join(t.TempDir(), "test.imba")

	if err := os.WriteFile(tmpfile, []byte(`###
This color is my favorite
I need several lines to really
emphasize this fact.
###
const color = "blue"

# this is line comment
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Imba", []string{"#"}, [][]string{{"###", "###"}})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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

func TestAnalyzeFile4Just(t *testing.T) {
	tmpfile := filepath.Join(t.TempDir(), "tmp.jf")

	if err := os.WriteFile(tmpfile, []byte(`polyglot: python js perl sh ruby nu

python:
  #!/usr/bin/env python3
  print('Hello from python!')

js:
  #!/usr/bin/env node
  console.log('Greetings from JavaScript!')  # with comment

# this is comment
`),
		0o600); err != nil {
		t.Fatalf("os.WriteFile() error. err=[%v]", err)
	}

	language := NewLanguage("Just", []string{"#"}, [][]string{{"", ""}}).
		WithRegexLineComments([]string{`^#[^!].*`})
	clocOpts := NewClocOptions()
	clocFile := AnalyzeFile(tmpfile, language, clocOpts)

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
