---
origin: vortex
destination: roundfix
type-hint: fix
created_at: 2026-10-20
capture: manual
---

# O `roundfix archive --qa-override` no item worktree falha ao reescrever um link relativo do PRD

## Observação

Num item da Delivery Queue parado em `qa-environment-partial`, o `roundfix archive <slug> --qa-override
--approval … --reason …`, rodado na raiz do item worktree como manda a próxima ação do `deliver
status`, recusou com:

```text
Preflight failed

Reason:
  rebase archive outcome: rebase citation "../../backlog/2026-10-20-o-login-gesttione-liga-o-oraculum-no-claude-web-e-no-chatgpt.md": Rel: can't make docs/backlog/2026-10-20-o-login-gesttione-liga-o-oraculum-no-claude-web-e-no-chatgpt.md relative to ../../../../../../../Users/marcio/.roundfix/worktrees/vortex-25dcc288/roundfix-deliver-0146-o-better-auth-sobe-para-a-1-7-dd9a0784e236ffe8-88a2767f/docs/history/specs
```

O link (`[backlog de 2026-10-20](../../backlog/<arquivo>.md)` no `_prd.md`) existia e resolvia no item
worktree. A mensagem mostra a base do `filepath.Rel` como um caminho relativo que sobe até `/Users/…`,
como se a raiz do repositório tivesse sido resolvida a partir de outro diretório (o checkout principal
`~/dev/vortex`?), enquanto o alvo era relativo à raiz do repositório.

## Evidência

- roundfix 0.95.0 (9448aa64), macOS. Item worktree em `~/.roundfix/worktrees/vortex-25dcc288/roundfix-deliver-0146-…-88a2767f`,
  Spec `0146-o-better-auth-sobe-para-a-1-7` do vortex, em 2026-10-20 por volta de 08:00 BRT.
- Contorno usado: trocar o link por uma citação de caminho em texto (sem link Markdown) num commit no item
  branch. Depois disso o mesmo comando arquivou normalmente.
- O mesmo tipo de link (`../../backlog/...`) arquivou sem problema em Specs anteriores do vortex, cujo
  archive rodou dentro do Run, e não no item worktree.

## Por que importa

O caminho documentado para `qa-environment-partial` (carry-forward, reparo, `archive --qa-override`,
`deliver retry`) trava em qualquer Spec cujo PRD tenha link relativo para fora da pasta dela, e o
contorno obriga a mexer no texto da Spec só para contornar a ferramenta.

## Ação sugerida

Resolver a base da reescrita de citações a partir da raiz do repositório do checkout em que o `archive`
roda (o item worktree) e cobrir com um teste que arquiva, a partir de um worktree fora do checkout
principal, uma Spec com link relativo para `docs/backlog/`.
