package cli

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strings"
)

type specFolderStatus struct {
	Changes []string // Tracked changes and untracked, non-ignored paths.
	Ignored []string // Ignored Spec Files, repository-relative and sorted.
}

func readSpecFolderStatus(ctx context.Context, repo, folder string) (specFolderStatus, error) {
	var status specFolderStatus
	cmd := exec.CommandContext(ctx, "git", "-C", repo, "-c", "core.fsmonitor=false", "status", "--porcelain=v1", "-z", "--no-renames", "--untracked-files=all", "--ignored=traditional", "--", folder)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		detail := strings.TrimSpace(stderr.String())
		if detail == "" {
			detail = err.Error()
		}
		return status, fmt.Errorf("read Spec folder status: %s: %w", detail, err)
	}
	// Porcelain -z keeps spaces and newlines in paths verbatim. --no-renames
	// ensures each entry carries exactly one path, including staged renames.
	for _, entry := range bytes.Split(stdout.Bytes(), []byte{0}) {
		if len(entry) == 0 {
			continue
		}
		if len(entry) < 4 || entry[2] != ' ' {
			return specFolderStatus{}, fmt.Errorf("read Spec folder status: invalid porcelain entry %q", entry)
		}
		path := string(entry[3:])
		if string(entry[:2]) == "!!" {
			status.Ignored = append(status.Ignored, path)
		} else {
			status.Changes = append(status.Changes, path)
		}
	}
	sort.Strings(status.Ignored)
	return status, nil
}

func refuseSpecFolderStatus(command, slug, folderRel string, status specFolderStatus) error {
	if len(status.Changes) > 0 {
		return fmt.Errorf("Spec %q has changes not committed at HEAD under %s; commit them before %s so the Archive Record's source_revision holds the Spec", slug, folderRel, command)
	}
	if len(status.Ignored) > 0 {
		return fmt.Errorf("Spec %q cannot %s while Git ignores %d file(s) under %s that no commit holds: %s; remove each file or commit it with git add -f, then retry", slug, command, len(status.Ignored), folderRel, strings.Join(status.Ignored, ", "))
	}
	return nil
}
