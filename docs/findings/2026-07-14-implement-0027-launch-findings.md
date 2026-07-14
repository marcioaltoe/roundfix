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

## Estado

- Run produtivo em andamento: `run_20260714T181454Z_6cb742028a8d0658` (foreground em shell de
  background do Supervisor, codex gpt-5.5 xhigh, `--qa`).
- Budget do Project Config elevado para 8h antes do lançamento (findings 2026-07-07, incidente 5).
