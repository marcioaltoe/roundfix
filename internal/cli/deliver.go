package cli

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

const deliverUsage = `Usage:
  roundfix deliver start <slug>...
  roundfix deliver status
  roundfix deliver resume
  roundfix deliver stop

Records and advances an ordered queue of Specs. The detached owner survives the
calling terminal. A blocked item is parked and the owner continues with the
next item.

Commands:
  start   Validate and record a new queue, then start its detached owner
  status  Print every queued Spec's stage and blocker
  resume  Start a detached owner for the persisted queue
  stop    Prove and terminate the persisted queue owner
`

type deliveryEngine interface {
	Run(context.Context, string) (delivery.EngineResult, error)
}

func runDeliverCommand(
	ctx context.Context,
	args []string,
	stdout, stderr io.Writer,
	detachChild *detachChild,
	environment commandEnvironment,
) int {
	if commandWantsHelp(args) || len(args) == 0 {
		fmt.Fprint(stdout, deliverUsage)
		return exitOK
	}

	subcommand := args[0]
	subcommandArgs := args[1:]
	switch subcommand {
	case "start":
		return runDeliverStart(ctx, subcommandArgs, stdout, stderr, environment)
	case "status":
		return runDeliverStatus(ctx, subcommandArgs, stdout, stderr, environment)
	case "resume":
		if detachChild != nil {
			return runDeliveryOwner(ctx, subcommandArgs, stdout, stderr, detachChild, environment)
		}
		return runDeliverResume(ctx, subcommandArgs, stdout, stderr, environment)
	case "stop":
		return runDeliverStop(ctx, subcommandArgs, stdout, stderr, environment)
	default:
		fmt.Fprintf(stderr, "roundfix: unknown deliver command %q\n", subcommand)
		fmt.Fprintln(stderr, "Run 'roundfix deliver --help' for usage.")
		return exitPreflight
	}
}

func runDeliverStart(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	slugs, err := parseDeliverStart(args)
	if err != nil {
		return printDeliverFailure("start", err, stderr)
	}
	loaded, specsRoot, err := loadDeliveryCommand(ctx, environment, stderr)
	if err != nil {
		return printDeliverFailure("start", err, stderr)
	}
	for _, slug := range slugs {
		if _, err := spec.Load(specsRoot.Path, slug); err != nil {
			return printDeliverFailure("start", err, stderr)
		}
	}

	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		return printDeliverFailure("start", err, stderr)
	}
	if existing, found, readErr := runStore.DeliveryQueue(ctx, loaded.GitRoot); readErr != nil {
		_ = runStore.Close()
		return printDeliverFailure("start", readErr, stderr)
	} else if found && existing.OwnerPID > 0 && !store.ProcessAlive(existing.OwnerPID) {
		released, releaseErr := runStore.ReleaseDeliveryQueueOwner(ctx, loaded.GitRoot, existing.OwnerPID, existing.OwnerIdentity)
		if releaseErr != nil || !released {
			_ = runStore.Close()
			if releaseErr == nil {
				releaseErr = errors.New("Delivery Queue owner changed while start was checking it")
			}
			return printDeliverFailure("start", releaseErr, stderr)
		}
	}
	if _, err := runStore.CreateDeliveryQueue(ctx, loaded.GitRoot, slugs); err != nil {
		_ = runStore.Close()
		return printDeliverFailure("start", err, stderr)
	}
	if err := runStore.Close(); err != nil {
		return printDeliverFailure("start", fmt.Errorf("close Run Database after recording Delivery Queue: %w", err), stderr)
	}
	return commandDependenciesForContext(ctx).startDeliveryOwner(ctx, loaded, environment, stdout, stderr)
}

func runDeliverStatus(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if err := parseDeliverNoArgs("status", args); err != nil {
		return printDeliverFailure("status", err, stderr)
	}
	loaded, _, err := loadDeliveryCommand(ctx, environment, stderr)
	if err != nil {
		return printDeliverFailure("status", err, stderr)
	}
	runStore, err := store.OpenReader(ctx, loaded.HomeDir)
	if err != nil {
		return printDeliverFailure("status", err, stderr)
	}
	defer func() {
		_ = runStore.Close()
	}()
	queue, found, err := runStore.DeliveryQueue(ctx, loaded.GitRoot)
	if err != nil {
		return printDeliverFailure("status", err, stderr)
	}
	if !found {
		return printDeliverFailure("status", fmt.Errorf("Delivery Queue for repository %q does not exist", loaded.GitRoot), stderr)
	}
	for _, item := range queue.Items {
		blocker := item.Blocker
		if blocker == "" {
			blocker = "-"
		}
		fmt.Fprintf(stdout, "%s\t%s\t%s\n", item.SpecSlug, item.Stage, blocker)
	}
	return exitOK
}

