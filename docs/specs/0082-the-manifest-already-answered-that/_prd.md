---
spec: 0082-the-manifest-already-answered-that
status: active
created: 2026-08-07
surfaces: [cli, docs]
---

# Baseline update

A repository that has already adopted the Context-Driven Baseline records every
answer it gave in its Setup Manifest — the profile, the ten decisions, the
managed artifacts and their digests, the verification gate. Refreshing that
repository against a newer Roundfix binary should be a mechanical act. Today it
is an interview: the interactive command reads the manifest, announces `update`
mode, and then asks all twelve questions anyway, each one offering the stored
value as a default that still has to be confirmed. It then spends a supervised
ACP turn re-segmenting and re-classifying the entire root instruction corpus
and asks for a rule-by-rule review of the result. The non-interactive path is
worse: it never opens the manifest at all, refuses with a list of the same
decisions the manifest holds, and — once those are supplied — refuses again for
an instruction-preservation Decision Document that, in Preservation mode, must
bind every Source Baseline Entry.

The consequence is a fleet that drifts because keeping it current is expensive.
Across nine adopted repositories there are six distinct catalog digests, and
only four of the nine sit on the one the current binary generates. This feature
makes a Baseline refresh a single command that answers itself from the manifest,
touches only what the Baseline owns, and carries the repository's skills along
with its guidance.

## Project Constraints

- Identifier strategy: not applicable — the `go-cli-tui` Baseline Profile does
  not select `identifier.strategy`, and this feature persists no record that
  needs a generated identifier; the Setup Manifest, Plan, and result documents
  are all identified by content digest. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the repository has no
  `docs/agents/backend.md` because the `go-cli-tui` Baseline Profile selects no
  backend module, and this feature exposes no HTTP surface and performs no
  authenticated request. External skill acquisition reuses the existing
  immutable Git source contract rather than introducing a new authenticated
  protocol. Source: `docs/agents/spec-routing.md`, which routes active
  obligation discovery in the absence of a backend guide.
- Active ADR obligations: applicable — ADR-0058 requires every Baseline
  transition to account for each prior managed Normative Clause and to block
  while any clause is unaccounted; ADR-0068 requires one confirmation-gated
  workflow bound to an approved Plan Digest; ADR-0070 limits automatic mutation
  to root carriers and keeps nested-carrier conflicts as warnings; ADR-0071
  requires plans to stay portable and preimage-bound; ADR-0073 requires apply to
  run as a recoverable transaction; ADR-0081 makes sanctioned digest
  regeneration fallout of the authorized edit; and ADR-0099, minted by this
  Spec, separates mechanical retention accounting from supervised instruction
  classification. Source: `docs/agents/spec-routing.md`.
- Tooling authority: applicable — this feature edits Roundfix-owned Skills and
  Baseline module assets, which `docs/agents/agent-instructions.md` places
  behind express maintainer authorization. Express maintainer authorization:
  recorded 2026-08-07 in
  `docs/workflow/authorizations/2026-08-07-baseline-update-command.md`, granted
  as "Skills + module assets, limitado ao propósito". Bounded files:
  `skills/setup-context-driven/**` and its `.agents/skills/setup-context-driven/**`
  mirror; `skills/roundfix/**` and its `.agents/skills/roundfix/**` mirror; and
  `internal/baseline/assets/modules/*.json` limited to teaching the update
  command to generated guidance. Deterministic pins rewritten by
  `make baseline-digests` are sanctioned fallout under ADR-0081 and need no
  separate grant. Source: `docs/agents/agent-instructions.md`.

## Goals

- An adopted repository refreshes its Baseline without answering any question
  it has already answered.
- A refresh leaves every repository-authored byte identical, including
  Repository-Specific Normative Rules and the prose surrounding managed blocks.
- A refresh carries the repository's skills with its guidance, so guidance and
  the skills it dispatches to cannot drift apart.
- A maintainer sweeps the whole fleet from a script, and any repository that
  genuinely needs a human decision says exactly which decision and stops
  without writing.
- Unaccounted removal of a managed Normative Clause still blocks the refresh.

## User Stories

1. As a maintainer sweeping the fleet after a Roundfix release, I want one
   non-interactive command per repository, so that bringing eight repositories
   current is a loop instead of eight interviews.
2. As a maintainer running the interactive command on an adopted repository, I
   want it to ask only what it does not already know, so that confirming twelve
   settled defaults stops being the price of a refresh.
3. As a maintainer whose repository carries hand-written rules, I want the
   refresh to leave those rules untouched and prove it, so that I can run it
   without reading a rule-by-rule diff.
4. As a maintainer, I want the refresh to update the repository's skills in the
   same act as its guidance, so that a generated guide never dispatches to a
   skill version the repository does not have.
5. As a maintainer whose repository is on an older catalog that introduced a new
   decision, I want the command to name that decision and write nothing, so that
   no policy enters my repository without my answer.
6. As a maintainer, I want a refresh that would drop a managed Normative Clause
   without accounting for it to fail rather than apply, so that a Baseline
   upgrade cannot quietly weaken governance.

## Core Features

1. **Manifest-driven refresh.** A new update path derives its complete input
   from the repository's Setup Manifest: the Baseline Profile, every recorded
   decision, and the managed-artifact ledger. It prompts for nothing that the
   manifest answers. A repository with no manifest, or with one the current
   Baseline cannot read, is not an update: the command refuses and directs the
   maintainer to first adoption.

