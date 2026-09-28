---
spec: 0176-baseline-follow-ups-and-the-incremental-tier
status: archived
created: 2026-09-28
surfaces: [backend, cli, docs]
archived: "2026-09-28"
source_slug: 0176-baseline-follow-ups-and-the-incremental-tier
---


# Baseline follow-ups and the incremental tier

The Context-Driven Baseline still has four open gaps, left behind by Specs 0121,
0163 and 0169:

- **The incremental tier is required but never published.** The generated
  agent instructions and Spec routing guide require "the incremental
  verification command declared by the active Baseline Profile" and call the
  clause unmet when none is declared. Only the Standard TypeScript Monorepo
  Profile declares one, the guide template renders only `{{verification.gate}}`,
  and the Setup Manifest records only the complete gate. So every adopter Spec
  writes a waiver by hand: the Fiscus mirror carries the waiver phrase in 25
  files. This is Spec 0121 Core Feature 6, which Spec 0163 left out because two
  governed test files were outside its authority. The maintainer authorized
  `internal/cli/baseline_plan_test.go` and
  `internal/cli/baseline_release_gate_test.go` on 2026-09-25.
- **The mandatory docs-layout rule still says to hand-stamp an override.** The
  clause `clause.spec.keep-artifacts-in-spec-folder` in
  `internal/baseline/assets/modules/spec-workflow.json`, rendered into
  `docs/agents/docs-layout.md`, tells an agent to "stamp `qa_override: true`" and
  record the approval itself. The archive-spec skill says an override runs only
  through `roundfix archive <slug> --qa-override --approval <source> --reason
  <text>`, the command that records the user authority for a QA archive
  override.
- **A lock rewrite during the fetch can be overwritten.** `buildSkillsReconcilePlan`
  in `internal/baseline/skills_reconcile.go` and `buildSkillsRestorePlan` in
  `internal/baseline/skills_restore.go` read `skills-lock.json`, then acquire
  the source commit, then capture the lock's transaction preimage. A rewrite
  during the acquisition becomes the preimage, so the apply-time check passes,
  and a result computed from the old bytes overwrites the new edit.
- **A blocked reconcile breaks its output contract.** When a Profile-required
  skill is absent at the selected revision, `baseline skills reconcile --format
  json` prints `"plannedChanges": null` and exits `3`. The help documents exit
  `3` only for plan confirmation, so an agent looks for a Plan Digest that does
  not exist.

This Spec takes over Spec 0121 Core Feature 6, the only part of Spec 0121 not
yet delivered: Specs 0148 and 0163 shipped the rest, and Core Feature 2 was cut
on 2026-09-24. On delivery, Spec 0121 is superseded by this Spec. The operator
runs `roundfix supersede --spec 0121-baseline-decisions-and-complete-regeneration
--by 0176-baseline-follow-ups-and-the-incremental-tier --reason <text>` and
`roundfix archive 0121-baseline-decisions-and-complete-regeneration` in the
implementation pull request's archive step. No Task of this Spec moves or edits
a Spec 0121 file.

## Project Constraints

