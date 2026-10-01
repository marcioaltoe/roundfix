---
task: task_01
spec: 0205-an-advisory-judge-for-spec-authoring
status: completed
type: backend
complexity: medium
---

# Task 01: The judge plans a Spec's judgments from its own artifacts, as measured

## Overview

Roundfix has no code that can ask a semantic question about a Spec. This Task
creates the `internal/judge` package without any network code: the embedded
question file, the only two readers of text the judge may send, the language
gate and the planner that turns one Spec into pending citation-support and
goal-mechanism judgments exactly as the 2026-09-30 measurement built its
datasets. It is verifiable on its own: a fixture Spec yields the claims, goal
pairs and skips the TechSpec's rules name, and the embedded file equals the
TechSpec's block byte for byte.

## Requirements

1. MUST create `internal/judge/questions.json` with bytes identical to the fenced `json` block under `_techspec.md` → Questions and thresholds, embed it, and load it through `Load` as `_techspec.md` → Interfaces sketches. `Load` MUST compile the four patterns (`accepted_model_pattern` and the three extraction patterns) and MUST be the package's only source of question texts, criteria, thresholds, limits, word lists, the pinned model, the accepted-model pattern, the transports with their key variables, endpoints and request model IDs, the price and the ceiling: no other file of the package may repeat one of those values as a literal.
2. MUST implement `readSpecArtifact` and `readADR` with the rules of `_techspec.md` → The only readers: regular files only, never following a symbolic link, at most 1 MiB, `_prd.md` and `_techspec.md` directly in the Spec directory, `docs/adr/<NNNN>-*.md` directly under the repository, and an ADR judged only when `accepted` or without a front matter status. `Source` MUST keep its fields unexported so no other code can build one.
3. MUST implement the language gate of `_techspec.md` → The language gate, over maximal runs of Unicode letters, lowercased.
4. MUST implement `PlanSpec` with every rule of `_techspec.md` → Citation claims and Goal pairs, for the stages `prd`, `techspec` and both, including each skip reason those sections name, the artifact skip `not English`, the artifact skip `not a regular file in the Spec directory`, and one pending judgment per distinct state. A pending judgment's state MUST be built only from `Source` text and the fixed phrase `the cited decision`, and MUST encode as the JSON objects those sections name, without HTML escaping.
5. MUST put the tests in `internal/judge/questions_test.go`, `internal/judge/source_test.go` and `internal/judge/pairs_test.go`, building every fixture Spec, ADR and file in `t.TempDir()`. They MUST cover each numbered rule of Citation claims and Goal pairs with a case that would pass if the rule were removed, an English and a Portuguese fixture for the language gate, and each refusal of The only readers.
6. MUST NOT open a network connection, read `ROUNDFIX_OPENROUTER_API_KEY`, `ROUNDFIX_TYPESAFE_API_KEY` or `OPENROUTER_API_KEY`, or write outside `t.TempDir()`.
7. MUST prove each new gate can fail. The Result MUST record one sabotage of the claim extraction (for example dropping the `(Spec NNNN)` removal) and one of a reader (for example following symbolic links), each with the test that failed, and that the code was restored.

## Subtasks

- [ ] Embed the question file and load it.
- [ ] Add the two readers and the language gate.
- [ ] Plan citation claims and goal pairs with their skips.
- [ ] Cover every rule with a failing-when-removed case.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] `internal/judge/questions.json` equals the TechSpec block byte for byte, loads, pins `jev-1.13`, compiles `accepted_model_pattern`, and lists the `openrouter` transport (`ROUNDFIX_OPENROUTER_API_KEY`, requesting `jev-1.13`) before the `typesafe` one (`ROUNDFIX_TYPESAFE_API_KEY`, requesting `jev-1.13.0`).
- [ ] A fixture sentence `ADR-NNNN keeps ... (Spec 0123).` with one accepted ADR becomes a claim with the token replaced and the Spec reference removed; a sentence with two ADR tokens, one inside a fence, one in a table row, one with a verb outside the pattern and one naming a proposed ADR plan nothing.
- [ ] A Coverage Map line `- Goals 1-2 → <title>.` pairs each goal with the longest qualifying named section, cut at 3500 characters, and a generic or short section is never named.
- [ ] A Portuguese PRD is skipped as `not English` and plans nothing; a symbolic-link PRD is skipped as `not a regular file in the Spec directory`.

## Context

