---
task: task_01
spec: 0233-a-light-tier-on-open-models
status: completed
type: docs
complexity: medium
---

# Task 01: The guides, the Roundfix Skill and the write-tasks skill describe the light tier and the judge's tasks stage

## Overview

Roundfix is about to run every `complexity: low` Task that is not QA and
declares no Governed Path on an open OpenRouter model by default, under a
monthly ceiling, with one escalation to the subscription, and its advisory
judge gains a `tasks` stage that suggests a tier. The repository's skill-sync
rule requires the guides and skills to describe that behavior in the same Pull
Request. This Task writes the record first, so the CLI change of the other
Tasks ships with it. It is verifiable on its own through the documents, their
mirrors and the recorded skill versions.

## Requirements

1. MUST add a light tier section to the configuration guide that states, in
   the PRD's words and ADR-0238's: which Tasks are `light` (Core Feature 1);
   the light selection, its default model and that it carries no reasoning
   effort (Core Feature 2); the key variable read through Roundfix's key
   helper, that only its name reaches OpenCode, and that the generic
   `OPENROUTER_API_KEY` is removed (Core Feature 3); the two User Config keys
   `openrouter.light_models` and `openrouter.implement_monthly_ceiling_usd`
   with their defaults, refusals and Project Config warning (API Contracts 1
   to 3), also as rows of the Key reference; the Light Spend Log, its path
   and fields (Core Feature 5); the skip warning with the exact phrase
   `light tier skipped for Task` and its three reasons (API Contract 4); and
   that every light Task's prompt, the files its agent reads and its
   diagnostics reach OpenRouter and the model's provider.
2. MUST add a light tier section to the Roundfix Skill's `runtime` reference
   that states the tier rule, the derived profile with source `light-tier`,
   the sentence `A light Task that fails Verification escalates once`, the
   Run Event phases `light_tier_skipped` and `light_tier_escalated` with their
   reason codes (API Contracts 4 to 6), and that a one-Run override turns the
   tier off.
3. MUST document the judge's `--stage tasks` in the Roundfix Skill's `spec`
   reference and in the `spec` command reference, including the
   `suggested model-tier` output line, that the suggestion is advisory and
   never read at dispatch, and the new unknown-stage message (API Contract 7,
   Surface Transcripts 1 and 2), and change both usage lines to
   `[--stage <prd|techspec|tasks>]`.
4. MUST add to the write-tasks skill, under its own heading, that a Task with
   `complexity: low` that is not QA and declares no Governed Path runs on the
   light tier, using the exact phrase `runs on the light tier`, and that the
   author may run `roundfix spec judge <slug> --stage tasks` for an advisory
   tier before the planning Pull Request merges.
5. MUST add to the model reference a short note that the light tier's
   default model is `deepseek/deepseek-v4.1-flash` under the `openrouter/`
   selection namespace and that every light model passes the subscription
   rule; the snapshot marker stays unchanged.
6. MUST change the Roundfix Skill and the write-tasks skill only through their
   canonical copies, raise each skill's version in both front-matter fields by
   recording `skills/testdata/owned-skill-versions.json` with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and then regenerate the mirrors with `make skills-sync`, so each mirror
   carries the raised version.
7. MUST NOT change the `### QA settlement` section of any skill, any Go
   source, `.roundfixrc.yml`, `CONTEXT.md` or any ADR.

## Subtasks

- [ ] Write the configuration guide section and Key reference rows.
- [ ] Write the Roundfix Skill `runtime` and `spec` reference text and the
      `spec` command reference text.
- [ ] Write the write-tasks skill heading and the model reference note.
- [ ] Sync the mirrors and record the raised versions.

## Acceptance Criteria

- [ ] The configuration guide names both keys, their defaults, the Light Spend
      Log and the skip warning.
- [ ] The `runtime` reference states the escalation and the Run Event phases.
- [ ] The `spec` references document `--stage tasks`.
- [ ] The write-tasks skill states which Tasks run on the light tier.
- [ ] Every mirror equals its canonical file and both versions are recorded.

## Context