- Identifier strategy: not applicable — no new persisted identifier. The new
  Baseline decision takes the catalog identity `verification.incremental`,
  beside the existing `verification.gate`, and lock entries keep their names.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files, Git and the embedded
  catalog only. Skill-lock acquisition keeps the unauthenticated public Git
  fetch `baseline skills restore` already uses, and adds no endpoint or
  credential. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0118 turns a repository constant
  that varies into a decision, and ADR-0144 has setup ask only for a missing or
  incompatible decision identifier. ADR-0066 keeps Baseline execution in the
  CLI, and ADR-0103 has an applied refresh republish the Setup Manifest.
  ADR-0058 has Baseline upgrades fail closed on an unaccounted rule removal,
  and ADR-0099 keeps retention accounting mechanical; the reworded clauses keep
  their identities. ADR-0081 makes sanctioned digest regeneration fallout of
  the authorized edit. ADR-0085 does not gate a regeneration run on the pins it
  rewrites. ADR-0149 has the grant name the command while the tree names its
  outputs. ADR-0071 keeps Baseline Plans portable and preimage-bound. ADR-0073
  has Baseline apply use a recoverable multi-file transaction. ADR-0072 has the
  Baseline Go cutover preserve the Python contracts, including the restore
  payload shape. ADR-0154 has a QA archive override record user authority, not
  a pass. ADR-0150 cites ADR-0118 but governs work and Run branch names, which
  this Spec does not touch, so it does not apply. ADR-0130 keeps a path
  governed once it is bounded. ADR-0093 checks Spec
  consistency by citation, and ADR-0104 accepts on evidence a Spec did not
  author. ADR-0155 makes the `qa` Task declare the matrix, and ADR-0156 makes a
  declared promise name a consuming Task. This Spec's gate is bound by ADR-0080,
  ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — two grants cover this Spec. The maintainer
  approved Onda 2 of the efficiency sequence in chat on 2026-09-28, and
  expressly authorized the two Baseline CLI test files on 2026-09-25. The other
  Go files and the Baseline assets ride the standing grant of 2026-09-21 for
  governed source, and the skill files ride the standing grant of 2026-09-18.
  All of it is recorded in [_authorization.md](_authorization.md); bounded
  files: `internal/baseline/assets/decisions.json`,
  `internal/baseline/assets/modules/core.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/go-cli-tui.json`,
  `internal/baseline/assets/profiles/rust-cli.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/baseline/assets/templates/index.json`,
  `internal/baseline/assets/templates/guides/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `docs/agents/agent-instructions.md`, `docs/agents/docs-layout.md`,
  `docs/agents/spec-routing.md`, `docs/agents/setup-context.json`,
  `internal/baseline/plan_test.go`, `internal/cli/baseline_human_test.go`,
  `internal/cli/baseline_plan_test.go`,
  `internal/cli/baseline_release_gate_test.go`,
  `internal/docscontract/publicdocs_test.go`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/setup-context-driven/SKILL.md`,
  `skills/setup-context-driven/SKILL.md`. Sanctioned regeneration:
  `make baseline-digests` and `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Goals

- A repository's incremental Verification is a recorded decision that the
  generated guidance names, so no adopter Spec needs a waiver.
- An existing single-gate Setup Manifest names the missing decision and never
  passes by omission.
- The generated docs-layout rule sends a QA Archive Override through the
  command, the way the archive-spec skill does.
- A `skills-lock.json` rewrite during the source fetch is never overwritten by
  a plan computed from the old bytes.
- A blocked reconcile prints an empty plan and a documented exit.

## Core Features

1. **A blocked reconcile honors the output contract.** The required-removed
   result prints `"plannedChanges": []` and `"planDigest": null`, and keeps exit
   `3` with finding `reconcile.required-removed`. The help, the user guide and
   the Roundfix skill document that exit next to the confirmation exit.
2. **The lock is read after the fetch.** Skill-lock reconciliation and
   restoration compute their plan from a `skills-lock.json` read taken after the
   last source acquisition. The recorded preimage equals the bytes the plan was
   computed from. A capture that differs is refused with
   `lock.changed-during-plan`, and a later rewrite still fails the apply-time
   preimage check.
3. **The docs-layout rule names the override command.** The clause says an
   override is performed only through the runtime's archive override command,
   and names `roundfix archive <slug> --qa-override --approval <source> --reason
   <text>`. It says the command records the approval and stamps the front matter,
   forbids a hand edit of the stamp, and keeps the "unsupported runtime is
   reported as unsupported" sentence. The guide, the Setup Manifest and the
   formatter golden are regenerated.
4. **The incremental tier is a published decision.** A new catalog decision,
   `verification.incremental`, is required by every built-in Profile and by the
   core module. Its suggestion is `rtk make verify-incremental`, and it has no
   default, so a human must answer it. The decision renders into the agent
   instructions as "The selected incremental Verification is ...". It projects
   into the Setup Manifest like `verification.gate`, and planning is blocked when
   the repository declares no matching command. The two-tier clauses point at the
   selected commands instead of the Profile. An existing single-gate manifest
   makes `baseline update` name the new decision. `--adopt-suggested` or the
   interactive workflow answers it, and the plan still refuses when the
   repository does not declare that command. This repository's own guides and
   manifest are regenerated. The user guide, the setup-context-driven skill and
   `CONTEXT.md` describe the decision. (Spec 0121 Core Feature 6.)

## Non-Goals / Out of Scope

- Changing any Fiscus file or any other adopter repository.
- Editing the frozen parity corpus under
  `internal/baseline/testdata/parity-corpus/`, or any Source Baseline.
- A new exit code for the reconcile command, or a new Baseline decision type.
- Changing when a QA Archive Override is accepted or refused.
- Moving, editing, superseding or archiving any Spec 0121 file inside a Task;
  that is the operator's archive step.

## Success Metrics

1. A reconcile whose selected revision lacks a Profile-required skill prints
   `plannedChanges` as an empty array and `planDigest` as `null`, and exits
   `3`. Today it prints `null` for `plannedChanges`.
2. A `skills-lock.json` rewritten during the source fetch keeps its new bytes
   after a confirmed apply of a plan previewed on the old bytes, for both
   reconcile and restore. Today the old-bytes result overwrites it.
3. The regenerated `docs/agents/docs-layout.md` names the override command and
   no longer contains "preserve any supplied reason", and a second `baseline
   update` changes nothing.
4. A Setup Manifest without `verification.incremental` makes `baseline update`
   exit `3` naming that decision and write nothing. After the decision is
   answered with a declared command, the rendered agent instructions name that
   command and the manifest records it, where today no incremental command is
   published anywhere.

## Recorded limits

- The decision's suggestion, `rtk make verify-incremental`, fits repositories
  that use this project's `rtk` and `make` conventions. Any other repository
  answers the decision itself. The declaration check refuses a suggestion that
  the repository does not declare.
- The `lock.changed-during-plan` refusal covers only the short window between
  the post-fetch read and the preimage capture. A rewrite after planning is
  caught by the existing apply-time preimage check, as before.

## Decisions

- **A decision, not a Profile expectation.** ADR-0118 makes a repository
  constant that varies a decision. The incremental command varies per
  repository, and the Go and Rust Profiles cannot name one for every adopter.
  The decision owns the `verification.incremental` projection. The Standard
  TypeScript Monorepo Profile's portable `incremental` expectation still
  projects when the decision is absent, and is not projected twice when the
  decision is present.
- **Suggest, never default.** The decision has a suggestion and no default, so
  the interactive workflow asks and `--adopt-suggested` is an explicit choice.
  The declaration check means a suggested command the repository does not
  declare blocks the plan, so absence never passes as compliance.
- **Read once, after the fetch.** The lock read that feeds the plan and the
  preimage happens after acquisition, so the preimage is the bytes the plan
  used. An earlier validation read may stay for fail-fast errors, but it never
  feeds the plan.
- **Keep exit 3 and document it.** A required-removed block is an action the
  maintainer must take, like an unconfirmed plan, and the finding code tells the
  two apart. A new exit code would break callers that already branch on `3`.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a blocked reconcile that still
prints `null`, a fetch-time rewrite that is overwritten, a guide that still
asks for a hand stamp, or a single-gate manifest that updates without naming the
new decision would each pass a happy-path test.

The outside-evidence row rests on a Setup Manifest this Spec did not write: the
Fiscus mirror at `/Users/marcio/dev/secondbrain/projects/fiscus/mirror/`, whose
`docs/agents/setup-context.json` records only `verification.gate` and whose
Specs carry the waiver. It is replayed in a disposable repository with the
built binary, and Fiscus is never changed. When the mirror is absent, or its
manifest cannot load under the current catalog for a reason this Spec does not
own, the row is recorded as blocked with that reason.

## Research basis

The waiver's cost was captured in the secondbrain inbox
(`inbox/roundfix/_triaged/2026-09-17-a-dispensa-do-comando-incremental-virou-paragrafo-em-dez-specs.md`),
and the original gap in
`inbox/roundfix/_triaged/2026-09-08-profile-exige-verificacao-incremental-que-nao-declara.md`.
The blast radius was measured on 2026-09-28 in a disposable copy of main
`0160f70a`. Adding the decision, the projection and the template token, then
running `make baseline-digests`, broke 163 top-level tests in 16 test files.
After decision lists and Makefile fixtures were updated, the failures that
remained were in `baseline_human_test.go` (prompt scripts),
`baseline_release_gate_test.go`, `baseline_update_test.go`,
`profile_alignment_test.go` and `publicdocs_test.go` (two decision-count
assertions). The three other items are recorded limits of the pre-PR reviews of
Specs 0163 and 0169. The adopted sources are indexed in
[references/_index.md](references/_index.md).

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
