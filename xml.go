package main

import (
	"encoding/xml"
	"fmt"
)

type XmlTotal struct {
	Code    int32 `xml:"code"`
	Comment int32 `xml:"comment"`
	Blank   int32 `xml:"blank"`
}
type XmlResultFiles struct {
	Files []ClocFile `xml:"file"`
	Total XmlTotal   `xml:"total"`
}
type XmlResult struct {
	XMLName xml.Name `xml:"results"`
	//XmlHeader XmlResultHeader
	XmlFiles XmlResultFiles `xml:"files"`
}

func (x *XmlResult) Encode() {
	if output, err := xml.MarshalIndent(x, "", "  "); err == nil {
		fmt.Printf(xml.Header)
		fmt.Println(string(output))
	}
}
