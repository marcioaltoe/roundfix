---
task: task_07
spec: 0082-the-manifest-already-answered-that
status: completed
type: docs
complexity: medium
---

# Task 07: Teach the update path to the docs and the owned skills

## Overview

The update command is only useful if the operating contract, the skill that
dispatches Baseline work, and the guidance generated into adopting repositories
all name it. This is an authorized protected-tooling task: it edits
Roundfix-owned Skills and Baseline module assets, and it may change nothing
else. Its boundary is the exact file list recorded in the authorization.

## Requirements

1. MUST update the durable operating contract in the user guide to document the
   update command, its flags, its exit categories, and the managed-refresh
   guarantee that non-managed bytes are preserved.
2. MUST update the `setup-context-driven` skill so its recipes route an already
   adopted repository to the update command and stop implying that every refresh
   requires the full interview.
3. MUST update the `roundfix` skill so the shipped CLI surface it describes
   matches the binary, satisfying the repository's CLI skill-sync rule.
4. MUST update any Baseline module asset that names the Baseline command family
   so generated guidance teaches the update path, and MUST change module assets
   for no other reason.
5. MUST regenerate derived digest pins through the sanctioned regeneration
   command rather than transcribing any value by hand.
6. MUST teach the published-example contract test to route the update
   subcommand to its own argument parser. Documenting a new subcommand makes a
   truthful Bash example unparseable until that routing exists, so this is part
   of documenting the command, not a separate concern. The test file is a test,
   not protected tooling, and needs no tooling authorization — only scope.
7. MUST change only these repository-relative paths plus this Task file:
   `skills/setup-context-driven/**` and `.agents/skills/setup-context-driven/**`;
   `skills/roundfix/**` and `.agents/skills/roundfix/**`;
   `internal/baseline/assets/modules/*.json`;
   `docs/user-guide/context-driven-development.md`;
   `internal/cli/baseline_documentation_contract_test.go`; and the pins the
   sanctioned regeneration command rewrites. Any other changed path fails this
   Task.
8. MUST NOT weaken any module's Normative Clauses, decisions, capabilities, or
   template selection.

## Subtasks

- [x] Document the update command in the durable operating contract.
- [x] Route adopted repositories to the update path in the setup skill.
- [x] Sync the roundfix skill with the shipped CLI surface.
- [x] Name the update path in the module assets that name the command family.
- [x] Route the update subcommand in the published-example contract test.
- [x] Run the sanctioned digest regeneration and keep its output unedited.
- [x] Confirm the changed-file set matches the authorized boundary exactly.

## Acceptance Criteria

- [x] The user guide documents the command, every flag, every exit category, and
      the preservation guarantee.
- [x] The setup skill no longer presents the full interview as the only refresh
      route for an adopted repository.
- [x] The roundfix skill's described CLI surface matches the binary's usage
      output for the Baseline command family.
- [x] Generated guidance produced from the module assets names the update path.
- [x] Every derived pin equals the value the sanctioned regeneration command
      produces; no pin was hand-edited.
- [x] Every published `roundfix baseline update` example in the user guide
      parses through the update command's own argument parser, so a documented
      example that the CLI would reject fails the contract test.
- [x] The changed-file set is a subset of the authorized boundary.

## Context

- instruction: `docs/workflow/authorizations/2026-08-07-baseline-update-command.md`
- instruction: `docs/agents/specific-repository.md`
- instruction: `docs/agents/skill-dispatch.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exits 0.
- `grep -q 'baseline update' docs/user-guide/context-driven-development.md` — expected: exits 0.
- `grep -q 'baseline update' skills/setup-context-driven/SKILL.md` — expected: exits 0.
- `git diff --name-only HEAD | grep -v -E '^(skills/(setup-context-driven|roundfix)/|\.agents/skills/(setup-context-driven|roundfix)/|internal/baseline/assets/|internal/baseline/testdata/|internal/cli/baseline_documentation_contract_test\.go|docs/user-guide/context-driven-development\.md|docs/specs/0082-the-manifest-already-answered-that/task_07\.md)' | grep . ; test $? -eq 1` — expected: exits 0, proving no path outside the authorized boundary changed.
- `go test ./internal/cli/ -run 'TestBaselineExamplesParse|TestBaselineDocumentationContractExamples' -v > /tmp/task_07-1.log 2>&1 && grep -q '^--- PASS: TestBaselineExamplesParse' /tmp/task_07-1.log` — expected: exits 0, proving the documented update examples parse rather than being routed to the interactive parser.
- `go test ./internal/baseline/ ./internal/cli/ ./skills/ -count=1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 9; Project Constraints: Tooling authority.
- `_prd.md` → Project Constraints: Tooling authority.
- ADR-0081.

