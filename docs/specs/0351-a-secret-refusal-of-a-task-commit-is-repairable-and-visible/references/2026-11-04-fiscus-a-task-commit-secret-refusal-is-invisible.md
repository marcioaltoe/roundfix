# Message from fiscus (ai-memory 01a1a775-04f2-7980-8817-539d9fc49c5f, 2026-11-04 13:47 BRT, nature: measured)

Subject: Sugestão: recusa do scanner de segredos no commit da Task não aparece em runs causes nem no park

Fiscus, Roundfix 0.130.1, 04/11/2026. Na Run run_20261104T163352Z_9eb4e9a0c3b72f10 (Spec 0077, task_01), a Verification da Task, o make verify e os dois settlement checks passaram. Mesmo assim a Task ficou failed e o item parou em `parked run-unresolved` com "finding: resolve the blocker". A causa só estava no run_events do roundfix.db: o scanner de segredos do Daemon recusou o commit pelo padrão credential-assignment, porque um teste novo trazia um valor de exemplo de BETTER_AUTH_SECRET com mais de 20 caracteres alfanuméricos.

Problemas:
1. `roundfix runs causes` não lista essa Run, porque só olha falhas de Verification.
2. A mensagem do park não diz que foi o scanner, nem aponta arquivo e linha.
3. As tasks dependentes foram puladas sem diagnóstico visível.

Sugestão: classificar a recusa do scanner como causa própria em `runs causes` (ex.: commit-secret-scan, com arquivo, linha e padrão), mostrá-la no park e mandá-la como Verification Feedback ao Agent no reparo. Hoje o Agent nem fica sabendo; um reparo de uma linha teria evitado a parada. Também ajudaria se o scanner considerasse que o mesmo valor já existe em arquivos rastreados (fixture antiga idêntica).

Evidence (sender): fiscus, Roundfix 0.130.1, run_20261104T163352Z_9eb4e9a0c3b72f10 (spec 0077, task_01); causa lida em run_events do roundfix.db
