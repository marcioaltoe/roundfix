package store

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"slices"
	"strings"
	"time"
)

const (
	ownerProcessGracePeriod  = 2 * time.Second
	ownerProcessStopWindow   = 10 * time.Second
	ownerProcessPollInterval = 25 * time.Millisecond
)

var (
	ErrOwnerProcessUnsupported      = errors.New("owner process control is unsupported")
	ErrOwnerIdentityUnreadable      = errors.New("owner process identity is unreadable")
	ErrOwnerProcessIdentityUnproven = errors.New("owner process identity cannot be proven")
	errOwnerProcessAlreadyAbsent    = errors.New("owner process is already absent")
	errGracefulSignalUnsupported    = errors.New("graceful owner process termination is unsupported")
)

type OwnerProcessControlError struct {
	PID  int
	Step string
	Err  error
}

func (err OwnerProcessControlError) Error() string {
	return fmt.Sprintf("owner process %d failed step %q: %v", err.PID, err.Step, err.Err)
}

func (err OwnerProcessControlError) Unwrap() error {
	return err.Err
}

type OwnerProcessControl struct {
	gracePeriod          time.Duration
	stopWindow           time.Duration
	pollInterval         time.Duration
	processAbsent        func(int) (bool, error)
	processStartIdentity func(context.Context, int) (string, error)
	ownedProcesses       func(int, string) ([]int, error)
	processInfo          func(context.Context, int) (OwnedProcess, error)
	signalProcess        func(int, bool) error
}

// OwnedProcess is one live process proven to belong to a recorded Run owner.
// Command is diagnostic context only and must not be used as ownership proof.
type OwnedProcess struct {
	PID     int
	Started time.Time
	CPUTime time.Duration
	Command string
}

// TerminationOutcome records whether one owned process was observed absent.
// Reason is empty when Proven is true and explains why absence could not be
// proven otherwise.
type TerminationOutcome struct {
	PID    int
	Proven bool
	Reason string
}

func NewOwnerProcessController() *OwnerProcessControl {
	return newOwnerProcessController(
		ownerProcessGracePeriod,
		ownerProcessStopWindow,
		ownerProcessPollInterval,
	)
}

func newOwnerProcessController(gracePeriod, stopWindow, pollInterval time.Duration) *OwnerProcessControl {
	return &OwnerProcessControl{
		gracePeriod:          gracePeriod,
		stopWindow:           stopWindow,
		pollInterval:         pollInterval,
		processAbsent:        processAbsent,
		processStartIdentity: processStartIdentity,
		ownedProcesses:       processTreePIDs,
		processInfo:          readOwnedProcess,
		signalProcess:        signalOwnerProcess,
	}
}

// ProveOwner proves that pid still identifies the process the Run recorded as
// its owner, without sending any signal or touching any other state. Callers
// run it before any destructive Force Stop step so a Run whose owner cannot be
// proven is left exactly as it was found. Success includes a proven-absent
// owner: absence is its own proof of exit.
func (controller *OwnerProcessControl) ProveOwner(ctx context.Context, pid int, recordedIdentity string) error {
	_, err := controller.proveOwner(ctx, pid, recordedIdentity)
	return err
}

