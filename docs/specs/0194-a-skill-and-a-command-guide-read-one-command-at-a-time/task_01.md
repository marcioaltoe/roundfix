---
task: task_01
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
status: completed
type: backend
complexity: medium
---

# Task 01: Every contract reads a document's entry file with its companion files

## Overview

The Roundfix skill and the command reference are about to be split into an entry file plus companion files. Several contracts pin a phrase in one of those documents by reading the single file: `roundfix skills check`, and tests in `internal/docscontract`, `internal/cli` and `skills`. After the split a phrase may live in a companion file. A required phrase would then be reported missing, and a forbidden phrase in a companion file would go unseen. This Task lands first, while both documents are still single files. It adds one reader that returns the entry file followed by its companion files, and switches every such contract to it. With no companion directory the reader returns the entry file alone, so every contract behaves exactly as before.

## Requirements

1. MUST add the package `internal/mdtree` with `Text(fsys fs.FS, entry, companion string) (string, error)`, as the TechSpec's "Interfaces" section states:
   - it returns the bytes of `entry`, then every `.md` file directly inside `companion` in lexical order, each preceded by one newline;
   - a missing `companion` directory yields `entry` alone;
   - a missing `entry` is an error;
   - nested directories and files without the `.md` suffix are ignored.
2. MUST make the required-wording and banned-branding check of `CheckReadiness` in `skills/skills.go` read the Roundfix skill through `mdtree.Text` (`roundfix/SKILL.md` with `roundfix/references`). The check moves into an unexported function over an `fs.FS`, so a test can drive it with `fstest.MapFS`. Every diagnostic keeps its text and its `roundfix/SKILL.md` path, and the `roundfix/agents/openai.yaml` check is unchanged.
3. MUST switch every Go reader of either document's content to `mdtree.Text`. This includes the two contract tests an earlier Spec added, `TestEveryCommandIsNamedInTheRoundfixSkill` in `internal/docscontract/command_documentation_test.go` and `TestEveryCommandIsNamedInTheUserGuide` in `internal/docscontract/user_guide_contract_test.go`: the first must read the skill's references and the second must read `docs/user-guide/commands/` as well as the top-level guide files, so neither goes blind when sections move. Those two files do not exist when this Spec is authored, so they are named here rather than declared in Context; the Daemon records them at the Task commit. The pairs are `.agents/skills/roundfix/SKILL.md` with `.agents/skills/roundfix/references`, `skills/roundfix/SKILL.md` with `skills/roundfix/references`, and `docs/user-guide/commands.md` with `docs/user-guide/commands`. The readers known when this Spec was written are the six tests in the TechSpec's "Pins read the tree" table. It MUST run the sweep command that section gives, switch every other reader the sweep finds, and name each file it changed in its Result.
4. MUST leave these readers on the entry file: `TestSettlementGuidanceIsOneTable`, `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical`, and any test that only checks that an installed or embedded `roundfix/SKILL.md` exists.
5. MUST NOT move, split or edit the Roundfix skill, the command reference or any other Markdown file.
6. MUST NOT rename or remove a top-level test, change an assertion's expected phrases, or change an exported function signature. It MUST edit no governed file other than `internal/docscontract/publicdocs_test.go`, `internal/cli/cli_test.go`, `internal/cli/baseline_documentation_contract_test.go` and `skills/baseline_skill_contract_test.go`. When the sweep finds a reader in another governed file, it MUST leave that file unchanged and name it in its Result.
7. MUST put the new tests in `internal/mdtree/mdtree_test.go` and `skills/roundfix_wording_test.go`, over `fstest.MapFS` only.
8. MUST NOT name any decision by its `ADR-` identifier in this Task's `## Result`.

## Subtasks

- [ ] Add the reader and its tests.
- [ ] Read the Roundfix skill through it in the skill check, and prove it with an in-memory bundle.
- [ ] Switch every test that reads either document, then run the sweep.

## Acceptance Criteria

- [ ] `Text` returns the entry alone without a companion directory, appends companion files in lexical order, ignores nested directories and other suffixes, and reports a missing entry.
- [ ] With an in-memory bundle, a required phrase that lives only in a reference file satisfies the skill check, a phrase missing from every file is reported for `roundfix/SKILL.md`, and banned branding in a reference file is reported.
- [ ] `roundfix skills check` still reports no diagnostic for the shipped bundle.
- [ ] No test in the four governed test files reads `docs/user-guide/commands.md` or the Roundfix skill's `SKILL.md` content directly, and every such test passes unchanged in what it asserts.

