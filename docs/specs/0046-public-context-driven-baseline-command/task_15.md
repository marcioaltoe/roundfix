---
task: task_15
spec: 0046-public-context-driven-baseline-command
status: completed
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

- [x] Update the Context-Driven user guide and CLI reference.
- [x] Publish automation schema, decision input, and cross-clone examples.
- [x] Publish migration, recovery, security, and troubleshooting guidance.
- [x] Rewrite the setup skill as thin public CLI guidance.
- [x] Add parser-backed documentation and skill-governance contract tests.

## Acceptance Criteria

- [x] A new user can complete every documented human and automation flow from public docs alone.
- [x] Every command example parses against the shipped CLI.
- [x] Every Decision Document example passes the strict runtime parser.
- [x] Recovery guidance covers stale plans, unsafe carriers, interrupted transactions, and incomplete rollback.
- [x] The setup skill invokes only public Go CLI behavior and describes no Python fallback.
- [x] Canonical and distributed skill guidance is byte-identical.
- [x] Documentation distinguishes profile expectations, executable commands, and recommendations.

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

## Result

### Behavior delivered

- Published the complete human and automation Baseline operating contract in
  the Context-Driven guide, command reference, and README, including adoption,
  update, profile changes, preservation, plan revision, recovery, migration,
  schemas, output ownership, exit categories, security limits, and explicit
  non-execution boundaries.
- Replaced the setup skill's independent script workflow with thin recipes for
  the public Go CLI and added matching Baseline guidance to the Roundfix skill.
- Added executable Go documentation contracts for shipped CLI parsing, strict
  Decision Document parsing, public help and schema vocabulary, portable
  examples, and canonical/distributed skill identity.
- Regenerated the owned skill parity corpus and synchronized embedded Baseline
  setup snapshots and compatibility fixtures after the two owned skill
  documents changed.

### Acceptance evidence

1. The public guide now presents first adoption, update, greenfield,
   preservation, profile revision, automation, restoration, asset
   synchronization, recovery, and migration as complete public-command flows.
2. `TestBaselineExamplesParse` dispatches every fenced Baseline shell example
   through the shipped CLI parsers and passed.
3. `TestBaselineDecisionExamples` extracts the published Decision Document,
   accepts it with the strict runtime parser, and verifies unknown-field
   rejection; it passed.
4. The recovery table has explicit stale-plan, unsafe-carrier,
   interrupted-transaction, incomplete-rollback, and cross-clone actions, and
   `TestBaselineDocumentationContract` passed.
5. The setup skill names the Roundfix binary as its only runtime authority,
   contains no Python invocation or independent fallback, and
   `TestBaselineSkillContract` passed.
6. `rtk make skills-sync-check` passed after regenerating the distributed skill
   bundle; the contract tests also compare both owned skill trees byte for
   byte.
7. The public guide and thin setup skill separately define profile
   expectations, locally bound executable repository commands, and
   recommendations that Baseline reports but never executes.

### Verification evidence

- `rtk env GOCACHE=/private/tmp/roundfix-task15-go-cache go test -count=1 ./internal/cli ./skills -run 'TestBaselineDocumentationContract|TestBaselineExamplesParse|TestBaselineDecisionExamples|TestBaselineSkillContract'`
  passed for both packages.
- `rtk make skills-sync-check` passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task15-go-cache make verify` passed:
  2,007 Go tests, 256 canonical setup tests, 256 distributed setup tests,
  setup asset validation, shipped skill checks, and the CLI build.
- `rtk git diff --check` passed.
