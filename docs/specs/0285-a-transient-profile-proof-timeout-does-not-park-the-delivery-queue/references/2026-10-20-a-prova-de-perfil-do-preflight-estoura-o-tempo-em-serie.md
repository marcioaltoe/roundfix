---
origin: oraculum
destination: roundfix
type-hint: fix
created_at: 2026-10-20
capture: manual
---

# A prova de perfil do preflight estoura o tempo em série e estaciona a Delivery Queue

## Observação

Entre 19/10 19h30 e 20/10 00h27 BRT, a Delivery Queue do oraculum estacionou
cinco vezes por timeout na prova de perfil. A falha caiu em runtimes e modelos
diferentes, e nenhuma criou Run. Toda vez, a mesma combinação passou logo
depois de um `roundfix deliver retry`. A máquina estava com carga baixa: load
average 1,6 às 23h27.

| Hora (BRT) | Fase | Tupla | Park |
| --- | --- | --- | --- |
| 19/10 19h30 | pre-PR review (0084) | claude/sonnet/high, review fallback[1] | `review-blocked` |
| 19/10 21h27 | implement preflight (0086) | claude/opus/high, fallback[1] de backend, docs e qa | `delivery-error` |
| 19/10 22h57 | implement preflight (0086) | codex/gpt-6.1-sol/high, qa preferred | `delivery-error` |
| 19/10 23h27 | implement preflight (0086) | claude/opus/high, qa fallback[1] | `delivery-error` |
| 20/10 00h27 | implement preflight (0086) | codex/gpt-6.1-sol/high, backend e qa preferred | `delivery-error` |

A mensagem é sempre a mesma: `profile proof setup timed out twice: context
deadline exceeded; classification: temporary; next: rerun the command when
load drops`. Em 22h57 veio junto
`disposable Agent Session cleanup found no Session to close ... No named
session "roundfix-preflight-..."`.

## Evidência

- `roundfix deliver status` do oraculum nos horários acima (Delivery Owners
  `delivery-d1255e752e9d-7546`, `-32619`, `-96819`, `-72986` e `-86228`). Os
  logs estão em
  `~/.roundfix/artifacts/d1255e752e9d9f49/delivery/<owner>/console.log` e
  `~/.roundfix/artifacts/d1255e752e9d9f49/pre-pr-review/<id>/pre-pr-review.json`.
- Versões: roundfix 0.91.0 e 0.93.0 (o binário subiu às 22h07), acpx 0.19.6,
  claude-agent-acp 0.93.0 e codex-acp 2.2.2 (`roundfix doctor`).

## Por que é acionável

- **O preflight prova toda tupla configurada**, a preferida e as de reserva.
  Basta uma das cinco ou seis tuplas demorar na criação da sessão para o item
  estacionar, e cada park custa um ciclo do Supervisor, cerca de 30 minutos.
- **A classificação já é `temporary`.** Uma nova tentativa automática dentro
  do preflight, ou um prazo configurável para a criação da sessão, evitaria o
  park. Também seria possível provar só a tupla preferida e deixar as reservas
  para quando forem usadas, como a ADR-0050 já faz com as Fallback Chains.
- **A documentação não oferece alternativa:**
  `.agents/skills/roundfix/references/profiles.md` não tem nenhuma opção para
  esse prazo.
