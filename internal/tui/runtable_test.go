package tui

import (
	"strings"
	"testing"
	"time"

	"roundfix/internal/store"
)

func TestShortRunIDIsTheFirstEightSuffixDigits(t *testing.T) {
	t.Parallel()
	for _, test := range []struct{ name, id, want string }{
		{"random suffix", "run_20261105T094412Z_93cee2a36ff33269", "93cee2a3"},
		{"short suffix", "run_x_abc", "run_x_abc"},
		{"legacy", "legacy", "legacy"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := ShortRunID(test.id); got != test.want {
				t.Fatalf("ShortRunID(%q)=%q, want %q", test.id, got, test.want)
			}
		})
	}
}

func TestRenderRunTableKeepsMissingFieldsAndLocalDates(t *testing.T) {
	t.Parallel()
	location := time.FixedZone("-03", -3*60*60)
	now := time.Date(2026, 7, 7, 1, 0, 0, 0, time.UTC)
	runs := []store.Run{
		{ID: "legacy", Kind: store.KindImplement, State: store.StateClean, CreatedAt: now.Add(-2 * time.Hour)},
		{ID: "run_x_123456789abcdef0", Kind: store.KindResolve, State: store.StateActive,
			RepositoryRoot: "/repos/fiscus", GitRoot: "/worktrees/task", PRNumber: "123",
			CreatedAt: now.Add(time.Minute), OwnerIdentityUnproven: true},
	}
	got := RenderRunTable(runs, RunTableOptions{Now: now, Location: location, Repo: true})
	lines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")
	if len(lines) != 3 {
		t.Fatalf("table=%q", got)
	}
	if strings.Join(strings.Fields(lines[0]), " ") != "RUN REPO STATE TARGET AGENT STARTED (-03) DURATION NOTE" {
		t.Fatalf("header=%q", lines[0])
	}
	if strings.Join(strings.Fields(lines[1]), " ") != "legacy - Clean - - 20:00 -" {
		t.Fatalf("missing fields or same local day=%q", lines[1])
	}
	if strings.Join(strings.Fields(lines[2]), " ") != "12345678 fiscus Active pr:123 - 22:01 0s owner identity unproven" {
		t.Fatalf("repository, skew or note=%q", lines[2])
	}
	for _, line := range lines {
		if strings.HasSuffix(line, " ") {
			t.Errorf("trailing spaces=%q", line)
		}
	}
	priorYear := now.AddDate(-1, 0, 0)
	runs[0].CreatedAt = priorYear
	got = RenderRunTable(runs[:1], RunTableOptions{Now: now, Location: location, Wide: true})
	if !strings.Contains(got, "Jul 06 22:00") || strings.Contains(got, "NOTE") {
		t.Fatalf("another year or absent note=%q", got)
	}
	if got := RenderRunTable(nil, RunTableOptions{}); got != "No Runs found.\n" {
		t.Fatalf("empty table=%q", got)
	}
}

func TestRenderRunTableCutsOnlyTheTargetToItsMinimumWidth(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 7, 6, 12, 0, 0, 0, time.UTC)
	got := RenderRunTable([]store.Run{{ID: "run_x_123456789abcdef0", Kind: store.KindImplement,
		State: store.StateActive, SpecSlug: "0123-long-target-for-a-narrow-terminal", Agent: "long-agent-name",
		CreatedAt: now}}, RunTableOptions{Now: now, Location: time.UTC, Width: 40})
	lines := strings.Split(strings.TrimSuffix(got, "\n"), "\n")
	if len(lines) != 2 || strings.Join(strings.Fields(lines[1]), " ") != "12345678 Active 0123-long-t… long-agent-name 12:00 0s" {
		t.Fatalf("minimum TARGET width must preserve other columns: %q", got)
	}
}
