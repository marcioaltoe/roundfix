---
type: refactor
status: open
created: 2026-11-16
spec: null
reason: null
waiting_on: "the maintainer's decision to schedule it (he said not now on 2026-11-16)"
since: 2026-11-16
---

# The documentation folders of adopters have unclear jobs, and reference files sit at the repository root

## Tangled

The maintainer listed on 2026-11-16 where the Baseline's documentation layout confuses adopters:

- `docs/agents/` holds rules, not agents, so its name misleads.
- `docs/backlog/` and `docs/findings/` both hold work to be done, with only a different scope, in two folders. He tried to merge them once before.
- `docs/design/` gets an API design in one repository and a front-end design in another; what belongs there is not clear.
- Almost every adopter keeps reference files such as `DESIGN.md` and `DATABASE.md` at the repository root instead of under `docs/`, for example in `docs/references/` or in a folder made for them.

## Cost

Each repository uses the Baseline without a common layout. An agent and a person have to learn each repository's placement, files end up in the wrong folder, and a Baseline update cannot rely on where things are.

## Shape

Non-binding: rename or regroup the folders so each name says its job, decide whether backlog and findings stay apart, state what `docs/design/` holds, give root reference files a home under `docs/`, and have `roundfix doctor` detect a layout that departs from the standard and offer the fix. The maintainer asked for this as future work, not now: "Não isso como algo a ser feito agora, mas eu quero que fique no backlog para futuro." He also doubted the benefit: "talvez a gente vai ter impacto zero de benefício."