// proveOwner is the single implementation of the ownership rule. It reports
// whether the owner is already absent, which every caller treats as success.
//
// When the Run recorded a start-time identity token produced by the current
// platform, the live process must present the same token, otherwise a reused
// PID would let Force Stop terminate an unrelated process. A token from an
// older or foreign identity source is unreadable rather than proof of reuse.
// An empty recorded token comes from a Run created before identity recording
// existed and keeps the legacy PID-only proof, mirroring the ADR-0044
// precedent that PID-less legacy Runs degrade gracefully.
func (controller *OwnerProcessControl) proveOwner(ctx context.Context, pid int, recordedIdentity string) (bool, error) {
	if pid <= 0 || pid == os.Getpid() {
		return false, ownerProcessControlError(pid, "prove owner process identity", ErrOwnerProcessIdentityUnproven)
	}
	absent, err := controller.processAbsent(pid)
	if err != nil {
		return false, ownerProcessControlError(pid, "prove owner process identity", err)
	}
	if absent {
		return true, nil
	}
	if recorded := strings.TrimSpace(recordedIdentity); recorded != "" {
		liveIdentity, identityErr := controller.processStartIdentity(ctx, pid)
		if identityErr != nil {
			// The owner may have exited between the liveness check and the
			// identity read; proven absence is the proof Force Stop needs.
			if absent, absentErr := controller.processAbsent(pid); absentErr == nil && absent {
				return true, nil
			}
			return false, ownerProcessControlError(pid, "prove owner process identity",
				fmt.Errorf("%w: read live owner identity: %w; resolve the host resource failure, then retry", ErrOwnerIdentityUnreadable, identityErr))
		}
		if !strings.HasPrefix(recorded, runtime.GOOS+":") {
			return false, ownerProcessControlError(pid, "prove owner process identity",
				fmt.Errorf("%w: recorded owner identity is not comparable on %s", ErrOwnerIdentityUnreadable, runtime.GOOS))
		}
		if liveIdentity != recorded {
			return false, ownerProcessControlError(pid, "prove owner process identity",
				fmt.Errorf("%w: live process start identity does not match the recorded owner identity", ErrOwnerProcessIdentityUnproven))
		}
	}
	return false, nil
}

// TerminateAndWait reuses the same ownership proof as ProveOwner before
// sending any signal, then terminates the owner and returns only once its exit
// is proven. Callers that already ran ProveOwner pay for a second read-only
// proof, which keeps this entry point safe on its own.
func (controller *OwnerProcessControl) TerminateAndWait(ctx context.Context, pid int, recordedIdentity string) error {
	absent, err := controller.proveOwner(ctx, pid, recordedIdentity)
	if err != nil {
		return err
	}
	if absent {
		return nil
	}

	stopCtx, cancel := context.WithTimeout(ctx, controller.stopWindow)
	defer cancel()

	gracefulErr := controller.signalProcess(pid, false)
	switch {
	case errors.Is(gracefulErr, errOwnerProcessAlreadyAbsent):
		return nil
	case gracefulErr == nil:
		absent, err = controller.waitForAbsence(stopCtx, pid, controller.gracePeriod)
		if err != nil {
			return ownerProcessControlError(pid, "prove exit after graceful termination", err)
		}
		if absent {
			return nil
		}
	case errors.Is(gracefulErr, errGracefulSignalUnsupported):
	default:
		return ownerProcessControlError(pid, "send graceful termination", gracefulErr)
	}

	if err := controller.signalProcess(pid, true); err != nil {
		if errors.Is(err, errOwnerProcessAlreadyAbsent) {
			return nil
		}
		return ownerProcessControlError(pid, "send force kill", err)
	}
	absent, err = controller.waitForAbsence(stopCtx, pid, 0)
	if err != nil {
		return ownerProcessControlError(pid, "prove exit after force kill", err)
	}
	if !absent {
		return ownerProcessControlError(pid, "prove exit after force kill", context.DeadlineExceeded)
	}
	return nil
}

// InspectTree proves the recorded owner without signalling it, then returns
// the live processes the platform still attributes to that owner. A host that
// cannot prove ownership, enumerate the tree, or observe liveness returns an
// error instead of an incomplete reclaimable set.
func (controller *OwnerProcessControl) InspectTree(
	ctx context.Context,
	pid int,
	recordedIdentity string,
) ([]int, error) {
	ownerAbsent, err := controller.proveOwner(ctx, pid, recordedIdentity)
	if err != nil {
		return nil, err
	}
	owned, err := controller.ownedProcesses(pid, recordedIdentity)
	if err != nil {
		return nil, ownerProcessControlError(pid, "enumerate owned process tree", err)
	}

	live := make([]int, 0, len(owned))
	for _, ownedPID := range normalizeProcessTree(pid, owned) {
		if ownedPID == pid && ownerAbsent {
			continue
		}
		absent, err := controller.processAbsent(ownedPID)
		if err != nil {
			return nil, ownerProcessControlError(ownedPID, "inspect owned process liveness", err)
		}
		if !absent {
			live = append(live, ownedPID)
		}
	}
	return live, nil
}

