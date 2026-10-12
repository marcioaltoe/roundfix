---
task: task_06
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: completed
type: docs
complexity: low
---

# Task 06: The runs guidance admits Source Revision Recovery and the delivery texts match the TechSpec

## Overview

The QA gate of 2026-10-12 (`qa/qa-report-2026-10-12.md`, in Git at run
`run_20261012T140955Z_23b372ef5c35c71c`) failed on two documentation findings.
QA-001: three `runs` texts still promise that the command writes nothing and
makes no network request, although Source Revision Recovery can fetch a pull
request head and add the missing Git objects. QA-002: five texts word the
archive's `--pull-request <n>` flag and the hand-delivery step differently from
the phrases `_techspec.md` → "Skill and document texts" requires verbatim. This
corrective Task fixes both. It changes no code.

## Requirements

1. MUST answer QA-001 of the QA report dated 2026-10-12: rewrite the paragraph
   at `docs/user-guide/commands/runs.md` (around line 96) and the paragraph in
   `.agents/skills/roundfix/references/runs.md` (around line 196) so that
   neither claims the command writes nothing or makes no network request. Each
   states, consistently with ADR-0261 and the later paragraphs of the same
   file, that the command changes no file and no ref, and that when a record's
   `source_revision` is missing locally and the record names its pull request,
   Source Revision Recovery fetches that pull request's head from the
   configured remote and adds the missing objects.
2. MUST answer QA-002 of the same report: make the archive reference
   (`.agents/skills/roundfix/references/archive.md`) and
   `docs/user-guide/commands/archive.md` contain, after whitespace
   normalization, the TechSpec phrase in which `--pull-request <n>`, in its own
   code span, "stores the delivery pull request in the Archive Record"; and make
   the hand-delivery reference
   (`.agents/skills/roundfix/references/spec-delivery.md`) contain the TechSpec
   phrase that begins with a lowercase "after the pull request opens, run" and
   names `roundfix archive <slug> --pull-request <n>`, then "commit the record
   and push". Keep every other sentence's meaning.
3. MUST sync the mirrors with `make skills-sync`, raise the `roundfix` skill's
   version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and never touch any `### QA settlement` section.

## Subtasks

- [ ] Correct the three runs paragraphs.
- [ ] Match the five TechSpec phrases.
- [ ] Sync mirrors and record the skill version.

## Acceptance Criteria

- [ ] No runs text promises no writes or no network request.
- [ ] The archive and hand-delivery texts carry the TechSpec phrases verbatim.
- [ ] Mirrors match and the skill version is recorded.

## Context

- instruction: `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `docs/user-guide/commands/runs.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `docs/user-guide/commands/archive.md`
- interface: `.agents/skills/roundfix/references/spec-delivery.md`
- interface: `skills/roundfix/references/spec-delivery.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `for f in docs/user-guide/commands/runs.md .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md; do if tr -s '[:space:]' ' ' < "$f" | grep -qE 'writes nothing and makes no network request|without writes or network requests'; then printf 'stale no-network claim in %s\n' "$f" >&2; exit 1; fi; done; for f in .agents/skills/roundfix/references/archive.md skills/roundfix/references/archive.md docs/user-guide/commands/archive.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE '[[:punct:]]--pull-request <n>[[:punct:]] stores the delivery pull request in the Archive Record' || { printf 'missing TechSpec flag phrase in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/spec-delivery.md skills/roundfix/references/spec-delivery.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE 'after the pull request opens, run [[:punct:]]roundfix archive <slug> --pull-request <n>[[:punct:]], commit the record and push' || { printf 'missing TechSpec hand-delivery phrase in %s\n' "$f" >&2; exit 1; }; done; for f in archive runs spec-delivery; do cmp -s ".agents/skills/roundfix/references/$f.md" "skills/roundfix/references/$f.md" || { printf 'reference mirror differs: %s\n' "$f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || { echo 'roundfix SKILL.md mirror differs' >&2; exit 1; }; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills` — expected: exit 0. Before this Task the three runs texts carry the stale claim and the five phrases differ, so it fails.

## References

- `qa/qa-report-2026-10-12.md` → QA-001, QA-002
- `_techspec.md` → Skill and document texts
- ADR-0261

## Result

Implementation and focused checks for daemon handoff:

- Reworded the canonical and user-guide `runs` paragraphs to state that the
  command changes no repository file or Git ref, while Source Revision Recovery
  can fetch a named pull request head from the configured remote and add missing
  Git objects.
- Matched the TechSpec archive phrase with `--pull-request <n>` in its own code
  span in both canonical references and the user guide, and matched the
  lowercase hand-delivery phrase in the canonical and mirrored references.
- Ran `make skills-sync`; the three reference mirrors and the `roundfix`
  `SKILL.md` mirror are synchronized.
- Ran `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions`: passed (`1 passed in 1 packages`). The `roundfix`
  skill is recorded at version `0.1.61`.
- Focused scans passed: stale runs wording absent from all three runs texts;
  whitespace-normalized TechSpec phrases present in all required files; mirror
  comparisons passed; `git diff --check` passed.

Acceptance evidence:

- No runs text retains either stale no-write/no-network claim, and both runs
  paragraphs document the bounded Source Revision Recovery behavior.
- The archive and hand-delivery texts contain the required whitespace-normalized
  TechSpec phrases.
- Canonical and mirrored skill files match, and the owned-skill version record
  contains the new `roundfix` version and digest.

Verification feedback repair:

- The Daemon's `make verify-changed` failed in the repository contract phase
  because the two edited user-guide behavior-surface fingerprints were stale.
- Reproduced the finding with `go test -tags docscontract ./internal/docscontract
  -run '^TestTheSkillCoverageMapIsCurrent$' -count=1`; it reported stale
  fingerprints only for the archive and runs guides.
- Ran the authorized regeneration command
  `go test -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`;
  it passed and updated only those two entries in
  `docs/references/behavior-surfaces.json`.
- Reran the focused docscontract test without the record flag; it passed.
  `git diff --check` also passed. The Task status remains daemon-owned and
  unchanged.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/references/behavior-surfaces.json`
