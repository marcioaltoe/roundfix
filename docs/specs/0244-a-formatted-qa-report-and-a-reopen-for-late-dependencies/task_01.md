---
task: task_01
spec: 0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies
status: completed
type: docs
complexity: low
---

# Task 01: The glossary, the guides and the Roundfix Skill describe the Format Command and the Late Dependency

## Overview

Describes the behavior that task_03 and task_04 build. The places are
`CONTEXT.md`, the configuration guide, the Spec workflow guide's evidence
snapshot passage, the reopen command reference, and the Roundfix Skill's
`implement` and `settle` references with their mirrors, so the skill ships
with the CLI behavior. This Task answers the Backlog Entry "An unformatted QA
report breaks the next Run's precondition" of 2026-10-06. It is verifiable on
its own by the phrases, the mirrors and the recorded skill version.

## Requirements

1. MUST add to `CONTEXT.md`, following the `domain-modeling` skill, the terms
   **Format Command** and **Late Dependency** in the glossary's `**Term**:` /
   definition / `_Avoid_:` shape. **Format Command** is the Project Config
   command `verification.format` the Daemon runs, outside the Agent sandbox,
   over the files under a Spec's `qa/` directory before the repository
   Verification precondition (an imported pass) and before the QA Report commit,
   restoring the original bytes when it fails (ADR-0249). **Late Dependency** is
   a Task in a completed QA gate's dependency closure that the Task Graph did
   not hold when the newest QA Report was committed; `roundfix reopen` returns
   the gate to `pending` over it whatever its status (ADR-0249).
2. MUST document `verification.format` in `docs/user-guide/configuration.md`:
   its default (empty, nothing runs), that Project Config replaces User Config,
   that each file is appended as its own argument and the command runs from the
   Run Worktree root, the two run points, that only regular files under the
   Spec's `qa/` directory are formatted, the 5-minute cap, the revert, and the
   stderr line starting `roundfix: QA format`. Add the key to the reference
   table beside the other `verification.*` keys and give stack examples that
   accept file arguments (for example `gofmt -w`, a Prettier `--write
   --ignore-unknown` invocation, `dprint fmt`), stating that a formatter which
   rejects unknown file types needs its own ignore option. The text describes
   `_techspec.md` → API Contract 1 and API Contract 2.
3. MUST describe in `docs/user-guide/context-driven-development.md`, in the
   evidence snapshot passage that names the `prior_report` phase, that an
   imported pass is copied and proven byte for byte and then formatted by the
   Format Command before the repository Verification precondition, and the
   `daemon.qa` phase `format` with its `stage` and `outcome` values
   (`_techspec.md` → Invariant 8).
4. MUST update `docs/user-guide/commands/reopen.md` and
   `.agents/skills/roundfix/references/settle.md` so each states that reopen
   also returns a completed gate to `pending` over a Late Dependency proven from
   the Task Graph manifest at the commit that added the newest QA Report, that
   the record line is `Dependencies added after the QA Report`, and that without
   that commit reopen refuses as before (`_techspec.md` → API Contract 3). Each
   file must contain the exact phrase "Late Dependency".
5. MUST add to `.agents/skills/roundfix/references/implement.md`, beside the
   QA Report commit boundary, that a configured `verification.format` runs over
   the QA directory before the QA Report commit and over an imported pass before
   the precondition. The file must contain the exact phrase
   "verification.format".
6. MUST raise the Roundfix Skill's version in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md` from the value it holds when this Task
   starts, to the next free version at record time, run `make skills-sync`, and
   record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Running those commands is an implementation step, never part of
   Verification, and no version is written by hand.
7. MUST NOT change any other command reference, the skill's other text, the
   `### QA settlement` section of any skill, the Baseline guides under
   `docs/agents/`, or another glossary entry.

## Subtasks

- [x] Add the two glossary terms.
- [x] Document the configuration key, its run points and its revert.
- [x] Describe the import formatting and the `format` phase.
- [x] Describe the Late Dependency reopen in the guide and the skill.
- [x] Raise and record the skill version and sync the mirrors.

