package exercises

import (
	"embed"
	"io/fs"
	"os"
	"path"

	"tests-runner/util/errors"
	"tests-runner/util/util/errutil"
)

//go:embed c go
var sources embed.FS

var cStdlib fs.FS

func init() {
	cStdlib = mustSub("c/stdlib")
}

func mustSub(dir string) fs.FS {
	sub, err := fs.Sub(sources, dir)
	errutil.FatalIf(err)
	return sub
}

func isDir(fsys fs.FS, name string) bool {
	st, err := fs.Stat(fsys, name)
	return err == nil && st.IsDir()
}

type Solution struct {
	Lang   string
	Dir    string
	Copies []Copy
}

type Copy struct {
	Src fs.FS
	Dst string
}

func (e *Exercise) FindSolution(repoDir string) *Solution {
	exists := func(p string) bool {
		_, err := os.Stat(path.Join(repoDir, p))
		return err == nil
	}

	var lang, dir string
	switch {
	case exists(path.Join(e.Lang, e.Name)):
		lang, dir = e.Lang, path.Join(e.Lang, e.Name)
	case exists(e.Name):
		lang, dir = e.Lang, e.Name
	default:
		return nil
	}

	s := &Solution{Lang: lang, Dir: path.Join(repoDir, dir)}
	stdlib := path.Join(path.Dir(s.Dir), "stdlib")
	s.Copies = []Copy{{e.AuxSrc, s.Dir}, {cStdlib, stdlib}}
	return s
}

func (e *Exercise) findSources() {
	for _, lang := range []string{"c", "go"} {
		if isDir(sources, path.Join(lang, e.Name)) {
			if e.AuxSrc != nil {
				errutil.FatalIf(errors.Errorf("exercise %s has sources in more than one language", e.Name))
			}
			e.Lang = lang
			e.AuxSrc = mustSub(path.Join(lang, e.Name))
		}
	}
	if e.AuxSrc == nil {
		errutil.FatalIf(errors.Errorf("exercise %s has no sources in c/ or go/", e.Name))
	}
}
