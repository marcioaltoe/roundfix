---
spec: 0163-baseline-decisions-and-regeneration
status: active
created: 2026-09-24
surfaces: [backend, cli, docs]
---

# Baseline decisions and regeneration

Four defects in the Baseline and its tooling remain from portfolio Spec 0121.

**Changing the HTTP mode discards the decision.** Answering "Change" on the
`http.contract` prompt returns `{"mode": ...}` from `promptBaselineDecision` in
`internal/cli/baseline_human.go`. Every recorded exception and the source are
dropped without warning. A repository observed on 2026-08-07 lost four
exceptions and its source this way.

**Greenfield asks for a step it never offers.** With stale managed carriers,
`planRootPreservationWithCatalog` in `internal/baseline/preservation.go` passes
the Greenfield early return and ends in a classification skeleton. But
`promptBaselineClassification` returns before classifying in every mode except
Preservation. So the plan fails and names a review the maintainer was never
offered.

**`make skills-sync` owns no declared outputs.** `baseline.OutputsFor` reads
ownership records only below `internal/baseline`, so a grant that names the
command covers none of the mirrors it writes. Every grant lists each mirror.

**An obsolete skill-lock entry has no safe removal.** `baseline skills restore`
restores entries but cannot retire one that upstream removed. The only way to
remove one is a hand edit to `skills-lock.json`, based on evidence nobody
recorded.

This Spec carries Spec 0121 Core Features 1, 3, 4, 5 and 7. Core Feature 2 was
cut in the 2026-09-24 restructure. Core Feature 6 stays with Spec 0121 for the
reason given under Non-Goals.

## Project Constraints

