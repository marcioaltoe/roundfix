---
status: done
created_at: 2026-09-08
updated_at: 2026-09-08
kind: finding
spec: 0121-baseline-decisions-and-complete-regeneration
---

# Baseline — generated incremental Verification has no declared command (2026-09-08)

The Fiscus capture reports that generated instructions require an incremental
Verification command while its active Profile and Setup Manifest only expose
the complete gate. A passing complete gate does not establish the missing
incremental declaration.

Source: Secondbrain `inbox/roundfix/2026-09-08-profile-exige-verificacao-incremental-que-nao-declara.md`.
The complete original capture remains the provenance record and will move to
that destination's `_triaged/` directory.

## 1. The required local tier is not representable

- Symptom / evidence: the captured public Profile output and generated guide disagree. Local inspection of `internal/baseline/assets/decisions.json` finds only `verification.gate`; `resolveVerificationProjection` in `internal/baseline/profile_alignment.go` projects that decision as `repository-gate`.
- Root cause: the generator publishes a two-tier obligation without a corresponding incremental Decision and Verification Projection.
- Action / suggestion: Spec 0121 must represent and validate both selected commands, preserve the complete gate identity, and provide an explicit migration for existing manifests. Reproduce the captured Fiscus shape in an isolated repository; this Finding grants no Fiscus tooling change.

This is public-output evidence supplied by Fiscus plus local source inspection,
not a claim that a migration or runtime repair has been executed.

## Addendum — 2026-09-08 — Implementation owner

The maintainer selected this remaining intent for the implementation queue.
[0121-baseline-decisions-and-complete-regeneration](../_prd.md) is its primary owner.
The source moves once into that Spec's reference index; other Specs link this
owned copy. Its lifecycle status records adoption, not implementation or QA
completion. Governed mutation and execution still require the consuming grant.

## Capture provenance after triage

The original Inbox path above records the observation at intake. The preserved
capture now lives at `inbox/roundfix/_triaged/2026-09-08-profile-exige-verificacao-incremental-que-nao-declara.md`
in Secondbrain, with `resolved_to` naming this owned Finding.
