---
type: feat
status: promoted
created: 2026-09-30
spec: 0195-owned-skills-and-a-release-step-that-follow-the-bundle
reason: null
---

# A release does not check skills and guides

## Problem

Two audits on 2026-09-30 found owned skills and Baseline guides that described an older product: command forms the Roundfix skill never named, a loop order the Delivery Queue does not follow, a refusal the audit stopped making. Each statement had been true when written. No step of a release re-reads them, so the drift ships and an Agent in an adopter's repository obeys it.

The maintainer decided the same day: "Essa validação e ajuste de skills e baseline deve acontecer junto com o release programado e ao final dos próximos release".

## Where

`docs/user-guide/release-runbook.md` (the "Cutting a release" steps) and the release clause `clause.core.plan-the-release-first` in `internal/baseline/assets/modules/core.json`.

## Expected

The release runbook carries a mandatory step, run before the release Pull Request, that checks owned skills and Baseline guides against shipped behavior. The step rests on checks that run without a model and without a network. A scan of the 138 Baseline clauses by a judging model on 2026-09-30 reached 17% precision on skill-against-clause conflicts, so a model scan is not the recurring step.

## Evidence

Specs 0192 and 0193 record the drift the audits found. The checks they add, and the two this Spec adds, are the ones the step names.
