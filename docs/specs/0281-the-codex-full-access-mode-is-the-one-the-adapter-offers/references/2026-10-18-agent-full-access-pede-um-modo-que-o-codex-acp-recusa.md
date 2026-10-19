---
origin: conexus
destination: roundfix
type-hint: finding
created_at: 2026-10-18
capture: manual
---

# `agent_full_access` pede o modo `full-access`, que o codex-acp 2.2.2 recusa

## Observação

No Conexus, Spec 0002, o QA gate (Roundfix 0.87.0, codex-acp 2.2.2, acpx 0.19.6) bloqueou três vezes seguidas todas as linhas de navegador com "browser socket directory /Users/marcio/.agent-browser is not writable: Operation not permitted". Dentro do sandbox do Codex, o daemon do agent-browser não sobe: o Chrome é barrado pelo Seatbelt, e um daemon iniciado fora do sandbox não é alcançado pelo socket. A saída prevista seria `defaults.agent_full_access: true`, mas ele continua inutilizável, como o comentário do `.roundfixrc.yml` registrou para o codex-acp 1.1.9.

## Evidência

- `internal/agent/agent.go` (origin/main `9e92ac39`): `spec.SupportedFullAccessMode = "full-access"` para `codex`.
- Teste direto no adaptador (`acpx codex set-mode -s <session> <mode>`): `full-access`, `danger-full-access`, `full-auto`, `yolo`, `auto` e `default` são recusados com ACP -32602; **`agent-full-access` e `read-only` são aceitos**.
- As sessões ACP ignoram `sandbox_mode` e `default_permissions` do `~/.codex/config.toml`. Mudar o config para `danger-full-access` não alterou nada: o relatório de QA da Run seguinte saiu byte a byte igual (`4c8983b1aeae`) e consumiu 10,3 milhões de tokens.
- O ramo do codex em `acpx_runner.go` só aplica o preset de sandbox `danger-full-access` quando `mode == "full-access"`, ou seja, depois de um set-mode que o adaptador atual recusa.

## Por que importa

Nenhum QA gate com navegador fecha sem override enquanto o modo estiver errado. Na Spec 0001, uma das execuções conseguiu o navegador, aparentemente por escalonamento de permissão aprovado pelo `--approve-all`, o que torna o resultado errático. Cada tentativa custa de 5 a 10 milhões de tokens.

## Ação sugerida

Mapear o codex para `agent-full-access`, ou descobrir o modo pelos `availableModes` da sessão em vez de fixar o nome, e manter a aplicação do preset `sandbox_mode=danger-full-access` nesse caminho. Atualizar o comentário do `.roundfixrc.yml` gerado pelo baseline, que ainda cita o codex-acp 1.1.9.

## Addendum — 2026-10-18 16:15 BRT — o chrome-devtools MCP contorna o bloqueio

Na Run `run_20261018T184649Z_aa55dbcad141babf`, a primeira depois de o QA Task da Spec 0002 passar a exigir o chrome-devtools MCP (commit `6df0e6e`), o gate dirigiu o navegador por ele. O servidor MCP roda fora do sandbox do Codex. O gate executou as linhas de navegador, inclusive com emulação de rede Fast 3G, e achou um defeito real de troca de Organização. Então o modo `full-access` errado bloqueia só o agent-browser, não qualquer navegador. Ainda vale corrigir o mapeamento, e a skill `qa-gate` poderia preferir o chrome-devtools MCP quando o agent-browser não sobe no sandbox.
