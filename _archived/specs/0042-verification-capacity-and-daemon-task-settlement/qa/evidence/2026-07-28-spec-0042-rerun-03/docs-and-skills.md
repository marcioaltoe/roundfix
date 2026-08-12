# Documentation and Skills

Current built help exited 0 for the root, Implement, and Attach commands.
Implement exposes no capacity flag and retains the profile-led Agent selection
contract. Attach documents:

```text
roundfix attach [<run-id>] [--no-input]
```

The documentation contract tests passed. The configuration guide defines
Verification Capacity default 1, strict positive validation, built-in to User
Config to Project Config precedence, independence from Task Capacity, and
per-Run scope. Command and usage guidance describe waiting/started events,
Daemon-only settlement, one Agent repair, project-authored exit 75, one
exclusive retry, retained diagnostics, exhaustion, and no log heuristic.

Task 07 commit inspection found no protected Skill path. Its documentation and
Agent Session wording match ADR-0051.

Both protected canonical/generated pairs compare byte-identically:

```text
cmp .agents/skills/implement-task/SKILL.md skills/implement-task/SKILL.md
cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md
```

Both exited 0 with no output. With the writable Go cache,
`make skills-sync-check` passed four policy tests, and
`go run -buildvcs=false ./cmd/roundfix skills check` passed all fourteen
shipped Skills. Task 08's exact commit paths match the PRD and TechSpec
authorization.
