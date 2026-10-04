---
status: approved
granted: 2026-09-30
action: make the monthly Jev ceiling the User Config value jev.monthly_ceiling_usd with the US$5 default, read by the judge, the Jev Router gate and its key-limit check, ignore it in Project Config, and describe it in the Roundfix Skill and the guides
consuming: 0225-a-jev-ceiling-the-maintainer-sets
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/runtime.md
  - .agents/skills/roundfix/references/spec.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0225

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-04 the maintainer asked for this change:
"A chave do openrouter está com limite de $50 mensal e quero que seja tudo
implementado e testado sem limitações. Depois vemos o custo e como reduzir o
custo."

The governed set was measured with `GovernedPath` on the authoring branch at
`6ea9e1e9`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.
- `.agents/skills/roundfix/references/spec.md` — the judge's ceiling.
- `.agents/skills/roundfix/references/runtime.md` — the Jev Router gate's
  ceiling and key limit.

## What is not governed

The Go sources and tests under `internal/config`, `internal/judge`,
`internal/daemon` and `internal/cli` that the Tasks declare,
`skills/roundfix/references/spec.md`, `skills/roundfix/references/runtime.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/configuration.md`, `docs/user-guide/commands/spec.md`, and
ADR-0201, ADR-0218 and ADR-0231 are ordinary.

## Sanctioned regeneration

The repository-owned command resolves the skill mirror. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, the CI workflows or `.roundfixrc.yml`.
- No file is written outside the repository: the operator's User Config is
  set by the operator after merge.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  `CONTEXT.md` or the `### QA settlement` section of any skill.
- No test, Verification command or QA row reaches OpenRouter, TypeSafe or the
  network, or reads or writes the real `~/.roundfix`; every User Config and
  Judge Log a test reads lives in a disposable Roundfix Home.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