- creates: `internal/judge/questions.json`
- creates: `internal/judge/questions.go`
- creates: `internal/judge/source.go`
- creates: `internal/judge/language.go`
- creates: `internal/judge/pairs.go`
- creates: `internal/judge/questions_test.go`
- creates: `internal/judge/source_test.go`
- creates: `internal/judge/pairs_test.go`
- instruction: `docs/history/findings/2026-09-30-jev-judgments-measured-against-the-spec-archive.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestQuestionFileLoads|TestReadersAcceptOnlySpecArtifactsAndADRs|TestLanguageGateSeparatesEnglishFromPortuguese|TestCitationClaimsFollowTheMeasuredExtraction|TestGoalPairsFollowTheMeasuredSelection|TestPlanSkipsANonEnglishArtifact|TestPlanSkipsASymbolicLinkArtifact)$" ./internal/judge 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQuestionFileLoads TestReadersAcceptOnlySpecArtifactsAndADRs TestLanguageGateSeparatesEnglishFromPortuguese TestCitationClaimsFollowTheMeasuredExtraction TestGoalPairsFollowTheMeasuredSelection TestPlanSkipsANonEnglishArtifact TestPlanSkipsASymbolicLinkArtifact; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the package does not exist, so the command fails.
- `tmp="$(mktemp -d)" || exit 1; awk 'BEGIN { fence = sprintf("%c%c%c", 96, 96, 96) } $0 == "### Questions and thresholds" { h = 1 } h && $0 == fence "json" { f = 1; next } f && $0 == fence { exit } f { print }' docs/specs/0205-an-advisory-judge-for-spec-authoring/_techspec.md > "$tmp/want.json" || exit 1; test -s "$tmp/want.json" || { printf 'no question block in the TechSpec\n' >&2; exit 1; }; cmp "$tmp/want.json" internal/judge/questions.json` — expected: exit 0; before this Task `internal/judge/questions.json` does not exist, so `cmp` fails.

## References

- [_prd.md](_prd.md) — Goals 3 and 5; User Stories 1, 2 and 5; Core Features 2, 3, 4 and 7; Success Metrics 5 and 6; Recorded limits
- [_techspec.md](_techspec.md) — Questions and thresholds; Interfaces; The only readers; Citation claims; Goal pairs; The language gate; Testing Approach 1; Testing Approach 2; Testing Approach 3; Build Order 1
- ADR-0200; ADR-0201; ADR-0176

## Result

Implemented the Task 01 local planning slice. `Load` embeds the exact authored
question block and compiles its four patterns. The two bounded readers are the
only constructors of `Source`; they reject links, non-regular files, oversized
files and inactive ADRs. The planner applies the measured extraction and language
rules, retains original line numbers, builds JSON without HTML escaping and
deduplicates identical states. No network or credential-reading code was added.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| Authored settings, model pin and ordered transports | `TestQuestionFileLoads` compares the embedded bytes with the TechSpec block, checks loaded transport order and fields against that block, and exercises all four compiled patterns. The block retains the authored model IDs, key-variable names, thresholds, price and ceiling; no Go file repeats those values. |
| Measured citation extraction | `TestCitationClaimsFollowTheMeasuredExtraction` covers paragraph boundaries, fences, tables, headings, list markers, sentence splitting, multiple tokens, attribution verbs, proposed and non-English ADRs, transformation, character boundaries, ADR truncation, front matter, line attribution, stage selection and state deduplication. Its exact-state case checks the replacement and removal against the resulting claim, including JSON escaping. |
| Measured goal selection | `TestGoalPairsFollowTheMeasuredSelection` covers PRD item continuation, numbered items, Coverage Map continuation and ranges, table mappings, section boundaries, normalized titles, longest-title selection, generic and short exclusions, Unicode character limits, exact goal states, line numbers, absent goals, unnamed sections, stage selection and deduplication. |
| Non-English and symbolic-link artifact skips | `TestPlanSkipsANonEnglishArtifact` checks Portuguese PRD and TechSpec skips and that a non-English PRD yields no goal pairs. `TestPlanSkipsASymbolicLinkArtifact` checks both artifact names and the required refusal reason. `TestReadersAcceptOnlySpecArtifactsAndADRs` and `TestLanguageGateSeparatesEnglishFromPortuguese` exercise the underlying gates and measured language boundaries. |

Every generated Spec, ADR and source fixture lives in `t.TempDir()`. The
question-file test reads the existing authored TechSpec to check its byte contract.
No live documentation fetch or API call was made, in accordance with this Task's
no-network requirement.

### Focused checks

- Before implementation, `rtk proxy ls internal/judge` exited 1: the package
  did not exist. The only pre-existing tracked change was the Daemon's
  `status: in_progress` in this Task file.
- `GOCACHE=/tmp/roundfix-task01-gocache GOPROXY=off GOSUMDB=off rtk proxy go test -count=1 ./internal/judge`
  exited 0 after the final implementation and restored sabotages.
- `GOOS=windows GOARCH=amd64 GOCACHE=/tmp/roundfix-task01-gocache GOPROXY=off GOSUMDB=off rtk proxy go build ./internal/judge`
  exited 0. The reader uses portable file APIs and checks the descriptor's
  regular-file type and identity before reading.
- The first `GOCACHE=/tmp/roundfix-task01-gocache GOPROXY=off GOSUMDB=off rtk make verify-incremental`
  exited 2. Two existing CLI Force Stop integration tests could not enumerate
  the process table (`operation not permitted`), and suiteguard caught an Agent
  edit to `internal/judge/source.go` while that check was running. No existing
  assertions or Verification configuration were changed.
- The same incremental command rerun with process-table access and a stable
  worktree exited 0: formatting, vet, package tests, skill checks and the CLI
  build passed. The Force Stop integration tests and suiteguard passed on
  this rerun.

### Sabotage evidence

- Claim extraction: temporarily removed the `(Spec NNNN)` stripping operation.
  `go test -count=1 -run TestCitationClaimsFollowTheMeasuredExtraction/rules_4_and_5_exact_state_and_Unicode_ADR_cut ./internal/judge`
  exited 1. The test reported a claim that still contained `(Spec 0123)`.
  Restored the stripping operation; the focused package check then exited 0.
- Reader: temporarily replaced `os.Lstat` with `os.Stat` to follow file links.
  `go test -count=1 -run TestReadersAcceptOnlySpecArtifactsAndADRs/spec_symlink ./internal/judge`
  exited 1 with `accepted=true error=<nil>`. Restored `os.Lstat`; the focused
  package check then exited 0. This sabotage was repeated after making the
  reader portable, with the same failure and restoration.

The sabotage commands used the same task-scoped cache and offline Go settings
as the focused package check. Task status, authored Verification, the Task Graph
and other Task files remain Daemon-owned. No commit, push or Pull Request was
made. Network transport, spending and CLI behavior belong to subsequent Tasks.

## Carry-forward provenance

- Source Run: `run_20261001T205638Z_5e1457a22648dd9b`
- Source commit: `5d53223a3476aa62bee1734fca2622145a5d719d`
