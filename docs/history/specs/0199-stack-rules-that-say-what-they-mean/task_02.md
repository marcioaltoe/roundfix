---
task: task_02
spec: 0199-stack-rules-that-say-what-they-mean
status: completed
type: backend
complexity: medium
---

# Task 02: The core rules hold in a repository with several languages

## Overview

Two clauses of the `core` module and one Skill Activation assume one language. The dependency clause binds "the repository's declared package manager", which a repository with Bun and Go cannot satisfy with one tool. The testing Skill Activation fires for "Writing or changing tests." and dispatches a Vitest skill, so a Go test change dispatches it too. Nothing says which governs when a dispatched skill's default disagrees with a rule. This Task applies the two clause changes and the trigger change the TechSpec gives, and adds the check that refuses a trigger that dispatches a technology skill without naming the technology.

This is an authorized tooling Task. It may change only the files in its Context, the derived pins the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST apply, in `internal/baseline/assets/modules/core.json`, the replacement the TechSpec's "Exact texts" gives for `clause.core.follow-dependency-workflow` and the appended sentence it gives for `clause.core.activate-matching-skills`. Both clauses stay `mandatory`, and every other clause stays byte-identical.
2. MUST change, in `internal/baseline/assets/skill-activations.json`, the `when` of `trigger.testing` to the TechSpec's text, and raise the file's `version` by one. MUST leave every other activation and every bundle as it is, including the owner of `trigger.production-code` and `trigger.debugging`.
3. MUST raise by one, from the value on the starting main, the versions the TechSpec's "Version changes" lists for task_02.
4. MUST create `internal/baseline/stack_core_wording_test.go` with `languageTriggerFindings` and the six tests the TechSpec's Testing Approach 2 names. It MUST reuse `stackWordingFindings` and `clauseForce` from task_01's file without editing that file.
5. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a
   golden or a generated guide, and MUST keep the module file's existing
   formatting by replacing strings and version numbers in place.
6. MUST NOT add a clause. The rule-over-skill-default sentence joins the existing skill clause.
7. MUST NOT edit any skill, any setup, or an existing test file, and MUST change no exported function signature.

## Subtasks

- [ ] Apply the two clause changes and the trigger change, and raise the versions.
- [ ] Create the trigger check, the composition test and the wording and force tests.
- [ ] Regenerate the pins, the two goldens, this repository's two guides and the Setup Manifest.

## Acceptance Criteria

