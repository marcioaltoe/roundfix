// Suite: User Config Run Slot Wait.
// Invariant: only the machine owner controls a nonnegative bounded wait.
// Boundary IN: layered YAML configuration in temporary homes.
// Boundary OUT: no host configuration or external tools.
package config

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunSlotWaitIsReadFromUserConfig(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		value string
		want  time.Duration
	}{{"", 2 * time.Hour}, {"90m", 90 * time.Minute}, {"0", 0}} {
		t.Run(tc.value, func(t *testing.T) {
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			if tc.value != "" {
				mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "runs:\n  slot_wait: "+tc.value+"\n")
			}
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work})
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Runs.SlotWait != tc.want {
				t.Fatalf("SlotWait = %v, want %v", loaded.Config.Runs.SlotWait, tc.want)
			}
		})
	}
}

func TestProjectConfigCannotSetTheRunSlotWait(t *testing.T) {
	t.Parallel()
	for _, value := range []string{"0", "9h", "-1h", "invalid"} {
		t.Run(value, func(t *testing.T) {
			home, work := t.TempDir(), t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustMkdir(t, filepath.Join(work, ".git"))
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "runs:\n  slot_wait: 90m\n")
			mustWrite(t, filepath.Join(work, ".roundfixrc.yml"), "runs:\n  slot_wait: "+value+"\n")
			var stderr bytes.Buffer
			loaded, err := Load(LoadOptions{HomeDir: home, WorkDir: work, Stderr: &stderr})
			if err != nil {
				t.Fatal(err)
			}
			if loaded.Config.Runs.SlotWait != 90*time.Minute {
				t.Fatalf("SlotWait = %v", loaded.Config.Runs.SlotWait)
			}
			want := "config: runs.slot_wait in Project Config is ignored; set runs.slot_wait in User Config\n"
			if stderr.String() != want {
				t.Fatalf("warning = %q, want %q", stderr.String(), want)
			}
		})
	}
}

func TestRunSlotWaitRefusesANegativeOrInvalidValue(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ value, want string }{{"-1s", "runs.slot_wait must be greater than or equal to 0"}, {"invalid", "invalid duration"}} {
		t.Run(tc.value, func(t *testing.T) {
			home := t.TempDir()
			mustMkdir(t, filepath.Join(home, ".roundfix"))
			mustWrite(t, filepath.Join(home, ".roundfix", "config.yml"), "runs:\n  slot_wait: "+tc.value+"\n")
			_, err := Load(LoadOptions{HomeDir: home, WorkDir: t.TempDir()})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
		})
	}
}
