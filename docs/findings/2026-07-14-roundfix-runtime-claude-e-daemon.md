# Roundfix — achados operacionais da sessão 2026-07-13/14 (specs 0009/0010)

Relato de campo do primeiro ciclo completo com dois runtimes (Codex e Claude) em specs reais.
Contexto: 0009 (backend) rodou impecável no Codex; 0010 (redesign) precisou de 5 lançamentos,
2 settles e 1 quarentena até fechar. Cada item abaixo tem reprodução observada nesta sessão.

## 1. Runtime Claude (`claude-code-acp`) — instabilidade grave

- **Adapter ausente não é diagnosticado cedo**: preflight falha com "Failed to spawn agent
  command: claude-code-acp" só depois de sondar modelos; a mensagem não diz que o binário não
  está instalado nem sugere `npm i -g @zed-industries/claude-code-acp`.
- **Nome de modelo divergente do doc**: o adapter sonda nomes curtos (`opus/fable/sonnet/haiku`);
  `--model claude-opus-4-8` (como documentado em `docs/agents/autonomous-work.md`) é rejeitado.
- **`session/set_config_option` não implementado** (ACP -32601): qualquer
  `runtimes.claude.reasoning_effort` não-vazio derruba o preflight. Só funciona com `""`.
  Sugestão: detectar a falta do método e degradar com aviso em vez de falhar.
- **Crash `acpx -32603 Internal error` no fecho de batches longos**: em 3 tasks distintas o agente
  terminou o trabalho, reportou `make verify` verde e o batch morreu no fechamento — task marcada
  `failed` com o trabalho inteiro não commitado no worktree. O `settle` recupera, mas o custo de
  diagnóstico é alto.
- **Escritas via subagentes (tool Task) não persistem no worktree**: a task_01 delegou edições a
  ~20 subagentes e nada chegou ao disco; 7 tool calls falharam silenciosamente na visão do log.
- **Alucinação de contexto**: agentes "leram" um DESIGN.md inexistente (Geist/oklch), arquivos de
  specs arquivadas como se fossem os atuais, e **inventaram escopo** (systems/dashboard e um
  "simulador tributário" com mock engine) — commits sem trailer Roundfix apareceram no worktree.
  Vale investigar o que o adapter fornece de cwd/contexto à sessão.

## 2. Daemon aceita task no-op como `completed`

Quando a Verification da task é genérica (`make verify` já verde no baseline), o daemon marca
`completed` com um Task commit contendo **apenas o flip do task file** (aconteceu 2× — task_01 do
run 1 e task_06 do run 2, ambos claude). Sugestões:

- exigir diff não-vazio fora de `docs/specs/` no Task commit (ou ao menos avisar);
- documentar no contrato de authoring que Verification deve provar o efeito da task
  (greps executáveis salvaram o replay — foi o que travou o no-op na re-execução).

## 3. Vocabulário de status frágil no reload pós-agente

Agente gravou `status: done` (em vez de `completed`) → "reload task file after the Agent" falha e
a task vira `failed` com todo o trabalho não commitado. Sugestões: normalizar sinônimos óbvios
(`done`→`completed`) ou validar/reescrever o frontmatter antes de encerrar o batch.

## 4. Stop Request órfão bloqueia relançamentos

Se o processo do run morre (kill) antes de processar o Stop Request, o registro fica `Active` e
bloqueia qualquer `implement` novo. O erro sugere `roundfix stop <run-id> --force`, mas a forma
posicional com `--force` é rejeitada ("unexpected argument") — só funciona
`roundfix stop --run-id <id> --force`. Corrigir a dica ou aceitar a forma posicional.

## 5. `settle` agrupa o worktree inteiro num commit só

`settle` commita "all Run Worktree changes plus the task file": com várias tasks failed no mesmo
worktree, o primeiro settle carrega o trabalho de todas sob a mensagem de uma. Não houve perda,
mas o histórico fica enganoso. Sugestão: avisar quando o commit do settle contém caminhos fora do
escopo aparente da task, ou permitir settle com pathspec.

## 6. Authoring em macOS — nota para o guia

Gate de Verification com `wc -l | grep -qx 0` falha no BSD `wc` (saída com padding). O Codex
diagnosticou e **não contornou** (comportamento correto), mas custou um ciclo. Vale uma nota no
guia de authoring: preferir formas portáveis (`[ "$(ls dir)" = "x" ]`) a pipelines com `wc`.

## 7. Review de PR (watch/fetch/resolve) — sem worktree e com guardrails de integridade da branch

Proposta de mudança de contrato para os Review Runs, motivada por dois riscos reproduzidos nesta
sessão: (a) worktrees mantidos com trabalho **não integrado** na branch (runs Unresolved/Stopped
deixaram commits e árvores sujas fora da branch da PR); (b) locks de Run ativos/órfãos sobre a
mesma working tree bloqueando ou concorrendo com operações na branch. Um Review Run que parte de
um HEAD que não contém todo o trabalho pendente revisa e "corrige" código defasado — e o Final
Push pode publicar uma HEAD que silenciosamente omite trabalho já concluído em worktrees.

