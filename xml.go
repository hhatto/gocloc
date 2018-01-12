package main

import (
	"encoding/xml"
	"fmt"
)

type XMLTotal struct {
	Code    int32 `xml:"code,attr"`
	Comment int32 `xml:"comment,attr"`
	Blank   int32 `xml:"blank,attr"`
}
type XMLResultFiles struct {
	Files []ClocFile `xml:"file"`
	Total XMLTotal   `xml:"total"`
}
type XMLResult struct {
	XMLName  xml.Name       `xml:"results"`
	XMLFiles XMLResultFiles `xml:"files"`
}

func (x *XMLResult) Encode() {
	if output, err := xml.MarshalIndent(x, "", "  "); err == nil {
		fmt.Printf(xml.Header)
		fmt.Println(string(output))
	}
}
