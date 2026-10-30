---
type: fix
status: promoted
created: 2026-10-30
spec: 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
reason: null
---

# The follow-up command rejects documented message types, misreads an archived-everything repository and documents the wrong JSON fields

## Symptom

Três atritos novos do `roundfix follow-up` (v0.126.0), vistos hoje:
1. `--messages` rejeita com exit 2 os tipos "other" e "chore" (`invalid type "other"`), e a referência da skill (follow-up.md) não lista os tipos aceitos. Só passou omitindo `type`. Sugestão: aceitar "other" (a referência diz que "the rest are other") ou listar os valores válidos na referência e no erro.
2. A fonte `specs` aparece `unavailable` ("specs.root resolved to .../docs/specs, which does not exist") num repositório saudável. O archive removeu a última pasta de Spec, e o Git não guarda diretório vazio, então todo repositório que arquivou tudo cai nesse estado. Sugestão: tratar a raiz ausente como `checked`, contagem 0, ou o archive deixar um `.gitkeep`.
3. A referência follow-up.md nomeia os campos dos itens como class, reference e requester, mas o JSON emite `priority_name`, `ref` e `requested_by`, além de `priority`, `basis` e `reason`. Um consumidor que segue a referência lê null.

## Where

- `roundfix follow-up --messages <path>`: the message parser in
  `internal/cli/follow_up_inventory.go` (`followUpMessages`) accepts only an
  empty `type` or `feat`, `fix`, `perf` and `refactor`, and its error names
  neither the accepted values nor the reference that should list them.
- The `specs` source of `roundfix follow-up` and `roundfix follow-up --status`
  (`readFollowUpSpecs`), through `config.ResolveSpecsRoot`, which refuses a
  Spec Root that does not exist.
- `.agents/skills/roundfix/references/follow-up.md` (and its mirror under
  `skills/`), which names the item fields with prose words instead of the
  JSON names the command emits.

## Expected

- Every message type the reference documents is accepted, the reference and
  the error list the same accepted set, and `other` is one of them.
- A repository whose last Spec was archived, so Git no longer holds its Spec
  Root, reports the `specs` source `checked` with count `0`.
- The reference names the JSON fields the command emits.

## Evidence

- Quaestor (ai-memory message `01a18c1a-90b6-7ca0-b720-b5c8672092b5` from
  `default/gesttione-solutions-quaestor`, 2026-10-30 06:18 BRT), Roundfix
  v0.126.0. The message also carried proof results of other Specs; only the
  three frictions above belong to this entry.
- Reproduced by the operator on 2026-10-30 in this repository: with types
  `fix` or omitted the command passes, and the JSON items carry
  `priority_name`, `ref`, `basis` and `reason`. Here `docs/specs/` exists
  (it holds a `.gitkeep`), so the second friction needs a fixture without it.
