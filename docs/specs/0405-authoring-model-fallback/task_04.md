---
task: task_04
spec: 0405-authoring-model-fallback
status: pending
type: backend
complexity: medium
---

# Task 04: The Baseline tells each root instruction file how an Authoring Phase keeps its model

## Overview

The Baseline's root block `root.spec-workflow`, which every built-in profile
writes into an adopting repository's `AGENTS.md`, gains one bullet scoped to
Claude Code: in Claude Code, a skill's `model` and `effort` header applies
only when a person types the skill; a skill an Agent activates runs on the
session's model; the Agent hands an Authoring Phase to its Author Subagent
and, on that model's usage limit, follows `roundfix phase fallback`. The
bullet says it does not apply to Codex or OpenCode, so `AGENTS.md` stays the
one instruction file and `CLAUDE.md` its link. The versions rise, the derived files are
regenerated, and this repository's `AGENTS.md` and Setup Manifest receive the
bullet through the Managed Refresh. No user-level file is read or written.

## Requirements

1. MUST change `internal/baseline/assets/templates/root/spec-workflow.md` to the text of `_techspec.md` → API Contract 10 → The template, word for word and with its line breaks. Its tokens stay `reference.spec-routing` and `reference.issue-tracker`.
2. MUST raise by one, in `internal/baseline/assets/templates/index.json`, the `version` of `template.root.spec-workflow`, and, in `internal/baseline/assets/modules/spec-workflow.json`, the `version` of the root block `root.spec-workflow`. MUST keep every id, token, reference and the empty `rules` list. MUST add no root block, clause or rule.
3. MUST regenerate, in this order and never by hand, with the commands of `_authorization.md` → Sanctioned regeneration:
   - the module record (`go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`), which writes the module's next version; the Spec names no number;
   - the derived Baseline files (`make baseline-digests`);
   - as the Task's last step, this repository's managed carriers (`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`), whose second run reports `current`.

   The files these commands rewrote in the measurement of `references/2026-11-18-root-statement-decision-and-measurement.md` are declared below; any other file they rewrite on the tree the Task runs on is recorded by the Daemon.
4. MUST add `internal/baseline/authoring_root_statement_test.go` with `TestTheSpecWorkflowRootBlockStatesHowAnAuthoringPhaseKeepsItsModel`, calling `t.Parallel()` first, with these subtests. Each compares whole sentences of the bullet, normalized to single spaces, never a word alone: the first sentence, which opens with "In Claude Code,", and the last, "This bullet does not apply to Codex or OpenCode, whose models come from the Agent Selection Profiles."
   - `plan_postimage`: a plan built on a temporary repository with the existing plan seams of `internal/baseline/plan_test.go` writes each of the two sentences exactly once into the `AGENTS.md` postimage, between the `root.spec-workflow` begin and end markers;
   - `formatter_golden`: the Standard TypeScript formatter golden `AGENTS.md` holds each of the two sentences exactly once;
   - `scoped_to_claude_code`: in the template, the bullet's first words are "In Claude Code,".
5. MUST NOT change, by hand or by command, any Source Baseline, Setup Snapshot, other module, other root block, guide template, profile other than `standard-typescript-monorepo.json`, any guide under `docs/agents/` other than `docs/agents/setup-context.json`, `CLAUDE.md`, `.roundfixrc.yml`, the Makefile, a workflow, `go.mod` or a lint configuration. Of the Governed Paths the Task changes only those it declares below, all listed in `_authorization.md`.
6. MUST NOT read or write any user-level file, such as `~/.claude/CLAUDE.md`, and no test may read the real home directory.

## Subtasks

- [ ] Add the bullet to the root block template and raise the template and root block versions.
- [ ] Add the root statement test.
- [ ] Record the module version and regenerate the derived files.
- [ ] Refresh this repository's managed carriers and confirm a second refresh reports `current`.

## Acceptance Criteria

