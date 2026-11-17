---
task: task_01
spec: 0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation
status: pending
type: backend
complexity: medium
---

# Task 01: The Model Catalog follows the installed Codex and Claude adapters

## Overview

The Model Catalog that `profiles configure` and the pickers offer follows
what Codex 0.169.0 and the Claude ACP adapter 0.107.3 advertise: it gains
`gpt-6.1-terra`, offers the `fable` alias in place of `claude-fable-5-1`, and
describes `haiku` as Haiku 5.5. The guides, the Roundfix skill's `profiles`
reference and the model selection reference list the same catalog. The
Recommended Profile does not change here; task 02 changes it once every model
it names is in the catalog.

## Requirements

1. MUST make the Codex catalog in `internal/agent/catalog.go`, in this order:
   `gpt-6.1-sol`, `gpt-6-astra`, `gpt-6.1-terra`, `gpt-6-sol`, `gpt-6-luna`,
   `gpt-5.6-sol`, `gpt-5.6-terra`, `gpt-5.6-luna`, each with label equal to
   value (`_techspec.md` → Model Catalog). `gpt-6.1-terra` is described as
   the fast, efficient writing model and `gpt-6-luna` as fast and affordable;
   "built-in review default" stays on `gpt-5.6-luna` until task 02 moves it.
2. MUST make the Claude catalog, in this order: `opus`, `sonnet`, `fable`,
   `haiku`, `default`; `fable` described as the alias for Fable 5.2, the most
   capable model at the highest latency and quota cost, and `haiku` as Haiku
   5.5, the fastest and cheapest. Its comment names adapter 0.107.3 and the
   date 2026-11-17. `claude-fable-5-1` MUST no longer be offered.
3. MUST NOT change the picker's reasoning efforts in `internal/cli/cli.go`
   and MUST NOT add `ultra` anywhere.
4. MUST update `TestModelCatalogOpensWithTheCurrentModels` to the two new
   lists, `TestModelCatalogsExposeOrderedPickerData` to the new entries with
   their descriptions, and add `gpt-5.5` and `claude-fable-5-1` to the
   retired values of `TestModelCatalogOffersNoRetiredModel`, so its subtests
   `codex/gpt-5.5` and `claude/claude-fable-5-1` exist. These are declared
   changes of those tests' expectations.
5. MUST change the catalog sentences of `docs/user-guide/usage.md`,
   `docs/user-guide/configuration.md`, the Roundfix skill's `profiles`
   reference and the `### Model Catalog and picker efforts` subsection of
   `docs/references/model-selection.md` to the new lists, naming
   `gpt-6.1-terra` and `fable`, and the adapter versions read on
   2026-11-17 (`codex-acp` 2.7.0, `claude-agent-acp` 0.107.3) in the
   reference. The adapter floors and the snapshot date stay.
6. MUST change `TestProfilesDocumentationContractMatchesPublicGuidance` in
   `internal/docscontract/publicdocs_test.go` to require the backticked
   `fable` (the string with its two backticks) where it required
   `claude-fable-5-1`, and change nothing else in that file.
7. MUST regenerate, never hand-edit, with the commands of
   `_authorization.md` → Sanctioned regeneration, in this order:
   `make skills-sync`; the `-record-skill-versions` step
   (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
   and the `-record-skill-coverage` step
   (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`).
8. MUST change no file outside the ones this Task declares; the Governed
   Paths it changes are in `_authorization.md` `paths`.

## Subtasks

- [ ] Change both catalog lists and their descriptions.
- [ ] Update the three catalog tests.
- [ ] Update the catalog sentences of the guides, the skill reference and the model selection reference.
- [ ] Update the docs contract test and run the sanctioned regeneration.

## Acceptance Criteria

- [ ] The Codex catalog lists `gpt-6.1-terra` after `gpt-6-astra`, and the Claude catalog lists `fable` and no `claude-fable-5-1`.
- [ ] No catalog offers a retired model, `gpt-5.5` and `claude-fable-5-1` included.
- [ ] The guides, both copies of the `profiles` reference and the model selection reference name the same catalog.
- [ ] The skill mirror, the version record and the Behavior Surface Record are current.

## Context

- instruction: `docs/specs/0403-the-update-check-names-each-profile-that-differs-from-a-refreshed-recommendation/references/2026-11-17-model-evidence.md`
- interface: `internal/agent/catalog.go`
- interface: `internal/agent/catalog_current_test.go`
- interface: `internal/agent/agent_test.go`
- interface: `internal/docscontract/publicdocs_test.go`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/model-selection.md`
- interface: `.agents/skills/roundfix/references/profiles.md`
- interface: `skills/roundfix/references/profiles.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `agent.ModelCatalog` — consumers: `./internal/config`, `./internal/cli`, `./internal/tui`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestModelCatalogOpensWithTheCurrentModels|TestModelCatalogsExposeOrderedPickerData|TestModelCatalogOffersNoRetiredModel)$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestModelCatalogOffersNoRetiredModel/codex/gpt-5.5 TestModelCatalogOffersNoRetiredModel/claude/claude-fable-5-1 TestModelCatalogOpensWithTheCurrentModels/codex TestModelCatalogOpensWithTheCurrentModels/claude; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; grep -qF -- 'Value: "gpt-6.1-terra"' internal/agent/catalog.go && grep -qF -- 'Value: "fable"' internal/agent/catalog.go && ! grep -qF -- 'claude-fable-5-1' internal/agent/catalog.go && go test -count=1 -timeout 30m ./internal/agent ./internal/config ./internal/cli ./internal/tui` — expected: exit 0. Before this Task the two retired subtests do not exist and the catalog lacks `gpt-6.1-terra` and `fable`, so the command fails; afterwards the catalog tests pass by name and every package that reads the catalog passes.
- `grep -qF -- '0.107.3' docs/references/model-selection.md && grep -qF -- 'gpt-6.1-terra' docs/references/model-selection.md && ! grep -qF -- 'claude-fable-5-1' docs/user-guide/usage.md && ! grep -qF -- 'claude-fable-5-1' .agents/skills/roundfix/references/profiles.md && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 20m ./skills && make skills-sync-check` — expected: exit 0. Before this Task the reference names adapter 0.84.0 and the guides list `claude-fable-5-1`; afterwards the guides, the docs contract, the mirror and the records pass.

## References

- `_prd.md` → Goals (second bullet); Core Features 2 and 4; Success Metric 2
- `_techspec.md` → Model Catalog; Interfaces (invariant 2); Testing Approach (Task 1); Build Order 1
- [`references/2026-11-17-model-evidence.md`](references/2026-11-17-model-evidence.md)
- ADR-0180, ADR-0235, ADR-0256
