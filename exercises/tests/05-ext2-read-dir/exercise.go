package exercise

import (
	_ "embed"

	"tests-runner/exercises"
)

var (
	//go:embed "testdata/002"
	t00dump  string
	t00dumpF = exercises.ExpectedOutputFile{
		Name:      "dump",
		IfMissing: "solution never created an output file",
		Verifier:  exercises.LinesEqualS(t00dump).Verify,
	}
	t00 = exercises.Testcase{
		Name:            "00",
		Args:            []string{"test", "002"},
		ExpectedOutputs: []exercises.ExpectedOutputFile{t00dumpF},
	}
)

func init() {
	exercises.RegisterExercise(exercises.Exercise{
		Name:      "05-ext2-read-dir",
		Testcases: []exercises.Testcase{t00},
	})
}
