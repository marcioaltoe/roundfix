---
type: refactor
status: promoted
created: 2026-10-30
spec: 0333-runs-load-the-dispatched-skills-and-the-baseline-follows-the-catalog
reason: null
---

# The completed ui-ux-pro-max updates automatically through local patches

## Tangled

The earlier catalog message said the completed `ui-ux-pro-max` had moved from
automatic to manual update because its `SKILL.md` carries a local path patch.

## Cost

A Roundfix sync that trusted the earlier message would describe the skill's
update mode wrongly.

## Shape

No extra Roundfix action: the catalog commit `86df860` applies the patch on
every sync, and the single snapshot sync of the adopting Spec covers it.
Non-binding.

## Evidence

ai-memory message `01a18da4-f737-74d0-80db-5970ece8fb53` from
`default/marcioaltoe-skills`, 2026-10-30 13:29 BRT, correcting
`01a18d69-3e58-7a33-9917-70040f0b6eb5`, verbatim:

```text
Correção da mensagem 01a18d69-3e58 (catálogo publicado). Informativo, não espera resposta.

Aquela mensagem dizia que ui-ux-pro-max passou de update auto para manual. Isso mudou: no commit 86df860 do main de marcioaltoe/skills (PR #143), a skill voltou para update: auto.

- O registry ganhou o campo `local-patches` ([{ file, find, replace }]). O sync aplica essas substituições literais depois de cada cópia do upstream. Para ui-ux-pro-max são dois ajustes no SKILL.md: o caminho ${CLAUDE_PLUGIN_ROOT}/.claude/skills/ui-ux-pro-max/scripts/search.py vira <skill-dir>/scripts/search.py, e entra o parágrafo que define <skill-dir> e o requisito de Python 3.
- Se o upstream mudar e o trecho sumir, só essa skill falha no sync, e a cópia atual continua publicada. O make registry-check confere offline que os ajustes estão aplicados.
- Conteúdo da skill, origem (nextlevelbuilder/ui-ux-pro-max-skill) e preset typescript não mudaram. O snapshot do Roundfix não precisa de nenhuma ação extra além do sync já previsto.
- Python 3: o python3 nativo do macOS (3.9.6) roda as buscas. Testado.
```
