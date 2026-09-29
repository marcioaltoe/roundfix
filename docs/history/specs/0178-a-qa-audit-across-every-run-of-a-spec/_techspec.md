---
spec: 0178-a-qa-audit-across-every-run-of-a-spec
status: active
created: 2026-09-28
surfaces: [backend, docs]
---

# A QA audit across every Run of a Spec

## Executive Summary

Resolve one Delivery Base per QA step, the merge base of the audited head and
the repository default branch. The mechanical stage reads every Task commit of
the Spec from it and reads each grant at it, and the report's authorization
audit table names each audited commit. The Daemon measures its own staleness
against the same base and, when it is `stale`, publishes and records a warning
while the gate proceeds. Settlement keeps
the Daemon's auditor fields as seeded and, in a Roundfix self-audit, requires
the QA Agent's `user_flow_binary` to be a build of the audited head.

## Project Constraints

- Identifier strategy: not applicable — no new identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git reads and files only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014, ADR-0015, ADR-0023, ADR-0057,
  ADR-0080, ADR-0089, ADR-0091, ADR-0093, ADR-0096, ADR-0097, ADR-0104,
  ADR-0117, ADR-0130, ADR-0132, ADR-0138, ADR-0155 and ADR-0156 hold; ADR-0020,
  ADR-0038, ADR-0056, ADR-0127, ADR-0141, ADR-0159 and ADR-0160 do not apply. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The Delivery Base

A new file, `internal/daemon/qa_delivery_base.go`, holds
`qaDeliveryBase(ctx context.Context, plan TaskPlan) (string, bool, error)`. It
reads the current branch of `plan.WorkDir` with
`git symbolic-ref --quiet --short HEAD` (empty when detached) and passes it to
`preflight.DetectDefaultBranch`. When the default branch is determined, it
resolves `refs/remotes/origin/<name>` and, when that ref is absent,
`refs/heads/<name>`, then returns `git merge-base <ref> HEAD` run in
`plan.WorkDir` with `true`. An undetermined default branch, a ref that resolves
under neither name, or a merge base that does not exist (exit `1`) returns
`("", false, nil)`; any other Git failure returns a wrapped error. The glossary
gains a **Delivery Base** entry.

## Every Run's Task commits

`qaMechanicalRequest` in `internal/daemon/task_engine.go` calls
`qaDeliveryBase` when `plan.HeadSHA` is set. With a resolved base it passes
that base to `speccheck.ResolveMechanicalAuthorization`, sets
`DeliveryTargetRevision` to it and reads Task commits from it. Unresolved, it
keeps `plan.HeadSHA` for all three, as today, and sets the new
`MechanicalRequest.TaskCommitsFromRunStart` to `true`.

`mechanicalTaskCommits` takes the range start as a parameter and reads
`git log --no-merges <start>..HEAD` as today, but keeps every commit of each
non-QA Task of the Spec, not only the first one per Task in log order. It still
drops a commit that intersects no governed path, and returns the kept commits
oldest first.

In `internal/speccheck/mechanical.go`, `detectMechanicalAuthPaths` records the
skip `(DetectorMechanicalAuthPaths, "Task commits of earlier Runs")` when
`TaskCommitsFromRunStart` is set, and still audits the commits it was given. The
input string is the exported constant `MechanicalSkipEarlierRunTaskCommits`.
Everything else in the detector is unchanged: the authorizing revision is still
`merge-base <delivery target> <commit>^1`, which for a commit made after the
Delivery Base is the Delivery Base itself, so each grant is read as it stood on
the default branch.

## The commit column

`MechanicalAuthorizationRead` gains `Commit string`. Every read
`detectMechanicalAuthPaths` appends — granted, refused and unresolved,
including an unavailable Task commit — sets it to the audited Task commit's
SHA. `WriteMechanicalResult` in `internal/speccheck/report.go` writes the table
as `| Task | Commit | Outcome | Record | Revision | Detail |`, with the full SHA
in the Commit cell. The table has no `Status` column, so it is never read as a
QA row.

The qa-gate skill's commit-dependent tooling audit states that the mechanical
stage audits every Task commit between the delivery base and the audited head,
that its authorization audit table names each audited commit, and that the gate
audits by command only a Task commit in that range the table does not list and
that changes a repository-tooling path, or every Task commit when the report's
mechanical skips name `Task commits of earlier Runs`.

## Staleness against the Delivery Base

