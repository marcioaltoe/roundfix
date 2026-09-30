---
task: task_01
spec: 0189-profiles-that-follow-the-current-models
status: pending
type: backend
complexity: medium
---

# Task 01: The Model Catalog, the picker efforts and the adapter floors follow today's adapters

## Overview

The Model Catalog in `internal/agent/catalog.go` still offers `gpt-5.4`, `gpt-5.4-mini` and `gpt-5.3-codex-spark`, which Codex has retired. It offers no GPT-6 model and no `claude-fable-5-1`. The picker's reasoning efforts in `reasoningEffortChoices` (`internal/cli/cli.go`) do not match what either adapter advertises. The adapter floors in `internal/agent/acpx_runner.go` are older than the versions the current tuples were proved on. This Task replaces all three with the values under "Reference data" in `_techspec.md`, which record what the installed adapters advertised on 2026-09-30. It is verifiable on its own: the picker data, Doctor's adapter refusal and the documents all state the new values.

## Requirements

1. MUST make `ModelCatalog("codex")` and `ModelCatalog("claude")` return the entries of `_techspec.md` → Reference data → Model Catalog, in that order, plus the two transitional entries at the positions that section gives. `Label` MUST equal `Value`, and the descriptions MUST be the ones given there. `gpt-5.4`, `gpt-5.4-mini` and `gpt-5.3-codex-spark` MUST be removed. `ModelCatalog("opencode")` MUST stay empty. The comment above the Claude catalog MUST state the adapter version and date of the PRD's recorded advertisement.
2. MUST make `reasoningEffortChoices` return `low`, `medium`, `high`, `xhigh`, `max` for `codex` and `default`, `low`, `medium`, `high`, `xhigh`, `max` for `claude`, and nothing for any other runtime.
3. MUST set `PinnedCodexAdapterVersion` to `2.0.1` and `PinnedClaudeAdapterVersion` to `0.84.0`. Every message, install action and generated override MUST keep reading the two constants, so no other production line repeats a version literal.
4. MUST update only the existing tests this change invalidates, and name each in the Result. Where a test means "the floor", it MUST read the constant instead of repeating the literal. It MUST rename or remove no top-level test and change no exported function signature.
5. MUST put the new tests in `internal/agent/catalog_current_test.go`, `internal/cli/effort_choices_test.go` and `internal/cli/adapter_floor_test.go`. The catalog tests MUST require each catalog to open with the models of the Reference data in order, and MUST NOT require the absence of the two transitional entries. The floor tests MUST drive Doctor through the existing fake health and adapter seams, never a real adapter. One MUST prove that an official Codex adapter reporting `1.1.5` and an official Claude adapter reporting `0.63.0` are each refused with an install action that names the floor. Another MUST prove that the same adapters reporting the floor pass.
6. MUST prove each new gate can fail. The Result MUST record one sabotage for the catalog tests (for example restoring `gpt-5.4`) and one for the floor tests (for example lowering a floor), each with the test that failed, and that the code was restored.
7. MUST update, in place, the sentences of `docs/user-guide/configuration.md`, `docs/user-guide/usage.md`, `docs/user-guide/commands.md` and `.agents/skills/roundfix/SKILL.md` that state the adapter floors, their `npm install -g` and `npx -y` commands, the official Codex and Claude identifiers, and the model `opus` resolves to. No sentence in those four files may still name `codex-acp@1.1.5` or `claude-agent-acp@0.63.0`. It MUST update the adapter version in `.agents/skills/roundfix/agents/openai.yaml` to `codex-acp 2.0.1`. It MUST add no text inside the `### QA settlement` section of the skill, then run `make skills-sync` and `make baseline-digests`, and name in the Result every file either command rewrote.
8. MUST NOT remove `gpt-5.5` or `claude-fable-5` from the catalog, and MUST NOT remove the strings `gpt-5.5`, `2026-08-07` or `roundfix/profiles/v1` from the guides or the skill. Later Tasks own those entries, those sentences and the contract test that pins them.
9. MUST NOT call a real ACP adapter, a provider or the network in any test.

## Subtasks

- [ ] Replace the two catalogs and the two effort lists.
- [ ] Raise the two adapter floors.
- [ ] Update the invalidated tests and add the new ones, each negative case separate.
- [ ] Update the guide, skill and manifest sentences, then sync the mirror and the digests.
- [ ] Record one sabotage per new gate in the Result.

## Acceptance Criteria

