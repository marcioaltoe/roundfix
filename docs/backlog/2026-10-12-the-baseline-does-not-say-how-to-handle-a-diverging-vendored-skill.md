---
type: fix
status: open
created: 2026-10-12
spec: null
---

# The Baseline does not say how to handle a vendored skill that diverges from it

## Symptom

Spec 0199 states that repository rules take precedence over skill defaults.
On 2026-10-09 the maintainer of `marcioaltoe/skills` sent a fuller rule for
the Baseline to apply when it is assembled and validated. A vendored skill is
never edited locally, because the sync replaces the folder. When it diverges,
the options apply in this order:

1. An instruction in the `roundfix` skill or in the Baseline text. This is
   the default for convention differences: package manager, test runner,
   library choice.
2. A complementary authored skill that adds or overrides only what is
   missing.
3. An authored skill that replaces the vendored one, only when the vendored
   skill is useless or almost useless to the Baseline.

Objective upstream defects go to the upstream as PRs and never justify
option 3. Two such PRs are open: `tumf/skills#4` (`rust-cli`) and
`fellipeutaka/leon#3` (`rust-expert`). `oss-publish` exists upstream but was
not vendored.

## Where

Baseline guidance about vendored skills and the `roundfix` skill.

## Expected

State the three options, in order, where the Baseline and the `roundfix` skill
talk about vendored skills. Apply option 1 to `bun`, `hono`, `drizzle-orm`,
`golang-*` and the Rust skills where the Baseline still lacks the overriding
instruction.

## Evidence

- Secondbrain `inbox/roundfix/_triaged/2026-10-09-como-lidar-com-skills-vendorizadas-que-divergem-do-baseline.md`
- Secondbrain `inbox/roundfix/_triaged/2026-10-09-pr-upstream-das-skills-vendorizadas-e-dois-itens-que-nao-eram-defeito.md`
