---
status: accepted
created_at: 2026-10-23T12:00:00Z
updated_at: 2026-10-23T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A Python profile takes the skills the Go and Rust setups share

Roundfix shipped no Baseline Profile for a Python repository, so Quaestor, a
Python 3.12 command-line tool run with uv, pytest and ruff, adopted the
Baseline on 2026-10-19 through a repository-owned profile that only removes the
`rust` module from a built-in one and carries no Python guidance. On
2026-10-22 the maintainer chose a full Python profile and allowed it to request
new recommended skills from the upstream skills repository. A built-in profile
names one Setup Snapshot, a Setup Snapshot mirrors one upstream setup list
(ADR-0191) or the union of several (ADR-0204), and the upstream repository has
no Python list and no Python skill: its lists are `go`, `rust`, `typescript`,
`devops` and `content`.

## Decision

**A Common Setup Snapshot.** A Setup Snapshot may name two or more component
snapshots and hold, in the first component's order, only the skills and
activation bundles that every component carries with equal entries. Like a
Composed Setup Snapshot it holds no list of its own choosing: the asset sync
writes it after it refreshes its components, and the catalog refuses one that
differs from its components' common part, names fewer than two components, or
names an unknown or derived component. The first one is `python`, the common
part of `go` and `rust`: the thirty language-neutral skills both lists carry,
the same set Quaestor installed by hand.

**The `python-cli` profile.** A built-in profile selects the universal,
context, Spec, autonomous-work, external-triage, Secondbrain and
repository-extension modules, the CLI surface and a new `python` module, on the
`python` setup. The module's guide, `docs/agents/python.md`, states with force
that uv owns dependencies and the project environment, that ruff's findings
fail the gate and are never hidden, that a Task Verification selects a new test
by its pytest node id alone, that tests stay hermetic and data from outside the
repository skips by name, that every operating-system branch is exercised, and
that a console script stays thin and is tested through its entry point. The
module requires and dispatches no skill until the upstream repository
publishes Python skills; the curation request for them is filed outside this
repository.

## Considered Options

- Composing `go` and `rust` was rejected: the union would restore every Go,
  TUI and Rust skill into a Python repository.
- A hand-written `python` snapshot was rejected: ADR-0191 makes the asset sync
  the only writer of a snapshot, and a hand-written list would drift from
  upstream silently.
- Taking the `go` setup alone was rejected for the same reason as the union:
  eleven Go and TUI skills a Python repository cannot use.
- Waiting for an upstream Python list was rejected by the maintainer's
  decision to ship the profile in cycle AP.

## Consequences

- When the upstream repository publishes a `python` list, a later Spec turns
  the `python` snapshot into a mirror of it by name, as ADR-0191 requires, and
  the module requires and dispatches the new skills; no adopter decision
  changes.
- A skill both `go` and `rust` add later joins the `python` snapshot at the
  next asset sync; a skill only one of them adds does not.
- The profile itself writes no Makefile, Project Config or CI workflow; its
  guidance never grants the authority to create repository tooling, which the
  core tooling clause reserves to the maintainer.
