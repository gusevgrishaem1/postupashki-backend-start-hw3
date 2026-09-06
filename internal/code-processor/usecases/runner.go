package usecases

import "context"

type RunRequest struct {
	Runtime string
	Code    string
	Input   string
}

type RunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type Runner interface {
	Run(context.Context, RunRequest) (RunResult, error)
}
