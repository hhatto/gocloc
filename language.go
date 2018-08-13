package gocloc

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"unicode"

	enry "gopkg.in/src-d/enry.v1"
)

// ClocLanguage is provide for xml-cloc format
type ClocLanguage struct {
	Name       string `xml:"name,attr"`
	FilesCount int32  `xml:"files_count,attr"`
	Code       int32  `xml:"code,attr"`
	Comments   int32  `xml:"comment,attr"`
	Blanks     int32  `xml:"blank,attr"`
}

type Language struct {
	Name         string
	lineComments []string
	multiLine    string
	multiLineEnd string
	Files        []string
	Code         int32
	Comments     int32
	Blanks       int32
	Total        int32
}
type Languages []Language

func (ls Languages) Len() int {
	return len(ls)
}
func (ls Languages) Swap(i, j int) {
	ls[i], ls[j] = ls[j], ls[i]
}
func (ls Languages) Less(i, j int) bool {
	if ls[i].Code == ls[j].Code {
		return ls[i].Name < ls[j].Name
	}
	return ls[i].Code > ls[j].Code
}

var reShebangEnv = regexp.MustCompile("^#! *(\\S+/env) ([a-zA-Z]+)")
var reShebangLang = regexp.MustCompile("^#! *[.a-zA-Z/]+/([a-zA-Z]+)")

var Exts = map[string]string{
	"as":          "ActionScript",
	"ada":         "Ada",
	"adb":         "Ada",
	"ads":         "Ada",
	"Ant":         "Ant",
	"adoc":        "AsciiDoc",
	"asciidoc":    "AsciiDoc",
	"asm":         "Assembly",
	"S":           "Assembly",
	"s":           "Assembly",
	"awk":         "Awk",
	"bat":         "Batch",
	"btm":         "Batch",
	"bb":          "BitBake",
	"cbl":         "COBOL",
	"cmd":         "Batch",
	"bash":        "BASH",
	"sh":          "Bourne Shell",
	"c":           "C",
	"carp":        "Carp",
	"csh":         "C Shell",
	"ec":          "C",
	"erl":         "Erlang",
	"hrl":         "Erlang",
	"pgc":         "C",
	"capnp":       "Cap'n Proto",
	"chpl":        "Chapel",
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
	"e":           "Eiffel",
	"elm":         "Elm",
	"el":          "LISP",
	"exp":         "Expect",
	"ex":          "Elixir",
	"exs":         "Elixir",
	"feature":     "Gherkin",
	"fish":        "Fish",
	"fr":          "Frege",
	"fst":         "F*",
	"F#":          "F#",   // deplicated F#/GLSL
	"GLSL":        "GLSL", // both use ext '.fs'
	"vs":          "GLSL",
	"shader":      "HLSL",
	"cg":          "HLSL",
	"cginc":       "HLSL",
	"hlsl":        "HLSL",
	"lean":        "Lean",
	"hlean":       "Lean",
	"lgt":         "Logtalk",
	"lisp":        "LISP",
	"lsp":         "LISP",
	"lua":         "Lua",
	"ls":          "LiveScript",
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
	"idr":         "Idris",
	"il":          "SKILL",
	"io":          "Io",
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
	"Objective-C": "Objective-C", // deplicated Obj-C/Matlab/Mercury
	"Matlab":      "MATLAB",      // both use ext '.m'
	"Mercury":     "Mercury",     // use ext '.m'
	"md":          "Markdown",
	"markdown":    "Markdown",
	"nix":         "Nix",
	"nsi":         "NSIS",
	"nsh":         "NSIS",
	"nu":          "Nu",
	"ML":          "OCaml",
	"ml":          "OCaml",
	"mli":         "OCaml",
	"mll":         "OCaml",
	"mly":         "OCaml",
	"mm":          "Objective-C++",
	"maven":       "Maven",
	"makefile":    "Makefile",
	"meson":       "Meson",
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
	"raml":        "RAML",
	"Rebol":       "Rebol",
	"red":         "Red",
	"Rmd":         "RMarkdown",
	"rake":        "Ruby",
	"rb":          "Ruby",
	"rkt":         "Racket",
	"rhtml":       "Ruby HTML",
	"rs":          "Rust",
	"rst":         "ReStructuredText",
	"sass":        "Sass",
	"scala":       "Scala",
	"scss":        "Sass",
	"scm":         "Scheme",
	"sed":         "sed",
	"stan":        "Stan",
	"sml":         "Standard ML",
	"sol":         "Solidity",
	"sql":         "SQL",
	"swift":       "Swift",
	"t":           "Terra",
	"tex":         "TeX",
	"thy":         "Isabelle",
	"tla":         "TLA",
	"sty":         "TeX",
	"tcl":         "Tcl/Tk",
	"toml":        "TOML",
	"ts":          "TypeScript",
	"mat":         "Unity-Prefab",
	"prefab":      "Unity-Prefab",
	"Coq":         "Coq",
	"vala":        "Vala",
	"Verilog":     "Verilog",
	"csproj":      "MSBuild script",
	"vcproj":      "MSBuild script",
	"vim":         "VimL",
	"vue":         "Vue",
	"xml":         "XML",
	"XML":         "XML",
	"xsd":         "XSD",
	"xsl":         "XSLT",
	"xslt":        "XSLT",
	"wxs":         "WiX",
	"yaml":        "YAML",
	"yml":         "YAML",
	"y":           "Yacc",
	"zep":         "Zephir",
	"zsh":         "Zsh",
}

