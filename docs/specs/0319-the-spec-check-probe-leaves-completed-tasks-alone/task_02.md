---
task: task_02
spec: 0319-the-spec-check-probe-leaves-completed-tasks-alone
status: pending
type: docs
complexity: medium
---

# Task 02: Document the probe's selection rule

## Overview

The `spec` guide and the Roundfix skill reference describe
`--run-verification` without saying which Tasks it probes, and the `qa-gate`
skill justifies its non-probing precondition with a sentence this Spec makes
false. This Task documents the selection rule and the `not probed` report,
replaces that sentence, amends ADR-0148 and revises the **Vacuous
Verification** glossary entry, then records both skills' versions and the
Behavior Surface fingerprint.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-25-the-spec-check-probe-re-runs-completed-tasks.md`](references/2026-10-25-the-spec-check-probe-re-runs-completed-tasks.md)
   of 2026-10-25 (`_prd.md` → Sources 1) by writing the documentation of
   `_techspec.md` → Interfaces.
2. MUST append the guide paragraph of `_techspec.md` → Interfaces to
   `docs/user-guide/commands/spec.md` and to the canonical
   `.agents/skills/roundfix/references/spec.md`, each directly after the
   paragraph ending "If the parser cannot start, the existing command probe
   still runs.".
3. MUST replace, in the canonical `.agents/skills/qa-gate/SKILL.md`, the
   paragraph after the `roundfix spec check <slug> --strict` block of the
   precondition bullet under `## Surface Transcripts at the gate` with the
   paragraph of `_techspec.md` → Interfaces. Never add or change text inside
   any `### QA settlement` section.
4. MUST append the ADR-0148 amendment of Interfaces to
   `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`
   and set its `updated_at` to `2026-10-26T00:00:00Z`, changing nothing else
   in that file.
5. MUST revise the **Vacuous Verification** entry of `CONTEXT.md` through the
   `domain-modeling` skill by adding the sentence of Interfaces after its
   second sentence, keeping its first two sentences and its `_Avoid_` line.
6. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   (it raises both version fields of the `roundfix` and `qa-gate` skills and
   records them), then `make skills-sync` again so the mirrors match.
7. MUST, as the Task's last step, after every guide and skill edit,
   re-record the Behavior Surface fingerprint ("re-record after the last
   guide edit") with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
8. MUST NOT change any Go source, any other skill, `docs/agents/`, any
   Baseline module or derived Baseline file, the Makefile or the CI
   workflows.

## Subtasks

- [ ] Add the guide and reference paragraph and replace the `qa-gate` sentence.
- [ ] Amend ADR-0148 and revise the glossary entry.
- [ ] Sync and record both skills, then re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] The guide, the canonical reference and its mirror state the selection rule and the `not probed` report.
- [ ] The `qa-gate` skill no longer says probing completed Tasks reports them vacuous, in both copies.
- [ ] ADR-0148 carries the dated amendment and `CONTEXT.md` the revised **Vacuous Verification** entry; skill versions and the Behavior Surface record are current.

## Consumers

None. This Task changes documents, skill text and records only; no Go
package reads them as a contract except the documentation contracts it runs.

## Context

- instruction: `skills/skills_test.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `go test -timeout 10m -count=1 ./skills && go test -timeout 15m -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'its command passing is the postcondition its work proved' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Vacuous Verification**' && for guide in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$guide" | grep -qF -- 'The probe runs only the commands of Tasks a Run would still execute' || { printf 'missing paragraph: %s\n' "$guide" >&2; exit 1; }; done && for skill in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md; do tr -s '[:space:]' ' ' < "$skill" | grep -qF -- 'The probe leaves completed Tasks unprobed' || { printf 'missing qa-gate sentence: %s\n' "$skill" >&2; exit 1; }; if tr -s '[:space:]' ' ' < "$skill" | grep -qF -- 'Probing completed Tasks would therefore report'; then printf 'stale qa-gate sentence: %s\n' "$skill" >&2; exit 1; fi; done && tr -s '[:space:]' ' ' < docs/adr/0148-one-verification-prober-refuses-vacuous-tasks-before-agent-and-at-authoring.md | grep -qF -- 'Amendment (2026-10-26): both callers probe the same Tasks'` — expected: exit 0. Before this Task the glossary sentence, the guide paragraph, the new `qa-gate` sentence and the ADR amendment are absent, so the first phrase check fails. Afterwards the owned-skill versions and mirrors are recorded, the Behavior Surface record is current, and every phrase is present while the stale `qa-gate` sentence is gone.

## References

- `_prd.md` → Sources; Goal 5; Core Feature 3; Success Metric 6; Glossary
- `_techspec.md` → Current behavior 6; Interfaces; ADR Alternatives; Build Order 1
