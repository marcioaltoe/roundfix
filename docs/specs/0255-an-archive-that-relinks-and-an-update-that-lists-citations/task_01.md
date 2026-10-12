---
task: task_01
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
status: completed
type: docs
complexity: medium
---

# Task 01: The glossary, the command references and the skills describe the Archive Relink and History Citations

## Overview

Writes ADR-0260 and the TechSpec into the durable guides before any code
changes. `CONTEXT.md` gains the four new terms and revises three. The archive,
baseline and deliver command references, the Roundfix Skill, the archive skill
and the four authoring skills describe the Archive Relink, the History
Citations of `roundfix baseline update`, `--promote` on the update and the
Stale History Citation report. The baseline references carry the adopter note
on citations an earlier update already broke. The Task is verifiable on its own
through phrase, mirror and version checks.

## Requirements

1. MUST add to `CONTEXT.md`, through the `domain-modeling` skill and in its
   existing entry format (bolded term, colon, one definition line, an `_Avoid_`
   line), the terms **Live Document**, **Archive Relink**, **History Citation**
   and **Stale History Citation**, as `_prd.md` → Core Feature 7 and ADR-0260
   define them. The definitions cite ADR-0260 and never a Spec.
2. MUST revise in `CONTEXT.md` the definition line of **Archive Record** to
   name its optional `relinked` list, of **Archive Command** to name the
   Archive Relink, and of **Managed Refresh** to name the History Citations it
   lists and binds into the Plan Digest and its `--promote` flag.
3. MUST state, in `write-idea`, `write-prd`, `write-techspec` and
   `write-tasks`, next to the existing sentence that a secondary Spec links the
   owner's copy, that the owner's archive performs an Archive Relink: it
   promotes that copy to `docs/references/` and rewrites the secondary Spec's
   links, so the source keeps one owner and the secondary Spec needs no copy
   (ADR-0260, ADR-0083). The `archive-spec` skill MUST say what the archive
   promotes automatically, that it rewrites Live Document citations, what it
   refuses and what it warns about. New text goes under existing headings or
   its own heading, never inside a `### QA settlement` section.
4. MUST describe in `.agents/skills/roundfix/references/archive.md` and
   `docs/user-guide/commands/archive.md` the behavior of `_techspec.md` → API
   Contracts 3, 4 and 5: the automatic promotion, the rewrite, the `relinked`
   key, the three refusals, the `; relinked <m> file(s)` suffix, the
   `unrelinked` warning line and the new `--plan` lines.
5. MUST describe in `.agents/skills/roundfix/references/baseline.md` and
   `docs/user-guide/commands/baseline.md` the behavior of `_techspec.md` → API
   Contracts 7 to 10: the History Citation lines and JSON fields, the offer, the
   repeatable `--promote` and its refusals, and the Stale History Citation
   lines. Both MUST carry an adopter note saying that a repository whose
   earlier update or Sanitize Batch already broke citations runs
   `roundfix baseline update` with the History Full Tag present, reads the
   `History stale citations` lines, and restores a target it keeps with
   `git show history-full:<path>` into `docs/references/` before fixing the
   link; and that Specs already written to link an owner's copy need nothing,
   because the owner's archive relinks them.
6. MUST describe in `.agents/skills/roundfix/references/deliver.md` and
   `docs/user-guide/commands/deliver.md` that the archive commit also carries
   the promoted copies and relinked files the record names, and that the queue
   accepts a relinked file only when it equals the archive's rewrite of the
   parent's bytes (`_techspec.md` → API Contract 6).
7. MUST copy every changed canonical skill file to its mirror under `skills/`
   with `make skills-sync`, and raise each changed owned skill's version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which rewrites both front-matter fields and
   `skills/testdata/owned-skill-versions.json`. Another Spec of the same cycle
   may record first; the record command picks the next free version.
