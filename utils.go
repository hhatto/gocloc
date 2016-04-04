package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func containComments(line, commentStart, commentEnd string) bool {
	inComments := 0
	for i := 0; i < len(line)/len(commentStart); i += len(commentStart) {
		section := line[i : i+len(commentStart)]

		if section == commentStart {
			inComments += 1
		} else if section == commentEnd {
			if inComments != 0 {
				inComments -= 1
			}
		}
	}
	return inComments != 0
}

func getAllFiles(paths []string, languages map[string]*Language) (filenum, maxPathLen int) {
	maxPathLen = 0
	for _, root := range paths {
		if _, err := os.Stat(root); err != nil {
			continue
		}
		walkCallback := func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}

			p := filepath.Join(root, rel)

			// check not-match directory
			if reNotMatchDir != nil && reNotMatchDir.MatchString(p) {
				return nil
			}

			if strings.HasPrefix(root, ".") || strings.HasPrefix(root, "./") {
				p = "./" + p
			}
			if ext, ok := getFileType(p); ok {
				if targetExt, ok := Exts[ext]; ok {
					// check exclude extension
					if _, ok := ExcludeExts[targetExt]; ok {
						return nil
					}

					languages[targetExt].files = append(languages[targetExt].files, p)
					filenum += 1
					l := len(p)
					if maxPathLen < l {
						maxPathLen = l
					}
				}
			}
			return nil
		}
		if err := filepath.Walk(root, walkCallback); err != nil {
			fmt.Println(err)
		}
	}
	return
}
