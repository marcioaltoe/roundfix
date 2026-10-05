---
status: approved
granted: 2026-10-04
action: let roundfix reconcile and the delivery owner release the Runs of a Spec that merge evidence shows merged even when its files diverged later, let the post-merge cleanup prove a squash merge onto a moved default branch, let reconcile report and release the item branches of merged Specs, and describe this in the Roundfix Skill and the reconcile and deliver guides
consuming: 0227-reconcile-releases-the-runs-of-merged-specs
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/reconcile.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0227

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". On 2026-10-04 the maintainer approved this cycle's
proposal of Specs O, P and Q with "Aprovado"; this Spec is item O, reconcile
releases the Runs of merged Specs.

The governed set was measured with `GovernedPath` on the authoring branch at
`7a1f564a`, through a `go test -overlay` probe that wrote nothing to the
repository, against every file the Tasks declare.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/reconcile.md` — the reference says today
  that "Other committed paths still require content comparison" and lists
  three debris candidate kinds; merge evidence and item branch candidates
  change both.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  — an owned skill's content changes only with its version, raised in both
  front-matter fields, and the repository's skill-sync rule requires a Pull
  Request that changes CLI behavior to ship the skill update.

## What is not governed

`internal/worktree/merged_head.go`, `internal/worktree/worktree.go`,
`internal/worktree/merge_evidence_test.go`,
`internal/worktree/merged_spec_leftovers_test.go`,
`internal/worktree/item_branch_reconcile.go`,
`internal/worktree/item_branch_reconcile_test.go`,
`internal/cli/reconcile.go`, `internal/cli/reconcile_merge_evidence_test.go`,
`internal/cli/reconcile_item_branch_test.go`,
`internal/cli/deliver_workflow.go`,
`internal/cli/deliver_release_moved_main_test.go`,
`skills/roundfix/references/reconcile.md`,
`skills/testdata/owned-skill-versions.json`,
`docs/user-guide/commands/reconcile.md`,
`docs/user-guide/commands/deliver.md` and ADR-0232 are ordinary.
`internal/cli/cli_test.go` and `.agents/skills/roundfix/references/deliver.md`
are governed and are not touched.

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
- No change to archived Specs, existing QA Reports, the Run Database schema,
  the `roundfix-reconcile/v1` schema name, `CONTEXT.md` or the
  `### QA settlement` section of any skill.
- No test, Verification command or QA row reaches GitHub, a provider or the
  network, or writes under the real `~/.roundfix`; every Run Database a test
  reads lives in a disposable Roundfix Home, and every repository a test
  mutates is disposable.
- No release, tag or deployment.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
