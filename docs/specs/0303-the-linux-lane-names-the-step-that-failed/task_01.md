---
task: task_01
spec: 0303-the-linux-lane-names-the-step-that-failed
status: completed
type: docs
complexity: low
---

# Task 01: The deliver guide, the skill reference and the glossary name a failed go list

## Overview

The `deliver` user guide and the Roundfix skill's `deliver` reference say
that every Lane Change Failure has the cause `go test exited <n>`, and the
glossary attributes a Lane Change Failure to the candidate's Go test run only.
This Task describes the lane script's `go list` branch, its exit `4`, its
three causes and their Park Classes in both documents, raises the Roundfix
skill's version, revises **Lane Change Failure** in the glossary, and
re-records the Behavior Surface Record last (`_prd.md` → Core Feature 3).

## Requirements

1. MUST make these edits in `docs/user-guide/commands/deliver.md` and, with
   the same text, in the canonical Roundfix reference
   `.agents/skills/roundfix/references/deliver.md`, inside the Linux Lane
   section; no text outside it changes:
   - in the Lane Container Failure cause list, right after the bullet that
     begins "Exit `3`:" and ends "`cannot unpack the candidate` or `cannot
     record the candidate`.", add the bullet "- Exit `4` whose last
     `roundfix-linux-lane: go list: exited <n>` output line names `<n>` of
     `128` or more: `go list exited <n>: the process ended by signal
     <n-128>`.";
   - replace the sentence "It parks in `finding`, with cause `go test exited
     <n>`." with "It parks in `finding`. The lane script selects the packages
     with `go list` before `go test`; when `go list` exits non-zero, the
     script prints `roundfix-linux-lane: go list: exited <n>` and exits `4`.
     Exit `4` has the cause `go list exited <n>` from the last such line
     whose `<n>` is below `128`, or `the lane script exited 4 without
     reporting the go list exit code` when no such line exists. Every other
     exit has the cause `go test exited <n>`." Neither document cites ADRs
     today, and this edit adds no ADR citation;
   - in the code block that follows ("The gate reason and parked blocker are
     respectively:"), replace `go test exited <n>` with `<cause>` on both
     lines;
   - in the lane log's outcome block, replace the line `Linux Lane change
     failure: go test exited <n>` with `Linux Lane change failure: <cause>`;
   - in the `roundfix deliver lane` output block, replace the line
     `roundfix: Linux Lane change failure: go test exited <n>; see <lane log>`
     with `roundfix: Linux Lane change failure: <cause>; see <lane log>`.
   No text inside any `### QA settlement` section of any skill may change.
2. MUST revise **Lane Change Failure** in `CONTEXT.md` through
   `domain-modeling`, replacing the first sentence of its definition, "A
   Linux Lane failure attributed to the candidate's Go test run.", with "A
   Linux Lane failure attributed to the candidate: its Go test run, or the
   `go list` run that selects its packages exiting below 128.", and adding a
   second ADR link right after the entry's existing ADR-0294 link, in the
   glossary's link form, to
   `docs/adr/0316-the-linux-lane-names-a-failed-go-list-apart-from-go-test.md`
   (ADR-0316). Its `_Avoid_` line and every other entry, **Lane Container
   Failure** included, keep their bytes.
3. MUST run these commands, because this Task edits an owned skill and a
   guide. Before them, raise both version fields (`metadata.version` and the
   top-level `version`) of the canonical `SKILL.md` of the Roundfix skill by
   one patch over the value the tree holds when this Task starts, because the
   record command refuses changed content under an unchanged version. The
   Agent runs them itself before handing back, because the Daemon never
   regenerates a record; no mirror or record is hand-edited. Re-record the
   Behavior Surface Record after the last guide edit: it is the last command
   of this Task, and no guide or reference changes after it.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     `deliver` guide's and reference's fingerprints change.
4. MUST search the guides, the agent guides and the skill copies for other
   wording that says every Lane Change Failure is `go test`'s, and report
   what it found.
