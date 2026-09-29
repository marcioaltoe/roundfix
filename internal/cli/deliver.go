package cli

import (
	"context"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/delivery"
	"roundfix/internal/spec"
	"roundfix/internal/store"
)

const deliverUsage = `Usage:
  roundfix deliver plan [--json] [<slug>...]
  roundfix deliver start <slug>...
  roundfix deliver status
  roundfix deliver resume
  roundfix deliver retry <slug>
  roundfix deliver stop

Records and advances an ordered queue of Specs. The detached owner survives the
calling terminal. A blocked item is parked and the owner continues with the
next item.

Commands:
  plan    Report which Specs are approved to run
  start   Validate and record a new queue, then start its detached owner
  status  Print every queued Spec's stage, blocker and item worktree
  resume  Start a detached owner for the persisted queue
  retry   Return one parked Spec to the queue and reach its owner
  stop    Prove and terminate the persisted queue owner

Flags:
  --max-duration <duration>  Set a positive queue duration
  --max-retries <n>          Set the per-item retry limit to at least 1
`

var deliverStartValueFlags = map[string]bool{
	"max-duration": true,
	"max-retries":  true,
}

type deliverStartOptions struct {
	Slugs       []string
	MaxDuration time.Duration
	MaxRetries  int
}

type deliveryEngine interface {
	Run(context.Context, string) (delivery.EngineResult, error)
	Retry(context.Context, string, string) (delivery.RetryResult, error)
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
	case "plan":
		return runDeliverPlan(ctx, subcommandArgs, stdout, stderr, environment)
	case "start":
		return runDeliverStart(ctx, subcommandArgs, stdout, stderr, environment)
	case "status":
		return runDeliverStatus(ctx, subcommandArgs, stdout, stderr, environment)
	case "resume":
		if detachChild != nil {
			return runDeliveryOwner(ctx, subcommandArgs, stdout, stderr, detachChild, environment)
		}
		return runDeliverResume(ctx, subcommandArgs, stdout, stderr, environment)
	case "retry":
		return runDeliverRetry(ctx, subcommandArgs, stdout, stderr, environment)
	case "stop":
		return runDeliverStop(ctx, subcommandArgs, stdout, stderr, environment)
	default:
		fmt.Fprintf(stderr, "roundfix: unknown deliver command %q\n", subcommand)
		fmt.Fprintln(stderr, "Run 'roundfix deliver --help' for usage.")
		return exitPreflight
	}
}

func runDeliverRetry(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	specSlug, err := parseDeliverRetry(args)
	if err != nil {
		return printDeliverFailure("retry", err, stderr)
	}
	loaded, _, err := loadDeliveryCommand(ctx, environment, stderr)
	if err != nil {
		return printDeliverFailure("retry", err, stderr)
	}
	runStore, err := store.Open(ctx, loaded.HomeDir)
	if err != nil {
		return printDeliverFailure("retry", err, stderr)
	}
	storeClosed := false
	defer func() {
		if !storeClosed {
			_ = runStore.Close()
		}
	}()
	closeStore := func() error {
		if err := runStore.Close(); err != nil {
			return fmt.Errorf("close Run Database after retrying Delivery Queue item: %w", err)
		}
		storeClosed = true
		return nil
	}

	dependencies := commandDependenciesForContext(ctx)
	engine := dependencies.newDeliveryEngine(runStore, loaded)
	if engine == nil {
		return printDeliverFailure("retry", errors.New("Delivery Engine is required"), stderr)
	}
	result, err := engine.Retry(ctx, loaded.GitRoot, specSlug)
	if err != nil {
		return printDeliverFailure("retry", err, stderr)
	}

	if result.OwnerPID > 0 {
		ownerState := "is not running"
		if store.ProcessAlive(result.OwnerPID) {
			if dependencies.ownerProcesses == nil {
				return printDeliverFailure("retry", errors.New("Delivery Queue owner process controller is required"), stderr)
			}
			proofErr := dependencies.ownerProcesses.ProveOwner(ctx, result.OwnerPID, result.OwnerIdentity)
			switch {
			case proofErr == nil && store.ProcessAlive(result.OwnerPID):
				if err := closeStore(); err != nil {
					return printDeliverFailure("retry", err, stderr)
				}
				printDeliverRetryResult(stdout, specSlug, result)
				fmt.Fprintf(stdout, "Handed %s to Delivery Queue owner PID %d.\n", specSlug, result.OwnerPID)
				return exitOK
			case proofErr == nil:
				// The process exited between the liveness check and identity proof.
			case errors.Is(proofErr, store.ErrOwnerProcessIdentityUnproven):
				ownerState = "has a different process identity"
			default:
				return printDeliverFailure("retry", proofErr, stderr)
			}
		}
		released, releaseErr := runStore.ReleaseDeliveryQueueOwner(
			ctx,
			loaded.GitRoot,
			result.OwnerPID,
			result.OwnerIdentity,
		)
		if releaseErr != nil {
			return printDeliverFailure("retry", releaseErr, stderr)
		}
		if !released {
			return printDeliverFailure("retry", errors.New("Delivery Queue owner changed while retry was checking it"), stderr)
		}
		fmt.Fprintf(
			stderr,
			"roundfix: Delivery Queue owner PID %d %s; reclaimed its owner record.\n",
			result.OwnerPID,
			ownerState,
		)
	}
	if err := closeStore(); err != nil {
		return printDeliverFailure("retry", err, stderr)
	}
	printDeliverRetryResult(stdout, specSlug, result)
	return dependencies.startDeliveryOwner(ctx, loaded, environment, stdout, stderr)
}