2. **Managed-only refresh scope.** The refresh regenerates exactly the artifacts
   the Baseline owns — declared root blocks, setup-owned guides, and the Setup
   Manifest itself. Every other byte in every touched file is preserved
   unchanged, including Repository-Specific Normative Rules, repository-rule
   blocks, and authored prose outside managed markers. The refresh never invokes
   the supervised semantic analyzer and never proposes a new destination for
   repository-authored content.

3. **Retention accounting without classification.** The refresh still maps every
   prior managed Normative Clause to a current clause, to Repository-Specific
   Normative Rules, to a recognized typed repository document, or to an explicit
   recorded rejection. Accounting is a mechanical comparison of generated
   artifacts and requires no analyzer. An unaccounted clause blocks both preview
   and apply, and the presented plan carries the accounting.

4. **Confirmation-gated apply.** The refresh produces the same portable,
   preimage-bound plan and Plan Digest as adoption, and mutates only after the
   exact current digest is confirmed. A non-interactive sweep supplies that
   confirmation up front; without it the command presents the plan and writes
   nothing. Apply remains one recoverable transaction.

5. **Skills refresh in the same act.** The refresh reinstalls the Roundfix-owned
   skill bundle carried in the binary and restores external Repository Skill Set
   members that are missing or drifted. When the immutable upstream source for
   an external skill cannot be reached, the refresh reports that skill's drift
   as a warning and completes; guidance refresh is never held hostage to network
   reachability. The result names every skill installed, every skill restored,
   and every skill left drifted with the reason.

6. **New decisions stop the sweep by default.** When the current catalog
   requires a decision the manifest does not record, the command exits without
   writing and names exactly which decisions are new. An explicit opt-in flag
   adopts the catalog's suggested value for those decisions instead, and the
   result lists every value adopted that way so it can be reviewed before merge.

7. **Interactive command asks only what is new.** The interactive workflow, run
   on a repository with a readable manifest, skips every question the manifest
   answers and presents the plan directly. It prompts only for decisions the
   manifest does not carry, and for the plan confirmation itself. Its behavior on
   a repository with no manifest — first adoption — is unchanged.

8. **Result document names the outcome.** The refresh emits a structured result
   naming the prior and current catalog identity, the artifacts rewritten, the
   clauses accounted for, the skills touched, the warnings raised, and the
   approved Plan Digest. Exit categories stay consistent with the existing
   Baseline command family so a sweep can branch on them.

## User Experience

The command surface is a CLI. The sweep case is the design center:

```
roundfix baseline update --repo <path> --yes --format json
```

On a current repository it reports no changes and exits successfully — running
it twice in a row is safe and the second run is a proven no-op. On a stale
repository it rewrites the managed artifacts, refreshes the skills, and reports
what moved. On a repository needing a new decision it writes nothing, names the
decision, and exits in the action-required category so the loop can collect
those repositories and hand them to a human.

Without `--yes`, the command presents the plan — file changes first, then the
retention accounting, then the Plan Digest — and writes nothing, so a maintainer
can inspect one repository before trusting the sweep across the rest.

## Non-Goals / Out of Scope

- First adoption. A repository with no Setup Manifest still goes through the
  full interactive interview, including instruction preservation and supervised
  classification. This feature changes nothing about that path.
- Profile changes. An update keeps the manifest's profile. Moving a repository
  to a different Baseline Profile remains the interactive workflow's job,
  because it is a policy decision with a full plan of its own.
- Re-filing repository-authored rules. Rules a maintainer hand-wrote into a root
  carrier after adoption stay exactly where they are. Moving them into a
  canonical carrier is classification, and this feature deliberately does not
  classify.
- Multi-repository orchestration. The command operates on one repository per
  invocation. Sweeping the fleet is a shell loop the maintainer owns; Roundfix
  ships no fleet driver here.
- Changing the external skill acquisition contract. The refresh reuses the
  existing immutable Git source and its provenance checks rather than
  introducing a new acquisition path.
- Repairing repositories whose managed markers a maintainer edited by hand. That
  is a conflict the existing warnings already surface, and it needs a human.

## Success Metrics

- A refresh of an already-current repository asks zero questions and reports
  zero file changes.
- A refresh of a stale repository asks zero questions when the catalog
  introduced no new decision.
- Every one of the nine adopted repositories reaches the current catalog digest
  through this command, with the repositories that stop for a new decision
  naming it rather than failing opaquely.
- A refresh leaves a repository's Repository-Specific Normative Rules
  byte-identical, proven by digest rather than asserted.
- A second refresh immediately after the first is a proven no-op.

## Decisions

- The update is a new command verb rather than a mode of the existing
  interactive command, and the interactive command is additionally taught to
  skip settled questions. A new verb keeps the automation contract additive; the
  interactive fix is what makes a one-off refresh bearable.
- An update never re-classifies repository-authored instructions; it regenerates
  only what the Baseline owns and preserves everything else byte-for-byte.
  See ADR-0099.
- Retention accounting survives the removal of classification and stays
  fail-closed, because it is what ADR-0058 actually requires and it needs no
  analyzer. See ADR-0099.
- The refresh covers both the binary-carried owned skills and the external
  Repository Skill Set, with external restoration best-effort so an unreachable
  upstream degrades to a warning instead of blocking the guidance refresh. This
  narrows the existing contract that adoption never restores skills as a side
  effect: an update restores them because keeping guidance and skills in step is
  the point of the command.
- A decision the manifest does not carry stops the command with a named exit
  rather than silently adopting a default; `--adopt-suggested` is the explicit
  opt-in for a batch sweep, and every value it adopts is reported.

## Open Questions

None.
