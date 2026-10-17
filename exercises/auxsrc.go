package exercises

import (
	"embed"
	"io/fs"
	"path"

	"tests-runner/util/errors"
	"tests-runner/util/util/errutil"
)

//go:embed c go
var sources embed.FS

var AuxSrcStdlib fs.FS

func init() {
	AuxSrcStdlib = mustSub("c/stdlib")
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