func printDeliverRetryResult(stdout io.Writer, specSlug string, result delivery.RetryResult) {
	for _, carried := range result.CarriedFrom.Runs {
		fmt.Fprintf(
			stdout,
			"Carried forward from Run %s: %s\n",
			carried.RunID,
			strings.Join(carried.Carried, ", "),
		)
	}
	fmt.Fprintf(stdout, "Retried %s: %s -> %s\n", specSlug, result.Blocker, result.Stage)
}

func runDeliverStart(ctx context.Context, args []string, stdout, stderr io.Writer, environment commandEnvironment) int {
	options, err := parseDeliverStart(args)
	if err != nil {
		return printDeliverFailure("start", err, stderr)
	}
	loaded, specsRoot, err := loadDeliveryCommand(ctx, environment, stderr)
	if err != nil {
		return printDeliverFailure("start", err, stderr)
	}
	var authorizationRefusals []string
	for _, slug := range options.Slugs {
		if _, err := spec.Load(specsRoot.Path, slug); err != nil {
			return printDeliverFailure("start", err, stderr)
		}
		if reasons := deliveryAuthorizationReasons(ctx, loaded, specsRoot, slug); len(reasons) > 0 {
			authorizationRefusals = append(authorizationRefusals, fmt.Sprintf("%s: %s", slug, strings.Join(reasons, "; ")))
		}
	}
	if len(authorizationRefusals) > 0 {
		return printDeliverFailure(
			"start",
			fmt.Errorf("Delivery Queue contains Specs without delivery authority:\n  %s\nRun 'roundfix deliver plan %s' to inspect the prepared queue", strings.Join(authorizationRefusals, "\n  "), strings.Join(options.Slugs, " ")),
			stderr,
		)
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
	limits := store.DeliveryQueueLimits{MaxRetries: options.MaxRetries}
	if options.MaxDuration > 0 {
		limits.Deadline = time.Now().UTC().Add(options.MaxDuration).Truncate(time.Second)
	}
	queue, err := runStore.CreateDeliveryQueueWithLimits(ctx, loaded.GitRoot, options.Slugs, limits)
	if err != nil {
		_ = runStore.Close()
		return printDeliverFailure("start", err, stderr)
	}
	if err := runStore.Close(); err != nil {
		return printDeliverFailure("start", fmt.Errorf("close Run Database after recording Delivery Queue: %w", err), stderr)
	}
	printDeliveryLimits(stdout, queue.Limits)
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
		itemWorktree := item.Worktree
		if itemWorktree == "" {
			itemWorktree = "-"
		}
		fmt.Fprintf(stdout, "%s\t%s\t%s\t%s\n", item.SpecSlug, item.Stage, blocker, itemWorktree)
	}
	for _, item := range queue.Items {
		if item.Warning != "" {
			fmt.Fprintf(stdout, "Warning: %s %s\n", item.SpecSlug, item.Warning)
		}
	}
	printDeliveryLimits(stdout, queue.Limits)
	if question, found := delivery.PendingQuestionFor(queue); found {
		fmt.Fprintf(stdout, "Pending question: %s parked %s\n", question.SpecSlug, question.Blocker)
		fmt.Fprintf(stdout, "Answer: %s\n", question.Answer)
		if question.Waiting > 0 {
			fmt.Fprintf(stdout, "Waiting behind it: %d parked item(s)\n", question.Waiting)
		}
	}
	return exitOK
}

