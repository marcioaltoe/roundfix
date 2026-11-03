---
task: task_04
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
status: pending
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
