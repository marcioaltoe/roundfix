---
origin: roundfix
destination: conexus
type-hint: chore
created_at: 2026-10-01
capture: manual
---

# Roteiro: atualizar o Roundfix, as skills, o Baseline e o `.roundfixrc.yml` em conexus

Siga as etapas na ordem. Cada etapa diz como confirmar que deu certo. Pare na primeira que falhar e resolva antes de seguir. Repita este roteiro a cada release do Roundfix. A versão mínima é a **0.23.0**.

Perfil Baseline deste repositório: `standard-typescript-monorepo`.

## Etapa 0 — Validar o ambiente antes de começar

Rode na máquina que vai trabalhar neste repositório:

| Requisito | Comando | Esperado |
| --- | --- | --- |
| Git | `git --version` | qualquer versão atual |
| GitHub CLI autenticado | `gh auth status` | logado na conta com push, PR e merge neste repo |
| Node.js com npm/npx | `node -v` | 22.13 ou mais |
| Roundfix | `roundfix --version` e `roundfix upgrade --check` | 0.23.0 ou mais. Se houver versão nova, `roundfix upgrade` |
| Codex CLI logado (conta Pro) | `codex --version` | instalado e logado; o `doctor` prova a sessão |
| Claude Code logado | `claude --version` | instalado e logado; o `doctor` prova a sessão |
| Toolchain da Verification do repo | `bun --version`, `make --version` e o que a Verification usar (banco, Docker) | rode a Verification do repo à mão uma vez e confirme exit 0 |

Variáveis de ambiente (só leia do ambiente, nunca grave em arquivo versionado):

| Variável | Para quê | Obrigatória? |
| --- | --- | --- |
| `EXA_API_KEY` | skill `exa-web-search` (pesquisa) | sim, se a skill for usada |
| `CONTEXT7_API_KEY` | skill `context7-cli` (documentação de bibliotecas) | sim, se a skill for usada |
| `FIRECRAWL_API_KEY` | skills `firecrawl*` | opcional |
| `ROUNDFIX_OPENROUTER_API_KEY` | todo uso do OpenRouter pelo Roundfix, a começar pelo juiz consultivo `roundfix spec judge` (a partir da v0.25); chave própria para separar o custo do Roundfix | opcional; sem ela o juiz usa a chave TypeSafe ou pula sem bloquear |
| `ROUNDFIX_TYPESAFE_API_KEY` | todo uso direto da TypeSafe pelo Roundfix; no juiz, é a alternativa usada só se a chave acima faltar. As genéricas `OPENROUTER_API_KEY` e `TYPESAFE_API_KEY` nunca são lidas pelo Roundfix | opcional |

O `gh` e as ferramentas da Verification ainda não são verificados pelo `roundfix doctor` (há uma Spec planejada para isso). Confira-os à mão.

Depois rode o setup e o doctor do Roundfix:

```bash
roundfix setup            # verifica Node e acpx, prova os adapters e oferece User/Project Config
roundfix doctor           # todas as linhas devem sair ok
```

No `doctor`, confira:
- `node:`, `acpx:` e `adapter:` com `ok`;
- `profiles:` com `ok`, o que prova que Codex e Claude abrem sessão com os modelos configurados;
- `codex:` com `ok`.

## Etapa 1 — Atualizar as skills (comece por aqui)

As skills do Roundfix (`roundfix`, `qa-gate`, `write-prd`, `write-techspec`, `write-tasks`, `implement-spec` e outras) vêm dentro do binário. A skill `roundfix` mudou de estrutura na 0.23.0: agora é um arquivo de entrada mais uma referência por família de comando. Atualize-a antes de tudo, para que os agentes leiam o produto como ele é hoje.

```bash
git switch -c chore/roundfix-0.23-update
roundfix skills check                      # valida os artefatos que o binário carrega
roundfix skills install --target project   # grava em .agents/skills (e liga .claude/skills se faltar)
```

Confirme com `grep -m1 'version:' .agents/skills/roundfix/SKILL.md`. Deve mostrar a versão que o binário carrega: 0.1.3 na 0.23.0. O `marcioaltoe/skills` já tem a 0.1.5, que descreve a próxima release. Instale sempre pelo binário, para que a skill descreva a CLI que está rodando.

