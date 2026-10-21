---
origin: pantheon
destination: roundfix
type-hint: fix
created_at: 2026-10-20
capture: manual
---

# O `roundfix archive` copia para o Archive Record os links relativos dos `satisfied-by` e o `spec check --strict` reprova

## Observação

Na Spec `0049-estatistica-por-consulta-e-alertas-do-cluster-vortex` do pantheon, o QA fechou `partial`
qualificado, com 12 linhas declaradas em `## Unreachable Acceptance`. Cada declaração tinha
`satisfied-by` com link relativo à pasta da Spec, como
`` [`references/registro-dos-alertas.md`](references/registro-dos-alertas.md) ``.

`roundfix archive <slug> --pull-request 91` gravou esses textos sem mudança no front matter do record
(`docs/history/specs/<slug>.md`, lista dos `satisfied-by` marcados `unproven`). Lá eles resolvem para
`docs/history/specs/references/...`, que não existe, e o `spec check --strict` falha com 15
`SC-HISTORY-RETAINED: ... links to a missing file or itself: docs/history/specs/references/<arquivo>.md`.

O `roundfix history sanitize`, em modo plano, lista as 15 reescritas (`record-links`). Ele aponta cada
arquivo removido para um Git Locator certo. O arquivo que o próprio archive tinha promovido para
`docs/references/medicao-de-sete-dias.md`, porém, ele reescreve para
`docs/history/specs/references/medicao-de-sete-dias.md` "(pull request #91)", um caminho que também não
existe. O `--apply` exige uma tag anotada `history-full`, que um repositório sem tags não tem.

## Evidência

- roundfix 0.97.0, macOS, pantheon, 2026-10-20 por volta de 16:40 BRT, branch
  `docs/0049-medicao-de-sete-dias` (PR #91). O archive também avisou que `source_revision` não é
  ancestral de `origin/main`, o esperado antes do squash.
- Contorno: as 15 reescritas foram aplicadas à mão no record, iguais ao plano do sanitize, com o
  arquivo promovido apontando para `../../references/medicao-de-sete-dias.md`. Isso gerou um commit
  próprio, e depois dele o `spec check --strict` deu exit 0.

## Por que importa

As Unreachable Acceptance são o caminho que a skill qa-gate recomenda para fechar `partial`
qualificado sem override, e o `satisfied-by` com link relativo é a forma natural de escrever. Todo
archive desse tipo sai reprovado pelo próprio strict. A correção só é possível editando à mão o front
matter que a regra manda não editar, ou com uma tag que o fluxo do repositório não usa. O archive
deveria reescrever esses links como já faz com o corpo do record: Git Locator para o arquivo removido e
link para `docs/references/` quando o arquivo foi promovido.
