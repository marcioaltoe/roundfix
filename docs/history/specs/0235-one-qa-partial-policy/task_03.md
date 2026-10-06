---
task: task_03
spec: 0235-one-qa-partial-policy
status: completed
type: docs
complexity: medium
---

# Task 03: The skills and user guides state the one partial policy and the network-denied marker

## Overview

The QA gate must write the marker that task_01 recognizes, and every reader
of the policy must find the same rule. This Task updates the qa-gate,
archive-spec and Roundfix skills and the user guides. It changes the
qualifying-partial row of the `### QA settlement` table identically in all
three skills. It then syncs the mirrors and records the raised skill
versions. Each copy of that table is identical by contract, and this Spec
owns the policy the row states.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-06
   ([settlement refuses a partial that archive accepts](references/2026-10-06-settlement-refuses-a-partial-that-archive-accepts.md))
   by setting the Settles cell of the "qualifying declared `partial`" row of
   `### QA settlement` to the text in the TechSpec's "Exact texts", in
   `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/archive-spec/SKILL.md`
   and `.agents/skills/roundfix/SKILL.md`. Keep the row label, the
   Archives cell and every other row byte-identical.
2. MUST tell the gate in `.agents/skills/qa-gate/SKILL.md`, outside
   `### QA settlement`, to record an outside-evidence row that is blocked
   only because the Run sandbox denied network access as
   `blocked (environment: network denied: <host>)`, naming the host the
   lookup tried, with an `outside-evidence row` item kept in its provenance,
   counted in `rows_blocked_environment`. Say that such a row records that
   the source was not reached and never decides a qualifying partial, and
   that any other blocked outside-evidence row still blocks Pull Request
   preparation. Rewrite the outside-evidence paragraph, the sandbox
   paragraph of "6. Close the report" and the outside-evidence decision
   example to match. Add one sentence saying that the Pull Request row's
   provenance item may carry a note after `Pull Request row`. Add one
   sentence stating the single verdict rule of ADR-0240: a Pull Request row
   with equivalent evidence for every control still allows `pass`, and
   without it the report closes `partial`, which qualifies when no other
   unmet row remains.
3. MUST state the same policy in `.agents/skills/roundfix/references/archive.md`
   and `.agents/skills/roundfix/references/settle.md`, replacing each
   sentence that names only the Pull Request row as exempt. MUST name the
   settle refusal's `(report <path>)` suffix in the settle reference.
4. MUST state the same policy, quoting the marker
   `blocked (environment: network denied: <host>)`, in
   `docs/user-guide/commands/qa-report.md`,
   `docs/user-guide/commands/archive.md`, `docs/user-guide/commands/settle.md`
   and `docs/user-guide/context-driven-development.md`, replacing each
   sentence that names only the Pull Request row as exempt.
5. MUST run `make skills-sync` so that every mirror equals its canonical
   file. MUST then re-record the owned skills' versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises both version fields of each changed `SKILL.md` to a free
   version. Running that record command is an implementation step of this
   Task, not part of its Verification: run it after the last skill edit and
   never write a digest or version by hand.
6. MUST NOT change the `pass` rule, the other rows of `### QA settlement`,
   any Baseline asset, `docs/agents/`, `CONTEXT.md` or `CHANGELOG.md`.

## Subtasks

- [ ] Change the qualifying-partial row in the three skills.
- [ ] Teach the qa-gate skill the marker, the note and the verdict rule.
- [ ] Update the Roundfix Skill's archive and settle references.
- [ ] Update the four user guides.
- [ ] Sync the mirrors and record the versions.

## Acceptance Criteria

- [ ] The three skills carry one identical `### QA settlement` section whose
      qualifying row names both exempt rows.
- [ ] The qa-gate skill tells the gate exactly how to write a network-denied
      outside-evidence row and which outside-evidence rows still block.
- [ ] No skill, reference or user guide still says that only the Pull Request
      row is exempt.
