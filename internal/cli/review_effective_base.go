package cli

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"roundfix/internal/preflight"
)

// reviewEffectiveBase returns the candidate commit's merge base with the review base.
func reviewEffectiveBase(ctx context.Context, git preflight.GitRunner, root, base, commit string) (string, error) {
	effectiveBase, err := git.RunGit(ctx, root, "merge-base", base, commit)
	if err != nil {
		return "", fmt.Errorf("find review effective base of %s and %s: %w", base, commit, err)
	}
	return strings.TrimSpace(effectiveBase), nil
}

// reviewMeasurement names the origin of a candidate-only diff.
type reviewMeasurement struct {
	From       string
	Carried    string
	Conflicted bool
	Conflicts  []string
}

// reviewMeasuredFrom carries from's own change onto to's Effective Base.
func reviewMeasuredFrom(ctx context.Context, git preflight.GitRunner, root, base, from, to string) (reviewMeasurement, error) {
	fromBase, err := reviewEffectiveBase(ctx, git, root, base, from)
	if err != nil {
		return reviewMeasurement{}, fmt.Errorf("measure review diff from %s to %s: %w", from, to, err)
	}
	toBase, err := reviewEffectiveBase(ctx, git, root, base, to)
	if err != nil {
		return reviewMeasurement{}, fmt.Errorf("measure review diff from %s to %s: %w", from, to, err)
	}
	if fromBase == toBase {
		return reviewMeasurement{From: from}, nil
	}
	if fromBase == from {
		return reviewMeasurement{From: toBase}, nil
	}
	tree, err := git.RunGit(ctx, root, "merge-tree", "--write-tree", from, toBase)
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
			return reviewConflictedCarry(ctx, root, from, to, toBase)
		}
		return reviewMeasurement{}, fmt.Errorf("carry review diff from %s to %s onto %s: %w", from, to, toBase, err)
	}
	firstLine, _, _ := strings.Cut(tree, "\n")
	firstLine = strings.TrimSpace(firstLine)
	if !reviewCarryTreeIDValid(firstLine) {
		return reviewMeasurement{}, fmt.Errorf("carry review diff from %s to %s onto %s: invalid tree %q", from, to, toBase, firstLine)
	}
	return reviewMeasurement{From: firstLine, Carried: from}, nil
}

// The runner folds stdout into its error on conflict; read Git's NUL-delimited
// tree and paths directly, with a fixed marker style independent of user config.
func reviewConflictedCarry(ctx context.Context, root, from, to, toBase string) (reviewMeasurement, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", root, "-c", "core.fsmonitor=false", "-c", "merge.conflictStyle=merge", "merge-tree", "--write-tree", "--name-only", "--no-messages", "-z", from, toBase)
	output, err := cmd.Output()
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) || exitErr.ExitCode() != 1 {
		if err == nil {
			err = errors.New("expected merge-tree to exit 1 for a conflicted carry")
		}
		return reviewMeasurement{}, fmt.Errorf("carry review diff from %s to %s onto %s: %w", from, to, toBase, err)
	}
	tree, paths, delimited := strings.Cut(string(output), "\x00")
	if !delimited || !reviewCarryTreeIDValid(tree) {
		return reviewMeasurement{}, fmt.Errorf("carry review diff from %s to %s onto %s: invalid conflicted tree %q", from, to, toBase, tree)
	}
	var conflicts []string
	for _, path := range strings.Split(paths, "\x00") {
		if path != "" {
			conflicts = append(conflicts, path)
		}
	}
	slices.Sort(conflicts)
	return reviewMeasurement{From: tree, Carried: from, Conflicted: true, Conflicts: slices.Compact(conflicts)}, nil
}

// The note sits outside the diff's untrusted-data boundary, so a candidate
// path with a control character or invalid UTF-8 is quoted: a line break in a
// file name can never start a new prompt line. Other paths stay as Git names them.
func reviewConflictedPathText(path string) string {
	if !utf8.ValidString(path) || strings.ContainsFunc(path, func(r rune) bool { return unicode.IsControl(r) || unicode.In(r, unicode.Zl, unicode.Zp) }) {
		return strconv.Quote(path)
	}
	return path
}

func reviewCarryTreeIDValid(tree string) bool {
	return (len(tree) == 40 || len(tree) == 64) && strings.Trim(tree, "0123456789abcdef") == ""
}

func reviewConflictedCarryNote(measurement reviewMeasurement) string {
	if !measurement.Conflicted {
		return ""
	}
	note := fmt.Sprintf("Carrying %s onto the default branch conflicted; where tree %s holds Git's conflict markers, the diff shows how the candidate resolved the conflict.\n", measurement.Carried, measurement.From)
	if len(measurement.Conflicts) > 0 {
		paths := make([]string, len(measurement.Conflicts))
		for index, path := range measurement.Conflicts {
			paths[index] = reviewConflictedPathText(path)
		}
		note += "Conflicted paths: " + strings.Join(paths, ", ") + ".\n"
	}
	return note
}
