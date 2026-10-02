Source: /Users/marcio/dev/secondbrain/inbox/conexus/2026-10-01-roteiro-de-atualizacao-do-roundfix-e-validacao-do-ambiente.md

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


Observation: the published adopter instruction explicitly required manual forge/toolchain checks; the new coded Doctor output replaces those manual diagnoses. Its key presence guidance remains supported by the sentinel public flow. The source does not contain the original maintainer request.
