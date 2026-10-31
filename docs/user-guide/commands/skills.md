### skills

```bash
roundfix skills list
roundfix skills check
roundfix skills install [--target codex|claude|opencode|all]
```

The binary ships 12 Roundfix-owned skills: the operational `roundfix` skill
plus the authorial workflow skills (`write-idea`, `write-prd`,
`write-techspec`, `write-tasks`, `setup-context-driven`, `implement-task`,
`brainstorming`, `business-analyst`,
`archive-spec`, `qa-gate`, `evidence-gate`). `skills list` also prints
recommended external skills, which install through your own tooling and are
never shipped. By default `skills install` writes to `<repo>/.agents/skills`;
`--target` selects user-scoped Agent skill directories instead.

Inside an adopted Git repository, `skills check` also inspects the skills the
repository's Baseline requires and fails with exit `1` for each Missing
Required Skill, printing
`.agents/skills/<name>: required by the Baseline and not installed; restore with roundfix baseline skills restore --profile <id> --skill <name>`
on stderr. Outside a Git repository, or without a Setup Manifest, it checks
only the shipped bundle.

After bundle installation, an existing `skills-lock.json` entry naming an
owned skill is updated with the installed folder's computed hash; every other
field is preserved. No entry or lock is created when no matching lock entry
exists. Each changed name is reported as `Updated skills-lock.json computedHash:
<name>`.
