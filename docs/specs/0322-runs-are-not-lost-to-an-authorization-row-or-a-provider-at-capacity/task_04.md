---
task: task_04
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
status: pending
type: docs
complexity: medium
---

# Task 04: The skills, guides and glossary say both rules

## Overview

The QA, archive and Roundfix skills, the settle, archive, qa-report and
context-driven guides and the glossary tell agents to write the network
sentence exactly, and the usage and configuration guides and the Roundfix
`profiles` reference say a Provider System Error falls back only before the
First Handoff or while a QA report is pending. This Task writes the subject
form and the fall-through into each of them, updates the skill contract test
that pins the offline-row section, records the owned-skill versions and
re-records the Behavior Surface record last.

## Requirements

1. MUST answer both Backlog Entries of 2026-10-26,
   [`references/2026-10-26-an-outside-evidence-row-that-names-what-the-authorization-forbids-does-not-qualify.md`](references/2026-10-26-an-outside-evidence-row-that-names-what-the-authorization-forbids-does-not-qualify.md)
   and
   [`references/2026-10-26-a-provider-at-capacity-parks-a-turn-the-fallback-chain-could-finish.md`](references/2026-10-26-a-provider-at-capacity-parks-a-turn-the-fallback-chain-could-finish.md)
   (`_prd.md` → Sources), by implementing API Contracts 7 and 8 of
   `_techspec.md`.
2. MUST replace the `### Rows the Spec's authorization keeps offline` section
   body in `.agents/skills/qa-gate/SKILL.md`,
   `.agents/skills/archive-spec/SKILL.md` and `.agents/skills/roundfix/SKILL.md`
   with API Contract 7's text, and the Network-Prohibited Row sentence of
   `.agents/skills/roundfix/references/archive.md` and `settle.md` with API
   Contract 7's reference sentence. Never add or change text inside any
   `### QA settlement` section.
3. MUST, in `.agents/skills/qa-gate/SKILL.md`, rewrite the outside-evidence
   paragraph of `### Row input declaration` and the outside-evidence bullet of
   `## Decision examples` so the gate records
   `blocked (environment: <subject> prohibited by Spec authorization)`, names
   what the authorization forbids in its words (`network` for a network
   lookup) within 1 to 80 characters and no parenthesis, keeps the
   `outside-evidence row` provenance item and names the authorization text in
   the report; drop the word "exactly" and cite ADR-0335 beside ADR-0276.
4. MUST replace the Network-Prohibited Row sentence of
   `docs/user-guide/commands/settle.md`, `docs/user-guide/commands/qa-report.md`
   and `docs/user-guide/context-driven-development.md` with API Contract 7's
   guide sentence, and the exact network status in
   `docs/user-guide/commands/archive.md`'s list of exempt rows with
   `blocked (environment: <subject> prohibited by Spec authorization)`.
5. MUST replace the Provider System Error sentence of
   `docs/user-guide/usage.md`, `docs/user-guide/configuration.md` and
   `.agents/skills/roundfix/references/profiles.md` with API Contract 8's
   sentence, and add to `docs/user-guide/usage.md`'s Provider System Error
   paragraph the reopened verdict and the `- discarded_verdict: <verdict>`
   line that API Contract 8 and `_techspec.md` → Data Models describe.
6. MUST revise, through `domain-modeling`, the `CONTEXT.md` entries
   **Network-Prohibited Row** (the subject form, its bounds, `network` for a
   network lookup, one exempt kind per row, ADR-0335),
   **Provider System Error** (the next Fallback Selection takes the work in
   any Task or QA turn, reopening a QA verdict the failed turn wrote; only an
   exhausted chain settles it as runtime infrastructure; ADR-0274 and
   ADR-0335), **Fallback Chain** (the exception: a Lost Rollout before the
   First Handoff or in a QA Task before its report, and a Provider System
   Error in any Task or QA turn) and **First Handoff** (after it, only a
   Provider System Error lets a Fallback Selection take over work a turn
   handed back), keeping each entry's `_Avoid_` line.
7. MUST set `networkProhibitedSettlementSection` and
   `networkProhibitedReferenceSentence` in `skills/acceptance_rules_test.go`
   to API Contract 7's texts and change nothing else in that file.
8. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   (it raises the version pair of `qa-gate`, `archive-spec` and `roundfix` and
   records them), and, as the last step after the last guide edit, re-record
   the Behavior Surface fingerprint with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
   Hand-edit none of their outputs.
9. MUST NOT change any Go production file, Baseline module, `docs/agents/`
   guide, ADR, `.roundfixrc.yml`, the Makefile or the CI workflows.

