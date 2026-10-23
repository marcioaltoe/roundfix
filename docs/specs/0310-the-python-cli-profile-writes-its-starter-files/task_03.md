---
task: task_03
spec: 0310-the-python-cli-profile-writes-its-starter-files
status: pending
type: backend
complexity: high
---

# Task 03: roundfix baseline starters previews and writes a profile's Starter Files

## Overview

This Task exposes task_02's operation as `roundfix baseline starters`, with
its help, text and JSON output and exit codes, and documents it: the command
guide, the Roundfix skill's Baseline reference, the Skill Coverage Map entry
and the glossary term **Starter File**. The two Surface Transcripts of the
TechSpec are asserted line by line (`_prd.md` → Core Features 3 and 5;
ADR-0323).

## Requirements

1. MUST answer the Backlog Entry Spec 0309 adopted,
   [`../0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md`](../0309-a-python-cli-baseline-profile-for-uv-pytest-and-ruff/references/2026-10-22-a-baseline-profile-for-python-repositories.md)
   of 2026-10-22, its expected item 4.
2. MUST create `internal/cli/baseline_starters.go` with the command of
   `_techspec.md` → API Contract 6: a `parseBaselineStartersCommand` for
   `--profile` (required), `--confirm-plan`, `--repo` (default `.`) and
   `--format` (`text` or `json`), the text renderer of the Surface
   Transcripts, the JSON payload on stdout, the stale line on stderr and the
   exit codes; route `starters` to it in `runBaselineCommand`
   (`internal/cli/baseline_profile.go`); add the synopsis line after
   `roundfix baseline skills reconcile` in the root usage and in
   `commandUsage("baseline")`, and a `commandUsage("baseline starters")` text
   in `internal/cli/cli.go` with the help of API Contract 6.
3. MUST create `internal/cli/baseline_starters_test.go` with the five tests of
   `_techspec.md` → Testing Approach (Build Order 3), calling the CLI in
   process over temporary repositories and homes. The two transcript tests
   assert every stdout line in order, empty stderr and the exit code, with
   only the digest as controlled variation after checking its shape
   (`sha256:` and 64 lowercase hex digits). Each top-level test calls
   `t.Parallel()` first.
4. MUST make the declared test changes of `_techspec.md` → Testing Approach:
   the `baseline` help case of `TestBaselineDocumentationContract` in
   `internal/docscontract/publicdocs_test.go` gains the snippet
   `roundfix baseline starters --profile <id>` and a new `baseline starters`
   case gains the snippets "writes only paths that do not exist", "A non-empty
   preview exits 3" and "An empty preview exits 0";
   `parsePublishedBaselineExample` in
   `internal/cli/baseline_documentation_contract_test.go` parses a `starters`
   example through `parseBaselineStartersCommand`.
5. MUST document the command: in `docs/user-guide/commands/baseline.md`, the
   synopsis line and a `#### Starter Files` section before `### Retired
   skills` that states what a Starter File is, the preview and its exit 3,
   `--confirm-plan`, `keep` for an existing path, the stale refusal, the
   removal on failure, that the repository owns the file afterwards and that
   confirming the preview is the maintainer's express authorization for
   exactly those paths, and the three `python-cli` Starter Files; in
   `.agents/skills/roundfix/references/baseline.md`, a short `### Starter
   Files` part with the same rules and the commands
   `roundfix baseline starters --profile python-cli` and
   `roundfix baseline starters --profile python-cli --confirm-plan <digest>`,
   and the line that an Agent runs the confirmation only with the
   maintainer's express authorization; name `baseline starters` in
   `.agents/skills/roundfix/SKILL.md` where the other Baseline commands are
   named. Neither the guide nor the skill contains `python3` or `Python
   fallback`.
6. MUST raise the Roundfix skill's version through the record command, copy
   the canonical skill with `make skills-sync`, and record it with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   (declaring `skills/roundfix/SKILL.md`,
   `skills/roundfix/references/baseline.md` and
   `skills/testdata/owned-skill-versions.json`); never add or change text
   inside a `### QA settlement` section.
7. MUST add to `docs/references/skill-coverage.json` the entry
   `command baseline starters` whose `skills` is
   `.agents/skills/roundfix/references/baseline.md` and whose `sources` are
   the other `command baseline` entries' sources plus
   `internal/cli/baseline_starters.go`, and add
   `internal/cli/baseline_starters.go` to the `sources` of every other
   `command baseline` entry.
