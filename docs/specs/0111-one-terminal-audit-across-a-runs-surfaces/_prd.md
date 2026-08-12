---
spec: 0111-one-terminal-audit-across-a-runs-surfaces
status: active
created: 2026-08-12
surfaces: [backend, cli]
---

# One terminal audit across a Run's surfaces

Two standing rollups converge on the same missing thing. A Run's lifecycle is
spread across process trees, Run Branches, Task and Run Worktrees, refs,
artifacts, notifications, and database storage, and no single audit answers what
became of everything one Run created — failures recur when one surface decides
terminal state without disposing or classifying what another surface made. The
same shape appears in execution environments: configuration validates, the runtime
launches, and the Task's filesystem access is decided in a third place, with
nothing proving the three describe one environment. This Spec is the audit that
closes both, and it is deliberately last: four earlier Specs each settle one
surface, and building the audit before them would model seams that are about to
move.

## Project Constraints

- Identifier strategy: applicable — Run, Run Branch, Run Worktree, Task Worktree,
  Agent Session, and Force Stop are glossary terms the audit reports against, and
  a terminal disposition vocabulary spanning them is what the audit adds. The
  closing node checks whether the work introduced or changed a term. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The audit reads local process state, Git refs, the
  filesystem, and the Run Database. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0014 gives the Daemon ownership of task
  verification and status settlement, which bounds the audit to reporting
  disposition rather than settling Tasks; ADR-0099 makes retention accounting
  mechanical while classification is not, which is the same split the audit must
  honour — it computes what exists and does not decide what a human should keep.
  The decisions extending the Daemon's ownership are accounted and none applies:
  ADR-0020 ranks a parsed prompt result above the runtime's exit code, ADR-0038
  allows one Verification repair, ADR-0057 gives the Daemon exclusive ownership of
  Implement Task status, ADR-0056 separates Task Capacity from Verification
  Capacity, and ADR-0096 with ADR-0117 place the gate's mechanical stage and its
  checks — the audit reports and settles nothing, so it changes none of them.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go across the run lifecycle, agent, and
  CLI packages plus their tests. Source: `docs/agents/agent-instructions.md`.

## Goals

1. Every resource a Run created has one reported terminal disposition.
2. The audit spans repositories, not only the current checkout.
3. Accepted configuration, runtime launch, and Task filesystem access are proven
   to describe one environment before a Run does work under them.
4. Evidence needed to check a cleanup survives the cleanup.

## User Stories

1. As an operator after a night of Runs, I want one report of what each Run
   created and what became of it, so that I do not check process tables, branch
   lists, worktrees, and storage separately and still miss one.
2. As an operator with Runs across several repositories, I want the audit to cover
   them, so that a machine-wide question does not need one command per checkout.
3. As a Supervisor configuring an execution policy, I want the accepted
   configuration proven against what the runtime will actually grant a Task, so
   that a policy the runtime cannot honour is a preflight defect rather than a
   Task failure.
4. As an operator acting on the audit, I want the evidence behind each disposition
   preserved, so that I can check a cleanup rather than trust it.

## Core Features

1. **One terminal audit.** A single report answers, per Run, the disposition of
   every resource it created across process tree, branches, worktrees, refs,
   artifacts, notifications, and storage.
2. **Machine-wide scope.** The audit reports across every repository the Run
   Database knows, with per-repository grouping, rather than only the current
   checkout.
3. **An execution-policy proof.** Accepted configuration, runtime launch, and the
   Task's filesystem access are proven to describe one environment, so a
   documented escape hatch or access level the runtime cannot honour refuses at
   preflight.
4. **Evidence that outlives the disposition.** Each reported disposition carries
   what it was decided from, and that evidence survives the cleanup it justifies.
5. **Reporting, not settling.** The audit reports and classifies; it does not
   settle Tasks, delete work, or make a verdict more permissive.

## User Experience

The audit prints one section per repository and one row per Run, each naming
resources and their dispositions, with anything unresolved grouped first because
that is what needs action. With nothing unresolved it says so rather than printing
an empty inventory.

## Non-Goals / Out of Scope

- Performing cleanup. Existing commands own reclamation; this audit answers what
  exists and what it is.
- Changing Run termination, garbage collection, or retention policy.
- Re-deciding the classifications the earlier Specs in this set establish; the
  audit consumes them.
- Changing verdict semantics or settlement.
- Any cross-machine or remote inventory.

## Success Metrics

- For a completed Run, every resource it created appears once with a disposition,
  proven against a Run whose resources span at least four of the named surfaces.
- The audit reports Runs from more than one repository in a single invocation.
- An execution policy the runtime cannot honour refuses at preflight rather than
  failing every Task, measured against the documented access level that validated
  and then failed every Task in a repository this Spec did not build.
- No disposition is reported without the evidence it was decided from.

## Decisions

- The audit reports and classifies rather than acting, so that a wrong
  classification costs a reading rather than lost work.
- This Spec is sequenced last in its set, because four earlier Specs each settle
  one of the surfaces it reports on.

## Open Questions

- Whether the audit is a new command or an extension of the existing
  reconciliation surface. The default until answered is an extension, since
  reconciliation already walks several of these surfaces and a second command
  reporting on the same ones would be a second answer to one question.
- Whether the execution-policy proof belongs in this Spec or with the preflight
  Spec that already proves selections. The default is here, because it is the
  cross-surface consistency claim rather than a per-selection proof.
