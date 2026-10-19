package exercises

import (
	"bytes"
	"io/fs"
	"net/http"
	"os"
	"path"
	"slices"

	"tests-runner/git"
	"tests-runner/util/errors"
)

const (
	cSolution    = "solution.c"
	rustTemplate = "lib.rs"
)

var (
	errSolutionDirMissing = errors.Error{
		HTTPStatus: http.StatusBadRequest,
		Code:       "solution.missing",
		Message:    "Solution directory is missing",
	}

	errMultipleSolutions = errors.Error{
		HTTPStatus: http.StatusConflict,
		Code:       "solution.multiple_languages",
		Message:    "Solution has versions in multiple languages: both solution.c and lib.rs differ from the templates, restore the one you don't want graded",
	}
)

type Solution struct {
	Dir      string
	Makefile string
	AuxSrc   fs.FS
}

func (e *Exercise) FindSolution(repoDir string) (*Solution, error) {
	dir := path.Join(repoDir, e.Name)
	if st, err := os.Stat(dir); err != nil || !st.IsDir() {
		return nil, errSolutionDirMissing.WithMessagef("directory not found in the solution repo: %q", e.Name)
	}

	if !e.rustIsModified(dir) {
		return &Solution{
			Dir:      dir,
			Makefile: "Makefile",
			AuxSrc:   withoutFS{e.AuxSrc, rustTemplate},
		}, nil
	}

	cIsModified, err := e.cIsModified(repoDir)
	if err != nil {
		return nil, err
	}
	if cIsModified {
		return nil, errMultipleSolutions.WithStack()
	}

	return &Solution{
		Dir:      dir,
		Makefile: "Makefile.rust",
		AuxSrc:   withoutFS{e.AuxSrc, rustTemplate},
	}, nil
}

func (e *Exercise) rustIsModified(dir string) bool {
	template, err := fs.ReadFile(e.AuxSrc, rustTemplate)
	if err != nil {
		return false
	}
	solution, err := os.ReadFile(path.Join(dir, rustTemplate))
	return err == nil && !bytes.Equal(solution, template)
}

func (e *Exercise) cIsModified(repoDir string) (bool, error) {
	file := path.Join(e.Name, cSolution)
	template, found, err := git.FirstVersion(repoDir, file)
	if err != nil {
		return false, err
	}
	solution, err := os.ReadFile(path.Join(repoDir, file))
	if errors.Is(err, fs.ErrNotExist) {
		return found, nil
	} else if err != nil {
		return false, err
	}
	return !found || !bytes.Equal(solution, template), nil
}

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