- [ ] Every mirror equals its canonical file, and every raised version is
      recorded.

## Context

- instruction: `docs/adr/0240-one-qa-partial-policy-and-rows-a-run-sandbox-cannot-reach.md`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/qa-report.md`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/settle.md`
- interface: `docs/user-guide/context-driven-development.md`

## Verification

- `for file in .agents/skills/qa-gate/SKILL.md .agents/skills/archive-spec/SKILL.md .agents/skills/roundfix/SKILL.md .agents/skills/roundfix/references/archive.md .agents/skills/roundfix/references/settle.md docs/user-guide/commands/qa-report.md docs/user-guide/commands/archive.md docs/user-guide/commands/settle.md docs/user-guide/context-driven-development.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'blocked (environment: network denied: <host>)' || { printf 'missing marker in %s\n' "$file" >&2; exit 1; }; done; for file in .agents/skills/qa-gate/SKILL.md .agents/skills/archive-spec/SKILL.md .agents/skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'a partial whose only unmet rows are such rows qualifies' || { printf 'missing table text in %s\n' "$file" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/settle.md | grep -qF -- '(report <path>)' || { printf 'missing settle suffix in the settle reference\n' >&2; exit 1; }; for pair in qa-gate/SKILL.md archive-spec/SKILL.md roundfix/SKILL.md roundfix/references/archive.md roundfix/references/settle.md; do cmp ".agents/skills/$pair" "skills/$pair" || exit 1; done; out="$(go test -count=1 -v -run '^(TestSettlementGuidanceIsOneTable|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSettlementGuidanceIsOneTable TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task no skill or guide quotes the network-denied marker, so the command fails; after it every file quotes it, the table is identical in the three skills, the mirrors equal their canonical files and the raised versions are recorded.

## References

- `_prd.md` → Goals; Core Features 3; Core Features 6; Success Metric 5
- `_techspec.md` → Exact texts; Vocabulary Contract; Invariants 1, 2 and 6; API Contract 4; Build Order 3
- ADR-0240; ADR-0189; ADR-0233; ADR-0104

## Result

Implementation evidence for daemon handoff; declared Verification remains for
the Daemon.

- Updated the qualifying declared `partial` settlement row identically in the
  three owned canonical skills, naming both the pre-PR Pull Request row and the
  network-denied outside-evidence row.
- Updated the qa-gate instructions for the exact network-denied marker, host
  naming, outside-evidence provenance, source-not-reached meaning, Pull Request
  provenance notes, and the single verdict rule. Updated the Roundfix archive
  and settle references and all four requested user guides with the same policy;
  the settle reference and guide name the `(report <path>)` suffix.
- `make skills-sync` — passed; canonical skill mirrors were synchronized.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` — passed after rerunning with the required filesystem
  access; raised versions are recorded in both skill frontmatter fields and
  `skills/testdata/owned-skill-versions.json`.
- Focused checks — `git diff --check` passed; SHA-256 comparisons confirmed
  every canonical skill and mirror pair matches; a targeted text sweep found no
  remaining Pull Request-only exemption wording in the affected surfaces.

Acceptance evidence:

- The three settlement tables contain the TechSpec exact row text, and their
  canonical/mirror content is synchronized.
- The qa-gate contains the marker, named host and `outside-evidence row`
  provenance instructions, states that the source was not reached, preserves
  the blocking rule for other outside-evidence rows, and states the verdict
  rule.
- The affected references and user guides contain the network-denied marker and
  the two-row exemption policy; settle guidance contains `(report <path>)`.
- Version recording changed only the three owned skills and the generated
  owned-version registry within this Task's declared documentation slice.

## Carry-forward provenance

- Source Run: `run_20261006T120012Z_3cae1c6ab9ca3edb`
- Source commit: `a8dd1be06fefb33043194eb7f4901d6aac0e15f3`
