package spec

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// Sequential: testing.AllocsPerRun sets runtime.GOMAXPROCS to 1 and counts process-wide allocations.
func TestHistoryCitationsInSourceWorkDoesNotGrowWithTheRemovedPaths(t *testing.T) {
	targets := make([]string, 2000)
	for i := range targets {
		targets[i] = fmt.Sprintf("docs/history/backlog/2026-01-01-a-removed-backlog-entry-%04d.md", i)
	}
	tokens := make([]string, 200, 201)
	for i := range tokens {
		tokens[i] = fmt.Sprintf("internal/spec/source-file-%03d.go", i)
	}
	last := targets[len(targets)-1]
	content := []byte(strings.Join(append(tokens, last), " "))
	measure := func(removed []string) float64 {
		return testing.AllocsPerRun(5, func() {
			citations, err := HistoryCitationsInSource(content, "Makefile", removed)
			if err != nil {
				t.Fatal(err)
			}
			if len(citations) != 1 {
				t.Fatalf("citations = %+v, want exactly one", citations)
			}
			if citations[0] != (HistoryCitation{Path: "Makefile", Line: 1, Target: last}) {
				t.Fatalf("citation = %+v, want last removed path on Makefile line 1", citations[0])
			}
		})
	}
	one := measure(targets[len(targets)-1:])
	all := measure(targets)
	t.Logf("allocations per call: one removed path = %.0f; 2000 removed paths = %.0f", one, all)
	if all > one+64 {
		t.Fatalf("2000 removed paths allocate %.0f more per call; maximum 64", all-one)
	}
}

func TestHistoryCitationsInSourceNamesEachRemovedPathAndItsDescendants(t *testing.T) {
	t.Parallel()
	targets := []string{
		"docs/history/adr/0001-old.md",
		"docs/history/specs/0001-old",
		"docs/history/specs/0001-old/references/source.md",
		"docs/history/backlog",
	}
	for _, tc := range []struct {
		name    string
		path    string
		token   string
		targets []string
	}{
		{name: "exact path", path: "Makefile", token: targets[0], targets: targets[:1]},
		{name: "descendant", path: "Makefile", token: targets[1] + "/_prd.md", targets: targets[1:2]},
		{name: "nested removed paths", path: "Makefile", token: targets[2], targets: targets[1:3]},
		{name: "relative prose path", path: "docs/agents/guide.txt", token: "../history/adr/0001-old.md", targets: targets[:1]},
		{name: "explicit dot prose path", path: "Makefile", token: "./" + targets[0], targets: targets[:1]},
		{name: "trailing slash", path: "Makefile", token: targets[1] + "/", targets: targets[1:2]},
		{name: "file prefix only", path: "Makefile", token: targets[0] + ".bak"},
		{name: "directory prefix only", path: "Makefile", token: "docs/history/backlog-old/x.md"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := HistoryCitationsInSource([]byte("first line\n"+tc.token), tc.path, targets)
			if err != nil {
				t.Fatal(err)
			}
			var want []HistoryCitation
			for _, target := range tc.targets {
				want = append(want, HistoryCitation{Path: tc.path, Line: 2, Target: target})
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("citations = %+v, want %+v", got, want)
			}
		})
	}
}
