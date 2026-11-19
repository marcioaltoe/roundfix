---
type: refactor
status: open
created: 2026-11-18
spec: null
reason: null
---

# `docs/agents/` holds rules under a name that says agents, and backlog and findings stay in two folders

## Tangled

The maintainer listed on 2026-11-16 where the Baseline's documentation layout confuses adopters. Spec 0410 took the list and states the job of each folder without renaming or moving anything. Two points of the list need a rename or a merge, and they stay here:

- `docs/agents/` holds rules, not agents, so its name misleads.
- `docs/backlog/` and `docs/findings/` both hold work to be done, with only a different scope, in two folders. He tried to merge them once before. ADR-0092 keeps them apart as intent and evidence, and it left "the possible future deprecation of `docs/findings/` into a backlog type" as "an explicitly deferred decision".

## Cost

A person or an Agent new to a repository reads `docs/agents/` as a folder about agents, and has to learn which of two folders takes a piece of work to be done.

Changing either is not small. Measured on 2026-11-18 in this repository at `d1c43f9ef26c64bf75000c71699e8bdaae784cf6`, as the number of files that name each path literally:

| Path | Go files, not tests | Go test files | Baseline modules and profiles | Files under `.agents/skills/` |
| --- | --- | --- | --- | --- |
| `docs/agents` | 18 | 104 | 18 | 17 |
| `docs/backlog` | 10 | 40 | 2 | 4 |
| `docs/findings` | 9 | 30 | 2 | 3 |

The count leaves out the derived Baseline files, which repeat every path, and the files of every adopter. Only the Spec Root can be configured (`internal/config/config.go:51`); the other paths under `docs/` are fixed in the code.

The maintainer doubted the benefit when he asked: "talvez a gente vai ter impacto zero de benefício."

## Shape

Non-binding. A rename of `docs/agents/` would reach every adopter through `roundfix baseline update`, as a move that keeps the old path readable for a time. One folder for backlog and findings would keep the difference between intent and evidence in a field of each file instead of in the folder, if ADR-0092's separation survives that. Each is its own decision, and neither should ride on a Spec about something else.

## Decisions

- 2026-11-18 (Idea of Spec 0410, AskUserQuestion): "Só a função de cada pasta (Recommended)" — o Baseline diz o que cada pasta guarda (inclusive docs/design/ e onde vão arquivos como DATABASE.md); o doctor aponta o que está fora; nada é renomeado nem movido; renomear/juntar volta ao backlog com o custo medido.
