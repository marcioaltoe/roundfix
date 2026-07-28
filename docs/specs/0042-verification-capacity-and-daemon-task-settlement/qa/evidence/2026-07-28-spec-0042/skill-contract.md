# Skill contract evidence

Both comparisons exited 0:

```text
cmp .agents/skills/implement-task/SKILL.md skills/implement-task/SKILL.md
cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md
```

With the Task-local Go cache:

```text
make skills-sync-check
Go test: 4 passed in 1 packages
```

The public shipped-Skill validation reported all fourteen bundled Skills
passed. The authorial wording reserves Task status, declared Verification,
settlement, and Task commits for the Daemon; documents both capacities,
waiting/started phases, one deterministic repair, exit `75`, one exclusive
retry, and Task Type-selected Agent Sessions.
