---
task: task_05
spec: 0078-roundfix-asks-for-the-review
status: completed
type: chore
complexity: low
---

# Task 05: Synchronise the Roundfix Skill

## Overview

This Spec adds a Review Source mutation, two configuration keys, and a
Preflight refusal, so the Skill must teach them before the Spec can close. This
is the authorized tooling Task.

## Requirements

1. MUST document that `watch` and `resolve` publish one review request per Round
   after the Final Push when `review_source.request_review` is enabled, and that
   `fetch` never publishes one.
2. MUST document both configuration keys with their defaults, and state that
   asking is never Evidence — a published request does not advance a Round.
3. MUST document the Preflight refusal and both refused combinations, so a
   reader meeting the exit `2` knows which pair produced it.
4. MUST state that no automatic retry, backoff, or capacity wait exists, and
   that a refused request ends the Run.
5. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`.
6. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

7. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`; any
   other path is out of scope — stop rather than widen it.
8. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the request, the two keys, the refusal, and the no-retry rule.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [ ] The Skill states that a Round's Final Push is followed by one review
      request, and that `fetch` publishes none.
- [ ] The Skill documents `review_source.request_review` and
      `review_source.request_command` with their defaults.
- [ ] The Skill states that a published request is not Evidence.
- [ ] The Skill names both refused configuration combinations.
- [ ] The Skill states that no automatic retry exists.
- [ ] `skills/roundfix/` is byte-identical to `.agents/skills/roundfix/`.
- [ ] `make verify` exits 0 after the regeneration chain.
- [ ] No Go source file changed.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --quiet HEAD -- '*.go'`
  — expected: exit 0; no Go source changed.
- `git diff --quiet HEAD -- . ':(exclude).agents/skills/roundfix/**' ':(exclude)skills/roundfix/**' ':(exclude)docs/specs/_archived/0078-roundfix-asks-for-the-review/task_05.md' ':(exclude)internal/baseline/assets/setups/**' ':(exclude)internal/baseline/assets/source-baselines/**' ':(exclude)internal/baseline/assets/formatter-fixtures/**' ':(exclude)internal/baseline/assets/profiles/**' ':(exclude)internal/baseline/testdata/**'`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 1, 3, 4
  and 5.
- `_techspec.md` → Integration Points; Build Order 5.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.

## Result

Prepared the canonical Roundfix Skill and its generated mirror for the shipped
Review Source request contract. The Skill now states the post-Final-Push
`watch` and `resolve` behavior, the `fetch` exemption, both configuration
defaults, the Evidence boundary, both Preflight refusal pairs, and the
no-retry rule. No product behavior or Go source changed.

Focused evidence by acceptance criterion:

- Final Push and `fetch`: the canonical and mirrored Skill state that enabled
  `watch` and `resolve` publish one idempotent review request after each Round's
  Final Push and that `fetch` publishes none; the focused `rtk rg` inspection
  found the contract in the canonical Skill.
- Configuration defaults: the same inspection found
  `review_source.request_review` default `false` and
  `review_source.request_command` default `@coderabbitai review`.
- Evidence boundary: the Skill states that publishing a request is not Review
  Source Evidence and does not advance the Round; `watch` and `resolve` still
  wait for head-bound Evidence.
- Preflight refusal pairs: the Skill names both
  `pushTriggersReview=false` with `request_review=false` and
  `pushTriggersReview=true` with `request_review=true`, states exit `2`, and
  explains which Project Config change repairs each pair.
- Refusal and retry policy: the Skill states that an explicit Review Source
  refusal ends the Run and that Roundfix does not retry, back off, or wait for
  review capacity.
- Mirror parity: `rtk make skills-sync` exited `0`, and
  `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  exited `0`.
- Derived artifacts: `rtk make baseline-digests` exited `0` with
  `ok=true, changed=true` and rewrote only paths under
  `DERIVED_DIGEST_PATHS`. The explicit baseline-plan characterization
  re-record then passed 7 tests, and the catalog-diagnostic re-record passed 2
  tests. The first baseline-plan attempt could not access the sandboxed global
  Go build cache; its approved rerun passed.
- Repository Verification: `make verify` was not run because the Daemon owns
  this Task's declared `## Verification` sequence.
- Go-source scope: the focused `rtk git diff --name-only HEAD` inspection
  listed no `.go` path. The changed paths are the canonical Skill, its mirror,
  this Task file, and ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`.
  The Task file's `pending` to `in_progress` frontmatter change was present at
  preflight and remains Daemon-owned; this Agent did not edit Task status.
