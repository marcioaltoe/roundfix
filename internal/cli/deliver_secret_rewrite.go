package cli

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"roundfix/internal/delivery"
	"roundfix/internal/secretscan"
)

var _ delivery.SecretRewriteProver = (*commandDeliveryWorkflow)(nil)

func (workflow *commandDeliveryWorkflow) ProveSecretRewrite(ctx context.Context, workDir, base, candidate, head string, named []string) (delivery.SecretRewrite, error) {
	refuse := func(reason string) (delivery.SecretRewrite, error) {
		return delivery.SecretRewrite{Reason: reason}, nil
	}
	if _, err := workflow.git.RunGit(ctx, workDir, "merge-base", "--is-ancestor", base, head); err != nil {
		var exitErr interface{ ExitCode() int }
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return refuse("item head does not descend from the Run start " + base)
		}
		return delivery.SecretRewrite{}, fmt.Errorf("check secret rewrite ancestry: %w", err)
	}
	raw, err := workflow.git.RunGit(ctx, workDir, "diff", "--raw", "--no-renames", "--no-ext-diff", "--no-abbrev", "-z", candidate, head, "--")
	if err != nil {
		return delivery.SecretRewrite{}, fmt.Errorf("list secret rewrite paths: %w", err)
	}
	allowed := make(map[string]bool, len(named))
	for _, path := range named {
		allowed[path] = true
	}
	type change struct {
		path          string
		fields        []string
		before, after []string
		matched       map[int]bool
	}
	var changes []change
	entries := strings.Split(strings.TrimSuffix(raw, "\x00"), "\x00")
	if raw == "" {
		entries = nil
	}
	for i := 0; i < len(entries); i += 2 {
		if i+1 >= len(entries) {
			return delivery.SecretRewrite{}, errors.New("read secret rewrite diff: incomplete path entry")
		}
		fields := strings.Fields(strings.TrimPrefix(entries[i], ":"))
		if len(fields) != 5 {
			return delivery.SecretRewrite{}, errors.New("read secret rewrite diff: invalid metadata")
		}
		path := entries[i+1]
		if !allowed[path] {
			return refuse(fmt.Sprintf("the rewrite changes %s, which the blocker does not name", path))
		}
		changes = append(changes, change{path: path, fields: fields})
	}
	// Apply each rule across all paths before moving to the next rule.
	for _, c := range changes {
		if c.fields[4] != "M" || c.fields[0] != c.fields[1] {
			return refuse(fmt.Sprintf("the rewrite adds, removes or changes the mode of %s", c.path))
		}
	}
	for i := range changes {
		c := &changes[i]
		before, err := workflow.secretRewriteBlob(ctx, workDir, c.fields[2])
		if err != nil {
			return delivery.SecretRewrite{}, err
		}
		after, err := workflow.secretRewriteBlob(ctx, workDir, c.fields[3])
		if err != nil {
			return delivery.SecretRewrite{}, err
		}
		c.before, c.after = strings.Split(before, "\n"), strings.Split(after, "\n")
		if len(c.before) != len(c.after) {
			return refuse(fmt.Sprintf("the rewrite changes the line count of %s", c.path))
		}
		c.matched = make(map[int]bool)
		for _, m := range secretscan.ScanContent(c.path, []byte(before)) {
			c.matched[m.Line] = true
		}
	}
	for _, c := range changes {
		for i, line := range c.before {
			if line != c.after[i] && !c.matched[i+1] {
				return refuse(fmt.Sprintf("the rewrite changes line %d of %s, which holds no Secret Match at the parked candidate", i+1, c.path))
			}
		}
	}
	// Unchanged named files must be clean too: a rewrite cannot leave a match
	// behind merely because that path has no entry in the tree diff.
	for _, path := range named {
		content, err := workflow.secretRewriteBlob(ctx, workDir, head+":"+path)
		if err != nil {
			return delivery.SecretRewrite{}, err
		}
		if matches := secretscan.ScanContent(path, []byte(content)); len(matches) > 0 {
			return refuse(fmt.Sprintf("the rewrite keeps a Secret Match at line %d of %s", matches[0].Line, path))
		}
	}
	return delivery.SecretRewrite{Accepted: true}, nil
}

func (workflow *commandDeliveryWorkflow) secretRewriteBlob(ctx context.Context, workDir, object string) (string, error) {
	body, err := workflow.git.RunGit(ctx, workDir, "cat-file", "blob", object)
	if err != nil {
		return "", fmt.Errorf("read secret rewrite blob: %w", err)
	}
	// ExecGitRunner trims trailing LF bytes. The immutable object's size restores
	// them, preserving the line count and every byte used as proof.
	sizeText, err := workflow.git.RunGit(ctx, workDir, "cat-file", "-s", object)
	if err != nil {
		return "", fmt.Errorf("read secret rewrite blob size: %w", err)
	}
	size, err := strconv.Atoi(strings.TrimSpace(sizeText))
	if err != nil {
		return "", fmt.Errorf("parse secret rewrite blob size: %w", err)
	}
	if size < len(body) {
		return "", errors.New("read secret rewrite blob: content exceeds object size")
	}
	return body + strings.Repeat("\n", size-len(body)), nil
}
