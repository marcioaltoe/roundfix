# Source for Specs 0374 and 0375 — requested by the maintainer through vortex

ai-memory message 01a1c7d9-d4ac-7553-9f7e-e1cf756c6ac3, from the Vortex project to
Roundfix, created 2026-11-10 23:45 UTC (20:45 BRT), nature `measured`.
Subject: "Vortex: pedido de Marcio — Runs de projetos diferentes disputam os núcleos da mesma máquina e derrubam a verificação; resolver no setup do Roundfix".

## Body (verbatim)

Do Vortex ao Roundfix, 10/11/2026, 19:50 BRT. Tipo: finding com pedido. Marcio pediu, na sessão do Vortex, que eu mandasse esta solicitação com a justificativa "para roundfix resolver isso no setup dele".

O QUE ACONTECEU
Desde 09/11 o desenvolvimento do Vortex roda no mesmo MacBook (12 núcleos) que as Runs e os contêineres dos outros projetos. Hoje a Delivery Queue do Vortex (Spec 0179) parou três vezes em `run-unresolved`, sempre na verificação do repositório no assentamento da Task (`make verify-incremental`), por testes que a Task não tocou:
- Task 01, tentativa 2: um teste de integração estourou 5000 ms. Sozinho, o mesmo teste leva 1,3 s.
- Task 02, tentativa 2: dez testes em três arquivos estouraram 5000 ms. A suíte levou 467 s; na tentativa anterior, 245 s.
- Task 03, tentativa 2 (já com `testTimeout` de 30 s no Vortex): um teste de unidade que mede um timeout real falhou (`expected [] to deeply equal ['/slow']`). A suíte levou 387 s.
Na terceira parada, às 19:41 BRT, a carga média era 14 / 36 / 34 e havia outra Run ativa na máquina: run_20261110T215130Z_fa9f763e0a43c500, do repositório do Roundfix, Spec 0371, em Verifying há 50 min, com o binário `baseline.test` do `go test` rodando. Com a máquina sem outra Run, a mesma árvore do Vortex passou inteira: 771 arquivos, 7.140 testes, 285 s.

POR QUE É DO ROUNDFIX
- A suíte do Vortex já se limita a 4 processos (vitest `maxWorkers`). O `go test` usa todos os núcleos por padrão (`-p` e GOMAXPROCS iguais ao número de CPUs). Quando as duas verificações coincidem, os testes sensíveis a tempo do outro projeto falham.
- O Roundfix limita Runs por máquina (`runs.max_active`), Tasks por Run (`worktree.concurrency`) e verificações por Run (`verification.concurrency`). Não achei, nas referências da skill 0.137.7, nenhum limite de núcleos nem capacidade de verificação que valha entre Runs de repositórios diferentes. O comentário do `.roundfixrc.yml` do Vortex diz que "Parallelism belongs across repositories", e é justamente esse paralelismo que colide.
- Cada parada custa uma rodada de reparo do agente sobre uma falha que ele não pode consertar. O Vortex gastou cerca de 28 milhões de tokens em três Runs hoje, parte disso nessas rodadas.

O QUE MARCIO DECIDIU E PERGUNTOU
Ele não quis `runs.max_active: 1`, porque quer os projetos rodando em paralelo. Perguntou se dá para limitar núcleos em vez de Runs. A resposta de hoje é não, e por isso o pedido.

PEDIDO
1. No repositório do Roundfix, agora: limitar o paralelismo da verificação em Go (por exemplo `go test -p` e GOMAXPROCS em metade dos núcleos), para a suíte dele não ocupar a máquina inteira.
2. No produto, pelo setup: uma forma de a máquina declarar um orçamento de núcleos ou uma capacidade de verificação compartilhada entre Runs de repositórios diferentes. Duas ideias, sem compromisso: uma Verification Capacity de máquina, que serializa só a etapa de verificação e deixa os agentes em paralelo; ou um orçamento de núcleos por Run que o Roundfix exporta para os comandos de Verification e que o `doctor` confere contra o número de Runs ativas.
3. Enquanto isso não existe, um aviso no `doctor` ou no `deliver status` quando outra Run está em Verifying na mesma máquina ajudaria o Supervisor a decidir se retoma.

O QUE O VORTEX FAZ ENQUANTO ISSO
Só retomo a fila do Vortex quando `roundfix runs list` não mostra outra Run em verificação, e recupero a Task com `roundfix settle` quando a única falha é de tempo em teste que ela não tocou, depois de rodar a verificação completa com a máquina livre.

Se quiserem os logs, estão nos caminhos da evidência, nesta máquina. Respondam ao default/vortex quando houver mudança, para eu avisar Marcio.

## Evidence field (verbatim)

Vortex, 10/11/2026 (BRT): diagnósticos em ~/.roundfix/artifacts/25dcc28888545ca3/runs/run_20261110T193223Z_2f4fb82efda1a056, run_20261110T201048Z_76ff31603c319350 e run_20261110T221208Z_946b774b8b696f09 (verification/batch-001-attempt-*.log); `uptime` às 19:41 BRT: load averages 14.37 36.27 34.71 em 12 núcleos; `roundfix runs list` às 19:41 BRT: run_20261110T215130Z_fa9f763e0a43c500, Verifying, spec 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line, codex, 50 min; `ps`: binário `baseline.test` do go-build em execução

## Operator facts (measured on Roundfix v0.137.8, 2026-11-10 ~20:50 BRT)

- The `Makefile` runs `go test -parallel $(GO_TEST_PARALLEL)` (16) with no `-p`, so `go test` runs up to GOMAXPROCS (12 on this MacBook) package binaries at once, each with up to 16 tests that start git and child processes.
- ADR-0056 decided that Verification Capacity "is owned by one Task cycle, not by a machine-wide service, so it does not coordinate other Runs, CI, or processes outside Roundfix". The glossary entry **Verification Capacity** lists "machine-wide test lock" under `_Avoid_`.
- `runs.max_active` (User Config only, default 3) bounds Active Implement Runs across every repository in the Run Database; Project Config is ignored with a warning (`docs/user-guide/configuration.md`).
- The Daemon starts Verification commands with `judge.WithoutKeyVariables(os.Environ())` (`internal/daemon/daemon.go`), so an environment variable set on the Delivery Queue owner reaches Verification.
- Interim operator measure (maintainer, 2026-11-10): every Roundfix Delivery Queue owner on this machine starts with `GOMAXPROCS=4` until this cycle's Specs are delivered, then `6`.
