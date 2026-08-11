---
task: task_03
spec: 0072-qa-is-a-task-not-a-flag
status: completed
type: backend
complexity: medium
---

# Task 03: Delete the QA parameter from the Implement Command

## Overview

`--qa` stops existing. What runs is what the graph says. The flag is removed
from the Implement Command's flag set and help text; passing it becomes an
unknown-flag error whose remediation names the authored-gate contract. The
"all Tasks completed" no-op rule inverts its special case: a graph whose
gate node is still unsettled is not a no-op — the Run starts and executes
the gate, which is the behavior the flag used to force. Agent selection
profile categories derive the `qa` category from the graph, not from the
request.

## Requirements

1. MUST remove the `qa` flag from the Implement Command's flag definitions,
   its help text, and the interactive-input QA choice merge; the request
   struct loses its QA field.
2. MUST make `--qa` an unknown-flag failure whose message tells the author
   the gate is declared in the Spec's Task Graph.
3. MUST treat an all-completed graph with an unsettled gate node as
   runnable: the Run starts and the cycle ends with the gate.
4. MUST treat an all-completed graph with a settled gate, a declined
   declaration, or a legacy shape as the no-op it is today.
5. MUST derive `implementProfileCategories`' `qa` category from the graph's
   gate declaration.
6. MUST update the CLI tests that pass `--qa` to graph-driven setups,
   keeping their assertions; add the unknown-flag remediation test.

## Subtasks

- [ ] Remove the flag, the request field, and the interactive merge path.
- [ ] Add the unknown-flag remediation message and its test.
- [ ] Invert the no-op rule around the unsettled gate node.
- [ ] Derive the profile category from the graph.
- [ ] Migrate the flag-driven CLI tests.

## Acceptance Criteria

- [ ] `roundfix implement --spec <slug> --qa` fails as an unknown flag and
      the message names the Task Graph as where the gate lives.
- [ ] The help text contains no QA parameter.
- [ ] An all-completed graph with an unsettled gate starts a Run that ends
      in the gate; with a settled gate or no gate it is a no-op.
- [ ] The `qa` selection profile category resolves for a graph with a gate
      and does not resolve for one without.
- [ ] `git status --porcelain` shows no path outside `internal/cli/` and
      this task file.

## Verification

- `go test ./internal/cli -count=1 -run 'Implement.*QA|QA.*Implement|UnknownFlag' -v | grep -q -- "--- PASS"`
  — expected: exit 0.
- `help="$(go run -buildvcs=false ./cmd/roundfix implement --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; case "$help" in *--qa*) printf '%s\n' "$help"; exit 1;; esac`
  — expected: exit 0; help generation must succeed and the parameter must be
  absent. A present parameter or a failed build exits nonzero with diagnostics.
- `go test ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Feature 3; Success Metric 2.
- `_techspec.md` → API Contracts; Integration Points (`internal/cli`).

## Result

### Implementation

- Removed the Implement Command's `--qa` definition, help entry,
  `commandRequest` field, Interactive Input value plumbing, and daemon-plan
  request value. The loaded Task Graph is now the only gate input.
- Added removed-flag remediation that identifies `--qa` as unknown and tells
  the author to declare the terminal QA Task in the Spec's Task Graph.
- Simplified the all-completed no-op rule: an unsettled authored gate remains
  runnable because it is a non-completed Task, while settled, declined, and
  legacy graphs have no runnable Task and remain no-ops.
- Changed `implementProfileCategories` to append `qa` only when the loaded
  graph declares `QATaskID`; QA Task rows no longer enter the ordinary Task
  Type category scan.
- Migrated the former flag-driven CLI fixtures to authored terminal QA Tasks,
  preserving the verdict, report, event, prompt, external Spec Root,
  auto-push, attach, and profile-fallback assertions.

### Focused checks

- Red signal: the first repository-local-cache run of the new removed-flag and
  profile tests exited 1 because `implementProfileCategories` still required
  the request-driven `includeQA` argument.
- `GOCACHE=<worktree>/.gocache go test ./internal/cli -run
  '^(TestRunImplementHelpListsExactlyImplementedFlags|TestRunImplementRemovedQAFlagExplainsTaskGraph|TestRunImplementInteractiveInputDoesNotChooseQAGate|TestImplementProfileCategoriesDerivesQAFromGraphDeclaration|TestRunImplementAllTasksCompletedReportsWithoutRun|TestRunImplementSettledQAGateReportsWithoutRun|TestRunImplementDeclinedQAGateReportsWithoutRun|TestRunImplementUsesConfiguredExternalSpecRootEndToEnd|TestRunImplementAutoPushOutcomeMatrix|TestRunImplementQAVerdictMatrix|TestRunImplementQAStepSkippedWhenAnyTaskFails|TestRunImplementQAOnlyRunSettlesOutcomeFromVerdict|TestRunImplementQAPromptStatesSpecTargetBranchFromRunRecord|TestAttachReplaysCompletedSpecRunReadOnly)$'
  -count=1` exited 0 with 30 passing tests.
- `GOCACHE=<worktree>/.gocache go test ./internal/cli -run
  '^TestAgentSelectionProfilesMacro$/mixed_profiles_configure_validate_fallback_persist_and_stream$'
  -count=1` exited 0 with both the parent and macro subtest passing against
  the built binary.
- `rtk git diff --check` exited 0.
- The commands under `## Verification` were not run; Daemon Verification owns
  them.

### Acceptance criteria evidence

1. `TestRunImplementRemovedQAFlagExplainsTaskGraph` passes with preflight exit
   2, no stdout, `unknown flag "--qa"`, and remediation naming the Spec's Task
   Graph.
2. `TestRunImplementHelpListsExactlyImplementedFlags` passes with the exact
   implemented flag vocabulary and no `--qa` entry.
3. `TestRunImplementQAOnlyRunSettlesOutcomeFromVerdict` starts a Run from
   completed implementation Tasks plus an unsettled authored gate and proves
   the cycle ends with that gate. The legacy, settled-gate, and declined-gate
   no-op tests pass without creating a Run.
4. `TestImplementProfileCategoriesDerivesQAFromGraphDeclaration` passes for
   both sides: a declared gate resolves `backend, qa`, while a legacy graph
   resolves only `backend`. The binary macro also proves the graph-selected QA
   profile participates in operational preflight and execution.
5. `rtk git -c core.fsmonitor=false status --porcelain` exited 0 and listed
   only `internal/cli/` paths and this Task file.

### Verification feedback repair — attempt 1

- Inspected the Daemon diagnostic artifact and traced command execution through
  `parseVerificationCommands` and `ExecVerifier`. The CLI help behavior was
  correct: the absence check produced no diagnostic text and exited 1, while
  the Daemon contract requires every successful Verification command to exit
  0. The prose-only expected exit code is not part of the parsed command.
- Replaced the inverted `grep` contract with a shell assertion that first
  requires `go run` to succeed, then exits nonzero only when help contains
  `--qa`. It prints captured help on either failure path for actionable Daemon
  diagnostics.
- Focused non-Verification check:
  `GOCACHE=<worktree>/.gocache go test ./internal/cli -run
  '^TestRunImplementHelpListsExactlyImplementedFlags$' -count=1` exited 0
  with one passing test.
- A parse-only `sh -n -c '<updated help assertion>'` check exited 0 without
  executing the declared Verification command.
- The updated command under `## Verification` was not run; the Daemon owns the
  one configured rerun after this repair turn.
