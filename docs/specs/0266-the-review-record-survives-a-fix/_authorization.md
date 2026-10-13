---
status: approved
granted: 2026-10-13
action: key the pre-PR review record by the resolved checkout path and compare every record repository through it, keep a Reviewed Head Record for each non-reused verdict and make the Delivery Retry's correction proof read the parked candidate's record and also admit a changed path the archived Spec's Tasks declared, and describe both in the review and deliver guides, the roundfix skill's review and deliver references with the recorded version raise, and the glossary
consuming: 0266-the-review-record-survives-a-fix
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/roundfix/references/review.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0266

On 2026-09-30 the maintainer asked for unattended work through every release
of the program, and on skills: "considere autorizado a ajustar todas as skills
se necessário". On 2026-10-11, answering "Grants dos caminhos governados" for
cycle AE, the maintainer wrote "Concedo" for the Governed Paths each Spec
declares. On 2026-10-13 the maintainer approved cycle AG (six Specs, order B,
A, C, D, E, F) and then answered "Pode iniciar", under the same operating
rules as AE and AF: the same operations and the same grant for the Governed
Paths each Spec declares, without the Makefile or the CI workflows. Cycle AG
item B is this Spec: the writer and the correction proof key the review record
by the same resolved path, a later review round keeps the earlier round's
findings and dispositions, every reader and writer of the record is swept, and
the fix is characterized end to end under a worktree root behind a symbolic
link. Asked whether the retry should also accept a code correction, the
maintainer answered "Relaxar, com limites" (AskUserQuestion, 2026-10-13): a
changed path among the files the archived Spec's Tasks declared is admitted,
review round 2 and the repository checks stay mandatory, and anything else
still requires a corrective Spec.

The governed set was measured on the authoring branch at `main` v0.76.0 with
`speccheck.GovernedPath`, through a `go test -overlay` probe in
`internal/speccheck` that wrote nothing to the repository. It ran against
every path task_01, task_02, task_03 and task_04 declare, and four are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/review.md` covers the `review` and
  `review dispose` Behavior Surfaces whose sources `internal/cli/review.go`,
  `internal/cli/review_lineage.go` and `internal/cli/review_validation.go`
  change, and it states where the record lives; ADR-0256's authoring rule and
  `docs/agents/cli.md` require it to change with the command.
- `.agents/skills/roundfix/references/deliver.md` covers the `deliver`
  Behavior Surfaces whose source `internal/cli/deliver_review_correction.go`
  changes, and it states what the moved-head retry proves.
- `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md` carry the
  owned skill's version fields, which the record command raises with any
  content change.

## What is not governed

These paths are ordinary: `internal/cli/review.go`,
`internal/cli/review_lineage.go`, `internal/cli/review_validation.go`,
`internal/cli/deliver_review_correction.go`,
`internal/cli/review_resolved_checkout_test.go`,
`internal/cli/review_head_record_test.go`,
`internal/cli/review_declared_correction_test.go`,
`internal/spec/declared_paths.go`, `internal/spec/declared_paths_test.go`,
`docs/user-guide/commands/review.md`, `docs/user-guide/commands/deliver.md`,
`skills/roundfix/references/review.md`,
`skills/roundfix/references/deliver.md`,
`skills/testdata/owned-skill-versions.json`, `CONTEXT.md` and
`docs/references/behavior-surfaces.json`.

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
mirror and `skills/testdata/owned-skill-versions.json`. The Behavior Surface
Record is re-recorded after the guides change with
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`;
that record is not a Governed Path.

No Baseline module changes, so neither `make baseline-digests` nor the module
version record command is declared.

## Limits

- No change to the Makefile, the CI workflows, `go.mod`, `go.sum`,
  `.roundfixrc.yml`, or lint, formatter or test-runner configuration, and no
  change to `internal/delivery`, `internal/config`, any other skill, the Skill
  Coverage Map or the `### QA settlement` section of any skill.
- No rewrite of a recorded review record, disposition ledger line or Delivery
  Queue item, and no Run Database migration.
- No test, Verification command or QA row opens a network connection, reads
  or writes the real `~/.roundfix`, follows the maintainer's real worktree
  link, or calls a real reviewer.
- No release, tag or deployment. Verification stays Daemon-owned, and Task
  status stays Daemon-written.
