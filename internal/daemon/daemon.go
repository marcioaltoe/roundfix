package daemon

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

type VerifyRequest struct {
	WorkDir string
	Command string
	Stream  io.Writer
}

type Verifier interface {
	Verify(context.Context, VerifyRequest) error
}

type ExecVerifier struct{}

func (ExecVerifier) Verify(ctx context.Context, req VerifyRequest) error {
	if strings.TrimSpace(req.WorkDir) == "" {
		return fmt.Errorf("run verification: git root is required")
	}
	if strings.TrimSpace(req.Command) == "" {
		return fmt.Errorf("run verification: command is required")
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", req.Command)
	cmd.Dir = req.WorkDir
	if req.Stream != nil {
		cmd.Stdout = req.Stream
		cmd.Stderr = req.Stream
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("verification command %q failed: %w", req.Command, err)
	}
	return nil
}

type CommitRequest struct {
	WorkDir string
	Message string
}

type Committer interface {
	Commit(context.Context, CommitRequest) error
}

type GitCommitter struct{}

func (GitCommitter) Commit(ctx context.Context, req CommitRequest) error {
	if strings.TrimSpace(req.WorkDir) == "" {
		return fmt.Errorf("create Batch commit: git root is required")
	}
	if strings.TrimSpace(req.Message) == "" {
		return fmt.Errorf("create Batch commit: message is required")
	}
	if err := runGit(ctx, req.WorkDir, "add", "-A"); err != nil {
		return err
	}
	if err := runGit(ctx, req.WorkDir, "commit", "-m", req.Message); err != nil {
		return err
	}
	return nil
}

func BatchCommitMessage(batchNumber int) string {
	return fmt.Sprintf("fix: resolve Roundfix batch %03d", batchNumber)
}

type PushRequest struct {
	WorkDir string
	Remote  string
	Branch  string
}

type Pusher interface {
	Push(context.Context, PushRequest) error
}

type GitPusher struct{}

func (GitPusher) Push(ctx context.Context, req PushRequest) error {
	if strings.TrimSpace(req.WorkDir) == "" {
		return fmt.Errorf("run Final Push: git root is required")
	}
	if strings.TrimSpace(req.Remote) == "" {
		return fmt.Errorf("run Final Push: remote is required")
	}
	if strings.TrimSpace(req.Branch) == "" {
		return fmt.Errorf("run Final Push: PR Head Branch is required")
	}
	return runGit(ctx, req.WorkDir, "push", req.Remote, "HEAD:"+req.Branch)
}

func runGit(ctx context.Context, workDir string, args ...string) error {
	cmdArgs := append([]string{"-C", workDir}, args...)
	cmd := exec.CommandContext(ctx, "git", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git %s failed: %w%s", strings.Join(args, " "), err, formatCommandOutput(output))
	}
	return nil
}

func formatCommandOutput(output []byte) string {
	trimmed := strings.TrimSpace(string(output))
	if trimmed == "" {
		return ""
	}
	return ": " + trimmed
}
