package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MachineVerificationHolderRun      = "run"
	MachineVerificationHolderDelivery = "delivery"
	machineVerificationSchema         = "roundfix/machine-verification-entry/v1"
	machineVerificationPollInterval   = 2 * time.Second
)

type MachineVerificationHolder struct {
	Kind       string
	RunID      string
	Repository string
	SpecSlug   string
	Step       string
}

type MachineVerificationEntry struct {
	Ticket        int64
	State         string
	Holder        MachineVerificationHolder
	OwnerPID      int
	OwnerIdentity string
	Since         time.Time
	AcquiredAt    time.Time
	Path          string
}

type MachineVerificationWait struct {
	Capacity int
	Holding  []MachineVerificationEntry
	Ahead    []MachineVerificationEntry
}

type MachineVerificationView struct {
	Entries    []MachineVerificationEntry
	Unreadable []string
}

// MachineVerificationLease owns one entry until Release removes it.
// The zero value is an off-capacity lease.
type MachineVerificationLease struct {
	mu      sync.Mutex
	release func() error
	waited  time.Duration
}

func (lease *MachineVerificationLease) Release() error {
	lease.mu.Lock()
	defer lease.mu.Unlock()
	if lease.release == nil {
		return nil
	}
	if err := lease.release(); err != nil {
		return err
	}
	lease.release = nil
	return nil
}

func (lease *MachineVerificationLease) Waited() time.Duration { return lease.waited }

type machineVerificationRecord struct {
	Schema        string     `json:"schema"`
	Ticket        int64      `json:"ticket"`
	State         string     `json:"state"`
	Kind          string     `json:"kind"`
	RunID         string     `json:"run_id"`
	Repository    string     `json:"repository"`
	Spec          string     `json:"spec"`
	Step          string     `json:"step"`
	OwnerPID      int        `json:"owner_pid"`
	OwnerIdentity string     `json:"owner_identity"`
	Since         time.Time  `json:"since"`
	AcquiredAt    *time.Time `json:"acquired_at,omitempty"`
}

// A closed generation wakes every local waiter. Capture the generation before
// inspecting the directory so a release between inspection and select is seen.
var machineVerificationWake = struct {
	sync.Mutex
	changed chan struct{}
}{changed: make(chan struct{})}

func machineVerificationChanged() <-chan struct{} {
	machineVerificationWake.Lock()
	defer machineVerificationWake.Unlock()
	return machineVerificationWake.changed
}

func wakeMachineVerificationWaiters() {
	machineVerificationWake.Lock()
	defer machineVerificationWake.Unlock()
	close(machineVerificationWake.changed)
	machineVerificationWake.changed = make(chan struct{})
}

func machineVerificationDirectory(homeDir string) string {
	return filepath.Join(homeDir, ".roundfix", "verification-capacity")
}

