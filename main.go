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

const FILE_HEADER string = "File"
const LANG_HEADER string = "Language"
const COMMON_HEADER string = "files          blank        comment           code"
const ROW string = "-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------"

var rowLen = 79

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

func getAllFiles(paths []string, languages map[string]*Language) (filenum, maxPathLen int) {
	maxPathLen = 0
	for _, root := range paths {
		if _, err := os.Stat(root); err != nil {
			continue
		}
		walkCallback := func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}

			p := filepath.Join(root, rel)

			// check not-match directory
			if reNotMatchDir != nil && reNotMatchDir.MatchString(p) {
				return nil
			}

			if strings.HasPrefix(root, ".") || strings.HasPrefix(root, "./") {
				p = "./" + p
			}
			if ext, ok := getFileType(p); ok {
				if targetExt, ok := Exts[ext]; ok {
					// check exclude extension
					if _, ok := ExcludeExts[targetExt]; ok {
						return nil
					}

					languages[targetExt].files = append(languages[targetExt].files, p)
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
	// parse command line options
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

	ExcludeExts = make(map[string]struct{})
	for _, ext := range strings.Split(opts.ExcludeExt, ",") {
		e, ok := Exts[ext]
		if ok {
			ExcludeExts[e] = struct{}{}
		} else {
			ExcludeExts[ext] = struct{}{}
		}
	}

	if opts.NotMatchDir != "" {
		reNotMatchDir = regexp.MustCompile(opts.NotMatchDir)
	}

	// define languages
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
	sh := NewLanguage("Bourne Shell", "#", "", "")
	sml := NewLanguage("Standard ML", "", "(*", "*)")
	sql := NewLanguage("SQL", "--", "/*", "*/")
	swift := NewLanguage("Swift", "//", "/*", "*/")
	tex := NewLanguage("TeX", "%", "", "")
	text := NewLanguage("Plain Text", "", "", "")
	toml := NewLanguage("TOML", "#", "", "")
	type_script := NewLanguage("TypeScript", "//", "/*", "*/")
	vim_script := NewLanguage("Vim script", "\"", "", "")
	xml := NewLanguage("XML", "<!--", "<!--", "-->")
	xsl := NewLanguage("XSLT", "<!--", "<!--", "-->")
	yaml := NewLanguage("YAML", "#", "", "")
	yacc := NewLanguage("Yacc", "//", "/*", "*/")
	zsh := NewLanguage("Zsh", "#", "", "")

	// value for language result
	languages := map[string]*Language{
		"as":       action_script,
		"s":        asm,
		"awk":      awk,
		"bat":      batch,
		"bash":     bash,
		"c":        c,
		"csh":      c_shell,
		"cs":       c_sharp,
		"clj":      clojure,
		"coffee":   coffee_script,
		"cfm":      cold_fusion,
		"cfc":      cf_script,
		"cpp":      cpp,
		"css":      css,
		"d":        d,
		"dart":     dart,
		"dts":      device_tree,
		"lua":      lua,
		"lisp":     lisp,
		"f77":      fortran_legacy,
		"f90":      fortran_modern,
		"go":       golang,
		"h":        c_header,
		"hs":       haskell,
		"hpp":      cpp_header,
		"html":     html,
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
		"ml":       ocaml,
		"mm":       objective_cpp,
		"makefile": makefile,
		"mustache": mustache,
		"php":      php,
		"pas":      pascal,
		"pl":       perl,
		"text":     text,
		"polly":    polly,
		"proto":    protobuf,
		"py":       python,
		"r":        r,
		"rb":       ruby,
		"rhtml":    ruby_html,
		"rs":       rust,
		"scss":     sass,
		"sh":       sh,
		"sml":      sml,
		"sql":      sql,
		"swift":    swift,
		"tex":      tex,
		"sty":      tex,
		"toml":     toml,
		"ts":       type_script,
		"vim":      vim_script,
		"xml":      xml,
		"xsl":      xsl,
		"yaml":     yaml,
		"y":        yacc,
		"zsh":      zsh,
	}

	total := NewLanguage("TOTAL", "", "", "")
	num, maxPathLen := getAllFiles(paths, languages)
	headerLen := 28
	header := LANG_HEADER

	// write header
	if opts.Byfile {
		headerLen = maxPathLen + 1
		rowLen = maxPathLen + len(COMMON_HEADER) + 2
		header = FILE_HEADER
	}
	if opts.OutputType == OutputTypeDefault {
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
		fmt.Printf("%-[2]*[1]s %[3]s\n", header, headerLen, COMMON_HEADER)
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
				Name: file,
			}
			isInComments := false

			func() {
				fp, err := os.Open(file)
				if err != nil {
					return // ignore error
				}
				defer fp.Close()

				isFirstLine := true
				scanner := bufio.NewScanner(fp)
				for scanner.Scan() {
					line := strings.TrimSpace(scanner.Text())

					if len(strings.TrimSpace(line)) == 0 {
						clocFiles[file].Blanks += 1
						continue
					}

					if language.multi_line != "" {
						if strings.HasPrefix(line, language.multi_line) {
							isInComments = true
						} else if containComments(line, language.multi_line, language.multi_line_end) {
							isInComments = true
							clocFiles[file].Code += 1
						}
					}

					if isInComments {
						if strings.Contains(line, language.multi_line_end) {
							isInComments = false
						}
						clocFiles[file].Comments += 1
						continue
					}

					// shebang line is 'code'
					if isFirstLine && strings.HasPrefix(line, "#!") {
						clocFiles[file].Code += 1
						isFirstLine = false
						continue
					}

					if language.line_comment != "" {
						single_comments := strings.Split(language.line_comment, ",")
						isSingleComment := false
						for _, single_comment := range single_comments {
							if strings.HasPrefix(line, single_comment) {
								clocFiles[file].Comments += 1
								isSingleComment = true
								break
							}
						}
						if isSingleComment {
							continue
						}
					}

					clocFiles[file].Code += 1
				}

				// uniq file detect & ignore
				// FIXME: not used, now
				if ret, err := fp.Seek(0, 0); ret != 0 || err != nil {
					panic(err)
				}
				if d, err := ioutil.ReadAll(fp); err == nil {
					hash := md5.Sum(d)
					c := fmt.Sprintf("%x", hash)
					fileCache[c] = true
				}
			}()

			language.code += clocFiles[file].Code
			language.comments += clocFiles[file].Comments
			language.blanks += clocFiles[file].Blanks
			clocFiles[file].Lang = language.name
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

	// write result
	if opts.Byfile {
		var sortedFiles ClocFiles
		for _, file := range clocFiles {
			sortedFiles = append(sortedFiles, *file)
		}
		sort.Sort(sortedFiles)

		switch opts.OutputType {
		case OutputTypeClocXml:
			t := XmlTotal{
				Code:    total.code,
				Comment: total.comments,
				Blank:   total.blanks,
			}
			f := XmlResultFiles{
				Files: sortedFiles,
				Total: t,
			}
			xmlResult := XmlResult{
				XmlFiles: f,
			}
			xmlResult.Encode()
		case OutputTypeSloccount:
			for _, file := range sortedFiles {
				p := ""
				if strings.HasPrefix(file.Name, "./") || string(file.Name[0]) == "/" {
					splitPaths := strings.Split(file.Name, string(os.PathSeparator))
					if len(splitPaths) >= 3 {
						p = splitPaths[1]
					}
				}
				fmt.Printf("%v\t%v\t%v\t%v\n",
					file.Code, file.Lang, p, file.Name)
			}
		default:
			for _, file := range sortedFiles {
				clocFile := file
				fmt.Printf("%-[1]*[2]s %21[3]v %14[4]v %14[5]v\n",
					maxPathLen, file.Name, clocFile.Blanks, clocFile.Comments, clocFile.Code)
			}
		}
	} else {
		var sortedLanguages Languages
		for _, language := range languages {
			if len(language.files) != 0 && language.printed {
				sortedLanguages = append(sortedLanguages, *language)
			}
		}
		sort.Sort(sortedLanguages)

		for _, language := range sortedLanguages {
			fmt.Printf("%-27v %6v %14v %14v %14v\n",
				language.name, len(language.files), language.blanks, language.comments, language.code)
		}
	}

	// write footer
	if opts.OutputType == OutputTypeDefault {
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
		if opts.Byfile {
			fmt.Printf("%-[1]*[2]v %6[3]v %14[4]v %14[5]v %14[6]v\n",
				maxPathLen, "TOTAL", total.total, total.blanks, total.comments, total.code)
		} else {
			fmt.Printf("%-27v %6v %14v %14v %14v\n",
				"TOTAL", total.total, total.blanks, total.comments, total.code)
		}
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	}
}
