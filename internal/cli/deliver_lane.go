package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/preflight"
)

// runDeliverLane uses the checkout configuration and Git objects, without
// resolving a Spec or opening the Delivery Queue or Run Database.
func runDeliverLane(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	fs := flagSet("deliver lane")
	all := fs.Bool("all", false, "test every package in the Linux Lane")
	if err := fs.Parse(hoistCommandFlags(args, nil)); err != nil {
		return printDeliverFailure("lane", err, stderr)
	}
	if fs.NArg() > 1 {
		return printDeliverFailure("lane", errors.New("deliver lane accepts at most one revision"), stderr)
	}
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	fail := func(err error) int {
		if ctx.Err() != nil {
			return 130
		}
		return printDeliverFailure("lane", err, stderr)
	}
	// Config warnings are not part of the lane's single-line outcome contract.
	loaded, err := loadCommandConfig(environment, io.Discard)
	if err != nil {
		return fail(err)
	}
	if loaded.GitRoot == "" {
		return fail(errors.New("deliver lane requires a git repository working tree"))
	}
	revision := "HEAD"
	if fs.NArg() == 1 {
		revision = fs.Arg(0)
	}
	git := preflight.ExecGitRunner{}
	head, err := git.RunGit(ctx, loaded.GitRoot, "rev-parse", "--verify", "--end-of-options", revision+"^{commit}")
	if err != nil {
		return fail(fmt.Errorf("resolve Linux Lane revision %q: %w", revision, err))
	}
	if !*all && preflight.DetectDefaultBranch(ctx, loaded.GitRoot, "", git).Source == preflight.DefaultBranchUndetermined {
		return fail(errors.New("Linux Lane: repository default branch is unknown"))
	}
	artifactDir, err := roundconfig.ValidateArtifactDirectory(loaded.Config.Defaults.ArtifactDir, loaded.GitRoot, loaded.HomeDir)
	if err != nil {
		return fail(err)
	}
	workflow := &commandDeliveryWorkflow{loaded: loaded, git: git, log: stdout}
	result, err := workflow.runLinuxLane(ctx, linuxLaneRequest{
		WorkDir: loaded.GitRoot, Head: strings.TrimSpace(head), Image: loaded.Config.Delivery.LinuxLaneImage,
		LogPath: filepath.Join(artifactDir, "delivery", "linux-lane.log"),
		Name:    fmt.Sprintf("roundfix-linux-lane-%d", os.Getpid()), All: *all,
	})
	if err != nil {
		return fail(err)
	}
	if result.Skip != "" {
		return 3
	}
	if result.Passed {
		return exitOK
	}
	fmt.Fprintln(stdout, "roundfix: "+result.Reason)
	if strings.HasPrefix(result.Blocker, delivery.BlockerLinuxLaneContainer+":") {
		return 3
	}
	return 1
}
