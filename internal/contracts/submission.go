package contracts

type Submission struct {
	ID       string `json:"task_id"`
	Code     string `json:"code"`
	Language string `json:"language"`
	Input    string `json:"input,omitempty"`
}

type Result struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}