## Acceptance Criteria

- [x] `CONTEXT.md` defines **Format Command** and **Late Dependency**.
- [x] The configuration guide documents `verification.format` and its stderr
      line.
- [x] The reopen reference and the skill's `settle` reference describe the Late
      Dependency.
- [x] Every skill mirror equals its canonical file, and the raised version is
      recorded.

## Result

Implemented the documentation contract for the Format Command and Late
Dependency. `CONTEXT.md` now defines both terms; the configuration guide covers
the empty default, Project Config replacement, file arguments, Run Worktree
root, both run points, regular-file scope, five-minute cap, revert behavior,
stderr diagnostic, stack examples, and TechSpec API Contracts 1 and 2. The
workflow guide records byte-for-byte import proof followed by formatting and
the `daemon.qa` `format` stage/outcome. The reopen guide and settle reference
document manifest-based Late Dependency reopening, its invalidation line, and
the refusal without the report-adding commit. The implement reference records
both `verification.format` run points.

Focused implementation checks:

- `make skills-sync` — passed after canonical reference edits and after the
  recorded skill version change.
- `GOCACHE=/tmp/roundfix-0244-gocache go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed;
  the Roundfix Skill is recorded at `0.1.50` in both front-matter fields and
  `skills/testdata/owned-skill-versions.json`.
- `make baseline-digests` — passed; no derived digest changes were needed.
- Focused phrase checks and `cmp` checks — passed for the glossary terms,
  configuration/stderr text, prior-report format phase, Late Dependency
  record line, implement reference, and all three skill mirrors.

The Daemon still owns this Task's status and declared Verification; those were
not run or changed here.

Verification Feedback repair:

- The diagnostic artifact reported `content changed under version 0.1.49;
  raise the version` after the later reference wording edit was synchronized.
- Re-ran the authorized version recorder with the task-scoped `GOCACHE`; it
  passed and raised the Roundfix Skill to `0.1.50` in both front-matter fields,
  with the new digest recorded in `skills/testdata/owned-skill-versions.json`.
- `make skills-sync` passed, and `make baseline-digests` passed with no
  derived changes.
- Focused `rg` phrase checks and `cmp` checks passed after the repair. The
  declared Verification command was not rerun in this feedback turn.

## Context

- interface: `CONTEXT.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/user-guide/commands/reopen.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0249-a-qa-step-formats-its-qa-directory-and-reopen-sees-a-late-dependency.md`
- instruction: `docs/adr/0194-the-daemon-records-what-a-qa-row-observed-and-hands-a-failed-pass-to-the-next.md`

## Verification

- `f() { tr -s '[:space:]' ' ' < "$1" | grep -qF -- "$2" || { printf 'missing in %s: %s\n' "$1" "$2" >&2; exit 1; }; }; f CONTEXT.md '**Format Command**:'; f CONTEXT.md '**Late Dependency**:'; f docs/user-guide/configuration.md 'verification.format'; f docs/user-guide/configuration.md 'roundfix: QA format'; f docs/user-guide/context-driven-development.md 'Format Command'; f docs/user-guide/commands/reopen.md 'Late Dependency'; f docs/user-guide/commands/reopen.md 'Dependencies added after the QA Report'; f .agents/skills/roundfix/references/settle.md 'Late Dependency'; f .agents/skills/roundfix/references/implement.md 'verification.format'; cmp .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md || exit 1; cmp .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md || exit 1; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0; before this Task no file carries the new terms, so the command fails at its first phrase check, and after it the raised version must be recorded for the skill test to pass.

## References

- `_prd.md` → Core Feature 4; Goals
- `_techspec.md` → Build Order 1; API Contract 1; API Contract 2; API Contract 3; Invariant 8
- ADR-0249
- ADR-0194

## Carry-forward provenance

- Source Run: `run_20261007T095139Z_d22b1ce63abf526a`
- Source commit: `597e12202dca52fbfcf2d953fb8ed65c71fb6561`