- Identifier strategy: not applicable — no new persisted entity; Baseline
  decision identities, lock entry names and ownership record paths keep their
  existing identities. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the CLI edits a repository-owned HTTP
  Contract, keeping its mode, typed exceptions and source; lock reconciliation
  fetches a public Git revision with the same unauthenticated acquisition
  `baseline skills restore` uses. No endpoint, credential or authentication
  policy is added. Source: `docs/agents/cli.md` and
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0081 makes sanctioned digest
  regeneration fallout of the authorized edit, ADR-0149 has the grant name the
  regeneration command while the tree names its outputs, ADR-0130 keeps a path
  governed once an authorization bounds it, ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0155 makes
  the `qa` Task declare the matrix and ADR-0156 makes a declared promise name a
  consuming Task. This Spec's gate is bound by ADR-0080, ADR-0091, ADR-0096,
  ADR-0097 and ADR-0117. All hold. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization granted
  2026-09-24, recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/cli/baseline_human_test.go`, `internal/baseline/plan_test.go`,
  `internal/baseline/derived_ownership_test.go`, `skills/_ownership.yml`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`. Sanctioned regeneration:
  `make baseline-digests` and `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Goals

- Changing only the HTTP mode retains every exception and the source.
- Greenfield either completes or refuses before asking for a step it cannot
  offer, and names Preservation as the route.
- A grant that names `make skills-sync` covers exactly the mirrors that command
  writes.
- An obsolete skill-lock entry is removed only on positive evidence, through a
  reviewed plan.

## Core Features

1. **A mode change keeps the decision.** Changing the HTTP mode replaces only
   `mode`. Exceptions and source stay as they were. A deliberate change to the
   exceptions still goes through an explicit decision value, which replaces the
   whole typed decision. The review shows which exceptions are kept. (Spec 0121
   Core Feature 1.)
2. **Greenfield refuses early.** When Greenfield meets retained managed source
   that needs classification, planning refuses before any classification or
   approval prompt. The next action names `preservation.mode=preservation`.
   Human and JSON planning refuse identically. No source is discarded and no
   classification is invented. (Spec 0121 Core Feature 3.)
3. **Skill regeneration declares its outputs.** `skills/_ownership.yml` assigns
   the owned skill mirrors to `make skills-sync`, and `OutputsFor` resolves that
   command's outputs from it. Authorial `.agents/skills`, Go files, `testdata`,
   `recommended.txt` and unowned skills are never among those outputs. (Spec
   0121 Core Feature 4.)
4. **Lock reconciliation on proven absence.** `roundfix baseline skills
   reconcile` classifies each `skills-lock.json` entry of one source repository
   at an explicitly selected immutable commit. It removes only an entry that is
   absent at that commit and not required by the Profile. The change runs
   through the existing preview, plan digest, preimage check and transaction.
   A required skill that is absent blocks. Removal is never authorized by an
   unreachable source, a missing local tree or a skill that is merely unneeded.
   Installed trees and unknown lock fields are kept. Doctor stays offline and
   read-only. (Spec 0121 Core Feature 7.)
5. **Shipped guidance follows.** The Roundfix and setup-context-driven skills
   describe the changed behaviour. Their mirrors are regenerated with
   `make skills-sync`. (Spec 0121 Core Feature 5.)

## Non-Goals / Out of Scope

- Spec 0121 Core Feature 2 (the unused HTTP Profile default), cut on 2026-09-24.
- Spec 0121 Core Feature 6, the publication of `verification.incremental` into
  the guide template, generated guides and Setup Manifest migration. Measured on
  2026-09-24: adding the decision to the Profiles' `entryDecisions` and running
  `make baseline-digests` breaks 84 tests in 11 files. Two of those files,
  `internal/cli/baseline_plan_test.go` and
  `internal/cli/baseline_release_gate_test.go`, are governed and outside this
  grant. The work stays with Spec 0121 until its authority covers them.
- Changing the suite guard's resolver in `internal/suiteguardcontract`, which is
  governed and outside this grant. The guard keeps exempting only Baseline
  outputs, a subset of what the audit accepts.
- Removing an installed skill tree, or changing any real repository's lock.

## Success Metrics

1. Changing only the mode of a decision with two exceptions and a source keeps
   both exceptions and the source, where today both are lost.
2. Greenfield planning with a stale managed carrier refuses before
   classification, and names Preservation, in both human and JSON output.
3. `OutputsFor(root, "make skills-sync")` returns every file of every owned
   mirror and nothing else, where today it returns none.
4. Reconciliation removes an entry absent at the selected commit and keeps
   present, moved, unrelated and unreachable-source entries.

## Decisions

- **Carry the untouched fields forward.** Of the three designs the 2026-08-07
  Finding named, a mode change that retains exceptions is the only one that
  cannot destroy data silently. Full exception editing stays explicit.
- **Refuse, don't reroute.** Greenfield does not quietly switch to
  Preservation; it refuses with the route named, so the maintainer chooses.
- **Absence is a fact at a commit.** A missing path at an immutable commit is
  evidence. A failed fetch is not. A skill found under another path is moved,
  not removed.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author. The
multi-exception HTTP decision recorded in the archived Finding
[2026-08-07-changing-the-http-contract-discards-its-exceptions.md](../../history/findings/2026-08-07-changing-the-http-contract-discards-its-exceptions.md)
is replayed with only its mode changed. The replay keeps that Finding as its
provenance and shows every exception and the source retained. Each Core Feature
needs positive and negative evidence in the Task Graph. The negative cases carry
the weight: a lock entry removed on a failed fetch, or a Greenfield refusal that
still prompts, would pass every happy path.

## Research basis

The HTTP and Greenfield defects are recorded in the archived Findings
[2026-08-07-changing-the-http-contract-discards-its-exceptions.md](../../history/findings/2026-08-07-changing-the-http-contract-discards-its-exceptions.md)
and
[2026-08-07-greenfield-adoption-cannot-satisfy-its-own-gate.md](../../history/findings/2026-08-07-greenfield-adoption-cannot-satisfy-its-own-gate.md).
The ownership gap is recorded in Spec 0121's
[reference](../0121-baseline-decisions-and-complete-regeneration/references/2026-09-08-skill-regeneration-declares-its-owned-outputs.md).
Each was re-read against the current source on 2026-09-24, and each is still
present. The measurement behind the Core Feature 6 exclusion was made on
2026-09-24 in a disposable worktree and then reverted.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
