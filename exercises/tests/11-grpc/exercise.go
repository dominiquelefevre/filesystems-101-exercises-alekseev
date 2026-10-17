package exercise

import (
	"tests-runner/exercises"
)

var (
	t00 = exercises.Testcase{
		Name:            "00",
		Args:            []string{},
		WithoutValgrind: true,
	}
)

func init() {
	exercises.RegisterExercise(exercises.Exercise{
		Name:      "11-grpc",
		Testcases: []exercises.Testcase{t00},
	})
}
