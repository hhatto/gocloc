package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

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

var fileCache map[string]struct{}

func analyzeFile(filename string, language *Language) *ClocFile {
	if opts.Debug {
		fmt.Printf("filename=%v\n", filename)
	}

	clocFile := &ClocFile{
		Name: filename,
	}

	fp, err := os.Open(filename)
	if err != nil {
		return clocFile // ignore error
	}
	defer fp.Close()

	isFirstLine := true
	isInComments := false
	isInCommentsSame := false
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		lineOrg := scanner.Text()
		line := strings.TrimSpace(lineOrg)

		if len(strings.TrimSpace(line)) == 0 {
			clocFile.Blanks += 1
			if opts.Debug {
				fmt.Printf("[BLNK,cd:%d,cm:%d,bk:%d,iscm:%v] %s\n",
					clocFile.Code, clocFile.Comments, clocFile.Blanks, isInComments, lineOrg)
			}
			continue
		}

		if language.multi_line != "" {
			if strings.HasPrefix(line, language.multi_line) {
				isInComments = true
			} else if containComments(line, language.multi_line, language.multi_line_end) {
				isInComments = true
				clocFile.Code += 1
			}
		}

		if isInComments {
			if language.multi_line == language.multi_line_end {
				if strings.Count(line, language.multi_line_end) == 2 {
					isInComments = false
					isInCommentsSame = false
				} else if strings.HasPrefix(line, language.multi_line_end) {
					if isInCommentsSame {
						isInComments = false
					}
					isInCommentsSame = !isInCommentsSame
				}
			} else {
				if strings.Contains(line, language.multi_line_end) {
					isInComments = false
				}
			}
			clocFile.Comments += 1
			if opts.Debug {
				fmt.Printf("[COMM,cd:%d,cm:%d,bk:%d,iscm:%v,iscms:%v] %s\n",
					clocFile.Code, clocFile.Comments, clocFile.Blanks, isInComments, isInCommentsSame, lineOrg)
			}
			continue
		}

		// shebang line is 'code'
		if isFirstLine && strings.HasPrefix(line, "#!") {
			clocFile.Code += 1
			isFirstLine = false
			if opts.Debug {
				fmt.Printf("[CODE,cd:%d,cm:%d,bk:%d,iscm:%v] %s\n",
					clocFile.Code, clocFile.Comments, clocFile.Blanks, isInComments, lineOrg)
			}
			continue
		}

		if language.line_comment != "" {
			single_comments := strings.Split(language.line_comment, ",")
			isSingleComment := false
			for _, single_comment := range single_comments {
				if strings.HasPrefix(line, single_comment) {
					clocFile.Comments += 1
					isSingleComment = true
					break
				}
			}
			if isSingleComment {
				if opts.Debug {
					fmt.Printf("[COMM,cd:%d,cm:%d,bk:%d,iscm:%v] %s\n",
						clocFile.Code, clocFile.Comments, clocFile.Blanks, isInComments, lineOrg)
				}
				continue
			}
		}

		clocFile.Code += 1
		if opts.Debug {
			fmt.Printf("[CODE,cd:%d,cm:%d,bk:%d,iscm:%v] %s\n",
				clocFile.Code, clocFile.Comments, clocFile.Blanks, isInComments, lineOrg)
		}
	}

	// uniq file detect & ignore
	// FIXME: not used, now
	if ret, err := fp.Seek(0, 0); ret != 0 || err != nil {
		panic(err)
	}
	if d, err := ioutil.ReadAll(fp); err == nil {
		hash := md5.Sum(d)
		c := fmt.Sprintf("%x", hash)
		fileCache[c] = struct{}{}
	}

	return clocFile
}
