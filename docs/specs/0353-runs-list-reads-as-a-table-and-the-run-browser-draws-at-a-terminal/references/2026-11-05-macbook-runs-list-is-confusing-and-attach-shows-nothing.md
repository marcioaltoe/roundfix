# MacBook evidence (reported by the maintainer, 2026-11-05 ~06:50 BRT, pasted terminal output)

Host `macbook`, Roundfix v0.132.0, cwd `~/Google Drive/My Drive` (not a repository).

```
❯ roundfix runs list
run_20261105T094412Z_93cee2a36ff33269  ResolvingWithAgent  implement  spec:0063-um-gtin-por-mestre-e-a-pendencia-de-identidade  codex  2026-11-05T09:44:12Z  running 3m  roundfix/deliver-0063-um-gtin-por-mestre-e-a-pendencia-de-identidade-ef2e7f03322aafc2  /Users/marcio/.roundfix/worktrees/fiscus-2f7556b8/roundfix-deliver-0063-um-gtin-por-mestre-e-a-pendencia-de-identidade-ef2e7f03322aafc2-411942b2
run_20261105T093207Z_cf0dfe58d848a15f  ResolvingWithAgent  implement  spec:0063-sequencia-pela-emissao-propria-e-pdf-ao-cliente  codex  2026-11-05T09:32:07Z  running 15m  roundfix/deliver-0063-sequencia-pela-emissao-propria-e-pdf-ao-cliente-d59d8051b42f9794  /Users/marcio/.roundfix/worktrees/quaestor-85a8e4c6/roundfix-deliver-0063-sequencia-pela-emissao-propria-e-pdf-ao-cliente-d59d8051b42f9794-1189ed6c
(9 terminal Run Worktrees retained; run 'roundfix reconcile' to inspect)
❯ roundfix attach
❯ roundfix attach
❯
```

Maintainer (verbatim): "A lista de runs está confusa, talves com informação em excesso e o roundfix attach não mostra nada."

Notes: bare `roundfix attach` at an interactive terminal should open the Run Browser (alt screen, `internal/tui/runbrowser.go` `RunBrowserSession`), or print "missing run id in non-interactive mode…" on stderr. Neither appeared: it returned to the prompt with no output twice. Exit code and terminal (TERM/TERM_PROGRAM) not yet known — the operator asked the maintainer for them; they will be sent to you.
