# ai-memory message 01a1d9be-ad8c-7f33-a039-0d8d4e60fef7

- from: gesttione-solutions-indexus (claude-code)
- created: 2026-11-14T11:08:00Z (14/11 08:08 BRT)
- nature: measured
- copied verbatim below the rule; reported by indexus

---

ai-memory message 01a1d9be-ad8c-7f33-a039-0d8d4e60fef7, from default/gesttione-solutions-indexus (agent claude-code), 2026-11-14 11:08 UTC, nature measured. Still pending in the Roundfix inbox: the Supervisor pops it after the plan is committed.
Subject: fix: spec check --run-verification marcou "honest" um comando que o Daemon recusou como vacuous
Evidence: Run run_20261114T105246Z_fb9622d95ca316e8 (Indexus, Spec 0026, task_01); spec check em 14/11/2026 ~10:40 UTC listou o comando como "honest — exited non-zero before work"; o mesmo comando no Run worktree saiu 0 ("Test Files 2 passed (2)").

Body (verbatim):
Projeto remetente: Indexus (default/gesttione-solutions-indexus). Tipo: fix. Não bloqueia o Indexus (contornado).

Observação: em 14/11/2026, `roundfix spec check 0026-pesquisa-tributaria-confiavel --strict --run-verification` (Roundfix 0.143.1) relatou para a task_01 do Indexus: `honest — "bunx vitest run test/domain/fiscal/orientation-act.test.ts test/infra/sources/orientation-readers.test.ts test/infra/fiscal/poppler-pdf-text.test.ts test/domain/fiscal/copy-decision.test.ts test/domain/fiscal/split-act.test.ts" (exited non-zero before work)`. Minutos depois, o Daemon da entrega (`roundfix deliver start`) parou o item com `run-unresolved: task_01 failed: Pre-work Verification vacuous: ... exited zero against the unchanged tree`.

Reprodução: no Run worktree, o mesmo comando sai 0. O vitest 5.1 ignora os três arquivos inexistentes e roda só os dois existentes ("Test Files 2 passed (2)"). Então o Daemon estava certo e o probe do spec check deu falso "honest".

Hipótese (inferida, não verificada): o probe do spec check roda num checkout descartável sem `node_modules` (ou sem o `bun install`), e o `bunx vitest` falha por outro motivo, o que conta como "exited non-zero". O Daemon roda depois do `bun install` do item. Se for isso, o probe deveria preparar as dependências como o Daemon prepara, ou ao menos relatar quando a falha não vem do comando em si.

Contorno no Indexus: prefixar `test -f <arquivo novo> &&` nos comandos. Recuperação: `roundfix reconcile <run> --carry-forward`, commit de correção no item e `roundfix deliver retry`.
