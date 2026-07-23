---
task: task_15
spec: 0046-public-context-driven-baseline-command
status: pending
type: docs
complexity: high
---

# Task 15: Document the public Baseline operating contract

## Overview

Publish the user and automation contract after every public command is stable.
The documentation becomes sufficient for first adoption, update, recovery,
and migration without reading implementation artifacts or invoking an
independent setup engine.

## Requirements

1. MUST document first adoption, update, greenfield, instruction preservation,
   profile changes, rejected-plan revision, automation, recovery,
   troubleshooting, and migration from the Python-backed skill.
2. MUST document exact command grammar, plan/result schemas, stdout/stderr
   ownership, exit categories, confirmation, stale-plan, manual fallback, and
   cross-clone safety.
3. MUST document profile authoring, Repository Skill Set restoration, asset
   synchronization, security limits, recommendations, and operations the
   Baseline Command never executes.
4. MUST convert the setup skill guidance to the public command family with no
   Python invocation or independent behavioral fallback.
5. MUST keep CLI help, user guide, automation examples, emitted Decision
   Document, recovery guidance, and skill recipes under executable
   documentation contract tests.
6. MUST use canonical glossary terms and contain no external repository
   identity or environment-specific path.

## Subtasks

- [ ] Update the Context-Driven user guide and CLI reference.
- [ ] Publish automation schema, decision input, and cross-clone examples.
- [ ] Publish migration, recovery, security, and troubleshooting guidance.
- [ ] Rewrite the setup skill as thin public CLI guidance.
- [ ] Add parser-backed documentation and skill-governance contract tests.

## Acceptance Criteria

- [ ] A new user can complete every documented human and automation flow from public docs alone.
- [ ] Every command example parses against the shipped CLI.
- [ ] Every Decision Document example passes the strict runtime parser.
- [ ] Recovery guidance covers stale plans, unsafe carriers, interrupted transactions, and incomplete rollback.
- [ ] The setup skill invokes only public Go CLI behavior and describes no Python fallback.
- [ ] Canonical and distributed skill guidance is byte-identical.
- [ ] Documentation distinguishes profile expectations, executable commands, and recommendations.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/agents/skill-governance.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `README.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`

## Verification

- `rtk go test -count=1 ./internal/cli ./skills -run 'TestBaselineDocumentationContract|TestBaselineExamplesParse|TestBaselineDecisionExamples|TestBaselineSkillContract'` — expected: help, docs, examples, schemas, and skill recipes match shipped behavior.
- `rtk make skills-sync-check` — expected: canonical and distributed setup and Roundfix skill guidance is synchronized.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 6; User Stories 9–10; Core Features 18, 21–22; Delivery note.
- `_techspec.md` → API Contracts; Integration Points: Skill distribution; Build Order 9.
- ADR-0066 → public CLI authority and thin skill.
- ADR-0068 → human and automation command contracts.