// InspectTreeProcesses returns details only for processes whose membership in
// the recorded owner's spawn lineage has already been proven by InspectTree.
// A process that exits between the lineage and detail reads is omitted; other
// detail failures are returned with every process that was readable.
func (controller *OwnerProcessControl) InspectTreeProcesses(
	ctx context.Context,
	pid int,
	recordedIdentity string,
) ([]OwnedProcess, error) {
	pids, err := controller.InspectTree(ctx, pid, recordedIdentity)
	if err != nil {
		return nil, err
	}

	processes := make([]OwnedProcess, 0, len(pids))
	var readErrors []error
	for _, ownedPID := range pids {
		process, err := controller.processInfo(ctx, ownedPID)
		if err == nil {
			process.PID = ownedPID
			processes = append(processes, process)
			continue
		}
		absent, absentErr := controller.processAbsent(ownedPID)
		if absentErr == nil && absent {
			continue
		}
		if absentErr != nil {
			err = errors.Join(err, fmt.Errorf("prove process absence after detail read: %w", absentErr))
		}
		readErrors = append(readErrors, ownerProcessControlError(ownedPID, "read owned process details", err))
	}
	return processes, errors.Join(readErrors...)
}

// TerminateTreeAndWait proves the recorded owner before signalling any
// process, then terminates every process the platform can still attribute to
// that owner. It reports proof of absence per PID; a process whose absence
// cannot be observed remains explicitly unproven in its outcome.
func (controller *OwnerProcessControl) TerminateTreeAndWait(
	ctx context.Context,
	pid int,
	recordedIdentity string,
) ([]TerminationOutcome, error) {
	ownerAbsent, err := controller.proveOwner(ctx, pid, recordedIdentity)
	if err != nil {
		return nil, err
	}

	outcomes := make([]TerminationOutcome, 0, 1)
	seen := make(map[int]struct{})
	if ownerAbsent {
		seen[pid] = struct{}{}
		outcomes = append(outcomes, TerminationOutcome{PID: pid, Proven: true})
	}

	for {
		if err := ctx.Err(); err != nil {
			return outcomes, ownerProcessControlError(pid, "terminate owned process tree", err)
		}
		owned, err := controller.ownedProcesses(pid, recordedIdentity)
		if err != nil {
			treeErr := ownerProcessControlError(pid, "enumerate owned process tree", err)
			if _, ownerRecorded := seen[pid]; !ownerRecorded {
				seen[pid] = struct{}{}
				outcomes = append(outcomes, unprovenTerminationOutcome(pid, treeErr))
			}
			return outcomes, treeErr
		}
		owned = normalizeProcessTree(pid, owned)

		discovered := 0
		for _, ownedPID := range owned {
			if _, ok := seen[ownedPID]; ok {
				continue
			}
			seen[ownedPID] = struct{}{}
			discovered++

			identity, absent, identityErr := controller.ownedProcessIdentity(ctx, pid, ownedPID, recordedIdentity)
			if identityErr != nil {
				outcomes = append(outcomes, unprovenTerminationOutcome(ownedPID, identityErr))
				continue
			}
			if absent {
				outcomes = append(outcomes, TerminationOutcome{PID: ownedPID, Proven: true})
				continue
			}
			if err := controller.TerminateAndWait(ctx, ownedPID, identity); err != nil {
				outcomes = append(outcomes, unprovenTerminationOutcome(ownedPID, err))
				continue
			}
			outcomes = append(outcomes, TerminationOutcome{PID: ownedPID, Proven: true})
		}
		if discovered == 0 {
			return outcomes, nil
		}
	}
}

