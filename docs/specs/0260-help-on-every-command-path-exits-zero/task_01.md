---
task: task_01
spec: 0260-help-on-every-command-path-exits-zero
status: pending
type: backend
complexity: medium
---

# Task 01: Help on every command path exits zero

## Overview

`roundfix qa-report accept --help` reads the flag as a report path and exits
1, and `baseline capabilities`, `baseline skills` and `baseline assets` answer
`--help` with a usage error and exit 2 (`_prd.md`; the adopted entry
[2026-10-09-qa-report-accept-help-exits-one.md](references/2026-10-09-qa-report-accept-help-exits-one.md),
recorded 2026-10-09). This Task makes all four answer `--help` and `-h` with a
usage block and exit 0, adds a test that asks every command path for help,
fingerprints `qa-report accept` from its own help, and corrects the
`roundfix` skill. It is verifiable on its own through the named tests and the
skill phrases.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-09,
   `2026-10-09-qa-report-accept-help-exits-one.md`: `qa-report accept` with a
   help flag prints its usage and exits 0, as `_techspec.md` → API Contract 1
   and Surface Transcripts 1-3 state.
2. MUST make the three `baseline` groups follow API Contract 2 and Surface
   Transcripts 4-6. The usage-error path of a bare group or an unknown command
   after it MUST keep its text, streams and exit 2 (Invariant 3).
3. MUST add `internal/cli/command_help_test.go` with
   `TestHelpSurfaceTranscripts`, which asserts Surface Transcripts 1-6 as
   `_techspec.md` → Testing Approach states, and
   `TestEveryCommandPathPrintsItsHelp`, which follows API Contract 3. Each
   calls `t.Parallel()` as its first statement (ADR-0259).
4. MUST add four rows to `TestRunCommandHelp`, named exactly
   `qa-report accept`, `baseline capabilities`, `baseline skills` and
   `baseline assets`, with arguments `qa-report accept --help`,
   `baseline capabilities --help`, `baseline skills -h` and
   `baseline assets --help`. The first row MUST require
   `roundfix qa-report accept <path>` in standard output, the others
   `roundfix baseline [--repo <path>] [--format <text|json>]`. Existing rows
   MUST NOT change.
5. MUST change `skillCoverageCLIHelp` and
   `TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent` as API
   Contract 4 states, including the failure message
   `command qa-report accept: fingerprint does not use its own help`, and
   remove the nearest-ancestor fallback and its comments.
   `docs/references/behavior-surfaces.json` MUST stay byte-identical
   (Invariant 4).
6. MUST replace, in `.agents/skills/roundfix/references/settle.md`, the
   sentence that says `accept` has no `--help` with: "Read its usage through
   `roundfix qa-report accept --help` or `-h`, which prints it and exits `0`
   without reading a report."
7. MUST add to `.agents/skills/roundfix/references/baseline.md`, after the
   paragraph on stdout and stderr, the sentence: "`roundfix baseline
   capabilities --help`, `roundfix baseline skills --help` and
   `roundfix baseline assets --help`, or `-h`, print the `baseline` usage,
   which lists their commands, and exit `0`."
8. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   to raise and record the `roundfix` skill version. MUST NOT edit text
   inside any `### QA settlement` section.
9. MUST NOT change the user guide, `docs/references/skill-coverage.json`, any
   usage text, or any other command's behavior.

## Subtasks

- [ ] Answer help in `qa-report accept` and the three `baseline` groups.
- [ ] Add the transcript and sweep tests and the four help-exit rows.
- [ ] Make the Behavior Surface Record contract read each command's own help.
- [ ] Correct the settle and baseline skill references and record the version.

## Acceptance Criteria

- [ ] `TestHelpSurfaceTranscripts` passes with Surface Transcripts 1-6 and
      fails if `qa-report accept` reads `--help` as a path again.
- [ ] `TestEveryCommandPathPrintsItsHelp` passes over at least 60 command
      paths and both flags, and names each failing command line otherwise.
- [ ] `TestRunCommandHelp` passes with its four new rows.
- [ ] The two Behavior Surface Record contract tests pass with the record
      unchanged.
- [ ] The skill references carry the new sentences, the mirrors match, and the
      owned-skill version record passes.
- [ ] Existing `qa_report_test.go` and `baseline_profile_test.go` cases pass
      unchanged.

## Context

- instruction: `docs/agents/cli.md`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- instruction: `docs/adr/0184-a-techspec-states-a-command-surface-as-a-transcript.md`
- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- interface: `internal/cli/qa_report.go`
- interface: `internal/cli/baseline_profile.go`
- creates: `internal/cli/command_help_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/docscontract/skill_coverage_test.go`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestHelpSurfaceTranscripts|TestEveryCommandPathPrintsItsHelp|TestRunCommandHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHelpSurfaceTranscripts TestEveryCommandPathPrintsItsHelp TestRunCommandHelp/qa-report_accept TestRunCommandHelp/baseline_capabilities TestRunCommandHelp/baseline_skills TestRunCommandHelp/baseline_assets; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two tests and four rows do not exist, so no PASS line is printed.
- `out="$(go test -count=1 -tags docscontract -v -run '^(TestTheSkillCoverageMapIsCurrent|TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheSkillCoverageMapIsCurrent TestBehaviorSurfaceBytesComeFromParsedValuesAndGuideContent; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < internal/docscontract/skill_coverage_test.go | grep -qF -- 'fingerprint does not use its own help' || { printf 'the contract does not assert the own-help fingerprint\n' >&2; exit 1; }; ! grep -qF -- 'nearest successful ancestor' internal/docscontract/skill_coverage_test.go || { printf 'the ancestor fallback remains\n' >&2; exit 1; }` — expected: exit 0; before this Task the own-help assertion is absent and the fallback remains.
- `for file in .agents/skills/roundfix/references/settle.md skills/roundfix/references/settle.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'which prints it and exits' || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; ! tr -s '[:space:]' ' ' < "$file" | grep -q -- 'has no .--help' || { printf 'stale help sentence in %s\n' "$file" >&2; exit 1; }; done; for file in .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'which lists their commands, and exit' || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; done; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1` — expected: exit 0; before this Task neither phrase exists and the stale sentence is present.

## References

- [_prd.md](_prd.md) — Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Invariants 1-4; Testing Approach; Build Order 1
- [2026-10-09-qa-report-accept-help-exits-one.md](references/2026-10-09-qa-report-accept-help-exits-one.md)
- ADR-0256; ADR-0184; ADR-0259; ADR-0189; ADR-0233
