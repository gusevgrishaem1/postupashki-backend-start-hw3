package runner

type RunRequest struct {
	Runtime string
	Code    string
}

type RunResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}
