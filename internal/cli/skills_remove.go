package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	roundskills "roundfix/skills"
)

func runSkillsRemove(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	fs := flag.NewFlagSet("skills remove", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	target := fs.String("target", "", "User-level skill target")
	dir := fs.String("dir", "", "Override target skills directory for a single target")
	apply := fs.Bool("apply", false, "Remove the listed copies")
	if err := fs.Parse(args); err != nil {
		if strings.Contains(err.Error(), "flag needs an argument: -target") {
			fmt.Fprintln(stderr, "roundfix: skills remove requires --target <claude|codex|opencode|all>")
		} else {
			fmt.Fprintf(stderr, "roundfix: %v\n", err)
		}
		return exitPreflight
	}
	if len(fs.Args()) > 0 {
		fmt.Fprintf(stderr, "roundfix: unexpected argument %q\n", fs.Args()[0])
		return exitPreflight
	}
	value := strings.TrimSpace(*target)
	switch value {
	case "project":
		fmt.Fprintln(stderr, "roundfix: skills remove refuses --target project: the project's copies are the Repository Skill Set")
		return exitPreflight
	case "claude", "codex", "opencode", "all":
	default:
		fmt.Fprintln(stderr, "roundfix: skills remove requires --target <claude|codex|opencode|all>")
		return exitPreflight
	}
	if value == "all" && strings.TrimSpace(*dir) != "" {
		fmt.Fprintln(stderr, "roundfix: --dir requires a single --target value")
		return exitPreflight
	}
	targets := []string{value}
	if value == "all" {
		targets = []string{"claude", "codex", "opencode"}
	}
	var known []roundskills.UserSkillDirectory
	if environment.dependencies.userSkillDirectories != nil {
		known = environment.dependencies.userSkillDirectories(environment.homeDir)
	}
	// Resolve every target before inspecting or removing any entry.
	directories := make([]string, 0, len(targets))
	for _, target := range targets {
		directory := strings.TrimSpace(*dir)
		if directory == "" {
			for _, candidate := range known {
				if candidate.Target == target {
					directory = candidate.Path
					break
				}
			}
		}
		if strings.TrimSpace(directory) == "" {
			fmt.Fprintf(stderr, "roundfix: skills remove failed: no user-level skill directory is known for %s\n", target)
			return exitRunFailed
		}
		directories = append(directories, directory)
	}
	names := roundskills.Names()
	sort.Strings(names)
	listed := false
	for _, directory := range directories {
		found := false
		for _, name := range names {
			if err := ctx.Err(); err != nil {
				fmt.Fprintf(stderr, "roundfix: skills remove failed: %v\n", err)
				return exitRunFailed
			}
			path := filepath.Join(directory, name)
			named, err := removeUserSkillCopy(path, *apply)
			if err != nil {
				fmt.Fprintf(stderr, "roundfix: skills remove failed: %v\n", err)
				return exitRunFailed
			}
			if named != "" {
				fmt.Fprintln(stdout, named)
				found, listed = true, true
			}
		}
		if !found {
			fmt.Fprintf(stdout, "No Roundfix-owned skill copy under %s\n", directory)
		}
	}
	if listed && !*apply {
		fmt.Fprintln(stdout, "Nothing was removed; run the command again with --apply.")
	}
	return exitOK
}

// Lstat and Readlink inspect only the entry: even a dangling link is unlinked
// without opening its target. Regular files at owned names are left alone.
func removeUserSkillCopy(path string, apply bool) (string, error) {
	info, err := os.Lstat(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("inspect %s: %w", path, err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		target, err := os.Readlink(path)
		if err != nil {
			return "", fmt.Errorf("read link %s: %w", path, err)
		}
		if apply {
			if err := os.Remove(path); err != nil {
				return "", fmt.Errorf("unlink %s: %w", path, err)
			}
			return fmt.Sprintf("Unlinked %s -> %s", path, target), nil
		}
		return fmt.Sprintf("Would unlink %s -> %s", path, target), nil
	}
	if !info.IsDir() {
		return "", nil
	}
	if apply {
		if err := os.RemoveAll(path); err != nil {
			return "", fmt.Errorf("remove %s: %w", path, err)
		}
		return "Removed " + path, nil
	}
	return "Would remove " + path, nil
}
