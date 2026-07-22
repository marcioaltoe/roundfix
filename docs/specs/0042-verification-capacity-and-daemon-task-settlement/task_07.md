---
task: task_07
spec: 0042-verification-capacity-and-daemon-task-settlement
status: pending
type: docs
complexity: medium
---

# Task 07: Align operator docs and shipped Agent Skills

## Overview

Publish one operator and Agent contract for independent capacities,
Daemon-owned Task status, Waiting for Verification, and the explicit exit-75
protocol. Configuration, command guidance, autonomous-work instructions, the
repository-owned authorial workflow, and canonical/embedded Roundfix Skills
must match shipped behavior and preserve skill ownership boundaries.

## Requirements

1. MUST document `verification.concurrency`, default `1`, strict positive
   validation, precedence, and its independence from `worktree.concurrency`.
2. MUST document the recommended Task Capacity `2` / Verification Capacity `1`
   flow without claiming machine-wide or external-process coordination.
3. MUST document Daemon-only Implement Task status, implementation-ready Agent
   handoff, focused checks, authoritative full Verification, one Agent repair,
   and capacity release around repair.
4. MUST document exit code `75` as a project-authored Temporary Verification
   Failure, its one exclusive retry, retained evidence, exhaustion behavior,
   and the prohibition on log heuristics.
5. MUST document Run Event Stream and Live Run View working, waiting,
   verifying, retry, and capacity evidence with profile-led command examples.
6. MUST update the canonical Roundfix Skill and repository-owned
   `implement-task` skill, regenerate the embedded Roundfix Skill, and leave
   every upstream-managed skill untouched.
7. MUST preserve ADR, finding, PRD, Tech Spec, glossary, and skill-governance
   traceability and keep all repository content in English.
8. MUST run the complete repository verification gate; any format, test, skill
   synchronization, or build failure is blocking.

## Subtasks

- [ ] Update configuration, command, usage, and autonomous-work guidance.
- [ ] Align implementation-ready handoff and Daemon Verification instructions.
- [ ] Document exit-75 project ownership, exclusive retry, and exhaustion.
- [ ] Document event and Live Run View capacity/phase evidence.
- [ ] Update canonical Roundfix and repository-owned workflow skills.
- [ ] Regenerate embedded skills and verify ownership boundaries.
- [ ] Resolve links, terminology, examples, and the full repository gate.

## Acceptance Criteria

- [ ] Users can configure both capacities and predict default, precedence,
      validation, and per-Run scope from one canonical page.
- [ ] No supported Agent instruction tells an Implement Agent to run declared
      Task Verification, edit Task status, or settle its own terminal verdict.
- [ ] Full-gate requirements remain explicit: the Agent hands off after focused
      checks and the Daemon runs the mandatory Task Verification before
      settlement.
- [ ] Exit `75` documentation clearly assigns classification to the project
      wrapper and promises only one observable exclusive retry.
- [ ] Events/Attach examples use canonical phase names and keep requested
      command output separate from diagnostics.
- [ ] Canonical and embedded Roundfix Skills are synchronized and their tests
      enforce the new configuration, status, and retry wording.
- [ ] No upstream-managed skill changes and every ADR/finding/Spec link resolves.
- [ ] `make verify` passes completely after all documentation and generated
      skill changes.

## Context

- instruction: `docs/agents/skill-governance.md`
- instruction: `docs/agents/autonomous-work.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `CONTEXT.md`
- interface: `README.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/agents/openai.yaml`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `rtk go test ./internal/cli -run 'Test(CommandUsage|DocumentationContract|RoundfixSkill|ImplementTaskSkill)' -count=1` — expected: user guidance, command examples, manifests, and shipped skill contracts match implemented behavior.
- `rtk make skills-sync-check` — expected: canonical and embedded Roundfix-owned skills have no drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: every shipped Roundfix Skill contract passes.
- `rtk git diff --check` — expected: Spec, ADR, finding, documentation, and generated skill files contain no whitespace errors.
- `rtk make verify` — expected: formatting, Go tests, setup-context checks, skill synchronization, shipped skill validation, and build all pass.

## References

- `_prd.md` → all Goals; User Stories 2–7; Core Features 1 and 3–10; User Experience; Decisions.
- `_techspec.md` → System Architecture; API Contracts; Integration Points; Risks & Considerations; Build Order 7.
- `../../adr/0056-spec-runs-separate-task-and-verification-capacity.md` → canonical capacity and temporary-failure contract.
- `../../adr/0057-daemon-exclusively-owns-implement-task-status.md` → canonical Task-status and Agent-handoff contract.
