package gocloc

import (
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func TestContainsComment(t *testing.T) {
	if !containsComment(`int a; /* A takes care of counts */`, [][]string{{"/*", "*/"}}) {
		t.Errorf("invalid")
	}
	if !containsComment(`bool f; /* `, [][]string{{"/*", "*/"}}) {
		t.Errorf("invalid")
	}
	if containsComment(`}`, [][]string{{"/*", "*/"}}) {
		t.Errorf("invalid")
	}
}

func TestCheckMD5SumIgnore(t *testing.T) {
	fileCache := make(map[string]struct{})

	if checkMD5Sum("./utils_test.go", fileCache) {
		t.Errorf("invalid sequence")
	}
	if !checkMD5Sum("./utils_test.go", fileCache) {
		t.Errorf("invalid sequence")
	}
}

func TestCheckDefaultIgnore(t *testing.T) {
	appFS := afero.NewMemMapFs()
	if err := appFS.Mkdir("/test", os.ModeDir); err != nil {
		t.Fatal(err)
	}
	_, _ = appFS.Create("/test/one.go")

	fileInfo, _ := appFS.Stat("/")
	if !checkDefaultIgnore("/", fileInfo, false) {
		t.Errorf("invalid logic: this is directory")
	}

	if !checkDefaultIgnore("/", fileInfo, true) {
		t.Errorf("invalid logic: this is vcs file or directory")
	}

	fileInfo, _ = appFS.Stat("/test/one.go")
	if checkDefaultIgnore("/test/one.go", fileInfo, false) {
		t.Errorf("invalid logic: should not ignore this file")
	}
}

type MockFileInfo struct {
	FileName    string
	IsDirectory bool
}

func (mfi MockFileInfo) Name() string       { return mfi.FileName }
func (mfi MockFileInfo) Size() int64        { return int64(8) }
func (mfi MockFileInfo) Mode() os.FileMode  { return os.ModePerm }
func (mfi MockFileInfo) ModTime() time.Time { return time.Now() }
func (mfi MockFileInfo) IsDir() bool        { return mfi.IsDirectory }
func (mfi MockFileInfo) Sys() interface{}   { return nil }

func TestCheckOptionMatch(t *testing.T) {
	opts := &ClocOptions{}
	fi := MockFileInfo{FileName: "/", IsDirectory: true}
	if !checkOptionMatch("/", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is nil")
	}

	opts.ReNotMatchDir = regexp.MustCompile("thisisdir-not-match")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if !checkOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is nil")
	}

	opts.ReNotMatchDir = regexp.MustCompile("thisisdir")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if checkOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is ignore")
	}

	opts = &ClocOptions{}
	opts.ReMatchDir = regexp.MustCompile("thisisdir")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if !checkOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is not ignore")
	}

	opts.ReMatchDir = regexp.MustCompile("thisisdir-not-match")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if checkOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is ignore")
	}

	opts = &ClocOptions{}
	opts.ReNotMatchDir = regexp.MustCompile("thisisdir-not-match")
	opts.ReMatchDir = regexp.MustCompile("thisisdir")
	fi = MockFileInfo{FileName: "one.go", IsDirectory: false}
	if !checkOptionMatch("/thisisdir/one.go", fi, opts) {
		t.Errorf("invalid logic: renotmatchdir is not ignore")
	}

	t.Run("--match option", func(t *testing.T) {
		opts = &ClocOptions{
			ReMatch: regexp.MustCompile("app.py"),
		}
		fi = MockFileInfo{FileName: "app.py", IsDirectory: false}
		if !checkOptionMatch("test_dir/app.py", fi, opts) {
			t.Errorf("invalid logic: match is not ignore")
		}
	})

	t.Run("--match option with --fullpath option", func(t *testing.T) {
		opts = &ClocOptions{
			ReMatch:  regexp.MustCompile("test_dir/app.py"),
			Fullpath: true,
		}
		fi = MockFileInfo{FileName: "app.py", IsDirectory: false}
		if !checkOptionMatch("test_dir/app.py", fi, opts) {
			t.Errorf("invalid logic: match(with fullpath) is not ignore")
		}
		if checkOptionMatch("app.py", fi, opts) {
			t.Errorf("invalid logic: match(with fullpath) is ignore")
		}
	})

	t.Run("--not-match option with --fullpath option", func(t *testing.T) {
		opts = &ClocOptions{
			ReNotMatch: regexp.MustCompile("test_dir/app.py"),
			Fullpath:   true,
		}
		fi = MockFileInfo{FileName: "app.py", IsDirectory: false}
		if checkOptionMatch("test_dir/app.py", fi, opts) {
			t.Errorf("invalid logic: not-match(with fullpath) is ignore")
		}
		if !checkOptionMatch("app.py", fi, opts) {
			t.Errorf("invalid logic: not-match(with fullpath) is not ignore")
		}
	})
}
