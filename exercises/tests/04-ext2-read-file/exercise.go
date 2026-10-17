package exercise

import (
	_ "embed"

	"tests-runner/exercises"
)

var (
	//go:embed "testdata/506-hello"
	t00dump  []byte
	t00dumpF = exercises.ExpectedOutputFile{
		Name:      "dump",
		IfMissing: "solution never created an output file",
		Verifier:  exercises.BytesEqual(t00dump),
	}
	t00 = exercises.Testcase{
		Name:            "00-helloworld",
		Args:            []string{"test", "506"},
		ExpectedOutputs: []exercises.ExpectedOutputFile{t00dumpF},
	}

	//go:embed "testdata/002-root"
	t01dump  []byte
	t01dumpF = exercises.ExpectedOutputFile{
		Name:      "dump",
		IfMissing: "solution never created an output file",
		Verifier:  exercises.BytesEqual(t01dump),
	}
	t01 = exercises.Testcase{
		Name:            "01-single-indirect",
		Args:            []string{"test", "002"},
		ExpectedOutputs: []exercises.ExpectedOutputFile{t01dumpF},
	}

	//go:embed "testdata/016-2ind"
	t02dump  []byte
	t02dumpF = exercises.ExpectedOutputFile{
		Name:      "dump",
		IfMissing: "solution never created an output file",
		Verifier:  exercises.BytesEqual(t02dump),
	}
	t02 = exercises.Testcase{
		Name:            "02-double-indirect",
		Args:            []string{"test", "016"},
		ExpectedOutputs: []exercises.ExpectedOutputFile{t02dumpF},
	}
)

func init() {
	exercises.RegisterExercise(exercises.Exercise{
		Name:      "04-ext2-read-file",
		Testcases: []exercises.Testcase{t00, t01, t02},
	})
}
