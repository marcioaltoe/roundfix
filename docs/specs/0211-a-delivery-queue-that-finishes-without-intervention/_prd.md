---
spec: 0211-a-delivery-queue-that-finishes-without-intervention
status: active
created: 2026-10-01
surfaces: [backend, cli, docs]
---

# A delivery queue that finishes without intervention

On 2026-10-01 the Delivery Queue needed an operator four times for items it
should have finished on its own, and every time the queue's own answer did
not work. Spec 0204 parked `qa-environment-partial`; the operator followed the
printed answer, and the Delivery Retry refused because it could not find the
Implement start head of a Run the queue owner had started. Spec 0203 parked
`checks-failed`; after the operator re-ran the failed job and retried, the
queue merged while the re-run was still pending, GitHub refused, and the item
parked `delivery-error`. A retry of that park refused even though the Pull
Request had turned green, and the operator merged Pull Request #317 by hand.
Spec 0210's Implement preflight then refused because a preload named in
`NODE_OPTIONS` had been removed from the temporary directory, which stops
every Node program, `acpx` and the ACP adapters included. This Spec makes a
Delivery Retry resume each of those parks at the stage its evidence supports,
makes the checking stage wait for what GitHub requires rather than for the
checks it happens to list, and keeps a dead preload out of the agent
environment.

## Prerequisites

None. Spec 0212 (gates and Run storage that let a correct delivery finish)
also raises the Roundfix Skill's version and names this Spec in its Task
Graph's `requires`, so the queue owner delivers this Spec first.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Delivery Queue
  items keep their Spec slug, Runs their Run identifiers and Park Classes
  their names. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the checking and merging stages keep
  reading GitHub only through the user's authenticated `gh`, and this Spec
  adds one field to a read the queue already makes. No credential is read,
  stored or printed, and no test reaches GitHub. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0211 (this Spec) removes a missing
  preload from the agent environment with a notice. The archived retry still
  re-enters review, because an override waives only the archive prerequisite,
  ADR-0154: "The exception can waive the terminal QA Task's archive completion/evidence".
  ADR-0140, ADR-0161, ADR-0170, ADR-0192, ADR-0193 and ADR-0199 hold
  unchanged, ADR-0187 and ADR-0189 govern the Roundfix Skill edit, and
  ADR-0184 has the TechSpec state the changed command surface as a
  transcript, ADR-0184: "A TechSpec now declares numbered Surface Transcripts".
  The gate is bound by ADR-0080, ADR-0088, ADR-0091, ADR-0104, ADR-0155,
  ADR-0156 and ADR-0167. ADR-0097's carry conditions, ADR-0194's recorded
  observations, ADR-0195's always-observed rows and ADR-0210's per-input
  digest hold unchanged for this Spec's own gate, and ADR-0096's mechanical
  stage keeps its role. ADR-0093, ADR-0117, ADR-0168, ADR-0176 and ADR-0183
  check this Spec's consistency by citation and receipt. ADR-0182 does not
  apply, because no Task settlement fact changes. ADR-0212 does not apply: it belongs to Spec 0212 and
  governs reconcile, which this Spec leaves unchanged. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the Roundfix Skill under `.agents/skills/`
  and its `SKILL.md` mirror are Governed Paths, and the maintainer authorized
  skill edits ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0211-a-delivery-queue-that-finishes-without-intervention/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- Each printed answer to a park the Delivery Queue recorded on 2026-10-01
  brings the item back to the stage its evidence supports and on to merge,
  with no hand-made Pull Request or merge.
- The queue merges only when GitHub reports the Pull Request mergeable with
  passing status, and a pending re-run keeps it waiting.
- A preload that cannot load never stops an agent process, and the operator
  is told which preload was dropped.

## User Stories

1. As the operator of an unattended queue, I want a Delivery Retry after an
   environment-only QA partial and a QA Archive Override to resume at review,
   so that the answer the queue prints works for a Run the queue started.
2. As the operator, I want the queue to wait while a re-run check is pending,
   so that it never asks GitHub to merge a Pull Request its branch policy
   still blocks.
3. As the operator, I want a retry of a `delivery-error` park whose Pull
   Request is green to resume checking and merge, so that I never merge by
   hand.
4. As the operator, I want a refused retry to say that the retry was refused
   and why, so that I do not read a usage error into a delivery decision.
5. As a user whose shell exports a `NODE_OPTIONS` preload from a temporary
   directory, I want Roundfix to drop that preload when the file is gone and
   tell me, so that cleanup of my temporary files does not stop a Run.

## Core Features

1. **The archived retry finds a queue-started Run.** The Delivery Retry of an
   operator-archived `qa-environment-partial` item reads the Implement start
   head of the item's Run by the repository the Run belongs to, not by the
   working directory the Run ran in, so a Run the queue owner started in the
   item worktree is found. The retry then accepts an item head that descends
   from that start head and returns the item to `reviewing`.
2. **Green means GitHub says mergeable and passing.** The checking stage
   reads the Pull Request's merge state with its checks. While GitHub reports
   the merge blocked or its state unknown, the item waits, even when every
   listed check passes, until the existing checks timeout parks it
   `checks-timeout`. Only a mergeable state with passing status moves the item
   to merge.
3. **A refused merge goes back to checking.** When GitHub refuses the merge at
   an unchanged head because its base-branch policy still prohibits it, the
   item returns to `checking` instead of parking `delivery-error`, and the
   checks timeout bounds the wait. Any other merge failure parks as today.