As skills externas, como `context7-cli`, `domain-modeling`, `tdd` e as do stack TypeScript, vêm do repositório `marcioaltoe/skills` no commit que o Roundfix fixa para o seu setup. O `marcioaltoe/skills` foi sincronizado em 01/10 com as skills do Roundfix. Os setups `typescript-bun`, `go-cli`, `go-tui` e `rust-cli` listam agora exatamente o que o Baseline exige. Restaure-as pelo comando de skills do Baseline. Ele mostra uma prévia e pede confirmação pelo digest:

```bash
roundfix baseline skills restore --profile standard-typescript-monorepo                          # prévia; sai 3 com o Plan Digest se houver mudança
roundfix baseline skills restore --profile standard-typescript-monorepo --confirm-plan <digest>  # aplica
```

Sem rede para o GitHub, use um clone local atualizado: `--source-dir ~/dev/skills`. O clone precisa conter o commit fixado (`git -C ~/dev/skills pull`).

Não instale versões mais novas dessas skills à mão. O `baseline update` compara cada skill com o digest fixado e propõe restaurar a versão fixada.

## Etapa 2 — Atualizar o Baseline

```bash
roundfix baseline update                 # mostra o plano, sem alterar nada
roundfix baseline update --yes           # aplica o plano calculado nesta execução
roundfix baseline update                 # deve reportar "Baseline update: current"
```

O que o update faz:
- reescreve só os blocos gerenciados, os guias do setup, as skills e o Setup Manifest;
- não toca nada fora dos marcadores, incluindo `docs/agents/specific-repository.md` e as decisões registradas.

O que chega nas releases recentes:
- **0.22.0:** as cláusulas descrevem o produto atual.
- **0.23.0:** as skills renomeadas (`context7` → `context7-cli`, `feature-systems-pattern` → `app-renderer-systems`, `rust` → `rust-expert`) e os guias citando só skills do setup.

Se o plano pedir uma decisão nova, decida conscientemente. `--adopt-suggested` só aceita a sugestão do catálogo, e o comando lista cada uma que aceitou.

A partir da v0.26.0, o layout de frontend vira uma decisão registrada. Enquanto nada for registrado, vale a sugestão `systems` e as cláusulas não mudam. conexus já nega warnings de lint na Verification.

## Etapa 3 — Atualizar o `.roundfixrc.yml`

```bash
roundfix profiles validate                          # leitura estrita: chave desconhecida é erro
roundfix profiles check                             # compara cada categoria com o Recommended Profile datado
roundfix profiles check --apply --scope project     # adota a recomendação nas categorias que diferem (com prova e prévia)
```

Ao adotar:
- se uma categoria precisar continuar diferente, registre uma Profile Deviation com o motivo, em vez de manter o desvio em silêncio;
- configs antigas costumam citar modelos que saíram (por exemplo `gpt-5.5`), e o `profiles check` mostra isso.

Revise também estas chaves:
- `budget.max_run_duration`: tempo máximo de um Run. O padrão é `2h`.
- `verification.repository_at_settlement`: o padrão é `true`. Cada Task de um grafo com gate de QA roda a Verification do repo ao assentar. Só desligue se a Verification for lenta demais, porque isso enfraquece o gate.
- `delivery.derived_paths`: **não** use ainda. A chave chega na v0.24, e um binário mais antigo recusa o arquivo.

O `.roundfixrc.yml` é um Governed Path. Mude-o num commit próprio e revisado.

## Etapa 4 — Verificar e entregar

```bash
roundfix doctor                          # tudo ok, incluindo skills: e recommendations:
roundfix baseline capabilities check     # capacidades do repositório; avisos "advisory" não bloqueiam
make verify                              # ou a Verification do repo; exit 0
```

Em seguida:
- abra um PR com as três mudanças (skills, Baseline e config), em commits separados;
- faça merge só com os checks verdes.

## Se algo der errado

- **`baseline update` recusa com "unaccounted clause".** É um defeito do Roundfix, não do repositório. Registre um inbox em `inbox/roundfix/` com a saída.
- **`profiles:` falha no `doctor`.** Confira o login do Codex e do Claude e rode `roundfix setup` de novo.
- **A fila (`roundfix deliver`) para.** Siga a resposta da Pending Question no `roundfix deliver status`. A partir da v0.24, cada parada também traz uma linha `Park:` com a próxima ação.
