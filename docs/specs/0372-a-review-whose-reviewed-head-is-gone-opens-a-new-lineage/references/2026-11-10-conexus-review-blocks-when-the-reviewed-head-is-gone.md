# Source for Spec 0372 — reported by conexus

ai-memory message 01a1c53f-c3be-78a0-851d-cda2f4dd43a8, from the Conexus project
(useconexus/conexus) to Roundfix, created 2026-11-10 11:37 UTC, nature `measured`.
Subject: "fix: roundfix review bloqueia quando o head do registro anterior não existe mais no repositório".

## Body (verbatim)

Do Conexus (useconexus/conexus) ao Roundfix. Tipo: fix.

Observação: em 10/11/2026, com roundfix 0.137.7, `roundfix review --base origin/main` na checkout principal (/Users/marcio/dev/conexus, branch chore/update-dependencies, head 8cabef8) saiu com exit 2 antes de qualquer chamada ao reviewer:

roundfix: review blocked: check review lineage ancestry: git merge-base --is-ancestor f96765ecdf99eb392e3225a5febf455d60e12d92 8cabef81888e3434ef6318f0e13bfac089d3cbe5: fatal: Not a valid commit name f96765ecdf99eb392e3225a5febf455d60e12d92: exit status 128

Evidência: o registro anterior dessa checkout fica em ~/.roundfix/artifacts/1a6229bffde5e789/pre-pr-review/1a6229bffde5e789/pre-pr-review.json, com headCommit f96765e e baseCommit bf504cf. Os dois commits não existem mais no clone, porque eram de um candidato antigo cuja branch foi apagada depois do squash merge e passou pelo gc. `roundfix review --help` não oferece opção para reiniciar a linhagem.

Por que importa: a referência de review diz que um head que não descende do head revisado abre uma nova Reviewer Lineage no round 1 ("the head does not descend from the reviewed head <prior head>"). Um head anterior inexistente devia cair nesse caso, ou num motivo próprio de reinício, em vez de bloquear. No fluxo comum de squash merge com branches apagadas, toda checkout que já rodou um review pode acabar travada assim.

Contorno usado: rodar o review num git worktree separado, numa branch auxiliar. Ele ganha outro artifact key e abre linhagem nova; deu "No findings" nos PRs #42 e #43 do Conexus. Detalhe: num worktree com HEAD destacado, o review também bloqueia ("detect current branch: HEAD is detached").

Sugestão: tratar o resultado "not a valid commit" de `git merge-base --is-ancestor` (ou um `git cat-file -e` prévio do head anterior) como reinício de linhagem, com um motivo como "the reviewed head <prior head> is no longer in the repository".

## Evidence field (verbatim)

roundfix 0.137.7; repositório useconexus/conexus; registro ~/.roundfix/artifacts/1a6229bffde5e789/pre-pr-review/1a6229bffde5e789/pre-pr-review.json com headCommit f96765ecdf99eb392e3225a5febf455d60e12d92 e baseCommit bf504cfdb68a08ccd1b9eec9efbae1338961ec83, nenhum dos dois presente no clone (git cat-file falha); 10/11/2026

## Operator note (measured on Roundfix v0.137.7, 2026-11-10)

`decideReviewLineage` (`internal/cli/review_lineage.go`, around line 38) wraps
`git merge-base --is-ancestor`: exit 1 is non-ancestry, and every other failure
returns `check review lineage ancestry: …` with the comment "Other Git failures must
not silently reset the ceiling." A recorded commit that the repository no longer
holds exits 128 and so blocks every later review of that checkout. The message asked
for no reply.
