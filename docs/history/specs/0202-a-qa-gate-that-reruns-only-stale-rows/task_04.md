---
task: task_04
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
status: completed
type: docs
complexity: medium
---

# Task 04: The qa-gate skill, the QA prompt and the guide teach the carry

## Overview

The gate Agent learns what to execute from the qa-gate skill and the QA
prompt. Today both tell it to rerun the whole matrix and treat input
declarations as optional, and nothing tells it that the Daemon writes the
snapshot. This Task teaches the carry. A carried row is kept and counts as
passed, every executed row declares its inputs, `commit_range` marks rows
that read Task commits, and the Agent never writes `evidence_snapshots`. It
also documents every emitted word of the TechSpec's Vocabulary Contract in the
Context-Driven Development guide.

## Requirements

1. MUST replace the section-1 sentence of `.agents/skills/qa-gate/SKILL.md`
   that starts `On a rerun, start with previously failed or blocked rows` with
   guidance that keeps every row the Daemon seeded as `carried (…)`. The
   guidance MUST say such a row counts as passed and is not executed again,
   and that every other row is executed, starting with the rows the seeded
   `Row carry-forward` section lists as re-run.
2. MUST rewrite the "Row input declaration" section so that every executed row
   declares `inputs:`, and add `commit_range` to its kind table as never
   carriable, for a row that reads Task commits, their authorization or their
   changed-file scope. It MUST tell the Agent to declare conservatively,
   including every source a built binary compiles from and the module
   manifest. It MUST state that the Agent never writes `evidence_snapshots`,
   which the Daemon records after the turn, and that rows naming the
   repository Verification or the Pull Request row are always observed.
3. MUST keep every byte of the `### QA settlement` section unchanged, and put
   any new heading outside it.
4. MUST raise both version fields of the qa-gate skill one patch step above
   the version on the Task's starting tree. It MUST run `make skills-sync`,
   then `make baseline-digests`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and name in the Result every path those commands rewrote.
5. MUST add to `qaGateContract` in `internal/agent/spec_prompt.go` these two
   lines, verbatim:
   `- A seeded row whose status is carried (established by: …; head: …) is not executed again: keep its identifier, status and provenance, and count it as passed. Execute every other row; the seeded Row carry-forward section names why each prior row re-runs.`
   and
   `- Declare inputs on every row you execute: repository_path for repository content, commit_range for a row that reads Task commits, their authorization or changed-file scope. Never write evidence_snapshots; the Daemon records it.`
6. MUST add a section to `docs/user-guide/context-driven-development.md`, after
   the paragraph on how the mechanical stage validates the newest report. It
   explains the Evidence Snapshot, the import of a failed pass, the carry proof
   and always-observed rows. It MUST name every emitted word of the TechSpec's
   Vocabulary Contract, so that `roundfix spec check` finds no undocumented
   token for either declaration.
7. MUST NOT edit `CONTEXT.md`, any other skill, or any existing test.

## Subtasks

- [ ] Rewrite the rerun sentence and the row input declaration in the skill.
- [ ] Raise the version, regenerate the mirror and record the version.
- [ ] Add the two lines to the QA prompt contract with a test.
- [ ] Document the carry and its vocabulary in the guide.

## Acceptance Criteria

- [ ] The canonical skill and its mirror name `commit_range`, `counts as
      passed`, `Row carry-forward` and `never writes
      \`evidence_snapshots\``, and are byte-identical.
- [ ] The `### QA settlement` section is unchanged, and the repository contract
      that pins it passes.
- [ ] The raised version is recorded, and the owned-skill version test passes.
- [ ] A QA prompt built by `BuildQAPrompt` contains both new lines exactly
      once.
- [ ] The guide names every emitted word of the Vocabulary Contract.

## Context

- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `internal/agent/spec_prompt.go`
- interface: `docs/user-guide/context-driven-development.md`
- creates: `internal/agent/qa_carry_prompt_test.go`
- instruction: `docs/adr/0194-the-daemon-records-what-a-qa-row-observed-and-hands-a-failed-pass-to-the-next.md`
- instruction: `docs/adr/0195-rows-that-read-the-gate-or-the-commits-are-observed-on-every-pass.md`

## Verification

