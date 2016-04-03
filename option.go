package main

type Options struct {
	Byfile  bool   `long:"by-file" description:"Report results for every source file encountered."`
	SortTag string `long:"sort" default:"code" description:"sort based on a certain column"`
}

var opts Options
