// Boundary: injected filesystem capacity and temporary paths; no host disk measurement.
package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	roundconfig "roundfix/internal/config"
)

func TestCommandFreeSpaceNamesEachLocationBelowTheFloor(t *testing.T) {
	t.Parallel()
	const mib int64 = 1024 * 1024
	for _, tc := range []struct {
		name       string
		available  int64
		unreadable bool
	}{
		{name: "below the floor", available: 10239 * mib},
		{name: "at the floor", available: 10240 * mib},
		{name: "unreadable", unreadable: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			home := t.TempDir()
			worktrees := filepath.Join(home, "worktrees", "not-created")
			cache := filepath.Join(home, "gocache", "not-created")
			loaded := roundconfig.Loaded{Config: roundconfig.Builtin(), GitRoot: t.TempDir(), HomeDir: home}
			loaded.Config.Worktree.Location = "~/worktrees/not-created"
			loaded.Config.Worktree.SharedCaches = []roundconfig.SharedCache{{Path: "{home}/gocache/not-created"}}
			var reads []string
			ctx := contextWithCommandDependencies(t.Context(), commandDependencies{
				availableBytes: func(path string) (int64, error) {
					reads = append(reads, path)
					if tc.unreadable {
						return 0, errors.New("capacity denied")
					}
					return tc.available, nil
				},
			})
			checker := commandDeliveryEngineDependencies(nil, loaded, io.Discard).FreeSpace
			if checker == nil {
				t.Fatal("command engine has no free-space checker")
			}
			got := checker.LowFreeSpace(ctx)
			var want []string
			if !tc.unreadable && tc.available < 10240*mib {
				want = []string{
					fmt.Sprintf("worktree location %s has 10239 MiB free, below 10240 MiB", worktrees),
					fmt.Sprintf("Shared Build Cache %s has 10239 MiB free, below 10240 MiB", cache),
					fmt.Sprintf("temporary directory %s has 10239 MiB free, below 10240 MiB", os.TempDir()),
				}
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("locations=%q want=%q", got, want)
			}
			if wantReads := []string{home, home, os.TempDir()}; !reflect.DeepEqual(reads, wantReads) {
				t.Fatalf("capacity reads=%q want=%q", reads, wantReads)
			}
		})
	}
}
