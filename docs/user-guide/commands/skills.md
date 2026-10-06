### skills

```bash
roundfix skills list
roundfix skills check
roundfix skills install [--target codex|claude|opencode|all]
```

The binary ships 13 Roundfix-owned skills: the operational `roundfix` skill
plus the authorial workflow skills (`write-idea`, `write-prd`,
`write-techspec`, `write-tasks`, `setup-context-driven`, `implement-task`,
`implement-spec`, `brainstorming`, `business-analyst`,
`archive-spec`, `qa-gate`, `evidence-gate`). `skills list` also prints
recommended external skills, which install through your own tooling and are
never shipped. By default `skills install` writes to `<repo>/.agents/skills`;
`--target` selects user-scoped Agent skill directories instead.
