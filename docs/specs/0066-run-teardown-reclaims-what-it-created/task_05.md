---
task: task_05
spec: 0066-run-teardown-reclaims-what-it-created
status: completed
type: chore
complexity: low
---

# Task 05: Synchronise the Roundfix Skill with the new contract

## Overview

The repository HARD RULE requires a Pull Request that changes CLI behaviour to
ship the Roundfix Skill update with it. This Spec changes what `reconcile`
offers and when Branch Integrity Preflight refuses, so the Skill must teach
both before the Spec can close. This is the authorized tooling Task.

## Requirements

1. MUST document the two new `reconcile` candidate kinds, their proofs, and
   that dry-run remains the default.
2. MUST document the Branch Integrity Preflight change: proven-superseded Run
   Branch work no longer blocks a review Run, and every other refusal stands.
3. MUST state that an unprovable termination is reported, never treated as
   success, so a reader does not mistake silence for a stopped process.
4. MUST regenerate the `skills/roundfix/**` mirror with `make skills-sync`.
5. MUST run `make baseline-digests`, then re-record the two characterization
   corpora that command does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

   Neither test is in `BASELINE_DIGEST_STEPS`, so the sanctioned command reports
   "no changes" while the gate stays red. The flags are verbatim because they do
   not match their test names.
6. MUST change only `.agents/skills/roundfix/**`, `skills/roundfix/**`, this
   Task file, and the ADR-0081 digest fallout under `DERIVED_DIGEST_PATHS`; any
   other path is out of scope — stop rather than widen it.
7. MUST NOT change behaviour. This Task documents what shipped.

## Subtasks

- [ ] Document the reconcile candidates and the Preflight change.
- [ ] Document the unprovable-termination reporting.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both re-records.

## Acceptance Criteria

- [ ] The Skill documents both new candidate kinds and the dry-run default.
- [ ] The Skill documents the Preflight relaxation and states the other
      refusals stand.
- [ ] The Skill states an unprovable termination is never reported as success.
- [ ] `skills/roundfix/` is byte-identical to `.agents/skills/roundfix/`.
- [ ] `make verify` exits 0 after the regeneration chain.
- [ ] No Go source file changed.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.
- `git diff --name-only HEAD | grep -vE "^(\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0066-run-teardown-reclaims-what-it-created/task_05\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Project Constraints (tooling authority); Core Features 2, 4, 5.
- `_techspec.md` → Integration Points; Build Order 5.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md` → standing
  Roundfix Skill CLI synchronisation grant.
- ADR-0081.

## Result

### Implementation

- The canonical Roundfix Skill now documents the additive `process` and
  `runBranch` reconcile candidates. It names the terminal Run and inspected
  owned process tree proof for `process`, and the superseding QA Report plus
  clean registered Run Worktree proof for `runBranch`.
- The Reconcile Command remains dry-run by default. The Skill states that
  ambiguous evidence is preserved, `--apply` is the only mutation switch, and
  process cleanup succeeds only when every reported process is proven absent.
- Branch Integrity Preflight now teaches the shipped failed-cycle behavior:
  the proven QA branch set is reported, left unchanged, and disregarded so a
  review Run can proceed. It states explicitly that this is the only
  relaxation and every other existing refusal still blocks.
- `make skills-sync` regenerated the embedded Skill mirror. The sanctioned
  digest chain then refreshed only ADR-0081-derived Baseline assets and the
  plan-characterization corpus named by the Task.

### Focused checks and regeneration

- `rtk make skills-sync` — exited 0.
- `rtk make baseline-digests` — exited 0 and reported five regenerated
  derived files under `internal/baseline/assets/setups` and
  `internal/baseline/testdata`.
- `rtk go test ./internal/baseline -count=1 -run TestBaselinePlanCharacterization -update-baseline-plan-characterization`
  — exited 0; seven characterization cases passed while re-recording.
- `rtk go test ./internal/baseline -count=1 -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics`
  — exited 0; two characterization cases passed while re-recording.
- `rtk go test ./internal/baseline -count=1 -run 'TestBaselinePlanCharacterization|TestCatalogDiagnosticCharacterization'`
  — exited 0 after regeneration; nine cases passed without update flags.
- `rtk proxy diff -qr .agents/skills/roundfix skills/roundfix` — exited 0;
  the two-file Skill trees are byte-identical.
- A literal-safe `rtk grep -n` inspection found both candidate kinds, the
  unprovable-termination rule, the Preflight relaxation, and the retained
  refusals. An earlier scan was discarded because unquoted Markdown backticks
  were interpreted by the shell; it did not mutate any file.
- The commands under `## Verification`, including `make verify`, were not run;
  the Daemon owns them.

### Acceptance criterion evidence

- Both candidate kinds and dry-run default: the Reconcile section names
  `process` and `runBranch`, states each proof, and says dry-run remains the
  default before showing `--apply`.
- Preflight relaxation and retained refusals: the review-run section states
  that the proven failed-cycle QA set is disregarded without Git mutation, and
  that non-fast-forward or ambiguous work, preserved branch-set evidence,
  another Active Run, a dirty tracked checkout, and every other refusal still
  block.
- Unprovable termination: the Reconcile section states that it is reported
  with its reason and never treated as success; host silence is not proof that
  a process stopped.
- Mirror parity: `rtk proxy diff -qr .agents/skills/roundfix skills/roundfix`
  exited 0. Both `SKILL.md` files also produced the same SHA-1 digest,
  `1c7757142ff79e90ab2de4931ded4e70673b85df`.
- Repository Verification: not observed in this Agent turn because `make
  verify` is a declared Verification command reserved for the Daemon.
- Go-source boundary: the post-regeneration `rtk git diff --name-only HEAD`
  inventory contains only the authorized Skill pair, this Task file, and
  ADR-0081-derived paths under `internal/baseline`; no `.go` path appears.
