---
task: task_04
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
status: completed
type: docs
complexity: low
---

# Task 04: The glossary and the repository rules describe the Full Contract Run

## Overview

Writes ADR-0253 into the durable documents. `CONTEXT.md` gains **Full
Contract Run** and revises **Repository Contract Test** and **Contract
Relevance**. The hard rule in `docs/agents/specific-repository.md` stops
claiming that `make verify-docs` runs every Repository Contract Test, which
was never true of `TestRegenerationIsDeclared`, and names the run that does.
It is verifiable on its own through phrase checks and the markdown contracts.

## Requirements

1. MUST, through the `domain-modeling` skill, revise the definition of
   **Repository Contract Test** in `CONTEXT.md` to this text, keeping its
   `_Avoid_` line:

   > A Go test under the `docscontract` or `repocontract` build tag that
   > checks the repository itself — its documents, derived artifacts or test
   > wiring — rather than one package's behavior. The Full Contract Run runs
   > every one; `make verify-changed` runs the ones its Contract Relevance
   > selects, and `make verify-docs` runs the `internal/docscontract` tests
   > and the contracts its `repo-test` step lists (ADR-0252, ADR-0253).

2. MUST revise **Contract Relevance** in `CONTEXT.md` so that its definition
   keeps every existing sentence and ends with these two sentences before the
   ADR reference, which becomes `(ADR-0252, ADR-0253)`:

   > The selector's summary line names each `relevant` contract it leaves out
   > and each `boundary` contract. The Full Contract Run ignores Contract
   > Relevance.

3. MUST add **Full Contract Run** to `CONTEXT.md`, directly after
   **Contract Relevance**, with this definition and `_Avoid_` line:

   > `make verify-contracts`: one run of every Repository Contract Test that
   > `cmd/verify-select` discovers, whatever its Contract Relevance, with no
   > hand-kept list. CI runs it on every push to main and before every
   > release, so a contract the selective gate leaves out of a pull request
   > still runs before the change ships (ADR-0253).
   > _Avoid_: full suite, nightly run, repo-test

4. MUST, in the hard rule "repository contracts validate at the pull request
   boundary" of `docs/agents/specific-repository.md`, replace the sentence
   that says `make verify-docs` still runs all Repository Contract Tests with:
   "`make verify-docs` runs the `internal/docscontract` tests, the contracts
   `repo-test` lists and `roundfix spec check`, and it **MUST** pass before
   any pull request opens. The Full Contract Run, `make verify-contracts`,
   runs every Repository Contract Test whatever its Contract Relevance, and CI
   runs it on every push to main and before every release (ADR-0253)." Every
   other sentence of the rule stays.
5. MUST NOT reference any Spec or Finding from `CONTEXT.md` or the agent
   guide; both cite ADRs only.

## Subtasks

- [ ] Revise the two existing glossary entries.
- [ ] Add the Full Contract Run entry.
- [ ] Correct the hard rule in the repository rules.
- [ ] Run the markdown contracts.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Full Contract Run** and no longer says `make verify-docs` runs every contract.
- [ ] The repository rules name `make verify-contracts` and keep the pull request boundary rule.

## Context

- interface: `CONTEXT.md`
- interface: `docs/agents/specific-repository.md`
- instruction: `docs/adr/0253-every-repository-contract-test-runs-on-main-and-before-a-release.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`

## Verification

- `t="$(tr -s '[:space:]' ' ' < CONTEXT.md)"; for term in '**Full Contract Run**:' '**Repository Contract Test**:' '**Contract Relevance**:'; do printf '%s\n' "$t" | grep -qF -- "$term" || { printf 'missing term in CONTEXT.md: %s\n' "$term" >&2; exit 1; }; done; for p in 'Full Contract Run[*][*]: .make verify-contracts.: one run of every Repository Contract Test' 'The Full Contract Run runs every one' 'names each .relevant. contract it leaves out and each .boundary. contract' 'The Full Contract Run ignores Contract Relevance' 'CI runs it on every push to main and before every release'; do printf '%s\n' "$t" | grep -q -- "$p" || { printf 'missing phrase in CONTEXT.md: %s\n' "$p" >&2; exit 1; }; done; ! printf '%s\n' "$t" | grep -q -- 'make verify-docs. runs every one' || { printf 'CONTEXT.md still says verify-docs runs every contract\n' >&2; exit 1; }; r="$(tr -s '[:space:]' ' ' < docs/agents/specific-repository.md)"; printf '%s\n' "$r" | grep -q -- 'The Full Contract Run, .make verify-contracts., runs every Repository Contract Test whatever its Contract Relevance' || { printf 'missing the Full Contract Run in the repository rules\n' >&2; exit 1; }; ! printf '%s\n' "$r" | grep -q -- 'still runs all Repository Contract Tests' || { printf 'the repository rules still claim verify-docs runs all contracts\n' >&2; exit 1; }; make verify-docs` — expected: exit 0. Before this Task `CONTEXT.md` has no **Full Contract Run** entry, so the first phrase check fails. After it, every glossary phrase is present, the false claims are gone from both files, and the markdown contracts pass.

## References

- `_prd.md` → Core Feature 4; Glossary
- `_techspec.md` → Glossary; Build Order 4
- ADR-0253; ADR-0252

## Result

- Revised `CONTEXT.md` so **Repository Contract Test** describes the Full
  Contract Run and the narrower `verify-docs` scope, preserved all existing
  **Contract Relevance** sentences, added the selector exclusions and Full
  Contract Run exception, and added the **Full Contract Run** definition.
- Revised the repository hard rule so `make verify-docs` names its actual
  contract scope and `roundfix spec check`, while `make verify-contracts` is
  identified as the discovery-driven run on pushes to main and before releases.
- Focused checks after the edit: inspected the diff and searched the changed
  documents for the required glossary headings, ADR-only references, the
  Full Contract Run phrases, and the removal of the obsolete `verify-docs`
  claim. The Task's declared `## Verification` command, including
  `make verify-docs`, was left for the Daemon as required.
- Acceptance criterion: `CONTEXT.md` defines **Full Contract Run** and no
  longer claims that `make verify-docs` runs every contract — implemented;
  focused phrase checks recorded above.
- Acceptance criterion: the repository rules name `make verify-contracts` and
  retain the pull request boundary rule — implemented; focused phrase checks
  recorded above.
