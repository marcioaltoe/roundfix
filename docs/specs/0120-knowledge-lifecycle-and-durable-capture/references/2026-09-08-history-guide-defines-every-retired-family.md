---
type: fix
status: promoted
created: 2026-09-08
spec: 0120-knowledge-lifecycle-and-durable-capture
reason: null
---

# The history resolver has a broader contract than the generated guide explains

## Symptom

The guide gives history examples for backlog, findings, and handoffs, but
does not clearly state the rule for other retired families or the location of
inactive ADRs. Fleet sessions have supplied conflicting conventions.

## Where

`docs/agents/docs-layout.md`, canonical
`internal/baseline/assets/modules/context-workflow.json`, and
`internal/spec/archive.go`. Archived Spec
`0094-one-history-root-under-docs` already implemented the family resolver and
retired ADR/review relocation.

## Expected

Document the existing archive-family behavior and the home of inactive ADRs
without inventing another root. Preserve the distinction between a proposed
ADR and a retired decision. Route this remaining contract gap to provisional P1.

## Evidence

Secondbrain `inbox/roundfix/_triaged/2026-08-27-o-guia-nao-diz-o-que-arquiva-nem-onde-e-a-raiz-de-historia-cresceu-sozinha.md` records the original
observation. Triage on 2026-09-08 checked the current local sources named above.
The Inbox body remains the observation's provenance; this entry records intent.

## Addendum — 2026-09-08 — Implementation owner

The maintainer selected this remaining intent for the implementation queue.
[0120-knowledge-lifecycle-and-durable-capture](../_prd.md) is its primary owner.
The source moves once into that Spec's reference index; other Specs link this
owned copy. Its lifecycle status records adoption, not implementation or QA
completion. Governed mutation and execution still require the consuming grant.
