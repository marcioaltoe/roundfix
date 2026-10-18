---
task: task_02
spec: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
status: pending
type: backend
complexity: high
---

# Task 02: The Secret Scan reads every Task commit, and test authors learn which fixture values pass

## Overview

On 2026-10-18 a Conexus test fixture matched `credential-assignment` only at
the final push, after every Task had committed it and after QA, review and
archive (`_prd.md` → Sources 2). The Secret Scan reads the QA Report commit and
every push, but never a Task commit.

This Task scans the lines each Task commit would add before the Task settles,
with the same added-line reading the push uses, so a fixture is refused by the
Task that adds it. It also documents fixture values the pattern does not match.
It is verifiable on its own through scanner and Daemon tests on temporary
repositories.

## Requirements

1. MUST answer the Finding of 2026-10-18,
   `references/2026-10-18-secret-found-asks-for-a-rewrite-the-retry-of-an-archived-item-refuses.md`
   (`_prd.md` → Sources 2, finding 2), by implementing API Contracts 4 and 5
   of `_techspec.md` (ADR-0292).
2. MUST add `ScanAddedLines(patch []byte) []Match` to
   `internal/secretscan/scan.go` (API Contract 4). It shares
   `ScanPushedHistory`'s hunk parser, reports added lines with their new line
   numbers and an empty `Commit`, and leaves `ScanPushedHistory`'s output
   unchanged. No Secret Pattern changes.
3. MUST scan the Task commit in `executeTask` as API Contract 5 states.
   - The scan runs after `prepareTaskCommit` and after the lost-path and
     authorization checks leave no failure, and before settlement.
   - It reads the prepared stageable paths:
     - a path that is absent contributes nothing;
     - a path `HEAD` does not hold is read whole with `ScanContent`;
     - every other path is read through one `git diff --no-color
       --no-ext-diff --no-renames -U0 HEAD --` and `ScanAddedLines`.
   - A Secret Match sets the failure `Task commit carries a secret: <match>;
     <match>`. The Daemon publishes one `daemon.task` event with phase
     `secret_scan` and each match's path, line and kind, and the Task settles
     `failed` through the existing path, with no Task commit.

   The data source is the Run Worktree and `HEAD`. The reason, the event and
   the console line are read by the operator and the Run Database, and they
   carry each match's path, line and kind only, never its value. A Git or
   read error MUST fail the Run as an infrastructure error, never settle the
   Task `completed`.
4. MUST add these tests, each calling `t.Parallel()` first and compiling for
   Windows. Every secret-shaped sample is built at runtime from parts
   (ADR-0273), so no file holds a literal value that matches a Secret Pattern.
   - In the new `internal/secretscan/added_lines_test.go`, which starts no
     process: `TestScanAddedLinesReportsOnlyAddedLines`. A patch with an
     added, a removed and a context line that each match reports only the
     added one, at its new line number.
   - In the new `internal/daemon/task_secret_scan_test.go`, reusing the Task
     fixtures already in the package:
     - `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit`: the Task
       settles `failed` with the exact reason naming `credential-assignment`,
       the event carries no value, and no Task commit exists;
     - `TestTaskCommitScansOnlyTheLinesItAdds`: a Task that edits another line
       of a file whose committed line already matches commits;
     - `TestTaskCommitAcceptsALettersOnlyFixture`: a Task that adds an
       assignment of a letters-only value of 30 characters commits.

   Every existing test MUST keep its name and pass, including the
   `secretscan`, push and QA Secret Scan tests.
5. MUST document the change:
   - `docs/user-guide/commands/implement.md` and
     `.agents/skills/roundfix/references/implement.md` each state that the
     Secret Scan reads the lines each Task commit adds, with the reason `Task
     commit carries a secret: <path>:<line>: <kind>`.
   - Each also gains a part headed `Fixture values` that says the
     `credential-assignment` pattern matches a name ending in `API_KEY`,
     `TOKEN`, `SECRET` or `PASSWORD`, an assignment, and a value of 20 or more
     characters that holds a letter and a digit. It also says that a fixture
     value of letters only, or of fewer than 20 characters, does not match.
   - `.agents/skills/implement-task/SKILL.md` gains the same `Fixture values`
     rule under its own heading.

   None of these texts may contain a value that matches a Secret Pattern.
   MUST search the guides and skill copies for text saying the Secret Scan
   covers only the QA Report commit and pushes, and qualify each. No text
   inside any `### QA settlement` section may change.
6. MUST revise **Secret Scan** in `CONTEXT.md` through domain-modeling, with
   ADR-0292's link, so it names the lines a Task commit would add.
7. MUST run these commands, because this Task edits a user guide page, the
   Roundfix reference and the `implement-task` skill. No record or mirror is
   hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to raise and record the Roundfix and `implement-task` skill versions;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add `ScanAddedLines`.
- [ ] Add the Task commit scan to settlement.
- [ ] Add the four tests.
- [ ] Document the scan and the fixture values in the guide, the reference, the `implement-task` skill and the glossary.
- [ ] Sync the mirrors, record the versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A Task that adds a matching line settles `failed` with its location and makes no commit.
- [ ] A Task that only edits beside an already committed match, or adds a letters-only fixture, commits.
- [ ] The guide, the reference, the skill and the glossary describe the scan and the fixture values.

## Consumers

- `secretscan.ScanAddedLines` — consumers: `./internal/daemon`

## Context

- instruction: `docs/adr/0292-a-task-commit-is-scanned-for-secrets-and-a-retry-accepts-a-secret-rewrite.md`
- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- instruction: `internal/secretscan/patterns.go`
- instruction: `internal/daemon/qa_secret_scan_test.go`
- creates: `internal/secretscan/added_lines_test.go`
- creates: `internal/daemon/task_secret_scan_test.go`
- interface: `internal/secretscan/scan.go`
- interface: `internal/daemon/task_engine.go`
- interface: `docs/user-guide/commands/implement.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/implement-task/SKILL.md`
- interface: `skills/implement-task/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^TestScanAddedLinesReportsOnlyAddedLines$' ./internal/secretscan 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestScanAddedLinesReportsOnlyAddedLines ' || { printf 'missing pass: TestScanAddedLinesReportsOnlyAddedLines\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit|TestTaskCommitScansOnlyTheLinesItAdds|TestTaskCommitAcceptsALettersOnlyFixture)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit TestTaskCommitScansOnlyTheLinesItAdds TestTaskCommitAcceptsALettersOnlyFixture; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/implement.md .agents/skills/roundfix/references/implement.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Task commit carries a secret' || { printf 'missing Task commit refusal in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/implement.md .agents/skills/roundfix/references/implement.md .agents/skills/implement-task/SKILL.md; do grep -qF -- 'Fixture values' "$f" || { printf 'missing Fixture values in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'a Task commit would add' || { printf 'missing revised glossary term Secret Scan\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Secret Scan**:' || { printf 'missing glossary term Secret Scan\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/secretscan || exit 1; go test -count=1 -run 'Task|Secret|QA' ./internal/daemon || exit 1; go test -count=1 -run 'Secret' ./internal/delivery || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the four tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the scanner and settlement tests, the documentation, the glossary, the skill versions and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 3; Core Feature 3; Core Feature 5; Success Metric 3; Glossary
- `_techspec.md` → API Contract 4; API Contract 5; API Contract 9; Invariant 2; Current behavior; Testing Approach; Build Order 2
- ADR-0292; ADR-0273; ADR-0117; ADR-0280; ADR-0189; ADR-0256; ADR-0259; ADR-0285