`spec.ResolveAuditorEvidence` in `internal/spec/auditor_evidence.go` becomes
`ResolveAuditorEvidence(ctx context.Context, repoRoot, deliveryBase string, binary app.AuditingBinary)`.
Ancestry compares the build commit with `deliveryBase` instead of `HEAD`: equal
is `AncestryNotOlder`, a strict ancestor is `AncestryOlder`, a non-ancestor is
`AncestryNotOlder`. An empty `deliveryBase` leaves ancestry `AncestryUnknown`,
and the declared-version fallback applies as today. `AuditorEvidence` gains
`Binary app.AuditingBinary` and `DeliveryBase string`, set to the binary and the
base it was resolved for.

`CompareToTree` in `internal/app/version.go` words its two ancestry reasons
`commit ancestry: build commit predates the delivery base` and
`commit ancestry: build commit does not predate the delivery base`.

`daemon.Dependencies` gains `Auditor func() app.AuditingBinary`; a nil value
means `app.Auditor`. `writeMechanicalQAReport` keeps its signature and resolves
evidence with the engine's auditor and the resolved Delivery Base, or `""` when
unresolved. `mechanicalQAReportContent` and `spec.WritePreconditionRefusalReport`
write `evidence.Binary` when its `Version` is set, and `app.Auditor()`
otherwise, so every existing caller keeps its output.

A `stale` auditor is a warning, never a refusal. When the evidence's state is
`stale`, `mechanicalQAReportContent` appends, in both its seed and its refusal
form, a `## Auditor staleness warning` section written as a list (never a
table, so it is never read as a QA row) with two items:
`- auditor_staleness: <staleness line>` and
`- action: rebuild roundfix from delivery base <base> and restart it before the next gate`.
After `writeMechanicalQAReport`, the QA step publishes exactly one
`runevent.KindDaemonQA` event whose payload carries `phase: auditor_staleness`,
`auditor_staleness` (the staleness line), `delivery_base`, `action` (the same
instruction) and `report`. `daemon.qa` is not projected onto the Supervisor Run
Event Stream, so `roundfix events` and its public contract are unchanged; the
event reaches the Run journal and the cockpit. Nothing else changes: the
mechanical result is not marked blocking, repository Verification and the Agent
turn run as for `current` and `unknown`, and neither state publishes the event
or writes the section. No `auditing binary age` Precondition Refusal and no
`AuditorStalenessCheck` constant exist.

## The user-flow binary

`spec.QAReport` gains `UserFlowBinary`, read from the optional front matter key
`user_flow_binary`. `AuditorEvidence` gains `SelfAudit bool`, true when the
binary's build commit, with `-dirty` removed, is a commit object in the audited
repository.

After `writeMechanicalQAReport`, the QA step reads the seeded report back and
keeps its `AuditingBinary` and `AuditorStaleness`. `settleQAVerdict` receives
them with the evidence and the audited head (`git rev-parse HEAD` in
`plan.WorkDir`). For a `pass` or `partial` that `QAReportEligibility` accepts,
it refuses:

- a present `auditing_binary` or `auditor_staleness` that differs from the
  seeded value, with a cause containing `auditor fields are Daemon-owned`;
- in a self-audit, a missing `user_flow_binary`, one with no build commit, or
  one whose build commit is not a prefix of the audited head, with a cause
  containing `user_flow_binary`. The build commit is the text after the first
  `(` up to the first `,` or `)`, with `-dirty` removed; it must be at least
  seven hexadecimal characters.

A refusal settles the QA Task with the existing reason form
`QA verdict <verdict> not accepted: <cause>`. Outside a self-audit the
`user_flow_binary` check does not run. In a self-audit, the QA prompt gains the
line
`Self-audit: build roundfix from this Run Worktree with make build, run every public-CLI row with ./bin/roundfix and never a roundfix found on PATH, and record its --version line as user_flow_binary.`

The qa-gate skill's report template carries `user_flow_binary`, and the skill
says to keep the seeded `auditing_binary` and `auditor_staleness` lines,
because the auditor fields are Daemon-owned. The QA Report and Auditing Binary
entries of `CONTEXT.md` say the same.

## API Contracts

1. The Delivery Base of a QA step is `git merge-base <default branch ref> HEAD`
   in the Run Worktree, preferring `refs/remotes/origin/<default>`.
2. The mechanical stage audits every governed commit of every non-QA Task
   between the Delivery Base and the audited head; without a Delivery Base it
   records the skip `Task commits of earlier Runs`.
3. The authorization audit table's columns are
   `Task | Commit | Outcome | Record | Revision | Detail`.
