# Roundfix — primeiro watch com checkout execution (PR #27, 2026-07-15)

Primeiro exercício real do contrato 0027/0029 em um watch de verdade: PR #27 (o próprio
release 0.3.0), 19 Review Issues do CodeRabbit, dois lançamentos. Contexto e histórico em
`2026-07-14-implement-0027-launch-findings.md` e na retrospectiva de 2026-07-15.

## 1. A flake do `gpt-5.6-sol` se apresenta de formas múltiplas — a classificação 0029 não pegou

- **Symptom / evidence**: run `run_20260715T155443Z_0690b6c0f333fb83` — `doctor` reportou
  `model: ok (gpt-5.6-sol)` às 15:54; o Batch 001 morreu às 15:55 com
  `Agent Batch failed after acpx exited with code 1: agent/protocol error` e o warning
  `Model metadata for gpt-5.6-sol not found` chegando como agent_message_chunk. Os 19 issues
  settlaram `failed` em ~60s.
- **Root cause**: a lista de modelos anunciada pelo codex-acp é montada por sessão (fetch
  dinâmico com fallback estático) — a sessão do doctor viu o modelo, a sessão do batch não.
  Desta vez o stderr do acpx **não** continha a frase `did not advertise that model`, então o
  `ModelNotAdvertisedError` (0029, best-effort por design) não classificou e a razão ficou
  genérica.
- **Action / suggestion**:
  1. Ampliar os padrões reconhecidos pela classificação (ex.: tratar o
     `Model metadata for <model> not found` seguido de exit 1 como o mesmo modo de falha).
  2. Estrutural: o probe de preflight valida numa **sessão descartável**, mas o Batch roda na
     sessão do Run (ADR-0018 — uma sessão por Run). Aplicar o Agent Model na **própria sessão
     do Run** durante o preflight (em vez da descartável) eliminaria a janela probe-ok/batch-fail
     por divergência entre sessões.
  3. Operacional: manter `gpt-5.5 @ xhigh` como cadeia de dogfood até a lista anunciada do
     codex-acp estabilizar; o relaunch com override (`run_20260715T160032Z_94b2d8811d607422`)
     entrou em `ResolvingWithAgent` normalmente.

## 2. Issues `failed` terminaram sem Outcome Comment (`commented: false`)

- **Symptom / evidence**: no run 155443Z, os 19 eventos `daemon.source_resolution` do
  `phase":"run_end` mostram `"action":"failed","commented":false` — nenhuma thread recebeu o
  comentário explicativo que a 0027 (Core Feature 7) promete para desfechos `failed`.
- **Root cause**: não estabelecida. Hipóteses: (a) o run-end pass comenta apenas issues
  `unresolved` e assume que `failed` foi comentado no settlement do Batch — mas quando o Batch
  inteiro morre por infraestrutura, o settlement de origem não publica comentários; (b) falha
  silenciosa de publicação. Investigar `resolveBatchSources`/run-end pass com um teste de
  batch-infra-failure.
- **Action / suggestion**: garantir que TODO desfecho `failed` visível no GitHub carregue o
  comentário com a razão (no caso, a razão de infra do Batch), inclusive quando o Batch morre
  antes do settlement normal; adicionar teste de engine para o caminho batch-failed→run_end.

## 3. Checkout execution: o operador não pode tocar no checkout — nem em arquivos não relacionados

- **Symptom / evidence**: este próprio arquivo foi redigido fora do repositório e movido para
  `docs/findings/` só no boundary do run. Com review Runs commitando direto na branch, um
  arquivo novo criado no meio de um Batch entra no snapshot-diff e pode ser varrido para o
  commit `fix: resolve Roundfix batch NNN` — o mesmo modo de falha que a transparência do
  settle (0028) mitigou, agora na superfície de batch.
- **Root cause**: o commit de Batch estagia o diff do snapshot inteiro, sem distinguir edições
  do Agente de edições do operador feitas durante o Batch (a premissa "árvore limpa no início ⇒
  tudo sujo é do Agente" só vale se o operador não tocar em nada).
- **Action / suggestion**: (a) curto prazo — documentar a regra "não edite o checkout durante um
  review Run" na skill/usage (feito na usage.md de 0.3.0, manter); (b) candidato de produto —
  restringir o staging do Batch aos paths que o Run Event Journal registrou como editados pelo
  Agente, ou ao menos avisar quando o commit inclui paths sem edição de Agente correspondente
  (paralelo direto da transparência do settle da 0028).

## 4. Desfecho do relaunch: `CleanUnverified` estreou corretamente — e a janela de 5m é curta para o CodeRabbit

- **Symptom / evidence**: o relaunch (`run_20260715T160032Z_94b2d8811d607422`) terminou
  `CleanUnverified after 1 Round(s)` — 21 resolved, 4 invalid, 0 failed — com Final Push feito
  e a instrução correta no report. `gh pr checks 27` minutos depois: `CodeRabbit pending —
  Review in progress`, ou seja, o check simplesmente ainda não existia dentro da janela.
- **Root cause**: `watch.check_grace_period` default de 5m é menor que o tempo real de o
  CodeRabbit criar e concluir o check de re-review nesta PR.
- **Action / suggestion**: subir o default (10–15m) ou tratar o estado "check criado mas
  pending" como continuação da espera pelo `review_timeout` (hoje só o MISSING usa a janela de
  graça; assim que o check aparece, o fluxo normal já espera). O desfecho em si é o contrato
  funcionando: nenhum falso Clean.

## 5. Ruído meta: CodeRabbit reage aos Outcome Comments do Roundfix

- **Symptom / evidence**: os Outcome Comments publicados pelo run falho (155443Z) viraram
  threads que o CodeRabbit comentou no round seguinte — issues 020–024 do relaunch eram o
  CodeRabbit analisando comentários administrativos do próprio Roundfix. O agente triou os
  quatro como `invalid` com razões corretas ("Administrative reply; the underlying finding was
  resolved in issue_NNN"), então o loop se autocorrigiu, mas gastou triagem.
- **Action / suggestion**: fazer o fetch reconhecer o marker de idempotência do Roundfix
  (`<!-- roundfix: ... -->`) e pular threads cujo corpo raiz é um comentário do próprio
  Roundfix, em vez de materializá-las como Review Issues.

## 6. O que funcionou — manter

- **`--detach` consertado provado em produção**: o mesmo comando que morria silencioso em
  2026-07-14 destacou normalmente nos dois lançamentos (handshake em duas fases da 0029
  cobrindo um probe de ~60s no primeiro run).
- **`doctor` com `adapter:`/`model:`** diagnosticou a seleção em segundos — a linha `model: ok`
  às 15:54 vs. batch falho às 15:55 é exatamente a evidência que faltava para o item 1.
- **Branch Integrity Preflight** passou silencioso com zero pendências (estado limpo real), e a
  atribuição pós-falha ("Uncommitted changes in the checkout are Agent work from this Run")
  apareceu no console como especificado.
- **Recuperação sem perda**: batch morto por infra deixou checkout limpo, artefatos dos 19
  issues em disco e threads abertas — o relaunch re-tenta tudo.