5. MUST NOT change any Go source or ADR, `docs/references/skill-coverage.json`
   or any other guide page.

## Subtasks

- [ ] Describe the `go list` branch, its causes and the generalized outcome lines in the guide and the reference.
- [ ] Revise **Lane Change Failure** in the glossary.
- [ ] Raise the skill version, sync the mirrors, record the version and, last, re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The guide and both copies of the reference name the `go list` marker, exit `4`, the change and signal causes, and `<cause>` in the outcome lines.
- [ ] `CONTEXT.md` revises **Lane Change Failure** and links ADR-0316.
- [ ] The skill version is raised and recorded, the mirrors are in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirrors, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `docs/adr/0316-the-linux-lane-names-a-failed-go-list-apart-from-go-test.md`
- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do n="$(tr -s '[:space:]' ' ' < "$f")"; printf '%s' "$n" | grep -qF -- 'roundfix-linux-lane: go list: exited <n>' || { printf 'missing go list marker in %s\n' "$f" >&2; exit 1; }; printf '%s' "$n" | grep -qF -- 'go list exited <n>: the process ended by signal <n-128>' || { printf 'missing go list signal cause in %s\n' "$f" >&2; exit 1; }; printf '%s' "$n" | grep -qF -- 'the lane script exited 4 without reporting the go list exit code' || { printf 'missing go list no-line cause in %s\n' "$f" >&2; exit 1; }; printf '%s' "$n" | grep -qF -- 'roundfix: Linux Lane change failure: <cause>; see <lane log>' || { printf 'missing generalized deliver lane line in %s\n' "$f" >&2; exit 1; }; printf '%s' "$n" | grep -qF -- 'gate-failed: Linux Lane change failure: <cause>; see <lane log>' || { printf 'missing generalized blocker in %s\n' "$f" >&2; exit 1; }; if printf '%s' "$n" | grep -q -- 'It parks in .finding., with cause'; then printf 'stale change failure sentence in %s\n' "$f" >&2; exit 1; fi; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -q -- 'its Go test run, or the .go list. run that selects its packages exiting below 128' || { echo 'missing Lane Change Failure revision' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '0316-the-linux-lane-names-a-failed-go-list-apart-from-go-test.md' || { echo 'missing ADR-0316 link in the glossary' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document names the `go list` marker, so the command fails at its first phrase check. Afterwards the guide and both reference copies name the marker, exit `4`, the signal and no-line causes and the generalized outcome lines, the glossary revises **Lane Change Failure** with ADR-0316, the mirrors are in sync, the skill version is recorded, and the documentation contracts, the Skill Coverage Map among them, pass.

## References

- `_prd.md` → Core Feature 3; Goals; Success Metric 6; Glossary
- `_techspec.md` → API Contract 5; Build Order 1
- ADR-0316

## Result

- Updated the Linux Lane sections in `docs/user-guide/commands/deliver.md` and
  `.agents/skills/roundfix/references/deliver.md` with the Exit `4` `go list`
  marker, signal cause, no-line cause, and generalized `<cause>` outcome
  lines. `make skills-sync` propagated the canonical reference and raised
  skill version `0.1.132` to `0.1.133` in both skill copies.
- Revised the `Lane Change Failure` glossary definition in `CONTEXT.md` to
  include the candidate's package-selecting `go list` run and linked ADR-0316
  after ADR-0294.
- Wording sweep: outside the edited guide/reference surfaces, the old
  `go test` wording remains only in ADR-0294 and ADR-0316 as historical or
  unchanged-outcome context; no other guide, agent guide, or skill copy has a
  claim that every Lane Change Failure belongs to `go test`.
- Focused evidence: `git diff --check` passed before regeneration; the three
  target references were inspected for the new marker and causes; and the
  canonical and mirror reference files were synchronized by `make skills-sync`.
- Required record commands passed:
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` (1 passed) and
  `go test -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  (1 passed). The latter was run after the final guide edit.
- Daemon-owned Verification was not run, and Task status remains unchanged.