- [ ] A plan for a built-in profile writes the statement once into `AGENTS.md`, inside the `root.spec-workflow` markers, and the Standard TypeScript formatter golden holds it once.
- [ ] This repository's `AGENTS.md`, which `CLAUDE.md` links to, holds the statement, and a second managed refresh reports `current`.
- [ ] The block keeps its id, tokens and references; no other block, clause, guide, Source Baseline, Setup Snapshot or profile changes.
- [ ] Removing the bullet from the template makes the three subtests fail, and removing its opening words "In Claude Code," or its last sentence makes at least one fail.

## Context

- instruction: `docs/adr/0431-the-baseline-tells-each-root-instruction-file-how-an-authoring-phase-keeps-its-model.md`
- instruction: `docs/specs/0405-authoring-model-fallback/references/2026-11-18-root-statement-decision-and-measurement.md`
- instruction: `internal/baseline/plan_test.go`
- instruction: `internal/baseline/repository_rule_placement_test.go`
- interface: `internal/baseline/assets/templates/root/spec-workflow.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- creates: `internal/baseline/authoring_root_statement_test.go`
- interface: `AGENTS.md`
- interface: `docs/agents/setup-context.json`

## Consumers

- The root block template, the module and the derived catalog records — consumers: `./internal/baseline`, `./internal/cli`
- This repository's root instruction file and Setup Manifest — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheSpecWorkflowRootBlockStatesHowAnAuthoringPhaseKeepsItsModel|TestEveryBaselineModuleVersionIsRecorded)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheSpecWorkflowRootBlockStatesHowAnAuthoringPhaseKeepsItsModel/plan_postimage TestTheSpecWorkflowRootBlockStatesHowAnAuthoringPhaseKeepsItsModel/formatter_golden TestTheSpecWorkflowRootBlockStatesHowAnAuthoringPhaseKeepsItsModel/scoped_to_claude_code TestEveryBaselineModuleVersionIsRecorded; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/baseline` — expected: exit 0. Before this Task the root statement test does not exist, so its pass lines are missing; afterwards the plan and the formatter golden carry the statement once, the module version is recorded, and the whole Baseline suite passes.
- `tr -s '[:space:]' ' ' < AGENTS.md | grep -qF -- "In Claude Code, a skill" || { echo 'AGENTS.md lacks the scoped header rule' >&2; exit 1; }; tr -s '[:space:]' ' ' < AGENTS.md | grep -qF -- "This bullet does not apply to Codex or OpenCode" || { echo 'AGENTS.md lacks the harness scope' >&2; exit 1; }; test -L CLAUDE.md || { echo 'CLAUDE.md is not a link' >&2; exit 1; }; tr -s '[:space:]' ' ' < CLAUDE.md | grep -qF -- "hand the phase to its Author Subagent" || { echo 'CLAUDE.md lacks the hand-over' >&2; exit 1; }; f="$(mktemp)"; go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format json > "$f"; status=$?; grep -q '"state":"current"' "$f" || { cat "$f"; rm -f "$f"; exit 1; }; rm -f "$f"; test "$status" -eq 0 && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract && go test -count=1 -timeout 30m ./internal/cli` — expected: exit 0. Before this Task this repository's root file lacks the statement; afterwards `AGENTS.md`, and `CLAUDE.md` through its link, carry it, a managed refresh reports `current`, and the documentation contracts and the whole CLI suite pass.

## References

- `_prd.md` → Goals (sixth bullet); User Story 8; Core Feature 9; Success Metric 6; Non-Goals; Decisions (2026-11-18, both)
- `_techspec.md` → Root instruction statement; API Contract 10; Testing Approach; Build Order 4; Risks & Considerations (other harnesses)
- [`references/2026-11-18-root-statement-decision-and-measurement.md`](references/2026-11-18-root-statement-decision-and-measurement.md)
- ADR-0431, ADR-0401, ADR-0186
