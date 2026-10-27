---
type: fix
status: promoted
created: 2026-10-27
spec: 0329-the-skills-lock-matches-what-is-installed-and-a-gate-less-run-settles-green
reason: null
---

# A Run without a QA gate ends Clean with the repository red

## Symptom

A Spec whose `_tasks.md` declares `qa: declined` ran one Task to `Clean`
although the repository Verification (`defaults.verification: make verify`)
was red: `make fmt-check` reported a file the Task created. The Task's own
Verification passed; the repository gate did not hold the settlement, or did
not run.

## Where

Settlement Checks, including the repository Verification at settlement
(`verification.repository_at_settlement`), apply only to Task Graphs with an
authored QA gate Task (Roundfix skill, "Settlement Checks"). A graph with a
declined QA gate has no point where the repository Verification runs before
`Clean`.

## Expected

The repository Verification runs at settlement whether or not the graph has a
QA gate, or the Run does not report `Clean` without it. Decide against the ADR
that scoped Settlement Checks to graphs with a QA gate.

## Evidence

Quaestor (ai-memory message `01a17b1a-38e7-7c82-b268-971c34cfcf4d`,
2026-10-26 23:04 BRT), friction 3:

```
3. O Run run_20261027T011645Z_5cdf7d60ca489eca (0023, qa declined, uma Task) chegou a Clean com o repositório vermelho: `make fmt-check` acusava tests/test_release_workflow.py, criado pela Task, sem ruff format. A Verification da Task passou, mas o portão do repositório (defaults.verification: make verify) não segurou o settle ou não rodou. Corrigimos com um commit de estilo no PR. Esperado: o settle rodar o portão do repositório também em grafo com QA declinado, ou o Clean não ser declarado.
```

Repeated on 0.125.0. Quaestor (ai-memory message
`01a17d0a-cd60-7191-a724-2988affea2d0`, 2026-10-27 08:07 BRT), item 2 under
"Os dois atritos do backlog se repetiram na 0.125.0":

```
2. O Run run_20261027T105442Z_3932740a0b7ec124 (0026, qa: declined) chegou a Clean com o make verify vermelho. A Task 01 mudou o /health, e tests/test_registro.py (Spec 0024) esperava 200 sem sentinela. O make verify na branch pegou o erro antes do PR.
```

Its "Outro ponto da 0024" paragraph belongs to Spec 0330 and is out of scope
here.