**Contrato proposto:**

1. **Review Runs não criam Run Worktree.** `watch`/`resolve` operam no checkout do usuário, na
   própria branch da PR. Racional: o fix de review é, por definição, um delta sobre a HEAD
   publicada da PR; a indireção por worktree é o que cria o problema de trabalho encalhado e os
   desfechos IntegrationPending. (O isolamento por worktree continua fazendo sentido para
   `implement`, onde há paralelismo de Tasks — não para o loop de review, que é sequencial e
   ancorado na branch.)
2. **Preflight determinístico de integridade da branch** — executado tanto no preflight do
   `watch` quanto no do `fetch` (o fetch é uma etapa do watch e pode ser invocado avulso; os dois
   precisam do mesmo guardrail):
   - **Worktrees pendentes**: enumerar todos os worktrees/branches `roundfix/run-*` cuja base é a
     branch da PR. Se algum tiver commits além da base, exigir a integração antes de iniciar —
     integrar automaticamente quando `git merge --ff-only` resolver, ou falhar o preflight
     nomeando cada worktree e o comando exato de integração. Só iniciar com zero worktrees
     pendentes relacionados à branch.
   - **Runs ativos**: nenhum Run (de qualquer tipo) vinculado à branch/working tree pode estar em
     execução. Falhar nomeando o `run-id` e o comando `roundfix stop <id>` (hoje o lock já
     bloqueia, mas locks órfãos de processos mortos ficam `Active` — ver item 4 — e o guardrail
     precisa diferenciar "ativo de verdade" de "lock órfão", oferecendo o `--force` do stop).
3. **Bypass somente explícito, com trilha de auditoria**: um parâmetro tipo `--force` (ou
   `--skip-branch-integrity`) pula os dois guardrails. Quando usado, o Roundfix **comenta na
   Pull Request** registrando a ação: quem/quando, run-id, quais guardrails foram pulados e o
   estado ignorado (worktrees pendentes e/ou runs ativos enumerados). Sem o comentário publicado,
   o bypass falha — a auditoria é parte do contrato, não cortesia.
