package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/hhatto/gocloc"
	flags "github.com/jessevdk/go-flags"
)

const OutputTypeDefault string = "default"
const OutputTypeClocXml string = "cloc-xml"
const OutputTypeSloccount string = "sloccount"

const FILE_HEADER string = "File"
const LANG_HEADER string = "Language"
const COMMON_HEADER string = "files          blank        comment           code"
const ROW string = "-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------" +
	"-------------------------------------------------------------------------"

var rowLen = 79

type CmdOptions struct {
	Byfile         bool   `long:"by-file" description:"report results for every encountered source file"`
	SortTag        string `long:"sort" default:"code" description:"sort based on a certain column"`
	OutputType     string `long:"output-type" default:"default" description:"output type [values: default,cloc-xml,sloccount]"`
	ExcludeExt     string `long:"exclude-ext" description:"exclude file name extensions (separated commas)"`
	IncludeLang    string `long:"include-lang" description:"include language name (separated commas)"`
	MatchDir       string `long:"match-d" description:"include dir name (regex)"`
	NotMatchDir    string `long:"not-match-d" description:"exclude dir name (regex)"`
	Debug          bool   `long:"debug" description:"dump debug log for developer"`
	SkipDuplicated bool   `long:"skip-duplicated" description:"skip duplicated files"`
	ShowLang       bool   `long:"show-lang" description:"print about all languages and extensions"`
}

func main() {
	var opts CmdOptions
	clocOpts := gocloc.NewClocOptions()
	// parse command line options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = "gocloc"
	parser.Usage = "[OPTIONS] PATH[...]"

	paths, err := flags.Parse(&opts)
	if err != nil {
		return
	}

	// value for language result
	languages := gocloc.NewDefinedLanguages()

	if opts.ShowLang {
		fmt.Println(languages.GetFormattedString())
		return
	}

	if len(paths) <= 0 {
		parser.WriteHelp(os.Stdout)
		return
	}

	// setup option for exclude extensions
	for _, ext := range strings.Split(opts.ExcludeExt, ",") {
		e, ok := gocloc.Exts[ext]
		if ok {
			clocOpts.ExcludeExts[e] = struct{}{}
		} else {
			clocOpts.ExcludeExts[ext] = struct{}{}
		}
	}

	// setup option for not match directory
	if opts.NotMatchDir != "" {
		clocOpts.ReNotMatchDir = regexp.MustCompile(opts.NotMatchDir)
	}
	if opts.MatchDir != "" {
		clocOpts.ReMatchDir = regexp.MustCompile(opts.MatchDir)
	}

	// setup option for include languages
	for _, lang := range strings.Split(opts.IncludeLang, ",") {
		if _, ok := languages.Langs[lang]; ok {
			clocOpts.IncludeLangs[lang] = struct{}{}
		}
	}

	clocOpts.Debug = opts.Debug
	clocOpts.SkipDuplicated = opts.SkipDuplicated

	processor := gocloc.NewProcessor(languages, clocOpts)
	result, err := processor.Analyze(paths)
	if err != nil {
		fmt.Printf("fail gocloc analyze. error: %v\n", err)
		return
	}

	total := result.Total
	maxPathLen := result.MaxPathLength
	clocFiles := result.Files
	clocLangs := result.Languages
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

	// write result
	if opts.Byfile {
		var sortedFiles gocloc.ClocFiles
		for _, file := range clocFiles {
			sortedFiles = append(sortedFiles, *file)
		}
		sort.Sort(sortedFiles)

		switch opts.OutputType {
		case OutputTypeClocXml:
			t := gocloc.XMLTotalFiles{
				Code:    total.Code,
				Comment: total.Comments,
				Blank:   total.Blanks,
			}
			f := &gocloc.XMLResultFiles{
				Files: sortedFiles,
				Total: t,
			}
			xmlResult := gocloc.XMLResult{
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
		var sortedLanguages gocloc.Languages
		for _, language := range clocLangs {
			if len(language.Files) != 0 {
				sortedLanguages = append(sortedLanguages, *language)
			}
		}
		sort.Sort(sortedLanguages)

		switch opts.OutputType {
		case OutputTypeClocXml:
			var langs []gocloc.ClocLanguage
			for _, language := range sortedLanguages {
				c := gocloc.ClocLanguage{
					Name:       language.Name,
					FilesCount: int32(len(language.Files)),
					Code:       language.Code,
					Comments:   language.Comments,
					Blanks:     language.Blanks,
				}
				langs = append(langs, c)
			}
			t := gocloc.XMLTotalLanguages{
				Code:     total.Code,
				Comment:  total.Comments,
				Blank:    total.Blanks,
				SumFiles: total.Total,
			}
			f := &gocloc.XMLResultLanguages{
				Languages: langs,
				Total:     t,
			}
			xmlResult := gocloc.XMLResult{
				XMLLanguages: f,
			}
			xmlResult.Encode()
		default:
			for _, language := range sortedLanguages {
				fmt.Printf("%-27v %6v %14v %14v %14v\n",
					language.Name, len(language.Files), language.Blanks, language.Comments, language.Code)
			}
		}
	}

	// write footer
	if opts.OutputType == OutputTypeDefault {
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
		if opts.Byfile {
			fmt.Printf("%-[1]*[2]v %6[3]v %14[4]v %14[5]v %14[6]v\n",
				maxPathLen, "TOTAL", total.Total, total.Blanks, total.Comments, total.Code)
		} else {
			fmt.Printf("%-27v %6v %14v %14v %14v\n",
				"TOTAL", total.Total, total.Blanks, total.Comments, total.Code)
		}
		fmt.Printf("%.[2]*[1]s\n", ROW, rowLen)
	}
}
