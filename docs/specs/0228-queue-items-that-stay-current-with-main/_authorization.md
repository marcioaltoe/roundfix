---
status: approved
granted: 2026-10-04
action: let the owned-skill record command raise a colliding version, let a derived merge resolve conflicts confined to declared version lines, let a Delivery Retry return a review-only correction after archive to review round 2, and describe the three rules in the Roundfix Skill, the implement-task skill, the guides and the repository's version rule
consuming: 0228-queue-items-that-stay-current-with-main
paths:
  - .agents/skills/implement-task/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/review.md
  - .roundfixrc.yml
  - docs/agents/specific-repository.md
  - skills/implement-task/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0228

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário", and of the Baseline source, the agent guides and
`.roundfixrc.yml` "Autorizar os dois". On 2026-10-04 the maintainer approved
this cycle's proposal of Specs O, P and Q with "Aprovado"; this Spec is P,
queue items that stay current with main.

The governed set was measured with `GovernedPath` on the authoring branch at
`7a1f564a`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.roundfixrc.yml` — the repository's record declaration gains the
  line-scoped version fields of every `SKILL.md`, so a derived merge resolves
  two items that raise the same skill.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, and the
  repository's skill-sync rule requires a Pull Request that changes CLI
  behavior to ship the skill update.
- `.agents/skills/roundfix/references/deliver.md` and
  `.agents/skills/roundfix/references/review.md` — the retry, derived-merge
  and review rules the Delivery Queue changes.
- `.agents/skills/implement-task/SKILL.md` and its mirror
  `skills/implement-task/SKILL.md` — a command a requirement names is part of
  the Task's work even when Verification runs the same test.
- `docs/agents/specific-repository.md` — the repository-owned rule that tells
  every author how an owned skill's version is raised and recorded.

## What is not governed

`skills/owned_skill_versions_test.go`,
`skills/owned_skill_version_raise_test.go`,
`skills/testdata/owned-skill-versions.json`, `internal/config/config.go`,
`internal/config/delivery.go`,
`internal/config/delivery_derived_lines_test.go`,
`internal/config/verification_tools_test.go`,
`internal/cli/deliver_workflow.go`,
`internal/cli/deliver_derived_lines_test.go`,
`internal/cli/deliver_review_correction.go`,
`internal/cli/deliver_review_correction_test.go`,
`internal/cli/review_lineage.go`, `internal/delivery/engine.go`,
`internal/delivery/review_correction_retry_test.go`,
`skills/roundfix/references/deliver.md`,
`skills/roundfix/references/review.md`,
`docs/user-guide/configuration.md`, `docs/user-guide/commands/deliver.md`,
`docs/user-guide/commands/review.md` and ADR-0233 are ordinary.

## Sanctioned regeneration

The repository-owned commands resolve the skill mirrors and the version
record. This declaration records the regeneration that follows the approved
edits and adds no source path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration or the CI workflows.
- No change to archived Specs, existing QA Reports, the Run Database schema,
  the review record or disposition ledger format, Park Class names, blocker
  names, retry limits, `CONTEXT.md` or the `### QA settlement` section of any
  skill.
- No test, Verification command or QA row reaches GitHub, a provider or the
  network, or writes under the real `~/.roundfix`; every artifact directory
  and Run Database a test reads lives in a disposable Roundfix Home.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
