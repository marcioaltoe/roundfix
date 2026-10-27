---
task: task_02
spec: 0330-settle-follows-the-branch-plan-and-declared-rows-match-their-criteria
status: completed
type: backend
complexity: medium
---

# Task 02: Settle runs from the user checkout

## Overview

`roundfix settle` run inside a kept Run Worktree takes that worktree for the
user checkout: it commits there and reports `integration pending`. This Task
makes settle refuse, before it writes anything, when its working directory's
branch is a Run Branch or Task Branch, and name the user checkout the Run
recorded; and writes the rule into the settle usage, the Roundfix `settle`
and `spec-delivery` references, the `settle` guide and `CONTEXT.md`.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-settle-run-inside-the-run-worktree-takes-it-for-the-user-checkout.md`](references/2026-10-26-settle-run-inside-the-run-worktree-takes-it-for-the-user-checkout.md)
   of 2026-10-26 by implementing API Contracts 5 and 9 of `_techspec.md`.
2. MUST change `preflightSettle` in `internal/cli/settle.go` to refuse right
   after it inspects Git when the branch starts with `store.RunBranchPrefix`,
   with exit `2` and the message of API Contract 5, naming the `GitRoot` of
   the Run whose ID the branch carries (read with `store.OpenReader`, the ID
   cut before any `-task_` suffix) or, without that record, the first of
   `roundconfig.RepositoryRoots` for the working directory's root.
3. MUST add `internal/cli/settle_checkout_test.go`, through CLI dispatch on
   the implement workspace fixture with `createImplementRunWorktreeFixture`
   and `setCommandWorkDirForTest`, with:
   - `TestSettleRefusesInsideTheRunWorktree`: a failed Task whose work is in
     the kept Run Worktree, run from that worktree; exit `2`, stderr names
     `settle runs from the user checkout` and the user checkout's path, and
     the worktree's HEAD, status and Task file are unchanged;
   - `TestSettleRefusesInsideATaskWorktree`: the same from a kept Task
     Worktree, naming the same checkout;
   - `TestSettleFromTheUserCheckoutStillSettlesTheRunWorktree`: the same
     Run settled from the user checkout exits `0` with `Run Worktree
     removed:` on stderr.
4. MUST add API Contract 9's sentence, verbatim, to `settleUsage`, the
   `## Settle Command` section of `.agents/skills/roundfix/references/settle.md`,
   the settle recovery step of
   `.agents/skills/roundfix/references/spec-delivery.md` and
   `docs/user-guide/commands/settle.md`, each with the refusal message; then
   run `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the `roundfix` version pair.
5. MUST revise **Settle Command** in `CONTEXT.md` through `domain-modeling`
   with "It runs from the user checkout.", keeping its `_Avoid_` line.
6. MUST, as the last step after the last guide edit, re-record the Behavior
   Surface fingerprint with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Hand-edit no mirror, version field or record.
7. MUST NOT change settle's surface order, its integration, its exit codes
   for any other case, the Daemon, any text inside a `### QA settlement`
   section, the Makefile or the CI workflows. Each new top-level test calls
   `t.Parallel()` first.

## Subtasks

- [ ] Refuse inside a Run or Task Worktree and name the user checkout.
- [ ] Add the three CLI tests.
- [ ] Write the rule into the usage, references, guide and glossary; sync and record; re-record the Behavior Surface last.

## Acceptance Criteria

- [ ] Settle inside a Run Worktree or Task Worktree exits 2, writes nothing and names the user checkout.
- [ ] Settle from the user checkout is unchanged.
- [ ] The usage, references, guide and glossary carry the rule, and the owned-skill and Behavior Surface records are current.

## Consumers

