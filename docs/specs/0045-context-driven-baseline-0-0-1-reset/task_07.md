---
task: task_07
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
type: backend
complexity: high
---

# Task 07: Resolve individual Readoption dispositions

## Overview

Turn the structural inventory into an explicit, reviewable Decision Plan.
Every source entry must receive a valid individual disposition, and any new
repository-specific rule must be proposed as exact bytes before setup may
write it.

## Requirements

1. MUST add `--decision-file <path>` support for the structured decision
   document defined by the TechSpec while retaining scalar `--decision` for
   decisions that remain scalar.
2. MUST require exactly one valid disposition for every incompatible Source
   Baseline Entry and reject missing, duplicate, unknown, or stale entry IDs.
3. MUST support only the typed dispositions and destination constraints
   declared by the TechSpec; no disposition may be inferred from source text.
4. MUST require exact proposed bytes, target path, and digest for content moved
   into a typed documentation destination.
5. MUST use `docs/agents/repository-rules.md` as the default home for
   Repository-Specific Normative Rules.
6. MUST require confirmation before the first Repository-Specific Normative
   Rules file is created, keep the file unmarked, and preserve it on later
   setup runs.
7. MUST expose stable rejection reasons and a deterministic plan preview while
   performing no writes.
8. MUST include the full normalized decision document in the plan digest
   input.

## Subtasks

- [x] Define and parse the strict structured decision-file schema.
- [x] Validate one-to-one coverage of the Source Baseline inventory.
- [x] Implement typed dispositions and destination restrictions.
- [x] Model exact proposed bytes and digests for documentation moves.
- [x] Add first-write confirmation and preservation rules for the repository
      rules file.
- [x] Add stale, incomplete, duplicate, and unsafe decision mutations.
- [x] Synchronize the canonical and distributed setup skill trees.

## Acceptance Criteria

- [x] A complete decision file maps every inventory entry exactly once and
      produces a stable preview without changing repository bytes.
- [x] Missing, duplicate, unknown, stale, or structurally invalid dispositions
      fail with an entry-specific reason.
- [x] Documentation moves expose their exact proposed bytes, target, and digest
      before confirmation and reject unsafe or untyped destinations.
- [x] The first repository-rules write requires explicit confirmation; an
      existing file is treated as repository-owned and preserved.
- [x] Changing any normalized disposition or proposed target changes the plan
      digest.
- [x] No source entry receives an automatic semantic classification.

## Context

- instruction: `docs/adr/0047-setup-decisions-declare-their-effects.md`
- instruction: `docs/adr/0063-repositories-own-their-http-contract.md`
- instruction: `docs/adr/0064-baseline-readoption-uses-byte-exhaustive-structural-inventory.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `.agents/skills/setup-context-driven/tests/test_decision_rendering.py`
- interface: `.agents/skills/setup-context-driven/tests/test_preview.py`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_readoption_decisions.py'` — expected: complete decisions preview deterministically and every incomplete, stale, duplicate, or unsafe mapping fails with the asserted reason.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_preview.py'` — expected: structured decisions affect the plan digest without producing writes.
- `rtk make skills-sync-check` — expected: canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Core Features 5, 7, 8, and 13; User Stories 1, 3, and 5.
- `_techspec.md` → Implementation Design: Interfaces, Data Models, and API
  Contracts; Build Order 7.
- ADR-0047 → declared decision effects and digest inputs.
- ADR-0063 → durable repository-owned HTTP policy.
- ADR-0064 → individual explicit Readoption dispositions.

## Result

Implemented a strict repeatable `--decision-file` contract for audit and apply.
The parser accepts only `setup-context-driven/decisions/0.0.1`, normalizes
ordered scalar or structured decisions, rejects duplicate JSON keys and
conflicting scalar inputs, and binds both normalized content and file digests
into the plan digest.

Baseline Readoption now validates one explicit classification and typed
disposition per current Source Baseline Entry. Managed entries must exist in
the current catalog; existing typed repository documents require a supported
type, safe path, and matching digest; rejections and non-governed evidence
require individual reasons. Repository-Specific Normative Rules expose exact
base64 proposed bytes, target, and digest at
`docs/agents/repository-rules.md`. An absent file produces a deterministic
confirmation-gated create preview, while an existing unmarked file remains
repository-owned and byte-preserved. Confirmed mutation and atomic Readoption
application remain Task 08's slice.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_readoption_decisions.py'` — passed, 9 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_preview.py'` — passed, 8 tests.
- `rtk make skills-sync-check` — passed; canonical and distributed setup skill trees are byte-identical.
- `rtk make verify` — passed after rerunning the unchanged command with host Go build-cache access. The gate reported 1,694 Go tests, 220 canonical setup tests, 220 distributed setup tests, valid setup assets, a passing Roundfix skill check, and a successful CLI build.

Acceptance evidence:

- Complete deterministic preview and zero writes: `test_complete_decisions_preview_exact_repository_rules_without_writes` compares repeated payloads and repository snapshots while asserting the exact proposed bytes and digest.
- Entry-specific rejection reasons: `test_incomplete_duplicate_unknown_stale_and_invalid_entries_fail_specifically` covers missing, duplicate, unknown, stale, and structurally invalid mappings by stable code and Source Baseline Entry id.
- Typed destinations: `test_typed_document_and_repository_rules_destinations_fail_closed` proves supported typed-document evidence is read-only and unsafe, untyped, or stale targets fail closed.
- First-write confirmation and preservation: the complete preview test requires `plan.confirmation.required`; `test_existing_repository_rules_are_preserved_without_confirmation` proves later audit never compares or replaces existing repository-owned bytes.
- Digest binding: `test_normalized_disposition_and_target_change_plan_digest` changes classification and typed target independently; `test_structured_scalar_decisions_are_normalized_and_digest_bound` proves structured scalar decisions also affect the setup plan digest.
- No semantic inference: every disposition record requires an explicit classification, and missing or invalid classifications fail instead of deriving meaning from `sourceBytes`.
