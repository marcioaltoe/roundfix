---
task: task_03
spec: 0247-a-run-gate-that-runs-the-repository-contracts
status: completed
type: docs
complexity: low
---

# Task 03: The glossary and the repository rules describe Contract Relevance

## Overview

Writes the two terms ADR-0252 introduces into `CONTEXT.md`. It also amends
the repository's hard rule on contracts, so that agents know the Run gate
already runs the relevant ones and `make verify-docs` still runs all of them.
This Task answers the CI failures of 2026-10-07 recorded in entries 237 and
238 of the operator's queue log, from the documentation side. It is
verifiable on its own through phrase checks and the docs contracts.

## Requirements

1. MUST add to `CONTEXT.md`, following the `domain-modeling` skill, the term
   **Repository Contract Test** in the glossary's `**Term**:` / definition /
   `_Avoid_:` shape, beside **Verification**. It is a Go test under the
   `docscontract` or `repocontract` build tag that checks the repository
   itself (its documents, derived artifacts or test wiring) rather than one
   package's behavior. `make verify-docs` runs every one, and
   `make verify-changed` runs those that its Contract Relevance selects
   (ADR-0252). `_Avoid_:` docs test, integration test, CI-only test.
2. MUST add **Contract Relevance** in the same shape, after it. It is the
   declaration in a Repository Contract Test's source that decides when
   `make verify-changed` runs it. The four cases:
   - `always`, on every change;
   - by default, when a file in its package directory changes;
   - `relevant`, when its package directory or a declared path changes;
   - `boundary`, never in that gate.

   A malformed declaration fails the gate. A change to `go.mod`, `go.sum` or
   `Makefile`, and a change list that Git cannot produce, each select every
   contract that is not `boundary` (ADR-0252). `_Avoid_:` test tags, impact
   map, test filter.
3. MUST amend the bullet "HARD RULE — repository contracts validate at the
   pull request boundary" in `docs/agents/specific-repository.md` to add three
   things:
   - `make verify-changed`, the Verification of Runs and of the QA gate, runs
     the Repository Contract Tests that their Contract Relevance selects;
   - `make verify-docs` still runs all of them and still MUST pass before any
     pull request opens;
   - platform-only failures stay a known limit, because CI is the Linux gate.

   The guide cites ADR-0252 and never a Spec.
4. MUST NOT change another glossary entry, any setup-owned guide under
   `docs/agents/` or any other bullet of `docs/agents/specific-repository.md`.

## Subtasks

- [ ] Add the two glossary terms.
- [ ] Amend the hard rule on repository contracts.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Repository Contract Test** and
      **Contract Relevance**.
- [ ] The repository rules say that `make verify-changed` runs the relevant
      contracts and that `make verify-docs` still runs all of them.

## Context

- interface: `CONTEXT.md`
- interface: `docs/agents/specific-repository.md`
- instruction: `docs/adr/0252-the-selective-gate-runs-the-repository-contracts-a-change-makes-relevant.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Verification

- `f() { tr -s '[:space:]' ' ' < "$1" | grep -qF -- "$2" || { printf 'missing in %s: %s\n' "$1" "$2" >&2; exit 1; }; }; f CONTEXT.md '**Repository Contract Test**:'; f CONTEXT.md '**Contract Relevance**:'; f CONTEXT.md 'ADR-0252'; f docs/agents/specific-repository.md 'make verify-changed'; f docs/agents/specific-repository.md 'Contract Relevance'; f docs/agents/specific-repository.md 'ADR-0252'; f docs/agents/specific-repository.md 'Linux gate'; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither file carries the new phrases, so the command fails at its first check. After it, the terms and the amended rule are present, and the docs contracts pass.

## References

- `_prd.md` → Core Feature 3; Glossary
- `_techspec.md` → Build Order 3; Glossary
- ADR-0252

## Result

Implementation:

- Added **Repository Contract Test** and **Contract Relevance** to the glossary
  immediately after **Verification**, including their selection rules, fallback
  behavior, and ADR-0252 reference.
- Amended only the repository-contract hard rule to describe the selective Run
  and QA gate, the all-contracts `make verify-docs` boundary, the Linux CI
  limit, and ADR-0252.

Focused checks:

- `git diff --check` — passed.
- Reviewed the post-edit diff to confirm the glossary terms use the required
  `**Term**:` / definition / `_Avoid_:` shape and that no other glossary entry
  or repository-rule bullet changed.

Acceptance evidence:

- `CONTEXT.md` contains both required definitions and their four Contract
  Relevance cases.
- `docs/agents/specific-repository.md` names `make verify-changed`,
  `Contract Relevance`, `make verify-docs`, and the Linux gate, with ADR-0252
  cited.

The Task's declared Verification commands were left for the Daemon.

## Carry-forward provenance

- Source Run: `run_20261007T202708Z_1511c9319dd8dbf6`
- Source commit: `81da1768f95f4e0f6e2e06d2519f2ba6678e83fb`
