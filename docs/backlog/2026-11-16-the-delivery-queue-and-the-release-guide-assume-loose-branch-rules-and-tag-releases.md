---
type: fix
status: open
created: 2026-11-16
spec: null
reason: null
---

# The Delivery Queue assumes loose branch rules

## Symptom

The Quaestor wrote on 2026-11-16 (ai-memory message `01a1e522-a608-7563-8568-289c6138d377`, information only, no reply asked) that its Spec 0115 (ADR-0115, ADR-0116) turns on strict required status checks on `main`, runs `make verify` only on `pull_request`, and publishes the tag and release on the push to `main` when `project.version` changes, with no manual tag. Its first observation about Roundfix, in its words:

1. "Com a exigência de branch atualizada, um PR de entrega que fica atrás da main é recusado com \"base branch policy prohibits the merge\" e, na segunda recusa, vai para delivery-error (.agents/skills/roundfix/references/deliver.md, seção Check Wait). A recuperação volta a reviewing, e revisão, gate e CI rodam de novo. A fila poderia atualizar sozinha uma branch que ficou para trás."

Its second observation, that the release guide assumes tag-triggered releases, was split from this entry on 2026-11-16 and adopted by Spec 0400 (`docs/references/2026-11-16-the-release-guide-assumes-tag-triggered-releases.md`), following the maintainer's decision of 2026-11-16 ~20:55 BRT, "Spec fix depois da 0399 (Recommended)", which keeps this half in the backlog.

## Where

The Delivery Queue's merge step and Check Wait (`references/deliver.md`).

## Expected

A delivery Pull Request that falls behind a base branch requiring up-to-date branches is brought up to date by the queue instead of ending in `delivery-error` and a full re-review.

## Evidence

The Quaestor's message above; not reproduced in Roundfix yet.

## Decisions

- 2026-11-18 (`/roundfix --follow-up`, AskUserQuestion): "Spec fix após a corretiva (Recommended)" — a fix Spec delivered right after the corrective qa-gate Spec; authorization "Mesmas do BE".
