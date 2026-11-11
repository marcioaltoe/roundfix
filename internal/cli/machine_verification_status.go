package cli

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

func machineVerificationCounts(view store.MachineVerificationView) (held, waiting int) {
	for _, entry := range view.Entries {
		switch entry.State {
		case "holding":
			held++
		case "waiting":
			waiting++
		}
	}
	return held, waiting
}

func printMachineVerificationStatus(ctx context.Context, loaded roundconfig.Loaded, stdout, stderr io.Writer) {
	view, err := store.ReadMachineVerification(ctx, loaded.HomeDir)
	if err != nil {
		fmt.Fprintf(stderr, "roundfix: warning: read Machine Verification Capacity: %v\n", err)
		return
	}
	if len(view.Entries) == 0 && len(view.Unreadable) == 0 {
		return
	}
	held, waiting := machineVerificationCounts(view)
	if loaded.Config.Runs.MaxVerifying == 0 {
		fmt.Fprintf(stdout, "Machine Verification Capacity: off, %d held, %d waiting\n", held, waiting)
	} else {
		fmt.Fprintf(stdout, "Machine Verification Capacity: %d of %d held, %d waiting\n", held, loaded.Config.Runs.MaxVerifying, waiting)
	}
	for _, entry := range view.Entries {
		if entry.State == "holding" {
			fmt.Fprintf(stdout, "  holding: %s, %s, since %s\n", store.DescribeMachineVerificationHolder(entry.Holder), entry.Holder.Step, entry.AcquiredAt.UTC().Format(time.RFC3339))
		}
	}
	n := 0
	for _, entry := range view.Entries {
		if entry.State == "waiting" {
			n++
			fmt.Fprintf(stdout, "  waiting %d: %s, %s, since %s\n", n, store.DescribeMachineVerificationHolder(entry.Holder), entry.Holder.Step, entry.Since.UTC().Format(time.RFC3339))
		}
	}
	for _, path := range view.Unreadable {
		fmt.Fprintf(stdout, "  unreadable: %s\n", path)
	}
}

func doctorVerificationCapacityCheck(ctx context.Context, loaded roundconfig.Loaded) CheckResult {
	result := CheckResult{Name: HealthCheckVerificationCapacity, Status: CheckStatusOK}
	view, err := store.ReadMachineVerification(ctx, loaded.HomeDir)
	if err != nil {
		result.Status = CheckStatusWarn
		result.Detail = "read Machine Verification Capacity: " + err.Error()
		return result
	}
	held, waiting := machineVerificationCounts(view)
	setting := fmt.Sprintf("runs.max_verifying %d", loaded.Config.Runs.MaxVerifying)
	if loaded.Config.Runs.MaxVerifying == 0 {
		setting = "off"
	}
	result.Detail = fmt.Sprintf("%s; %d held, %d waiting", setting, held, waiting)
	for _, entry := range view.Entries {
		if entry.State == "holding" {
			result.Detail += "; holding: " + store.DescribeMachineVerificationHolder(entry.Holder) + ", " + entry.Holder.Step
		}
	}
	n := 0
	for _, entry := range view.Entries {
		if entry.State == "waiting" {
			n++
			result.Detail += fmt.Sprintf("; waiting %d: %s, %s", n, store.DescribeMachineVerificationHolder(entry.Holder), entry.Holder.Step)
		}
	}
	var actions []string
	for _, path := range view.Unreadable {
		result.Status = CheckStatusWarn
		result.Detail += "; unreadable entry " + path
		actions = append(actions, "inspect "+path+" and remove it if no Roundfix process wrote it")
	}
	result.NextAction = strings.Join(actions, "; ")
	return result
}
