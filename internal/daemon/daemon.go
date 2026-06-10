package daemon

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
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
	WorkDir      string
	Message      string
	ExcludePaths []string
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
	addArgs := []string{"add", "-A", "--", "."}
	addArgs = append(addArgs, gitExcludePathspecs(req.WorkDir, req.ExcludePaths)...)
	if err := runGit(ctx, req.WorkDir, addArgs...); err != nil {
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

func gitExcludePathspecs(workDir string, paths []string) []string {
	if len(paths) == 0 {
		return nil
	}
	root := filepath.Clean(workDir)
	pathspecs := make([]string, 0, len(paths))
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		cleanPath := filepath.Clean(path)
		if filepath.IsAbs(cleanPath) {
			rel, err := filepath.Rel(root, cleanPath)
			if err != nil || rel == "." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
				continue
			}
			cleanPath = rel
		}
		cleanPath = filepath.ToSlash(cleanPath)
		if cleanPath == "." || cleanPath == ".." || strings.HasPrefix(cleanPath, "../") {
			continue
		}
		pathspecs = append(pathspecs, ":(exclude)"+cleanPath)
	}
	return pathspecs
}

func runGit(ctx context.Context, workDir string, args ...string) error {
	// fsmonitor is disabled per invocation so daemon warnings such as
	// fsmonitor_ipc__send_query never reach parsed output.
	cmdArgs := append([]string{"-C", workDir, "-c", "core.fsmonitor=false"}, args...)
	cmd := exec.CommandContext(ctx, "git", cmdArgs...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		detail := formatCommandOutput(stderr.Bytes())
		if detail == "" {
			detail = formatCommandOutput(stdout.Bytes())
		}
		return fmt.Errorf("git %s failed: %w%s", strings.Join(args, " "), err, detail)
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