## Result

Implementation-ready work is present for Daemon Verification. The declared
`## Verification` commands were not run in this Agent turn.

### Implementation

- The durable user guide now routes compatible adopted repositories through
  `roundfix baseline update`, documents all seven flags and exit categories
  `0`, `1`, `2`, `3`, and `130`, and states the byte-identical preservation
  guarantee for every non-managed region.
- The canonical `setup-context-driven` skill now separates managed refresh from
  first adoption and Profile changes, and the canonical `roundfix` skill names
  the update command, flags, result schema, preservation behavior, and exit
  categories. `make skills-sync` regenerated both embedded mirrors.
- The `context-workflow` module's existing `setup-context-driven` dispatch
  trigger now teaches `roundfix baseline update`; no Normative Clause,
  decision, capability, or template changed.
- The published-example contract now dispatches `baseline update` examples to
  `parseBaselineUpdateCommand`, and its help contract covers the shipped update
  usage, flags, result schema, and exit categories.
- `make baseline-digests` regenerated every derived pin and fixture; no digest
  value was transcribed manually.

### Focused checks

- Initial signal: repository search found no `baseline update` occurrence in
  the user guide or either owned skill pair, and
  `parsePublishedBaselineExample` had no update dispatch branch.
- `rtk make skills-sync` — exit `0`; `cmp -s` confirmed both canonical/embedded
  skill pairs are byte-identical.
- `rtk make baseline-digests` — exit `0`, eight regeneration tests and the
  strict catalog validation passed; emitted
  `{"schemaVersion":1,"type":"baseline-digests","ok":true,"changed":true}`.
- `GOCACHE=<worktree>/.gocache rtk go test ./internal/cli -run '^(TestBaselineDocumentationContract|TestBaselineExamplesParse)$' -count=1`
  — 15 subtests passed. The first attempt did not compile because the sandbox
  denied the shared macOS Go cache; the repository-local ignored cache removed
  that environment boundary.
- `GOCACHE=<worktree>/.gocache rtk go test ./skills -run '^(TestThinSetupSkill|TestAuthorialSkillSync)$' -count=1`
  — 19 subtests passed.
- `GOCACHE=<worktree>/.gocache rtk go run -buildvcs=false ./cmd/roundfix baseline update --help`
  — exit `0`; output listed every documented flag, schema
  `roundfix/baseline-update-result/v1`, and exits `0/1/2/3/130`.

### Acceptance-criterion evidence

1. The user-guide option and exit tables are grounded in the focused shipped
   help output, and the guide explicitly preserves Repository-Specific
   Normative Rules, repository-rule blocks, and authored prose byte-for-byte.
2. The setup skill's first recipe is now the dedicated managed-refresh command;
   the interview is reserved for first adoption, Profile changes, or genuinely
   new decisions.
3. The roundfix skill records the Baseline update usage, all flags, result
   schema, preservation behavior, and every exit category observed in shipped
   help.
4. Sanctioned regeneration rewrote the formatter fixture's generated
   `docs/agents/skill-dispatch.md`, whose setup trigger now names
   `roundfix baseline update`.
5. The successful `baseline-digests` result regenerated all catalog, profile,
   formatter, diagnostic, and characterization pins from canonical sources.
6. `TestBaselineExamplesParse` passed after the update branch was added, so
   every published update Bash example was accepted by
   `parseBaselineUpdateCommand` rather than the interactive parser.
7. Changed-path postflight is recorded below after the final Task-file update.

### Changed-path postflight

- `rtk git -c core.fsmonitor=false status --short --untracked-files=all` — exit
  `0`; every changed source path is one of the two authorized skill pairs, the
  command-family module asset, the user guide, the published-example contract
  test, or this Task file. Every remaining changed Baseline profile, formatter,
  diagnostic, catalog, and characterization path is deterministic output named
  by the successful sanctioned regeneration. No other tracked, staged, or
  untracked path was present.
- `rtk git diff --check` — exit `0` after the final implementation edits.

### Verification Feedback repair — attempt 1

- The Daemon's diagnostic artifact identified one failing documentation
  composition assertion: the user guide no longer contained the contiguous
  contract phrase `exact source bytes`.
- Root cause: the managed-refresh clarification retained those words but line
  wrapping split `exact` from `source bytes`, so the established documentation
  contract could not find the phrase. The repair reflows that sentence only;
  it does not change the new update semantics or weaken preservation guidance.
- `GOCACHE=<worktree>/.gocache rtk go test ./internal/cli -run '^TestGuidanceCompositionDocumentation$' -count=1`
  — 3 subtests passed after the repair.
- The failed declared Verification command was not rerun; the Daemon owns its
  one configured retry.