func runDeliverResume(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if err := parseDeliverNoArgs("resume", args); err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	loaded, _, err := loadDeliveryCommand(ctx, environment, stderr)
	if err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	defer func() {
		_ = runStore.Close()
	}()
	queue, found, err := runStore.DeliveryQueue(ctx, loaded.GitRoot)
	if err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	if !found {
		return printDeliverFailure("resume", fmt.Errorf("Delivery Queue for repository %q does not exist", loaded.GitRoot), stderr)
	}
	if queue.OwnerPID > 0 {
		ownerState := "is not running"
		if store.ProcessAlive(queue.OwnerPID) {
			controller := commandDependenciesForContext(ctx).ownerProcesses
			if controller == nil {
				return printDeliverFailure("resume", errors.New("Delivery Queue owner process controller is required"), stderr)
			}
			proofErr := controller.ProveOwner(ctx, queue.OwnerPID, queue.OwnerIdentity)
			switch {
			case proofErr == nil && store.ProcessAlive(queue.OwnerPID):
				return printDeliverFailure("resume", fmt.Errorf("Delivery Queue already has owner PID %d; run 'roundfix deliver stop' first", queue.OwnerPID), stderr)
			case proofErr == nil:
				// The process exited between the liveness check and identity proof.
			case errors.Is(proofErr, store.ErrOwnerProcessIdentityUnproven):
				ownerState = "has a different process identity"
			default:
				return printDeliverFailure("resume", proofErr, stderr)
			}
		}
		released, releaseErr := runStore.ReleaseDeliveryQueueOwner(ctx, loaded.GitRoot, queue.OwnerPID, queue.OwnerIdentity)
		if releaseErr != nil {
			return printDeliverFailure("resume", releaseErr, stderr)
		}
		if !released {
			return printDeliverFailure("resume", errors.New("Delivery Queue owner changed while resume was checking it"), stderr)
		}
		fmt.Fprintf(stderr, "roundfix: Delivery Queue owner PID %d %s; reclaimed its owner record.\n", queue.OwnerPID, ownerState)
	}
	return commandDependenciesForContext(ctx).startDeliveryOwner(ctx, loaded, environment, stdout, stderr)
}

func runDeliverStop(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	if err := parseDeliverNoArgs("stop", args); err != nil {
		return printDeliverFailure("stop", err, stderr)
	}
	loaded, _, err := loadDeliveryCommand(ctx, environment, stderr)
	if err != nil {
		return printDeliverFailure("stop", err, stderr)
	}
	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		return printDeliverFailure("stop", err, stderr)
	}
	defer func() {
		_ = runStore.Close()
	}()
	queue, found, err := runStore.DeliveryQueue(ctx, loaded.GitRoot)
	if err != nil {
		return printDeliverFailure("stop", err, stderr)
	}
	if !found {
		return printDeliverFailure("stop", fmt.Errorf("Delivery Queue for repository %q does not exist", loaded.GitRoot), stderr)
	}
	if queue.OwnerPID == 0 {
		fmt.Fprintln(stdout, "Delivery Queue owner is not running.")
		return exitOK
	}
	controller := commandDependenciesForContext(ctx).ownerProcesses
	if controller == nil {
		return printDeliverFailure("stop", errors.New("Delivery Queue owner process controller is required"), stderr)
	}
	if err := controller.ProveOwner(ctx, queue.OwnerPID, queue.OwnerIdentity); err != nil {
		return printDeliverFailure("stop", err, stderr)
	}
	if _, err := controller.TerminateTreeAndWait(ctx, queue.OwnerPID, queue.OwnerIdentity); err != nil {
		return printDeliverFailure("stop", err, stderr)
	}
	released, err := runStore.ReleaseDeliveryQueueOwner(context.WithoutCancel(ctx), loaded.GitRoot, queue.OwnerPID, queue.OwnerIdentity)
	if err != nil {
		return printDeliverFailure("stop", err, stderr)
	}
	if !released {
		current, currentFound, readErr := runStore.DeliveryQueue(ctx, loaded.GitRoot)
		if readErr != nil {
			return printDeliverFailure("stop", readErr, stderr)
		}
		if currentFound && current.OwnerPID != 0 {
			return printDeliverFailure("stop", errors.New("Delivery Queue owner changed while stop was in progress"), stderr)
		}
	}
	fmt.Fprintf(stdout, "Stopped Delivery Queue owner PID %d.\n", queue.OwnerPID)
	return exitOK
}