## Context

- creates: `internal/mdtree/mdtree.go`
- creates: `internal/mdtree/mdtree_test.go`
- interface: `skills/skills.go`
- creates: `skills/roundfix_wording_test.go`
- interface: `internal/docscontract/publicdocs_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/baseline_documentation_contract_test.go`
- interface: `skills/baseline_skill_contract_test.go`
- instruction: `skills/skills_test.go`
- instruction: `skills/settlement_guidance_repocontract_test.go`
- instruction: `docs/adr/0187-the-roundfix-skill-and-the-command-reference-are-read-one-command-at-a-time.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTextReturnsTheEntryAloneWithoutCompanions|TestTextAppendsCompanionFilesInLexicalOrder|TestTextIgnoresNestedDirectoriesAndOtherSuffixes|TestTextReportsAMissingEntry|TestRoundfixWordingIsSatisfiedByAReferenceFile|TestRoundfixWordingReportsAPhraseMissingFromEveryFile|TestRoundfixWordingReportsBannedBrandingInAReference|TestCheckValidatesRoundfixSkillArtifacts|TestNoPythonBaselineRuntime|TestEventsHelpDocumentsAgentSelectionFilter|TestBaselineExamplesParse)$" ./internal/mdtree ./skills ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTextReturnsTheEntryAloneWithoutCompanions TestTextAppendsCompanionFilesInLexicalOrder TestTextIgnoresNestedDirectoriesAndOtherSuffixes TestTextReportsAMissingEntry TestRoundfixWordingIsSatisfiedByAReferenceFile TestRoundfixWordingReportsAPhraseMissingFromEveryFile TestRoundfixWordingReportsBannedBrandingInAReference TestCheckValidatesRoundfixSkillArtifacts TestNoPythonBaselineRuntime TestEventsHelpDocumentsAgentSelectionFilter TestBaselineExamplesParse; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the package `internal/mdtree` and the seven new named tests do not exist, so the command fails.
- `docs="$(go test -count=1 -tags docscontract -v -run "^(TestBaselineDocumentationContract|TestProfilesDocumentationContractMatchesPublicGuidance|TestReleasePlanDocumentationContract)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$docs"; exit 1; }; for name in TestBaselineDocumentationContract TestProfilesDocumentationContractMatchesPublicGuidance TestReleasePlanDocumentationContract; do printf "%s\\n" "$docs" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for file in internal/docscontract/publicdocs_test.go internal/cli/cli_test.go internal/cli/baseline_documentation_contract_test.go skills/baseline_skill_contract_test.go; do grep -q -e "mdtree" "$file" || { printf 'no tree reader in %s\n' "$file" >&2; exit 1; }; if grep -q -e '"user-guide", "commands.md"' -e '"roundfix", "SKILL.md"' "$file"; then printf 'direct document read in %s\n' "$file" >&2; exit 1; fi; done` — expected: exit 0; before this Task the four files read both documents directly and none uses the reader, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 4; Core Feature 3; Success Metric 3
- [_techspec.md](_techspec.md) — Interfaces; Pins read the tree; API Contract 1; Testing Approach 1, Testing Approach 2; Build Order 1
- ADR-0187

## Result

Implemented the Task slice without splitting or editing either source document.
Task status and the declared Verification commands remain Daemon-owned; no
commit, push or pull request was made.

### Implementation and acceptance evidence

- `internal/mdtree/mdtree.go` adds `Text(fs.FS, entry, companion)`: it preserves
  entry bytes, adds one newline before each direct `.md` companion in lexical
  order, ignores nested directories and other suffixes, treats a missing
  companion directory as the entry alone, and wraps other read errors.
  `internal/mdtree/mdtree_test.go` exercises all four required reader cases
  with `fstest.MapFS`, including preservation of existing trailing newlines
  and matching the missing-entry error with `errors.Is`.
- `skills/skills.go` delegates operational wording checks to the unexported
  `checkRoundfixWording(fs.FS)`. Required wording and branding cover the entry
  and `roundfix/references`; diagnostic messages and the
  `roundfix/SKILL.md` diagnostic path are preserved. The manifest still reads
  only `roundfix/agents/openai.yaml` and uses the existing manifest checks.
  `skills/roundfix_wording_test.go` uses self-contained `fstest.MapFS` bundles
  to prove reference-only required wording, exact missing-phrase diagnostics,
  and both banned-branding spellings in a reference.
