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
		Name:      "12-grpc-prometheus",
		Testcases: []exercises.Testcase{t00},
	})
}
