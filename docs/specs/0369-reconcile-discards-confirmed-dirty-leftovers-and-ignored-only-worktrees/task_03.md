---
task: task_03
spec: 0369-reconcile-discards-confirmed-dirty-leftovers-and-ignored-only-worktrees
status: completed
type: docs
complexity: low
---

# Task 03: The amended ADRs and the Roundfix skill's reconcile reference describe both

## Overview

ADR-0379 amends ADR-0053, ADR-0289 and ADR-0290, and the Roundfix skill's
`reconcile` reference still says three mutation switches exist and that an
unregistered worktree is `unknown`. This Task writes the three amendment
notices and documents the Dirty Discard, the ignored-only release and the
empty repository directory sweep that task_01 and task_02 delivered in the
skill reference, so agents reading the skill act as the command does. It
answers the Fiscus report of 2026-11-09
([`references/2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md`](references/2026-11-09-fiscus-reconcile-leaves-abandoned-run-leftovers-without-a-path.md),
reported by fiscus), whose cleanup the reference forbade by hand.

## Requirements

1. MUST add to ADR-0053, after its decision paragraph, the paragraph
   `**Amended by [ADR-0379](0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md) (2026-11-09).**`
   followed by: dirty work of an ended Run whose Spec's Archive Record is on
   the default branch and whose commit proof holds now has one explicit exit,
   the Dirty Discard (`roundfix reconcile --discard-dirty <digest>`), which
   records a Branch Disposition naming the discarded paths before it removes
   anything; every other dirty, ambiguous or unintegrated work stays intact,
   and `--apply` still releases only proven work. Set its `updated_at` to
   `2026-11-09T15:00:00Z`.
2. MUST add to ADR-0289, after its Decision section's last paragraph, the
   same `Amended by [ADR-0379]` paragraph lead as Requirement 1, followed by: a kept
   dirty Run whose only failing proof is a path outside the archived Spec's
   scope can now be removed by an operator-confirmed Dirty Discard; the first
   Considered Option, releasing every dirty Run of a merged Spec, stays
   rejected for every automatic path, because the confirmation is the
   judgement it asked for. Set its `updated_at` to `2026-11-09T15:00:00Z`.
3. MUST add to ADR-0290, after its Decision table, the same paragraph lead
   as Requirement 1, followed by: the Run Worktree and Run Branch row gains the Dirty Discard as a
   reclaimer that acts on an operator's confirmed digest, and a recorded Run
   Worktree that Git does not register and that holds only ignored files is
   reclaimed by `roundfix reconcile --apply` on proof; neither reclaims by
   age. Set its `updated_at` to `2026-11-09T15:00:00Z`. No other line of the
   three ADRs changes.
4. MUST, in the canonical `.agents/skills/roundfix/references/reconcile.md`:
   - in the state table, say in the `unknown` row that an unregistered Run
     Worktree that holds only files the repository ignores is `safe`
     instead, and in the `dirty` row that a Run the dry run lists with
     `discard-path:` lines can leave through `--discard-dirty`;
   - add, under its own heading `### Dirty Discard` placed before "After
     reviewing the dry-run, apply cleanup explicitly:", the bounds of a Dirty
     Discard, the `discard-path:` lines and `Discard dirty leftovers with:`
     line, the JSON fields `discardPaths`, `discardDirtyDigest` and
     `discardDirtyCommand`, that a changed set or no candidate exits 2 with
     nothing on stdout, and that the switch writes the Branch Disposition
     record with `discardedPaths` and `discardDigest` before it removes the
     Run Worktree and Run Branch;
   - add a paragraph on the ignored-only release and the full-scan empty
     repository directory sweep with `emptyDirectoryCandidates` and
     `emptyDirectoriesRemoved`;
   - add `roundfix reconcile --discard-dirty <digest>` to the mutation
     command block, change "These three mutation switches are mutually
     exclusive." to "These four mutation switches are mutually exclusive.",
     and keep the closing rule that forbids manual Git deletion, adding that
     a Dirty Discard is the supported exit for an ended Run's leftovers.
5. MUST run `make skills-sync` to write the mirror
   `skills/roundfix/references/reconcile.md`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   to raise both version fields of the Roundfix `SKILL.md` pair and record
   them in `skills/testdata/owned-skill-versions.json`; never hand-edit the
   mirror, a version field or the record.
