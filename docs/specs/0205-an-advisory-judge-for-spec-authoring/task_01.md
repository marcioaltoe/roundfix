---
task: task_01
spec: 0205-an-advisory-judge-for-spec-authoring
status: pending
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
6. MUST NOT open a network connection, read `ROUNDFIX_JEV_OPENROUTER_API_KEY`, `TYPESAFE_API_KEY` or `OPENROUTER_API_KEY`, or write outside `t.TempDir()`.
7. MUST prove each new gate can fail. The Result MUST record one sabotage of the claim extraction (for example dropping the `(Spec NNNN)` removal) and one of a reader (for example following symbolic links), each with the test that failed, and that the code was restored.

## Subtasks

- [ ] Embed the question file and load it.
- [ ] Add the two readers and the language gate.
- [ ] Plan citation claims and goal pairs with their skips.
- [ ] Cover every rule with a failing-when-removed case.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] `internal/judge/questions.json` equals the TechSpec block byte for byte, loads, pins `jev-1.13`, compiles `accepted_model_pattern`, and lists the `openrouter` transport (`ROUNDFIX_JEV_OPENROUTER_API_KEY`, requesting `jev-1.13`) before the `typesafe` one (`TYPESAFE_API_KEY`, requesting `jev-1.13.0`).
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
