---
spec: 0232-a-queue-that-sees-a-pull-request-merged-by-hand
status: archived
created: 2026-10-05
surfaces: [backend, cli, docs]
archived: "2026-10-05"
source_slug: 0232-a-queue-that-sees-a-pull-request-merged-by-hand
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: the Run sandbox could not reach GitHub or the web (rows 7a-7c, 7f) and the original 0231 queue row is gone (7e); the operator confirmed with gh on 2026-10-05 that #404, #382 and #396 are MERGED from their item branches (merged 2026-10-05T20:20:08Z, 2026-10-04T21:03:19Z, 2026-10-05T15:52:46Z); row 10 has no Pull Request yet. Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 477a1ebe287fcd8887123fdd4542a36552d40200
---


# A queue that sees a Pull Request merged by hand

This is a bug fix adopted from one Backlog Entry,
[a Pull Request merged by hand stays parked](references/2026-10-05-a-pull-request-merged-by-hand-stays-parked.md),
recorded on 2026-10-05 from the operator's intervention log. When the operator
merges a parked item's Pull Request by hand, the Delivery Queue never learns
it, and no command answers the item's Pending Question.

Measured cause, on the tree at `69462a0c`:

1. **The owner never looks at a parked item.** `Engine.Run` skips every item
   whose stage is `parked`, so `deliver resume` cannot notice a merge, and
   `deliver status` reads only the Run Database.
