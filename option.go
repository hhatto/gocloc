package gocloc

import "regexp"

type Options struct {
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

const OutputTypeDefault string = "default"
const OutputTypeClocXml string = "cloc-xml"
const OutputTypeSloccount string = "sloccount"

var ExcludeExts map[string]struct{}
var IncludeLangs map[string]struct{}
var reNotMatchDir *regexp.Regexp
var reMatchDir *regexp.Regexp
