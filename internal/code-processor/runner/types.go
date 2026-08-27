package runner

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
