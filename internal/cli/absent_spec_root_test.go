package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/store"
)

func TestAnAbsentSpecRootReadsLikeAnEmptyOne(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		args []string
		exit int
	}{
		{"follow-up", []string{"follow-up", "--format", "json"}, 0},
		{"follow-up --status", []string{"follow-up", "--status", "--format", "json"}, 0},
		{"spec check", []string{"spec", "check"}, 0},
		{"deliver plan", []string{"deliver", "plan"}, 0},
		{"deliver status", []string{"deliver", "status"}, 2},
		{"history sanitize", []string{"history", "sanitize"}, 0},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := newFollowUpFixture(t)
			if tt.name == "history sanitize" {
				followUpWrite(t, f.repo, "docs/history/findings/retired.md", "---\nstatus: closed\nreason: done\n---\n\n# Retired\n\nPreserved outcome.\n")
				gitImplement(t, f.repo, "add", ".")
				gitImplement(t, f.repo, "commit", "-m", "docs: seed retired history")
				gitImplement(t, f.repo, "tag", "-a", "history-full", "-m", "Full history")
			}
			db, err := store.Open(t.Context(), f.home)
			if err != nil {
				t.Fatal(err)
			}
			if err := db.Close(); err != nil {
				t.Fatal(err)
			}
			root := filepath.Join(f.repo, "docs/specs")
			if err := os.Remove(root); err != nil {
				t.Fatal(err)
			}
			var absentOut, absentErr bytes.Buffer
			absentCode := runCLI(t, tt.args, &absentOut, &absentErr)
			if strings.Contains(absentOut.String()+absentErr.String(), "which does not exist") {
				t.Error("absent Spec Root was refused")
			}
			if strings.HasPrefix(tt.name, "follow-up") {
				var report struct {
					Sources []followUpSource `json:"sources"`
				}
				if err := json.Unmarshal(absentOut.Bytes(), &report); err != nil {
					t.Fatal(err)
				}
				found := false
				for _, source := range report.Sources {
					if source.Name == "specs" {
						found = true
						if source.State != "checked" || source.Count != 0 {
							t.Errorf("specs source = %+v, want checked with count 0", source)
						}
					}
				}
				if !found {
					t.Error("specs source missing")
				}
			}
			if err := os.MkdirAll(root, 0755); err != nil {
				t.Fatal(err)
			}
			var emptyOut, emptyErr bytes.Buffer
			emptyCode := runCLI(t, tt.args, &emptyOut, &emptyErr)
			if emptyCode != tt.exit {
				t.Fatalf("empty root: exit=%d stdout=%s stderr=%s", emptyCode, &emptyOut, &emptyErr)
			}
			if absentCode != emptyCode || absentOut.String() != emptyOut.String() || absentErr.String() != emptyErr.String() {
				t.Errorf("absent: exit=%d stdout=%q stderr=%q; empty: exit=%d stdout=%q stderr=%q", absentCode, &absentOut, &absentErr, emptyCode, &emptyOut, &emptyErr)
			}
		})
	}
}

func TestBaselineUpdateHistoryPlansAnAbsentSpecRootLikeAnEmptyOne(t *testing.T) {
	t.Parallel()
	f := newFollowUpFixture(t)
	root := filepath.Join(f.repo, "docs/specs")
	if err := os.Remove(root); err != nil {
		t.Fatal(err)
	}
	env := commandEnvironmentForTest(t)
	absent := planBaselineUpdateHistory(t.Context(), f.repo, env)
	if absent.Report.Status == "blocked" {
		t.Errorf("absent root blocked: %+v", absent.Report)
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		t.Fatal(err)
	}
	empty := planBaselineUpdateHistory(t.Context(), f.repo, env)
	if empty.Report.Status == "blocked" {
		t.Fatalf("empty root blocked: %+v", empty.Report)
	}
	if !reflect.DeepEqual(absent.Report, empty.Report) || absent.specRootRel != empty.specRootRel || absent.archiveRoot != empty.archiveRoot {
		t.Errorf("absent plan = %+v; empty plan = %+v", absent, empty)
	}
}

func TestANamedSpecInAnAbsentSpecRootKeepsTheRefusal(t *testing.T) {
	t.Parallel()
	for _, command := range []string{"spec check", "deliver plan", "deliver start", "deliver retry"} {
		t.Run(command, func(t *testing.T) {
			f := newFollowUpFixture(t)
			root := filepath.Join(f.repo, "docs/specs")
			if err := os.Remove(root); err != nil {
				t.Fatal(err)
			}
			args := append(strings.Fields(command), "0001-widget")
			var out, stderr bytes.Buffer
			code := runCLI(t, args, &out, &stderr)
			reason := fmt.Errorf("specs.root resolved to %q, which does not exist; create the directory or update specs.root", root)
			var wantErr bytes.Buffer
			if command == "spec check" {
				fmt.Fprintf(&wantErr, "roundfix: spec check failed: %v\nRun 'roundfix spec check --help' for usage.\n", reason)
			} else {
				printPreflightFailure(command, reason, &wantErr)
			}
			if code != 2 || out.Len() != 0 || stderr.String() != wantErr.String() {
				t.Fatalf("exit=%d stdout=%q stderr=%q; want exit=2, no stdout, stderr=%q", code, &out, &stderr, &wantErr)
			}
		})
	}
}
