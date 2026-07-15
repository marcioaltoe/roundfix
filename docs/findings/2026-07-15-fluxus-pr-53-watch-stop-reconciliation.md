# Dogfood no Fluxus PR #53 — watch, stop e reconciliação (2026-07-15)

Sessão de uso do Roundfix no PR
`gesttione-solutions/fluxus#53`, branch `ma/fix-validation-errors`. O comando
`roundfix watch --source coderabbit --pr 53 --agent codex --until-clean --max-rounds 3 --detach`
criou o Run `run_20260715T160701Z_19300a52ae4882c9`, com Codex `gpt-5.5`,
`xhigh` e HEAD `a8eccf9c4968ef8ea29e63b89e4baff95f7536d3`.

O PR tinha 306 arquivos analisáveis. CodeRabbit publicou um status de sucesso,
mas pulou a revisão porque o limite é 300 arquivos. O Run buscou um Round sem
Review Issues e terminou `CleanUnverified` após cinco minutos. Durante a
limpeza também foi removido manualmente um Run Worktree antigo, já integrado,
que continuava registrado como `IntegrationPending`.

## 1. Status CodeRabbit presente foi reportado como ausente

- **Sintoma / evidência**: `gh pr checks 53 -R gesttione-solutions/fluxus`
  mostrou `CodeRabbit pass` com `Review skipped: 306 files exceed the limit
  of 300`. O console do Run terminou com
  `CleanUnverified: Merge-Ready was not confirmed because the Review Source
  check never appeared within the grace period`.
- **Causa raiz**: desconhecida. O status existia no mesmo HEAD enquanto o
  Roundfix aguardava. Falta confirmar se o matcher ignorou um Status Context,
  não reconheceu o detalhe de revisão pulada ou observou uma projeção diferente
  da usada pelo GitHub CLI.
- **Ação / sugestão**: persistir no Run Event Journal a identidade e o resumo
  de cada check observado. Adicionar um teste com Status Context CodeRabbit
  bem-sucedido e detalhe `Review skipped`. O desfecho deve distinguir check
  ausente de revisão não executada.

## 2. Revisão pulada parece revisão limpa

- **Sintoma / evidência**: o console registrou
  `Review Source status: settled` e `Fetched Round 001 with 0 Review
  Issue(s)`, embora CodeRabbit não tenha revisado o PR por exceder o limite de
  arquivos.
- **Causa raiz**: o contrato atual reduz o status externo a presença, sucesso e
  issues encontradas. Ele não representa “check verde, mas revisão pulada”.
- **Ação / sugestão**: introduzir uma classificação explícita, por exemplo
  `ReviewSkipped`, sem tratá-la como Merge-Ready. A saída deve incluir o
  motivo e a ação concreta: reduzir o PR, dividir a mudança ou solicitar outra
  revisão. Um Round com zero issues só prova revisão limpa quando o Review
  Source confirma que analisou o HEAD.

## 3. Espera do check não mostra prazo nem progresso

- **Sintoma / evidência**: depois de buscar o Round 001, o Run permaneceu
  `Active` e a última linha continuou `Fetched Round 001 with 0 Review
  Issue(s).`. A configuração usa cinco minutos de grace period, mas a tela
  permaneceu em `WaitingForReview` sem prazo, contagem regressiva ou último
  check observado.
- **Causa raiz**: a espera faz parte do contrato de `CleanUnverified`, mas a
  projeção do Run não expõe essa fase.
- **Ação / sugestão**: representar a fase como `WaitingForReviewCheck` e
  mostrar início, tempo restante, HEAD esperado e último estado observado.
  Publicar Run Events periódicos ou nas mudanças de estado, sem poluir o
  journal com cada poll.

## 4. Stop Request não interrompeu a espera do Review Source

- **Sintoma / evidência**: `roundfix stop
  run_20260715T160701Z_19300a52ae4882c9` respondeu
  `Stop Request recorded; the Run stops after the current Work Item settles`.
  Cinco segundos depois, `runs list` ainda mostrava o Run `Active`, sem Work
  Item em execução e sem novo evento.
- **Causa raiz**: desconhecida. A evidência indica que a espera do grace period
  não observa a Stop Request com a mesma frequência dos limites entre Work
  Items.
- **Ação / sugestão**: tornar todo wait do Review Source cancelável e consultar
  a Stop Request durante o polling. Um watch sem Work Item ativo deve terminar
  `Stopped` prontamente, sem esperar o grace period completo.

## 5. Force stop liberou o lock, mas o processo alterou o desfecho depois

- **Sintoma / evidência**: `roundfix stop --force
  run_20260715T160701Z_19300a52ae4882c9` informou `State: Stopped` e
  `released its Active Run locks`. O processo detached continuou e o registro
  final passou para `CleanUnverified`; `runs list --all --state all` mostra
  esse último estado.
