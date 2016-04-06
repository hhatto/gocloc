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
			inComments += 1
		} else if section == commentEnd {
			if inComments != 0 {
				inComments -= 1
			}
		}
	}
	return inComments != 0
}

func checkMD5Sum(filename string) (ignore bool) {
	fp, err := os.Open(filename)
	if err != nil {
		// because don't open file
		fmt.Printf("os.Open() error. err=[%v]\n", err)
		return true
	}
	defer fp.Close()

	// uniq file detect & ignore
	d, err := ioutil.ReadAll(fp)
	if err != nil {
		// because don't read file
		fmt.Printf("ioutil.ReadAll() error. err=[%v]\n", err)
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
	reVCS := regexp.MustCompile("(.bzr|.cvs|.hg|.git|.svn)")
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

			if strings.HasPrefix(root, ".") || strings.HasPrefix(root, "./") {
				p = "./" + p
			}
			if ext, ok := getFileType(p); ok {
				if targetExt, ok := Exts[ext]; ok {
					// check exclude extension
					if _, ok := ExcludeExts[targetExt]; ok {
						return nil
					}

					ignore := checkMD5Sum(p)
					if !opts.SkipUniqueness && ignore {
						if opts.Debug {
							fmt.Printf("[ignore=%v] find same md5\n", p)
						}
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
