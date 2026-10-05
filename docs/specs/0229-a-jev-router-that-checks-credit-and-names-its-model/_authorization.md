---
status: approved
granted: 2026-09-30
action: make the Jev Router gate read the OpenRouter account credit and refuse below the User Config floor jev.router_min_credit_usd, name an OpenRouter credit refusal openrouter_credit_refused, carry routed requests through a loopback relay that records the routed model on the Judge Log, and describe it in the Roundfix Skill and the guides
consuming: 0229-a-jev-router-that-checks-credit-and-names-its-model
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/runtime.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0229

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-04 the maintainer asked for the router work
to run unconstrained, "quero que seja tudo implementado e testado sem
limitações. Depois vemos o custo e como reduzir o custo.", and approved the
cycle that includes these cost controls: "Aprovado".

The governed set was measured with `GovernedPath` on the authoring branch at
`7a1f564a`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.
- `.agents/skills/roundfix/references/runtime.md` — the Jev Router gate's
  credit floor, the named credit refusal, the relay and the Judge Log fields.

## What is not governed

The Go sources and tests under `internal/config`, `internal/jevrouter`,
`internal/agent`, `internal/daemon` and `internal/cli` that the Tasks declare,
`skills/roundfix/references/runtime.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/configuration.md`, ADR-0218 and ADR-0234 are ordinary.

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
- No change to the `docs` or `chore` default profiles or any Recommended
  Profile; the router stays project-selected.
- No file is written outside the repository: the operator's User Config is
  set by the operator after merge.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  `CONTEXT.md` or the `### QA settlement` section of any skill.
- No test, Verification command or QA row sends a prompt to OpenRouter or
  reaches TypeSafe; every OpenRouter answer a test reads comes from a local
  stand-in, and no test reads or writes the real `~/.roundfix`.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.

## Maintainer decisions — 2026-10-04

Asked through AskUserQuestion on 2026-10-04, the maintainer approved the local
relay ("Aprovo o relay"): each routed session reaches OpenRouter through a
relay on 127.0.0.1 that forwards requests unchanged and notes the routed
model, provider and any HTTP 402, with the key held only in Roundfix's memory.
The maintainer set the credit floor's default to US$15 ("US$15").
