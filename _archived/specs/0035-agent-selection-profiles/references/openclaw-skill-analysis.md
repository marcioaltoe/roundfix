# OpenClaw skill analysis

Reviewed: 2026-07-16  
Scope: transferable workflow and model-selection improvements for Roundfix

## Sources reviewed

- acpx repository and built-in agent documentation: https://github.com/openclaw/acpx
- acpx Codex agent: https://github.com/openclaw/acpx/blob/main/agents/Codex.md
- acpx Claude agent: https://github.com/openclaw/acpx/blob/main/agents/Claude.md
- autoreview: https://github.com/openclaw/agent-skills/blob/main/skills/autoreview/SKILL.md
- behavior-validator: https://github.com/openclaw/agent-skills/blob/main/skills/behavior-validator/SKILL.md
- crabbox: https://github.com/openclaw/agent-skills/blob/main/skills/crabbox/SKILL.md
- agent-transcript: https://github.com/openclaw/agent-skills/blob/main/skills/agent-transcript/SKILL.md
- handoff: https://github.com/openclaw/agent-skills/blob/main/skills/handoff/SKILL.md
- session-viewer: https://github.com/openclaw/agent-skills/blob/main/skills/session-viewer/SKILL.md

This analysis paraphrases workflow contracts. It does not copy OpenClaw names, prompts, package identifiers, or generated artifacts into Roundfix production behavior.

## Findings adopted by this Spec

### 1. Review and QA need separate Agent Work Categories

The review workflow is source-aware: it examines changed code, findings, and convergence. The behavior-validation workflow is intentionally source-blind and tests a prewritten behavior contract through user-visible surfaces. Treating both as one model route would hide materially different work. This Spec therefore requires separate `review` and `qa` Agent Selection Profiles and separate selection evidence.

The model profile chooses the Agent Session only. It does not make review proof equivalent to runtime behavior proof and does not replace Roundfix's existing review and qa-gate governance.

### 2. Review has a useful official-id precedent

The reviewed autoreview skill uses `gpt-5.6-sol` with `high` as its default and names `gpt-5.6-terra` as an access-related alternative. It also uses the official Claude id `claude-fable-5`. That supports the official identifiers and review default in this Spec, but Roundfix retains its own explicit profile, preflight, persistence, and fallback boundary.

### 3. Fallback reasons must be classified and visible

The reviewed workflow distinguishes model-access failure from capacity or rate-limit failure rather than treating every agent error as equivalent. Roundfix should likewise normalize the reason that prevented a live selection from starting, record the failed and next tuples, and notify before activation. General prompt, tool, verification, and post-start session failures remain Work Item failures rather than fallback triggers.

### 4. Isolation and proof are different concerns

The reviewed workflows separate an isolated review bundle, remote-environment proof, and source-blind behavior validation. Roundfix should not treat a model preflight as evidence that a Task, review, or QA result is correct. Disposable selection proof answers only: “Can this installed runtime apply this exact model and reasoning tuple now?”

## Improvements recorded as follow-up, not expanded into this Spec

| Improvement | Why it matters | Disposition |
| --- | --- | --- |
| Prewritten source-blind QA behavior contract | Prevents implementation knowledge from weakening user-visible validation. | Evaluate as a future qa-gate Spec; this Spec only gives QA a dedicated profile/session. |
| Frozen review scope baseline and finding classification | Prevents review fixes from expanding without an explicit stop/escalate boundary. | Compare with current Roundfix Review Issue contracts in a separate review-workflow audit. |
| Stop after two non-converging review cycles | Bounds churn when a reviewer and fixer cannot converge. | Existing Run Budget and Max Rounds need a separate product decision; do not couple it to models. |
| Isolated, redacted review bundle | Reduces unrelated context and secret exposure. | Evaluate against the existing Spec Context Bundle and Review Issue artifacts. |
| Remote proof with exact provider/id/command/result | Makes heavy, cross-OS, or live E2E evidence auditable. | Candidate improvement for qa-gate and evidence-gate, not Agent Selection Profiles. |
| Sanitized local transcript provenance | Preserves agent evidence while failing closed on credentials and requiring approval before public use. | Candidate Run Event/Agent log privacy Spec; do not expand selection persistence into prompt storage. |
| Portable path-free handoff | Lets a fresh Agent resume without assuming one machine layout. | Existing handoff skill concern; no model-routing dependency. |

## Roundfix-specific guardrails

- The CLI, not a skill, owns profile configuration and recommendation display.
- Official ids are stored and rendered exactly; runtime aliases supplied by a user remain custom values.
- Preflight proves both the Preferred Selection and every Fallback Selection relevant to the Run.
- Fallback activation is automatic only after a notification and only before Agent work begins.
- A proven model does not prove good output. Task Verification, QA, review convergence, and evidence gates remain mandatory.
- Selection persistence contains no prompt, transcript, credential, cookie, token, or runtime-owned configuration.
- Benchmark rank is never a fallback order unless the user explicitly configures the same order in a Fallback Chain.