- Settle preflight — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0343-settle-follows-the-branch-plan-and-a-declared-row-matches-the-criterion-it-names.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/cli/settle_test.go`
- interface: `internal/cli/settle.go`
- creates: `internal/cli/settle_checkout_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/references/spec-delivery.md`
- interface: `skills/roundfix/references/spec-delivery.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/settle.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestSettleRefusesInsideTheRunWorktree|TestSettleRefusesInsideATaskWorktree|TestSettleFromTheUserCheckoutStillSettlesTheRunWorktree)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSettleRefusesInsideTheRunWorktree TestSettleRefusesInsideATaskWorktree TestSettleFromTheUserCheckoutStillSettlesTheRunWorktree; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < internal/cli/settle.go | grep -qF -- 'Settle runs from the user checkout: inside a Run Worktree or Task Worktree it refuses and names the checkout.' || { echo 'settle usage lacks the checkout rule' >&2; exit 1; }; go test -count=1 -timeout 15m -run 'Settle' ./internal/cli` — expected: exit 0. Before this Task the three tests do not exist and the usage does not state the rule, so the first pass check fails. Afterwards settle refuses inside a Run or Task Worktree and every settle test of the package passes.
- `for file in .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md .agents/skills/roundfix/references/spec-delivery.md skills/roundfix/references/spec-delivery.md docs/user-guide/commands/settle.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'Settle runs from the user checkout: inside a Run Worktree or Task Worktree it refuses and names the checkout.' || { printf '%s lacks the checkout rule\n' "$file" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'It runs from the user checkout.' || { echo 'CONTEXT.md lacks the revised **Settle Command**' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Settle Command**:' || { echo 'CONTEXT.md lacks **Settle Command**' >&2; exit 1; }; diff -r .agents/skills/roundfix skills/roundfix && go test -timeout 15m -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills && go test -timeout 15m -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no reference or guide states the rule, so the first phrase check fails. Afterwards the skill and its mirror agree, the versions are recorded and the documentation contracts pass with the Behavior Surface record re-recorded last.

## References

- `_prd.md` → Goals 4 and 6; Core Features 4 and 6; Success Metrics 4 and 7; Glossary
- `_techspec.md` → Current behavior 5; API Contracts 5 and 9; Declared breaks 3 and 6; Invariant 3; Testing Approach; Build Order 2
- ADR-0343


## Result

Implemented the checkout guard immediately after Git inspection. A Run Branch
or Task Branch refuses with exit 2 and the API Contract 5 recovery message.
The branch's Run ID is cut before `-task_`; a read-only Run lookup supplies
its recorded GitRoot, with RepositoryRoots' main worktree as the fallback
when the database or record is absent. Surface selection and integration
remain downstream of this guard.

Acceptance evidence:

- Run/Task Worktree refusal: the CLI tests assert exit 2, the recovery message
  including the user checkout, empty stdout, unchanged HEAD, porcelain status
  and Task file, and preservation of the Run Worktree. Additional subtests
  prove a recorded linked checkout takes precedence over the main worktree,
  an unknown Run falls back, and a missing database stays absent.
- User checkout behavior: the CLI test settles the same failed Task from the
  user checkout with exit 0, reports `Run Worktree removed:`, verifies removal
  and reads the recovered file from the user checkout.
- Written rule: added API Contract 9's exact sentence and refusal to usage,
  both canonical references and the settle guide. Revised Settle Command
  through domain-modeling with `It runs from the user checkout.` while
  retaining its Avoid line. Generated mirrors match their canonical sources;
  the version recorder raised both roundfix versions to 0.1.153. The Behavior
  Surface fingerprint was regenerated after the final guide edit.

Focused checks and regeneration:

- Red check: `go test ./internal/cli -run '^TestSettleRefusesInsideTheRunWorktree$' -count=1 -timeout 3m`
  reproduced exit 1 and `integration pending` before the implementation.
- `go test ./internal/cli -run '^(TestSettleRefusesInside.*|TestSettleFromTheUserCheckout.*|TestSettleCheckoutRecordAndFallback)$' -count=1 -timeout 3m`
  exited 0 after the final code edit, covering the three requested tests and
  all record/fallback subtests.
- `make skills-sync` exited 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and wrote the version pair and record through the sanctioned recorder.
- `make baseline-digests` exited 0 with `changed:false`; no derived Baseline
  file changed.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
  exited 0 after the final guide edit.
- `git diff --check` exited 0. Scope inspection found only Task 02's declared
  paths, including its new test file; the pre-existing Task status was retained.

Declared Verification and repository Verification remain unrun for the
Daemon. No Task status, other Task, Task Graph, QA settlement section,
Makefile or CI edit was made. No commit, push or Pull Request was made.
