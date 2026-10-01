---
task: task_04
spec: 0211-a-delivery-queue-that-finishes-without-intervention
status: pending
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
