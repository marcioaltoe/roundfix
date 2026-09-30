---
spec: 0187-a-queue-that-recovers-without-a-supervisor
status: active
created: 2026-09-30
surfaces: [backend, cli, docs]
---

# A queue that recovers without a supervisor

The first two full Delivery Queue trials, on 2026-09-29, still needed a
Supervisor to finish their items. Four of the remaining interventions came
from Roundfix reading a stale or partial fact:

- After the queue squash-merged Spec 0183 (#282), releasing the Spec's Runs
  failed with `merge commit "f3cbf0c4…" does not resolve to a commit`, and
  later `is not on default branch`. The queue checked the recorded merge commit
  against local refs it had never refreshed. The item worktree and the Runs
  were left for a manual `reconcile`.
- A `deliver retry` after a Spec amendment on the item branch refused the whole
  carry-forward set, because every Task's declared inputs had moved. Its next
  action named only `reconcile --carry-forward`, which refuses the same way
  while the amendment sits under it.
- The Wave 5 queue owner was built before Spec 0181 merged, and its items
  started from a main that held 0181. The owner kept the defects 0181 fixed for
  the whole queue, which caused four parks, and nothing said the owner was
  older than its items' main.
- The QA authorization audit refused Spec 0181's task_07 even though the
  widened grant was on main (#277) and on the item branch before the Task ran.
  The audit reads the grant where the item branch forked from main, so only a
  rebase fixed it.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; a Delivery Queue
  item keeps its Spec slug, branch and Run ID, and the owner warning names
  existing commit IDs. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git (including one
  `git fetch` of the delivery remote the queue already fetches) and the Run
  Database only; no credential is read and no new network endpoint is added.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0178 (this Spec) accepts the grant
  a Task commit ran under when the delivery target's history holds the same
  record. ADR-0161 releases a merged Spec's Runs on the merged head, and this
  Spec only refreshes the default branch first. ADR-0170 treats a Task already
  completed on its target as nothing to carry, and ADR-0053 keeps
  reconciliation proof-based; the retry refusal keeps every proof and only
  names its cause. ADR-0090 forbids reusing a repository fact across a
  mutation, which is why the cleanup refreshes the default branch before
  reading it. ADR-0158 keeps a BudgetExceeded Run recoverable through
  carry-forward. ADR-0160 keeps the frozen authorization as the only authority
  that opens a red repository gate, and ADR-0081 and ADR-0149 keep sanctioned
  regeneration outputs inside a grant; none changes. ADR-0166 records a Task's
  undeclared paths, and every Task here declares its paths. ADR-0167 keeps the
  pre-PR Pull Request row from deciding a qualifying partial; this Spec's gate
  aims at `pass`. ADR-0169 reviews a candidate from its merge base. ADR-0176
  reads citations only from authored text. This Spec's gate is bound by
  ADR-0080, ADR-0088, ADR-0091, ADR-0096, ADR-0104, ADR-0117, ADR-0155 and
  ADR-0156, and ADR-0093 and ADR-0094 check its consistency by citation and
  artifact presence. ADR-0097 cites ADR-0080 but carries a QA row forward, not
  a grant. ADR-0164 cites ADR-0158 but renews a Run's budget at each Task
  settlement. ADR-0168 cites ADR-0093 but narrows the related-ADR check. This
  Spec changes none of them, so they do not apply. All the others hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer asked on 2026-09-30 to
  "Continue até o final para o release de todas as implementações e ajustes";
  the Roundfix skill files ride the standing grant of 2026-09-18 for keeping
  the shipped skills true to the CLI, recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A merged item's Runs and worktree are released without a warning, whatever
  the state of the operator's local default branch.
- A retry refused because the item branch amended the Spec tells the operator
  how to recover it, with the exact commands.
- An operator learns, while the queue runs, that its owner binary is older than
  the Roundfix source its items started from.
- A grant that the default branch holds, and that a Task ran under, authorizes
  that Task's commit whether the item branch merged or rebased main.

## Core Features

1. **Cleanup refreshes the default branch first.** Before it releases a merged
   Spec's Runs, the queue fetches the default branch from the delivery remote
   when that remote exists. It then checks the recorded merge commit against
   the refreshed remote-tracking branch. A failed fetch is still a cleanup
   warning that names the fetch, and it is retried on the next owner pass.
2. **A refusal caused by an amendment says so.** When every Task a
   carry-forward refuses was refused for moved inputs, and each moved input was
   changed only by non-Task commits on the target after the Run started, the
   refusal names those amending commits. Its next action lists the commands in
   order: save the branch, reset it to the commit before the first amendment,
   carry forward, cherry-pick the amendments, and retry. Roundfix prints these
   commands and never runs them.
3. **An owner older than its items' main is named.** When an item starts, and
   the owner binary's build commit is an ancestor of the item's starting main
   that changed Roundfix source (`cmd/`, `internal/`, `go.mod`, `go.sum`) after
   it, the item records the warning
   `owner-older-than-main: owner build <commit> predates starting main <commit>`.
   `deliver status` prints it. The queue never stops for it, and the check
   applies only when the build commit exists in the repository, which is the
   Roundfix repository itself.
4. **A Task commit is authorized by the grant it ran under.** The QA
   authorization audit reads the grant at the fork point, as today. When that
   grant does not cover the commit, it reads the record in the commit's parent,
   and uses it if the delivery target's history holds a byte-identical record
   at the same path. A path outside every such grant still refuses, and a
   self-approving commit is still refused (ADR-0178).

## Non-Goals / Out of Scope

- Running Git recovery commands for the operator, or changing what a
  carry-forward proves.
- Stopping, parking or refusing a queue because of an older owner binary, or
  rebuilding it.
- Accepting a grant that exists only on the item branch, or one that differs
  from every version on the default branch.
- A Run Database schema change, a new command or flag, or a change to the
  retry limit.

## Success Metrics

1. In a repository whose clone has a delivery remote, an item merged on the
   remote but not yet fetched locally has its Runs and worktree released with
   no warning. A repository without that remote keeps today's behavior.
2. A carry-forward refused only because a later non-Task commit amended the
   Spec's `_prd.md` names that commit and prints the five ordered commands.
   A refusal caused by a Task's own input change prints today's next action
   unchanged.
3. An item whose starting main changed `internal/` after the owner's build
   commit records `owner-older-than-main` and still runs. An owner built from
   the starting main, or a repository where the build commit is absent,
   records nothing.
4. A Task commit whose parent carries a grant byte-identical to one on the
   delivery target is authorized for the paths that grant bounds, and its
   `Revision` names the delivery-target commit that holds it. The same commit
   with a parent grant that exists only on the item branch still refuses.

## Recorded limits

- The owner warning needs the build commit in the repository. A released
  binary without a build commit never warns, as Auditor Staleness already
  reports an unknown ancestry.
- The amendment hint names amending commits only when every refusal is a
  moved-input refusal explained by them. A mixed refusal keeps today's text.

## Decisions

- **Refresh, then read.** The cleanup reads the refreshed remote-tracking
  branch, the same fact the item worktree was created from. A local
  `refs/heads/<default>` is used only when the repository has no delivery
  remote.
- **Print the recovery, never run it.** The recovery resets the item branch, a
  destructive Git operation that stays with the operator.
- **Warn on source age, not version.** The owner's build commit and Git
  ancestry are what Auditor Staleness already trusts. Diffing only Roundfix
  source keeps a docs-only main from warning.
- **The grant the Task ran under.** See ADR-0178.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in
the Task Graph. The outside-evidence rows rest on sources this Spec did not
produce:

- This session's recorded queue evidence, quoted in the adopted entries:
  `deliver status` on 2026-09-29, 19:26–22:59, showing the cleanup warning for
  0183. The 0181 QA reports on Run branches
  `roundfix/run-run_20260929T200238Z_ca4a5877c9dc9d5c` and
  `roundfix/run-run_20260929T201521Z_3d1a54f26096ba81` record revision
  `6784210b` for task_07.
- An independent project hit the same class: worktrunk issue #3519,
  "wt merge trusts stale local default-branch ref"
  (<https://github.com/max-sixty/worktrunk/issues/3519>). Its maintainer
  traced the defect to resolving a target against a local ref that was never
  refreshed.

## Research basis

The four adopted Backlog Entries are indexed in
[references/_index.md](references/_index.md). The Secondbrain was consulted
through `wiki/index.md` and the queries
`qmd query "delivery queue cleanup after squash merge fetch merge commit"` and
`qmd query "authorization grant widened mid delivery merge base task commit"`.
They returned only this repository's mirrors: Spec 0175's cleanup
references, Spec 0178's audit Tasks and the queue-trial finding. They add
nothing beyond the adopted entries. Exa located worktrunk issue #3519 and a
merge-verification guide stating that remote-tracking refs must be fetched
before an ancestry check is trusted. Both confirm this Spec's first premise.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
