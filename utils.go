package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var fileCache map[string]struct{}

func containComments(line, commentStart, commentEnd string) bool {
	inComments := 0
	for i := 0; i < len(line)/len(commentStart); i += len(commentStart) {
		section := line[i : i+len(commentStart)]

		if section == commentStart {
			inComments++
		} else if section == commentEnd {
			if inComments != 0 {
				inComments--
			}
		}
	}
	return inComments != 0
}

func getContents(path string) ([]byte, error) {
	fp, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	return ioutil.ReadAll(fp)
}

func checkMD5Sum(filename string) (ignore bool) {
	d, err := getContents(filename)
	if err != nil {
		return true
	}

	// calc md5sum
	hash := md5.Sum(d)
	c := fmt.Sprintf("%x", hash)
	if _, ok := fileCache[c]; ok {
		return true
	}

	fileCache[c] = struct{}{}
	return false
}

func getAllFiles(paths []string, languages map[string]*Language) (filenum, maxPathLen int) {
	reVCS := regexp.MustCompile("\\.(bzr|cvs|hg|git|svn)")
	maxPathLen = 0
	for _, root := range paths {
		if _, err := os.Stat(root); err != nil {
			continue
		}
		vcsInRoot := reVCS.MatchString(root)
		walkCallback := func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(root, path)
			if err != nil {
				return nil
			}

			if !vcsInRoot && reVCS.MatchString(path) {
				return nil
			}

			p := filepath.Join(root, rel)

			// check not-match directory
			if reNotMatchDir != nil && reNotMatchDir.MatchString(p) {
				return nil
			}

			if reMatchDir != nil && !reMatchDir.MatchString(filepath.Dir(p)) {
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

					if len(IncludeLangs) != 0 {
						if _, ok = IncludeLangs[targetExt]; !ok {
							return nil
						}
					}

					ignore := checkMD5Sum(p)
					if !opts.SkipUniqueness && ignore {
						if opts.Debug {
							fmt.Printf("[ignore=%v] find same md5\n", p)
						}
						return nil
					}

					languages[targetExt].files = append(languages[targetExt].files, p)
					filenum++
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