- [ ] The Codex catalog opens with the seven models of the Reference data in order and ends with `gpt-5.5`.
- [ ] The Claude catalog is `opus`, `sonnet`, `claude-fable-5-1`, `claude-fable-5`, `haiku`, `default`.
- [ ] No catalog offers `gpt-5.4`, `gpt-5.4-mini` or `gpt-5.3-codex-spark`.
- [ ] The documentation contract still passes.
- [ ] The picker offers the two effort lists of Requirement 2, and no `maximum`.
- [ ] Doctor refuses `codex-acp` 1.1.5 and `claude-agent-acp` 0.63.0 with install actions naming `2.0.1` and `0.84.0`, and accepts both at the floor.
- [ ] The three guides, the skill, its mirror and both manifests carry the new floors and identifiers, none names an old floor, and `make skills-sync-check` passes.

## Context

- interface: `internal/agent/catalog.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/cli/cli.go`
- interface: `internal/agent/agent_test.go`
- interface: `internal/agent/acpx_runner_test.go`
- instruction: `internal/docscontract/publicdocs_test.go`
- creates: `internal/agent/catalog_current_test.go`
- creates: `internal/cli/effort_choices_test.go`
- creates: `internal/cli/adapter_floor_test.go`
- instruction: `internal/cli/doctor.go`
- instruction: `internal/cli/setup.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/agents/openai.yaml`
- interface: `skills/roundfix/agents/openai.yaml`

## Verification

- `out="$(go test -count=1 -v -run "^(TestModelCatalogOpensWithTheCurrentModels|TestModelCatalogOffersNoRetiredModel|TestReasoningEffortChoicesFollowTheAdapters|TestDoctorRefusesAnAdapterBelowTheFloor|TestDoctorAcceptsAnAdapterAtTheFloor|TestModelCatalogsExposeOrderedPickerData|TestCheckAdapterProvesOfficialClaudePackageAndVersion|TestModelRecommendationsUseOfficialCatalogModels)$" ./internal/agent ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestModelCatalogOpensWithTheCurrentModels TestModelCatalogOffersNoRetiredModel TestReasoningEffortChoicesFollowTheAdapters TestDoctorRefusesAnAdapterBelowTheFloor TestDoctorAcceptsAnAdapterAtTheFloor TestModelCatalogsExposeOrderedPickerData TestCheckAdapterProvesOfficialClaudePackageAndVersion TestModelRecommendationsUseOfficialCatalogModels; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the five new named tests exists, so the command fails.
- `for pair in "docs/user-guide/configuration.md|codex-acp@2.0.1" "docs/user-guide/configuration.md|claude-agent-acp@0.84.0" "docs/user-guide/configuration.md|claude-fable-5-1" "docs/user-guide/configuration.md|gpt-6.1-sol" "docs/user-guide/commands.md|codex-acp@2.0.1" "docs/user-guide/commands.md|claude-agent-acp@0.84.0" "docs/user-guide/usage.md|codex-acp@2.0.1" "docs/user-guide/usage.md|claude-agent-acp@0.84.0" "docs/user-guide/usage.md|claude-fable-5-1" "docs/user-guide/usage.md|gpt-6.1-sol" ".agents/skills/roundfix/SKILL.md|codex-acp@2.0.1" ".agents/skills/roundfix/SKILL.md|claude-agent-acp@0.84.0" ".agents/skills/roundfix/SKILL.md|claude-fable-5-1" ".agents/skills/roundfix/SKILL.md|gpt-6.1-sol" "skills/roundfix/SKILL.md|codex-acp@2.0.1" "skills/roundfix/SKILL.md|claude-agent-acp@0.84.0" "skills/roundfix/SKILL.md|claude-fable-5-1" "skills/roundfix/SKILL.md|gpt-6.1-sol" ".agents/skills/roundfix/agents/openai.yaml|codex-acp 2.0.1" "skills/roundfix/agents/openai.yaml|codex-acp 2.0.1"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && ! grep -rqE 'codex-acp@1\.1\.5|claude-agent-acp@0\.63\.0' docs/user-guide .agents/skills/roundfix skills/roundfix && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check && out="$(go test -count=1 -tags docscontract -v -run "^(TestProfilesDocumentationContractMatchesPublicGuidance)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProfilesDocumentationContractMatchesPublicGuidance; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task no document names `codex-acp@2.0.1`, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 2; Core Features 1 and 2; Success Metrics 2 and 4; Acceptance evidence
- [_techspec.md](_techspec.md) — Reference data; API Contract 2; API Contract 5; Testing Approach 1–2; Build Order 1
- ADR-0180; ADR-0079; ADR-0147

## Result
