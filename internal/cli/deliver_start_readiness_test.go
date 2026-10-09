// Suite: Delivery publish readiness through public CLI and the real queue store.
// Boundary OUT: forge reads and detached owner launch are injected; no network.
package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	roundconfig "roundfix/internal/config"
	"roundfix/internal/store"
)

func readyDeliveryReadiness(context.Context, roundconfig.Loaded) []CheckResult {
	return []CheckResult{{Name: HealthCheckGH, Status: CheckStatusOK}, {Name: HealthCheckRemote, Status: CheckStatusOK}}
}

func TestDeliverStartRefusesWhenThisMachineCannotPublish(t *testing.T) {
	t.Parallel()
	for _, multiple := range []bool{false, true} {
		t.Run(map[bool]string{false: "transcript", true: "multiple failures"}[multiple], func(t *testing.T) {
			home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
			setImplementFixtureAuthorizationOperations(t, repo, "implement", "commit", "push", "pull_request", "merge")
			started, checked := false, false
			updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
				deps.deliveryReadiness = func(_ context.Context, loaded roundconfig.Loaded) []CheckResult {
					checked = true
					if loaded.GitRoot != repo {
						t.Fatalf("repo=%q", loaded.GitRoot)
					}
					results := []CheckResult{readinessResult("gh", "gh version", []readinessFinding{{"DR-GH-UNAUTHENTICATED", CheckStatusFailed, "gh has no account for github.com", "gh auth login --hostname github.com"}})}
					if multiple {
						results = append(results, readinessResult("remote", "origin", []readinessFinding{{"DR-REMOTE-MISSING", CheckStatusFailed, "origin is missing", "add origin"}}))
					}
					return results
				}
				deps.startDeliveryOwner = func(context.Context, roundconfig.Loaded, commandEnvironment, io.Writer, io.Writer) int {
					started = true
					return exitOK
				}
			})
			var out, diag bytes.Buffer
			code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &out, &diag)
			reason := "  gh: DR-GH-UNAUTHENTICATED: gh has no account for github.com; next: gh auth login --hostname github.com\n"
			if multiple {
				reason += "  remote: DR-REMOTE-MISSING: origin is missing; next: add origin\n"
			}
			want := "Preflight failed\n\nReason:\n  Delivery Queue cannot publish from this machine:\n" + reason + "\nNo side effects:\n  Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.\n\nUsage:\n  Run 'roundfix deliver start --help' for usage.\n"
			if code != exitPreflight || out.Len() != 0 || diag.String() != want || started || !checked {
				t.Fatalf("exit=%d stdout=%q stderr=%q started=%v checked=%v", code, out.String(), diag.String(), started, checked)
			}
			// Even the database must not have been opened by a refused start.
			if _, err := os.Stat(store.DatabasePath(home)); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("database exists before inspection: %v", err)
			}
			db, err := store.Open(t.Context(), home)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			if _, found, err := db.DeliveryQueue(t.Context(), repo); err != nil || found {
				t.Fatalf("queue found=%v err=%v", found, err)
			}
		})
	}
}

func TestDeliverStartProceedsOnAForgeWarning(t *testing.T) {
	t.Parallel()
	home, repo := newImplementWorkspace(t, []implementSeed{{id: "task_01"}})
	setImplementFixtureAuthorizationOperations(t, repo, "implement", "commit", "push", "pull_request", "merge")
	started := false
	updateCommandDependenciesForTest(t, func(deps *commandDependencies) {
		deps.deliveryReadiness = func(context.Context, roundconfig.Loaded) []CheckResult {
			return []CheckResult{readinessResult("gh", "", []readinessFinding{{"DR-GH-UNREACHABLE", CheckStatusWarn, "could not confirm the github.com login", "re-run roundfix doctor when github.com is reachable"}})}
		}
		deps.startDeliveryOwner = func(ctx context.Context, loaded roundconfig.Loaded, _ commandEnvironment, _, _ io.Writer) int {
			db, err := store.Open(ctx, home)
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			queue, found, err := db.DeliveryQueue(ctx, loaded.GitRoot)
			if err != nil || !found || len(queue.Items) != 1 || queue.Items[0].SpecSlug != implementTestSlug {
				t.Fatalf("queue=%+v found=%v err=%v", queue, found, err)
			}
			started = true
			return exitOK
		}
	})
	var out, diag bytes.Buffer
	code := runCLI(t, []string{"deliver", "start", implementTestSlug}, &out, &diag)
	if code != exitOK || !started || diag.Len() != 0 || strings.Contains(out.String(), "DR-GH-UNREACHABLE") {
		t.Fatalf("exit=%d started=%v stdout=%q stderr=%q", code, started, out.String(), diag.String())
	}
}
