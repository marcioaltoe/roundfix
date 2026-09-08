---
type: fix
status: open
created: 2026-09-08
spec: null
reason: null
---

# Skill regeneration has no declared ownership for the authorization audit

## Symptom

The audit understands sanctioned regeneration fallout, but skill-sync
outputs have no ownership declaration. Grants therefore enumerate mirrored
files individually, contrary to the recorded command-owned derivation policy.

## Where

`skills/` has no `_ownership.yml`; current declarations live under
`internal/baseline/`. `internal/speccheck/mechanical.go` and
`internal/baseline/derived_ownership.go` consume ownership. Spec 0114 added
fallout-aware auditing; ADR-0149 and ADR-0081 explain why grants name the cause.

## Expected

Give the sanctioned skill regeneration command a machine-readable output
contract that the audit can resolve, retaining refusal of manual derived edits.
The concrete ownership declaration belongs to provisional P2, Baseline
decisions and complete regeneration.

## Evidence

Secondbrain `inbox/roundfix/_triaged/2026-08-30-skills-nao-declara-ownership-e-forca-todo-grant-a-enumerar-espelhos.md` records the original
observation. Triage on 2026-09-08 checked the current local sources named above.
The Inbox body remains the observation's provenance; this entry records intent.