var shebang2ext = map[string]string{
	"gosh":    "scm",
	"make":    "make",
	"perl":    "pl",
	"rc":      "plan9sh",
	"python":  "py",
	"ruby":    "rb",
	"escript": "erl",
}

func getShebang(line string) (shebangLang string, ok bool) {
	ret := reShebangEnv.FindAllStringSubmatch(line, -1)
	if ret != nil && len(ret[0]) == 3 {
		shebangLang = ret[0][2]
		if sl, ok := shebang2ext[shebangLang]; ok {
			return sl, ok
		}
		return shebangLang, true
	}

	ret = reShebangLang.FindAllStringSubmatch(line, -1)
	if ret != nil && len(ret[0]) >= 2 {
		shebangLang = ret[0][1]
		if sl, ok := shebang2ext[shebangLang]; ok {
			return sl, ok
		}
		return shebangLang, true
	}

	return "", false
}

func getFileTypeByShebang(path string) (shebangLang string, ok bool) {
	f, err := os.Open(path)
	if err != nil {
		return // ignore error
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	line, err := reader.ReadBytes('\n')
	if err != nil {
		return
	}
	line = bytes.TrimLeftFunc(line, unicode.IsSpace)

	if len(line) > 2 && line[0] == '#' && line[1] == '!' {
		return getShebang(string(line))
	}
	return
}

func getFileType(path string, opts *ClocOptions) (ext string, ok bool) {
	ext = filepath.Ext(path)
	base := filepath.Base(path)

	switch ext {
	case ".m", ".v", ".fs", ".r":
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return "", false
		}
		lang := enry.GetLanguage(path, content)
		if opts.Debug {
			fmt.Printf("path=%v, lang=%v\n", path, lang)
		}
		return lang, true
	}

	switch base {
	case "meson.build", "meson_options.txt":
		return "meson", true
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
	case "nukefile":
		return "nu", true
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

func NewLanguage(name string, lineComments []string, multiLine, multiLineEnd string) *Language {
	return &Language{
		Name:         name,
		lineComments: lineComments,
		multiLine:    multiLine,
		multiLineEnd: multiLineEnd,
		Files:        []string{},
	}
}

func lang2exts(lang string) (exts string) {
	es := []string{}
	for ext, l := range Exts {
		if lang == l {
			switch lang {
			case "Objective-C", "MATLAB", "Mercury":
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

type DefinedLanguages struct {
	Langs map[string]*Language
}

func (langs *DefinedLanguages) GetFormattedString() string {
	var buf bytes.Buffer
	printLangs := []string{}
	for _, lang := range langs.Langs {
		printLangs = append(printLangs, lang.Name)
	}
	sort.Strings(printLangs)
	for _, lang := range printLangs {
		buf.WriteString(fmt.Sprintf("%-30v (%s)\n", lang, lang2exts(lang)))
	}
	return buf.String()
}

func NewDefinedLanguages() *DefinedLanguages {
	return &DefinedLanguages{
		Langs: map[string]*Language{
			"ActionScript":        NewLanguage("ActionScript", []string{"//"}, "/*", "*/"),
			"Ada":                 NewLanguage("Ada", []string{"--"}, "", ""),
			"Ant":                 NewLanguage("Ant", []string{"<!--"}, "<!--", "-->"),
			"AsciiDoc":            NewLanguage("AsciiDoc", []string{}, "", ""),
			"Assembly":            NewLanguage("Assembly", []string{"//", ";", "#", "@", "|", "!"}, "/*", "*/"),
			"Awk":                 NewLanguage("Awk", []string{"#"}, "", ""),
			"Batch":               NewLanguage("Batch", []string{"REM", "rem"}, "", ""),
			"BASH":                NewLanguage("BASH", []string{"#"}, "", ""),
			"BitBake":             NewLanguage("BitBake", []string{"#"}, "", ""),
			"C":                   NewLanguage("C", []string{"//"}, "/*", "*/"),
			"C Header":            NewLanguage("C Header", []string{"//"}, "/*", "*/"),
			"C Shell":             NewLanguage("C Shell", []string{"#"}, "", ""),
			"Cap'n Proto":         NewLanguage("Cap'n Proto", []string{"#"}, "", ""),
			"Carp":                NewLanguage("Carp", []string{";"}, "", ""),
			"C#":                  NewLanguage("C#", []string{"//"}, "/*", "*/"),
			"Chapel":              NewLanguage("Chapel", []string{"//"}, "/*", "*/"),
			"Clojure":             NewLanguage("Clojure", []string{"#", "#_"}, "", ""),
			"COBOL":               NewLanguage("COBOL", []string{"*", "/"}, "", ""),
			"CoffeeScript":        NewLanguage("CoffeeScript", []string{"#"}, "###", "###"),
			"Coq":                 NewLanguage("Coq", []string{"(*"}, "(*", "*)"),
			"ColdFusion":          NewLanguage("ColdFusion", []string{}, "<!---", "--->"),
			"ColdFusion CFScript": NewLanguage("ColdFusion CFScript", []string{"//"}, "/*", "*/"),
			"CMake":               NewLanguage("CMake", []string{"#"}, "", ""),
			"C++":                 NewLanguage("C++", []string{"//"}, "/*", "*/"),
			"C++ Header":          NewLanguage("C++ Header", []string{"//"}, "/*", "*/"),
			"Crystal":             NewLanguage("Crystal", []string{"#"}, "", ""),
			"CSS":                 NewLanguage("CSS", []string{"//"}, "/*", "*/"),
			"Cython":              NewLanguage("Cython", []string{"#"}, "\"\"\"", "\"\"\""),
			"CUDA":                NewLanguage("CUDA", []string{"//"}, "/*", "*/"),
			"D":                   NewLanguage("D", []string{"//"}, "/*", "*/"),
			"Dart":                NewLanguage("Dart", []string{"//"}, "/*", "*/"),
			"DTrace":              NewLanguage("DTrace", []string{}, "/*", "*/"),
			"Device Tree":         NewLanguage("Device Tree", []string{"//"}, "/*", "*/"),
			"Eiffel":              NewLanguage("Eiffel", []string{"--"}, "", ""),
			"Elm":                 NewLanguage("Elm", []string{"--"}, "{-", "-}"),
			"Elixir":              NewLanguage("Elixir", []string{"#"}, "", ""),
			"Erlang":              NewLanguage("Erlang", []string{"%"}, "", ""),
			"Expect":              NewLanguage("Expect", []string{"#"}, "", ""),
			"Fish":                NewLanguage("Fish", []string{"#"}, "", ""),
			"Frege":               NewLanguage("Frege", []string{"--"}, "{-", "-}"),
			"F*":                  NewLanguage("F*", []string{"(*", "//"}, "(*", "*)"),
			"F#":                  NewLanguage("F#", []string{"(*"}, "(*", "*)"),
			"Lean":                NewLanguage("Lean", []string{"--"}, "/-", "-/"),
			"Logtalk":             NewLanguage("Logtalk", []string{"%"}, "", ""),
			"Lua":                 NewLanguage("Lua", []string{"--"}, "--[[", "]]"),
			"LISP":                NewLanguage("LISP", []string{";;"}, "#|", "|#"),
			"LiveScript":          NewLanguage("LiveScript", []string{"#"}, "/*", "*/"),
			"FORTRAN Legacy":      NewLanguage("FORTRAN Legacy", []string{"c", "C", "!", "*"}, "", ""),
			"FORTRAN Modern":      NewLanguage("FORTRAN Modern", []string{"!"}, "", ""),
			"Gherkin":             NewLanguage("Gherkin", []string{"#"}, "", ""),
			"GLSL":                NewLanguage("GLSL", []string{"//"}, "/*", "*/"),
			"Go":                  NewLanguage("Go", []string{"//"}, "/*", "*/"),
			"Groovy":              NewLanguage("Groovy", []string{"//"}, "/*", "*/"),
			"Haskell":             NewLanguage("Haskell", []string{"--"}, "{-", "-}"),
			"Haxe":                NewLanguage("Haxe", []string{"//"}, "/*", "*/"),
			"HLSL":                NewLanguage("HLSL", []string{"//"}, "/*", "*/"),
			"HTML":                NewLanguage("HTML", []string{"//", "<!--"}, "<!--", "-->"),
			"Idris":               NewLanguage("Idris", []string{"--"}, "{-", "-}"),
			"Io":                  NewLanguage("Io", []string{"//", "#"}, "/*", "*/"),
			"SKILL":               NewLanguage("SKILL", []string{";"}, "/*", "*/"),
			"JAI":                 NewLanguage("JAI", []string{"//"}, "/*", "*/"),
			"Java":                NewLanguage("Java", []string{"//"}, "/*", "*/"),
			"JavaScript":          NewLanguage("JavaScript", []string{"//"}, "/*", "*/"),
			"Julia":               NewLanguage("Julia", []string{"#"}, "#:=", ":=#"),
			"Jupyter Notebook":    NewLanguage("Jupyter Notebook", []string{"#"}, "", ""),
			"JSON":                NewLanguage("JSON", []string{}, "", ""),
			"JSX":                 NewLanguage("JSX", []string{"//"}, "/*", "*/"),
			"Kotlin":              NewLanguage("Kotlin", []string{"//"}, "/*", "*/"),
			"LD Script":           NewLanguage("LD Script", []string{"//"}, "/*", "*/"),
			"LESS":                NewLanguage("LESS", []string{"//"}, "/*", "*/"),
			"Objective-C":         NewLanguage("Objective-C", []string{"//"}, "/*", "*/"),
			"Markdown":            NewLanguage("Markdown", []string{}, "", ""),
			"Nix":                 NewLanguage("Nix", []string{"#"}, "/*", "*/"),
			"NSIS":                NewLanguage("NSIS", []string{"#", ";"}, "/*", "*/"),
			"Nu":                  NewLanguage("Nu", []string{";", "#"}, "", ""),
			"OCaml":               NewLanguage("OCaml", []string{}, "(*", "*)"),
			"Objective-C++":       NewLanguage("Objective-C++", []string{"//"}, "/*", "*/"),
			"Makefile":            NewLanguage("Makefile", []string{"#"}, "", ""),
			"MATLAB":              NewLanguage("MATLAB", []string{"%"}, "%{", "}%"),
			"Mercury":             NewLanguage("Mercury", []string{"%"}, "/*", "*/"),
			"Maven":               NewLanguage("Maven", []string{"<!--"}, "<!--", "-->"),
			"Meson":               NewLanguage("Meson", []string{"#"}, "", ""),
			"Mustache":            NewLanguage("Mustache", []string{}, "{{!", "}}"),
			"M4":                  NewLanguage("M4", []string{"#"}, "", ""),
			"Nim":                 NewLanguage("Nim", []string{"#"}, "#[", "]#"),
			"lex":                 NewLanguage("lex", []string{}, "/*", "*/"),
			"PHP":                 NewLanguage("PHP", []string{"#", "//"}, "/*", "*/"),
			"Pascal":              NewLanguage("Pascal", []string{"//"}, "{", ")"),
			"Perl":                NewLanguage("Perl", []string{"#"}, ":=", ":=cut"),
			"Plain Text":          NewLanguage("Plain Text", []string{}, "", ""),
			"Plan9 Shell":         NewLanguage("Plan9 Shell", []string{"#"}, "", ""),
			"Pony":                NewLanguage("Pony", []string{"//"}, "/*", "*/"),
			"PowerShell":          NewLanguage("PowerShell", []string{"#"}, "<#", "#>"),
			"Polly":               NewLanguage("Polly", []string{"<!--"}, "<!--", "-->"),
			"Protocol Buffers":    NewLanguage("Protocol Buffers", []string{"//"}, "", ""),
			"Python":              NewLanguage("Python", []string{"#"}, "\"\"\"", "\"\"\""),
			"R":                   NewLanguage("R", []string{"#"}, "", ""),
			"Rebol":               NewLanguage("Rebol", []string{";"}, "", ""),
			"Red":                 NewLanguage("Red", []string{";"}, "", ""),
			"RMarkdown":           NewLanguage("RMarkdown", []string{}, "", ""),
			"RAML":                NewLanguage("RAML", []string{"#"}, "", ""),
			"Racket":              NewLanguage("Racket", []string{";"}, "#|", "|#"),
			"ReStructuredText":    NewLanguage("ReStructuredText", []string{}, "", ""),
			"Ruby":                NewLanguage("Ruby", []string{"#"}, ":=begin", ":=end"),
			"Ruby HTML":           NewLanguage("Ruby HTML", []string{"<!--"}, "<!--", "-->"),
			"Rust":                NewLanguage("Rust", []string{"//", "///", "//!"}, "/*", "*/"),
			"Scala":               NewLanguage("Scala", []string{"//"}, "/*", "*/"),
			"Sass":                NewLanguage("Sass", []string{"//"}, "/*", "*/"),
			"Scheme":              NewLanguage("Scheme", []string{";"}, "#|", "|#"),
			"sed":                 NewLanguage("sed", []string{"#"}, "", ""),
			"Stan":                NewLanguage("Stan", []string{"//"}, "/*", "*/"),
			"Solidity":            NewLanguage("Solidity", []string{"//"}, "/*", "*/"),
			"Bourne Shell":        NewLanguage("Bourne Shell", []string{"#"}, "", ""),
			"Standard ML":         NewLanguage("Standard ML", []string{}, "(*", "*)"),
			"SQL":                 NewLanguage("SQL", []string{"--"}, "/*", "*/"),
			"Swift":               NewLanguage("Swift", []string{"//"}, "/*", "*/"),
			"Terra":               NewLanguage("Terra", []string{"--"}, "--[[", "]]"),
			"TeX":                 NewLanguage("TeX", []string{"%"}, "", ""),
			"Isabelle":            NewLanguage("Isabelle", []string{}, "(*", "*)"),
			"TLA":                 NewLanguage("TLA", []string{"/*"}, "(*", "*)"),
			"Tcl/Tk":              NewLanguage("Tcl/Tk", []string{"#"}, "", ""),
			"TOML":                NewLanguage("TOML", []string{"#"}, "", ""),
			"TypeScript":          NewLanguage("TypeScript", []string{"//"}, "/*", "*/"),
			"Unity-Prefab":        NewLanguage("Unity-Prefab", []string{}, "", ""),
			"MSBuild script":      NewLanguage("MSBuild script", []string{"<!--"}, "<!--", "-->"),
			"Vala":                NewLanguage("Vala", []string{"//"}, "/*", "*/"),
			"Verilog":             NewLanguage("Verilog", []string{"//"}, "/*", "*/"),
			"VimL":                NewLanguage("VimL", []string{`"`}, "", ""),
			"Vue":                 NewLanguage("Vue", []string{"<!--"}, "<!--", "-->"),
			"WiX":                 NewLanguage("WiX", []string{"<!--"}, "<!--", "-->"),
			"XML":                 NewLanguage("XML", []string{"<!--"}, "<!--", "-->"),
			"XSLT":                NewLanguage("XSLT", []string{"<!--"}, "<!--", "-->"),
			"XSD":                 NewLanguage("XSD", []string{"<!--"}, "<!--", "-->"),
			"YAML":                NewLanguage("YAML", []string{"#"}, "", ""),
			"Yacc":                NewLanguage("Yacc", []string{"//"}, "/*", "*/"),
			"Zephir":              NewLanguage("Zephir", []string{"//"}, "/*", "*/"),
			"Zsh":                 NewLanguage("Zsh", []string{"#"}, "", ""),
		},
	}
}
