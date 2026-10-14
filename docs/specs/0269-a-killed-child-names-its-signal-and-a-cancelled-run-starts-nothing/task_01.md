---
task: task_01
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
status: completed
type: docs
complexity: low
---

# Task 01: The deliver reference, the user guide and the glossary name a Signal Exit

## Overview

task_02 changes what the Delivery Queue prints when the item binary's
`migrate --check` or a step child is ended by a signal. This Task brings the
`roundfix` skill's `deliver` reference, the `deliver` user guide and
`CONTEXT.md` in line, and opens the Backlog Entry that carries the trace of
F6's sender, which this Spec leaves out by the maintainer's decision. It is
verifiable by phrase checks, the mirror check, the skill version record and
the Behavior Surface Record.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [a test child is killed by a signal from an unknown sender](references/2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md),
   for the documentation half: the item-binary notice and the step failure
   follow `_techspec.md` → API Contract 1 and API Contract 2.
2. MUST change, in `.agents/skills/roundfix/references/deliver.md` and in
   `docs/user-guide/commands/deliver.md`, the sentence before the fallback
   notice to cover a probe that "exits with any other code or is killed by a
   signal", replace `exited <n>` in the notice's `text` block with
   `<outcome>`, and add after the block one paragraph that defines
   `<outcome>` as `exited <n>` or, for a Signal Exit,
   `was killed by signal <name> (<number>), core dumped` or
   `was killed by signal <name> (<number>), no core dump`, and that states a
   step child ended by a signal fails the step with
   `<operation> was killed by signal <name> (<number>)` and its core state.
   No text inside any `### QA settlement` section may change.
3. MUST run `make skills-sync`, raise and record the skill's version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`;
   no record or mirror is hand-edited.
4. MUST add **Signal Exit** to `CONTEXT.md` through `domain-modeling`, next to
   **Proven Signal**: the end of a child process that a signal caused instead
   of an exit status, which Roundfix and its test helpers report by the
   signal's name and number and whether a core was dumped, never as exit code
   `-1` or `1`; `_Avoid_: exit -1, crash, negative exit code`. It links no
   Spec (`docs/agents/docs-layout.md`).
