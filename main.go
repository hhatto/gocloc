package main

import (
	"fmt"
	"os"
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

func main() {
	// parse command line options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gocloc"
	parser.Usage = "[OPTIONS] PATH[...]"

	paths, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	if opts.ShowLang {
		PrintDefinedLanguages()
		return
	}

	if len(paths) <= 0 {
		parser.WriteHelp(os.Stdout)
		return
	}

	// setup option for exclude extensions
	ExcludeExts = make(map[string]struct{})
	for _, ext := range strings.Split(opts.ExcludeExt, ",") {
		e, ok := Exts[ext]
		if ok {
			ExcludeExts[e] = struct{}{}
		} else {
			ExcludeExts[ext] = struct{}{}
		}
	}

	// setup option for not match directory
	if opts.NotMatchDir != "" {
		reNotMatchDir = regexp.MustCompile(opts.NotMatchDir)
	}
	if opts.MatchDir != "" {
		reMatchDir = regexp.MustCompile(opts.MatchDir)
	}

	// value for language result
	languages := GetDefinedLanguages()

	// setup option for include languages
	IncludeLangs = make(map[string]struct{})
	for _, lang := range strings.Split(opts.IncludeLang, ",") {
		if _, ok := languages[lang]; ok {
			IncludeLangs[lang] = struct{}{}
		}
	}

	total := NewLanguage("TOTAL", []string{}, "", "")
	num, maxPathLen := getAllFiles(paths, languages)
	headerLen := 28
	header := LANG_HEADER
	clocFiles := make(map[string]*ClocFile, num)

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

	for _, language := range languages {
		if language.printed {
			continue
		}

		for _, file := range language.files {
			cf := analyzeFile(file, language)
			cf.Lang = language.name

			language.code += cf.Code
			language.comments += cf.Comments
			language.blanks += cf.Blanks
			clocFiles[file] = cf
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
			t := XMLTotalFiles{
				Code:    total.code,
				Comment: total.comments,
				Blank:   total.blanks,
			}
			f := &XMLResultFiles{
				Files: sortedFiles,
				Total: t,
			}
			xmlResult := XMLResult{
				XMLFiles: f,
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

		switch opts.OutputType {
		case OutputTypeClocXml:
			var langs []ClocLanguage
			for _, language := range sortedLanguages {
				c := ClocLanguage{
					Name:       language.name,
					FilesCount: int32(len(language.files)),
					Code:       language.code,
					Comments:   language.comments,
					Blanks:     language.blanks,
				}
				langs = append(langs, c)
			}
			t := XMLTotalLanguages{
				Code:     total.code,
				Comment:  total.comments,
				Blank:    total.blanks,
				SumFiles: total.total,
			}
			f := &XMLResultLanguages{
				Languages: langs,
				Total:     t,
			}
			xmlResult := XMLResult{
				XMLLanguages: f,
			}
			xmlResult.Encode()
		default:
			for _, language := range sortedLanguages {
				fmt.Printf("%-27v %6v %14v %14v %14v\n",
					language.name, len(language.files), language.blanks, language.comments, language.code)
			}
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
