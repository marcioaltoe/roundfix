# Roundfix 0.1.0 — incidentes do refino 0004/0005 e melhorias

Registro completo dos problemas encontrados ao executar as specs `0004-revisao-vs-verdade-e-catalogo`
(13 tasks) e `0005-exportacao-dinamica` (4 tasks) com `roundfix implement` + Codex `gpt-5.5 @ xhigh`,
em 2026-07-06/07. Resultado final: **17/17 tasks entregues, QA pass nas duas specs, ambas arquivadas**
— mas 6 runs morreram por problemas de ferramenta/autoria, todos recuperados sem perda de trabalho.
O procedimento operacional consolidado vive em `docs/agents/autonomous-work.md` → "Operational
procedure proven on roundfix 0.1.0".

## Incidentes (ordem cronológica)

| # | Incidente | Causa raiz | Contorno aplicado |
|---|---|---|---|
| 1 | Preflight: `_prd.md frontmatter status is ""` | Roundfix exige `status: active` no PRD; nossa convenção de spec não previa o campo | Estampar `status: active` nos PRDs ativos (e manter no `write-prd`) |
| 2 | Preflight: spec "not found" em `.knowledge/docs/specs/...` | O loader realpatha `docs/` e espera `docs/specs` na **raiz** do checkout `.knowledge`; não conhece o layout `projects/<name>/` | Symlink de compatibilidade `.knowledge/docs → projects/tax-poc/docs` (untracked, `/docs` no info/exclude) + recriação no `worktree.bootstrap` |
| 3 | Bootstrap do worktree: `docker compose up` conflita com containers `tax-poc-*` do host | `container_name` fixos no compose; segundo projeto compose no worktree tenta recriar os mesmos nomes | Bootstrap sem docker: `bun install && db:migrate && db:seed` reutilizando a infra do host |
| 4 | Bootstrap: `db:seed` aborta exigindo env vars da e-Auditoria | `worktree.copy` só copiava o `.env` raiz; as vars vivem em `packages/backend/.env` | `copy: [".env", "packages/backend/.env"]` |
| 5 | Budget default de 2h insuficiente para spec Runs | User Config `max_run_duration: 2h` vs. 10+ tasks × `make verify` | Project Config: `budget.max_run_duration: 8h` |
| 6 | **Bug principal**: Run morre no commit da task com `fatal: pathspec 'docs/specs/…' is beyond a symbolic link` — **após** a verificação passar | O daemon roda `git -C <worktree> add -- docs/specs/<slug>/task_NN.md <código…>`; o pathspec atravessa o symlink `docs` para outro repo git | Shim `~/.roundfix/bin-shim/git` que filtra pathspecs `docs/*` do `git add` (matou 3 runs — tasks 04/05/06 da 0004 — settladas manualmente antes do shim estabilizar) |
| 7 | Shim ignorado no primeiro uso | `--detach` re-executa o daemon **sem herdar o PATH** do chamador | Lançar sem `--detach`, com self-detach: `nohup env PATH="$HOME/.roundfix/bin-shim:$PATH" roundfix implement … &` |
| 8 | Shim v1 ainda ignorado | O daemon chama `git -C <path> add …` — o subcomando vem **depois** dos flags globais; o shim v1 só interceptava `add` como 1º argumento | Shim v2 percorre os flags globais (`-C`, `-c`, `--git-dir`, …) antes de identificar o subcomando |
| 9 | Run termina "implement failed" na **limpeza** do worktree, com tudo integrado | O reap usa `git worktree remove` sem `--force`; artefatos de bootstrap (`.knowledge`, `node_modules`, `.env`) são untracked e bloqueiam | Remoção manual `git worktree remove --force` + `git branch -D`; conferir o branch antes de assumir perda |
| 10 | Flips de status/QA evidence presos no `.knowledge` do worktree | O roundfix escreve `docs/specs/…` no clone do worktree, fora do Run Branch; some com o worktree | Monitor de sync periódico (worktree → origin, `projects/tax-poc` inteiro) + push manual no boundary |
| 11 | Flip da task_03 (0005) perdido → run seguinte re-executou a task do zero | Correção de status editada no worktree e o worktree removido **antes** do sync commitar (erro de sequência do orquestrador); runs novos re-executam tasks `failed` | Regra: flip + commit + push **antes** de qualquer remoção de worktree; sync passou a confirmar o push explicitamente |
| 12 | task_03 (0005) settlada `failed` 2× com o trabalho pronto e `make verify` verde | O task file exigia `make knip`, que sai com exit 2 por **baseline global pré-existente** (~332 achados de dead code não relacionados); `roundfix settle` não ajuda porque re-roda o mesmo comando verbatim | Settle manual do orquestrador com evidência (exit 2 também **sem** as mudanças; 0 resíduos no output; verify verde). Follow-up: sanear o baseline do Knip |
| 13 | task_04 (0005) settlada `failed` com 100% do trabalho entregue e QA pass | Critério de aceitação "docs commitadas no `.knowledge`" conflita com o contrato do agente ("never commit") | Settle manual; regra de autoria: commits/pushes nunca entram em critérios de task |

Menores: `roundfix stop` em run já terminal falha com exit 2 (esperado, mas ruim para automação);
eco de `tail -F` ao recriar console logs; preflight resolve o repo pelo cwd — lançar **sempre** da
raiz do repo de código (um lançamento a partir de `.knowledge` foi rejeitado por "branch main").

## Melhorias

### Upstream (roundfix)

1. **Suporte nativo a knowledge workspace**: resolver `docs/specs` através de symlinks/layout
   `projects/<name>/` e commitar task files no repositório git correto (ou ao menos pular pathspecs
   que cruzam symlinks em vez de matar o run pós-verificação).
2. **`--detach` deve preservar o ambiente** (PATH) do chamador.
3. **Reap de worktrees com `--force`** (ou limpeza prévia dos artefatos de bootstrap); nunca
   converter um desfecho Clean em Failed por erro de cleanup.
4. **Verificação de task com política de exit code** (allow-list/baseline) e `settle --skip-verify`
   ou verificação alternativa para recuperação supervisionada.
5. Reportar no stdout final o motivo por task (hoje o motivo do `failed` só existe no task file/log).

### Nossas (autoria e operação)

1. **Verificação de task deve ser hermética e satisfazível**: nada de `make knip` (ou qualquer
   comando com baseline sujo) até o baseline ser saneado; nada de critérios que exijam commit/push
   do agente.
2. Manter o procedimento operacional do `autonomous-work.md` (shim + nohup + sync + cleanup manual)
   enquanto o upstream não resolver #1–#3.
3. `write-prd` passa a estampar `status: active` por padrão.
4. Specs candidatas derivadas: saneamento do baseline do Knip; validação de path-params UUID
   (id inválido → 400, hoje 500).
