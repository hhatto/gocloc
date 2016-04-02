package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Language struct {
	name           string
	line_comment   string
	multi_line     string
	multi_line_end string
	files          []string
	code           int32
	comments       int32
	blanks         int32
	lines          int32
	total          int32
	printed        bool
}

type ClocFile struct {
	name     string
	code     int32
	comments int32
	blanks   int32
	lines    int32
}

const LANG_HEADER string = "Language                     files          blank        comment           code"
const ROW string = "-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------"

var rowLen = 79

func NewLanguage(name, line_comment, multi_line, multi_line_end string) *Language {
	return &Language{
		name:           name,
		line_comment:   line_comment,
		multi_line:     multi_line,
		multi_line_end: multi_line_end,
		files:          []string{},
	}
}

func containComments(line, commentStart, commentEnd string) bool {
	inComments := 0
	for i := 0; i < len(line)/len(commentStart); i += len(commentStart) {
		section := line[i : i+len(commentStart)]

		if section == commentStart {
			inComments += 1
		} else if section == commentEnd {
			if inComments != 0 {
				inComments -= 1
			}
		}
	}
	return inComments != 0
}

func getFileType(path string) string {
	ext := filepath.Ext(path)
	if strings.ToLower(filepath.Base(path)) == "makefile" {
		return "makefile"
	}
	if len(ext) >= 2 {
		return ext[1:]
	}
	return ext
}

