package gocloc

import (
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func trimBOM(line string) string {
	l := len(line)
	if l >= 3 {
		if line[0] == 0xef && line[1] == 0xbb && line[2] == 0xbf {
			trimLine := line[3:]
			return trimLine
		}
	}
	return line
}

func containsComment(line string, multiLines [][]string) bool {
	for _, comments := range multiLines {
		for _, comm := range comments {
			if strings.Contains(line, comm) {
				return true
			}
		}
	}
	return false
}

func nextRune(s string) rune {
	for _, r := range s {
		return r
	}
	return 0
}

func checkMD5Sum(path string, fileCache map[string]struct{}) (ignore bool) {
	content, err := os.ReadFile(path)
	if err != nil {
		return true
	}

	// calc md5sum
	hash := md5.Sum(content)
	c := fmt.Sprintf("%x", hash)
	if _, ok := fileCache[c]; ok {
		return true
	}

	fileCache[c] = struct{}{}
	return false
}

func isVCSDir(path string) bool {
	if len(path) > 1 && path[0] == os.PathSeparator {
		path = path[1:]
	}
	vcsDirs := []string{".bzr", ".cvs", ".hg", ".git", ".svn"}
	for _, dir := range vcsDirs {
		if strings.Contains(path, dir) {
			return true
		}
	}
	return false
}

func checkDefaultIgnore(path string, info os.FileInfo, isVCS bool) bool {
	if info.IsDir() {
		// directory is ignored
		return true
	}
	if !isVCS && isVCSDir(path) {
		// vcs file or directory is ignored
		return true
	}

	return false
}

func checkOptionMatch(path string, opts *ClocOptions) bool {
	// check match directory & file options
	if opts.ReNotMatch != nil && opts.ReNotMatch.MatchString(path) {
		return false
	}
	if opts.ReMatch != nil && !opts.ReMatch.MatchString(path) {
		return false
	}

	dir := filepath.Dir(path)
	if opts.ReNotMatchDir != nil && opts.ReNotMatchDir.MatchString(dir) {
		return false
	}

	if opts.ReMatchDir != nil && !opts.ReMatchDir.MatchString(dir) {
		return false
	}

	return true
}

// getAllFiles return all the files to be analyzed in paths.
func getAllFiles(paths []string, languages *DefinedLanguages, opts *ClocOptions) (result map[string]*Language, err error) {
	result = make(map[string]*Language, 0)
	fileCache := make(map[string]struct{})

	for _, root := range paths {
		vcsInRoot := isVCSDir(root)
		err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err)
				return nil
			}
			if ignore := checkDefaultIgnore(path, info, vcsInRoot); ignore {
				return nil
			}

			// check match & not-match directory
			if match := checkOptionMatch(path, opts); !match {
				return nil
			}

			if ext, ok := getFileType(path, opts); ok {
				if targetExt, ok := Exts[ext]; ok {
					// check exclude extension
					if _, ok := opts.ExcludeExts[targetExt]; ok {
						return nil
					}

					if len(opts.IncludeLangs) != 0 {
						if _, ok = opts.IncludeLangs[targetExt]; !ok {
							return nil
						}
					}

					if !opts.SkipDuplicated {
						ignore := checkMD5Sum(path, fileCache)
						if ignore {
							if opts.Debug {
								fmt.Printf("[ignore=%v] find same md5\n", path)
							}
							return nil
						}
					}

					if _, ok := result[targetExt]; !ok {
						definedLang := NewLanguage(
							languages.Langs[targetExt].Name,
							languages.Langs[targetExt].lineComments,
							languages.Langs[targetExt].multiLines,
						)
						if len(languages.Langs[targetExt].regexLineComments) > 0 {
							definedLang.regexLineComments = languages.Langs[targetExt].regexLineComments
						}
						result[targetExt] = definedLang
					}
					result[targetExt].Files = append(result[targetExt].Files, path)
				}
			}
			return nil
		})
	}
	return
}
