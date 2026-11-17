---
type: fix
status: open
created: 2026-11-16
spec: null
reason: null
---

# The Delivery Queue and the release guide assume loose branch rules and tag-triggered releases

## Symptom

The Quaestor wrote on 2026-11-16 (ai-memory message `01a1e522-a608-7563-8568-289c6138d377`, information only, no reply asked) that its Spec 0115 (ADR-0115, ADR-0116) turns on strict required status checks on `main`, runs `make verify` only on `pull_request`, and publishes the tag and release on the push to `main` when `project.version` changes, with no manual tag. Its two observations about Roundfix, in its words:

1. "Com a exigência de branch atualizada, um PR de entrega que fica atrás da main é recusado com \"base branch policy prohibits the merge\" e, na segunda recusa, vai para delivery-error (.agents/skills/roundfix/references/deliver.md, seção Check Wait). A recuperação volta a reviewing, e revisão, gate e CI rodam de novo. A fila poderia atualizar sozinha uma branch que ficou para trás."
2. "`.agents/skills/roundfix/references/release.md` diz \"Preserve the existing tag-triggered workflow\", o que supõe que todo repositório adotante publica release por tag. O Quaestor passa a publicar no merge da versão (ADR-0116)."

## Where

The Delivery Queue's merge step and Check Wait (`references/deliver.md`), and the release guide (`references/release.md`).

## Expected

A delivery Pull Request that falls behind a base branch requiring up-to-date branches is brought up to date by the queue instead of ending in `delivery-error` and a full re-review. The release guide describes a release published on the merge of the version change as well as one triggered by a tag.

## Evidence

The Quaestor's message above; not reproduced in Roundfix yet.
