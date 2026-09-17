package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"

	"roundfix/internal/daemon"
)

var trackedExecutables = []string{
	".agents/skills/systematic-debugging/find-polluter.sh",
	".githooks/commit-msg",
	".githooks/pre-commit",
}

func main() {
	if len(os.Args) != 2 {
		panic("usage: tracked_executable_replay <repository-root>")
	}
	sourceRoot := os.Args[1]
	fixture, err := os.MkdirTemp("", "roundfix-0141-tracked-executable-")
	must(err)
	defer os.RemoveAll(fixture)

	run(fixture, "git", "init", "--quiet")
	run(fixture, "git", "config", "user.name", "Roundfix QA")
	run(fixture, "git", "config", "user.email", "qa@example.invalid")
	run(fixture, "git", "config", "commit.gpgsign", "false")
	for _, path := range trackedExecutables {
		source := filepath.Join(sourceRoot, path)
		info, err := os.Stat(source)
		must(err)
		content, err := os.ReadFile(source)
		must(err)
		destination := filepath.Join(fixture, path)
		must(os.MkdirAll(filepath.Dir(destination), 0o755))
		must(os.WriteFile(destination, content, info.Mode().Perm()))
	}
	run(fixture, "git", "add", "--", ".")
	run(fixture, "git", "commit", "--quiet", "-m", "seed tracked executables")

	for _, path := range trackedExecutables {
		file, err := os.OpenFile(filepath.Join(fixture, path), os.O_APPEND|os.O_WRONLY, 0)
		must(err)
		_, err = file.WriteString("\n# Spec 0141 QA replay\n")
		must(err)
		must(file.Close())
	}

	kept, dropped := daemon.FilterStageablePaths(context.Background(), fixture, trackedExecutables)
	if len(dropped) != 0 || !slices.Equal(kept, trackedExecutables) {
		panic(fmt.Sprintf("kept=%v dropped=%+v", kept, dropped))
	}
	args := append([]string{"add", "--"}, kept...)
	run(fixture, "git", args...)
	run(fixture, "git", "commit", "--quiet", "-m", "carry tracked executable edits")
	committed := output(fixture, "git", "diff-tree", "--no-commit-id", "--name-only", "-r", "HEAD")
	fmt.Printf("kept: %v\n", kept)
	fmt.Printf("dropped: %v\n", dropped)
	fmt.Printf("settlement commit paths:\n%s", committed)
}

func run(dir string, name string, args ...string) {
	command := exec.Command(name, args...)
	command.Dir = dir
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	must(command.Run())
}

func output(dir string, name string, args ...string) string {
	command := exec.Command(name, args...)
	command.Dir = dir
	value, err := command.Output()
	must(err)
	return string(value)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
