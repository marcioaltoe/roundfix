---
task: task_06
spec: 0235-one-qa-partial-policy
status: completed
type: docs
complexity: low
---

# Task 06: The archived baseline names every accepted partial, and the settle reference gives the suffix to settle

## Overview

The first QA gate of this Spec failed two rows on documentation, not on
behavior. Its archived-report replay found that the PRD's baseline is wrong:
it says archived Spec 0179 is the only accepted archived `partial`, but the
newest reports of 0073 and 0079 are accepted too. Its guidance audit found
that the Roundfix settle reference gives the `(report <path>)` suffix to
`roundfix qa-report accept`, while only the Settle Command prints it. This
Task corrects both statements, syncs the skill mirror and records the raised
Roundfix Skill version. It changes no code and no archived report.

## Requirements

1. MUST correct the archived-report baseline in `_prd.md` → Acceptance
   evidence to the measured truth. Replace the clause
   "and 0179 is the only one accepted." with
   "and 0073, 0079 and 0179, whose only unmet rows are declared rows their
   Specs cover, are accepted. Every other archived `partial` report is
   refused." Keep every other sentence of that bullet byte-identical. The
   measurement, repeated during authoring on 2026-10-06 with the shipped
   binary (v0.46.0, `08660fcc`), is that `roundfix qa-report accept` exits 0
   on `docs/history/specs/0073-skill-versions-decoupled-from-the-binary/qa/qa-report-2026-08-06-01.md`
   (one declared row), on
   `docs/history/specs/0079-one-door-for-fleet-knowledge/qa/qa-report-2026-08-06-04.md`
   (five declared rows) and on the 0179 report
   `qa-report-2026-09-29-04.md` (four declared rows). Each has
   `rows_blocked_environment: 0` and `rows_blocked_finding: 0`, so the
   pre-change policy already accepted it. The QA gate's replay of all 28
   newest archived `partial` reports refused every other one.
2. MUST NOT change `_techspec.md`, the other task files or
   `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
   for this claim. None of them states it: ADR-0240 says only that archived
   reports keep reading as before, and `task_05.md` Requirement 2 compares
   against the PRD's measurement, so the corrected PRD is the one baseline
   the gate replays against.
3. MUST correct the last sentence of `## QA Report acceptance` in
   `.agents/skills/roundfix/references/settle.md`. Replace
   "When this eligibility check refuses, the stderr reason ends with
   `(report <path>)` naming the report it judged." with
   "When this eligibility check refuses, `qa-report accept` prints the reason
   alone. Only the Settle Command appends `(report <path>)` to its refusal,
   naming the report it judged." A fresh `roundfix qa-report accept` refusal
   prints `roundfix: QA Report eligibility refused: rows_blocked_environment
   is 1; expected 0` with no suffix, and `docs/user-guide/commands/settle.md`
   already gives the suffix to settle. Keep every other line of the reference
   byte-identical.
4. MUST run `make skills-sync` so that `skills/roundfix/references/settle.md`
   equals its canonical file, and MUST then run
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   That record command raises both version fields of the Roundfix Skill's
   `SKILL.md` to a free version and records it in
   `skills/testdata/owned-skill-versions.json`. Running it is an
   implementation step of this Task, not only a Verification. Run
   `make skills-sync` again if the record command changed only the canonical
   `SKILL.md`, so both copies stay equal. Never write a digest or a version by
   hand.
5. MUST NOT change the `### QA settlement` section of any skill, any other
   skill, any archived Spec or QA Report, `CONTEXT.md`, `CHANGELOG.md` or any
   Go source. The two corrections introduce no domain term: Settle Command is
   already in `CONTEXT.md`.

## Subtasks

- [ ] Correct the archived baseline sentence in the PRD.
- [ ] Correct the suffix sentence in the settle reference.
- [ ] Sync the mirror and record the raised Roundfix Skill version.

## Acceptance Criteria

- [ ] The PRD names 0073, 0079 and 0179 as the accepted archived partials and
      no Spec artifact or ADR-0240 says 0179 is the only one.
- [ ] Both copies of the settle reference give the `(report <path>)` suffix
      to the Settle Command only, and the mirrors equal their canonical
      files.
- [ ] The raised Roundfix Skill version is recorded.

## Context

- instruction: `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- interface: `docs/specs/0235-one-qa-partial-policy/_prd.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `prd=docs/specs/0235-one-qa-partial-policy/_prd.md; tr -s '[:space:]' ' ' < "$prd" | grep -qF -- '0073, 0079 and 0179, whose only unmet rows are declared rows their Specs cover, are accepted' || { printf 'missing corrected archived baseline in the PRD\n' >&2; exit 1; }; for file in "$prd" docs/specs/0235-one-qa-partial-policy/_techspec.md docs/specs/0235-one-qa-partial-policy/task_05.md docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md; do if tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'is the only one accepted'; then printf 'stale archived baseline claim in %s\n' "$file" >&2; exit 1; fi; done; for file in .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'prints the reason alone. Only the Settle Command appends' || { printf 'missing corrected suffix sentence in %s\n' "$file" >&2; exit 1; }; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'When this eligibility check refuses, the stderr reason ends with'; then printf 'stale suffix sentence in %s\n' "$file" >&2; exit 1; fi; done; for pair in roundfix/SKILL.md roundfix/references/settle.md; do cmp ".agents/skills/$pair" "skills/$pair" || exit 1; done; out="$(go test -count=1 -v -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestEveryOwnedSkillVersionIsRecorded ' || { printf 'missing pass: TestEveryOwnedSkillVersionIsRecorded\n' >&2; exit 1; }` — expected: exit 0; before this Task the PRD says 0179 is the only accepted archived partial and the settle reference gives the suffix to `qa-report accept`, so the command fails; after it the PRD names the three accepted partials, both copies of the reference give the suffix to settle only, the mirrors match and the raised Roundfix Skill version is recorded. Editing the reference without raising the version fails the record test.

## References

- `_prd.md` → Acceptance evidence; Success Metric 1
- `_techspec.md` → API Contract 4; Surface Transcript 3
- ADR-0240; ADR-0189

## Result

Implementation evidence:

- Updated the PRD's measured archived `partial` baseline to accept 0073,
  0079 and 0179, with every other archived `partial` refused.
- Updated the canonical Roundfix settle reference, synchronized its mirror,
  and recorded Roundfix Skill version 0.1.39 with the generated digest.

Focused checks:

- `make skills-sync` — passed before and after owned-skill version recording.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` — reached the test after switching to the
  task-scoped Go cache and passed with exit 0. The initial host-cache attempt
  was blocked by cache permissions.
- `cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` and the
  equivalent settle-reference comparison — passed; both mirrors are equal.
- Targeted text inspection confirmed the corrected PRD baseline and the
  corrected `qa-report accept` / Settle Command suffix wording. No stale
  `0179 is the only one accepted` claim remains in the declared Spec and ADR
  scope.

Acceptance evidence:

- The PRD acceptance criterion is supported by the corrected measured clause
  and the targeted text inspection; `_techspec.md`, `task_05.md` and
  ADR-0240 were left unchanged.
- The settle-reference criterion is supported by both `cmp` checks and the
  targeted wording inspection; the suffix is assigned to Settle Command only.
- The raised-version criterion is supported by the recorder test passing and
  the generated 0.1.39 entry in `skills/testdata/owned-skill-versions.json`.

## Carry-forward provenance

- Source Run: `run_20261006T131510Z_7cff0c6eeb3ab564`
- Source commit: `c3e658a0c83ba081e3d64346b096dd5bda6a09ba`
