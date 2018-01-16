package gocloc

import (
	"encoding/xml"
	"fmt"
)

type XMLTotal struct {
	Code    int32 `xml:"code"`
	Comment int32 `xml:"comment"`
	Blank   int32 `xml:"blank"`
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
