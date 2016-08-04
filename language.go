package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/generaltso/linguist"
)

type Language struct {
	name         string
	lineComment  string
	multiLine    string
	multiLineEnd string
	files        []string
	code         int32
	comments     int32
	blanks       int32
	lines        int32
	total        int32
	printed      bool
}
type Languages []Language

var reShebangEnv *regexp.Regexp = regexp.MustCompile("^#! *(\\S+/env) ([a-zA-Z]+)")
var reShebangLang *regexp.Regexp = regexp.MustCompile("^#! *[.a-zA-Z/]+/([a-zA-Z]+)")

func (ls Languages) Len() int {
	return len(ls)
}
func (ls Languages) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}
func (ls Languages) Less(i, j int) bool {
	return ls[i].code > ls[j].code
}

var Exts map[string]string = map[string]string{
	"as":          "ActionScript",
	"Ant":         "Ant",
	"asm":         "Assembly",
	"S":           "Assembly",
	"s":           "Assembly",
	"awk":         "Awk",
	"bat":         "Batch",
	"btm":         "Batch",
	"cbl":         "COBOL",
	"cmd":         "Batch",
	"bash":        "BASH",
	"sh":          "Bourne Shell",
	"c":           "C",
	"csh":         "C Shell",
	"ec":          "C",
	"erl":         "Erlang",
	"hrl":         "Erlang",
	"pgc":         "C",
	"capnp":       "Cap'n Proto",
	"cs":          "C#",
	"clj":         "Clojure",
	"coffee":      "CoffeeScript",
	"cfm":         "ColdFusion",
	"cfc":         "ColdFusion CFScript",
	"cmake":       "CMake",
	"cc":          "C++",
	"cpp":         "C++",
	"cxx":         "C++",
	"pcc":         "C++",
	"c++":         "C++",
	"cr":          "Crystal",
	"css":         "CSS",
	"cu":          "CUDA",
	"d":           "D",
	"dart":        "Dart",
	"dtrace":      "DTrace",
	"dts":         "Device Tree",
	"dtsi":        "Device Tree",
	"elm":         "Elm",
	"el":          "LISP",
	"exp":         "Expect",
	"ex":          "Elixir",
	"exs":         "Elixir",
	"F#":          "F#",   // deplicated F#/GLSL
	"GLSL":        "GLSL", // both use ext '.fs'
	"vs":          "GLSL",
	"shader":      "HLSL",
	"cg":          "HLSL",
	"cginc":       "HLSL",
	"hlsl":        "HLSL",
	"lisp":        "LISP",
	"lsp":         "LISP",
	"lua":         "Lua",
	"sc":          "LISP",
	"f":           "FORTRAN Legacy",
	"f77":         "FORTRAN Legacy",
	"for":         "FORTRAN Legacy",
	"ftn":         "FORTRAN Legacy",
	"pfo":         "FORTRAN Legacy",
	"f90":         "FORTRAN Modern",
	"f95":         "FORTRAN Modern",
	"f03":         "FORTRAN Modern",
	"f08":         "FORTRAN Modern",
	"go":          "Go",
	"groovy":      "Groovy",
	"gradle":      "Groovy",
	"h":           "C Header",
	"hs":          "Haskell",
	"hpp":         "C++ Header",
	"hh":          "C++ Header",
	"html":        "HTML",
	"hx":          "Haxe",
	"hxx":         "C++ Header",
	"il":          "SKILL",
	"ipynb":       "Jupyter Notebook",
	"jai":         "JAI",
	"java":        "Java",
	"js":          "JavaScript",
	"jl":          "Julia",
	"json":        "JSON",
	"jsx":         "JSX",
	"kt":          "Kotlin",
	"lds":         "LD Script",
	"less":        "LESS",
	"Objective-C": "Objective-C", // deplicated Obj-C/Matlab
	"Matlab":      "MATLAB",      // both use ext '.m'
	"md":          "Markdown",
	"markdown":    "Markdown",
	"ML":          "OCaml",
	"ml":          "OCaml",
	"mli":         "OCaml",
	"mll":         "OCaml",
	"mly":         "OCaml",
	"mm":          "Objective-C++",
	"maven":       "Maven",
	"makefile":    "Makefile",
	"mustache":    "Mustache",
	"m4":          "M4",
	"l":           "lex",
	"nim":         "Nim",
	"php":         "PHP",
	"pas":         "Pascal",
	"PL":          "Perl",
	"pl":          "Perl",
	"pm":          "Perl",
	"plan9sh":     "Plan9 Shell",
	"pony":        "Pony",
	"ps1":         "PowerShell",
	"text":        "Plain Text",
	"txt":         "Plain Text",
	"polly":       "Polly",
	"proto":       "Protocol Buffers",
	"py":          "Python",
	"pxd":         "Cython",
	"pyx":         "Cython",
	"r":           "R",
	"R":           "R",
	"Rmd":         "RMarkdown",
	"rake":        "Ruby",
	"rb":          "Ruby",
	"rkt":         "Racket",
	"rhtml":       "Ruby HTML",
	"rs":          "Rust",
	"sass":        "Sass",
	"scala":       "Scala",
	"scss":        "Sass",
	"scm":         "Scheme",
	"sed":         "sed",
	"sml":         "Standard ML",
	"sql":         "SQL",
	"swift":       "Swift",
	"t":           "Terra",
	"tex":         "TeX",
	"thy":         "Isabelle",
	"sty":         "TeX",
	"tcl":         "Tcl/Tk",
	"toml":        "TOML",
	"ts":          "TypeScript",
	"mat":         "Unity-Prefab",
	"prefab":      "Unity-Prefab",
	"Coq":         "Coq",
	"Verilog":     "Verilog",
	"csproj":      "MSBuild script",
	"vcproj":      "MSBuild script",
	"vim":         "VimL",
	"xml":         "XML",
	"XML":         "XML",
	"xsd":         "XSD",
	"xsl":         "XSLT",
	"xslt":        "XSLT",
	"wxs":         "WiX",
	"yaml":        "YAML",
	"yml":         "YAML",
	"y":           "Yacc",
	"zsh":         "Zsh",
}

