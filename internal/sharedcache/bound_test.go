package sharedcache_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"roundfix/internal/sharedcache"
)

func TestBoundMainTrimsOnlyIdleEntriesAboveTheBound(t *testing.T) {
	t.Parallel()
	const entryBytes = 600 * 1024
	for _, tc := range []struct {
		name     string
		maxMB    string
		ages     []time.Duration
		removed  int
		noReadme bool
	}{
		{name: "within its bound", maxMB: "2", ages: []time.Duration{10 * time.Hour, 5 * time.Hour, 30 * time.Minute}},
		{name: "above its bound with idle entries", maxMB: "1", ages: []time.Duration{10 * time.Hour, 5 * time.Hour, 30 * time.Minute}, removed: 2},
		{name: "above its bound with no idle entry", maxMB: "1", ages: []time.Duration{90 * time.Minute, 30 * time.Minute}},
		{name: "not a Go build cache", maxMB: "1", ages: []time.Duration{10 * time.Hour, 5 * time.Hour, 30 * time.Minute}, noReadme: true},
		// Input order differs from age order, and trimming stops with idle entries left.
		{name: "oldest first until within bound", maxMB: "2", ages: []time.Duration{5 * time.Hour, 30 * time.Minute, 10 * time.Hour, 3 * time.Hour}, removed: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			now := time.Date(2026, 11, 6, 12, 0, 0, 0, time.UTC)
			readme := sharedcache.GoBuildReadme + "\n"
			if !tc.noReadme {
				if err := os.WriteFile(filepath.Join(dir, "README"), []byte(readme), 0600); err != nil {
					t.Fatal(err)
				}
			}
			var paths []string
			for i, age := range tc.ages {
				path := filepath.Join(dir, "aa", fmt.Sprintf("%064x-d", i))
				if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, bytes.Repeat([]byte("x"), entryBytes), 0600); err != nil {
					t.Fatal(err)
				}
				if err := os.Chtimes(path, now.Add(-age), now.Add(-age)); err != nil {
					t.Fatal(err)
				}
				paths = append(paths, path)
			}
			var stderr bytes.Buffer
			code := sharedcache.BoundMain([]string{"-dir", dir, "-max-mb", tc.maxMB, "-idle", "2h"}, &stderr, func() time.Time { return now })
			if code != 0 {
				t.Fatalf("exit = %d; stderr: %s", code, &stderr)
			}
			want := ""
			switch {
			case tc.noReadme:
				_, err := sharedcache.MeasureContext(t.Context(), dir, 1<<20, 2*time.Hour, now)
				if !errors.Is(err, sharedcache.ErrNotGoBuildCache) {
					t.Fatalf("refusal = %v", err)
				}
				want = fmt.Sprintf("gocache-bound: warning: %v; the build continues\n", err)
			case tc.removed > 0:
				want = fmt.Sprintf("gocache-bound: %s: removed %d idle entries (%d bytes)\n", dir, tc.removed, tc.removed*entryBytes)
			case tc.maxMB == "1":
				want = fmt.Sprintf("gocache-bound: warning: %s holds %d bytes, above its 1048576-byte bound, and no other entry is idle for 2h0m0s; the build continues\n", dir, len(tc.ages)*entryBytes+len(readme))
			}
			if stderr.String() != want {
				t.Fatalf("stderr = %q, want %q", stderr.String(), want)
			}
			for i, path := range paths {
				removed := tc.removed == 2 && i < 2 || tc.removed == 1 && i == 2
				_, err := os.Stat(path)
				if removed {
					if !errors.Is(err, os.ErrNotExist) {
						t.Errorf("entry %d must be removed: %v", i, err)
					}
				} else if err != nil {
					t.Errorf("entry %d must remain: %v", i, err)
				}
			}
			if !tc.noReadme {
				got, err := os.ReadFile(filepath.Join(dir, "README"))
				if err != nil || string(got) != readme {
					t.Fatalf("README changed: %q, %v", got, err)
				}
			}
		})
	}
}

func TestBoundMainRefusesAnInvalidRequest(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name string
		args []string
		want string
	}{
		{"relative dir", []string{"-dir", "relative", "-max-mb", "1", "-idle", "2h"}, "-dir must be an absolute path"},
		{"missing dir", []string{"-max-mb", "1", "-idle", "2h"}, "-dir must be an absolute path"},
		{"zero bound", []string{"-max-mb", "0"}, "-max-mb must be a positive MiB count whose bytes fit an int64"},
		{"negative bound", []string{"-max-mb", "-1"}, "-max-mb must be a positive MiB count whose bytes fit an int64"},
		{"overflowing bound", []string{"-max-mb", "8796093022208"}, "-max-mb must be a positive MiB count whose bytes fit an int64"},
		{"short idle", []string{"-max-mb", "1", "-idle", "1h"}, "-idle must be at least 2h0m0s"},
		{"extra argument", []string{"-max-mb", "1", "extra"}, "unexpected argument \"extra\""},
		{"unknown flag", []string{"-unknown"}, "flag provided but not defined: -unknown"},
		{"malformed bound", []string{"-max-mb", "no"}, "invalid value \"no\" for flag -max-mb: parse error"},
		{"malformed idle", []string{"-idle", "no"}, "invalid value \"no\" for flag -idle: parse error"},
		{"missing flag value", []string{"-idle"}, "flag needs an argument: -idle"},
		{"help is a usage error", []string{"-h"}, "flag: help requested"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			dir := t.TempDir()
			args := append([]string{"-dir", dir}, tc.args...)
			if tc.name == "missing dir" {
				args = tc.args
			}
			var stderr bytes.Buffer
			if code := sharedcache.BoundMain(args, &stderr, func() time.Time { t.Fatal("invalid request reached measurement"); return time.Time{} }); code != 2 {
				t.Fatalf("exit = %d, want 2; stderr: %s", code, &stderr)
			}
			if want := "gocache-bound: " + tc.want + "\n"; stderr.String() != want {
				t.Fatalf("stderr = %q, want %q", stderr.String(), want)
			}
			entries, err := os.ReadDir(dir)
			if err != nil || len(entries) != 0 {
				t.Fatalf("invalid request changed cache: %v, %v", entries, err)
			}
			if strings.Count(stderr.String(), "\n") != 1 {
				t.Fatalf("usage must be one line: %q", stderr.String())
			}
		})
	}
}

func TestBoundMainReportsItsBudgetAndKeepsTheCache(t *testing.T) {
	t.Parallel()
	dir, now, paths := fixture(t)
	synctest.Test(t, func(t *testing.T) {
		var stderr bytes.Buffer
		code := sharedcache.BoundMain([]string{"-dir", dir, "-max-mb", "1", "-idle", "2h"}, &stderr, func() time.Time {
			time.Sleep(sharedcache.BoundBudget + time.Nanosecond)
			return now
		})
		want := fmt.Sprintf("gocache-bound: warning: %s: stopped after its 30s budget; the next make continues\n", dir)
		if code != 0 || stderr.String() != want {
			t.Fatalf("exit = %d; stderr = %q, want %q", code, stderr.String(), want)
		}
	})
	for _, path := range paths {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("budget changed cache: %v", err)
		}
	}
}