- Migrated every wording reader in the four authorized governed test files:
  `internal/docscontract/publicdocs_test.go`, `internal/cli/cli_test.go`,
  `internal/cli/baseline_documentation_contract_test.go` and
  `skills/baseline_skill_contract_test.go`. Their expected phrases and
  top-level test names are unchanged.
- Also migrated `internal/docscontract/command_documentation_test.go`
  (including its deliberately removed-command check) and
  `internal/docscontract/user_guide_contract_test.go`. Command coverage reads
  the command entry with companions alongside the other top-level guides;
  link and forbidden-flag checks enumerate companion files separately so
  source-relative links and file/line diagnostics retain their meaning.
- The shipped bundle remains accepted: `go run ./cmd/roundfix skills check`
  exited 0 and reported the skill check passed, with no diagnostics.

### Focused checks

All Go commands below used
`GOCACHE=/private/tmp/roundfix-0194-task01-gocache`; the initial attempts using
Go's default cache were refused by the sandbox before compilation.

- Pre-change `rtk proxy go test ./internal/mdtree`: exit 1, package directory
  absent, establishing the missing reader before implementation.
- `rtk proxy go test -count=1 ./internal/mdtree ./skills -run
  'TestText|TestRoundfixWording|TestCheckValidatesRoundfixSkillArtifacts|TestNoPythonBaselineRuntime'`:
  exit 0 for both packages.
- `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run
  'TestEveryCommandIsNamed|TestAnUndocumentedCommand|TestBaselineDocumentationContract|TestProfilesDocumentationContractMatchesPublicGuidance|TestReleasePlanDocumentationContract|TestUserGuide'`:
  exit 0.
- `rtk proxy go test -count=1 ./internal/cli -run
  'TestEventsHelpDocumentsAgentSelectionFilter|TestBaselineExamplesParse'`:
  exit 0.
- After the final fixture and path cleanup, `rtk proxy go test -count=1
  ./skills -run 'TestRoundfixWording|TestNoPythonBaselineRuntime'`: exit 0.
- `rtk proxy go run ./cmd/roundfix skills check`: exit 0.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exit 0.

- `rtk proxy go test -count=1 ./internal/spec -run
  '^TestRepositorySpecCorpusStillLoads$'`: exit 0 after restoring the Task's
  authored projection.
- `rtk make verify-incremental`: final exit 0 with process-table permission;
  formatting, vet, package tests, skill mirror checks, shipped skill checks
  and the CLI build passed. Earlier attempts exposed edits made while the
  mutation guard was active, sandbox-denied process-table access, and a
  Result-writer error that temporarily truncated the Task after an earlier
  mention of its Result heading. The authored sections were restored from
  HEAD with the current status preserved, and the final run held the tree
  stable throughout. No check or assertion was suppressed.

### Sweep and scope

Ran the TechSpec sweep before and after migration:
`grep -rn --include='*.go' -e 'commands\.md' -e '"roundfix", "SKILL.md"'
-e 'roundfix/SKILL\.md' .`.

Every additional wording reader found was migrated in the two newer
`internal/docscontract` files listed above. No additional governed wording
reader required an out-of-grant edit. Left
`skills/settlement_guidance_repocontract_test.go` and
`skills/owned_skill_edit_repocontract_test.go` byte-identical as required.
Other remaining hits in `skills/skills_test.go`, `skills/repository_test.go`,
`internal/baseline/derived_ownership_test.go`, `internal/cli/doctor_test.go`,
`internal/speccheck/{mechanical,surface,verification_wrap}_test.go`,
`internal/verifyselect/verifyselect_test.go` and
`internal/releaseplan/classifier_test.go` are artifact existence checks,
entry/version mutations or path/classification fixtures, rather than
whole-document wording contracts, and remain unchanged.

The implementation adds the two `internal/mdtree` files and
`skills/roundfix_wording_test.go`, changes `skills/skills.go` and the six test
files named above, and records this Result in this Task file. No other Task,
Task Graph, skill Markdown, command-reference Markdown, build configuration,
manifest or derived artifact was edited.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/docscontract/command_documentation_test.go`
- `internal/docscontract/user_guide_contract_test.go`

## Carry-forward provenance

- Source Run: `run_20261001T010622Z_8bfe4bebad66f971`
- Source commit: `c084d62af1ad6e760a41166558c2a209d0f93d3f`
