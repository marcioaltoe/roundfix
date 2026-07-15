# Roundfix — achados no lançamento do implement da spec 0027 (2026-07-14)

Registro dos incidentes ao lançar `roundfix implement --spec 0027-review-loop-integrity` nesta
sessão, antes do Run produtivo `run_20260714T181454Z_6cb742028a8d0658`.

## 1. Modelo pinado no Project Config não é anunciado pelo codex-acp

`runtimes.codex.model: gpt-5.6-sol` derruba todo Batch com
`Agent Batch failed after acpx exited with code 1: agent/protocol error`
(run `run_20260714T180906Z_501e63c7c921d85d`: 4 tasks failed em ~1s cada, 7 skipped, Unresolved).

- Causa raiz (reproduzida com `acpx --verbose --format json` fora do Roundfix):
  `Cannot apply --model "gpt-5.6-sol": the ACP agent did not advertise that model.
  Available models: gpt-5.5, gpt-5.4, gpt-5.4-mini, gpt-5.3-codex-spark.`
- O codex CLI (0.144.4) executa `gpt-5.6-sol` normalmente (`codex exec -m gpt-5.6-sol` → ok);
  o gargalo é o `@zed-industries/codex-acp` **0.16.0 (latest)**, que ainda não anuncia a família 5.6.
- **Gap de guardrail**: o preflight de disponibilidade de modelo (ADR-0039/0041) passou e o Run
  iniciou; a rejeição só apareceu no primeiro prompt do Batch. O probe deveria aplicar o modelo na
  sessão descartável do mesmo jeito que o Batch aplica (ou validar contra a lista anunciada), e a
  única pista no console — "Model metadata for `gpt-5.6-sol` not found" — vinha do codex-acp, não
  do Roundfix.
- Contorno: exceção one-Run `--model gpt-5.5 --reasoning-effort xhigh` (cadeia comprovada).
  Repin do config para a família 5.6 só quando o codex-acp anunciá-la.
- **Diagnóstico adicional (2026-07-15):** a lista anunciada é montada por sessão (fetch dinâmico
  com fallback) — em 15/07 ela passou a incluir `gpt-5.6-sol` e o preflight voltou a rejeitar
  modelos inválidos corretamente. Sessões criadas com minutos de diferença podem ver listas
  diferentes, então o probe nunca garante a visão da sessão do Batch; o fix durável é a rejeição
  em tempo de Batch virar falha acionável (modelo + lista anunciada + recovery) em vez de
  `agent/protocol error`. Specado em `docs/specs/0029-launch-and-recovery-fixes/`.

## 2. `--detach` morre no preflight sem stderr

Com os mesmos flags que funcionam em foreground, `roundfix implement ... --detach` saiu com
exit 1 imprimindo apenas as linhas de prune/reap — sem mensagem de erro, sem Run row criado,
sem console.log. O contrato diz que o pai deve retransmitir o stderr do filho verbatim quando o
filho morre antes do handshake; aqui o stderr chegou vazio. Reproduzido 2×. O mesmo comando sem
`--detach` criou o Run e entrou em `ResolvingWithAgent` normalmente.

Sugestão upstream: quando o filho detached morre antes do handshake com stderr vazio, o pai
deveria reportar ao menos o exit code e o caminho de qualquer log parcial; hoje o chamador fica
sem qualquer pista. Investigar por que o preflight do filho detached falha onde o foreground
passa (possível interação com o probe do modelo/reasoning na re-execução self-detach).

**Root cause (2026-07-15, repro mínimo de spawn):** `detachHandshakeTimeout = 10s` é menor que
o Preflight Validation real — o probe de modelo mediu 11.4s nesta máquina. O pai estoura o
timeout, mata o filho no meio do preflight (que ainda não escreveu nada no console temp) e
retransmite um arquivo vazio, silenciosamente. O branch de timeout nem imprime diagnóstico.
Fix specado em `docs/specs/0029-launch-and-recovery-fixes/` (handshake em duas fases).

## 3. Settle resolve o worktree kept errado como superfície

Com dois Runs kept do mesmo Spec (um force-stopped antigo com worktree presente, um recém-Clean
com worktree removido), `roundfix settle` resolveu o worktree **antigo e obsoleto** como
superfície e recusou a task como `pending` (naquele worktree a task nunca tinha rodado), embora o
task file autoritativo no checkout estivesse `failed`. Contorno: remover o worktree/branch
integrados do run antigo (`git worktree remove --force` + `git branch -d`) e reexecutar o settle,
que então caiu na superfície do repositório atual e passou. Sugestão upstream: a resolução de
superfície deveria preferir a superfície cujo task file está `failed` (ou ao menos avisar qual
superfície foi selecionada e por quê).

## 4. Reproduções confirmadas de findings anteriores

- **Settle agrupa o worktree inteiro** (findings 2026-07-14 §5): o settle da task_09 levou junto
  toda a implementação da task_10 sob a mensagem da task_09.
- **Verificação não-hermética custa ciclos** (§6): `go build ./...` sem `-buildvcs=false` falha
  nos worktrees sob `~/.roundfix` (marcador `.git` inválido em `/Users/marcio`); tasks 09/10
  falharam com o trabalho 100% pronto. Os agentes das tasks 01–08 emendaram o comando no task
  file; o da 09/10 recusou corretamente e settlou failed. Regra de autoria reforçada: sempre
  `go build -buildvcs=false ./...` e `grep` (não `rg`) em comandos de Verification.

## Estado

- Run produtivo em andamento: `run_20260714T181454Z_6cb742028a8d0658` (foreground em shell de
  background do Supervisor, codex gpt-5.5 xhigh, `--qa`).
- Budget do Project Config elevado para 8h antes do lançamento (findings 2026-07-07, incidente 5).