func getAllFiles(paths []string, languages map[string]*Language) {
	for _, root := range paths {
		walkCallback := func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}

			if ext := getFileType(rel); ext != "" {
				if _, ok := languages[ext]; ok {
					languages[ext].files = append(languages[ext].files, path)
					//fmt.Printf("match rel=[%v] path=[%v] ext=[%v]\n", rel, path, ext)
				}
			}
			return nil
		}
		if err := filepath.Walk(root, walkCallback); err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	flag.Parse()
	paths := flag.Args()
	if len(paths) <= 0 {
		flag.PrintDefaults()
		return
	}

	action_script := NewLanguage("ActionScript", "//", "/*", "*/")
	asm := NewLanguage("Assembly", "", "", "")
	awk := NewLanguage("Awk", "#", "", "")
	bash := NewLanguage("BASH", "#", "", "")
	batch := NewLanguage("Batch", "REM", "", "")
	c := NewLanguage("C", "//", "/*", "*/")
	c_header := NewLanguage("C Header", "//", "/*", "*/")
	c_sharp := NewLanguage("C#", "//", "/*", "*/")
	c_shell := NewLanguage("C Shell", "#", "", "")
	clojure := NewLanguage("Clojure", ",#,#_", "", "")
	coffee_script := NewLanguage("CoffeeScript", "#", "###", "###")
	cold_fusion := NewLanguage("ColdFusion", "", "<!---", "--->")
	cf_script := NewLanguage("ColdFusion CFScript", "//", "/*", "*/")
	cpp := NewLanguage("C++", "//", "/*", "*/")
	cpp_header := NewLanguage("C++ Header", "//", "/*", "*/")
	css := NewLanguage("CSS", "//", "/*", "*/")
	d := NewLanguage("D", "//", "/*", "*/")
	dart := NewLanguage("Dart", "//", "/*", "*/")
	device_tree := NewLanguage("Device Tree", "//", "/*", "*/")
	lisp := NewLanguage("LISP", "", "#|", "|#")
	fortran_legacy := NewLanguage("FORTRAN Legacy", "c,C,!,*", "", "")
	fortran_modern := NewLanguage("FORTRAN Modern", "!", "", "")
	golang := NewLanguage("Go", "//", "/*", "*/")
	haskell := NewLanguage("Haskell", "--", "", "")
	html := NewLanguage("HTML", "<!--", "<!--", "-->")
	jai := NewLanguage("JAI", "//", "/*", "*/")
	java := NewLanguage("Java", "//", "/*", "*/")
	java_script := NewLanguage("JavaScript", "//", "/*", "*/")
	julia := NewLanguage("Julia", "#", "#:=", ":=#")
	json := NewLanguage("JSON", "", "", "")
	jsx := NewLanguage("JSX", "//", "/*", "*/")
	less := NewLanguage("LESS", "//", "/*", "*/")
	linker_script := NewLanguage("LD Script", "//", "/*", "*/")
	lua := NewLanguage("Lua", "--", "--[[", "]]")
	makefile := NewLanguage("Makefile", "#", "", "")
	markdown := NewLanguage("Markdown", "", "", "")
	mustache := NewLanguage("Mustache", "", "{{!", "))")
	objective_c := NewLanguage("Objective C", "//", "/*", "*/")
	objective_cpp := NewLanguage("Objective C++", "//", "/*", "*/")
	ocaml := NewLanguage("OCaml", "", "(*", "*)")
	php := NewLanguage("PHP", "#,//", "/*", "*/")
	pascal := NewLanguage("Pascal", "//,(*", "{", ")")
	polly := NewLanguage("Polly", "<!--", "<!--", "-->")
	perl := NewLanguage("Perl", "#", ":=", ":=cut")
	protobuf := NewLanguage("Protocol Buffers", "//", "", "")
	python := NewLanguage("Python", "#", "'''", "'''")
	r := NewLanguage("R", "#", "", "")
	ruby := NewLanguage("Ruby", "#", ":=begin", ":=end")
	ruby_html := NewLanguage("Ruby HTML", "<!--", "<!--", "-->")
	rust := NewLanguage("Rust", "//,///,//!", "/*", "*/")
	sass := NewLanguage("Sass", "//", "/*", "*/")
	sml := NewLanguage("Standard ML", "", "(*", "*)")
	sql := NewLanguage("SQL", "--", "/*", "*/")
	swift := NewLanguage("Swift", "//", "/*", "*/")
	tex := NewLanguage("TeX", "%", "", "")
	text := NewLanguage("Plain Text", "", "", "")
	toml := NewLanguage("TOML", "#", "", "")
	type_script := NewLanguage("TypeScript", "//", "/*", "*/")
	vim_script := NewLanguage("Vim script", "\"", "", "")
	xml := NewLanguage("XML", "<!--", "<!--", "-->")
	yaml := NewLanguage("YAML", "#", "", "")
	yacc := NewLanguage("Yacc", "//", "/*", "*/")
	zsh := NewLanguage("Zsh", "#", "", "")

	languages := map[string]*Language{
		"as":       action_script,
		"s":        asm,
		"awk":      awk,
		"bat":      batch,
		"btm":      batch,
		"cmd":      batch,
		"bash":     bash,
		"sh":       bash,
		"c":        c,
		"csh":      c_shell,
		"ec":       c,
		"pgc":      c,
		"cs":       c_sharp,
		"clj":      clojure,
		"coffee":   coffee_script,
		"cfm":      cold_fusion,
		"cfc":      cf_script,
		"cc":       cpp,
		"cpp":      cpp,
		"cxx":      cpp,
		"pcc":      cpp,
		"c++":      cpp,
		"css":      css,
		"d":        d,
		"dart":     dart,
		"dts":      device_tree,
		"dtsi":     device_tree,
		"el":       lisp,
		"lisp":     lisp,
		"lsp":      lisp,
		"lua":      lua,
		"sc":       lisp,
		"f":        fortran_legacy,
		"f77":      fortran_legacy,
		"for":      fortran_legacy,
		"ftn":      fortran_legacy,
		"pfo":      fortran_legacy,
		"f90":      fortran_modern,
		"f95":      fortran_modern,
		"f03":      fortran_modern,
		"f08":      fortran_modern,
		"go":       golang,
		"h":        c_header,
		"hs":       haskell,
		"hpp":      cpp_header,
		"hh":       cpp_header,
		"html":     html,
		"hxx":      cpp_header,
		"jai":      jai,
		"java":     java,
		"js":       java_script,
		"jl":       julia,
		"json":     json,
		"jsx":      jsx,
		"lds":      linker_script,
		"less":     less,
		"m":        objective_c,
		"md":       markdown,
		"markdown": markdown,
		"ml":       ocaml,
		"mli":      ocaml,
		"mm":       objective_cpp,
		"makefile": makefile,
		"mustache": mustache,
		"php":      php,
		"pas":      pascal,
		"pl":       perl,
		"text":     text,
		"txt":      text,
		"polly":    polly,
		"proto":    protobuf,
		"py":       python,
		"r":        r,
		"rake":     ruby,
		"rb":       ruby,
		"rhtml":    ruby_html,
		"rs":       rust,
		"sass":     sass,
		"scss":     sass,
		"sml":      sml,
		"sql":      sql,
		"swift":    swift,
		"tex":      tex,
		"sty":      tex,
		"toml":     toml,
		"ts":       type_script,
		"vim":      vim_script,
		"xml":      xml,
		"yaml":     yaml,
		"yml":      yaml,
		"y":        yacc,
		"zsh":      zsh,
	}

	total := NewLanguage("TOTAL", "", "", "")
	fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	fmt.Printf("%.[2]*[1]s\n", LANG_HEADER, rowLen)
	fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)

	getAllFiles(paths, languages)
	for _, language := range languages {
		if language.printed {
			continue
		}

		for _, file := range language.files {
			isInComments := false

			func(filename string) {
				fp, err := os.Open(filename)
				if err != nil {
					panic(err)
				}
				defer fp.Close()

				scanner := bufio.NewScanner(fp)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())

					if len(strings.TrimSpace(line)) == 0 {
						language.blanks += 1
						continue
					}

					if language.multi_line != "" {
						if strings.HasPrefix(line, language.multi_line) {
							isInComments = true
						} else if containComments(line, language.multi_line, language.multi_line_end) {
							isInComments = true
							language.code += 1
						}
					}

					if isInComments {
						if strings.Contains(line, language.multi_line_end) {
							isInComments = false
						}
						language.comments += 1
						continue
					}

					if language.line_comment != "" {
						single_comments := strings.Split(language.line_comment, ",")
						isSingleComment := false
						for _, single_comment := range single_comments {
							if strings.HasPrefix(line, single_comment) {
								language.comments += 1
								isSingleComment = true
								break
							}
						}
						if isSingleComment {
							continue
						}
					}

					language.code += 1
				}
			}(file)
		}

		language.printed = true

		files := int32(len(language.files))
		if len(language.files) <= 0 {
			continue
		}

		fmt.Printf("%-27v %6v %14v %14v %14v\n",
			language.name, files, language.blanks, language.comments, language.code)
		total.total += files
		total.blanks += language.blanks
		total.comments += language.comments
		total.code += language.code
	}

	fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	fmt.Printf("%-27v %6v %14v %14v %14v\n",
		"TOTAL", total.total, total.blanks, total.comments, total.code)
	fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
}