func printDeliveryLimits(output io.Writer, limits store.DeliveryQueueLimits) {
	deadline := "none"
	if !limits.Deadline.IsZero() {
		deadline = limits.Deadline.UTC().Format(time.RFC3339)
	}
	retries := "none"
	if limits.MaxRetries > 0 {
		retries = strconv.Itoa(limits.MaxRetries)
	}
	fmt.Fprintf(
		output,
		"Limits: deadline %s, retries per item %s, concurrency 1, spend not measured\n",
		deadline,
		retries,
	)
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
	for {
		if _, err := engine.Run(ctx, loaded.GitRoot); err != nil {
			return printDeliverFailure("resume", err, stderr)
		}
		released, err := runStore.ReleaseIdleDeliveryQueueOwner(ctx, loaded.GitRoot, pid, identity)
		if err != nil {
			return printDeliverFailure("resume", err, stderr)
		}
		if released {
			return exitOK
		}
	}
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

func parseDeliverStart(args []string) (deliverStartOptions, error) {
	fs := flagSet("deliver start")
	var options deliverStartOptions
	fs.DurationVar(&options.MaxDuration, "max-duration", 0, "maximum queue duration")
	fs.IntVar(&options.MaxRetries, "max-retries", 0, "maximum retries per item")
	if err := fs.Parse(hoistCommandFlags(args, deliverStartValueFlags)); err != nil {
		return deliverStartOptions{}, validationError{message: err.Error()}
	}
	setFlags := make(map[string]bool, len(deliverStartValueFlags))
	fs.Visit(func(flag *flag.Flag) {
		setFlags[flag.Name] = true
	})
	if setFlags["max-duration"] && options.MaxDuration <= 0 {
		return deliverStartOptions{}, validationError{message: "max-duration must be greater than zero"}
	}
	if setFlags["max-retries"] && options.MaxRetries < 1 {
		return deliverStartOptions{}, validationError{message: "max-retries must be at least 1"}
	}
	options.Slugs = fs.Args()
	if len(options.Slugs) == 0 {
		return deliverStartOptions{}, validationError{message: "missing required Spec slug; pass roundfix deliver start <slug>..."}
	}
	for index := range options.Slugs {
		options.Slugs[index] = strings.TrimSpace(options.Slugs[index])
		if options.Slugs[index] == "" {
			return deliverStartOptions{}, validationError{message: "Spec slug cannot be empty"}
		}
	}
	return options, nil
}

func parseDeliverRetry(args []string) (string, error) {
	fs := flagSet("deliver retry")
	if err := fs.Parse(hoistCommandFlags(args, nil)); err != nil {
		return "", validationError{message: err.Error()}
	}
	slugs := fs.Args()
	if len(slugs) == 0 {
		return "", validationError{message: "missing required Spec slug; pass roundfix deliver retry <slug>"}
	}
	if len(slugs) > 1 {
		return "", validationError{message: fmt.Sprintf("unexpected argument %q", slugs[1])}
	}
	specSlug := strings.TrimSpace(slugs[0])
	if specSlug == "" {
		return "", validationError{message: "Spec slug cannot be empty"}
	}
	return specSlug, nil
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