func runDeliveryOwner(
	ctx context.Context,
	args []string,
	_ io.Writer,
	stderr io.Writer,
	detachChild *detachChild,
	environment commandEnvironment,
) int {
	if err := parseDeliverNoArgs("resume", args); err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	loaded, _, err := loadDeliveryCommand(ctx, environment, stderr)
	if err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	defer func() {
		_ = runStore.Close()
	}()

	pid := os.Getpid()
	identity, err := store.OwnerProcessIdentity(ctx, pid)
	if err != nil {
		return printDeliverFailure("resume", fmt.Errorf("read Delivery Queue owner identity: %w", err), stderr)
	}
	if err := runStore.ClaimDeliveryQueueOwner(ctx, loaded.GitRoot, pid, identity); err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	defer func() {
		_, _ = runStore.ReleaseDeliveryQueueOwner(context.WithoutCancel(ctx), loaded.GitRoot, pid, identity)
	}()

	artifactDir, err := roundconfig.ValidateArtifactDirectory(loaded.Config.Defaults.ArtifactDir, loaded.GitRoot, loaded.HomeDir)
	if err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	ownerID := deliveryOwnerID(loaded.GitRoot, pid)
	consoleLog := filepath.Join(artifactDir, "delivery", ownerID, "console.log")
	if err := detachChild.reportStarted(ownerID, consoleLog); err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	engine := commandDependenciesForContext(ctx).newDeliveryEngine(runStore, loaded)
	if engine == nil {
		return printDeliverFailure("resume", errors.New("Delivery Engine is required"), stderr)
	}
	if _, err := engine.Run(ctx, loaded.GitRoot); err != nil {
		return printDeliverFailure("resume", err, stderr)
	}
	return exitOK
}

func startDetachedDeliveryOwner(
	ctx context.Context,
	loaded roundconfig.Loaded,
	environment commandEnvironment,
	stdout, stderr io.Writer,
) int {
	req := commandRequest{name: "deliver", artifactDir: loaded.Config.Defaults.ArtifactDir}
	return runDetachedCommandWithReport(
		[]string{"deliver", "resume"},
		req,
		loaded,
		stdout,
		stderr,
		environment.environ,
		environment.workDir,
		environment.dependencies.detachTimeouts,
		printDetachedDeliveryReport,
	)
}

func printDetachedDeliveryReport(stdout io.Writer, ownerID string, consoleLog string) {
	fmt.Fprintf(stdout, "Delivery Owner: %s\n", ownerID)
	fmt.Fprintf(stdout, "Console Log: %s\n", consoleLog)
	fmt.Fprintln(stdout, "Status: roundfix deliver status")
	fmt.Fprintln(stdout, "Stop: roundfix deliver stop")
}

func deliveryOwnerID(gitRoot string, pid int) string {
	digest := sha256.Sum256([]byte(filepath.Clean(gitRoot)))
	return fmt.Sprintf("delivery-%x-%d", digest[:6], pid)
}

func loadDeliveryCommand(
	ctx context.Context,
	environment commandEnvironment,
	stderr io.Writer,
) (roundconfig.Loaded, roundconfig.SpecsRoot, error) {
	if err := ctx.Err(); err != nil {
		return roundconfig.Loaded{}, roundconfig.SpecsRoot{}, err
	}
	loaded, err := loadCommandConfig(environment, stderr)
	if err != nil {
		return roundconfig.Loaded{}, roundconfig.SpecsRoot{}, err
	}
	if strings.TrimSpace(loaded.GitRoot) == "" {
		return roundconfig.Loaded{}, roundconfig.SpecsRoot{}, errors.New("deliver requires a git repository working tree")
	}
	specsRoot, err := roundconfig.ResolveSpecsRoot(loaded, loaded.GitRoot)
	if err != nil {
		return roundconfig.Loaded{}, roundconfig.SpecsRoot{}, err
	}
	return loaded, specsRoot, nil
}

func parseDeliverStart(args []string) ([]string, error) {
	fs := flagSet("deliver start")
	if err := fs.Parse(hoistCommandFlags(args, nil)); err != nil {
		return nil, validationError{message: err.Error()}
	}
	slugs := fs.Args()
	if len(slugs) == 0 {
		return nil, validationError{message: "missing required Spec slug; pass roundfix deliver start <slug>..."}
	}
	for index := range slugs {
		slugs[index] = strings.TrimSpace(slugs[index])
		if slugs[index] == "" {
			return nil, validationError{message: "Spec slug cannot be empty"}
		}
	}
	return slugs, nil
}

func parseDeliverNoArgs(subcommand string, args []string) error {
	fs := flagSet("deliver " + subcommand)
	if err := fs.Parse(hoistCommandFlags(args, nil)); err != nil {
		return validationError{message: err.Error()}
	}
	if remaining := fs.Args(); len(remaining) > 0 {
		return validationError{message: fmt.Sprintf("unexpected argument %q", remaining[0])}
	}
	return nil
}

func printDeliverFailure(subcommand string, err error, stderr io.Writer) int {
	printPreflightFailure("deliver "+subcommand, err, stderr)
	return exitPreflight
}