6. MUST, as the last edit of the Task, re-record the Behavior Surface Record
   with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
7. MUST NOT change any code, test, guide, `CONTEXT.md`, another skill or
   reference, or text inside the `### QA settlement` section of any skill.

## Subtasks

- [ ] Write the three amendment notices.
- [ ] Document the switch, the release and the sweep in the canonical reference.
- [ ] Sync the mirror, record the skill version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] ADR-0053, ADR-0289 and ADR-0290 each carry the `Amended by [ADR-0379]` paragraph and the new `updated_at`.
- [ ] The canonical reference and its mirror describe four mutation switches, the Dirty Discard and the ignored-only release, and the skill version record is current.

## Context

- interface: `docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md`
- interface: `docs/adr/0289-a-delivered-specs-dirty-run-is-superseded-under-its-declared-directories-and-a-kept-one-names-why.md`
- interface: `docs/adr/0290-merged-items-give-back-their-item-binaries-and-run-worktrees-leave-only-by-proof.md`
- interface: `.agents/skills/roundfix/references/reconcile.md`
- interface: `skills/roundfix/references/reconcile.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- instruction: `docs/adr/0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md`
- instruction: `docs/user-guide/commands/reconcile.md`

## Consumers

- The Roundfix skill's `reconcile` reference and its version record — consumers: `./skills`

## Verification

- `for f in docs/adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md docs/adr/0289-a-delivered-specs-dirty-run-is-superseded-under-its-declared-directories-and-a-kept-one-names-why.md docs/adr/0290-merged-items-give-back-their-item-binaries-and-run-worktrees-leave-only-by-proof.md; do grep -qF -- '**Amended by [ADR-0379](0379-an-archived-specs-ended-run-gives-up-its-dirty-leftovers-only-on-a-confirmed-discard.md) (2026-11-09).**' "$f" || { echo "missing amendment notice: $f" >&2; exit 1; }; grep -qx 'updated_at: 2026-11-09T15:00:00Z' "$f" || { echo "missing updated_at: $f" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the three ADRs carries the notice.
- `for f in .agents/skills/roundfix/references/reconcile.md skills/roundfix/references/reconcile.md; do grep -qx '### Dirty Discard' "$f" || { echo "missing heading: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'These four mutation switches are mutually exclusive.' || { echo "missing switch count: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'discardDirtyDigest' || { echo "missing digest field: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'emptyDirectoriesRemoved' || { echo "missing sweep field: $f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/reconcile.md skills/roundfix/references/reconcile.md || exit 1; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — expected: exit 0. Before this Task the reference has no Dirty Discard heading; afterwards both copies match, the owned skill version record and the skill contracts pass, and the Behavior Surface Record is current.

## References

- `_prd.md` → Goals; Core Feature 5
- `_techspec.md` → ADR Alternatives; Build Order 3
- ADR-0379; ADR-0053; ADR-0289; ADR-0290

## Result

Skills loaded: implement-task, tech-writer, evidence-gate

- Added the ADR-0379 amendment notices to ADR-0053, ADR-0289 and ADR-0290 and set each `updated_at` to `2026-11-09T15:00:00Z`. Focused inspection confirmed the required notice and timestamp in all three files.
- Updated the canonical Roundfix `reconcile` reference with the Dirty Discard bounds, `discard-path:` output, digest JSON fields and refusal behavior, Branch Disposition ordering, ignored-only release, full-scan empty-directory sweep, four mutation switches and the manual-deletion rule. `git diff --check` passed.
- Ran `make skills-sync`; `cmp -s .agents/skills/roundfix/references/reconcile.md skills/roundfix/references/reconcile.md` passed.
- Ran `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` successfully with the required elevated filesystem permission after the sandbox denied the canonical skill write. The generated result raised both Roundfix `SKILL.md` copies to `0.1.207` and recorded the new digest in `skills/testdata/owned-skill-versions.json`.
- Acceptance evidence: the three ADR notices/timestamps are present; the canonical and mirrored references are identical and contain `### Dirty Discard`, `These four mutation switches are mutually exclusive.`, `discardDirtyDigest` and `emptyDirectoriesRemoved`; the owned-skill version record is current.
- Ran `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` successfully (1 passed); this re-record is repeated after this Result update as the final edit of the Task. Daemon Verification remains outside this turn.
