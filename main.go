package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	flags "github.com/jessevdk/go-flags"
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
type Languages []Language

func (ls Languages) Len() int {
	return len(ls)
}
func (ls Languages) Swap(i, j int) {
	ls[i].code, ls[j].code = ls[j].code, ls[i].code
}
func (ls Languages) Less(i, j int) bool {
	return ls[i].code < ls[i].code
}

type ClocFile struct {
	name     string
	code     int32
	comments int32
	blanks   int32
	lines    int32
}

const FILE_HEADER string = "File                         "
const LANG_HEADER string = "Language                     "
const COMMON_HEADER string = "files          blank        comment           code"
const ROW string = "-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------"

var ReShebangEnv *regexp.Regexp = regexp.MustCompile("^#!(\\S+/env) ([a-zA-Z]+)")
var ReShebangLang *regexp.Regexp = regexp.MustCompile("^#!/[\\w/]([a-zA-Z]+)")
var rowLen = 79
var LanguageByScript map[string]string

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

func getFileTypeByShebang(path string) (shebangLang string, ok bool) {
	func() {
		fp, err := os.Open(path)
		if err != nil {
			return // ignore error
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			line := scanner.Text()
			l := strings.TrimSpace(line)

			if ReShebangEnv.MatchString(l) {
				ret := ReShebangEnv.FindAllStringSubmatch(l, -1)
				if len(ret[0]) == 3 {
					shebangLang = ret[0][2]
					break
				}
			}

			if ReShebangLang.MatchString(l) {
				ret := ReShebangLang.FindAllStringSubmatch(l, -1)
				if len(ret) == 3 {
					shebangLang = ret[0][1]
					break
				}
			}

			break
		}
	}()

	sl, ok := LanguageByScript[shebangLang]
	if ok {
		return sl, ok
	}
	return shebangLang, false
}

func getFileType(path string) (ext string, ok bool) {
	ext = filepath.Ext(path)
	if strings.ToLower(filepath.Base(path)) == "makefile" {
		return "makefile", true
	}
	if len(ext) >= 2 {
		return ext[1:], true
	}
	ext, ok = getFileTypeByShebang(path)
	return ext, ok
}

func getAllFiles(paths []string, languages map[string]*Language) (filenum, maxPathLen int) {
	maxPathLen = 0
	for _, root := range paths {
		walkCallback := func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}

			p := filepath.Join(root, rel)
			if ext, ok := getFileType(p); ok {
				if _, ok := languages[ext]; ok {
					languages[ext].files = append(languages[ext].files, p)
					filenum += 1
					l := len(p)
					if maxPathLen < l {
						maxPathLen = l
					}
				}
			}
			return nil
		}
		if err := filepath.Walk(root, walkCallback); err != nil {
			fmt.Println(err)
		}
	}
	return
}

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gocloc"
	parser.Usage = "[OPTIONS] PATH[...]"

	paths, err := flags.Parse(&opts)
	if err != nil {
		parser.WriteHelp(os.Stdout)
		return
	}
	if len(paths) <= 0 {
		parser.WriteHelp(os.Stdout)
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
	LanguageByScript = map[string]string{
		"perl":   "pl",
		"python": "py",
		"ruby":   "rb",
	}

	total := NewLanguage("TOTAL", "", "", "")
	num, maxPathLen := getAllFiles(paths, languages)
	headerLen := 40

	if opts.Byfile {
		headerLen := maxPathLen
		rowLen = maxPathLen + len(COMMON_HEADER) + 2
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
		fmt.Printf("%-[2]*[1]s  %[3]s\n", FILE_HEADER, headerLen, COMMON_HEADER)
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	} else {
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
		fmt.Printf("%.[2]*[1]s%[3]s\n", LANG_HEADER, headerLen, COMMON_HEADER)
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	}

	clocFiles := make(map[string]*ClocFile, num)
	fileCache := make(map[string]bool)

	for _, language := range languages {
		if language.printed {
			continue
		}

		for _, file := range language.files {
			clocFiles[file] = &ClocFile{
				name: file,
			}
			isInComments := false

			func() {
				fp, err := os.Open(file)
				if err != nil {
					return // ignore error
				}
				defer fp.Close()

				scanner := bufio.NewScanner(fp)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())

					if len(strings.TrimSpace(line)) == 0 {
						clocFiles[file].blanks += 1
						continue
					}

					if language.multi_line != "" {
						if strings.HasPrefix(line, language.multi_line) {
							isInComments = true
						} else if containComments(line, language.multi_line, language.multi_line_end) {
							isInComments = true
							clocFiles[file].code += 1
						}
					}

					if isInComments {
						if strings.Contains(line, language.multi_line_end) {
							isInComments = false
						}
						clocFiles[file].comments += 1
						continue
					}

					if language.line_comment != "" {
						single_comments := strings.Split(language.line_comment, ",")
						isSingleComment := false
						for _, single_comment := range single_comments {
							if strings.HasPrefix(line, single_comment) {
								clocFiles[file].comments += 1
								isSingleComment = true
								break
							}
						}
						if isSingleComment {
							continue
						}
					}

					clocFiles[file].code += 1
				}

				if ret, err := fp.Seek(0, 0); ret != 0 || err != nil {
					panic(err)
				}
				if d, err := ioutil.ReadAll(fp); err == nil {
					hash := md5.Sum(d)
					c := fmt.Sprintf("%x", hash)
					fileCache[c] = true
				}
			}()

			language.code += clocFiles[file].code
			language.comments += clocFiles[file].comments
			language.blanks += clocFiles[file].blanks

			if opts.Byfile {
				clocFile := clocFiles[file]
				fmt.Printf("%-[1]*[2]s %21[3]v %14[4]v %14[5]v\n",
					maxPathLen, file, clocFile.blanks, clocFile.comments, clocFile.code)
			}
		}

		files := int32(len(language.files))
		if len(language.files) <= 0 {
			continue
		}

		language.printed = true

		total.total += files
		total.blanks += language.blanks
		total.comments += language.comments
		total.code += language.code
	}

	if !opts.Byfile {
		var sortedLanguages Languages
		for _, language := range languages {
			if len(language.files) != 0 && language.printed {
				fmt.Println(language)
				sortedLanguages = append(sortedLanguages, *language)
			}
		}
		sort.Sort(sortedLanguages)

		for _, language := range sortedLanguages {
			fmt.Printf("%-27v %6v %14v %14v %14v\n",
				language.name, len(language.files), language.blanks, language.comments, language.code)
		}
	}

	if opts.Byfile {
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
		fmt.Printf("%-[1]*[2]v %6[3]v %14[4]v %14[5]v %14[6]v\n",
			maxPathLen, "TOTAL", total.total, total.blanks, total.comments, total.code)
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	} else {
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
		fmt.Printf("%-27v %6v %14v %14v %14v\n",
			"TOTAL", total.total, total.blanks, total.comments, total.code)
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	}
}