## Subtasks

- [ ] Rewrite the offline-row section, reference sentences and qa-gate rules.
- [ ] Rewrite the guides' Network-Prohibited Row and Provider System Error sentences.
- [ ] Revise the four glossary entries.
- [ ] Update the skill contract test, sync and record the skills, re-record the Behavior Surface record last.

## Acceptance Criteria

- [ ] No skill, guide or glossary entry this Task owns says the authorization "forbids network access" or tells the gate to write the network sentence exactly.
- [ ] The three skills carry API Contract 7's section, and `TestTheSettlementSkillsNameTheNetworkProhibitedRow` passes.
- [ ] The usage, configuration and `profiles` texts carry API Contract 8's sentence, and the owned-skill and Behavior Surface records are current.

## Consumers

None. The Task changes documents, skills and their contract test; no package consumes them at run time.

## Context

- instruction: `docs/adr/0335-a-stop-that-is-not-the-work-s-fault-does-not-cost-the-run.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `skills/owned_skill_versions_test.go`
- instruction: `skills/settlement_guidance_repocontract_test.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/settle.md`
- interface: `.agents/skills/roundfix/references/profiles.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/settle.md`
- interface: `skills/roundfix/references/profiles.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/acceptance_rules_test.go`
- interface: `docs/user-guide/commands/settle.md`
- interface: `docs/user-guide/commands/qa-report.md`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/configuration.md`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `for file in .agents/skills/qa-gate/SKILL.md .agents/skills/archive-spec/SKILL.md .agents/skills/roundfix/SKILL.md .agents/skills/roundfix/references/archive.md .agents/skills/roundfix/references/settle.md skills/qa-gate/SKILL.md skills/archive-spec/SKILL.md skills/roundfix/SKILL.md docs/user-guide/commands/settle.md docs/user-guide/commands/qa-report.md docs/user-guide/commands/archive.md docs/user-guide/context-driven-development.md CONTEXT.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- '<subject> prohibited by Spec authorization)' || { printf '%s lacks the subject form\n' "$file" >&2; exit 1; }; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'authorization forbids network access'; then printf '%s keeps the network-only wording\n' "$file" >&2; exit 1; fi; done; if grep -qF -- 'record the row exactly as' .agents/skills/qa-gate/SKILL.md; then echo 'qa-gate still says exactly' >&2; exit 1; fi; for file in docs/user-guide/usage.md docs/user-guide/configuration.md .agents/skills/roundfix/references/profiles.md skills/roundfix/references/profiles.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'A Provider System Error lets the next Fallback Selection take the Task in any Task or QA turn' || { printf '%s lacks the fall-through sentence\n' "$file" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/usage.md | grep -qF -- '- discarded_verdict: <verdict>' || { echo 'usage guide lacks the discarded verdict line' >&2; exit 1; }; for phrase in '**Network-Prohibited Row**' '**Provider System Error**' '**Fallback Chain**' '**First Handoff**'; do grep -qF -- "$phrase" CONTEXT.md || { printf 'CONTEXT.md lacks %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'reopening a QA verdict the failed turn wrote' || { echo 'Provider System Error entry not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'only a Provider System Error lets a Fallback Selection take over work a turn handed back' || { echo 'First Handoff entry not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'and a Provider System Error in any Task or QA turn, let the next Fallback Selection take the Task' || { echo 'Fallback Chain entry not revised' >&2; exit 1; }` — expected: exit 0. Before this Task every one of these files carries the network-only wording or the First Handoff window, so the first check fails. Afterwards the skills and their mirrors, the guides and the four glossary entries carry both rules.
- `tr -s '[:space:]' ' ' < skills/acceptance_rules_test.go | grep -qF -- 'forbids in 1 to 80 characters without parentheses (ADR-0276, ADR-0335).' || { echo 'skill contract test not updated' >&2; exit 1; }; go test -timeout 15m -count=1 -run '^(TestTheSettlementSkillsNameTheNetworkProhibitedRow|TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills && go test -timeout 15m -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the contract test pins the network sentence, so the first check fails. Afterwards the three skills match the pinned section, the owned-skill versions are recorded, the `### QA settlement` table is unchanged, and the documentation contracts pass with the Behavior Surface record re-recorded last.

## References

- `_prd.md` → Sources; Goal 5; Core Feature 4; Success Metric 7; Glossary
- `_techspec.md` → API Contracts 7-8; Data Models; Testing Approach; Build Order 4; Glossary
- ADR-0335; ADR-0276; ADR-0274
