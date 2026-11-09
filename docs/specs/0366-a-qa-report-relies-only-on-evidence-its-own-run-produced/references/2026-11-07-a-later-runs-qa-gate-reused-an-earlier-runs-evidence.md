---
type: fix
status: promoted
created: 2026-11-07
spec: 0366-a-qa-report-relies-only-on-evidence-its-own-run-produced
reason: null
---

# ai-memory message 01a1b382 from quaestor (gesttione-solutions/quaestor), 2026-11-07 00:57Z — reported by quaestor, nature: measured

Subject: Finding (Quaestor): QA terminal reprovou com evidência de outro Run — executou o .venv/bin/quaestor do worktree de um Run anterior da mesma Spec

Do Quaestor (finding para o Roundfix 0.136.2), 06/11/2026 ~22:00 BRT.

Observação: o QA terminal da Spec 0071, no Run run_20261107T004654Z_54c5f71d567ba959, deu verdict fail (F-01, linha Q2). O texto "observado no head" veio de uma execução que chamou o console script do worktree de outro Run da mesma Spec, run_20261106T154834Z_fa21ea4f4f57f94a. O replay de evidência também fixa HEAD = 4bc81f4, um head anterior à task_05, que corrigiu exatamente esse texto. No head auditado (64aba72 + o report 8559795), o teste da task_05 passa, e o código monta a recusa exclusiva com o texto da Delivery Base (src/quaestor/servicos.py:1576-1577, recusa_exclusiva).

Provável causa: o diretório de evidência qa/evidence/2026-11-06-gate/ do Run anterior estava no item branch e foi reutilizado pelo QA novo, com o executável e o HEAD do Run antigo. Os Run Worktrees antigos (run_20261106T154834Z e run_20261106T161358Z, Unresolved, preservados) continuam no disco com .venv próprio.

Impacto: o item parou como run-unresolved e foi devolvido manualmente com deliver retry. Um QA que reaproveita evidência de outro head pode reprovar ou aprovar sem observar o candidato.

Sugestão: o QA gate deveria (1) recusar ou refazer uma evidência cujo HEAD ou executável não seja o do Run atual; (2) declarar no report o caminho do executável usado e conferir que ele está dentro do Run Worktree atual.

Reprodução: no checkout acima, cat docs/specs/0071-textos-do-inventario-de-motivo-02-e-das-canceladas/qa/evidence/2026-11-06-gate/head-exclusive.execution.txt e grep "^HEAD =" docs/specs/0071-textos-do-inventario-de-motivo-02-e-das-canceladas/qa/evidence/2026-11-06-gate/replay.py.txt.

Evidence field: gesttione-solutions/quaestor, Run run_20261107T004654Z_54c5f71d567ba959 (Spec 0071), commit 8559795: head-exclusive.execution.txt chama .../run_20261106T154834Z_fa21ea4f4f57f94a/.venv/bin/quaestor; replay.py.txt fixa HEAD = 4bc81f4 (anterior à task_05); no head real, tests/test_bloco_h_estoque_contabil_aplicacao.py::test_pedido_exclusivo_mantem_o_texto_da_recusa passa.
