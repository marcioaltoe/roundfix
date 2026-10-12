---
task: task_01
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: completed
type: docs
complexity: medium
---

# Task 01: The skills, guides and glossary describe the record's pull request, its recovery and the archive warning

## Overview

The CLI guide requires the skills and user guides to change with the
behavior, and the Skill Coverage Map names the Roundfix references that cover
the archive, review, runs and deliver commands. This Task writes the texts of
`_techspec.md` → Skill and document texts, before any code Task, so every
code Task depends on the Task that documents its surface. It adds the
glossary terms **Source Revision Recovery** and **Pull Request Record Commit**
and revises **Archive Record** and **Archive Command** (ADR-0261). It is
verifiable on its own through phrase checks, mirror equality, the owned-skill
version record and the Behavior Surface Record.

This is an authorized tooling Task. It may change only the files in its
Context and this Task file.

## Requirements

1. MUST answer the two adopted Backlog Entries of 2026-10-12,
   [the squash entry](references/2026-10-12-an-archive-record-points-at-a-branch-commit-a-squash-merge-drops.md)
   and [the review entry](references/2026-10-12-review-reads-any-markdown-in-the-archive-root-as-an-archive-record.md),
   on the documentation side, with every text of `_techspec.md` → Skill and
   document texts placed in the file it names.
2. MUST revise **Archive Record** and **Archive Command** and add **Source
   Revision Recovery** and **Pull Request Record Commit** in `CONTEXT.md`
   through `domain-modeling`, with the texts and `_Avoid_` lines of
   `_techspec.md` → Skill and document texts. Each existing `_Avoid_` line
   MUST stay.
3. MUST write the archive-spec text under a new `### Squash merges` heading.
   It MUST NOT add or change any text inside the `### QA settlement` section
   of any skill.
4. MUST replace, not merely add to, the sentences that say the review skips
   context and the Run causes report marks a revision unavailable "if the
   commit is absent". Each sentence becomes: the reader first fetches the pull
   request head ref, and only when that fails does it skip or mark
   unavailable. The search MUST cover the canonical skill, its mirror and the
   user guides.
5. MUST describe the Pull Request Record Commit as a step of the Delivery
   Queue between opening the pull request and judging checks, in the deliver
   reference and user guide. It MUST describe the hand-delivery step in
   `spec-delivery.md`.
6. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises both version fields of the Roundfix and archive-spec skills,
   then
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   It MUST NOT hand-edit a version, a mirror or the Behavior Surface Record.
7. MUST NOT describe behavior beyond `_techspec.md` → API Contracts 1-5. It
   MUST NOT change any Go file.

## Subtasks

- [ ] Write the archive, review, runs, deliver and spec-delivery reference texts.
- [ ] Write the archive-spec squash section.
- [ ] Write the four user guide texts and replace the "commit is absent" sentences.
- [ ] Add and revise the four glossary entries.
- [ ] Sync the mirrors, record the skill versions and the Behavior Surface Record.

## Acceptance Criteria

- [ ] Every phrase of `_techspec.md` → Skill and document texts is present in
      the canonical skill, its mirror and the matching user guide.
- [ ] `CONTEXT.md` carries the four entries.
- [ ] The `### QA settlement` tables are byte-identical.
- [ ] Mirrors, skill versions and the Behavior Surface Record are current.

## Context

- instruction: `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/references/2026-10-12-squash-merged-archive-source-revisions.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/spec-delivery.md`
- interface: `skills/roundfix/references/spec-delivery.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/review.md`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md docs/user-guide/commands/archive.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "stores the delivery pull request in the Archive Record" || { printf 'missing pull request phrase in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qE -- "warns when .source_revision. is not an ancestor of the default branch" || { printf 'missing warning phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md docs/user-guide/commands/review.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE -- "that declares the Archive Record schema is read as a record" || { printf 'missing record identity phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md docs/user-guide/commands/review.md .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md docs/user-guide/commands/runs.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE -- "fetches .refs/pull/<N>/head. when the commit is absent" || { printf 'missing recovery phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md docs/user-guide/commands/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "records the pull request number in the Archive Record in a Pull Request Record Commit" || { printf 'missing record commit phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/spec-delivery.md skills/roundfix/references/spec-delivery.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- "commit the record and push" || { printf 'missing hand step in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md; do grep -q '^### Squash merges$' "$f" || { printf 'missing squash heading in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qE -- "A squash merge keeps .source_revision. only in the pull request's head ref" || { printf 'missing squash phrase in %s\n' "$f" >&2; exit 1; }; done` — expected: exit 0; before this Task the archive reference has no pull request phrase, so the command fails at its first check.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Archive Record**:" || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "any other Markdown file in that root is not a record" || { printf 'CONTEXT.md lacks the revised **Archive Record**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Archive Command**:" || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "stores the delivery pull request in the record, also in a record the current branch already added" || { printf 'CONTEXT.md lacks the revised **Archive Command**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Source Revision Recovery**: A reader's fetch of" || { printf 'CONTEXT.md lacks **Source Revision Recovery**\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Pull Request Record Commit**: The Delivery Queue's commit, made after it opens a Spec's pull request" || { printf 'CONTEXT.md lacks **Pull Request Record Commit**\n' >&2; exit 1; }` — expected: exit 0; before this Task neither new term exists, so the command fails.
- `grep -q '^### Squash merges$' .agents/skills/archive-spec/SKILL.md || exit 1; out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync|TestArchiveSpecContract|TestSettlementGuidanceIsOneTable)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync TestArchiveSpecContract TestSettlementGuidanceIsOneTable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -tags docscontract -v -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheSkillCoverageMapIsCurrent" || exit 1` — expected: exit 0; before this Task the archive-spec skill has no squash heading, so the command fails at its first check.

## References

- [_prd.md](_prd.md) — Goal 5; Core Feature 6; User Story 2
- [_techspec.md](_techspec.md) — Skill and document texts; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Build Order 1
- ADR-0261; ADR-0247; ADR-0256; ADR-0189; ADR-0233; ADR-0244

## Result

Implemented the documentation slice for the archive pull request identity and
squash recovery:

- Added the archive warning and `--pull-request <n>` record guidance.
- Added Archive Record identity filtering and pull request head recovery to the
  review guidance, and recovery-first Run causes guidance.
- Added the Delivery Queue's Pull Request Record Commit between pull request
  creation and check judgment, plus the hand-delivery command sequence.
- Added the archive-spec `### Squash merges` section above the unchanged QA
  settlement table.
- Revised the four glossary entries in `CONTEXT.md), preserving all existing
  `_Avoid_` lines and adding Source Revision Recovery and Pull Request Record
  Commit.
- Generated the skill mirrors, owned-skill version records
  (`archive-spec` 0.0.9 and `roundfix` 0.1.60), and Behavior Surface Record.

Focused evidence:

- `make skills-sync` — passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1` — passed.
- Focused whitespace-normalized phrase checks passed for all canonical
  references, mirrors, and matching user guides.
- Focused byte comparisons passed for the `### QA settlement` tables in
  `archive-spec` and `roundfix` canonical and mirror skills.
- `git diff --check` — passed.

The declared Task Verification commands were left for the Daemon. Task status
remains Daemon-owned and unchanged.

## Carry-forward provenance

- Source Run: `run_20261012T140955Z_23b372ef5c35c71c`
- Source commit: `ef8fe2fa6cb74a040a4b4ed23a0fe969ac8c2c25`
