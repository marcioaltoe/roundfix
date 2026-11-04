# ai-memory message 01a1a32d-4113-7043-9940-998aca690524 (from gesttione-solutions-pantheon, 2026-11-03T20:50:14Z, nature: measured)

Subject: Defeito 0.129.1: deliver retry reimporta QA Report 'fail' de precondição e não roda o QA gate de novo

Pantheon, roundfix 0.129.1, 03/11/2026.

O que aconteceu:
1. Run 1 da Spec 0059 (`run_20261103T200657Z_a437bac8bed64b17`): as Tasks 01 a 04 resolveram. O QA gate (task_05) recusou na precondição `make verify`, porque o `terraform validate` deu timeout no registry.terraform.io (falha de rede transitória). O Daemon gravou e commitou (73d704a) um QA Report `verdict: fail`, com `rows_blocked_precondition: 1` e `precondition_reason` apontando o log do Run 1. A fila parou com `run-unresolved`.
2. `roundfix deliver retry 0059-...` levou adiante task_01 a task_04 e abriu o Run 2 (`run_20261103T204416Z_73e61e8cd94f990d`). Sem nenhuma sessão de agente, o Run 2 recommitou (a2c3e68) o mesmo report `fail` do Run 1 e rodou só a Verification da task_05. Ela falhou em 0,25 s no `oxfmt --check`, porque o report gerado pelo Daemon não estava formatado (tabelas `| - | --- |`). A task_05 fechou `failed`, e a fila parou de novo com `run-unresolved`.

Esperado: um retry depois de uma recusa de precondição deveria rodar o QA gate de novo (precondição, depois agente), e não reimportar o report da recusa. Além disso, o report de precondição gerado pelo Daemon deveria passar pelo `verification.format` antes do commit, como os outros QA Reports.

Contorno usado aqui: apaguei o Run Branch e o worktree do Run 2 e voltei o Run Branch do Run 1 para o commit da task_04 (de4de1d), tirando o commit do report `fail`. Depois rodei `deliver retry` de novo. Esse retry não imprimiu `Carried forward from Run ...`, então talvez refaça as Tasks 01 a 04.

A referência `implement.md` não documenta a importação de QA Report em retry (só cita "imported pass"). Se o comportamento for intencional, falta documentar como forçar um QA novo.

Evidence: pantheon, Spec 0059-backend-do-fiscus-no-mini; Runs run_20261103T200657Z_a437bac8bed64b17 (commit 73d704a, QA fail) e run_20261103T204416Z_73e61e8cd94f990d (commit a2c3e68, mesma fail reimportada; events cursor 3-13 em 20:44:17Z, sem sessão de agente); verification/batch-001-attempt-1.log do Run 2 com oxfmt --check falhando em qa/qa-report-2026-11-03.md
