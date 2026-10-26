package sharedcache_test

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/sharedcache"
)

func fixture(t *testing.T) (string, time.Time, []string) {
	t.Helper()
	dir := t.TempDir()
	now := time.Date(2026, 10, 24, 12, 0, 0, 0, time.UTC)
	names := []string{"aa/" + strings.Repeat("a", 64) + "-d", "aa/" + strings.Repeat("b", 64) + "-a", "bb/" + strings.Repeat("c", 64) + "-d", "bb/" + strings.Repeat("d", 64) + "-d"}
	ages := []time.Duration{10 * time.Hour, 5 * time.Hour, 3 * time.Hour, 30 * time.Minute}
	write := func(name, content string, age time.Duration) {
		t.Helper()
		path := filepath.Join(dir, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, now.Add(-age), now.Add(-age)); err != nil {
			t.Fatal(err)
		}
	}
	for i, name := range names {
		write(name, strings.Repeat("x", 1000), ages[i])
		names[i] = filepath.Join(dir, filepath.FromSlash(name))
	}
	write("README", sharedcache.GoBuildReadme+"\n", 48*time.Hour)
	write("trim.txt", "trim\n", 48*time.Hour)
	write("aa/notes.txt", strings.Repeat("x", 1000), 48*time.Hour)
	return dir, now, names
}

func paths(plan sharedcache.Plan) []string {
	result := make([]string, 0, len(plan.Remove))
	for _, entry := range plan.Remove {
		result = append(result, entry.Path)
	}
	return result
}

func TestMeasurePlansIdleEntriesOldestFirstUntilWithinTheLimit(t *testing.T) {
	t.Parallel()
	dir, now, names := fixture(t)
	plan, err := sharedcache.Measure(dir, 3500, sharedcache.MinimumIdle, now)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(paths(plan), names[:2]) {
		t.Fatalf("planned %v, want %v", paths(plan), names[:2])
	}
	wantBytes := int64(5000 + len(sharedcache.GoBuildReadme) + 1 + len("trim\n"))
	if plan.Bytes != wantBytes || plan.RemoveSize != 2000 || plan.Limit != 3500 || plan.Path != dir {
		t.Fatalf("unexpected plan: %+v", plan)
	}
	for _, name := range names {
		if _, err := os.Stat(name); err != nil {
			t.Fatalf("Measure changed entry: %v", err)
		}
	}
	count, bytes, err := sharedcache.Apply(plan, sharedcache.MinimumIdle, func() time.Time { return now })
	if err != nil || count != 2 || bytes != 2000 {
		t.Fatalf("Apply = %d, %d, %v", count, bytes, err)
	}
	for _, name := range names[:2] {
		if _, err := os.Stat(name); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("removed entry %q: %v", name, err)
		}
	}
	for _, name := range append(names[2:], filepath.Join(dir, "README"), filepath.Join(dir, "trim.txt"), filepath.Join(dir, "aa", "notes.txt")) {
		if _, err := os.Stat(name); err != nil {
			t.Errorf("kept file %q: %v", name, err)
		}
	}
}

func TestMeasureKeepsEveryEntryWithinTheLimitOrInUse(t *testing.T) {
	t.Parallel()
	dir, now, names := fixture(t)
	plan, err := sharedcache.Measure(dir, 1<<20, sharedcache.MinimumIdle, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Remove) != 0 || plan.RemoveSize != 0 {
		t.Fatalf("within limit: %+v", plan)
	}
	count, bytes, err := sharedcache.Apply(plan, sharedcache.MinimumIdle, func() time.Time { return now })
	if err != nil || count != 0 || bytes != 0 {
		t.Fatalf("within limit Apply = %d, %d, %v", count, bytes, err)
	}
	plan, err = sharedcache.Measure(dir, 0, sharedcache.MinimumIdle, now)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(paths(plan), names[:3]) || plan.RemoveSize != 3000 {
		t.Fatalf("zero limit: %+v", plan)
	}
	// The cutoff is strict: exactly two hours old is still protected.
	if err := os.Chtimes(names[2], now.Add(-sharedcache.MinimumIdle), now.Add(-sharedcache.MinimumIdle)); err != nil {
		t.Fatal(err)
	}
	plan, err = sharedcache.Measure(dir, 0, sharedcache.MinimumIdle, now)
	if err != nil || !reflect.DeepEqual(paths(plan), names[:2]) {
		t.Fatalf("cutoff plan = %+v, %v", plan, err)
	}
}

func TestMeasureRefusesWhatIsNotAGoBuildCache(t *testing.T) {
	t.Parallel()
	dir, now, _ := fixture(t)
	link := filepath.Join(t.TempDir(), "cache")
	if err := os.Symlink(dir, link); err != nil {
		t.Fatal(err)
	}
	invalid := t.TempDir()
	if err := os.WriteFile(filepath.Join(invalid, "README"), []byte("unrelated directory"), 0600); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		dir  string
		idle time.Duration
	}{
		{t.TempDir(), sharedcache.MinimumIdle}, {link, sharedcache.MinimumIdle}, {dir, time.Hour}, {invalid, sharedcache.MinimumIdle}, {filepath.Join(dir, "README"), sharedcache.MinimumIdle},
	} {
		if _, err := sharedcache.Measure(tc.dir, 0, tc.idle, now); !errors.Is(err, sharedcache.ErrNotGoBuildCache) {
			t.Errorf("Measure(%q,%s): %v", tc.dir, tc.idle, err)
		}
	}
}

func TestApplySkipsAnEntryUsedSinceItWasPlanned(t *testing.T) {
	t.Parallel()
	dir, now, names := fixture(t)
	plan, err := sharedcache.Measure(dir, 0, sharedcache.MinimumIdle, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Remove) != 3 {
		t.Fatalf("planned %d entries", len(plan.Remove))
	}
	if err := os.Chtimes(plan.Remove[0].Path, now, now); err != nil {
		t.Fatal(err)
	}
	count, bytes, err := sharedcache.Apply(plan, sharedcache.MinimumIdle, func() time.Time { return now })
	if err != nil || count != 2 || bytes != 2000 {
		t.Fatalf("Apply = %d, %d, %v", count, bytes, err)
	}
	for _, name := range []string{names[0], names[3]} {
		if _, err := os.Stat(name); err != nil {
			t.Errorf("protected entry %q: %v", name, err)
		}
	}
	for _, name := range names[1:3] {
		if _, err := os.Stat(name); !errors.Is(err, os.ErrNotExist) {
			t.Errorf("idle entry %q: %v", name, err)
		}
	}
	count, bytes, err = sharedcache.Apply(plan, sharedcache.MinimumIdle, func() time.Time { return now })
	if err != nil || count != 0 || bytes != 0 {
		t.Fatalf("repeat Apply = %d, %d, %v", count, bytes, err)
	}
}
