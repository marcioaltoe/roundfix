---
spec: 0149-a-supported-way-to-reopen-a-settled-gate
status: active
created: 2026-09-19
surfaces: [backend, cli, docs]
---

# A supported way to reopen a settled gate

## Executive Summary

A new `reopen` command returns a Spec's terminal QA Task to `pending` when, and
only when, the loader's own stale-gate condition already holds. It reuses the
status writer `settle` uses, keeps the QA Report untouched, and records the
invalidation in the QA Task beside the Result it invalidated.

## Project Constraints

- Identifier strategy: applicable — the command addresses a Spec by slug and the
  QA Task by its Task Graph id, and mints no identity of its own. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads and writes only.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — the terminal QA gate contract and Task
  status ownership both hold; this Spec adds an operation under them. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — the shipped Roundfix skill documents the CLI
  surface. Express maintainer authorization: granted 2026-09-18 as a standing
  grant, consumed in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## The condition

`spec.Load` already computes exactly the predicate this command needs. When the
QA Task's status is `StatusCompleted` and any Task in its dependency closure is
not completed, it returns `StaleGateError{QATaskID, TaskIDs}`, wrapped as
`validate qa gate: %w`.

That error is the authority for the reopen. The command does not re-derive
staleness from its own reading of the graph — it asks the loader, and acts only
on the loader's answer. A second implementation of the predicate is a second
thing that can disagree with the gate.

The consequence is a load seam: `Load` returns `(nil, err)` today, so a caller
that recognizes `StaleGateError` still cannot see the graph it describes. The
command needs the QA Task's file path, which the error does not carry. The seam
resolves that without weakening `Load` for anyone else — every existing caller
keeps refusing on a stale gate.

## The command

```
roundfix reopen --spec <slug>
```

- **Exit 0.** The QA Task's status is written `pending`, and an invalidation
  record is appended to the QA Task file.
- **Exit 2.** Preflight Validation failed: the Spec has no terminal QA Task, the
  QA Task is not `completed`, or every dependency is completed so the verdict
  still describes the graph. The reason names which condition failed. No file is
  written.

It creates no Run, writes no Run Event Journal entry, commits nothing and pushes
nothing.

## What is written, and what is not

- **The QA Task's status** becomes `pending`, through `spec.SetStatus`, which
  rewrites only the status value and is already covered by its own tests.
- **The QA Task's body** gains an invalidation record naming the QA Report that
  was invalidated, the dependency Task ids that made the verdict stale, and the
  date. The prior `## Result` is left in place above it.
- **The QA Report** is not opened for writing. It is the evidence of what was
  verified; a reopen that rewrites it destroys the record it exists to keep.
- **No other Task file** is touched. The Tasks below the gate are in whatever
  state made the gate stale, and that state is the maintainer's, not this
  command's.

## API Contracts

1. `reopen` accepts `--spec <slug>` and refuses unknown flags, matching the
   surrounding commands.
2. A successful reopen leaves the Spec loadable: the commands that refused with
   the stale-gate reason load it afterwards.
3. A refused reopen leaves the Spec directory byte-identical.

## Coverage Map

- Goal 1 → The condition; The command.
- Goal 2 → The command (refusals); API Contract 3.
- Goal 3 → What is written, and what is not.
- Core Feature 1 → The condition; The command; API Contract 2.
- Core Feature 2 → The command (refusals); API Contract 3.
- Core Feature 3 → What is written, and what is not.
- Core Feature 4 → The command (envelope).
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 3.
- Success Metric 3 → Testing Approach 4.
- API Contracts 1-3 → The command, What is written, and what is not.

## Integration Points

- **`spec.Load` and `StaleGateError`.** The predicate and its message stay where
  they are. This Spec adds the seam that lets one caller act on the error.
- **`spec.SetStatus`.** Reused as the status writer; not modified.
- **`settle`.** The envelope this command copies — one Task's status, no Run, no
  journal entry, no push.
- **Spec 0129.** The amendment record for premise falsification and supersession,
  the coverage extension and the temporal-prerequisite representation stay
  there.

## Testing Approach

1. **The stale case, at the command seam.** A fixture Spec whose QA Task is
   completed above a pending dependency: the command exits 0, the QA Task status
   is `pending`, and `spec.Load` on the same fixture afterwards returns no
   error. Fails on the tree as it stands, where the command does not exist.
2. **The not-settled refusal.** A fixture whose QA Task is `pending`: the command
   refuses and names that condition.
3. **The not-stale refusal.** A fixture whose QA Task is completed above
   completed dependencies: the command refuses, and a directory-wide digest
   taken before and after is unchanged.
4. **Evidence retention.** After a successful reopen, the QA Report file's bytes
   are unchanged and the QA Task file still contains its prior `## Result` text,
   with the invalidation record present.
5. **The skill is true.** The shipped skill and its mirror name the command and
   its refusals, and the mirror matches the canonical file.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The load seam that exposes the QA Task path alongside a stale-gate error,
   with unit tests (depends on: none).
2. The `reopen` command, its refusals and its invalidation record, with tests at
   the command seam (depends on: 1).
3. The shipped skill and the user guide (depends on: 2).
4. Terminal QA (depends on: 1, 2, 3).

## Risks & Considerations

- **A reopen that becomes a discard.** If the command ever accepts a healthy
  gate, it stops being a recovery path and becomes a supported way to throw a
  verdict away. The not-stale refusal is the control, and it is a Success
  Metric, not a nicety.
- **A load seam that weakens the loader.** Exposing a graph alongside an error
  invites a caller to use it. The seam is named for recovery and every existing
  caller keeps its current behavior; the tests assert that `Load` still refuses.
- **Records that drift from evidence.** The invalidation record names a report
  by path. If the report is later renamed the record points nowhere, which is
  why the report is never rewritten and the record is dated.
