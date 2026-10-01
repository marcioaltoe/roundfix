---
task: task_04
spec: 0209-sources-that-share-a-context-share-a-spec
status: pending
type: docs
complexity: low
---

# Task 04: The authoring skills teach grouping and extending, and answer a suggestion

## Overview

task_01 makes grouping and extend-before-minting Baseline rules, and task_03
prints grouping suggestions, but the authoring skills still read as one Spec
per source and say nothing about a suggestion. This Task adds one section,
`## Sources that share a context`, to `write-idea`, `write-prd` and
`write-techspec`, placed where the author reads it before choosing what a
Spec adopts (ADR-0208, ADR-0209). It is verifiable on its own: the three
skills and their mirrors carry the section and its phrases, and each raised
version is recorded.

## Requirements

1. MUST add to `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md` and `.agents/skills/write-techspec/SKILL.md` a new second-level section `## Sources that share a context`, placed immediately before `## Process`.
2. Each section MUST state, in these words or longer sentences that contain them: that a Spec may adopt several Inbox Entries, Backlog Entries and Findings whose context is similar or complementary, and that `one Spec per source is neither required nor preferred`; that before minting a Spec the author looks among the open Backlog Entries and unresolved Findings for sources that share the context; that grouping stops at `four implementation Tasks plus its QA gate`, past which the scope splits and each source keeps one owning Spec; and that before minting a new Finding or Backlog Entry the author extends one that fits: `revise an open Backlog Entry in place`, or give an unresolved Finding a `dated addendum`.
3. The sections of `write-prd` and `write-techspec` MUST also state that the author must `answer each suggested pair` that `roundfix spec judge` prints, by adopting the open source within the bound or by stating in the report why it stays apart, that a suggestion never gates, and that no suggestion does not mean that nothing fits. The section of `write-techspec` MUST state that a refactor or bug fix answering a Finding or Backlog Entry looks for the others that share its context before it mints its minimal PRD.
4. MUST change no other text of the three skills, of their references, or of any other skill, and MUST add no text inside the `### QA settlement` section of any skill.
5. MUST raise both version fields of each of the three skills by one patch step from the value on this Task's starting commit, then run `make skills-sync`, record the versions with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, run `make baseline-digests`, and name in the Result every file those commands rewrote.

## Subtasks

- [ ] Add the section to `write-idea`.
- [ ] Add the section, with the suggestion answer, to `write-prd` and `write-techspec`.
- [ ] Raise the three versions, sync the mirrors and record the versions.

## Acceptance Criteria

- [ ] The three skills carry `## Sources that share a context` immediately before `## Process`, with the required phrases.
- [ ] `write-prd` and `write-techspec` tell the author how to answer a suggested pair.
- [ ] Each mirror equals its canonical copy, and `TestEveryOwnedSkillVersionIsRecorded` passes with the three new versions recorded.

## Context

- instruction: `docs/adr/0208-work-items-that-share-a-context-share-a-spec.md`
- instruction: `docs/adr/0209-the-judge-suggests-which-sources-share-a-spec.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `for skill in write-idea write-prd write-techspec; do file=".agents/skills/$skill/SKILL.md"; awk '$0 == "## Sources that share a context" { s = NR; next } s && n == "" && substr($0, 1, 3) == "## " { n = $0 } END { exit (s > 0 && n == "## Process" ? 0 : 1) }' "$file" || { printf 'section missing or after ## Process in %s\n' "$file" >&2; exit 1; }; for phrase in "one Spec per source is neither required nor preferred" "four implementation Tasks plus its QA gate" "revise an open Backlog Entry in place" "dated addendum"; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; diff -r ".agents/skills/$skill" "skills/$skill" >/dev/null || { printf 'mirror differs: skills/%s\n' "$skill" >&2; exit 1; }; done; for skill in write-prd write-techspec; do tr -s '[:space:]' ' ' < ".agents/skills/$skill/SKILL.md" | grep -qF -- "answer each suggested pair" || { printf 'missing suggestion answer in %s\n' "$skill" >&2; exit 1; }; done; out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync|TestWritePRDProjectConstraints|TestWriteTechSpecProjectConstraints)$" ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync TestWritePRDProjectConstraints TestWriteTechSpecProjectConstraints; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task no skill carries `## Sources that share a context`, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1, 2 and 3; User Stories 1, 2, 3 and 4; Core Feature 5; Success Metric 6
- [_techspec.md](_techspec.md) — Exact clause texts; Integration Points; Testing Approach 5; Build Order 4
- ADR-0208; ADR-0209; ADR-0189