8. MUST add **Starter File** to `CONTEXT.md` through `domain-modeling`, right
   before the **Retired Skill** entry, as "**Starter File**:" followed on the
   next line by "A repository tooling file, such as a Makefile, a Project
   Config or a CI workflow, that a Baseline module offers and `roundfix
   baseline starters` creates only where the path does not exist, after the
   maintainer confirms the exact preview; once written the repository owns
   it, and no refresh changes it", a space, a link to ADR-0323 written as the
   neighbouring entries write theirs (link text `ADR-0323`, target
   `docs/adr/0323-a-profile-starter-file-is-written-only-where-absent-and-only-on-a-confirmed-preview.md`,
   in parentheses) and a closing period, then the line "_Avoid_: Scaffold,
   template file, managed tooling".
9. MUST, as its LAST step and after the last edit of any user guide page,
   re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   which rewrites `docs/references/behavior-surfaces.json`; re-record again if
   a guide changes after it.
10. MUST NOT change any Baseline asset, template, module, profile or setup,
    any other skill, `.roundfixrc.yml` or this repository's Makefile.

## Subtasks

- [ ] Add the command, its routing, usage and help.
- [ ] Add the five CLI tests and the two declared test changes.
- [ ] Document the command in the guide, the Roundfix reference and the glossary; raise, sync and record the skill version.
- [ ] Add the coverage entry, then re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] The preview and apply transcripts match line by line, and the JSON payload, the stale refusal and the invalid inputs behave as API Contract 6 states.
- [ ] A written Project Config loads through Roundfix's loader with the starter's values.
- [ ] The guide, the Roundfix reference, the coverage map and the glossary describe the command, and the documentation contracts pass.

## Consumers

- `roundfix baseline starters` — consumers: `./internal/cli`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0323-a-profile-starter-file-is-written-only-where-absent-and-only-on-a-confirmed-preview.md`
- instruction: `internal/cli/baseline_skills_restore.go`
- creates: `internal/cli/baseline_starters.go`
- interface: `internal/cli/baseline_profile.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/baseline_starters_test.go`
- interface: `internal/cli/baseline_documentation_contract_test.go`
- interface: `internal/docscontract/publicdocs_test.go`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestBaselineStartersPreviewTranscript|TestBaselineStartersApplyTranscript|TestBaselineStartersWritesAProjectConfigRoundfixLoads|TestBaselineStartersJSONAndRefusals|TestBaselineStartersHelp|TestEveryCommandPathPrintsItsHelp|TestBaselineExamplesParse)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineStartersPreviewTranscript TestBaselineStartersApplyTranscript TestBaselineStartersWritesAProjectConfigRoundfixLoads TestBaselineStartersJSONAndRefusals TestBaselineStartersHelp TestEveryCommandPathPrintsItsHelp TestBaselineExamplesParse; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the five starters tests do not exist, so the command fails at its first pass check. Afterwards the transcripts, the written Project Config, the JSON and refusals, the help, every command path's help and the published examples pass, and the package's tests compile for Windows.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Starter File**:' || { echo 'missing glossary term' >&2; exit 1; }; grep -qF -- '"command baseline starters"' docs/references/skill-coverage.json || { echo 'missing coverage entry' >&2; exit 1; }; for f in .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md docs/user-guide/commands/baseline.md; do grep -qF -- 'roundfix baseline starters --profile' "$f" || { printf 'missing command in %s\n' "$f" >&2; exit 1; }; done; go test -count=1 -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestEveryCommandIsNamedInTheRoundfixSkill|TestEveryCommandIsNamedInTheUserGuide|TestBaselineDocumentationContract)$' ./internal/docscontract && go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestNoPythonBaselineRuntime|TestRoundfixSkillIndexNamesEveryReference)$' ./skills && make skills-sync-check` — expected: exit 0. Before this Task the glossary lacks the term, so the command fails at its first check. Afterwards the glossary, the coverage entry and the three documents name the command, the documentation contracts and the skill version record pass, and the mirrors match.

## References

- `_prd.md` → Core Features 3 and 5; Success Metrics 2, 3 and 5; User Stories 1-3; Glossary
- `_techspec.md` → Current behavior 5-8; API Contract 6; Surface Transcript 1; Surface Transcript 2; Invariants 4 and 6; Testing Approach; Build Order 3
- ADR-0323
