---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# `reconcile --apply` keeps Runs it already classified superseded

## Symptom

A Run whose Tasks a later Clean Run delivered is classified superseded, but only `--discard-superseded` acts on it; `--apply` keeps it, so merged Specs leave branches behind.

## Where

`internal/cli/reconcile.go` (legacy `--apply` classification).

## Expected

Release proven-superseded Runs under `--apply`, or fold them into the post-merge release (B02).

## Evidence

Session 2026-09-25 reconcile dry run: `run_20260924T205403Z_4627184d18e75ee8` (Spec 0159, merged in #247) classified superseded and retained.
