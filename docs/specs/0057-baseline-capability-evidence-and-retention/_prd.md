---
spec: 0057-baseline-capability-evidence-and-retention
status: active
created: 2026-07-28
surfaces: [backend, cli, docs]
---

# Baseline capability evidence and retention accounting

Two live Baseline sessions — a profile refresh in fluxus and a greenfield
alignment in vortex — exposed the same defect family: the Baseline Command's
divergences are individually correct but do not carry enough evidence to be
acted on, and a changed Profile under an unchanged Baseline identity bypasses
retention accounting entirely, letting managed Normative Clauses disappear
with an empty Upgrade Retention Contract ledger. Capability probes reject
executable symlinks (the norm on Homebrew and Docker Desktop machines),
divergence output hides the probe it evaluated, there is no read-only way to
re-check capabilities after remediation, and a clean adoption warns about the
thirteen files it just wrote. Evidence:
[profile refresh applied without semantic retention accounting](../../findings/2026-07-26-baseline-profile-refresh-retention-gap.md)
and
[capability divergences do not carry enough evidence to be remediated](../../findings/2026-07-26-vortex-baseline-capability-remediation.md).
Specs 0049–0051 fixed Preservation idempotency and Doctor hardening; none of
this report's items were covered there.

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier
  is created; Profile identifiers, digests, and diagnostic codes keep their
  existing contracts. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Baseline planning stays local,
  offline, and read-only until digest-confirmed apply; no authentication or
  HTTP surface. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0058 requires upgrades to fail
  closed on unaccounted managed-rule removal, which the same-identity drift
  path currently bypasses and this Spec enforces; ADR-0068 keeps one
  confirmation-gated workflow; ADR-0070 keeps the audit byte-exhaustive
  while preserving root instructions; ADR-0071 keeps plans portable and
  preimage-bound; ADR-0075 keeps divergence resolution a confirmed
  repository-owned adaptation. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`, plus the deterministic Skill-digest fallout in
  exactly `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- A changed Profile can never produce a ready plan with empty retention
  accounting: every previous managed Normative Clause gets one explicit
  disposition or planning stops action-required.
- Every unsatisfied capability divergence renders the probe it evaluated,
  so remediation never requires reading catalog assets or Go source.
- Capability evidence accepts the installs that actually exist on
  maintainer machines — executable symlink chains — and can be re-checked
  read-only after remediation.
- Baseline results separate what was proven from what was not: managed
  carriers stop warning about themselves, and apply success stops implying
  semantic readiness.

## User Stories

1. As a maintainer updating a repository whose Profile changed under the
   same Baseline identity, I want a plan that accounts for every previous
   managed clause — retained, moved, replaced, or explicitly rejected — or
   an action-required stop, so that rules cannot silently disappear.
2. As a maintainer reading a blocking capability divergence, I want the
   evaluated probe rendered — each inspected path with its state and the
   expected content, or the inspected PATH candidate — so that I can
   remediate without reading source.
3. As a maintainer whose tools are installed through Homebrew or Docker
   Desktop symlinks, I want executable discovery to resolve a bounded
   symlink chain, so that installed working tools stop reporting as
   missing.
4. As a maintainer who just remediated a blocking divergence, I want a
   read-only capability re-check that needs no decisions, so that the
   remediate-and-re-check loop exists.
5. As a maintainer facing a mixed divergence set, I want an explicit
   "remediate in the repository and re-run" outcome distinct from decline,
   so that pausing for repository work is not recorded as refusing the
   Baseline.
6. As a maintainer re-planning after a verified apply, I want zero warnings
   about the managed guides the apply just wrote, so that a real nested
   carrier conflict is not buried in noise.
7. As a maintainer reading the final result, I want transaction, retention,
   alignment, Verification, and idempotence reported as separate states, so
   that "verified postimages" is never read as "update complete".

## Core Features

