---
task: task_02
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
status: completed
type: backend
complexity: medium
---

# Task 02: Go, CLI and TUI rules carry force

## Overview

The Go, CLI and TUI modules ship their rules as paragraphs with no force. This Task splits each into clauses with a stated force, adds the recorded-reason Go dependency preference, the cross-platform build for build-constrained files and the CLI clause that ships a command skill with the behavior it describes, and moves the three modules to the clause-bearing module schema. It is verified through the Go CLI/TUI profile's rendered guides, which this repository also carries.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST replace the rules of `internal/baseline/assets/modules/go.json`, `cli-surface.json` and `tui-surface.json` with the clauses the TechSpec's "Exact texts" gives for task_02, keeping every rule identifier, and MUST move each module's `schemaVersion` to `setup-context-driven/module-v3` with `"repositoryExtensions": []`.
2. MUST NOT write a Go clause that forbids a named library, and MUST NOT edit any skill, the Go guide template or this repository's Repository-Specific Normative Rules.
3. MUST raise by one, from the value on the starting main, the versions the TechSpec's "Version changes" lists for task_02, leaving the modules' skill lists as they are.
4. MUST add the task_02 clauses to the force record, and MUST add `docs/agents/cli.md` and `docs/agents/tui.md` to the list of guides that grew past the frozen parity record in `internal/baseline/plan_test.go`, changing no other line of that file. Spec 0207 already lists `docs/agents/go.md` there; MUST NOT add it a second time.
5. MUST create `internal/baseline/stack_force_go_cli_tui_test.go` with `ruleLevelGuidanceFindings`, `namedLibraryBanFindings` and the six tests the TechSpec's Testing Approach 2 names; each negative test feeds its check a literal.
6. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`;
   a second refresh MUST report `File changes: 0`. MUST NOT hand-edit a snapshot, a golden or a generated guide.
7. MUST change no exported function signature and rename or remove no top-level test.

## Subtasks

- [ ] Split the Go, CLI and TUI rules into clauses with force.
- [ ] Raise versions and move the schema.
- [ ] Update the force record and the parity list; create the new test file.
- [ ] Regenerate and refresh twice.

## Acceptance Criteria

- [ ] No rule of the Go, CLI or TUI module carries rule-level guidance, and a literal rule-level rule is reported.
- [ ] The rendered Go, CLI and TUI guides label every rule line with its force and state the recorded-reason exception, the cross-platform build and the skill-ships-with-behavior clause.
- [ ] The rendered Go guide names no library it forbids, and a literal ban on Cobra is reported.
- [ ] A second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0202-a-stack-rule-carries-force-and-a-preference-names-what-licenses-the-exception.md`
- instruction: `docs/adr/0190-a-stack-rule-names-the-language-or-workspace-it-governs.md`
- interface: `internal/baseline/assets/modules/go.json`
- interface: `internal/baseline/assets/modules/cli-surface.json`
- interface: `internal/baseline/assets/modules/tui-surface.json`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/go.md`
- interface: `docs/agents/cli.md`
- interface: `docs/agents/tui.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/stack_force_go_cli_tui_test.go`

## Scope after Spec 0207

The composed profile `go-cli-typescript-monorepo` that Spec 0207 shipped selects the Go and CLI modules, so the clauses this Task writes also render in its `go.md` and `cli.md`, which already open with the Go scope sentence. No file of that profile changes: it has no formatter golden and no Source Baseline. A rehearsal on 2026-10-02, on top of task_01 at `18ef15eb`, applied this Task's module, force-record and parity-list changes, regenerated, converged on a second Managed Refresh, and passed the Baseline and skills tests, including the composed profile's repeated-clause and convergence tests.

## Verification

- `out="$(go test -count=1 -v -run "^(TestEveryGoCLIAndTUIRuleCarriesForce|TestARuleWithoutForceIsReported|TestTheGoCLIAndTUIGuidesStateTheirClauses|TestTheGoGuideBansNoLibraryByName|TestALibraryBanInTheGoGuideIsReported|TestTheGoCLIAndTUIClausesCarryTheirForce|TestPlanDeterminismMatchesMaintainedManagedEntryFixture|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestEveryGoCLIAndTUIRuleCarriesForce TestARuleWithoutForceIsReported TestTheGoCLIAndTUIGuidesStateTheirClauses TestTheGoGuideBansNoLibraryByName TestALibraryBanInTheGoGuideIsReported TestTheGoCLIAndTUIClausesCarryTheirForce TestPlanDeterminismMatchesMaintainedManagedEntryFixture TestBaselineClauseForceIsCharacterized; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && for pair in "docs/agents/go.md|**mandatory**: Prefer the Go standard library." "docs/agents/go.md|record that reason in an ADR or another recorded repository decision" "docs/agents/go.md|for every operating system its constraints name" "docs/agents/go.md|**prohibited**: Do not hand-edit" "docs/agents/cli.md|updates that skill or reference in the same pull request." "docs/agents/cli.md|**mandatory**: Make automation deterministic and non-interactive." "docs/agents/tui.md|**mandatory**: Use terminal emulation only when"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && if tr -s '[:space:]' ' ' < docs/agents/go.md | grep -qF -- "Use stdlib"; then printf 'stale Go testing sentence\n' >&2; exit 1; fi && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the new tests do not exist and the Go guide carries no force labels, so the command fails.

## References

- `_prd.md` → Goal 1; Goal 2; Story 1; Story 2; Core Feature 1; Success Metric 1; Success Metric 2
- `_techspec.md` → Candidate rules; Exact texts (task_02); Version changes; Existing tests that change; API Contract 2; Testing Approach 2; Testing Approach 5; Build Order 2
- ADR-0059, ADR-0081, ADR-0149, ADR-0186, ADR-0190, ADR-0202

## Result

Implemented the task_02 slice for Daemon Verification. The Go, CLI and TUI
modules now carry the TechSpec's exact 17 clauses (nine Go, five CLI, three
TUI), use module-v3 with empty repository extensions, and retain every rule
identifier and skill list. Module, guide and rule versions each rose by one:
Go module 4 → 5, Go guide 3 → 4, and all other affected versions 2 → 3.
The force record gained all 17 clauses. Only the CLI and TUI paths were added
to the frozen-parity exception list; the existing Go entry was retained.

