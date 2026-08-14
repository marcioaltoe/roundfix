---
task: task_06
spec: 0085-what-an-agent-reads-before-it-decides
status: completed
type: infra
complexity: medium
---

# Task 06: Make the Secondbrain consultation unconditional

## Overview

Consulting the Secondbrain before design and architecture decisions, and before
writing an Idea, PRD, or TechSpec, is currently conditional wording in the
Baseline catalog. This Task makes it unconditional at those two named moments and
renders it into the two guides that carry it. Authorized tooling work with exact
bounded files.

## Requirements

1. MUST make the consultation clause unconditional at the two named moments:
   before a design or architecture proposal, and before authoring an Idea, PRD,
   or TechSpec.
2. MUST bind that clause to the authoring stages in the workflow module, so the
   obligation reaches the stage that must obey it.
3. MUST render the change into the two setup-owned guides through their managed
   regions, never by editing the rendered text directly.
4. MUST run the sanctioned regeneration so derived Baseline pins match; those
   pins are sanctioned fallout under ADR-0081, not separate targets.
5. MUST break the conditional-clause case Task 01 declared, and update it in the
   same commit.
6. MUST NOT touch any path outside the bounded list below.

## Subtasks

- [ ] Make the clause unconditional in the catalog.
- [ ] Bind it to the authoring stages.
- [ ] Render and regenerate the derived pins.

## Acceptance Criteria

- [ ] The clause states the obligation without a condition.
- [ ] It is bound to the Idea, PRD, and TechSpec authoring stages.
- [ ] Both guides render the new text through their managed regions.
- [ ] Derived pins match after regeneration.

## Bounded scope

Authorized by
`docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md`.
This Task may create or modify only:

- `internal/baseline/assets/modules/secondbrain.json`, limited to the
  consultation trigger clause
- `internal/baseline/assets/modules/context-workflow.json`, limited to binding
  that consultation to the Idea, PRD, and TechSpec authoring stages
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/secondbrain.md`
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`,
  limited to the clause entries the edits above require
- `docs/agents/secondbrain.md`
- `internal/baseline/testdata/**` and `docs/references/coverage-record.json`, only
  as sanctioned regeneration fallout under ADR-0081
- `docs/specs/0085-what-an-agent-reads-before-it-decides/task_06.md`

Any other path is out of scope; stop and fail the Task rather than widen it.

## Verification

- `grep -q 'secondbrain' internal/baseline/assets/modules/context-workflow.json` — expected: exits 0, proving the consultation is bound in the workflow module.

## Context

- instruction: `docs/workflow/authorizations/2026-08-09-what-an-agent-reads-before-it-decides.md`

A guard asserting the conditional wording is absent is deliberately absent too:
the module carries no such wording today, so the check passes before any work.
Requirement 1 is what obliges the unconditional clause.

Asserting the layout characterization still passes is deliberately absent: this
Task changes no observable path, so the case passes before and after, and a
command that cannot fail cannot prove anything. Keeping it passing is the
Run-level gate's job.

Whole-package sweeps, `go build`, `go clean -testcache` and `make verify` are
deliberately absent: each passes against a tree where no work has happened, so
it approves the Task before it starts. Regression is the Run-level gate's job.

## References

- `_prd.md` → the consultation obligation.
- `_techspec.md` → Build Order 6.
- ADR-0081.

## Result

Implemented the unconditional Secondbrain consultation clause at the two
authorized moments and bound `knowledge-workspace` activation to Idea, PRD,
and TechSpec authoring through an explicit Secondbrain authoring trigger.
Regenerated the source-baseline manifest, catalog pins, plan fixtures, and
rendered guide fixtures with the sanctioned Baseline digest target.

Focused checks:

- Exact `jq -e` assertions for the consultation clause and authoring trigger:
  passed.
- `cmp -s docs/agents/secondbrain.md internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`:
  passed; the repository guide matches the generated managed guide.
- `GOCACHE=/private/tmp/roundfix-task06-gocache go test ./internal/baseline -run '^(TestCatalogCompatibility|TestFormatterComposition|TestBaselineCompatibilityCorpus)$' -count=1`:
  passed. The first attempt used the host Go cache and was blocked by sandbox
  permissions; the writable-cache retry passed.
- `make baseline-digests`: passed and regenerated the derived artifacts. A
  second run with the writable Go cache reported `changed:false`, proving the
  derived pins match their canonical sources.

Acceptance evidence:

- The catalog clause now requires consultation before a design or architecture
  proposal and before authoring an Idea, PRD, or TechSpec, with no
  local-documentation condition.
- `context-workflow` now carries
  `trigger.context-workflow.knowledge-workspace-secondbrain-authoring` for
  Secondbrain consultation before Idea, PRD, and TechSpec authoring.
- The generated Secondbrain and skill-dispatch guide fixtures carry the new
  clause and authoring-stage binding, and the repository Secondbrain managed
  region matches its generated fixture.
- Sanctioned regeneration is idempotent after the edit, so the derived pins
  match.

Verification feedback repair, attempt 1:

- Inspected the Daemon diagnostic artifact; it exists but contains no log
  bytes. The Daemon report identifies the failed probe and exit status.
- A pre-repair structured JSON check confirmed that `context-workflow` had no
  lowercase `secondbrain` token: the binding used only the generic
  `knowledge-workspace` name, so it could not satisfy the authored probe.
- Renamed the binding to
  `trigger.context-workflow.knowledge-workspace-secondbrain-authoring` and made
  its activation text explicitly require consulting the Secondbrain before
  authoring an Idea, PRD, or TechSpec.
- Post-repair `jq -e` assertions confirmed both the lowercase `secondbrain`
  token and the exact authoring trigger. The generated skill-dispatch guide
  contains the same trigger.
- The focused Baseline contract tests passed, and the post-repair sanctioned
  regeneration rerun reported `changed:false`.
