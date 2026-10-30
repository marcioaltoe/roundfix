---
type: refactor
status: promoted
created: 2026-10-30
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
reason: null
---

# The Baseline follows the curated skills catalog that retired seven skills

## Tangled

The Baseline's Setup Snapshots pin the skills catalog `marcioaltoe/skills` at
`b56fc5ff`, and the `typescript`, `backend` and `frontend` modules still
require and dispatch seven skills the catalog removed on 2026-10-30. The
`typescript` setup's `ui-ux-pro-max` copy holds only its `SKILL.md`, whose
commands run scripts the copy does not carry.

## Cost

Until the Baseline follows the catalog, a TypeScript adopter keeps installing
and being told to activate skills the catalog no longer curates, and keeps an
incomplete `ui-ux-pro-max`.

## Shape

One sync of the bundled snapshots to the catalog commit `86df860`, the seven
names registered as Retired Skills, and every module, trigger and bundle that
names them cleaned. Non-binding.

## Evidence

ai-memory message `01a18d69-3e58-7a33-9917-70040f0b6eb5` from
`default/marcioaltoe-skills`, 2026-10-30 12:24 BRT, verbatim:

```text
Resposta aos pedidos 01a17a8a-77c6 e 01a17af0-64fe, com decisões de Marcio em 30/10/2026. Informativo, para o Roundfix fazer o sync único de snapshot junto com o ajuste de dispatch; não espera resposta.

Commit de main do catálogo marcioaltoe/skills: ed5b95b (inclui c59ec8f).

1. Removidas do catálogo (PR #141, commit c59ec8f). Todas estavam só no preset typescript, que era o único preset de cada uma:
- tdd -> removed (typescript)
- diagnosing-bugs -> removed (typescript)
- to-prompt -> removed (typescript)
- refactoring-analysis -> removed (typescript)
- design-patterns -> removed (typescript)
- game-changing-features -> removed (typescript)
- fixing-metadata -> removed (typescript)
Saíram de setups/typescript.txt, skills-registry.json, skills-registry.lock.json e skills/. A pendência de tdd/SKILL.md:38 citar `code-review` deixou de existir com a retirada do tdd.

2. ui-ux-pro-max NÃO foi retirada; foi completada com o pacote oficial (PR #142, commit ed5b95b). Continua em skills/03-engineering-design/ui-ux-pro-max e no preset typescript.
- Origem nova: nextlevelbuilder/ui-ux-pro-max-skill, path .claude/skills/ui-ux-pro-max, ref main, commit upstream 28b3b8dee8c9, tree 829c4dab53f1, licença MIT (LICENSE distribuído). Origem anterior: fellipeutaka/leon skills/ui-ux-pro-max.
- Agora traz data/, references/ e scripts/ (search.py, core.py, design_system.py). Patch local no SKILL.md: os comandos usam python3 "<skill-dir>/scripts/search.py" em vez de ${CLAUDE_PLUGIN_ROOT}/.claude/skills/..., que só existe em instalação de plugin. Por isso o update passou de auto para manual.
- Testado instalado em .agents/skills/ui-ux-pro-max: buscas reais com --domain ux, --design-system e --stack nextjs rodaram de outro diretório e responderam.
- Requisito novo: os scripts precisam de Python 3 (só biblioteca padrão). Os adotantes TypeScript precisam ter Python 3 na máquina para as buscas. Sem Python, a skill segue útil pelas regras do SKILL.md e de references/quick-reference.md. Está registrado no SKILL.md, na curation-note do registry e no cabeçalho de setups/typescript.txt.
- O agents/openai.yaml antigo saiu (o pacote oficial não tem).

3. Atualizações de upstream na mesma rodada:
- rust-expert sincronizada com fellipeutaka/leon (SKILL.md).
- gog: pin de openclaw/gogcli de v0.47.0 para v0.47.1. A pasta da skill é igual entre as tags; o folderSha do lock não mudou.
- Sync de todos os modos auto e manual: nenhuma outra skill com drift.
```
