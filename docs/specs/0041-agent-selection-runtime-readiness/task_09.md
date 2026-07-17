---
task: task_09
spec: 0041-agent-selection-runtime-readiness
status: pending
type: docs
complexity: medium
---

# Task 09: Synchronize Agent Selection readiness guidance

## Overview

Publish the implemented adapter, profile, override, Setup, configuration, and
Doctor contracts through canonical vocabulary, user documentation, release
guidance, and the Roundfix Skill. The canonical and embedded skill copies must
remain synchronized and externally managed skills must remain untouched.

## Requirements

1. MUST update `CONTEXT.md` with the accepted adapter-readiness and exact Agent
   Selection Profile proof vocabulary without inventing parallel terms.
2. MUST document Sol/high and GPT-5.5/xhigh generated defaults, official model
   validity, advisory recommendations, and environment-specific exact proof.
3. MUST document supported adapter provisioning, legacy override diagnosis,
   authorization before migration, and fail-before-mutation behavior.
4. MUST document profile-led command invocation and complete all-or-none
   one-Run overrides across command help, user guides, examples, and manifests.
5. MUST document Setup, `profiles configure`, `profiles validate`, and Doctor
   readiness behavior with deterministic next actions.
6. MUST call out the partial-override rejection as an intentional CLI contract
   correction in release guidance.
7. MUST update the canonical Roundfix Skill, regenerate its embedded copy, and
   preserve the repo-owned versus external skill boundary.
8. MUST keep ADR, finding, Spec 0036, and validation-record traceability intact
   and keep all repository content in English.

## Subtasks

- [ ] Update canonical glossary and agent guidance.
- [ ] Update configuration, usage, command, and Setup/Doctor documentation.
- [ ] Replace bare-agent examples with profile-led or complete-override forms.
- [ ] Add release guidance for the override contract correction.
- [ ] Update the canonical Roundfix Skill and OpenAI manifest.
- [ ] Regenerate and verify the embedded Roundfix Skill.
- [ ] Verify links, terminology, command examples, and ownership boundaries.

## Acceptance Criteria

- [ ] Supported docs present the same defaults, adapter identity contract,
      capability proof, fallback, and no-mutation behavior as the CLI.
- [ ] Every canonical Agent-starting example either omits all selection flags
      or provides runtime, model, and reasoning effort together.
- [ ] Users can distinguish valid model identifiers, advisory rankings, and
      proven operational availability without reading implementation code.
- [ ] Setup and Doctor guidance name the official adapter recovery path and do
      not recommend model-managed reasoning for rejected explicit `high`.
- [ ] Release guidance identifies partial overrides as rejected usage and gives
      both supported invocation forms.
- [ ] Canonical and embedded Roundfix Skill copies are synchronized and shipped
      skill contract tests enforce the new wording.
- [ ] No externally managed skill receives an authorial edit.
- [ ] Spec 0036 remains ordered after profile-aware Doctor readiness and every
      durable evidence link resolves.

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `CONTEXT.md`
- interface: `README.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/agents/openai.yaml`

## Verification

- `rtk go test ./internal/cli -run 'Test(CommandUsage|DocumentationContract|RoundfixSkill)' -count=1` — expected: command help, documentation, manifest, and skill contracts match shipped behavior.
- `rtk make skills-sync-check` — expected: canonical and embedded Roundfix-owned skills have no drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: every shipped Roundfix Skill contract passes.
- `rtk git diff --check` — expected: documentation and generated skill files contain no whitespace errors.
- `rtk make verify` — expected: formatting, tests, setup-context checks, skill synchronization, shipped skill validation, and build pass.

## References

- `_prd.md` → User Stories 1–9; Core Features 1–10; User Experience; Decisions.
- `_techspec.md` → Documentation and Roundfix Skill surfaces in System
  Architecture; Build Order 9; Risks and Considerations; Decisions.
- `../../adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md`
  → durable adapter/capability decision.
- `references/validation.md` → Secondbrain, primary-source, and live validation
  evidence.