4. **Documentação junto**: refletir o novo contrato na skill do Roundfix (seções "User-Facing
   Review Runs" e "Run Worktree Isolation") e nos resources de watch/fetch/resolve — hoje a skill
   afirma que review Runs executam em Run Worktree; se o contrato mudar, a skill e os exemplos de
   comando precisam mudar na mesma entrega, senão os agentes seguem o contrato antigo.

## 8. Propagação de status para o Review Source (GitHub) — mais cedo e com motivo

Hoje a resolução dos threads no Review Source acontece tarde no ciclo (após verificação/push do
Round), e statuses não-resolvidos ficam invisíveis para quem acompanha a PR pelo GitHub — o
revisor humano não sabe se um apontamento foi descartado, falhou ou ainda vai ser tratado.

**Proposta:**

1. **Atualizar o status do thread no GitHub o quanto antes**: ao final do fix de cada Review
   Issue; se a latência das chamadas por issue afetar significativamente a performance do Round,
   degradar para o final do Batch (nunca mais tarde que isso). O objetivo é o estado no GitHub
   acompanhar o progresso real do Run, não só o desfecho.
2. **Comentar no thread todo desfecho que não for `resolved`**:
   - `invalid`/ignorada → comentário breve com o motivo da triagem (por que o apontamento não se
     aplica);
   - `failed` → comentário breve com a causa da falha (o que impediu o fix seguro);
   - `duplicated` → apontar o thread/occurrence canônico.
   Isso fecha o loop com o revisor (humano ou bot) sem exigir acesso aos artefatos locais do
   Roundfix, e deixa a PR auto-auditável.

## 9. Auditoria do fluxo do watch (PR #17, run_20260714T100348Z — Clean em 1 round)

Round observado de ponta a ponta com `watch --source coderabbit --pr 17 --agent codex
--until-clean`. O fluxo esperado (fetch → resolve → commit → push → novo review) foi seguido em
quase tudo: espera pelo Review Source assentar → fetch do Round 001 com 50 Review Issues →
resolve em 3 Batches (triagem + agente + `make verify` como gate de cada batch) → 1 commit por
batch (`fix: resolve Roundfix batch 00N`) → resolução dos threads no GitHub **por batch**
(20+20+10) → integração fast-forward na branch → Final Push. Resultado: 48 resolved, 2 invalid,
0 failed, 0 unresolved. Dois desvios e uma observação:

- **Desvio A — a etapa "novo review" não aconteceu.** O run se declarou Clean imediatamente após
  o Final Push com a nota `Review Source check missing for the pushed HEAD; treating Run as
  Clean`. É comportamento documentado, mas é uma corrida: o CodeRabbit ainda nem havia começado a
  reanalisar a HEAD nova — "check ausente" foi tratado como "check verde", e issues novos sobre os
  commits de fix ficarão sem ninguém observando. Solução: período de graça com polling pelo check
  do Review Source no commit pushado antes de declarar Clean (o mesmo mecanismo de espera que o
  fetch já usa no início do round); alternativamente um modo estrito em que `--until-clean` exige
  sucesso afirmativo do check, e/ou um desfecho distinto (`CleanUnverified`) em vez de nota em
  stderr — hoje o chamador de script não consegue distinguir os dois Cleans pelo exit code.
- **Desvio B — Roundfix commitou os artefatos de review na branch** (`docs: review rounds for pr
  17`, com `docs/specs/_reviews/pr-17/round-001/issue_*.md`), enquanto a skill afirma "Roundfix
  never commits or gitignores review artifacts; repository owners decide". Comportamento e
  documentação precisam convergir: se commitar é o contrato novo (armazenamento na árvore de
  specs, ADR-0029), atualizar a skill e tornar o commit opt-out; senão, é bug.
- **Observação C — threads `invalid` resolvidos em silêncio.** Os 50 threads foram resolvidos no
  GitHub, incluindo os 2 triados como `invalid` — o revisor humano vê o apontamento "resolvido"
  sem saber que foi descartado nem por quê. Reforça o item 8: `invalid` deveria comentar o motivo
  da triagem em vez de resolver silenciosamente (ou nem resolver, deixando o comentário).
- **Positivo:** a resolução de threads por batch já existe (calibra o item 8 — o gap real é a
  granularidade por issue e os comentários nos desfechos ≠ resolved), o verify por batch segurou
  o gate de qualidade, e a integração ff-only + Final Push mantiveram branch local e PR
  perfeitamente sincronizadas.

## 10. Comportamento observado de `fetch` e `resolve` avulsos (PR #17, round 002)

Exercício deliberado das superfícies separadas após o watch, para semear o comportamento real:

- **`fetch` é exemplar**: Run próprio com desfecho `Fetched`, relatório determinístico com
  "No side effects" explícito (sem agente, worktree de execução, commit, push ou resolução de
  threads). Round 002 materializado com 6 issues em `docs/specs/_reviews/pr-17/round-002/`.
- **Composição fetch→resolve funciona**: o `resolve` consumiu os 6 issues **baixados** ("selected
  6 downloaded Unresolved Review Issue(s) from 1 Compatible Artifact Round(s)") sem re-buscar —
  as etapas são de fato separáveis, como o modo watch as encadeia.
- **Política de versionamento dos artefatos é inconsistente por superfície**: o `fetch` deixa o
  round untracked no checkout; o daemon de `resolve`/`watch` **commita** os artefatos no fim
  (`docs: review rounds for pr 17`). A skill afirma que o Roundfix "never commits review
  artifacts" — três comportamentos declarados/observados diferentes para a mesma coisa. Decidir
  um contrato (sugestão: commitar sempre que o destino é a árvore de specs versionada, com
  opt-out) e alinhar skill + superfícies.
- **Mensagem de commit dos artefatos não identifica o round**: dois commits idênticos
  `docs: review rounds for pr 17` (rounds 001 e 002). Incluir o número do round na mensagem
  (`docs: review round 002 for pr 17`) tornaria o histórico auto-explicativo.
- **Experimento de colisão**: deixamos o `round-002/` untracked no checkout enquanto o run
  commitava os mesmos paths no worktree. A integração **absorveu** sem recusar — o checkout
  terminou com os arquivos autoritativos (`status: resolved`), árvore limpa e branch sincronizada.
  O risco de IntegrationPending por artefato untracked não se reproduziu; a motivação do item 7
  permanece sendo o caso de worktrees com **commits** encalhados, não este.
- **Contagem do relatório é cumulativa entre rounds**: o resolve deste run (6 issues) reportou
  `Clean after 1 Round(s): 54 resolved, 2 invalid...` — soma dos rounds 001+002 do artefato, não
  do run. Para scripts que parseiam o relatório, "quantos issues ESTE run resolveu" não é
  derivável; vale separar contagem do run e contagem acumulada da PR.

## O que já foi mitigado neste repo

- `.roundfixrc.yml` com `runtimes.claude.reasoning_effort: ""` comentado com o porquê.
- Guarda-corpos + gates executáveis de efeito nos task files da 0010 (padrão a repetir).
- Memória do agente orquestrador com o procedimento de `settle` e a rota de fallback para Codex
  (`docs/agents/autonomous-work.md`, "when in doubt, use Codex" — confirmado na prática).
