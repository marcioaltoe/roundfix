---
task: task_04
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
status: completed
type: test
complexity: medium
---

# Task 04: A skill-contract test holds the step generic and every owned skill free of personal names

## Overview

The step stays generic only while nothing re-adds one user's tool to a skill
that ships to every adopter. A new skill-contract test holds the five
surfaces, canonical and mirror, to the generic paragraph and their own
events, and fails when any file of any owned skill names the assistant, the
tracker, the channel or the block of the adopted Backlog Entry of 2026-11-03
([references](references/2026-11-03-the-spec-skills-announce-each-spec-event.md)),
reported by jarvis.

## Requirements

1. MUST create `skills/spec_event_announcement_test.go` in package
   `skills_test` with the names of `_techspec.md` → Interfaces and the two
   tests of API Contract 9: `TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines`
   lists every sentence of API Contracts 1 to 6 as its expected text, and
   `TestNoOwnedSkillNamesAPersonalAssistant` reads the owned skills from the
   `skills` keys of `skills/testdata/owned-skill-versions.json` and walks
   every file of each under `.agents/skills/` and `skills/`.
2. MUST call `t.Parallel()` as the first statement of each top-level test and
   subtest, start no process, read nothing outside this repository, and
   report the repository-relative path with each failure.
3. MUST NOT edit `skills/roundfix_modes_test.go`, which Spec 0343 may change,
   or any skill file.
4. MUST prove each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy, never in the checkout, and record in the Task Result the
   failing test and message for each.

## Subtasks

- [ ] Write the two tests with the contract sentences and the forbidden words.
- [ ] Characterize both tests and run the four sabotages on temporary copies.

## Acceptance Criteria

- [ ] Both tests pass on the tree with the five sections.
- [ ] Removing a contract sentence, duplicating the heading, or adding a forbidden word to any owned skill file makes the named test fail with the path.

## Context

- instruction: `skills/glossary_definition_test.go`
- creates: `skills/spec_event_announcement_test.go`

## Consumers

- The new skill-contract test — consumers: `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -run '^(TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines|TestNoOwnedSkillNamesAPersonalAssistant)$' -v ./skills 2>&1)"; status=$?; printf '%s\n' "$out"; test "$status" -eq 0 || exit 1; for name in TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines TestNoOwnedSkillNamesAPersonalAssistant; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name (" || { printf 'missing pass of %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task no test has either name, so `go test` reports no test run and the pass check fails; afterwards both tests pass by name and the whole skills suite passes.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 3; Success Metric 1; Success Metric 2
- `_techspec.md` → Interfaces; API Contract 9; Invariant 3; Testing Approach; Build Order 4
- ADR-0354

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, tech-writer, evidence-gate.

Added `skills/spec_event_announcement_test.go` in `skills_test`, with the
three interface variables and both named tests from API Contract 9. The
announcement test lists all eight generic sentences and every surface's event
sentences from API Contracts 1–6, requires exactly one heading in each of the
ten canonical/mirror files, and limits sentence matching to that section.
The guard reads the ownership record's `skills` keys and walks every file of
all twelve owned skills in both trees, checking all four forbidden strings
without regard to case. Every top-level test and subtest starts with
`t.Parallel()`. Repository reads use `os.OpenRoot("..")` to prevent escaping
the checkout, and the tests start no process. Failures name repository-relative
paths. No skill file or `skills/roundfix_modes_test.go` was edited.

Focused evidence for the acceptance criteria:

- Both tests pass with the five sections: after `gofmt -w
  skills/spec_event_announcement_test.go`, ran `go test -count=1 -timeout 2m
  -v skills/spec_event_announcement_test.go` from the repository root; exit 0,
  both named top-level tests passed, with ten announcement subtests and
  twenty-four owned-skill subtests. The test file was absent before this
  Task's edit.
- Characterization and sabotage checks ran with `python3
  /private/tmp/roundfix-0344-task04-arur85w7/probe.py`; the harness exited 0
  after checking each expected result and saving every full transcript.
  Each case extracted its own `git archive HEAD` copy of audited revision
  `762e8938face672457cd98ed9d76f375595de10d`, then copied in the new test
  file. The harness ran `go test -count=1 -timeout 2m -v
  spec_event_announcement_test.go` from that copy's `skills/` directory.
  The unchanged positive copy passed both tests. A separate characterization
  copy with the sections removed failed the announcement test for all ten
  paths with `has 0 Spec Event Announcement headings, want 1`; the guard
  passed. All mutations stayed in temporary copies.

Each required sabotage independently exited 1, with the other top-level test
passing:

| Sabotage | Failing test/subtest | Observed diagnostic |
| --- | --- | --- |
| Remove the Run prohibition from canonical archive skill | `TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines/.agents/skills/archive-spec/SKILL.md` | `.agents/skills/archive-spec/SKILL.md is missing sentence "An Agent inside a Daemon-assigned Run never announces."` |
| Append `Send it to Jarvis.` to canonical PRD template | `TestNoOwnedSkillNamesAPersonalAssistant/.agents/skills/write-prd` | `.agents/skills/write-prd/references/prd-template.md contains forbidden word "jarvis"` |
| Append `Track it in Todoist.` to mirror QA skill | `TestNoOwnedSkillNamesAPersonalAssistant/skills/qa-gate` | `skills/qa-gate/SKILL.md contains forbidden word "todoist"` |
| Append a second heading to mirror implement reference | `TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines/skills/roundfix/references/implement.md` | `skills/roundfix/references/implement.md has 2 Spec Event Announcement headings, want 1` |

Local transcripts: `/private/tmp/roundfix-0344-task04-arur85w7/positive.log`,
`no-sections.log`, `missing-sentence.log`, `canonical-template-word.log`,
`mirror-other-skill-word.log`, and `duplicate-heading.log` in that directory.
The audited test file's SHA-256 was
`6caa9300324e843e64a2d9379801045f6894a4ce0ef0ea76a6a59356fe534f80`.

The declared Verification and repository Verification were not run; the
Daemon owns them and Task settlement. No follow-up scope was discovered.
