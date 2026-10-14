---
status: approved
granted: 2026-10-13
action: reject a turn whose stream reports a provider system error and recover it through the Fallback Chain under the Lost Rollout rule, park an unrecovered one as runtime infrastructure, refuse Claude background tasks in Run sessions and settle a QA turn that leaves one as an Unfinished Turn, and describe them in the glossary, the guides and the roundfix skill's profiles and deliver references with the recorded version raise
consuming: 0268-a-turn-that-did-not-finish-is-not-accepted
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/profiles.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0268

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, and on skills: "considere autorizado a ajustar todas as skills
se necessário". On 2026-10-11, answering "Grants dos caminhos governados" for
cycle AE, the maintainer wrote "Concedo" for the Governed Paths each Spec
declares. On 2026-10-13 the maintainer approved cycle AG (six Specs, order
B → A → C → D → E → F) and answered "Pode iniciar" to start them under the
same operating rules as AE and AF, which carries the same operations and the
same grant for the Governed Paths each AG Spec declares. Cycle AG item C is
this Spec: a turn that did not finish is not accepted.

The governed set was measured on the authoring branch at `9963ffd3` with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository. It ran against
every path task_01 to task_04 declare under `interface:` or `creates:`, and
four are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/deliver.md` covers the `deliver`
  Behavior Surfaces whose source `internal/cli/deliver_workflow.go` changes;
  ADR-0256's authoring rule and `docs/agents/cli.md` require the skill to
  change with the command, and it gains the Provider System Error park.
- `.agents/skills/roundfix/references/profiles.md` describes the Fallback
  Chain's Lost Rollout exception, which this Spec extends to the Provider
  System Error.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  owned skill's version fields, which the record command raises with any
  content change.

## What is not governed

Every other declared path is ordinary: `CONTEXT.md`,
`docs/user-guide/configuration.md`, `docs/user-guide/usage.md`,
`docs/user-guide/commands/deliver.md`, the two reference mirrors under
`skills/roundfix/references/`, `skills/testdata/owned-skill-versions.json`,
and the Go sources and tests under `internal/agent`, `internal/daemon` and
`internal/cli` that task_02 to task_04 declare.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors after the approved
skill edit and adds no source path.

```yaml
command: make skills-sync
```

task_01 also runs the owned-skill record command, whose package installs no
suite guard, so it needs no declared output list here:
`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
rewrites the version lines of `.agents/skills/roundfix/SKILL.md` and its
mirror and `skills/testdata/owned-skill-versions.json`.

No Baseline module changes, so neither `make baseline-digests` nor the module
version record command is declared, and the Behavior Surface Record keeps its
bytes.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`,
  `.roundfixrc.yml`, or lint, formatter or test-runner configuration.
- No change to the `qa-gate` skill, any `### QA settlement` section,
  `docs/references/behavior-surfaces.json`,
  `docs/references/skill-coverage.json`, archived Specs or Archive Records.
- No test, Verification command or QA row starts a real Agent Session, sends a
  request to a model provider, opens a network connection, reads `~/.acpx`,
  reads or writes the real `~/.roundfix`, or persists an environment value.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
