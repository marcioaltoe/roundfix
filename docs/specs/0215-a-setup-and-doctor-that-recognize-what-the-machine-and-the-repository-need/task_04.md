---
task: task_04
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
status: pending
type: backend
complexity: medium
---

# Task 04: Setup reports the readiness lines and deliver start refuses when this machine cannot publish

## Overview

Setup is the command a person runs on a new machine, and `deliver start` is
where an unattended queue begins; neither looks at the forge. This Task has
Setup print the five readiness lines Doctor prints (API Contract 2), has
`deliver start` refuse on a `failed` `gh` or `remote` finding before it
creates a queue (API Contract 3, Surface Transcript 3), and ships the Roundfix
Skill update the repository's skill-sync rule requires for task_01 to this
Task. It is verifiable alone through `runCLI` with fake readiness.

## Requirements

1. MUST add the `readiness` field to `setupDependencies`, defaulting to the
   five lines task_01 and task_02 built, and MUST print them, with the same
   text as Doctor, after `acpx` and before the adapter work, on both the path
   where acpx is ready and the path where it is not.
2. MUST offer no install or change for these lines, and MUST make Setup exit
   `1` at its end when one is `failed`, as API Contract 2 states.
3. MUST print no readiness line when `setupDependencies.readiness` is unset,
   so the Setup tests in `internal/cli/cli_test.go`, which this Task does not
   change, keep their output.
4. MUST add the `deliveryReadiness` field to `commandDependencies`,
   defaulting to the `gh` and `remote` lines, and MUST run it in
   `deliver start` after the authorization check and before the Run Database
   is opened.
5. MUST refuse `deliver start` on any `failed` finding with exit `2` and the
   `Preflight failed` output of Surface Transcript 3, one reason line per
   finding, and MUST leave no queue record and start no owner; a `warn`
   finding MUST NOT refuse or print.
6. MUST give each existing `deliver start` test that passes the
   authorization check a ready fake for `deliveryReadiness`, without changing
   what it asserts.
7. MUST describe the Setup lines in `docs/user-guide/commands/setup.md` and
   the refusal in `docs/user-guide/commands/deliver.md`.
8. MUST describe in the Roundfix Skill the five lines, their codes and the
   `warn` status in `.agents/skills/roundfix/references/setup.md`, the
   `deliver start` refusal in `.agents/skills/roundfix/references/deliver.md`,
   and the restore of a trailing upstream skill in
   `.agents/skills/roundfix/references/baseline.md`; raise the skill's version
   by one patch level in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md`; run `make skills-sync`; and re-record
   the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
9. MUST NOT edit `internal/cli/cli_test.go`, the `### QA settlement` section
   of any skill, any other command reference, `commands.md` or `CONTEXT.md`.
10. MUST add the tests named in Verification: Setup in the new file
    `internal/cli/setup_readiness_test.go`, and Surface Transcript 3 and the
    `warn` case in the new file
    `internal/cli/deliver_start_readiness_test.go`.

## Subtasks

- [ ] Print the readiness lines in Setup and fold a failure into its exit code.
- [ ] Run the forge lines in `deliver start` and refuse on a failure.
- [ ] Give the existing `deliver start` tests ready fakes.
- [ ] Describe Setup and the refusal in the command guides.
- [ ] Update the Roundfix Skill, sync the mirrors and record the version.
- [ ] Add the tests.

## Acceptance Criteria

- [ ] Setup prints `gh`, `git`, `remote`, `toolchain` and `environment` after
      `acpx` with Doctor's text, offers nothing for them, and exits `1` when
      one is `failed`.
- [ ] `deliver start` with `DR-GH-UNAUTHENTICATED` exits `2` with Surface
      Transcript 3's output and leaves no Delivery Queue; with only
      `DR-GH-UNREACHABLE` it proceeds to the owner.
- [ ] The Roundfix Skill references carry the codes and the refusal, each
      mirror equals its canonical file, and the raised version is recorded.

## Context

- creates: `internal/cli/setup_readiness_test.go`
- creates: `internal/cli/deliver_start_readiness_test.go`
- interface: `internal/cli/setup.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/deliver_test.go`
- interface: `internal/cli/deliver_limits_test.go`
- interface: `internal/cli/deliver_plan_test.go`
- interface: `internal/cli/deliver_plan_empty_paths_test.go`
- interface: `internal/cli/deliver_token_ceiling_test.go`
- interface: `internal/cli/deliver_usage_test.go`
- interface: `docs/user-guide/commands/setup.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/cli/deliver_prerequisite_test.go`
- instruction: `docs/adr/0220-doctor-checks-what-delivery-needs-and-only-a-definite-answer-fails.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSetupPrintsTheReadinessLinesAndOffersNothing|TestDeliverStartRefusesWhenThisMachineCannotPublish|TestDeliverStartProceedsOnAForgeWarning)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestSetupPrintsTheReadinessLinesAndOffersNothing TestDeliverStartRefusesWhenThisMachineCannotPublish TestDeliverStartProceedsOnAForgeWarning; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three tests do not exist, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestDeliverStartRefusesAPrerequisiteCycle|TestRunSetupFreshMachineAcceptsOffers|TestDeliverStartRefusesWhenThisMachineCannotPublish)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDeliverStartRefusesAPrerequisiteCycle TestRunSetupFreshMachineAcceptsOffers TestDeliverStartRefusesWhenThisMachineCannotPublish; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; the existing prerequisite refusal and Setup tests pass unedited beside a new test, which does not exist before this Task.
- `for phrase in "DR-GH-UNAUTHENTICATED" "DR-TOOL-MISSING" "warn"; do tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/setup.md | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/setup.md "$phrase" >&2; exit 1; }; done; for file in .agents/skills/roundfix/references/deliver.md docs/user-guide/commands/deliver.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "cannot publish from this machine" || { printf 'missing phrase in %s: %s\n' "$file" "cannot publish from this machine" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- "DR-SKILL-TRAILS-SNAPSHOT" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/baseline.md "DR-SKILL-TRAILS-SNAPSHOT" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/setup.md | grep -qF -- "DR-GH-UNAUTHENTICATED" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands/setup.md "DR-GH-UNAUTHENTICATED" >&2; exit 1; }` — expected: exit 0; before this Task none of these files carries these phrases, so the command fails.
- `for name in SKILL.md references/setup.md references/deliver.md references/baseline.md; do cmp -s ".agents/skills/roundfix/$name" "skills/roundfix/$name" || { printf 'mirror differs: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- "DR-GH-UNAUTHENTICATED" skills/roundfix/references/setup.md || { printf 'mirror lacks the new text\n' >&2; exit 1; }; out="$(go test -count=1 -v -run "^TestEveryOwnedSkillVersionIsRecorded$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; printf "%s\\n" "$out" | grep -q -- "--- PASS: TestEveryOwnedSkillVersionIsRecorded" || { printf 'missing pass: %s\n' TestEveryOwnedSkillVersionIsRecorded >&2; exit 1; }` — expected: exit 0; before this Task the mirror lacks the new text, so the command fails; after it, each mirror equals its canonical file and the raised version is recorded.

## References

- `_prd.md` → User Stories 1, 3; Core Features 1, 7, 9; Success Metric 3
- `_techspec.md` → Interfaces; API Contract 2; API Contract 3; Surface Transcript 3; Testing Approach 3-4; Build Order 4
- ADR-0220; ADR-0187; ADR-0189