- `for pair in ".agents/skills/qa-gate/SKILL.md|commit_range" ".agents/skills/qa-gate/SKILL.md|counts as passed" ".agents/skills/qa-gate/SKILL.md|Row carry-forward" ".agents/skills/qa-gate/SKILL.md|never writes $(printf '\140')evidence_snapshots$(printf '\140')" "skills/qa-gate/SKILL.md|commit_range" "skills/qa-gate/SKILL.md|counts as passed" "skills/qa-gate/SKILL.md|Row carry-forward" "skills/qa-gate/SKILL.md|never writes $(printf '\140')evidence_snapshots$(printf '\140')"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task the skill names none of the four phrases, so the command fails.
- `for phrase in 'evidence_snapshots' 'commit_range' 'Row carry-forward' 're-run: not pass' 'no inputs' 'non-repository input' 'always observed: repository Verification' 'always observed: Pull Request row' 'always observed: commit_range input' 'no evidence snapshot' 'establishing report unavailable' 'establishing head unproven' 'input moved:' 'evidence differs' 'prior_report' 'carried_rows' 'rerun_rows'; do tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- "$phrase" || { printf 'missing phrase in guide: %s\n' "$phrase" >&2; exit 1; }; done` — expected: exit 0; before this Task the guide names none of these words, so the command fails.
- `make skills-sync-check && go test -count=1 ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' && go test -count=1 -tags repocontract -run '^TestSettlementGuidanceIsOneTable$' ./skills && out="$(go test -count=1 -v -run '^TestQAContractKeepsCarriedRowsAndDeclaresInputs$' ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestQAContractKeepsCarriedRowsAndDeclaresInputs" || { printf 'missing pass: TestQAContractKeepsCarriedRowsAndDeclaresInputs\n' >&2; exit 1; }` — expected: exit 0; before this Task the named prompt test does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 1; User Story 4; Core Feature 6
- [_techspec.md](_techspec.md) — System Architecture; Vocabulary Contract; API Contract 2; API Contract 3; API Contract 4; Testing Approach 4; Build Order 4
- ADR-0194; ADR-0195; ADR-0097; ADR-0155

## Result

Implemented the QA carry guidance and prompt contract for this Task. The
canonical and mirrored `qa-gate` skills now teach carried rows as passed and
not re-executed, require conservative `inputs:` declarations including
`commit_range`, identify always-observed rows, and reserve
`evidence_snapshots` for the Daemon. Both skill version fields are `0.0.6`.
The Context-Driven Development guide documents Evidence Snapshot, failed-pass
import, carry proof, Carry Disposition reasons, always-observed rows, and the
`daemon.qa` event vocabulary. Added the focused prompt contract test
`TestQAContractKeepsCarriedRowsAndDeclaresInputs`.

Focused implementation evidence:

- `make skills-sync` succeeded and rewrote `skills/qa-gate/SKILL.md` from the
  canonical `.agents/skills/qa-gate/SKILL.md`.
- `make baseline-digests` succeeded and reported no changed derived paths.
- `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` passed and recorded `0.0.6` in
  `skills/testdata/owned-skill-versions.json`.
- With `GOCACHE="$PWD/.gocache"`,
  `rtk proxy go test -count=1 -v ./internal/agent -run
  '^TestQAContractKeepsCarriedRowsAndDeclaresInputs$'` passed.
- Exact comparison confirms `.agents/skills/qa-gate/SKILL.md` and
  `skills/qa-gate/SKILL.md` are byte-identical. The SHA-256 of the
  `### QA settlement` section through the section before `## 1.` remains
  `7f3a75e01068656ed3b840b011106b35305a55ee8b308bbdd104d8ea2e76128d` in
  both files.
- The guide contains every declared carry vocabulary token, including
  `evidence_snapshots`, `commit_range`, `## Row carry-forward`, all closed
  `re-run:` reasons, `prior_report`, `carried_rows`, and `rerun_rows`.

The initial focused Agent test attempt used the host Go cache and was blocked by
cache permissions; rerunning with the task-scoped `.gocache` passed. The
Daemon-owned `status: in_progress` field was left unchanged. The Task's
declared Verification remains for the Daemon.

## Carry-forward provenance

- Source Run: `run_20261001T133454Z_7983d7d312c8ec5c`
- Source commit: `a5924f26eb040b981c285bb624b7e6151060b1d1`
