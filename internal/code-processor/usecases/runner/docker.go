package runner

import "github.com/moby/moby/client"

type Runner struct {
	docker *client.Client
}

func New() (*Runner, error) {
	cli, err := client.New(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &Runner{docker: cli}, nil
}

func (r *Runner) Close() error {
	return r.docker.Close()
}
