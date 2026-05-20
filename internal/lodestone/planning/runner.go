package planning

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
)

type Runner interface {
	Run(ctx context.Context, model, prompt string) (string, error)
}

type ClaudeRunner struct {
	Binary string
}

func NewClaudeRunner() *ClaudeRunner {
	return &ClaudeRunner{Binary: "claude"}
}

var ErrClaudeNotFound = errors.New("claude CLI not found in PATH; install from https://docs.anthropic.com/en/docs/claude-code")

func (r *ClaudeRunner) Run(ctx context.Context, model, prompt string) (string, error) {
	bin := r.Binary
	if bin == "" {
		bin = "claude"
	}
	if _, err := exec.LookPath(bin); err != nil {
		return "", ErrClaudeNotFound
	}

	args := []string{"--print"}
	if model != "" {
		args = append(args, "--model", model)
	}
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdin = bytes.NewBufferString(prompt)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("claude run: %w; stderr: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

type FakeRunner struct {
	Output string
	Err    error
	Calls  []FakeCall
}

type FakeCall struct {
	Model  string
	Prompt string
}

func (r *FakeRunner) Run(_ context.Context, model, prompt string) (string, error) {
	r.Calls = append(r.Calls, FakeCall{Model: model, Prompt: prompt})
	if r.Err != nil {
		return "", r.Err
	}
	return r.Output, nil
}
