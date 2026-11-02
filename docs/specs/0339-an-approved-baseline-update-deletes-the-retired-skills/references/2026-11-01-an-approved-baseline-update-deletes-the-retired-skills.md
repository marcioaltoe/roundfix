---
type: fix
status: promoted
created: 2026-11-01
spec: 0339-an-approved-baseline-update-deletes-the-retired-skills
reason: null
---

# An approved Baseline update deletes the Retired Skills

## Symptom

`roundfix baseline update` lists each Retired Skill the repository still holds
under `Skills retired` with the paths to delete, and never deletes them
(ADR-0246, "A retired skill leaves the Baseline and the adopter removes its
copy"; `.agents/skills/roundfix/references/baseline.md`). The v0.127.0
CHANGELOG and the adopter note said the update removes them, and adopters who
trusted that kept the copies, the `.claude/skills` links and the
`skills-lock.json` entries.

## Decision (maintainer, AskUserQuestion, 2026-11-01)

"Apply apaga, valide se as skills context7 e exa estão fora da baseline e
sendo removidas, deveriam ficar." Option text: "Nova Spec substitui a
ADR-0246: o apply aprovado (--confirm-plan) apaga .agents/skills/<x>, o link
em .claude/skills e a entrada no skills-lock.json."

## Expected

- An approved apply (`--confirm-plan`) deletes exactly the skills listed under
  `Skills retired`: `.agents/skills/<skill>`, `.claude/skills/<skill>` and the
  `skills-lock.json` entry, and reports each deletion. A new ADR supersedes
  ADR-0246.
- Nothing outside that list is deleted. `context7-cli` and `exa-web-search`
  are in every Setup Snapshot (`internal/baseline/assets/setups/*.json`) and
  are not retired; a regression test proves they survive an apply, as does any
  skill the repository installed outside the Baseline.
- The v0.127.0 CHANGELOG wording is corrected in the release that ships this.

## Evidence

Reported by Argus (`default/gesttione-solutions-argus`, ai-memory message
`01a19718`, 2026-11-01, roundfix 0.127.0): the approved apply
(`--confirm-plan sha256:0cfdb20a…`) printed `Skills retired: 7 … delete
.agents/skills/<skill>, .claude/skills/<skill> and its skills-lock.json entry`
and deleted nothing; Argus removed them by hand in PR #15 (`72929f5`). The
operator confirmed on Vortex (2026-11-01) that the plan lists only the seven
retired skills and `implement-spec`, and that `context7-cli` and
`exa-web-search` are required by all five Setup Snapshots.
