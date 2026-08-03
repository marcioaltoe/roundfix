---
task: task_06
spec: 0072-qa-is-a-task-not-a-flag
status: completed
type: docs
complexity: low
---

# Task 06: Align the agent guides with the authored gate

## Overview

`docs/agents/autonomous-work.md` still teaches the flag: it shows
`roundfix implement --spec <slug> --qa --detach`, tells the Supervisor when
to pass `--qa`, and reasons about re-requesting it per cycle. All of that
lives in the repository-owned section outside the setup markers (the
baseline assets carry no `--qa` reference, verified during design). The
guidance moves to the authored-gate contract: the gate is declared at
decomposition, the graph runs it, and a gate that reported on a graph that
later grew is invalidated by load-time validation.

## Requirements

1. MUST replace every `--qa` reference in `docs/agents/autonomous-work.md`
   with the authored-gate contract, keeping the section's existing rules
   about corrective-Task caps and serial-chain warnings intact.
2. MUST NOT edit inside setup-owned marker blocks; the changes live in the
   repository-owned section only.
3. MUST leave no `--qa` reference anywhere under `docs/agents/`.
4. MUST keep the guidance consistent with what task_03 shipped: the
   parameter does not exist, and an unsettled authored gate makes an
   all-completed graph runnable.

## Subtasks

- [ ] Rewrite the Supervisor guidance around the authored gate.
- [ ] Sweep `docs/agents/` for residual flag references.

## Acceptance Criteria

- [ ] `docs/agents/autonomous-work.md` describes declaring the gate at
      decomposition and never mentions the flag.
- [ ] No file under `docs/agents/` contains `--qa`.
- [ ] Setup-owned marker blocks are byte-identical.
- [ ] `git status --porcelain` shows no path outside `docs/agents/` and
      this task file.

## Verification

- `grep -rL -- "--qa" docs/agents/*.md | grep -q autonomous-work` —
  expected: exit 0; the guide no longer references the flag.
- `matches="$(grep -rn -- "--qa" docs/agents/)"; grep_status=$?; if [ "$grep_status" -eq 0 ]; then printf '%s\n' "$matches"; exit 1; fi; [ "$grep_status" -eq 1 ]`
  — expected: exit 0 when nothing under the guides references the flag;
  matches or a scan error exit nonzero, with matches printed as diagnostics.
- `go build -buildvcs=false ./...` — expected: exit 0.

## References

- `_prd.md` → Goals (declared once, not per invocation).
- `_techspec.md` → Build Order 6.

## Result

### Implementation

- Replaced the Supervisor's per-invocation QA guidance with the authored-gate
  contract: `write-tasks` declares one terminal `qa` Task covering every
  non-QA leaf, or records `qa: declined` with a non-empty reason.
- Updated profile routing, runtime ownership, and the end-to-end Implement
  invocation so the Task Graph is the sole gate input. The guide states that
  an unsettled authored gate keeps an otherwise all-completed graph runnable.
- Preserved the serial full-gate-cycle warning and the cap of two corrective
  Tasks, while documenting load-time invalidation when a reported gate's
  graph later grows.

### Focused checks

- Red signal: `rtk grep -n -- "--qa" docs/agents/autonomous-work.md`
  exited 0 before the edit and found six stale references at lines 64, 68,
  131, 134, 145, and 156.
- `rtk rg -n --fixed-strings -- '--qa' docs/agents` exited 1 after the edit
  with no matches, the expected absence result.
- The setup-owned `guide.autonomous-work` block from `HEAD` and the working
  tree both hashed to
  `641ad7e3478fcbc0d757372d47db709a391152e839e6073cbaa3843a1b013fa0`.
- `rtk git -c core.fsmonitor=false diff --check` exited 0.
- `rtk git -c core.fsmonitor=false status --porcelain` exited 0 and listed
  only `docs/agents/autonomous-work.md` and this Task file.
- The commands under `## Verification` were not run; Daemon Verification owns
  them.

### Acceptance criteria evidence

1. The `Author the QA gate once per Spec` section declares the gate during
   decomposition, records the decline alternative, and states that an
   unsettled authored gate remains runnable after implementation Tasks settle;
   the focused absence scan found no flag reference in the guide.
2. The focused `rtk rg` scan covered all of `docs/agents/` and returned no
   occurrence of the removed parameter.
3. The extracted setup-owned block has the same SHA-256 digest in `HEAD` and
   the working tree; all guide edits are in the repository-owned section.
4. The focused status check lists no path outside `docs/agents/` and this Task
   file.

### Verification feedback repair — attempt 1

- Inspected the Daemon diagnostic artifact; it contained no output because the
  absence scan found no matches, but the command returned exit 1.
- Traced the contract through `parseVerificationCommands` and `ExecVerifier`.
  The parser extracts only the backticked command, and the executor treats any
  nonzero exit other than the typed temporary exit 75 as Verification failure;
  the prose-only expected exit code cannot change settlement.
- Replaced the negative scan with a shell assertion that returns 0 only for
  grep's no-match status, prints forbidden matches before failing, and also
  fails on grep errors.
- A parse-only `rtk sh -n -c '<updated absence assertion>'` check exited 0
  without executing the declared Verification command.
- The updated command under `## Verification` was not run; the Daemon owns the
  single configured rerun after this repair turn.