- [ ] `docs/agents/agent-instructions.md` binds each language's declared package manager, and `docs/agents/skill-dispatch.md` ends the skill clause with the rule-over-skill-default sentence.
- [ ] A composition of the `core`, `go` and `typescript` modules renders the testing trigger for TypeScript tests and keeps the Go testing trigger; a composition of `core` and `go` renders neither the Vitest skill nor the testing trigger.
- [ ] A trigger that dispatches the Vitest skill without naming TypeScript is reported by the trigger check.
- [ ] Both reworded clauses keep their identifier and stay `mandatory`.
- [ ] A second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0190-a-stack-rule-names-the-language-or-workspace-it-governs.md`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/skill-activations.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/skill-dispatch.md`
- creates: `internal/baseline/stack_core_wording_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheCoreGuidesNameEachLanguageAndTheGoverningRule|TestTheCoreWordingCheckReportsTheWordingItReplaced|TestTheRewordedCoreClausesKeepTheirForce|TestEveryTriggerThatDispatchesALanguageSkillNamesItsLanguage|TestATestingTriggerThatNamesNoLanguageIsReported|TestAComposedProfileDispatchesVitestOnlyForTypeScriptTests|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheCoreGuidesNameEachLanguageAndTheGoverningRule TestTheCoreWordingCheckReportsTheWordingItReplaced TestTheRewordedCoreClausesKeepTheirForce TestEveryTriggerThatDispatchesALanguageSkillNamesItsLanguage TestATestingTriggerThatNamesNoLanguageIsReported TestAComposedProfileDispatchesVitestOnlyForTypeScriptTests TestFormatterComposition TestCatalogCompatibility TestBaselinePlanCharacterization; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && golden=internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents && for pair in "docs/agents/agent-instructions.md|Use each language's declared package manager and lockfile workflow." "docs/agents/skill-dispatch.md|Repository-Specific Normative Rule, follow the rule." "$golden/agent-instructions.md|Use each language's declared package manager and lockfile workflow." "$golden/skill-dispatch.md|Repository-Specific Normative Rule, follow the rule." "$golden/skill-dispatch.md|Writing or changing TypeScript tests."; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && for pair in "docs/agents/agent-instructions.md|Use the repository's declared package manager" "$golden/agent-instructions.md|Use the repository's declared package manager" "$golden/skill-dispatch.md|Writing or changing tests."; do file="${pair%%|*}"; phrase="${pair#*|}"; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase"; then printf 'stale phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; fi; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task none of the six new named tests exists and the guides still carry the two replaced sentences, so the command fails.

## References

- [_techspec.md](_techspec.md) — Exact texts; Version changes; Testing Approach 2; Testing Approach 5; Build Order 2
- `_prd.md` → Goal 2; Goal 4; Goal 5; Core Feature 2; Success Metric 2; Success Metric 3; Success Metric 7; Success Metric 8
- `_techspec.md` → API Contract 2
- ADR-0067, ADR-0073, ADR-0081, ADR-0103, ADR-0149, ADR-0186, ADR-0190

## Result

Implemented the task_02 slice for Daemon Verification. The dependency clause
now binds each language's package manager; the existing skill clause ends with
the rule-over-skill-default sentence; and `trigger.testing` names TypeScript.
The core module, both affected rules, both affected guides and activation file
versions each increased by one. No clause was added and no skill was edited.

Added the six authored tests in `stack_core_wording_test.go`, reusing
`stackWordingFindings` and `clauseForce` without changing task_01's test file.
`languageTriggerFindings` follows each activation's bundle and reports its
trigger identifier when a mapped technology skill lacks its technology in
the trigger text. The coverage assertion counts mapped skills reached by real
activations, including Vitest, React and Hono.

### Focused checks

All Go commands used `GOCACHE=/private/tmp/roundfix-0199-task02-gocache`.

- Before source edits, `rtk proxy go test -count=1 -run
  'Test(TheCore|EveryTrigger|ATestingTrigger|AComposedProfile|TheRewordedCore)'
  ./internal/baseline` exited 1: rendered core wording, `trigger.testing`
  scope and the mixed Go/TypeScript dispatch exposed the old statements.
- After source edits and digest regeneration, the same focused selection with
  `-v` exited 0; all six new tests passed, including both composition cases.
- `rtk make baseline-digests` exited 0 and reported `ok: true, changed: true`.
  It regenerated the two formatter goldens and declared derived pins and
  snapshots. No generated file was hand-edited.
- `rtk proxy go run -buildvcs=false ./cmd/roundfix baseline update --repo .
  --no-skills --yes --format text` initially encountered a sandbox restriction
  opening the Git-private transaction lock. The authorized elevated rerun
  exited 0, verified approved postimages and updated exactly the two managed
  guides and Setup Manifest. A second invocation exited 0, reported
  `File changes: 0` and `Idempotence: verified`.
- A read-only JSON comparison against HEAD confirmed unchanged clause
  identities, exactly two changed clauses with `mandatory` force, all six
  version increments, unchanged bundles and unchanged other activations.
- Changed-file postflight found 17 paths, all within this Task's Context,
  created test file and assigned Task file. `git diff --check` exited 0.
- The repository-selected `rtk make verify-incremental` was interrupted by
  the sandbox blocking network access to `cafe.github.com`. Its authorized
  elevated rerun exited 0: formatting, vet, package tests, skill synchronization
  checks, skill validation and the CLI build passed. This incremental check
  does not replace the Daemon's declared Verification.

### Acceptance evidence

| Criterion | Implementation and focused evidence |
| --- | --- |
| Generated core guides bind each language and end the skill clause with rule precedence | Managed Refresh verified both postimages; guide diff inspection and `TestTheCoreGuidesNameEachLanguageAndTheGoverningRule` confirm exact wording. |
| Mixed composition scopes Vitest to TypeScript; Go-only composition omits it | `TestAComposedProfileDispatchesVitestOnlyForTypeScriptTests` passed both compositions and preserves the Go testing trigger in each. |
| Unscoped Vitest trigger is reported | `TestATestingTriggerThatNamesNoLanguageIsReported` returned exactly `trigger.testing`; the embedded activation sweep passed with three covered technology skills. |
| Both clause identifiers and mandatory force survive | `TestTheRewordedCoreClausesKeepTheirForce` passed for both identifiers; the JSON comparison confirmed no added or removed clauses. |
| Second Managed Refresh is a no-op | Second update reported zero file changes and verified idempotence at exit 0. |

The declared `## Verification` command was not run. Task status and settlement
remain Daemon-owned; no commit, push or Pull Request was made. The starting
worktree already contained the Daemon's `pending` to `in_progress` status
change in this Task file. No follow-up implementation was added.