The six new tests exercise embedded modules and rendered temporary-adopter
plans. Both negative tests feed literal input to their pure check. No skill,
Go guide template, repository-specific rule, production signature, existing
top-level test, Task Graph or other Task file was changed.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| No rule-level guidance; a literal rule is reported | `TestEveryGoCLIAndTUIRuleCarriesForce` and `TestARuleWithoutForceIsReported` passed. Before changing the modules, the first test failed and identified all six original rule-level paragraphs. |
| Rendered guides label every rule and state the three obligations | `TestTheGoCLIAndTUIGuidesStateTheirClauses` passed against all 17 exact texts, rejecting the old Go sentences and any unlabelled rule line. `TestTheGoCLIAndTUIClausesCarryTheirForce` passed against the literal force table. The Managed Refresh rendered the recorded-reason preference, constrained-platform build and skill-ships-with-behavior obligation into this repository's guides. |
| Go guide bans no named library; literal Cobra ban is reported | `TestTheGoGuideBansNoLibraryByName` passed for Cobra, testify and Viper; `TestALibraryBanInTheGoGuideIsReported` passed with the literal `Do not use Cobra.`. The repository's own library prohibition remains unchanged. |
| Second Managed Refresh is a no-op | First confirmed refresh exited 0 with four file changes: Go, CLI and TUI guides plus Setup Manifest. The second exited 0 with `File changes: 0`, `approved Baseline Plan is already applied`, and `Idempotence: verified`. |

### Checks run

- `go test ./internal/baseline -run '^Test(EveryGoCLIAndTUIRuleCarriesForce|ARuleWithoutForceIsReported|ALibraryBanInTheGoGuideIsReported)$' -count=1`
  — expected red before the module changes: all three modules reported their
  rule-level guidance.
- `GOCACHE=/private/tmp/roundfix-task02-go-cache go test ./internal/baseline -count=1 -v -run '^Test(EveryGoCLIAndTUIRuleCarriesForce|ARuleWithoutForceIsReported|TheGoCLIAndTUIGuidesStateTheirClauses|TheGoGuideBansNoLibraryByName|ALibraryBanInTheGoGuideIsReported|TheGoCLIAndTUIClausesCarryTheirForce)$'`
  — exit 0; all six tests reported PASS on the regenerated tree.
- `rtk make baseline-digests` — exit 0, `ok: true`, `changed: true`;
  regenerated only the catalog digest, normalized catalog and four declared
  plan-characterization goldens. No generated artifact was hand-edited.
- `GOCACHE=/private/tmp/roundfix-task02-go-cache go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
  — run twice with sandbox escalation for the Git-private transaction lock;
  both exited 0, with four then zero file changes. Skills were skipped.
- `GOCACHE=/private/tmp/roundfix-task02-go-cache rtk make verify-incremental`
  — sandboxed attempt exited 2 because two CLI force-stop tests could not
  read the process table. The rerun with sandbox escalation exited 0,
  including format checks, vet, package tests, skill checks and build.
- Source-contract inspection — confirmed exact authored texts, schema,
  version increments, retained rule IDs and unchanged skills/root blocks.
- Scope postflight and `git diff --check` — all 17 changed/untracked paths
  belong to task_02's Context or assigned Task file; no whitespace errors.

The first unscoped-cache refresh attempt could not read the host Go cache;
the task-scoped cache resolved that restriction. The next sandboxed refresh
could not open the Git-private transaction lock; escalation resolved it.
Both successful refreshes retained the reported nested-carrier warnings for
the formatter fixture and Source Baseline corpus; those carriers were not
changed. No implementation follow-up was needed.

Task status remains Daemon-owned. The authored Verification command was not
run; no terminal Task verdict, commit, push or pull request is claimed.
