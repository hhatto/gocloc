package main

type Options struct {
	Byfile bool `long:"by-file" description:"Report results for every source file encountered."`
}

var opts Options
