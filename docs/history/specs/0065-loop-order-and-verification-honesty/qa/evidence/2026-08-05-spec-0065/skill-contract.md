# Skill contract

- `rtk make skills-sync-check` — exit 0; 4 Skill integrity tests passed.
- `rtk bin/roundfix skills check` — exit 0; all 14 required owned Skills
  passed.
- `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` —
  exit 0.
- `rtk cmp -s .agents/skills/write-tasks/SKILL.md
  skills/write-tasks/SKILL.md` — exit 0.

Fresh text inspection confirmed:

- both `write-tasks` copies name the repository-wide-gates-plus-clean-tree
  refused shape, the same-subject contradiction rule, and the exact
  `## Rehearsal Cases` syntax;
- both `roundfix` copies list all four new stable `SC-*` identifiers;
- both canonical Skills state the settled loop order identically.