2. **A Delivery Retry never asks whether the work is merged.** `Engine.Retry`
   first requires the item branch, then applies the archived-head rules of
   ADR-0223 and ADR-0229. For Spec 0231 (Pull Request #404, merged at head
   `dca3aa19`, the operator's CI fix on top of candidate `09a30447`) the
   local item branch and worktree are gone, so the retry stops at
   `item branch ... is missing`. Had the branch remained, the moved head
   descends from the candidate, so the retry would return the item to
   `reviewing`; publication would then push the branch again and, because
   `FindOrCreatePullRequest` lists only open Pull Requests, open a second Pull
   Request for merged work. Specs 0225 (#382) and 0229 (#396) were refused
   and then merged through Pull Requests the queue never recorded, so their
   items carry no Pull Request number at all.

## Prerequisites

None. No other Spec is active, and no Verification here depends on another
Spec's artifacts.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; queue
  items, stages, blockers and Pull Requests keep their names and numbers.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — a Delivery Retry gains one GitHub
  read through the authenticated `gh` CLI the queue already uses, `gh pr
  view` of the item's recorded Pull Request, in the repository's own
  checkout. Roundfix reads, stores and sends no credential of its own, and
  the destination is the repository's own GitHub API. No test or Verification
  command reaches GitHub; every `gh` call in a test goes through a scripted
  command runner. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0237 (this Spec) decides the rule:
  "A Delivery Retry of a parked item now first asks whether the item is
  already merged". ADR-0232 supplies the merge evidence it reuses: "A Spec now
  carries merge evidence when the default-branch head holds its archived
  `_prd.md`", and its cleanup proof stays as it is. For every item that is
  not merged, ADR-0223 keeps its rule for a delivery that keeps its work, ADR-0223: "records that head as a new candidate". ADR-0229 keeps its operator-archive rule, ADR-0229: "It
  records that head as the candidate and resumes". The owner
  still releases waiting items after a merge, as ADR-0193 decides: "The owner
  itself returns such an item to `queued` once every prerequisite is met".
  ADR-0199 acts "when a parked item is retried"; ADR-0237 refines it so that
  recording a merge, which resumes no work, is not refused. ADR-0184 binds the changed retry output as Surface
  Transcripts. ADR-0179 bounds the Governed Paths in `_authorization.md`, and
  ADR-0189 and ADR-0233 bind the Roundfix Skill's version raise. The authored
  QA gate follows ADR-0080: "QA verdicts distinguish environment-blocked
  rows", and ADR-0091: "required to be terminal and to depend on every leaf";
  ADR-0096, ADR-0097 and ADR-0167 bind its machine stage, its row carry and
  its pre-PR Pull Request row, ADR-0104: "Every Spec therefore rests at least
  one named" acceptance row on outside evidence, ADR-0194, ADR-0195 and
  ADR-0210 bind what a QA row records, when it is observed again and its
  evidence snapshot, and ADR-0182 settles each Task on the facts its gate
  checks. ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check
  this Spec's consistency by citation and receipt. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill's deliver reference and
  both copies of its `SKILL.md` are Governed Paths. The maintainer's standing
  grant of 2026-09-30, "considere autorizado a ajustar todas as skills se
  necessário", covers them; this Spec changes no other Governed Path. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0232-a-queue-that-sees-a-pull-request-merged-by-hand/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Goals

- A Delivery Retry of a parked item whose recorded Pull Request was merged
  from the item branch records the item `merged`, and the owner runs the
  normal post-merge cleanup, whatever the park, the item head or the item
  branch's presence.
- A Delivery Retry of a parked item without that answer records it `merged`
  when the local default branch already archives its Spec through a delivery
  commit outside the item branch.
- A recorded Pull Request closed without merging, with no merge evidence,
  keeps the item parked and says why.
- Every item that is not merged retries exactly as before, and `deliver
  status` and `deliver resume` keep their behavior.

## Core Features

1. **A retry that sees the merge.** Before the Delivery Queue Limits and the
   item workspace, the retry reads the item's recorded Pull Request; if it is
   merged from the item branch, the retry records the item `merged` with that
   merge commit and the Pull Request's head as its newest candidate.
2. **Merge evidence for unrecorded Pull Requests.** Otherwise the retry looks
   for ADR-0232's merge evidence for the item's Spec and records the item
   `merged` with the delivery commit as its merge commit.
3. **The owner finishes the item.** The retry hands the queue to its owner as
   it does today; the owner's pass removes the item worktree and branch,
   releases the Spec's Runs and releases items waiting on it. The retry prints
   which evidence it used.
4. **The guides say it.** The deliver command reference and the Roundfix
   Skill's deliver reference describe the rule, its output and the
   closed-unmerged refusal.

## Non-Goals / Out of Scope

- Changing `deliver status`, `deliver resume`, the owner's pass or the
  post-merge cleanup's proofs.
- Probing GitHub for parked items without an operator's retry, or recognizing
  a merge in any stage other than `parked`.
- Reopening, closing or merging Pull Requests, or pushing to an item branch.
- The Run Database schema, the Makefile, `go.mod`, CI workflows and
  `CONTEXT.md`.

## Success Metrics

1. Success Metric: replaying Spec 0231 with a scripted GitHub that reports
   Pull Request #404 merged at a head that descends from the candidate, and
   with the item branch and worktree removed, `deliver retry` exits `0`,
   prints the merge line of Surface Transcript 1, records the item `merged`
   with that merge commit, leaves its retry count unchanged, and the owner's
   next pass leaves no cleanup warning (before: the retry is refused with
   `item branch ... is missing`).
2. Success Metric: replaying Spec 0225 in a disposable repository whose
   default branch archives the Spec through a squash commit the item branch
   lacks, with no Pull Request recorded, `deliver retry` records the item
   `merged` with that commit (before: the archived-head rules refuse or
   resume the item).
3. Success Metric: an item whose recorded Pull Request is closed without
   merging and whose Spec is not archived on the default branch is refused
   with the reason of Surface Transcript 2 and left unchanged.
4. Success Metric: every existing retry test passes unchanged, and a retry
   with no merge observer behaves exactly as before.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- GitHub's records of Pull Requests #404, #382 and #396, read with
  `gh pr view <n> --json number,state,headRefName,headRefOid,mergeCommit,mergedAt`
  on 2026-10-05: each is `MERGED`; #404's head is `dca3aa19` on the 0231 item
  branch and its merge commit is `4b5ea48b`.
- The operator's intervention log (`~/.roundfix-operator/queue-interventions.md`),
  entries 158, 169 and 175 to 177, and the Run Database's queue row for 0231
  (stage `parked`, blocker `checks-failed`, Pull Request `404`, candidates
  `ea6135c3`, `09a30447`), read on 2026-10-05.
- Mergify's merge queue documentation
  (<https://docs.mergify.com/merge-queue/deploy/>), read 2026-10-05: in
  hybrid mode "Mergify and manual merges coexist". Mergify's published
  merge-queue agent skill
  (<https://www.claudepluginhub.com/skills/mergifyio-mergify/mergify-merge-queue>),
  read the same day, lists the dequeue code `PR_MANUALLY_MERGED` as "Merged
  outside the queue" with the next action "Nothing": a normal exit, not a
  failure that waits for the operator.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "merge queue pull request merged manually outside the queue"`
(`--all --files --min-score 0.3`); it returned this repository's mirrors of
Specs 0175, 0187, 0211 and 0228 and of the decision on stale checks, none of which recognizes a
merge made outside the queue. Exa found Mergify's merge queue documentation and
its published agent skill, which treat a Pull Request merged outside the
queue as a normal exit, and
GitHub's merge queue documentation, which lists only removal reasons that
keep the Pull Request open. The one open Backlog Entry, a judge-assigned
model tier per Task, shares no context with this source, and there is no
unresolved Finding.

## Decisions

- Recognize the merge in `deliver retry`, the operator's explicit answer to a
  park, and not in `status` or in the owner's pass. See ADR-0237.
- Reuse ADR-0232's merge evidence for items without a recorded Pull Request.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
