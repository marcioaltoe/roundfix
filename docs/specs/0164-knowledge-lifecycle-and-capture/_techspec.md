---
spec: 0164-knowledge-lifecycle-and-capture
status: active
created: 2026-09-24
surfaces: [backend, docs]
---

# Knowledge lifecycle and capture

## Executive Summary

Supersede ADR-0123 with ADR-0163, then make Review retirement read a recorded
outcome instead of local Git, license a terminal record by its closure fields,
name every history family from the resolver, and give capture a named
publication owner with three evidence levels.

## Project Constraints

- Identifier strategy: applicable — existing basenames, slugs, Review Artifact
  names and ADR numbers are kept; the new decision is ADR-0163. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — recorded files only; no provider
  read, credential or network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — This Spec adds ADR-0163, superseding
  ADR-0123 and the proposed ADR-0152. ADR-0081, ADR-0092, ADR-0120, ADR-0122 and
  ADR-0149 hold, as do ADR-0080, ADR-0091, ADR-0093, ADR-0096, ADR-0097,
  ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization on
  2026-09-24, recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/secondbrain.json`,
  `internal/spec/archive.go`, `internal/spec/archive_test.go`,
  `internal/speccheck/backlog.go`, `internal/speccheck/backlog_test.go`,
  `internal/docscontract/publicdocs_test.go`, `docs/agents/docs-layout.md`,
  `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`. Sanctioned
  regeneration: `make baseline-digests`, `make skills-sync`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## The decision

`docs/adr/0163-review-artifact-retirement-reads-recorded-evidence.md` states
that an orphan Review Artifact retires on its recorded outcome, that a missing
outcome is an explicit unknown, that relocating the legacy `_reviews` root is
independent of liveness, and — carried over from ADR-0123 — that the Review
Artifact root resolver never resolves into the history root. ADR-0123 and
ADR-0152 move to `docs/history/adr/` with `status: superseded`,
`superseded_by: ADR-0163` and a new `updated_at`, and ADR-0163 records
"Supersedes ADR-0123 and ADR-0152", the consolidation form earlier
superseding decisions use.

## Review retirement

`ClassifyReview` reads `outcome.md` at the Review Artifact root. Its front
matter carries `pull_request_state` (`merged`, `closed` or `open`), `merge_commit`
and `recorded_at`. `merged` with a hexadecimal `merge_commit` is a squash
receipt and answers `ReviewFinished`; `closed` answers `ReviewFinished`; `open`
answers `ReviewLive`. A missing or malformed record, or `merged` without a merge
commit, answers `ReviewUndecidable` with a reason naming what is missing. No Git
command runs. In `internal/baseline/history_layout.go`, a live or undecidable
artifact under `docs/specs/_reviews/` relocates to `docs/specs/reviews/`, and a
finished one to `docs/history/reviews/` as today, through the existing collision
checks.

## Terminal dispositions

`parseFindingFrontmatter` reads `closure_reason` and `closure_evidence`.
`detectArchiveLicenses` accepts an archived Finding with no `absorbed_by` when
both are non-empty, and still reports a present, unresolvable `absorbed_by`.
`ClassifyBacklogEntry` keeps retiring `declined`, and retires `done`,
`deprecated`, `superseded`, `closed` and `cancelled` when the entry names its
`spec` or a non-empty `reason`. `detectBacklogPromotion` reports, under
`SC-BACKLOG-UNMOVED`, a terminal entry still in `docs/backlog/`, naming
`docs/history/backlog/` and, without a `spec`, the missing `reason`.

## History families

`internal/spec/archive.go` gains `ArchiveKindHandoff` (`docs/history/handoffs`)
and `ArchiveKinds()`, which lists every kind. The clause
`clause.context.docs-one-job-per-directory` names each directory those kinds
resolve and how each family records its disposition: a dated disposition
addendum for Findings and Backlog Entries, lifecycle front matter for ADRs, the
archive stamp for Specs, and a byte-identical move for Review Artifacts and
handoffs. The guide and manifest are rendered by the public Baseline update.

## Capture ownership

The capture sentence in `clause.secondbrain.capture-self-contained` becomes a
publication contract. The publication owner is an installed scheduler that
commits and pushes the checkout, or the capturing session when none is
installed. Under a scheduler the session writes and does not commit. Without one
it commits and pushes only its own entry's path, and never installs a second
schedule. A session reports only the level it observed: captured locally,
committed locally, or confirmed on the remote.

## API Contracts

1. A Review Artifact's `outcome.md` front matter — `pull_request_state`,
   `merge_commit`, `recorded_at` — is the only input to its retirement.
2. An archived Finding is licensed by a resolvable `absorbed_by` or by non-empty
   `closure_reason` and `closure_evidence`.
3. `spec.ArchiveKinds()` lists every retired family, including
   `ArchiveKindHandoff`.
4. `SC-BACKLOG-UNMOVED` reports a terminal Backlog Entry left in `docs/backlog/`.

## Coverage Map

- Goal 1 → The decision; Review retirement; API Contract 1.
- Goal 2 → Terminal dispositions; API Contracts 2 and 4.
- Goal 3 → History families; API Contract 3.
- Goal 4 → Capture ownership.
- Core Feature 1 → The decision; Review retirement.
- Core Feature 2 → Terminal dispositions.
- Core Feature 3 → History families.
- Core Feature 4 → Capture ownership.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 1 and 3.
- Success Metric 4 → Testing Approach 4.
- API Contracts 1-4 → Review retirement, Terminal dispositions, History families.

## Integration Points

- **Spec 0120.** The portfolio Spec whose Core Features 2, 3, 4 and 8 this Spec
  delivers.
- **`internal/baseline` history migration.** Consumes `ClassifyReview`,
  `ClassifyBacklogEntry` and `ArchiveDir`.
- **Secondbrain.** The companion Inbox contract follows in its own repository.

## Testing Approach

1. **Review retirement.** With real Git, one artifact classifies the same with
   its head present, absent and reachable only from fetched refs; a squash
   receipt retires it with the merge commit absent; no record is unknown; a
   legacy `_reviews` artifact relocates whether live, unknown or finished.
2. **Closure.** Closure fields license an absorber-less archived Finding;
   closure fields beside an invalid `absorbed_by` still fail; each terminal
   Backlog status retires with `spec` or `reason`, `open` and unknown statuses do
   not; a terminal entry in `docs/backlog/` is reported.
3. **History families.** `ArchiveKinds()` includes handoffs, and the rendered
   guide names every directory it resolves.
4. **Capture.** The module and the rendered guide carry the owner and the three
   levels and not the old commit command; the catalog stays compatible.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Supersede ADR-0123 with ADR-0163 (depends on: none).
2. Review retirement reads recorded evidence (depends on: 1).
3. A terminal record closes without a Spec (depends on: 2).
4. Every retired family has one documented home (depends on: 3).
5. Capture names its publication owner (depends on: 4).
6. Terminal QA (depends on: 1, 2, 3, 4, 5).

## Risks & Considerations

- **Retaining more than today.** Without recorded outcomes, orphan Reviews stay
  in the live root; that costs a directory and is declared, where a wrong
  retirement breaks a running Round.
- **Hiding a broken license.** Closure fields must never mask an invalid
  `absorbed_by`; a negative test pins it.
- **Shared generated files.** Tasks 04 and 05 both rewrite the manifest and the
  derived pins, so they run in sequence.