- instruction: `docs/adr/0238-a-light-tier-runs-low-complexity-tasks-on-an-open-model.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/references/model-selection.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `tr -s '[:space:]' ' ' <docs/user-guide/configuration.md | grep -qF -- 'openrouter.implement_monthly_ceiling_usd' && tr -s '[:space:]' ' ' <docs/user-guide/configuration.md | grep -qF -- 'openrouter.light_models' && tr -s '[:space:]' ' ' <docs/user-guide/configuration.md | grep -qF -- 'light tier skipped for Task' && tr -s '[:space:]' ' ' <.agents/skills/roundfix/references/runtime.md | grep -qF -- 'A light Task that fails Verification escalates once' && tr -s '[:space:]' ' ' <.agents/skills/roundfix/references/runtime.md | grep -qF -- 'light_tier_escalated' && tr -s '[:space:]' ' ' <.agents/skills/roundfix/references/spec.md | grep -qF -- 'suggested model-tier' && tr -s '[:space:]' ' ' <docs/user-guide/commands/spec.md | grep -qF -- '--stage tasks' && tr -s '[:space:]' ' ' <.agents/skills/write-tasks/SKILL.md | grep -qF -- 'runs on the light tier' && tr -s '[:space:]' ' ' <docs/references/model-selection.md | grep -qF -- 'light tier' && cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp -s .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md && cmp -s .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md && cmp -s .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md` — expected: exit 0; before this Task none of the phrases exists, so the command fails.
- `tr -s '[:space:]' ' ' <skills/write-tasks/SKILL.md | grep -qF -- 'runs on the light tier' && tr -s '[:space:]' ' ' <skills/roundfix/references/runtime.md | grep -qF -- 'light_tier_escalated' && go test -count=1 ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'` — expected: exit 0; before this Task the mirrors lack the new text, so the command fails.

## References

- `_prd.md` → Core Features 1-8; User Stories 3 and 6
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; Surface Transcript 1; Surface Transcript 2; Vocabulary Contract; Build Order 1
- ADR-0238; ADR-0235; ADR-0187; ADR-0189

## Result

Implemented the light-tier documentation across the configuration guide, model
reference, Spec command reference, canonical Roundfix and write-tasks skills,
and regenerated the embedded skill mirrors. The canonical skills now carry
Roundfix `0.1.35` and write-tasks `0.0.9`; both front-matter fields and the
owned-skill version record were updated by the recorder.

Focused checks:

- `make skills-sync` — passed; canonical and embedded skill trees were synchronized.
- `GOCACHE=/private/tmp/roundfix-task-0233-gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — passed.
- `make baseline-digests` — passed; no derived digest changes were needed.
- Fresh phrase and scope probes — passed for both configuration keys, the Light Spend Log and skip warning, tier and event terms, judge stage/output/refusal terms, write-tasks guidance, and model namespace/default; `cmp` passed for every canonical/mirror skill pair.

Acceptance evidence:

- Configuration guide: names `openrouter.light_models` and `openrouter.implement_monthly_ceiling_usd` with defaults, validation/refusal behavior, Project Config warnings, the Light Spend Log path and fields, the three skip reasons, key-variable handling, and provider-bound data flow.
- Runtime reference: records the light rule, `light-tier` profile source, no reasoning effort, one escalation, `light_tier_skipped` reason codes, and `light_tier_escalated` event data plus the one-Run override.
- Spec references: document `--stage tasks`, advisory `suggested model-tier` output, dispatch independence, and the exact unknown-stage message in both references.
- Write-tasks and model references: state which Tasks run on the light tier, the advisory planning judge command, and the `openrouter/deepseek/deepseek-v4.1-flash` default selection namespace.
- Mirror/version contract: canonical and embedded skill files compare equal, and `skills/testdata/owned-skill-versions.json` records the raised versions.

The Daemon-owned Task status remains `in_progress`. The declared Verification
commands were not run in this Agent turn.

## Carry-forward provenance

- Source Run: `run_20261006T004151Z_0e07759cbdbec8e5`
- Source commit: `c8dec6dd788a4a0a0d5a1f90e9105e9b50a4b09a`
