package apply

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type RealGit struct{}

func runGit(ctx context.Context, repoRoot string, args ...string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = repoRoot
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w (stderr: %s)", strings.Join(args, " "), err, errBuf.String())
	}
	return out.String(), nil
}

func (RealGit) Status(ctx context.Context, repoRoot string) (string, error) {
	return runGit(ctx, repoRoot, "status", "--porcelain")
}
func (RealGit) CreateBranch(ctx context.Context, repoRoot, branch string) error {
	_, err := runGit(ctx, repoRoot, "switch", "-c", branch)
	return err
}
func (RealGit) Add(ctx context.Context, repoRoot string, paths ...string) error {
	args := append([]string{"add"}, paths...)
	_, err := runGit(ctx, repoRoot, args...)
	return err
}
func (RealGit) Commit(ctx context.Context, repoRoot, message string) error {
	_, err := runGit(ctx, repoRoot, "commit", "-m", message)
	return err
}
func (RealGit) Push(ctx context.Context, repoRoot, branch string) error {
	_, err := runGit(ctx, repoRoot, "push", "-u", "origin", branch)
	return err
}
func (RealGit) DeleteBranchLocal(ctx context.Context, repoRoot, branch string) error {
	_, err := runGit(ctx, repoRoot, "branch", "-D", branch)
	return err
}
func (RealGit) DeleteBranchRemote(ctx context.Context, repoRoot, branch string) error {
	_, err := runGit(ctx, repoRoot, "push", "origin", "--delete", branch)
	return err
}

type RealPR struct{}

func (RealPR) CreateDraftPR(ctx context.Context, repoRoot, branch, title, body string) (int, string, error) {
	if _, err := exec.LookPath("gh"); err != nil {
		return 0, "", fmt.Errorf("gh CLI not in PATH; branch %s ist gepusht, PR bitte manuell öffnen", branch)
	}
	cmd := exec.CommandContext(ctx, "gh", "pr", "create", "--draft",
		"--base", "main", "--head", branch, "--title", title, "--body", body)
	cmd.Dir = repoRoot
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return 0, "", fmt.Errorf("gh pr create: %w (stderr: %s)", err, errBuf.String())
	}
	url := strings.TrimSpace(out.String())
	number := parsePRNumber(url)
	return number, url, nil
}

func (RealPR) ClosePR(ctx context.Context, repoRoot string, number int) error {
	if _, err := exec.LookPath("gh"); err != nil {
		return fmt.Errorf("gh CLI not in PATH")
	}
	cmd := exec.CommandContext(ctx, "gh", "pr", "close", strconv.Itoa(number), "--delete-branch")
	cmd.Dir = repoRoot
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("gh pr close: %w (stderr: %s)", err, errBuf.String())
	}
	return nil
}

func parsePRNumber(url string) int {
	idx := strings.LastIndex(url, "/")
	if idx < 0 || idx == len(url)-1 {
		return 0
	}
	n, err := strconv.Atoi(url[idx+1:])
	if err != nil {
		return 0
	}
	return n
}
