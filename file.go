package main

type ClocFile struct {
	name     string
	code     int32
	comments int32
	blanks   int32
	lines    int32
}

type ClocFiles []ClocFile

func (cf ClocFiles) Len() int {
	return len(cf)
}
func (cf ClocFiles) Swap(i, j int) {
	cf[i], cf[j] = cf[j], cf[i]
}
func (cf ClocFiles) Less(i, j int) bool {
	return cf[i].code > cf[j].code
}
