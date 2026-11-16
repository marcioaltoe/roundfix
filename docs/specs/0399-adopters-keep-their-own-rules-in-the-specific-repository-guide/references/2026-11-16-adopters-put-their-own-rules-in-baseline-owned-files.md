---
type: fix
status: promoted
created: 2026-11-16
spec: 0399-adopters-keep-their-own-rules-in-the-specific-repository-guide
reason: null
---

# Adopters put their own rules in Baseline-owned files instead of `specific-repository.md`

## Symptom

Repositories that adopted the Context-Driven Baseline change `docs/agents/skill-dispatch.md` and other files the Baseline owns, or add their own rules to the root `AGENTS.md`, instead of putting each variation in `docs/agents/specific-repository.md`. A later `roundfix baseline update` then meets local edits in files it owns, so updating the fleet stops being easy and the repositories drift from the standard rules.

The maintainer reported it on 2026-11-16 as a frequent problem and asked for a Baseline rule delivered before he updates every project: "Houvesse uma ênfase maior nisso, e como usar o specific repository para as variações ou regras específicas de cada um repositório."

## Where

The Baseline guides that adopters receive (`docs/agents/agent-instructions.md`, `docs/agents/skill-dispatch.md`, `docs/agents/specific-repository.md`), the root instruction files the Baseline writes, and `roundfix baseline update` and `roundfix doctor`, which are the places that could say where a local rule belongs.

## Expected

The Baseline says, where an agent reads it before editing, that a repository's own rules and variations go only in `docs/agents/specific-repository.md` and never in a Baseline-owned file or in the root instruction file. It shows how to write such a rule there. An adopter that already broke the rule learns which edits to move.

## Evidence

The maintainer's report of 2026-11-16 in the Supervisor session. No Finding yet; the Spec that adopts this entry inventories the adopters' edits to Baseline-owned files.
