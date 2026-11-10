package cli

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

type reviewCommitPresence int

const (
	reviewCommitNotProvenGone reviewCommitPresence = iota
	reviewCommitGone
	reviewCommitHistoryIncomplete
)

// reviewRecordedCommitPresence proves absence only from an unambiguous answer
// in a complete repository. Probing a partial clone could trigger a lazy fetch.
func reviewRecordedCommitPresence(ctx context.Context, root, commit string) reviewCommitPresence {
	if len(commit) != 40 && len(commit) != 64 {
		return reviewCommitNotProvenGone
	}
	for _, digit := range commit {
		if !(digit >= '0' && digit <= '9' || digit >= 'a' && digit <= 'f') {
			return reviewCommitNotProvenGone
		}
	}
	run := func(args ...string) (string, string, error) {
		cmd := exec.CommandContext(ctx, "git", append([]string{"-C", root, "-c", "core.fsmonitor=false"}, args...)...)
		var stdout, stderr bytes.Buffer
		cmd.Stdout, cmd.Stderr = &stdout, &stderr
		err := cmd.Run()
		return stdout.String(), stderr.String(), err
	}
	shallow, _, err := run("rev-parse", "--is-shallow-repository")
	shallow = strings.TrimSpace(shallow)
	if err != nil || shallow != "true" && shallow != "false" {
		return reviewCommitNotProvenGone
	}
	config, _, err := run("config", "--get-regexp", `^(remote\..+\.promisor|extensions\.partialclone)$`)
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
			return reviewCommitNotProvenGone
		}
	} else {
		for _, line := range strings.Split(config, "\n") {
			fields := strings.Fields(line)
			if len(fields) > 0 && (fields[0] == "extensions.partialclone" || len(fields) == 2 && fields[1] == "true") {
				return reviewCommitHistoryIncomplete
			}
		}
	}
	stdout, stderr, err := run("cat-file", "-e", commit)
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 && stdout == "" && stderr == "" {
		if shallow == "true" {
			return reviewCommitHistoryIncomplete
		}
		return reviewCommitGone
	}
	return reviewCommitNotProvenGone
}

func reviewHistoryIncompleteError(commit string) error {
	return fmt.Errorf("a shallow or partial clone cannot prove that the recorded commit %s is gone; fetch the full history and try again", commit)
}

func reviewedHeadGoneRefusal(ctx context.Context, root, reviewed string) error {
	switch reviewRecordedCommitPresence(ctx, root, reviewed) {
	case reviewCommitGone:
		return fmt.Errorf("reviewed head %q is no longer in the repository; run roundfix review to open a new Reviewer Lineage", reviewed)
	case reviewCommitHistoryIncomplete:
		return reviewHistoryIncompleteError(reviewed)
	default:
		return nil
	}
}
