---
task: task_04
spec: 0211-a-delivery-queue-that-finishes-without-intervention
status: completed
type: backend
complexity: medium
---

# Task 04: The agent environment drops a NODE_OPTIONS preload whose file is missing

## Overview

Every agent process Roundfix starts, `acpx` and the ACP adapters, is a Node
program that inherits the user's `NODE_OPTIONS`. When a preload it names no
longer exists, Node exits before running anything and the profile proof
reports an adapter failure. This Task removes, from the agent environment
only, each preload whose absolute path does not exist, keeps every other
option, and prints one notice per removed path (ADR-0211).

## Requirements

1. MUST add `agentNodeOptions` in the new file `internal/agent/node_options.go`
   with the signature and rules of "The agent environment": word splitting on
   unquoted spaces with double quotes and backslash escapes; the preload forms
   `--require <v>`, `-r <v>`, `--import <v>`, `--require=<v>` and
   `--import=<v>`; dropping only an absolute path or an absolute `file://` URL
   that does not exist; keeping relative paths, package names and existing
   paths; re-quoting a kept word that holds a space; `ok` false for an
   unbalanced quote; and removal of the variable when nothing is kept.
2. MUST apply it in `ACPXRunner.baseEnv` to the last `NODE_OPTIONS` entry,
   for the process environment and an explicit `Environment`, so the `acpx`
   probe, adapter checks, profile proof and sessions all receive the result.
3. MUST add the `Notices io.Writer` field to `ACPXRunner` and write, once per
   process per dropped path, the notice line of API Contract 4 to it, or to
   standard error when it is nil.
4. MUST NOT change any other variable, the `CODEX_PATH` and `CLAUDECODE`
   stripping, or the user's environment, and MUST pass a value it cannot split
   through unchanged.
5. MUST add the tests named in Verification to the new file
   `internal/agent/node_options_test.go`: the splitter's cases of Testing
   Approach 4, and a runner test with an explicit `Environment`, a `Notices`
   buffer and a fake `acpx` script that prints the `NODE_OPTIONS` it
   received, asserting `--max-old-space-size=4096` alone, one notice across
   two probes, and the hygiene variables still stripped.

## Subtasks

- [ ] Add the preload filter.
- [ ] Apply it in the runner's base environment.
- [ ] Write the notice once per dropped path.
- [ ] Add the splitter and runner tests.

## Acceptance Criteria

- [ ] `--require=/missing.cjs --max-old-space-size=4096` yields
      `--max-old-space-size=4096` and one dropped path; the space form, `-r`,
      and `--import` with a `file://` URL behave the same.
- [ ] An existing preload, a relative path and a package name are kept; a
      quoted path with a space round-trips; an unbalanced quote is passed
      through unchanged; a value of only missing preloads removes the
      variable.
- [ ] The fake `acpx` receives only the kept options, the notice appears
      exactly once across two probes, and `CODEX_PATH` and `CLAUDECODE` are
      still stripped.

## Context

- interface: `internal/agent/codex_spawn.go`
- interface: `internal/agent/acpx_runner.go`
- creates: `internal/agent/node_options.go`
- creates: `internal/agent/node_options_test.go`
- instruction: `internal/agent/codex_spawn_test.go`
- instruction: `internal/agent/acpx_runner_test.go`
- instruction: `internal/agent/selection_assignment.go`
- instruction: `docs/adr/0211-the-agent-environment-drops-a-node-preload-whose-file-is-missing.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAgentNodeOptionsDropsAMissingPreload|TestAgentNodeOptionsKeepsWhatItCannotSplit|TestACPXRunnerDropsAMissingPreloadWithOneNotice)$" ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAgentNodeOptionsDropsAMissingPreload TestAgentNodeOptionsKeepsWhatItCannotSplit TestACPXRunnerDropsAMissingPreloadWithOneNotice; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestACPXCommandEnvStripsGuardAndHygieneVariables|TestAgentNodeOptionsDropsAMissingPreload)$" ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestACPXCommandEnvStripsGuardAndHygieneVariables TestAgentNodeOptionsDropsAMissingPreload; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing environment hygiene test runs unedited beside the splitter test, which does not exist before this Task.

## References

- `_prd.md` → User Story 5; Core Feature 6; Success Metric 5; Acceptance evidence
- `_techspec.md` → The agent environment; Interfaces; API Contract 4; Testing Approach 4; Build Order 4
- ADR-0211; ADR-0140

## Result

Implemented the Task 04 environment filter. `agentNodeOptions` removes only
missing absolute preload paths and absolute `file://` URL paths, preserves
option order, and writes kept words with the quotes and escapes needed to
round-trip. An unsplittable value passes through unchanged.