- **Causa raiz**: o force-stop completa o Run e libera o lock, mas o proprietário
  detached ainda consegue executar a conclusão normal e sobrescrever um estado
  terminal. O encerramento da Agent Session não controla esse proprietário
  quando nenhuma sessão existe.
- **Ação / sugestão**: proteger transições terminais com compare-and-set:
  `CompleteRun` não pode trocar um estado terminal por outro. O force-stop
  também deve cancelar o processo proprietário pelo PID registrado e aguardar
  confirmação limitada. Adicionar um teste de corrida em que o proprietário
  tenta concluir após `Stopped`; o estado, outcome event e lock devem
  permanecer `Stopped`.

## 6. Force stop tentou fechar uma Agent Session que nunca existiu

- **Sintoma / evidência**: o Run buscou zero Review Issues e nunca iniciou
  Agent work. Mesmo assim, o force-stop tentou fechar
  `roundfix-run_20260715T160701Z_19300a52ae4882c9` e emitiu
  `No named session ... for cwd /Users/marcio/dev/fluxus and agent codex`.
- **Causa raiz**: o cleanup deriva o nome esperado da Agent Session sem
  confirmar se o Run chegou a criá-la.
- **Ação / sugestão**: persistir a criação da Agent Session e fechar somente
  sessões registradas. “Sessão não encontrada” deve ser idempotência, não
  warning, quando nenhum Batch iniciou.

## 7. Round sem issues deixou o checkout do usuário sujo

- **Sintoma / evidência**: o watch escreveu o arquivo untracked
  `docs/specs/_reviews/pr-53/round-001/round.md`, com
  `issue_count: 0`. O operador precisou removê-lo antes do squash merge para
  recuperar um checkout limpo.
- **Causa raiz**: Review Runs escrevem artefatos no checkout do usuário. Um
  desfecho `CleanUnverified` ou interrompido não consolidou nem removeu o
  Round vazio.
- **Ação / sugestão**: não persistir um Round sem Review Issues quando o Review
  Source não executou a revisão. Alternativamente, manter artefatos
  intermediários no Artifact Directory e promover apenas Rounds relevantes ao
  Spec Root. Se o arquivo permanecer, o relatório terminal deve nomear o path e
  dizer se ele foi commitado ou ficou untracked.

## 8. Estado IntegrationPending não reconcilia integração manual já concluída

- **Sintoma / evidência**: o Run
  `run_20260715T130045Z_efce6cafab78b0b1` continuou
  `IntegrationPending` com Run Worktree e Run Branch, embora
  `roundfix/run-run_20260715T130045Z_efce6cafab78b0b1` já fosse ancestral da
  branch do PR. O operador verificou zero commits exclusivos, removeu o
  worktree e apagou a branch manualmente após o squash merge. O Run Database
  ainda conserva `IntegrationPending`.
- **Causa raiz**: a integração foi concluída com o comando Git impresso pelo
  Roundfix, mas não existe uma etapa que reconcilie o resultado de volta ao Run
  Database e ao lifecycle do Run Worktree.
- **Ação / sugestão**: oferecer `roundfix integrate <run-id>` ou
  `roundfix reconcile <run-id>`. O comando deve verificar ancestralidade,
  atualizar o outcome sem perder o histórico e remover Run Worktree e Run
  Branch quando não houver trabalho exclusivo. O preflight de um Run seguinte
  também pode fazer essa reconciliação de forma segura.

## 9. Skill do projeto consumidor divergiu do CLI usado

- **Sintoma / evidência**: os arquivos têm hashes diferentes. A skill em
  `/Users/marcio/dev/fluxus/.agents/skills/roundfix/SKILL.md` afirma que
  `resolve` e `watch` executam em Run Worktree. O `roundfix watch --help`
  usado na sessão afirma que watch executa no checkout do usuário e não cria
  Run Worktree. A skill canônica deste repositório já descreve o contrato
  atual.
- **Causa raiz**: a cópia repo-local do consumidor ficou atrás da versão
  instalada/canônica, sem falha antecipada no início do comando.
- **Ação / sugestão**: adicionar uma versão ou hash de contrato à skill e fazer
  `roundfix skills check` comparar o pacote canônico com cópias repo-locais.
  Comandos operacionais podem emitir um warning acionável quando detectarem
  drift: `roundfix skills install --target project`.

## O que funcionou — manter

- O relatório detached identificou Run ID, console log, comando de attach e
  comando de stop.
- O cabeçalho registrou Review Source, Agent Model, Default Reasoning Effort,
  HEAD e `Max Rounds: 3`.
- O Run Database e o console log preservaram evidência suficiente para
  reconstruir a corrida entre stop, force-stop e conclusão.
- `roundfix doctor` validou Node, acpx, adapter, Agent Model e higiene do
  binário antes do watch.
- O PR não recebeu commits ou pushes de correção quando nenhum Review Issue
  existia.
