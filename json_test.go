package gocloc

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestOutputJSON(t *testing.T) {
	total := &Language{}
	files := []ClocFile{
		{Name: "one.go", Lang: "Go"},
		{Name: "two.go", Lang: "Go"},
	}
	jsonResult := NewJSONFilesResultFromCloc(total, files)
	if jsonResult.Files[0].Name != "one.go" {
		t.Errorf("invalid result. Name: one.go")
	}
	if jsonResult.Files[1].Name != "two.go" {
		t.Errorf("invalid result. Name: two.go")
	}
	if jsonResult.Files[1].Lang != "Go" {
		t.Errorf("invalid result. lang: Go")
	}

	// check output json text
	buf, err := json.Marshal(jsonResult)
	if err != nil {
		fmt.Println(err)
		t.Errorf("json marshal error")
	}

	actualJSONText := `{"files":[{"code":0,"comment":0,"blank":0,"name":"one.go","language":"Go"},{"code":0,"comment":0,"blank":0,"name":"two.go","language":"Go"}],"total":{"files":0,"code":0,"comment":0,"blank":0}}`
	resultJSONText := string(buf)
	if actualJSONText != resultJSONText {
		t.Errorf("invalid result. '%s'", resultJSONText)
	}
}