`ACPXRunner.baseEnv` filters the last `NODE_OPTIONS` entry in a copied explicit
environment or a fresh process-environment snapshot. It removes all entries
for that variable when no option remains, so a shadowed entry cannot become
active. An `os.Stat` failure removes a preload only when it proves absence.
Notices use `Notices` or standard error and are serialized and deduplicated
by path across runners in the same Roundfix process. The version probe now
uses the existing `commandEnv` hygiene policy; `acpxCommandEnv` itself and
its existing tests are unchanged.

Acceptance evidence from focused checks:

1. `TestAgentNodeOptionsDropsAMissingPreload` covers equals and space forms,
   `-r`, `--import`, file URLs, percent-encoded paths, quoted missing paths,
   and preservation of other options in order. Each removal returns the
   expected absolute path and kept value.
2. The same test keeps existing paths, existing file URLs, relative paths,
   package names, quoted spaces, escaped quotes and backslashes, and unrelated
   options. `TestAgentNodeOptionsKeepsWhatItCannotSplit` checks unchanged
   malformed input without any existence checks. Runner environment tests
   use a real existing file, check last-entry precedence and complete variable
   removal, and prove that the process environment stays unchanged.
3. `TestACPXRunnerDropsAMissingPreloadWithOneNotice` runs a fake `acpx` script
   twice. Both invocations record only `--max-old-space-size=4096`, stripped
   `CODEX_PATH` and `CLAUDECODE`, and the retained unrelated variable. The
   notice matches API Contract 4 exactly once; a second runner does not repeat
   it, and the supplied environment remains unchanged. The concurrent-notice
   test exercises repeated occurrences of one path across eight callers.

Focused commands and observed outcomes:

- Before implementation:
  `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test -count=1 -run 'NodeOptions|MissingPreload' ./internal/agent`
  exited 1 because `agentNodeOptions` and `ACPXRunner.Notices` did not exist.
- After implementation:
  `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test -count=1 -v -run 'NodeOptions|MissingPreload|ACPXCommandEnv|CommandEnvDefaults' ./internal/agent`
  exited 0; all selected tests and subtests passed, including the unedited
  environment hygiene tests.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test -race -count=1 -run 'NodeOptions|MissingPreload|ACPXCommandEnv|CommandEnvDefaults' ./internal/agent`
  exited 0 with no race report.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test -count=1 -run '^(TestACPXProbe|TestACPXConfigPathUsesOnlyTheExplicitEnvironment)' ./internal/agent`
  exited 0; existing probe and explicit-environment regression checks passed.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The authored Verification commands and repository-wide Verification were not
run in this Agent turn. Task status remains `in_progress` as supplied by the
Daemon. No other Task, Task Graph, tooling, commit, push, or Pull Request was
changed or created. No follow-up work was identified.

### Verification Feedback repair — attempt 1

Inspected the Daemon diagnostic at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261002T075527Z_5bea04bbe00b3ea6/verification/batch-002-attempt-1.log`.
The configured `make verify-changed` rejected the fake `acpx` test fixture's
direct executable write under the repository's frozen executable inventory.
A focused run of the repository inventory subtest reproduced the rejection.

Repaired only `internal/agent/node_options_test.go`: the fake shell script is
now non-executable data with mode `0o600`, invoked by a launcher compiled
through the existing `testfixture.FixtureBinary` helper. The Go tool owns the
executable output. The script still observes both real probe environments,
and all kept-option, notice, hygiene, and environment-preservation assertions
remain in place. The executable inventory and its guard were not changed.

Fresh focused evidence after repair:

- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test -count=1 -run '^TestNoTestWritesAnExecutableOutsideTheResidue$' ./internal/testfixture`
  exited 0; the unchanged inventory guard and its seeded detector cases passed.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test -race -count=1 -run 'NodeOptions|MissingPreload|ACPXCommandEnv|CommandEnvDefaults' ./internal/agent`
  exited 0 with no race report; all acceptance-related focused checks passed
  with the repaired compiled launcher and fake shell script.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.

The configured Verification sequence and authored Verification commands were
not rerun. Task status remains Daemon-owned and unchanged. No commit, push,
Pull Request, other Task, or Task Graph change was made during this repair.

## Carry-forward provenance

- Source Run: `run_20261002T075527Z_5bea04bbe00b3ea6`
- Source commit: `e6d482d5ef1129ef137431d9ab22395fac316afe`