8. MUST search the documentation, agent guides and skill copies for the old
   wording ("A secondary Spec links", "promote with:", "no PRD or link is
   rewritten") and change every live occurrence this Spec contradicts. MUST
   NOT edit ADR-0083, ADR-0247, ADR-0254, any archived Spec or any Archive
   Record.

## Subtasks

- [ ] Add and revise the glossary terms through `domain-modeling`.
- [ ] Update the four authoring skills and the archive skill.
- [ ] Update the archive, baseline and deliver references and user guides, with the adopter note.
- [ ] Sync the mirrors and record the skill versions.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines the four new terms and the three revised ones name what changed.
- [ ] Each authoring skill and the archive skill name the Archive Relink, and every mirror equals its canonical copy.
- [ ] The baseline references tell an adopter how to find and restore citations an earlier update broke.
- [ ] The owned-skill version record and the QA settlement contract tests pass.

## Context

- interface: `CONTEXT.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
- instruction: `docs/adr/0083-adopted-sources-move-to-their-owning-spec.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Verification

- `f=CONTEXT.md; n="$(tr -s '[:space:]' ' ' < "$f")"; for t in '**Live Document**:' '**Archive Relink**:' '**History Citation**:' '**Stale History Citation**:' '**Archive Record**:' '**Archive Command**:' '**Managed Refresh**:'; do printf '%s' "$n" | grep -qF -- "$t" || { printf 'missing glossary term %s\n' "$t" >&2; exit 1; }; done; grep -A1 -F -- '**Archive Record**:' "$f" | grep -q 'relinked' || { echo 'Archive Record does not name relinked' >&2; exit 1; }; grep -A1 -F -- '**Archive Command**:' "$f" | grep -q 'Archive Relink' || { echo 'Archive Command does not name the Archive Relink' >&2; exit 1; }; grep -A1 -F -- '**Managed Refresh**:' "$f" | grep -q 'History Citation' || { echo 'Managed Refresh does not name History Citations' >&2; exit 1; }` — expected: exit 0. Before this Task the four new terms are absent and the three revised definitions do not name what changed, so it fails.
- `for s in write-idea write-prd write-techspec write-tasks archive-spec; do tr -s '[:space:]' ' ' < ".agents/skills/$s/SKILL.md" | grep -qF -- 'Archive Relink' || { printf 'skill %s does not name the Archive Relink\n' "$s" >&2; exit 1; }; cmp -s ".agents/skills/$s/SKILL.md" "skills/$s/SKILL.md" || { printf 'skill mirror differs: %s\n' "$s" >&2; exit 1; }; done; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills` — expected: exit 0. Before this Task no skill names the Archive Relink, so it fails; after it, each skill does, each mirror matches, every raised version is recorded and the QA settlement sections are unchanged.
- `check() { tr -s '[:space:]' ' ' < "$1" | grep -qF -- "$2" || { printf 'missing phrase in %s: %s\n' "$1" "$2" >&2; exit 1; }; }; for d in .agents/skills/roundfix/references docs/user-guide/commands; do check "$d/archive.md" 'Archive Relink'; check "$d/archive.md" 'unrelinked'; check "$d/baseline.md" 'History Citation'; check "$d/baseline.md" 'Stale History Citation'; check "$d/baseline.md" 'git show history-full:'; check "$d/baseline.md" '--promote'; check "$d/deliver.md" 'relinked'; done; for f in archive baseline deliver; do cmp -s ".agents/skills/roundfix/references/$f.md" "skills/roundfix/references/$f.md" || { printf 'reference mirror differs: %s\n' "$f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || { echo 'roundfix SKILL.md mirror differs' >&2; exit 1; }` — expected: exit 0. Before this Task the references name neither the Archive Relink nor History Citations, so it fails.

## References

- `_prd.md` → Core Feature 7; User Story 6; Goals; Acceptance evidence
- `_techspec.md` → Build Order 1; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; API Contract 10
- ADR-0260; ADR-0083; ADR-0247; ADR-0254

## Result

Implementation evidence for the daemon handoff:

- Glossary: `CONTEXT.md` now defines Live Document, Archive Relink, History
  Citation and Stale History Citation with ADR-0260 references. Archive Record
  names its optional `relinked` list, Archive Command names Archive Relink,
  and Managed Refresh names History Citations, the Plan Digest and repeatable
  `--promote`.
- Skills and mirrors: the four authoring skills and archive-spec now describe
  the owner's Archive Relink; the Roundfix Skill and archive, baseline and
  deliver references describe the new surfaces. `make skills-sync` completed,
  and a focused equality check reported `canonical/mirror equality: ok`.
- Baseline adopter guidance: both baseline references list the History Citation
  and Stale History Citation text forms, JSON fields, promotion refusals and
  the `git show history-full:<path>` restoration procedure.
- Version evidence: `GOCACHE=/private/tmp/roundfix-task-0255-gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  completed with `1 passed in 1 packages`; it recorded the next versions and
  digests for the six changed owned skills.
- Focused hygiene: `git diff --check` completed successfully; the changed
  documentation and generated version record contain no whitespace errors.
- Verification feedback repair: Attempt 1 of `make verify-changed` passed all
  earlier packages and checks, then failed only because
  `TestTheSkillCoverageMapIsCurrent` reported stale fingerprints for the
  archive, baseline and deliver guides. The repository-authorized record
  command
  `GOCACHE=/private/tmp/roundfix-task-0255-gocache go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract -record-skill-coverage`
  regenerated `docs/references/behavior-surfaces.json`, and the focused
  follow-up
  `GOCACHE=/private/tmp/roundfix-task-0255-gocache go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract`
  passed with `1 passed in 1 packages`.

The Task-declared Verification commands, including the QA settlement contract
test, remain for the Daemon as required by the execution invariants. No Task
status was changed, and no commit, push or PR was created.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/references/behavior-surfaces.json`

## Carry-forward provenance

- Source Run: `run_20261012T110224Z_380dd005763aa861`
- Source commit: `314f67220c0bf1ad3e0fe49b099783d0eaa681c3`