4. `auditor_staleness` compares the Daemon's build commit with the Delivery
   Base, with reasons `commit ancestry: build commit predates the delivery base`
   and `commit ancestry: build commit does not predate the delivery base`.
5. A `stale` auditor publishes one `daemon.qa` event with
   `phase: auditor_staleness`, `auditor_staleness`, `delivery_base`, `action`
   and `report`, and the seeded report carries an
   `## Auditor staleness warning` list; the gate proceeds.
6. A QA Report may carry `user_flow_binary`, the `--version` line of the binary
   its public rows ran. Daemon settlement refuses a `pass` or `partial` that
   changes a seeded auditor field and, in a self-audit, one whose
   `user_flow_binary` is not a build of the audited head.

## Coverage Map

- Goal 1 → Every Run's Task commits; API Contracts 1-2.
- Goal 2 → Every Run's Task commits; API Contract 1.
- Goal 3 → The commit column; API Contract 3.
- Goal 4 → Staleness against the Delivery Base; API Contracts 4-5.
- Goal 5 → The user-flow binary; API Contract 6.
- Core Feature 1 → The Delivery Base, Every Run's Task commits.
- Core Feature 2 → The commit column.
- Core Feature 3 → Staleness against the Delivery Base.
- Core Feature 4 → The user-flow binary.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 1.
- Success Metric 3 → Testing Approach 1.
- Success Metric 4 → Testing Approach 3.
- Success Metric 5 → Testing Approach 4.
- API Contracts 1-2 → The Delivery Base, Every Run's Task commits.
- API Contract 3 → The commit column.
- API Contracts 4-5 → Staleness against the Delivery Base.
- API Contract 6 → The user-flow binary.

## Integration Points

- **Spec 0172.** Seeded QA Reports as `pending` and routes a refused `pass` or
  `partial` through `QA verdict <verdict> not accepted: <cause>`, which this
  Spec's settlement checks reuse.
- **Spec 0138.** Left the gate auditing Task commits by command until the stage
  names its commits; the commit column lets the gate stop repeating the stage.

## Testing Approach

1. **Every Run's commits.** In disposable repositories with an
   `origin/HEAD` symbolic ref: an earlier Run's out-of-grant commit is found;
   an older commit of a Task is found; a commit already on the default branch is
   not audited; a grant widened on the Spec branch does not authorize and the
   same grant landed on the default branch and merged does; an undetermined
   default branch falls back with the skip. Each case runs
   `speccheck.RunMechanicalStage` on the request `qaMechanicalRequest` built.
2. **Commit column.** Each read carries its commit, including an unavailable
   one, and the written table carries the column and one row per commit.
3. **Staleness.** `ResolveAuditorEvidence` against a base: a build of the base
   is current while HEAD carries more commits, an ancestor of the base is stale,
   no base is unknown. A QA step with an injected stale auditor publishes the
   warning event once, records the warning section and still runs repository
   Verification and the Agent; current and unknown auditors publish no warning
   and write no section.
4. **User-flow binary.** Settlement accepts the seeded fields and a
   `user_flow_binary` of the audited head, refuses a rewritten field, a missing
   or foreign `user_flow_binary`, ignores the key outside a self-audit, and the
   self-audit prompt names it; `TestArchivedPassCorpusRemainsArchiveEligible`
   stays green.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The Delivery Base and every Run's Task commits (depends on: none).
2. The commit column (depends on: 1).
3. Staleness against the Delivery Base (depends on: 2).
4. The user-flow binary (depends on: 3).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A mid-Spec grant widening now needs its own pull request.** Spec 0172
  widened its grant inside the consuming pull request; under this Spec the
  widening must land on the default branch and be merged into the Spec branch.
- **A stale Daemon keeps auditing.** Merging the default branch into a Spec
  branch, or a queued Spec starting from a newer main, moves the Delivery Base,
  so a Daemon built earlier reads `stale`. The gate still runs with it, so its
  mechanical detectors may lack rules the default branch already holds; the
  warning names the base to rebuild from, and the public rows stay covered by
  the `user_flow_binary` check.
- **Shared files.** Tasks 1, 3 and 4 edit `internal/daemon/task_engine.go` and
  `CONTEXT.md`, Tasks 1 and 2 edit `internal/speccheck/mechanical.go`, Tasks 2
  and 4 edit the qa-gate skill, and Tasks 3 and 4 edit
  `internal/spec/auditor_evidence.go` and `internal/spec/qa.go`, and every Task
  shares the implement-task instruction, so the graph is a chain. Each Task puts
  its new tests in a file of its own.
