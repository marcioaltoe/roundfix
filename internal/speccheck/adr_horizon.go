package speccheck

import (
	"context"
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
)

type adrHorizon struct {
	repoRoot  string
	prdCommit string
	addedBy   map[string]string
}

// newADRHorizon returns false when the full related-ADR check must apply.
func newADRHorizon(repoRoot, prdPath string) (adrHorizon, bool) {
	repoRoot = filepath.Clean(repoRoot)
	prdCommit, readable := prdAddingCommit(repoRoot, prdPath)
	if !readable || prdCommit == "" {
		return adrHorizon{}, false
	}

	adrLog, err := adrHorizonGitOutput(
		repoRoot,
		"log", "--diff-filter=A", "--format=%x00%H", "--name-only", "--", "docs/adr",
	)
	if err != nil {
		return adrHorizon{}, false
	}
	return adrHorizon{
		repoRoot:  repoRoot,
		prdCommit: prdCommit,
		addedBy:   adrAddingCommits(adrLog),
	}, true
}

// prdAddingCommit distinguishes unreadable history from an uncommitted PRD.
func prdAddingCommit(repoRoot, prdPath string) (string, bool) {
	repoRoot = filepath.Clean(repoRoot)
	ctx := context.Background()
	gitRoot, err := mechanicalRepositoryRoot(ctx, repoRoot)
	if err != nil || !adrHorizonSamePath(repoRoot, gitRoot) {
		return "", false
	}
	repoCommonDir, err := mechanicalRepositoryCommonDir(ctx, repoRoot)
	if err != nil {
		return "", false
	}
	prdCommonDir, err := mechanicalRepositoryCommonDir(ctx, filepath.Dir(prdPath))
	if err != nil || repoCommonDir != prdCommonDir {
		return "", false
	}

	shallow, err := adrHorizonGitOutput(repoRoot, "rev-parse", "--is-shallow-repository")
	if err != nil || strings.TrimSpace(string(shallow)) != "false" {
		return "", false
	}

	prdRelativePath, err := filepath.Rel(repoRoot, prdPath)
	if err != nil {
		return "", false
	}
	prdRelativePath = filepath.ToSlash(filepath.Clean(prdRelativePath))
	if prdRelativePath == ".." || strings.HasPrefix(prdRelativePath, "../") {
		return "", false
	}
	prdCommitOutput, err := adrHorizonGitOutput(
		repoRoot,
		"log", "-1", "--diff-filter=A", "--format=%H", "--", prdRelativePath,
	)
	if err != nil {
		return "", false
	}
	prdCommit := strings.TrimSpace(string(prdCommitOutput))
	if strings.ContainsAny(prdCommit, "\r\n") {
		return "", false
	}

	return prdCommit, true
}

func adrHorizonSamePath(first, second string) bool {
	first, firstErr := filepath.Abs(first)
	second, secondErr := filepath.Abs(second)
	if firstErr != nil || secondErr != nil {
		return false
	}
	if resolved, err := filepath.EvalSymlinks(first); err == nil {
		first = resolved
	}
	if resolved, err := filepath.EvalSymlinks(second); err == nil {
		second = resolved
	}
	return filepath.Clean(first) == filepath.Clean(second)
}

func (horizon adrHorizon) predates(adrPath string) bool {
	adrPath = filepath.ToSlash(filepath.Clean(filepath.FromSlash(adrPath)))
	adrCommit := horizon.addedBy[adrPath]
	if adrCommit == "" {
		return false
	}

	_, err := adrHorizonGitOutput(
		horizon.repoRoot,
		"merge-base", "--is-ancestor", adrCommit, horizon.prdCommit,
	)
	if err == nil {
		return true
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false
	}
	return true
}

func adrAddingCommits(output []byte) map[string]string {
	addedBy := make(map[string]string)
	for _, block := range strings.Split(string(output), "\x00") {
		lines := strings.Split(strings.TrimSpace(block), "\n")
		if len(lines) < 2 {
			continue
		}
		commit := strings.TrimSpace(lines[0])
		if commit == "" {
			continue
		}
		for _, line := range lines[1:] {
			path := filepath.ToSlash(filepath.Clean(filepath.FromSlash(strings.TrimSpace(line))))
			if path == "." {
				continue
			}
			if _, exists := addedBy[path]; !exists {
				addedBy[path] = commit
			}
		}
	}
	return addedBy
}

func adrHorizonGitOutput(repoRoot string, args ...string) ([]byte, error) {
	commandArgs := make([]string, 0, len(args)+4)
	commandArgs = append(commandArgs, "-C", repoRoot, "-c", "core.fsmonitor=false")
	commandArgs = append(commandArgs, args...)
	command := exec.Command("git", commandArgs...)
	command.Env = mechanicalGitEnvironment()
	return command.CombinedOutput()
}
