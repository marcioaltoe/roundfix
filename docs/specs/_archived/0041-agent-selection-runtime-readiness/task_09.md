---
task: task_09
spec: 0041-agent-selection-runtime-readiness
status: completed
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

- [x] Update canonical glossary and agent guidance.
- [x] Update configuration, usage, command, and Setup/Doctor documentation.
- [x] Replace bare-agent examples with profile-led or complete-override forms.
- [x] Add release guidance for the override contract correction.
- [x] Update the canonical Roundfix Skill and OpenAI manifest.
- [x] Regenerate and verify the embedded Roundfix Skill.
- [x] Verify links, terminology, command examples, and ownership boundaries.

## Acceptance Criteria

- [x] Supported docs present the same defaults, adapter identity contract,
      capability proof, fallback, and no-mutation behavior as the CLI.
- [x] Every canonical Agent-starting example either omits all selection flags
      or provides runtime, model, and reasoning effort together.
- [x] Users can distinguish valid model identifiers, advisory rankings, and
      proven operational availability without reading implementation code.
- [x] Setup and Doctor guidance name the official adapter recovery path and do
      not recommend model-managed reasoning for rejected explicit `high`.
- [x] Release guidance identifies partial overrides as rejected usage and gives
      both supported invocation forms.
- [x] Canonical and embedded Roundfix Skill copies are synchronized and shipped
      skill contract tests enforce the new wording.
- [x] No externally managed skill receives an authorial edit.
- [x] Spec 0036 remains ordered after profile-aware Doctor readiness and every
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

## Result

Published one Agent Selection readiness contract across the glossary, command
help, user guides, release guidance, agent guidance, and the shipped Roundfix
Skill surfaces.

Acceptance evidence:

- `CONTEXT.md`, the configuration and usage guides, and command help now use
  the accepted Adapter Readiness, Exact Agent Selection Proof, and Agent
  Selection Profile Readiness vocabulary. They document Sol/high and
  GPT-5.5/xhigh generated defaults, official model validity, advisory ranking,
  exact environment proof, fallback, and fail-before-mutation behavior.
- Documentation-contract tests scan canonical Agent-starting examples and
  reject partial selections. Supported examples now either use the configured
  profile with no selection flags or pass `--agent`, `--model`, and
  `--reasoning-effort` together.
- Setup, `profiles configure`, `profiles validate`, and Doctor guidance now
  identifies `@agentclientprotocol/codex-acp` 1.1.4 or newer as the official
  Codex adapter recovery path, diagnoses the legacy override, requires
  migration authorization, and preserves rejected explicit `high` rather than
  recommending model-managed reasoning.
- The release runbook records partial overrides as intentionally rejected CLI
  usage, includes both supported invocation forms, and states the exit-2,
  pre-mutation behavior.
- The canonical Roundfix Skill and OpenAI manifest were regenerated into their
  embedded copies. Shipped-skill requirements now enforce the adapter,
  profile-readiness, and all-or-none wording. The related repo-owned
  setup-context snapshot was refreshed; no externally managed skill was
  authorially edited.
- ADR-0055, the Spec 0041 validation record, the profile-preflight finding, and
  Spec 0036 all resolve. The guidance keeps Spec 0036 ordered after
  profile-aware Doctor readiness.

Verification evidence:

- `rtk go test ./internal/cli -run 'Test(CommandUsage|DocumentationContract|RoundfixSkill)' -count=1` passed.
- `rtk go test ./skills -count=1` passed.
- `rtk make skills-sync-check` passed.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed all 14 shipped
  skill contracts.
- `rtk make setup-context-check` passed 79 tests and both asset-catalog checks.
- `rtk git diff --check` passed.
- The four durable link targets were resolved from the repository.
- Daemon verification attempt 1 exposed a Doctor help compatibility regression:
  the established `codex runtime hygiene` phrase had been split by a newline in
  the raw help string. The source help text was reflowed without weakening the
  existing assertion.
- `rtk go test ./internal/cli -run 'TestRunCommandHelp/doctor$' -count=1`
  reproduced the failure before the repair and passed afterward.
- `rtk go test ./...` passed after the repair, covering the repository test
  target that failed inside the first `rtk make verify` attempt.
- The Daemon owns the single authoritative full Verification rerun after this
  repair turn, as required by the task execution contract.

Follow-ups: none for Task 09. Spec 0036 retains its declared later ordering.