func AcquireMachineVerification(ctx context.Context, homeDir string, holder MachineVerificationHolder, capacity int, onWait func(MachineVerificationWait)) (*MachineVerificationLease, error) {
	if capacity <= 0 {
		return &MachineVerificationLease{}, nil
	}
	dir := machineVerificationDirectory(homeDir)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("wait for the Machine Verification Capacity: create directory: %w", err)
	}
	identity, identityErr := OwnerProcessIdentity(ctx, os.Getpid())
	if identityErr != nil {
		// An unreadable identity never proves death.
		identity = ""
	}
	entry := MachineVerificationEntry{State: "waiting", Holder: holder, OwnerPID: os.Getpid(), OwnerIdentity: identity, Since: time.Now().UTC()}
	err := withMachineVerificationLock(ctx, dir, func() error {
		ticket, err := nextMachineVerificationTicket(dir)
		if err != nil {
			return err
		}
		entry.Ticket = ticket
		entry.Path = filepath.Join(dir, fmt.Sprintf("%020d.json", ticket))
		return writeMachineVerificationEntry(entry)
	})
	if err != nil {
		return nil, fmt.Errorf("wait for the Machine Verification Capacity: %w", err)
	}
	lease := &MachineVerificationLease{release: func() error { return removeMachineVerificationEntry(ctx, entry.Path) }}
	var waitingSince time.Time
	ticker := time.NewTicker(machineVerificationPollInterval)
	defer ticker.Stop()
	for {
		changed := machineVerificationChanged()
		var granted bool
		wait := MachineVerificationWait{Capacity: capacity}
		err := withMachineVerificationLock(ctx, dir, func() error {
			view, err := readMachineVerificationEntries(ctx, dir)
			if err != nil {
				return err
			}
			for _, other := range view.Entries {
				if other.Path != entry.Path {
					if gone, _ := OwnerGone(ctx, other.OwnerPID, other.OwnerIdentity); gone {
						if err := os.Remove(other.Path); err != nil && !errors.Is(err, os.ErrNotExist) {
							return fmt.Errorf("remove entry %q: %w", other.Path, err)
						}
						wakeMachineVerificationWaiters()
						continue
					}
				}
				if other.State == "holding" {
					wait.Holding = append(wait.Holding, other)
				}
				if other.State == "waiting" && other.Ticket < entry.Ticket {
					wait.Ahead = append(wait.Ahead, other)
				}
			}
			if len(wait.Holding)+len(wait.Ahead) >= capacity {
				return nil
			}
			if err := ctx.Err(); err != nil {
				return err
			}
			entry.State = "holding"
			entry.AcquiredAt = time.Now().UTC()
			if err := writeMachineVerificationEntry(entry); err != nil {
				return err
			}
			granted = true
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("wait for the Machine Verification Capacity: %w", errors.Join(err, lease.Release()))
		}
		if granted {
			if !waitingSince.IsZero() {
				lease.waited = time.Since(waitingSince)
			}
			return lease, nil
		}
		if waitingSince.IsZero() {
			waitingSince = time.Now()
			if onWait != nil {
				onWait(wait)
			}
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("wait for the Machine Verification Capacity: %w", errors.Join(ctx.Err(), lease.Release()))
		case <-changed:
		case <-ticker.C:
		}
	}
}

func withMachineVerificationLock(ctx context.Context, dir string, action func() error) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}
	file, err := os.OpenFile(filepath.Join(dir, "lock"), os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		return fmt.Errorf("open capacity lock: %w", err)
	}
	defer func() { err = errors.Join(err, file.Close()) }()
	if err := acquireWriteLock(ctx, file); err != nil {
		return err
	}
	defer func() { err = errors.Join(err, releaseWriteLock(file)) }()
	if err := ctx.Err(); err != nil {
		return err
	}
	return action()
}

func nextMachineVerificationTicket(dir string) (int64, error) {
	path := filepath.Join(dir, "sequence")
	data, err := os.ReadFile(path)
	var last int64
	if err == nil {
		last, err = strconv.ParseInt(strings.TrimSpace(string(data)), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("read ticket sequence: %w", err)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return 0, fmt.Errorf("read ticket sequence: %w", err)
	}
	if last < 0 || last == math.MaxInt64 {
		return 0, fmt.Errorf("invalid ticket sequence %d", last)
	}
	ticket := last + 1
	if err := os.WriteFile(path, []byte(strconv.FormatInt(ticket, 10)+"\n"), 0o600); err != nil {
		return 0, fmt.Errorf("write ticket sequence: %w", err)
	}
	return ticket, nil
}

func writeMachineVerificationEntry(entry MachineVerificationEntry) (err error) {
	record := machineVerificationRecord{Schema: machineVerificationSchema, Ticket: entry.Ticket, State: entry.State, Kind: entry.Holder.Kind, RunID: entry.Holder.RunID, Repository: entry.Holder.Repository, Spec: entry.Holder.SpecSlug, Step: entry.Holder.Step, OwnerPID: entry.OwnerPID, OwnerIdentity: entry.OwnerIdentity, Since: entry.Since}
	if !entry.AcquiredAt.IsZero() {
		record.AcquiredAt = &entry.AcquiredAt
	}
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("encode entry %q: %w", entry.Path, err)
	}
	file, err := os.CreateTemp(filepath.Dir(entry.Path), ".entry-*")
	if err != nil {
		return fmt.Errorf("create entry %q: %w", entry.Path, err)
	}
	defer func() {
		if removeErr := os.Remove(file.Name()); removeErr != nil && !errors.Is(removeErr, os.ErrNotExist) {
			err = errors.Join(err, removeErr)
		}
	}()
	_, writeErr := file.Write(append(data, '\n'))
	if err := errors.Join(writeErr, file.Close()); err != nil {
		return fmt.Errorf("write entry %q: %w", entry.Path, err)
	}
	if err := os.Rename(file.Name(), entry.Path); err != nil {
		return fmt.Errorf("publish entry %q: %w", entry.Path, err)
	}
	return nil
}

