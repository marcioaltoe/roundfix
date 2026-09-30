---
status: approved
granted: 2026-09-30
action: make an approved Spec-contained record whose paths is an explicit empty list an operations-only grant, and describe it in the delivery guides
consuming: 0188-a-grant-that-authorizes-delivery-without-governed-paths
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0188

The maintainer asked on 2026-09-30 to "Continue até o final para o release de
todas as implementações e ajustes": to deliver and release every remaining
implementation and adjustment of this cycle. The 2026-09-29 "Autonomia ampla"
standing authority covers authoring, corrective Tasks and delivery through
merge. The skill files ride the standing grant of 2026-09-18 for keeping the
shipped skills true to the CLI. The set was measured with `GovernedPath` on
the authoring branch at `9bd35366`, through a `go test -overlay` probe that
wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the Delivery
  queue section states what a start requires of a Spec's authorization, and
  must say how a Spec with no Governed Path records its grant (task_01).
  `.agents/skills/` is canonical and `skills/` its mirror.

## What is not governed

`internal/authorization/authorization.go`, the new tests
`internal/authorization/empty_paths_test.go`,
`internal/cli/deliver_plan_empty_paths_test.go` and
`internal/speccheck/mechanical_empty_grant_test.go`, and
`docs/user-guide/commands.md` are ordinary. The governed files
`internal/speccheck/constraints.go`,
`internal/speccheck/governed_repocontract_test.go`,
`internal/cli/cli_test.go` and `docs/references/coverage-record.json` are not
touched, and the setup-owned guides `docs/agents/spec-routing.md` and
`docs/agents/docs-layout.md` are not edited.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to archived Specs, existing QA Reports, ADR-0130's governed set,
  the operations vocabulary or any command's flags and exit codes.
- The live Run Database under `~/.roundfix` is never opened for writing by a
  Task or the gate, and no test reaches GitHub or a provider.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
