---
status: approved
granted: 2026-09-30
action: record the tokens and reported cost of every prompt a Run sends, print them in `runs show`, the Implement Run summary, the Run Event Stream and `deliver status`, add the `deliver start --max-tokens` ceiling, and describe them in the Roundfix skill
consuming: 0204-a-run-that-reports-the-tokens-and-spend-it-used
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/events.md
  - .agents/skills/roundfix/references/implement.md
  - .agents/skills/roundfix/references/runs.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0204

On 2026-09-30 the maintainer approved the program this Spec belongs to and
decided its content: Roundfix records the token usage the ACP adapters
report, prints it in `deliver status`, the Run's events and the Run summary,
and adds an optional per-queue token ceiling that parks the queue and never a
Run in the middle of a Task. A metering gateway may be used, but Roundfix must
not reimplement its logic, and its use stays an operator option in the user
guide. The same day the maintainer expressly authorized the skill files:
"considere autorizado a ajustar todas as skills se necessário".

The set was measured with `GovernedPath` on `30e8504f`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — a change
  to an owned skill's content raises both of its version fields, and the
  mirror's front matter is governed as well.
- `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/events.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/runs.md` — the Roundfix skill must
  describe the CLI behavior this Spec ships: the `Usage:` line and
  `--max-tokens` of `deliver`, the `usage` stream category of `events`, the
  `Tokens:` line of the Implement Run summary and the new `runs show`. Spec
  0194 creates these per-command files, and this Spec is delivered after it.
  All of the text goes under a new heading, `### Token usage`, in each file.

## What is not governed

Every Go source and test file this Spec's Tasks create or change under
`internal/agent`, `internal/store`, `internal/runevent`, `internal/daemon`,
`internal/delivery` and `internal/cli`, the user-guide files, the mirror
reference files under `skills/roundfix/references/` and
`skills/testdata/owned-skill-versions.json` are ordinary.

`internal/cli/cli_test.go`, `docs/references/coverage-record.json`,
`.roundfixrc.yml`, the skill manifests, `Makefile`, the CI workflows and
`go.mod` are governed and are not touched. The help text changes keep every
string `TestRunCommandHelp` and `TestEventsHelpDocumentsAgentSelectionFilter`
require.

## Sanctioned regeneration

The repository-owned commands resolve their generated outputs. These
declarations record the regeneration that follows the approved skill edits and
add no source paths.

```yaml
command: make skills-sync
```

```yaml
command: make baseline-digests
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, or the CI workflows.
- No change to archived Specs, existing QA Reports, `skills/_ownership.yml`,
  `CONTEXT.md` or the `### QA settlement` section of any skill.
- No code that proxies, meters or rewrites model traffic, and no change to
  the metering gateway.
- No test reaches a real ACP adapter, a provider, a gateway or the network,
  and the live Run Database under `~/.roundfix` is never opened for writing.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