func (controller *OwnerProcessControl) ownedProcessIdentity(
	ctx context.Context,
	ownerPID int,
	pid int,
	recordedOwnerIdentity string,
) (string, bool, error) {
	if pid == ownerPID {
		return recordedOwnerIdentity, false, nil
	}
	identity, err := controller.processStartIdentity(ctx, pid)
	if err == nil {
		return identity, false, nil
	}
	if errors.Is(err, ErrOwnerProcessUnsupported) {
		return "", false, nil
	}
	absent, absentErr := controller.processAbsent(pid)
	if absentErr == nil && absent {
		return "", true, nil
	}
	if absentErr != nil {
		return "", false, fmt.Errorf("read owned process identity: %w; prove absence after identity read: %v", err, absentErr)
	}
	return "", false, fmt.Errorf("read owned process identity: %w", err)
}

func normalizeProcessTree(ownerPID int, pids []int) []int {
	unique := make(map[int]struct{}, len(pids)+1)
	unique[ownerPID] = struct{}{}
	for _, pid := range pids {
		if pid > 0 {
			unique[pid] = struct{}{}
		}
	}
	result := make([]int, 0, len(unique))
	for pid := range unique {
		if pid != ownerPID {
			result = append(result, pid)
		}
	}
	slices.Sort(result)
	return append([]int{ownerPID}, result...)
}

func unprovenTerminationOutcome(pid int, err error) TerminationOutcome {
	return TerminationOutcome{PID: pid, Reason: err.Error()}
}

type processParent struct {
	pid       int
	parentPID int
}

type processParentWithStart struct {
	pid       int
	parentPID int
	started   uint64
}

func descendantProcessPIDs(ownerPID int, processes []processParent) []int {
	owned := map[int]struct{}{ownerPID: {}}
	for {
		added := false
		for _, process := range processes {
			if _, ok := owned[process.pid]; ok {
				continue
			}
			if _, ok := owned[process.parentPID]; !ok {
				continue
			}
			owned[process.pid] = struct{}{}
			added = true
		}
		if !added {
			break
		}
	}
	result := make([]int, 0, len(owned))
	for pid := range owned {
		result = append(result, pid)
	}
	return result
}

// descendantProcessPIDsAfterStart accepts only parent-child links whose child
// was created no earlier than its claimed parent. This prevents a stale parent
// PID left by process reuse from attributing an older, unrelated process tree
// to the Run owner.
func descendantProcessPIDsAfterStart(
	ownerPID int,
	ownerStarted uint64,
	processes []processParentWithStart,
) []int {
	owned := map[int]uint64{ownerPID: ownerStarted}
	for {
		added := false
		for _, process := range processes {
			if _, ok := owned[process.pid]; ok {
				continue
			}
			parentStarted, ok := owned[process.parentPID]
			if !ok || process.started < parentStarted {
				continue
			}
			owned[process.pid] = process.started
			added = true
		}
		if !added {
			break
		}
	}
	result := make([]int, 0, len(owned))
	for pid := range owned {
		result = append(result, pid)
	}
	slices.Sort(result)
	return result
}

func (controller *OwnerProcessControl) waitForAbsence(ctx context.Context, pid int, window time.Duration) (bool, error) {
	pollInterval := controller.pollInterval
	if pollInterval <= 0 {
		pollInterval = ownerProcessPollInterval
	}
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	var timer *time.Timer
	var windowElapsed <-chan time.Time
	if window > 0 {
		timer = time.NewTimer(window)
		windowElapsed = timer.C
		defer timer.Stop()
	}

	for {
		absent, err := controller.processAbsent(pid)
		if err != nil {
			return false, err
		}
		if absent {
			return true, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-windowElapsed:
			return false, nil
		case <-ticker.C:
		}
	}
}

func ownerProcessControlError(pid int, step string, err error) error {
	return OwnerProcessControlError{PID: pid, Step: step, Err: err}
}

// ProcessAlive reports whether pid is alive or cannot be proven absent.
func ProcessAlive(pid int) bool {
	if pid <= 0 {
		return true
	}
	absent, err := processAbsent(pid)
	return err != nil || !absent
}

// OwnerProcessIdentity returns the opaque start-time identity token for pid.
// Callers compare tokens verbatim and must not parse them.
func OwnerProcessIdentity(ctx context.Context, pid int) (string, error) {
	if pid <= 0 {
		return "", fmt.Errorf("read owner process identity: pid %d is invalid", pid)
	}
	return processStartIdentity(ctx, pid)
}