5. MUST create `docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md`
   with the `type: fix` Backlog Operational Contract of
   `docs/agents/docs-layout.md`, `status: open`, `created: 2026-10-13` (the
   maintainer's decision date), `spec: null` and `reason: null`. Its Symptom
   names the F6 test and its three occurrences; its Where names
   `internal/cli` delivery tests that run the item binary; its Expected says
   the sender is found and fixed so ten full suites pass; and its Evidence
   says that since this Spec the failure reports the Signal Exit (the phrase
   `names the signal` and the term `Signal Exit` appear), that every recorded
   occurrence was the `stdout` subtest, the second subtest of its test and
   the one whose build copies the fixture over an existing `bin/roundfix`
   (the phrase `stdout subtest` appears),
   recorded as an observation and not a cause, and that `eslogger` needs
   `sudo`, which only the maintainer can run. It links no Spec folder.

## Subtasks

- [ ] Update the `deliver` reference and the `deliver` user guide.
- [ ] Sync the mirror, record the skill version and re-record the Behavior Surface Record.
- [ ] Add **Signal Exit** to `CONTEXT.md`.
- [ ] Write the Backlog Entry for the sender.

## Acceptance Criteria

- [ ] Both documents state the `<outcome>` placeholder and both Signal Exit forms, and the mirror matches.
- [ ] The skill version is raised and recorded, and the Behavior Surface Record is current.
- [ ] `CONTEXT.md` defines **Signal Exit**.
- [ ] The new Backlog Entry is open and states the observation and the `sudo` limit.

## Consumers

- Owned skill version of `roundfix` and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/agents/docs-layout.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- creates: `docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md`

## Verification

- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "exits with any other code or is killed by a signal" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "exits with any other code or is killed by a signal" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "the item binary's migrate --check <outcome>:" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "the item binary's migrate --check <outcome>:" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "was killed by signal <name> (<number>), core dumped" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "was killed by signal <name> (<number>), core dumped" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/deliver.md | grep -qF -- "was killed by signal <name> (<number>), no core dump" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/deliver.md "was killed by signal <name> (<number>), no core dump" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "exits with any other code or is killed by a signal" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "exits with any other code or is killed by a signal" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "the item binary's migrate --check <outcome>:" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "the item binary's migrate --check <outcome>:" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "was killed by signal <name> (<number>), core dumped" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "was killed by signal <name> (<number>), core dumped" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- "was killed by signal <name> (<number>), no core dump" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/deliver.md "was killed by signal <name> (<number>), no core dump" >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Signal Exit**:' || { printf 'CONTEXT.md does not define **Signal Exit**\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract` — expected: exit 0. Before this Task neither document names the outcome placeholder and the glossary lacks the term, so the command fails.
- `test -f docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf 'missing %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md >&2; exit 1; }; grep -q -- '^type:.fix$' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf 'missing line in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md 'type: fix' >&2; exit 1; }; grep -q -- '^status:.open$' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf 'missing line in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md 'status: open' >&2; exit 1; }; grep -q -- '^created:.2026-10-13$' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf 'missing line in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md 'created: 2026-10-13' >&2; exit 1; }; grep -q -- '^spec:.null$' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf 'missing line in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md 'spec: null' >&2; exit 1; }; grep -q -- '^##.Symptom$' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf 'missing line in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md '## Symptom' >&2; exit 1; }; grep -q -- '^##.Where$' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf 'missing line in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md '## Where' >&2; exit 1; }; grep -q -- '^##.Expected$' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf 'missing line in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md '## Expected' >&2; exit 1; }; grep -q -- '^##.Evidence$' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf 'missing line in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md '## Evidence' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md | grep -qF -- "Signal Exit" || { printf 'missing phrase in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md "Signal Exit" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md | grep -qF -- "names the signal" || { printf 'missing phrase in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md "names the signal" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md | grep -qF -- "stdout subtest" || { printf 'missing phrase in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md "stdout subtest" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md | grep -qF -- "sudo" || { printf 'missing phrase in %s: %s\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md "sudo" >&2; exit 1; }; ! grep -q 'docs/specs/' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md || { printf '%s links a Spec folder\n' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md >&2; exit 1; }` — expected: exit 0. Before this Task the entry does not exist.

## References

- [_prd.md](_prd.md) — Core Feature 5; User Story 2; Glossary; Decisions
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; Invariant 1; Build Order 1
- [2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md](references/2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md)
- ADR-0256; ADR-0189; ADR-0233; ADR-0244; ADR-0266

## Result

The deliver reference and user guide now describe the fallback `<outcome>`
placeholder, both Signal Exit forms, and signal-ended step failures. The
canonical skill mirror, owned skill version record, and Behavior Surface Record
were regenerated by their repository-owned commands. `CONTEXT.md` now defines
Signal Exit beside Proven Signal, and the open Backlog Entry records the three
F6 observations, the stdout subtest, and the sudo-only sender trace.

Focused checks run after the edits:

- `git diff --check` — passed.
- `cmp -s .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md` — passed after `make skills-sync`.
- The owned skill version and Behavior Surface Record commands completed successfully and wrote their generated records.

Acceptance evidence:

- The canonical deliver reference and user guide contain `<outcome>`, both Signal Exit forms, and the step-child failure wording; the mirror matches the canonical skill reference.
- The owned skill version and Behavior Surface Record were regenerated by their record commands.
- `CONTEXT.md` defines Signal Exit adjacent to Proven Signal and excludes exit codes `-1` and `1`.
- The new open Backlog Entry has the required operational front matter, F6 observation, `stdout subtest` detail, and maintainer-only `sudo` limit, with no Spec-folder link.

Verification feedback repair:

- The Daemon diagnostic identified the literal phrase check for `stdout subtest`; the Backlog Entry now uses that exact unquoted phrase.
- Focused check `rg -n --fixed-strings 'stdout subtest' docs/backlog/2026-10-13-trace-the-sender-of-the-f6-signal.md` — passed.