4. **A `delivery-error` retry resumes its stage.** A Delivery Retry of a
   `delivery-error` park on an archived candidate whose head is unchanged and
   whose Pull Request is recorded returns the item to `checking`, which merges
   once GitHub reports it mergeable. A refusal still names the rule that
   refused it.
5. **A refused retry reads as a refusal.** `roundfix deliver retry` reports a
   refused retry under `Retry refused` with its reason, the item's recorded
   stage and blocker, and a statement that nothing changed, with no usage
   hint. Argument errors keep the usage hint,
   and the exit code stays `2` for both.
6. **A missing preload is dropped from the agent environment.** Roundfix
   removes, from the environment of every agent process it starts, each
   `NODE_OPTIONS` `--require`, `-r` or `--import` whose value is a file path
   that does not exist, keeps every other option, and prints one notice naming
   each removed path (ADR-0211). The user's environment is not edited.
7. **The skill and the guides say so.** The Roundfix Skill's `deliver` and
   `runtime` references and the `deliver` and `doctor` command guides describe
   the new waiting, the merge refusal, the retry resumption and its refusal
   wording, and the preload notice.

## User Experience

The operator sees fewer parks and no new command. A retry that is refused
prints `Retry refused`, the reason, the item's stage and blocker, and that
nothing changed. When a preload is dropped, the command that
started agent processes prints, once on stderr, a notice naming the
`NODE_OPTIONS` preload that does not exist and saying it was left out of the
agent environment. `roundfix deliver status` keeps its fields; an item waiting
on a blocked merge state stays in `checking`.

## Non-Goals / Out of Scope

- Retrying failed checks automatically beyond the existing single re-run of a
  check that failed outside the item's packages.
- Reading GitHub's required-check configuration or branch-protection rules.
- Interpreting any `NODE_OPTIONS` option other than a preload, editing any
  other variable, or editing the user's shell configuration.
- The pre-PR review on a large diff, the QA import, the Task Context kinds and
  reconcile: Spec 0212 owns them.
- Changing Park Class names, retry limits, the Queue Token Ceiling or the
  QA Archive Override contract.
- The evidence snapshot that broke the review of Spec 0203 on 2026-10-01: Spec
  0210 already fixes it, and this Spec does not repeat it.

## Success Metrics

1. Success Metric: a Run created with the item worktree as its working root
   and the checkout as its repository yields its start head to the archived
   retry, which returns the item to `reviewing`; before this Spec the same
   fixture refuses with `has no Implement start head`.
2. Success Metric: with every listed check passing and the merge state
   `BLOCKED`, the queue issues no merge and keeps the item in `checking` until
   the merge state is `CLEAN`, then merges exactly once.
3. Success Metric: a merge refused with `the base branch policy prohibits the
   merge` returns the item to `checking`, and the item merges on the next
   clean read with no `delivery-error` park.
4. Success Metric: a retry of a `delivery-error` park at an unchanged archived
   head with a recorded Pull Request returns the item to `checking`, and a
   refused retry prints `Retry refused` without `Usage:`.
5. Success Metric: an agent process started with `NODE_OPTIONS` holding a
   missing `--require` path and `--max-old-space-size=4096` receives only
   `--max-old-space-size=4096`, and one notice names the path; an existing
   preload and a package-name preload reach the process unchanged.

## Acceptance evidence

The outside-evidence row rests on records this Spec did not produce:

- The operator's intervention log of 2026-10-01, entries 61 to 63, and the
  Backlog Entry this Spec adopts, which record the three parks and their
  hand answers.
- GitHub's record of Pull Request #317, read on 2026-10-01: head
  `d46674d16d4a0bf2a3f280f9cd32ebd061a21b94`, workflow `CI | Verify` attempt 2
  started at 18:05:24Z and completed at 18:10:24Z, and the Pull Request merged
  by hand at 18:11:44Z.
- The Run Database, read-only on 2026-10-01: Runs the queue owner started
  record the item worktree as their working root and `/Users/marcio/dev/roundfix`
  as their repository root.
- GitHub's GraphQL reference for `MergeStateStatus`
  (<https://docs.github.com/en/graphql/reference/enums#mergestatestatus>):
  "BLOCKED: The merge is blocked." and "CLEAN: Mergeable and passing commit
  status." The `gh pr checks` manual
  (<https://cli.github.com/manual/gh_pr_checks>) puts each check into a
  `pass`, `fail`, `pending`, `skipping` or `cancel` bucket.
- Node.js's command-line reference (<https://nodejs.org/api/cli.html>):
  `NODE_OPTIONS` is "A space-separated list of command-line options", an
  option value with a space "can be escaped using double quotes", and
  `--require` and `--import` "Preload the specified module at startup."
- A reproduction during authoring on 2026-10-01: with the shell's
  `NODE_OPTIONS` naming a removed preload, a Node command-line tool died with
  `Error: Cannot find module` and `Require stack: - internal/preload`, and the
  same command succeeded with `NODE_OPTIONS` unset.

## Decisions

- A missing preload is dropped with a notice rather than refused. The
  maintainer preferred dropping on 2026-10-01, and ADR-0211:
  "Roundfix now removes, from the environment it gives an agent process, each
  `NODE_OPTIONS` preload whose value is a file path that does not exist". See
  ADR-0211.
- The checking stage trusts GitHub's merge state over the listed checks, and
  a base-branch policy refusal is waiting, not an error. A listed check set
  can omit a re-run that has not registered yet; the merge state cannot.
- The Run lookup matches by repository, as every other Run listing in the
  Run Database already does.
- The three interventions extend the adopted Backlog Entry instead of
  minting new ones, under the rule that sources sharing a context share a
  Spec.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
