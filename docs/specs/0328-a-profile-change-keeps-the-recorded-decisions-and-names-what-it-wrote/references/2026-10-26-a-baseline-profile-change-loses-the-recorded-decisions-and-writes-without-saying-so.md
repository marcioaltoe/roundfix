---
type: fix
status: promoted
created: 2026-10-26
spec: 0328-a-profile-change-keeps-the-recorded-decisions-and-names-what-it-wrote
reason: null
---

# A Baseline profile change loses the recorded decisions and writes without saying so

## Symptom

Quaestor moved its adopted repository to the `python` Baseline Profile with
Roundfix v0.123.0 and met five frictions on the way, reported in part 1 of
the message quoted under Evidence:

- (a) `roundfix baseline update --repo . --decision profile=python` refuses
  with "unknown decision profile"; only `baseline plan --profile python`
  followed by `baseline apply --plan --confirm-plan` changed the profile.
- (b) `baseline plan` did not take the decisions the Setup Manifest records:
  eleven were passed by hand, then `preservation.mode`, one refusal per round.
- (c) `baseline apply` created `AGENTS.<sha256>.md` at the root, a copy of the
  old `AGENTS.md`, without naming it in its output or in
  `docs/agents/setup-context.json`.
- (d) `baseline skills reconcile --revision b56fc5f` answered "no changes"
  while the installed `knowledge-workspace` stayed at the Setup Snapshot's
  older upstream revision.
- (e) `baseline skills restore` installed the `ai-memory-*` skills under
  `.agents/skills` but created no link for them under `.claude/skills`.

## Where

`roundfix baseline update`, `roundfix baseline plan`, `roundfix baseline
apply`, `roundfix baseline skills reconcile` and `roundfix baseline skills
restore`; packages `internal/baseline` and `internal/cli`.

## Expected

A profile change names its working command when asked through the update,
carries every recorded decision the new profile still selects, and asks for
everything else in one refusal; an apply names each root backup it writes;
reconcile says what it compared against; restore keeps the repository's
per-skill `.claude/skills` links.

## Evidence

Provenance: ai-memory message `01a17abd-20fe-7213-ab5c-f1a1bdb043fe` from
`default/gesttione-solutions-quaestor`, created 2026-10-26 21:22 BRT. Part 1
belongs to this Spec; part 2 belongs to Spec 0327 and is out of scope here.
Part 1, verbatim (Portuguese, as received):

```text
1. Perfil python: adotado e verde na v0.123.0. PR gesttione-solutions/quaestor#46 (37c11cf); CI Verify verde na branch e na main; make verify com 341 passed e 4 skipped; `roundfix baseline starters --profile python` mostra Makefile, .roundfixrc.yml e verify.yml como keep, State current.
Atritos:
(a) `roundfix baseline update --repo . --decision profile=python` recusa com "unknown decision profile". A troca só funcionou com `baseline plan --profile python` e `baseline apply --plan --confirm-plan`.
(b) O plan não herdou as decisões do Setup Manifest. Tivemos de passar à mão 11 decisões (autonomous.enabled, domain.layout, language.generated, repository.extension.enabled, secondbrain.enabled, spec.scaffold, triage.external, verification.gate, verification.incremental, runtime.backend, runtime.design) e depois preservation.mode, um erro por rodada.
(c) O apply criou na raiz AGENTS.<sha256>.md, cópia do AGENTS.md antigo, sem citação no setup-context.json e sem aviso.
(d) O preset fixava marcioaltoe/skills 65d5f36 com knowledge-workspace 0.6.1, mas a curadoria estava em b56fc5f (0.6.2). `baseline skills reconcile --revision b56fc5f` respondeu "no changes".
(e) `baseline skills restore` instalou as ai-memory-* em .agents/skills, mas não criou os links em .claude/skills.
A v0.124.0 ainda não foi aplicada aqui.
```

Reproduced by the operator on 2026-10-26 with Roundfix built from
`6a2630e5` in temporary repositories (see the Spec's PRD, Validation).
