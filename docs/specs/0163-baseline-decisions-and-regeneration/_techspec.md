---
spec: 0163-baseline-decisions-and-regeneration
status: active
created: 2026-09-24
surfaces: [backend, cli, docs]
---

# Baseline decisions and regeneration

## Executive Summary

Keep the untouched HTTP fields on a mode change. Refuse Greenfield before an
unreachable classification. Declare `make skills-sync` outputs in the tree.
Add a lock reconciliation that removes only entries proven absent at an
immutable commit. Then update the two shipped skills.

## Project Constraints

- Identifier strategy: not applicable — existing decision, lock entry and
  ownership record identities are kept. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  authentication policy is added; the HTTP Contract is edited as repository
  data, and lock reconciliation reuses the unauthenticated public Git fetch of
  `baseline skills restore`. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0081, ADR-0091, ADR-0093,
  ADR-0096, ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0149, ADR-0155 and
  ADR-0156 hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization granted
  2026-09-24, recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/cli/baseline_human_test.go`, `internal/baseline/plan_test.go`,
  `internal/baseline/derived_ownership_test.go`,
  `internal/speccheck/mechanical_test.go`, `skills/_ownership.yml`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`. Sanctioned regeneration:
  `make baseline-digests` and `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## The HTTP mode change

In the `http-contract` branch of `promptBaselineDecision`
(`internal/cli/baseline_human.go`), when a valid current value exists and the
maintainer chooses "Change", the new mode is written into a clone of that
value. `exceptions` and `source` are kept, and the result is validated through
`baseline.ValidateDecisionValue` (which reaches `normalizeHTTPContract`). The
review line names each kept exception scope. With no current value, the prompt
behaves as it does today. A value supplied with `--decision` or
`--decision-file` still replaces the whole decision.

## The Greenfield refusal

In `planRootPreservationWithCatalog` (`internal/baseline/preservation.go`), when
Greenfield meets a non-empty `staleManagedSources`, the plan becomes
`PreservationStateBlocked`. It carries the finding
`baseline.preservation.greenfield.managed-source-retained` and a `NextAction`
naming `preservation.mode=preservation`. No decision skeleton is emitted.
`promptBaselineClassification` plans Greenfield too, and returns that action
before any classification or approval prompt. `BuildPlan` surfaces the same
next action in JSON. Greenfield without stale managed source is unchanged.

## Skill regeneration ownership

`skills/_ownership.yml` declares `owner: dedicated` and
`command: make skills-sync`. `OutputsFor` in
`internal/baseline/derived_ownership.go` becomes command-aware: for
`make skills-sync` it reads that record and the `OWNED_SKILLS` list from the
Makefile, which it reads and never writes. It returns every regular file under
`skills/<owned>/`. For every other command, resolution is unchanged. The audit
in `internal/speccheck/mechanical.go` already consumes `OutputsFor`. The mirror
drift check (`make skills-sync-check`) keeps refusing a hand-edited mirror.

## Lock reconciliation

`baseline.ReconcileSkillsLock` lives in the new
`internal/baseline/skills_reconcile.go`. It takes a repository, a built-in
Profile, a source repository, a 40-hex commit, an optional `--source-dir` and an
optional confirmation. It acquires the commit through the fetch and commit
check in `acquireRestoreGroup`, and reads lock entries with the ordered JSON in
`skills_restore_lock.go`. Each entry whose `source` equals the selected
repository gets one disposition:

- `present` when the entry's `skillPath` exists at the commit.
- `moved` when that path is absent but a `SKILL.md` for the same name exists
  elsewhere.
- `obsolete` when neither exists and the Profile does not require the skill.
- `required-removed` when neither exists and the Profile requires it. This
  blocks the whole plan.

Entries from other sources are `unrelated`. Only `obsolete` entries become
planned lock removals. The installed tree of each one is reported as retained.
A failed fetch or commit mismatch writes nothing and classifies nothing. The
preview returns a plan digest. Apply requires `--confirm-plan` with that digest
and runs through the restore transaction with its preimage check.
`roundfix baseline skills reconcile` in the new
`internal/cli/baseline_skills_reconcile.go` exposes this, with the same exit
categories and output shape as `baseline skills restore`. Doctor is not
changed.

## API Contracts

1. A mode-only HTTP change returns the stored value with only `mode` replaced.
2. Greenfield with retained managed source refuses with
   `baseline.preservation.greenfield.managed-source-retained` before any
   classification prompt, in human and JSON planning.
3. `OutputsFor(root, "make skills-sync")` returns exactly the regular files
   under `skills/<owned>/` for each `OWNED_SKILLS` member.
4. `roundfix baseline skills reconcile` removes only `obsolete` entries, after a
   confirmed plan digest, and writes nothing on a blocked plan, a failed fetch
   or a stale preimage.

## Coverage Map

- Goal 1 → The HTTP mode change; API Contract 1.
- Goal 2 → The Greenfield refusal; API Contract 2.
- Goal 3 → Skill regeneration ownership; API Contract 3.
- Goal 4 → Lock reconciliation; API Contract 4.
- Core Feature 1 → The HTTP mode change.
- Core Feature 2 → The Greenfield refusal.
- Core Feature 3 → Skill regeneration ownership.
- Core Feature 4 → Lock reconciliation.
- Core Feature 5 → Testing Approach 5.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 4.
- API Contracts 1-4 → The HTTP mode change, The Greenfield refusal, Skill
  regeneration ownership, Lock reconciliation.

## Integration Points

- **Spec 0121.** The portfolio these Core Features come from; it keeps Core
  Feature 6.
- **Spec 0148.** Its `verification.incremental` Profile declaration is not
  touched here.
- **`internal/suiteguardcontract`.** Its resolver stays Baseline-only (governed,
  outside this grant); it exempts a subset of what the audit accepts.

## Testing Approach

1. **HTTP.** Replay the archived Finding's multi-exception decision with a
   source, change only the mode, and compare every exception and the source
   field by field. An explicit decision value that drops an exception still
   replaces the decision.
2. **Greenfield.** A repository with a stale managed carrier refuses in
   Greenfield before classification, with the finding and next action, in the
   plan result and in the human command. Greenfield without one still plans.
3. **Ownership.** `OutputsFor` for `make skills-sync` equals the file set of the
   owned mirrors in a fixture tree, excludes `.agents/skills`, Go files,
   `testdata`, `recommended.txt` and an unowned skill directory, and leaves
   `make baseline-digests` outputs unchanged. The audit accepts a changed
   mirror under a grant naming `make skills-sync` and refuses it under one that
   does not.
4. **Reconciliation.** With a local bare repository through `--source-dir`:
   present, moved, obsolete, required-removed and unrelated entries; an
   unreachable source; a mutable ref; a stale preimage; installed trees and
   unknown lock fields kept; Doctor leaves an obsolete entry untouched.
5. **Guidance.** Both skills name the changed behaviour, and each mirror equals
   its canonical copy.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The HTTP mode change (depends on: none).
2. The Greenfield refusal (depends on: 1).
3. Skill regeneration ownership (depends on: 2).
4. Lock reconciliation engine (depends on: 3).
5. Lock reconciliation command (depends on: 4).
6. Shipped guidance and sanctioned regeneration (depends on: 5).
7. Terminal QA (depends on: 1, 2, 3, 4, 5, 6).

## Risks & Considerations

- **Removing a lock entry on weak evidence.** Only a path absent at a verified
  immutable commit, with no same-named skill elsewhere, counts; every failure
  mode writes nothing.
- **Over-broad ownership.** A `make skills-sync` output set that reached
  `.agents/skills` or Go files would let a grant cover authorial source; the
  exclusion cases pin it.
