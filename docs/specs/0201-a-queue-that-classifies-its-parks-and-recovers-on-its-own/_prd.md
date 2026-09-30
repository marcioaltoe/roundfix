---
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
status: active
created: 2026-09-30
surfaces: [backend, cli, docs]
---

# A queue that classifies its parks and recovers on its own

The v0.22.0 Delivery Queue ran from 11:41 to 18:26 on 2026-09-30 and needed
twelve manual interventions. Only Spec 0199 went from start to merge through
the queue, and it still needed two manual steps. The finding
[2026-09-30-the-v0-22-0-queue-needed-twelve-manual-interventions.md](../../findings/2026-09-30-the-v0-22-0-queue-needed-twelve-manual-interventions.md)
groups them by class. Three classes are the queue's own:

- **Order is not a dependency.** The owner started 0195 while 0192 and 0193
  were parked. 0195's Verification named tests those two Specs create, so
  task_04 failed twice and the later rebase cost a corrective task_06.
- **A conflicting Pull Request looks like slow checks.** PR #298 was
  `CONFLICTING` after 0192 and 0199 merged. GitHub never started its checks,
  and the owner parked it `checks-timeout` thirty minutes later. Every
  conflicted file but one was derived and regenerable: the digest pin, the
  catalog snapshots, the four plan goldens and the Setup Manifest. The source
  conflict was a Baseline profile.