var shebang2ext map[string]string = map[string]string{
	"gosh":    "scm",
	"make":    "make",
	"perl":    "pl",
	"rc":      "plan9sh",
	"python":  "py",
	"ruby":    "rb",
	"escript": "erl",
}

func getShebang(line string) (shebangLang string, ok bool) {
	if reShebangEnv.MatchString(line) {
		ret := reShebangEnv.FindAllStringSubmatch(line, -1)
		if len(ret[0]) == 3 {
			shebangLang = ret[0][2]
			if sl, ok := shebang2ext[shebangLang]; ok {
				return sl, ok
			}
			return shebangLang, true
		}
	}

	if reShebangLang.MatchString(line) {
		ret := reShebangLang.FindAllStringSubmatch(line, -1)
		if len(ret[0]) >= 2 {
			shebangLang = ret[0][1]
			if sl, ok := shebang2ext[shebangLang]; ok {
				return sl, ok
			}
			return shebangLang, true
		}
	}

	return "", false
}

func getFileTypeByShebang(path string) (shebangLang string, ok bool) {
	line := ""
	func() {
		fp, err := os.Open(path)
		if err != nil {
			return // ignore error
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		for scanner.Scan() {
			l := scanner.Text()
			line = strings.TrimSpace(l)
			break
		}
	}()

	shebangLang, ok = getShebang(line)
	return shebangLang, ok
}

func getFileType(path string) (ext string, ok bool) {
	ext = filepath.Ext(path)
	base := filepath.Base(path)

	switch ext {
	case ".m", ".v", ".fs":
		// TODO: this is slow. parallelize...
		hints := linguist.LanguageHints(path)
		cont, err := getContents(path)
		if err != nil {
			return "", false
		}
		lang := linguist.LanguageByContents(cont, hints)
		if opts.Debug {
			fmt.Printf("path=%v, lang=%v\n", path, lang)
		}
		return lang, true
	}

	switch base {
	case "CMakeLists.txt":
		return "cmake", true
	case "configure.ac":
		return "m4", true
	case "Makefile.am":
		return "makefile", true
	case "build.xml":
		return "Ant", true
	case "pom.xml":
		return "maven", true
	}

	switch strings.ToLower(base) {
	case "makefile":
		return "makefile", true
	case "rebar": // skip
		return "", false
	}

	shebangLang, ok := getFileTypeByShebang(path)
	if ok {
		return shebangLang, true
	}

	if len(ext) >= 2 {
		return ext[1:], true
	}
	return ext, ok
}

func NewLanguage(name, lineComment, multiLine, multiLineEnd string) *Language {
	return &Language{
		name:         name,
		lineComment:  lineComment,
		multiLine:    multiLine,
		multiLineEnd: multiLineEnd,
		files:        []string{},
	}
}

func lang2exts(lang string) (exts string) {
	es := []string{}
	for ext, l := range Exts {
		if lang == l {
			switch lang {
			case "Objective-C", "MATLAB":
				ext = "m"
			case "F#":
				ext = "fs"
			case "GLSL":
				if ext == "GLSL" {
					ext = "fs"
				}
			}
			es = append(es, ext)
		}
	}
	return strings.Join(es, ", ")
}

func PrintDefinitionLanguages() {
	printLangs := []string{}
	for _, lang := range GetDefinitionLanguages() {
		printLangs = append(printLangs, lang.name)
	}
	sort.Strings(printLangs)
	for _, lang := range printLangs {
		fmt.Printf("%-30v (%s)\n", lang, lang2exts(lang))
	}
}

func GetDefinitionLanguages() map[string]*Language {
	// define languages

	return map[string]*Language{
		"ActionScript":        NewLanguage("ActionScript", "//", "/*", "*/"),
		"Ant":                 NewLanguage("Ant", "<!--", "<!--", "-->"),
		"Assembly":            NewLanguage("Assembly", "//,;,#,@,|,!", "/*", "*/"),
		"Awk":                 NewLanguage("Awk", "#", "", ""),
		"Batch":               NewLanguage("Batch", "REM,rem", "", ""),
		"BASH":                NewLanguage("BASH", "#", "", ""),
		"C":                   NewLanguage("C", "//", "/*", "*/"),
		"C Header":            NewLanguage("C Header", "//", "/*", "*/"),
		"C Shell":             NewLanguage("C Shell", "#", "", ""),
		"Cap'n Proto":         NewLanguage("Cap'n Proto", "#", "", ""),
		"C#":                  NewLanguage("C#", "//", "/*", "*/"),
		"Clojure":             NewLanguage("Clojure", ",#,#_", "", ""),
		"COBOL":               NewLanguage("COBOL", "*,/", "", ""),
		"CoffeeScript":        NewLanguage("CoffeeScript", "#", "###", "###"),
		"Coq":                 NewLanguage("Coq", "(*", "(*", "*)"),
		"ColdFusion":          NewLanguage("ColdFusion", "", "<!---", "--->"),
		"ColdFusion CFScript": NewLanguage("ColdFusion CFScript", "//", "/*", "*/"),
		"CMake":               NewLanguage("CMake", "#", "", ""),
		"C++":                 NewLanguage("C++", "//", "/*", "*/"),
		"C++ Header":          NewLanguage("C++ Header", "//", "/*", "*/"),
		"Crystal":             NewLanguage("Crystal", "#", "", ""),
		"CSS":                 NewLanguage("CSS", "//", "/*", "*/"),
		"Cython":              NewLanguage("Cython", "#", "\"\"\"", "\"\"\""),
		"CUDA":                NewLanguage("CUDA", "//", "/*", "*/"),
		"D":                   NewLanguage("D", "//", "/*", "*/"),
		"Dart":                NewLanguage("Dart", "//", "/*", "*/"),
		"DTrace":              NewLanguage("DTrace", "#", "/*", "*/"),
		"Device Tree":         NewLanguage("Device Tree", "//", "/*", "*/"),
		"Elm":                 NewLanguage("Elm", "--", "{-", "-}"),
		"Elixir":              NewLanguage("Elixir", "#", "", ""),
		"Erlang":              NewLanguage("Erlang", "%", "", ""),
		"Expect":              NewLanguage("Expect", "#", "", ""),
		"F#":                  NewLanguage("F#", "(*", "(*", "*)"),
		"Lua":                 NewLanguage("Lua", "--", "--[[", "]]"),
		"LISP":                NewLanguage("LISP", ";;", "#|", "|#"),
		"FORTRAN Legacy":      NewLanguage("FORTRAN Legacy", "c,C,!,*", "", ""),
		"FORTRAN Modern":      NewLanguage("FORTRAN Modern", "!", "", ""),
		"GLSL":                NewLanguage("GLSL", "//", "/*", "*/"),
		"Go":                  NewLanguage("Go", "//", "/*", "*/"),
		"Groovy":              NewLanguage("Groovy", "//", "/*", "*/"),
		"Haskell":             NewLanguage("Haskell", "--", "", ""),
		"Haxe":                NewLanguage("Haxe", "//", "/*", "*/"),
		"HLSL":                NewLanguage("HLSL", "//", "/*", "*/"),
		"HTML":                NewLanguage("HTML", "//,<!--", "<!--", "-->"),
		"SKILL":               NewLanguage("SKILL", ";", "/*", "*/"),
		"JAI":                 NewLanguage("JAI", "//", "/*", "*/"),
		"Java":                NewLanguage("Java", "//", "/*", "*/"),
		"JavaScript":          NewLanguage("JavaScript", "//", "/*", "*/"),
		"Julia":               NewLanguage("Julia", "#", "#:=", ":=#"),
		"Jupyter Notebook":    NewLanguage("Jupyter Notebook", "#", "", ""),
		"JSON":                NewLanguage("JSON", "", "", ""),
		"JSX":                 NewLanguage("JSX", "//", "/*", "*/"),
		"Kotlin":              NewLanguage("Kotlin", "//", "/*", "*/"),
		"LD Script":           NewLanguage("LD Script", "//", "/*", "*/"),
		"LESS":                NewLanguage("LESS", "//", "/*", "*/"),
		"Objective-C":         NewLanguage("Objective-C", "//", "/*", "*/"),
		"Markdown":            NewLanguage("Markdown", "", "", ""),
		"OCaml":               NewLanguage("OCaml", "", "(*", "*)"),
		"Objective-C++":       NewLanguage("Objective-C++", "//", "/*", "*/"),
		"Makefile":            NewLanguage("Makefile", "#", "", ""),
		"MATLAB":              NewLanguage("MATLAB", "%", "%{", "}%"),
		"Maven":               NewLanguage("Maven", "<!--", "<!--", "-->"),
		"Mustache":            NewLanguage("Mustache", "", "{{!", "}}"),
		"M4":                  NewLanguage("M4", "#", "", ""),
		"Nim":                 NewLanguage("Nim", "#", "#[", "]#"),
		"lex":                 NewLanguage("lex", "", "/*", "*/"),
		"PHP":                 NewLanguage("PHP", "#,//", "/*", "*/"),
		"Pascal":              NewLanguage("Pascal", "//,(*", "{", ")"),
		"Perl":                NewLanguage("Perl", "#", ":=", ":=cut"),
		"Plain Text":          NewLanguage("Plain Text", "", "", ""),
		"Plan9 Shell":         NewLanguage("Plan9 Shell", "#", "", ""),
		"Pony":                NewLanguage("Pony", "//", "/*", "*/"),
		"PowerShell":          NewLanguage("PowerShell", "#", "<#", "#>"),
		"Polly":               NewLanguage("Polly", "<!--", "<!--", "-->"),
		"Protocol Buffers":    NewLanguage("Protocol Buffers", "//", "", ""),
		"Python":              NewLanguage("Python", "#", "\"\"\"", "\"\"\""),
		"R":                   NewLanguage("R", "#", "", ""),
		"RMarkdown":           NewLanguage("RMarkdown", "", "", ""),
		"Racket":              NewLanguage("Racket", ";", "#|", "|#"),
		"Ruby":                NewLanguage("Ruby", "#", ":=begin", ":=end"),
		"Ruby HTML":           NewLanguage("Ruby HTML", "<!--", "<!--", "-->"),
		"Rust":                NewLanguage("Rust", "//,///,//!", "/*", "*/"),
		"Scala":               NewLanguage("Scala", "//", "/*", "*/"),
		"Sass":                NewLanguage("Sass", "//", "/*", "*/"),
		"Scheme":              NewLanguage("Scheme", ";", "#|", "|#"),
		"sed":                 NewLanguage("sed", "#", "", ""),
		"Bourne Shell":        NewLanguage("Bourne Shell", "#", "", ""),
		"Standard ML":         NewLanguage("Standard ML", "", "(*", "*)"),
		"SQL":                 NewLanguage("SQL", "--", "/*", "*/"),
		"Swift":               NewLanguage("Swift", "//", "/*", "*/"),
		"Terra":               NewLanguage("Terra", "--", "--[[", "]]"),
		"TeX":                 NewLanguage("TeX", "%", "", ""),
		"Isabelle":            NewLanguage("Isabelle", "", "(*", "*)"),
		"Tcl/Tk":              NewLanguage("Tcl/Tk", "#", "", ""),
		"TOML":                NewLanguage("TOML", "#", "", ""),
		"TypeScript":          NewLanguage("TypeScript", "//", "/*", "*/"),
		"Unity-Prefab":        NewLanguage("Unity-Prefab", "", "", ""),
		"MSBuild script":      NewLanguage("MSBuild script", "<!--", "<!--", "-->"),
		"Verilog":             NewLanguage("Verilog", "//", "/*", "*/"),
		"VimL":                NewLanguage("VimL", "\"", "", ""),
		"WiX":                 NewLanguage("WiX", "<!--", "<!--", "-->"),
		"XML":                 NewLanguage("XML", "<!--", "<!--", "-->"),
		"XSLT":                NewLanguage("XSLT", "<!--", "<!--", "-->"),
		"XSD":                 NewLanguage("XSD", "<!--", "<!--", "-->"),
		"YAML":                NewLanguage("YAML", "#", "", ""),
		"Yacc":                NewLanguage("Yacc", "//", "/*", "*/"),
		"Zsh":                 NewLanguage("Zsh", "#", "", ""),
	}
}
