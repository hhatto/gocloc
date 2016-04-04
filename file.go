package main

type ClocFile struct {
	Code     int32  `xml:"code,attr"`
	Comments int32  `xml:"comment,attr"`
	Blanks   int32  `xml:"blank,attr"`
	Name     string `xml:"name,attr"`
	Lang     string `xml:"language,attr"`
}

type ClocFiles []ClocFile

func (cf ClocFiles) Len() int {
	return len(cf)
}
func (cf ClocFiles) Swap(i, j int) {
	cf[i], cf[j] = cf[j], cf[i]
}
func (cf ClocFiles) Less(i, j int) bool {
	return cf[i].Code > cf[j].Code
}