- **Environment failures park like defects.** 0192's QA gate was `partial`
  only because the QA sandbox's proxy refused GitHub. The operator satisfied
  the row and archived with a QA Archive Override. `deliver retry` then refused
  with "archived item head differs", and the item was delivered by hand (#296).
  Two time-bound tests failed under load in packages the failing items never
  touched.

A park today carries only a blocker string. `deliver status` prints it, and the
Pending Question gives a specific answer for two blockers and a generic one for
the rest. Nothing tells a dependency from a defect, or an environment failure
from a finding.

This Spec lets a Spec name its prerequisites, reports a conflicting Pull
Request at once and resolves it when only derived paths conflict, lets a retry
resume an operator-archived item, re-runs a failed check once when it failed
outside the item's change, and gives every park a class and a next command.

## Project Constraints

- Identifier strategy: applicable — no generated identifier is added. New
  stable names are added in the existing kebab-case forms: the blockers
  `prerequisite-unmerged`, `pull-request-conflict`, `qa-environment-partial`
  and `flaky-check`; the Park Classes `dependency`, `conflict`, `environment`,
  `flaky-check`, `finding`, `budget`, `review`, `authorization` and
  `unclassified`; the Task Graph manifest key `requires`; and the Project
  Config key `delivery.derived_paths`. Items keep their Spec slug, branch and
  Run ID. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git, the Run Database
  and the `gh` CLI the queue already uses, with the operator's existing `gh`
  session. Three `gh` reads or actions are added on the same repository (a Pull
  Request's `mergeable` field, a failed workflow run's attempt and log, and a
  re-run of its failed jobs). No credential is read, stored or printed, and no
  new endpoint or service is added. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0192 (this Spec) resolves a
  conflict confined to declared derived paths by regeneration, and ADR-0193
  (this Spec) lets a Spec name its prerequisites and has the owner wait for
  them. ADR-0149 keeps one regeneration declaration for grants, and ADR-0192
  adds a merge-recovery declaration that grants nothing. ADR-0178 keeps the
  audit of a Task commit on the grant it ran under. ADR-0154 makes a QA Archive
  Override a record of user authority, not a pass; a retry after it still runs
  the review, the repository gate and the checks. ADR-0165 keeps a blocking
  review after archive on a corrective Spec. ADR-0170 treats a Task already
  completed on its target as nothing to carry. ADR-0053 keeps terminal Run
  Worktree reconciliation explicit and proof-based. This Spec changes none of
  them. ADR-0158 keeps a budget-expired Run recoverable, and a `budget` park
  answers with a retry. ADR-0164 renews an Implement Run Budget at each Task
  settlement and is unchanged. ADR-0161 releases a merged
  Spec's Runs on the merged head. ADR-0090 forbids reusing a repository fact
  across a mutation, which is why the owner fetches the default branch before
  it reads a prerequisite or merges. ADR-0153, ADR-0169, ADR-0174 and ADR-0196
  govern the Pre-PR Review the resumed items pass through; ADR-0197 bounds its
  rounds. ADR-0160 keeps the frozen authorization as the only opener of a red
  repository gate. ADR-0081 makes sanctioned digest regeneration fallout of
  the authorized edit, and a derived merge regenerates the same way.
  ADR-0166 records a Task's undeclared paths, ADR-0167 keeps the pre-PR
  Pull Request row from deciding a qualifying partial, and ADR-0176 reads
  citations only from authored text. This Spec's gate is bound by ADR-0080,
  ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and
  ADR-0093 and ADR-0094 check its consistency by citation and artifact
  presence. ADR-0179 lets a grant list operations with its paths. ADR-0184
  has this TechSpec state its two command surfaces as Surface Transcripts.
  ADR-0097 cites ADR-0080 but carries a QA row forward, and this Spec's gate
  carries none. ADR-0168 cites ADR-0093 but narrows the related-ADR check,
  which this Spec does not change. ADR-0182 cites ADR-0096 but decides how the
  Daemon settles a Task, which this Spec does not change. ADR-0183 cites
  ADR-0093 but binds receipts to a Spec whose PRD follows its guide, and this
  PRD precedes it. ADR-0194 and ADR-0195 cite ADR-0097 but decide how the QA
  gate records and observes its rows again, and this Spec uses the gate as it
  stands. None of these six applies. ADR-0198 counts the tokens a Run's
  adapters report, and ADR-0199 sets a queue token ceiling; this Spec adds no
  limit and leaves both unchanged. The other ADRs present in this tree decide
  nothing about delivery, review or authorization and do not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer said on 2026-09-30
  "considere autorizado a ajustar todas as skills se necessário", which covers
  the Roundfix skill's delivery reference and its version, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync`.
  `.roundfixrc.yml`, the `Makefile`, CI workflows and `go.mod` do not change.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- The owner never starts a Spec before the Specs it names as prerequisites are
  merged into the default branch.
- A conflicting Pull Request is reported as a conflict when GitHub reports it,
  and a conflict confined to declared derived paths is resolved without an
  operator.
- An item the operator archived with a QA Archive Override, after an
  environment-only partial, can be retried instead of delivered by hand.
- Every park names its class and the next command, and a check that failed
  only outside the item's change is re-run once before it parks.

## User Stories

1. As an operator, I want a Spec to name the Specs it depends on, so that the
   queue does not start it on a main where its Verification cannot pass.
2. As an operator, I want a conflicting Pull Request reported as a conflict at
   once, with the conflicted paths, so that I do not wait for a check timeout.
3. As an operator, I want a conflict on derived files resolved by regeneration,
   so that parallel items that touch the Baseline merge without me.
4. As an operator who satisfied an environment-blocked QA row and archived with
   a QA override, I want `deliver retry` to continue delivery, so that I do not
   finish the review, gate, Pull Request and merge by hand.
5. As an operator reading `deliver status`, I want each park's class and the
   command that answers it, so that I know whether to wait, fix the
   environment, fix code or retry.
6. As an operator, I want a check that failed in a package the item did not
   change to be re-run once, so that one flaky test does not park an item.

## Core Features

1. **Prerequisite Specs.** A Spec's Task Graph manifest may carry a `requires`
   list of Spec slugs. A prerequisite is met when its archived Spec folder is
   on the refreshed default branch. The owner never starts an item while a
   prerequisite is unmet:
   - while the prerequisite is ahead in the same queue and not parked, the item
     stays `queued` and the owner advances other items;
   - when the prerequisite is parked, or is not in the queue, the item parks as
     `prerequisite-unmerged: <slug>, …`;
   - the owner returns such an item to `queued`, with no retry counted, as soon
     as every prerequisite is met, and a Delivery Retry does the same.
   `deliver start` refuses a queue whose `requires` entry names an unknown Spec
   or the Spec itself, or closes a cycle among the queued Specs (ADR-0193).
2. **A conflicting Pull Request is its own park.** While the owner waits for
   checks, it reads the Pull Request's `mergeable` state. `UNKNOWN` keeps it
   waiting; `CONFLICTING` stops the wait at once. The owner merges the
   refreshed default branch into the item branch in the item worktree:
   - when every conflicted path matches a declaration in the Project Config's
     `delivery.derived_paths`, it takes the default branch's bytes for them,
     runs each matched declaration's regeneration command, commits the merge
     if the regeneration changed only declared paths, and sends the new head
     through the repository gate, the push and the checks again;
   - otherwise it aborts the merge and parks as
     `pull-request-conflict: <path>, …`, naming the conflicted paths no
     declaration matches (ADR-0192).
   A Delivery Retry of that park accepts the operator's merge commits and
   resumes at the review.
3. **A retry resumes an operator-archived item.** A Run that ends unresolved
   with a newest QA Report of verdict `partial`, no finding-blocked row and at
   least one environment-blocked row parks as `qa-environment-partial` instead
   of `run-unresolved`. When the operator then archives the Spec with a QA
   Archive Override, `deliver retry` accepts the item if its head descends from
   the parked candidate (or, before any candidate, from the item's Run start),
   and resumes at the review. The archive stage sees the Spec already archived
   and continues to the repository gate. Every other archived item whose head
   moved is refused as today.
4. **Park Classes and the next command.** Every blocker maps to one Park Class:
   `dependency`, `conflict`, `environment`, `flaky-check`, `finding`, `budget`,
   `review`, `authorization`, or `unclassified` for a blocker the table does
   not know. `deliver status` prints one `Park:` line per parked item with its
   class and next command, and the Pending Question answers with the same next
   command.
5. **One re-run of a check that failed outside the change.** When a required
   check fails, the owner reads the failed log of its GitHub Actions run. When
   every failing Go package it names lies outside the paths the item changed
   against the default branch, and the run is on its first attempt, the owner
   re-runs its failed jobs once and keeps waiting. A second failure parks as
   `flaky-check: <package>, …`. A failure it cannot attribute, or one in a
   package the item changed, parks as `checks-failed`, as today.

## User Experience

`deliver status` keeps its table and its lines. After the Warning lines it adds
one line per parked item:

```text
Park: 0205-example conflict: merge the default branch into the item branch in /worktrees/0205-example, resolve internal/baseline/assets/profiles/standard-typescript-monorepo.json, commit, then run roundfix deliver retry 0205-example
```

The Pending Question's `Answer:` line prints the same next command for the
lowest parked item. The queue owner's console log names each owner action this
Spec adds: a waiting item, a dependency park released, a derived merge and its
regeneration commands, and a check re-run.

## Non-Goals / Out of Scope

- Fixing the two time-bound tests or the leaked detached child the finding
  names. They go to a Backlog Entry; this Spec only classifies and re-runs.
- Retrying any failure blindly, re-running a check more than once, or re-running
  a check that failed in a package the item changed.
- Rebasing an item branch, forcing a push or resolving a source conflict.
- Declaring this repository's own derived paths in `.roundfixrc.yml`. A binary
  older than this Spec refuses an unknown config key, so the declaration lands
  after a release holds this Spec.
- Showing prerequisites in `deliver plan`, or teaching the authoring skills to
  write `requires`.
- A Run Database schema change, a new command or flag, or a change to the retry
  limit.
- Accepting an operator-moved archived head for any park other than
  `qa-environment-partial` after a QA Archive Override and
  `pull-request-conflict`.

## Success Metrics

1. A queue holding a Spec that requires a parked Spec parks it as
   `prerequisite-unmerged` without creating its worktree, and returns it to
   `queued` with an unchanged retry count once the prerequisite's archive is on
   the default branch. A Spec without `requires` starts exactly as today.
2. A Pull Request reported `CONFLICTING` stops the check wait on the first
   read. A conflict confined to declared derived paths produces a merge commit
   whose declared paths equal a fresh regeneration and returns the item to the
   repository gate. A conflict on an undeclared path parks as
   `pull-request-conflict` and names exactly that path.
3. A retry of an item archived with a QA Archive Override after a
   `qa-environment-partial` park resumes at `reviewing`, passes the archive
   stage without a new archive commit, and reaches `gating`. An archived item
   whose head moved without that override is still refused with today's text.
4. Every blocker constant the delivery engine defines maps to a Park Class
   other than `unclassified`, and `deliver status` prints one `Park:` line for
   each parked item.
5. A failed check whose log names only packages the item did not change is
   re-run exactly once. A failure in a changed package, or with no Go package
   in its log, parks as `checks-failed` with no re-run.

## Declared breaks

- `deliver status` gains `Park:` lines for parked items. The test that pins the
  full status output of a parked item changes on purpose; the status of a queue
  without a parked item is byte-identical.
- The Pending Question answers from the Park Class table. Every blocker that
  exists today keeps its answer byte for byte; only the four new blockers get
  new answers.
- The delivery authorization is read at the parent of the commit that archived
  the Spec instead of at `HEAD^`. The two are the same commit whenever the
  archive commit is the head, which is every case today.

## Prerequisites

- Spec 0194 is delivered first. It splits the Roundfix skill and the command
  reference by command, and this Spec edits the delivery files that split
  creates: `.agents/skills/roundfix/references/deliver.md`, its mirror and
  `docs/user-guide/commands/deliver.md`.
- Specs 0189, 0190 and 0196 are delivered first, because they edit
  `internal/config/config.go` and the Roundfix skill before this Spec does.
- The queue that delivers this Spec runs an owner built before it, so that
  owner does not enforce this Spec's own ordering. The operator orders the
  queue.

## Recorded limits

- Package attribution reads Go test output only (`FAIL\t<import path>`). A
  check of another language, or one outside GitHub Actions, is never re-run.
- A path declared derived that also carries hand-edited source is resolved
  with the default branch's bytes plus regeneration. The repository gate that
  follows is the net, and a failed gate parks the item before anything is
  pushed.
- The owner releases a dependency park only while it runs. When every item is
  parked it exits, and the retry that restarts it releases the dependent item
  once the prerequisite merges.

## Decisions

- **Merge, never rebase.** The owner's conflict resolution is a merge commit
  pushed as a fast-forward. See ADR-0192.
- **No second review for a derived merge.** The merge adds only default-branch
  commits and declared regeneration output. An operator's hand resolution goes
  through the review again. See ADR-0192.
- **Prerequisites in the Task Graph manifest.** Dependencies already live in
  `_tasks.md`; a Spec's prerequisites join them. See ADR-0193.
- **A dependency park is released by the owner, not by a retry.** The item
  never started, so there is nothing to retry or carry. See ADR-0193.
- **Merged means archived on the default branch.** A squash merge of a
  delivered Spec always carries its archive, and a Spec merged by hand does
  too.
- **One re-run, attributed by package.** A failing test that runs none of the
  change is the published definition of a flaky failure; the re-run is bounded
  to one and to GitHub Actions runs still on their first attempt.
- **Declared derived paths, adopted later here.** The declaration is Project
  Config, so any repository can use it; this repository adopts it after a
  release, because older binaries refuse the key.

## Acceptance evidence

Each Core Feature needs positive and negative public-contract evidence in the
Task Graph. The outside-evidence rows rest on sources this Spec did not
produce:

- Recorded session evidence in the finding cited above: PR #298's conflicted
  file list and its `checks-timeout` park; Runs
  `run_20260930T162501Z_8ec6b68021df9cef` and
  `run_20260930T203751Z_2b4c20bff8fbb572` for 0195's failed task_04; and the
  hand delivery of 0192 (#296).
- GitHub's GraphQL reference defines `MergeableState`: `CONFLICTING` means the
  Pull Request "cannot be merged due to merge conflicts", and `UNKNOWN` means
  its mergeability "is still being calculated"
  (<https://docs.github.com/en/graphql/reference/enums#mergeablestate>).
- Bell et al., "DeFlaker: Automatically Detecting Flaky Tests", ICSE 2018
  (<https://dl.acm.org/doi/10.1145/3180155.3180164>), marks as flaky a newly
  failing test that executed none of the latest change, with a 1.5% false
  alarm rate over 4,846 failures.
- An independent project's contributor guide resolves a conflict in generated
  files by accepting the upstream copy and regenerating
  (<https://hummingbird-project.io/docs/contributing/resolve-merge-conflicts-generated/>).
- The GitHub CLI manual documents `gh run rerun --failed` and
  `gh run view --log-failed` (<https://cli.github.com/manual/gh_run_rerun>).

## Research basis

The Secondbrain was consulted through `wiki/index.md` and the queries
`qmd query "delivery queue park blocker classification retry"`,
`qmd query "merge conflict generated files regenerate rebase"` and
`qmd query "flaky test rerun CI"`. They returned this repository's mirrors of
Specs 0173, 0177 and 0180, the 2026-09-29 queue finding, the command guide and
the vendored `git-rebase` and `resolving-merge-conflicts` skills. None records
a queue design that classifies parks or resolves derived conflicts, so they add
nothing beyond the finding. Exa located the GitHub GraphQL `MergeableState`
reference, the DeFlaker paper, an independent project's rule for generated-file
conflicts, the GitHub CLI run manuals, and merge-queue vendor guidance that
bounds automatic re-runs to one to three. No name, text or code is copied from
any of them.

## Open Questions

- This repository's own `delivery.derived_paths` declaration waits for a
  release that holds this Spec. Until then a conflict on the Baseline's derived
  files parks as `pull-request-conflict` with its paths, which is still earlier
  and clearer than `checks-timeout`.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