1. Full Setup Manifest validity is the first compatibility test: a matching
   Baseline identifier with changed Profile or catalog digests requires a
   retention transition keyed by the source tuple, or planning exits
   action-required. An empty retention ledger can never accompany a ready
   update plan (per ADR-0058), covered by a regression fixture with
   unchanged identity, changed digests, and a disappearing clause.
2. Carrier classification distinguishes current managed artifacts (no
   warning), stale managed artifacts (upgrade input requiring clause-level
   retention), recognized repository extensions, and unmanaged nested
   carriers (the only ones that warn), with warnings naming the managed
   versus preserved byte boundary.
3. A clause-level semantic delta renders before final confirmation: every
   previous clause classified retained, moved, replaced,
   repository-document, repository-extension, reasoned-rejection, or
   unaccounted, with counts; apply is not offered while any clause is
   unaccounted.
4. Executable capability discovery resolves a bounded symlink chain to a
   regular executable, rejecting cycles, broken links, and non-executable
   targets with distinct diagnostics, and always reports the inspected
   candidate or the absence of one.
5. Every unsatisfied divergence renders its probe: declared-file probes
   list each inspected path with state and expected content; executable
   probes name the PATH candidate; required stack capabilities state the
   selected technology, both resolutions, and whether removal cascades to
   any decision.
6. A read-only capability re-check is available without resolving
   decisions, matching the outcomes a full plan would produce.
7. The divergence prompt gains a fourth outcome — exit without writing,
   print per-divergence remediation, name the re-run command — journaled
   distinctly from decline; the adaptation option documents its
   removal-only constraint.
8. Divergences group by requirement strength — blocking, advisory,
   informational — and every advisory states that it does not block
   readiness or apply before any optional next action.
9. A portable Verification role can be satisfied by mapping it to a
   declared repository command, so a selected repository gate stops
   re-appearing as an unresolved workspace divergence.
10. The final result is a status matrix — approved postimages, semantic
    retention, profile alignment, repository Verification, idempotence —
    and completion language is used only when retention is verified and
    the idempotence check passed.

## User Experience

- A blocking divergence reads as probe evidence, not prose: inspected
  paths with states, expected content, and the exact re-check command.
- The consolidated review shows the clause-level delta and its counts
  before the file ledger; the ledger remains for machine review.
- The re-check loop is one read-only command; the prompt's new outcome
  prints it.
- The final matrix distinguishes `verified` from `not run` per axis.

## Non-Goals / Out of Scope

- Changing the fail-closed apply, digest confirmation, or preimage binding.
- Auto-inferring capability satisfaction, decisions, or adaptations —
  explicit confirmation stays per ADR-0075.
- Executing repository Verification, formatters, or builds during plan or
  apply.
- Re-litigating Preservation idempotency or Doctor readiness (Specs
  0049–0051).
- Profile catalog content changes beyond what probe rendering requires.

## Success Metrics

- The same-identity drift fixture with a disappearing clause exits
  action-required with an explicit unaccounted count; no ready plan carries
  an empty retention ledger when clauses changed.
- An idempotent re-plan after a verified apply reports zero file changes
  and zero warnings; an unmanaged nested carrier still warns.
- `rtk` and Docker installed through symlinks report satisfied evidence
  without executing the binary; a broken symlink and a cycle produce
  distinct diagnostics.
- A declared-file divergence renders every inspected path and state; a
  capability re-check is obtainable with zero decisions supplied and
  matches full-plan outcomes.
- Selecting the remediate-and-re-run outcome is distinguishable in the
  journal from a decline.
- A mapped repository gate satisfies its portable Verification role and the
  divergence disappears.

## Decisions

- Same-identity Profile drift is an upgrade, not a refresh: retention
  accounting or action-required, enforcing ADR-0058 on the path that
  bypassed it.
- Probe evidence is rendered from the catalog's own probe definitions —
  the diagnostic and the evaluation share one source.
- The remediation loop is read-only by construction; only apply writes.

## Open Questions

None.
