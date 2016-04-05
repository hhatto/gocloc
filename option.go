package main

import "regexp"

type Options struct {
	Byfile      bool   `long:"by-file" description:"report results for every source file encountered."`
	SortTag     string `long:"sort" default:"code" description:"sort based on a certain column"`
	OutputType  string `long:"output-type" default:"default" description:"output type [values: default,cloc-xml,sloccount]"`
	ExcludeExt  string `long:"exclude-ext" description:"exclude file name extensions (separated commas)"`
	NotMatchDir string `long:"not-match-d" description:"exclude dir name (regex)"`
	Debug       bool   `long:"debug" description:"dump debug log for developer"`
}

const OutputTypeDefault string = "default"
const OutputTypeClocXml string = "cloc-xml"
const OutputTypeSloccount string = "sloccount"

var opts Options
var ExcludeExts map[string]struct{}
var reNotMatchDir *regexp.Regexp
