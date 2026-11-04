---
type: fix
status: promoted
created: 2026-11-04
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
reason: null
---

# The vendored catalog skills trail the catalog's main and the lock names branches and old commits

## Symptom

This repository's `skills-lock.json` holds 33 entries with
`source: marcioaltoe/skills`. They name seven refs: `main` for 6 entries and
six older commits for the rest (`b3c45a4` ×10, `65d5f36` ×9, `a4e18e4` ×3,
`86df860` ×3, `b56fc5f` ×1, `14fdf46` ×1). An entry that names a branch cannot
be reproduced, and no command moves the 33 entries together. 32 of the 33
installed folders under `.agents/skills/` match the catalog's `main`
(`c8b2e85`); `systematic-debugging` still carries five files the catalog
removed: `CREATION-LOG.md`, `test-academic.md`, `test-pressure-1.md`,
`test-pressure-2.md` and `test-pressure-3.md`.

Reported by marcioaltoe-skills (the skills catalog repository).

## Where

`skills-lock.json` at the repository root, the vendored folders
`.agents/skills/<slug>` (aliased by `.claude/skills`), the five upstream
Setup Snapshots under `internal/baseline/assets/setups/` that pin `86df860`,
and the Makefile, whose `skills-update` target floats the entries through
another installer.

## Expected

The 33 catalog entries follow the catalog's `main` through one full commit,
with a `computedHash` that matches the installed folder; every installed
folder of the 33 equals the catalog's tree at that commit, so
`systematic-debugging` loses the five files; the 12 Roundfix-authored skills
(archive-spec, brainstorming, business-analyst, evidence-gate,
implement-task, qa-gate, roundfix, setup-context-driven, write-idea,
write-prd, write-tasks, write-techspec) stay out of the lock and are never
pulled from the catalog; the update is repeatable each cycle; and nobody
edits a vendored copy here: corrections go to `marcioaltoe/skills`.
`roundfix doctor` stays `ok`.

## Evidence

The ai-memory message below, received by this project on 2026-11-04,
reproduced verbatim. The operator re-measured its counts on `6a0ed904`
(v0.130.1) against the catalog's `main` `c8b2e85` the same day: the seven
refs and their counts, and git tree comparisons showing every folder but
`systematic-debugging` identical to `c8b2e85`; that folder differs only by
the five files and the mode of `find-polluter.sh` (`100755` upstream,
`100644` here), which neither the lock's `computedHash` nor the Setup
Snapshot's tree digest reads.

> # ai-memory message 01a1a608-3426-7bb1-884e-0c6f4b2975ac (from marcioaltoe-skills, 2026-11-04T10:08:37Z, nature: measured)
>
> Subject: Atualizar as 33 skills de terceiros pela main de marcioaltoe/skills (exceto as 12 autoradas pelo Roundfix)
>
> OBJETIVO: O Roundfix deve receber as atualizações das skills de terceiros que ele usa em `.agents/skills` a partir do catálogo marcioaltoe/skills, acompanhando a `main` do catálogo. A exceção são as 12 skills autoradas pelo próprio Roundfix (archive-spec, brainstorming, business-analyst, evidence-gate, implement-task, qa-gate, roundfix, setup-context-driven, write-idea, write-prd, write-tasks, write-techspec): continuam nascendo e sendo editadas no Roundfix, nunca puxadas do catálogo (o catálogo as copia de marcioaltoe/roundfix/.agents/skills/<slug>, ref main; puxá-las criaria ciclo). Hoje estão fora do skills-lock.json do Roundfix. Mantenha assim.
>
> ESTADO ATUAL (04/11/2026): skills-lock.json do Roundfix tem 33 entradas source marcioaltoe/skills (agentic-cli-design, ai-memory-* (6), bubbletea, coding-guidelines, context7-cli, conventional-commits, crafting-effective-readmes, domain-modeling, exa-web-search, github-pr-workflow, golang-* (9), grill-with-docs, grilling, knowledge-workspace, no-workarounds, systematic-debugging, tech-writer, testing-boss, tui-design, typesafe-ai). Só 6 usam ref main; 27 fixadas em commits antigos (65d5f36, b3c45a4, 86df860, a4e18e4, 14fdf46, b56fc5f). 32 das 33 pastas idênticas à main do catálogo (c8b2e85); divergente: systematic-debugging, com 5 arquivos que o upstream removeu (CREATION-LOG.md, test-academic.md, test-pressure-1.md, test-pressure-2.md, test-pressure-3.md).
>
> O QUE FAZER: 1) as 33 entradas acompanham a main (ref main ou sem ref) com computedHash coerente; 2) reinstalar as 33 do catálogo e confirmar que systematic-debugging perdeu os 5 arquivos; 3) atualização repetível (alvo de Makefile, passo de CI ou rotina); 4) não editar localmente as cópias dessas 33 skills — correções vão para marcioaltoe/skills.
> Atenção: o relato do baseline update que mantém o ref antigo (mensagem 01a1a5f7) é o mesmo problema visto do lado dos adotantes.
>
> RESULTADO ESPERADO: skills-lock.json com as 33 acompanhando a main do catálogo e sem nenhuma das 12 autoradas; .agents/skills/<slug> das 33 idêntico a skills/<collection>/<slug> na main do catálogo; as 12 intactas; roundfix doctor em ok.
>
> Evidence: medido contra marcioaltoe/roundfix main 6a0ed90 (v0.130.1) e marcioaltoe/skills main c8b2e85: skills-lock.json (33 entradas), git tree SHA de cada .agents/skills/<slug> comparado com HEAD:<local-path> do catálogo; skills-registry.json do catálogo.
