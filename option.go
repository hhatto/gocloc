package main

type Options struct {
	Byfile     bool   `long:"by-file" description:"report results for every source file encountered."`
	SortTag    string `long:"sort" default:"code" description:"sort based on a certain column"`
	OutputType string `long:"output-type" default:"default" description:"output type [values: default,cloc-xml,sloccount]"`
}

var opts Options
