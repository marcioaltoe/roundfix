package cli

import (
	"context"
	"fmt"
	"io"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

// implementMachineVerificationGate attributes every lease to the checkout
// that started Implement, rather than to a transient Run or Task Worktree.
type implementMachineVerificationGate struct {
	homeDir    string
	repository string
	capacity   int
	stderr     io.Writer
}

func newImplementMachineVerificationGate(loaded roundconfig.Loaded, repository string, stderr io.Writer) *implementMachineVerificationGate {
	return &implementMachineVerificationGate{
		homeDir: loaded.HomeDir, repository: repository,
		capacity: loaded.Config.Runs.MaxVerifying, stderr: stderr,
	}
}

func (gate *implementMachineVerificationGate) Acquire(ctx context.Context, holder store.MachineVerificationHolder, onWait func(store.MachineVerificationWait)) (func(), error) {
	holder.Repository = gate.repository
	lease, err := store.AcquireMachineVerification(ctx, gate.homeDir, holder, gate.capacity, onWait)
	if err != nil {
		return nil, fmt.Errorf("acquire Implement Machine Verification Capacity: %w", err)
	}
	return func() {
		if err := lease.Release(); err != nil {
			fmt.Fprintf(gate.stderr, "roundfix: warning: %v\n", err)
		}
	}, nil
}
