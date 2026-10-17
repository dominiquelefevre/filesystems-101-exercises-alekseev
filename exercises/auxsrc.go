package exercises

import (
	"embed"
	"io/fs"
	"os"
	"path"
	"slices"

	"tests-runner/util/errors"
	"tests-runner/util/util/errutil"
)

//go:embed c go rust
var sources embed.FS

var (
	cStdlib    fs.FS
	rustStdlib fs.FS
)

func init() {
	cStdlib = mustSub("c/stdlib")
	rustStdlib = mustSub("rust/stdlib")
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

const RustMarker = ".selected"

type withoutFS struct {
	fs.FS
	name string
}

func (w withoutFS) ReadDir(name string) ([]fs.DirEntry, error) {
	entries, err := fs.ReadDir(w.FS, name)
	if name == "." {
		entries = slices.DeleteFunc(entries, func(e fs.DirEntry) bool { return e.Name() == w.name })
	}
	return entries, err
}

func (e *Exercise) FindSolution(repoDir string) *Solution {
	exists := func(p string) bool {
		_, err := os.Stat(path.Join(repoDir, p))
		return err == nil
	}

	var lang, dir string
	switch {
	case e.RustAuxSrc != nil && exists(path.Join("rust", e.Name, RustMarker)):
		lang, dir = "rust", path.Join("rust", e.Name)
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
	if lang == "rust" {
		s.Copies = append(s.Copies, Copy{e.RustAuxSrc, s.Dir}, Copy{rustStdlib, stdlib})
	}
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

	if isDir(sources, path.Join("rust", e.Name)) {
		if e.Lang != "c" {
			errutil.FatalIf(errors.Errorf("exercise %s: Rust solutions need a C harness", e.Name))
		}
		e.RustAuxSrc = withoutFS{mustSub(path.Join("rust", e.Name)), "lib.rs"}
	}
}
