package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/moby/moby/api/pkg/stdcopy"
	containertypes "github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/mount"
	"github.com/moby/moby/client"
)

type runtimeConfig struct {
	Image    string
	FileName string
	Command  []string
}

var runtimes = map[string]runtimeConfig{
	"python": {
		Image:    "python:3.13-alpine",
		FileName: "main.py",
		Command:  []string{"python3", "/code/main.py"},
	},
	"gcc": {
		Image:    "gcc:15",
		FileName: "main.cpp",
		Command:  []string{"sh", "-c", "g++ /code/main.cpp -std=c++20 -O2 -o /tmp/main && /tmp/main"},
	},
	"clang": {
		Image:    "silkeh/clang:20",
		FileName: "main.cpp",
		Command:  []string{"sh", "-c", "clang++ /code/main.cpp -std=c++20 -O2 -o /tmp/main && /tmp/main"},
	},
}

func (r *Runner) Run(ctx context.Context, request RunRequest) (RunResult, error) {
	config, ok := runtimes[request.Runtime]
	if !ok {
		return RunResult{}, fmt.Errorf("unsupported runtime %q", request.Runtime)
	}
	if err := r.ensureImage(ctx, config.Image); err != nil {
		return RunResult{}, err
	}

	directory, err := os.MkdirTemp("", "runner-*")
	if err != nil {
		return RunResult{}, err
	}
	defer os.RemoveAll(directory)
	if err := os.WriteFile(filepath.Join(directory, config.FileName), []byte(request.Code), 0o644); err != nil {
		return RunResult{}, err
	}

	runCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	response, err := r.docker.ContainerCreate(runCtx, client.ContainerCreateOptions{
		Config: &containertypes.Config{
			Image:        config.Image,
			Cmd:          config.Command,
			AttachStdout: true,
			AttachStderr: true,
		},
		HostConfig: &containertypes.HostConfig{
			NetworkMode: "none",
			Resources: containertypes.Resources{
				Memory:   256 * 1024 * 1024,
				NanoCPUs: 1_000_000_000,
			},
			Mounts: []mount.Mount{{Type: mount.TypeBind, Source: directory, Target: "/code", ReadOnly: true}},
		},
	})
	if err != nil {
		return RunResult{}, err
	}
	defer r.removeContainer(response.ID)

	attached, err := r.docker.ContainerAttach(runCtx, response.ID, client.ContainerAttachOptions{
		Stream: true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		return RunResult{}, err
	}
	defer attached.Close()

	var stdout, stderr bytes.Buffer
	outputDone := make(chan error, 1)
	go func() {
		_, err := stdcopy.StdCopy(&stdout, &stderr, attached.Reader)
		outputDone <- err
	}()
	if _, err := r.docker.ContainerStart(runCtx, response.ID, client.ContainerStartOptions{}); err != nil {
		return RunResult{}, err
	}
	wait := r.docker.ContainerWait(runCtx, response.ID, client.ContainerWaitOptions{Condition: containertypes.WaitConditionNotRunning})
	var exitCode int
	select {
	case err := <-wait.Error:
		if err != nil {
			return RunResult{}, err
		}
	case status := <-wait.Result:
		exitCode = int(status.StatusCode)
	case <-runCtx.Done():
		_, _ = r.docker.ContainerKill(context.Background(), response.ID, client.ContainerKillOptions{Signal: "SIGKILL"})
		return RunResult{}, errors.New("execution timed out")
	}
	if err := <-outputDone; err != nil {
		return RunResult{}, err
	}
	return RunResult{Stdout: stdout.String(), Stderr: stderr.String(), ExitCode: exitCode}, nil
}

func (r *Runner) ensureImage(ctx context.Context, image string) error {
	if _, err := r.docker.ImageInspect(ctx, image); err == nil {
		return nil
	}
	pullCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()
	response, err := r.docker.ImagePull(pullCtx, image, client.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer response.Close()
	return response.Wait(pullCtx)
}

func (r *Runner) removeContainer(id string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, _ = r.docker.ContainerRemove(ctx, id, client.ContainerRemoveOptions{Force: true})
}