func removeMachineVerificationEntry(ctx context.Context, path string) error {
	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
	defer cancel()
	err := withMachineVerificationLock(cleanupCtx, filepath.Dir(path), func() error {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		wakeMachineVerificationWaiters()
		return nil
	})
	if err != nil {
		return fmt.Errorf("remove entry %q: %w", path, err)
	}
	return nil
}

// readMachineVerificationEntries preserves unreadable entries without counting
// them. Its caller decides whether a proven-gone entry is omitted or removed.
func readMachineVerificationEntries(ctx context.Context, dir string) (MachineVerificationView, error) {
	var view MachineVerificationView
	files, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return view, nil
	}
	if err != nil {
		return view, fmt.Errorf("read capacity directory: %w", err)
	}
	for _, file := range files {
		if err := ctx.Err(); err != nil {
			return view, err
		}
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}
		path := filepath.Join(dir, file.Name())
		data, readErr := os.ReadFile(path)
		var record machineVerificationRecord
		if readErr != nil || json.Unmarshal(data, &record) != nil || record.Schema != machineVerificationSchema {
			view.Unreadable = append(view.Unreadable, path)
			continue
		}
		entry := MachineVerificationEntry{Ticket: record.Ticket, State: record.State, Holder: MachineVerificationHolder{Kind: record.Kind, RunID: record.RunID, Repository: record.Repository, SpecSlug: record.Spec, Step: record.Step}, OwnerPID: record.OwnerPID, OwnerIdentity: record.OwnerIdentity, Since: record.Since, Path: path}
		if record.AcquiredAt != nil {
			entry.AcquiredAt = *record.AcquiredAt
		}
		view.Entries = append(view.Entries, entry)
	}
	slices.SortFunc(view.Entries, func(a, b MachineVerificationEntry) int {
		if a.Ticket < b.Ticket {
			return -1
		}
		if a.Ticket > b.Ticket {
			return 1
		}
		return 0
	})
	return view, nil
}

// ReadMachineVerification takes no lock and never changes the entry directory.
func ReadMachineVerification(ctx context.Context, homeDir string) (MachineVerificationView, error) {
	view, err := readMachineVerificationEntries(ctx, machineVerificationDirectory(homeDir))
	if err != nil {
		return view, err
	}
	live := make([]MachineVerificationEntry, 0, len(view.Entries))
	for _, entry := range view.Entries {
		if err := ctx.Err(); err != nil {
			return view, err
		}
		if gone, _ := OwnerGone(ctx, entry.OwnerPID, entry.OwnerIdentity); !gone {
			live = append(live, entry)
		}
	}
	view.Entries = live
	return view, nil
}

func DescribeMachineVerificationHolder(holder MachineVerificationHolder) string {
	if holder.Kind == MachineVerificationHolderRun {
		return fmt.Sprintf("Run %s of %s in %s", holder.RunID, holder.SpecSlug, holder.Repository)
	}
	return fmt.Sprintf("Delivery Queue item %s in %s", holder.SpecSlug, holder.Repository)
}
