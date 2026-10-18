---
task: task_02
spec: 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
status: pending
type: backend
complexity: high
---

# Task 02: A delivered Spec's leftovers beneath its declared paths are superseded, and a kept dirty Run names why

## Overview

After Spec 0272 merged (pull request #555), reconcile kept three of its Runs
`dirty` with only "Run Worktree has tracked or untracked changes", and the
maintainer had them removed by hand on 2026-10-18 (`_prd.md` → Sources 2). The
leftover rule compares each dirty path with the declared and recorded paths by
exact equality, so a declaration such as `deletes: docs/history/findings`
covers no file beneath it, and the bare reason hides which proof failed. This
Task makes a declared or recorded path cover the paths beneath it and makes a
kept dirty Run name the failing proof or the first path outside the scope. It
is verifiable on its own through reconciliation tests on temporary
repositories.

## Requirements

1. MUST answer the operator observation of 2026-10-18 (Spec 0272, pull request
   #555) cited in `_prd.md` → Sources 2 by implementing API Contracts 6 and 7
   of `_techspec.md` (ADR-0289).
2. MUST make the leftover scope cover a dirty path that equals or lies beneath
   (`<declared>/` prefix) a path a Task of the archived Spec declares under
   `interface:`, `creates:` or `deletes:` or records under `## Recorded paths`,
   read at the Archive Record's `source_revision` as today. A sibling that only
   shares the prefix stays out of scope.
3. MUST keep every other input of the classification unchanged: the commit
   proof at the merged head or by the delivery, the archive requirement, and
   the apply-time re-proof of the exact dirty path set. A dirty path outside
   the scope still keeps the Run, so no change outside what the delivered Spec
   declared, recorded or holds in its own directories is ever released.
4. MUST build the dirty reason of API Contract 7: when the Spec is archived at
   the chosen merged head or at the default branch and the Run stays `dirty`,
   the constant, `; ` and either the commit proof's own reason or
   `<path> is outside the archived Spec's declared and recorded paths` for the
   first out-of-scope path in byte order, bounded to 160 bytes without cutting
   the constant. A Run whose Spec is not archived there keeps the bare
   constant, and an unavailable archive source keeps today's reason.
5. MUST add these tests in the new
   `internal/worktree/merged_spec_leftover_scope_test.go`, each calling
   `t.Parallel()` first, compiling for Windows and using the fixtures of
   `merged_spec_leftovers_test.go`:
   - `TestMergedSpecRunWithLeftoversBeneathADeclaredDirectoryIsSuperseded`: a
     dirty deletion beneath a declared `deletes:` directory and a dirty file
     beneath a recorded directory are `superseded`, and `ApplyTerminalRun`
     removes the worktree and branch; a sibling sharing the prefix stays
     `dirty`;
   - `TestMergedSpecRunNamesTheLeftoverOutsideItsScope`: the reason equals the
     constant followed by `; unrelated.txt is outside the archived Spec's
     declared and recorded paths`;
   - `TestMergedSpecRunNamesTheUnrepresentedCommit`: an unrepresented Task
     commit plus a scoped leftover keeps the Run `dirty` with a reason that
     begins with the constant and names the commit.
6. MUST declare and make this break: in
   `internal/worktree/merged_spec_leftovers_test.go`,
   `TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty` stops asserting the
   bare constant and asserts the named reason for each of its subtests.
   `TestMergedHeadKeepsADirtyRunWorktree` keeps asserting the bare constant,
   because its Spec is not archived. No test is renamed or removed.
7. MUST document the change, each naming `beneath a declared`:
   - in `docs/user-guide/commands/reconcile.md`, the paragraph on dirty Run
     Worktrees and the `dirty` row of its state table, with the reason forms;
   - in `.agents/skills/roundfix/references/reconcile.md`.

   MUST search the guides and skill copies for text saying a dirty path must
   equal a declared path, and qualify each. No text inside any
   `### QA settlement` section may change.
8. MUST revise **Run Worktree Reconciliation** in `CONTEXT.md` through
   domain-modeling with ADR-0289's link, so it states that a declared or
   recorded path covers the dirty paths beneath it and that a kept dirty Run
   names why; the revised entry contains `names the failing proof`.
9. MUST run these commands, because this Task edits a user guide and a Roundfix
   reference. No record or mirror is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to raise and record the Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Widen the leftover scope to paths beneath a declared or recorded path.
- [ ] Build the named dirty reason.
- [ ] Add the three tests and move the declared break.
- [ ] Document the change in the guide, the reference and the glossary.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] A leftover beneath a declared directory is `superseded` and released by apply, and a prefix-sharing sibling stays `dirty`.
- [ ] A kept dirty Run of an archived Spec names the out-of-scope path or the unrepresented commit.
- [ ] A dirty Run of an unarchived Spec keeps the bare reason.

## Consumers

- Run Worktree Reconciliation's state and reason — consumers: `./internal/worktree`, `./internal/cli`

## Context

- instruction: `docs/adr/0289-a-delivered-specs-dirty-run-is-superseded-under-its-declared-directories-and-a-kept-one-names-why.md`
- instruction: `docs/adr/0212-a-merged-spec-supersedes-its-runs-spec-directory-work-and-declared-leftovers.md`
- instruction: `docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md`
- interface: `internal/worktree/merged_head.go`
- interface: `internal/worktree/worktree.go`
- interface: `internal/worktree/merged_spec_leftovers_test.go`
- instruction: `internal/worktree/merged_head_test.go`
- creates: `internal/worktree/merged_spec_leftover_scope_test.go`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `.agents/skills/roundfix/references/reconcile.md`
- interface: `skills/roundfix/references/reconcile.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestMergedSpecRunWithLeftoversBeneathADeclaredDirectoryIsSuperseded|TestMergedSpecRunNamesTheLeftoverOutsideItsScope|TestMergedSpecRunNamesTheUnrepresentedCommit|TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty|TestMergedHeadKeepsADirtyRunWorktree)$' ./internal/worktree 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMergedSpecRunWithLeftoversBeneathADeclaredDirectoryIsSuperseded TestMergedSpecRunNamesTheLeftoverOutsideItsScope TestMergedSpecRunNamesTheUnrepresentedCommit TestMergedSpecRunWithAnUndeclaredLeftoverStaysDirty TestMergedHeadKeepsADirtyRunWorktree; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- 'is outside the archived Spec' internal/worktree/merged_spec_leftovers_test.go || { printf 'the declared break was not made\n' >&2; exit 1; }; for f in docs/user-guide/commands/reconcile.md .agents/skills/roundfix/references/reconcile.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'beneath a declared' || { printf 'missing directory scope in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'names the failing proof' || { printf 'missing revised Run Worktree Reconciliation\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Run Worktree Reconciliation**:' || { printf 'missing glossary term Run Worktree Reconciliation\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/worktree || exit 1; go test -count=1 -run 'Reconcile|Release|Merged' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check. Afterwards these checks pass: the reconciliation tests and the declared break, the documentation, the glossary, the skill version and the Behavior Surface Record.

## References

- `_prd.md` → Sources; Goals; User Story 3; User Story 4; Core Feature 3; Core Feature 4; Core Feature 7; Success Metric 2; Glossary
- `_techspec.md` → API Contract 6; API Contract 7; API Contract 10; Invariants 3 and 4; Current behavior; Testing Approach; Build Order 2
- ADR-0289; ADR-0212; ADR-0053; ADR-0115; ADR-0280; ADR-0189; ADR-0256; ADR-0259
