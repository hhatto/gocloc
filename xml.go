package gocloc

import (
	"encoding/xml"
	"fmt"
)

// XMLResultType is the result type in XML format.
type XMLResultType int8

const (
	// XMLResultWithLangs is the result type for each language in XML format
	XMLResultWithLangs XMLResultType = iota
	// XMLResultWithFiles is the result type for each file in XML format
	XMLResultWithFiles
)

// XMLTotalLanguages is the total result in XML format.
type XMLTotalLanguages struct {
	SumFiles int32 `xml:"sum_files,attr"`
	Code     int32 `xml:"code,attr"`
	Comment  int32 `xml:"comment,attr"`
	Blank    int32 `xml:"blank,attr"`
}

// XMLResultLanguages stores the results in XML format.
type XMLResultLanguages struct {
	Languages []ClocLanguage    `xml:"language"`
	Total     XMLTotalLanguages `xml:"total"`
}

// XMLTotalFiles is the total result per file in XML format.
type XMLTotalFiles struct {
	Code    int32 `xml:"code,attr"`
	Comment int32 `xml:"comment,attr"`
	Blank   int32 `xml:"blank,attr"`
}

// XMLResultFiles stores per file results in XML format.
type XMLResultFiles struct {
	Files []ClocFile    `xml:"file"`
	Total XMLTotalFiles `xml:"total"`
}

// XMLResult stores the results in XML format.
type XMLResult struct {
	XMLName      xml.Name            `xml:"results"`
	XMLFiles     *XMLResultFiles     `xml:"files,omitempty"`
	XMLLanguages *XMLResultLanguages `xml:"languages,omitempty"`
}

// Encode outputs XMLResult in a human readable format.
func (x *XMLResult) Encode() {
	if output, err := xml.MarshalIndent(x, "", "  "); err == nil {
		fmt.Printf(xml.Header)
		fmt.Println(string(output))
	}
}

// NewXMLResultFromCloc returns XMLResult with default data set.
func NewXMLResultFromCloc(total *Language, sortedLanguages Languages, _ XMLResultType) *XMLResult {
	var langs []ClocLanguage
	for _, language := range sortedLanguages {
		c := ClocLanguage{
			Name:       language.Name,
			FilesCount: int32(len(language.Files)),
			Code:       language.Code,
			Comments:   language.Comments,
			Blanks:     language.Blanks,
		}
		langs = append(langs, c)
	}
	t := XMLTotalLanguages{
		Code:     total.Code,
		Comment:  total.Comments,
		Blank:    total.Blanks,
		SumFiles: total.Total,
	}
	f := &XMLResultLanguages{
		Languages: langs,
		Total:     t,
	}

	return &XMLResult{
		XMLLanguages: f,
	}
}
