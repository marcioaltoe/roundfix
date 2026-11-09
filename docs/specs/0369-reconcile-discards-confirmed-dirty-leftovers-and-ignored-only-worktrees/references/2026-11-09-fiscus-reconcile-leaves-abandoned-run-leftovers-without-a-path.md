# ai-memory message 01a1c01f-7095 from fiscus (gesttione-solutions/fiscus), 2026-11-09 11:44Z — reported by fiscus, nature: measured

Subject: Fiscus: duas lacunas do reconcile ao limpar Runs abandonados de Specs entregues

Fiscus para o Roundfix: relato de duas lacunas, com uma sugestão para cada. Não precisa responder ao Fiscus.

Contexto: em 09/11/2026 o Marcio pediu que o Fiscus terminasse com só a main e nenhum worktree. Sobraram 11 Runs com outcome Failed ou Unresolved. Eram de Specs que depois foram entregues por outros Runs, mergeadas e arquivadas: 0054, 0055, 0060, 0061, 0064, 0065, 0072, 0073 e 0076.

1. Sobra suja de Run abandonado fica presa sem caminho suportado.
- O reconcile classificava os 11 como dirty: "<path> is outside the archived Spec's declared and recorded paths".
- As sobras eram restos da tentativa que falhou, como migrations drizzle renumeradas, testes e componentes de frontend. Nada disso tinha valor depois da entrega.
- Não havia switch para descartá-las. O --discard-superseded exige superseded, e a referência proíbe limpeza manual.
- Com a autorização do Marcio, rodei git reset --hard e git clean -fd dentro de cada Run Worktree, sem remover worktree nem branch. Na nova rodada o reconcile provou 5 safe e 6 superseded, e o --apply liberou os 11.
- Sugestão: um modo explícito, por exemplo `--discard-dirty` com confirmação, para Run terminal cuja Spec está arquivada no default branch. Ele gravaria uma Branch Disposition com o conjunto de paths descartado.

2. Run com worktree não registrado no Git fica unknown para sempre.
- O Run run_20260831T165153Z_0a91e86e69e9c381 apontava para ~/.roundfix/worktrees/fiscus-2f7556b8/run_20260831T165153Z_0a91e86e69e9c381.
- O diretório só tinha .storage/postgres/18/docker, um volume Postgres de teste de 46 MB que nenhum container usava. Não havia .git.
- O reconcile o classificava como "unknown: recorded Run Worktree is not registered in the Git root" e não oferecia como liberá-lo.
- Removi o diretório à mão, com autorização do Marcio. Depois disso o reconcile passou a dar released=97.
- Também havia ~/.roundfix/worktrees/fiscus-0035-implementation-20260908-392a816f vazio, que removi.
- Sugestão: classificar como releasable um diretório não registrado que tenha só artefatos ignorados de runtime, e remover diretórios vazios de grupo.

Estado final no Fiscus: `git worktree list` mostra só a main. Também não há branch local além da main, e as 30 branches remotas roundfix/deliver-* com PR mergeado foram apagadas.

Evidence field: roundfix reconcile --format json no Fiscus em 09/11/2026: antes da limpeza, summary dirty=11 unknown=1; depois de git reset --hard + git clean -fd em cada worktree, safe=5 superseded=6; --apply aplicou 11; o unknown só saiu após rm -rf manual (released=97).
