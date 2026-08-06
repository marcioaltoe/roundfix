---
status: pending
created_at: 2026-08-06
updated_at: 2026-08-06
---

# Minting an ADR opens gaps no one can ever close

**Date:** 2026-08-06
**Found by:** `make verify` failing on `TestCheckCorpusGolden` immediately
after three new ADRs landed with the Spec 0080 and 0081 tech specs.

## What happened

The corpus characterization moved by exactly one number:

```
archived  SC-ADR-RELATED  35 -> 70
```

`SC-ADR-RELATED` reports, at gap severity, that an ADR citing an ADR a Spec
lists is not itself listed by that Spec. ADR-0096, ADR-0097, and ADR-0098 each
cite ADRs that many archived Specs list — ADR-0014, ADR-0080, ADR-0008 — so
every one of those archived Specs immediately gained a gap it did not have
before its own work ended.

The count doubled from one commit that added no Spec and changed no archived
artifact.

## Why the gap is unclosable

Archived Specs are immutable by contract: the archive is the record of what was
built and why, and both `write-prd`'s adoption step and `archive-spec` exclude
`docs/specs/_archived/` from rewrites. So the only two ways to close one of
these gaps are forbidden — edit an archived Spec, or unmint the ADR.

The gap is therefore permanent by construction, and its count can only grow.
Every future ADR that cites a well-referenced predecessor will move this number
again, and each time a maintainer must decide whether the movement is the
harmless kind or a real regression hiding inside it. That is the exact
signal-to-noise failure a characterization test exists to prevent.

## What is not wrong here

The detector is not misfiring by its own definition, and the severity is right:
`gap` reports and exits zero, so nothing is blocked. ADR-0093's non-regression
clause is also honored in the letter — no Spec is *blocked*. What breaks is the
characterization's usefulness: a golden that must be bumped on unrelated
changes stops being evidence and becomes a chore.

## The shape of a fix

`SC-ADR-RELATED` is a suggestion to the author of an *active* Spec — "you
listed one of these, consider the other". An archived Spec has no author left
to advise. Scoping the detector to active Specs, the way ADR-0094 already
scopes detectors by artifact presence, would make the archived count structurally
zero and let the number mean something again.

This belongs to Spec 0080's mechanical-stage work, which is already touching
detector placement and already carries the principle that a check must fail
only where someone can act on it.
