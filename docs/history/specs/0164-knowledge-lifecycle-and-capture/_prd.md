---
spec: 0164-knowledge-lifecycle-and-capture
status: archived
created: 2026-09-24
surfaces: [backend, docs]
archived: "2026-09-25"
source_slug: 0164-knowledge-lifecycle-and-capture
---


# Knowledge lifecycle and capture

Delivery 6 of the restructured queue carries four Core Features of the portfolio
Spec 0120 (its Core Features 8, 3, 4 and 2, in that order here). Four defects
remain:

- **A Review retires on whatever the object store holds.** `ClassifyReview` in
  `internal/spec/review_liveness.go` asks local Git about the newest Round's
  head. Measured on 2026-08-14, the same fifty orphan Review Artifacts gave
  three answers in one session as refs were fetched and deleted, and a
  squash-merging repository can never retire one by ancestry. ADR-0123 made that
  rule, so it must be revised before the behaviour changes.
- **Completed work has no truthful terminal state without a Spec.**
  `detectArchiveLicenses` in `internal/speccheck/citations.go` demands an
  `absorbed_by` for every archived Finding, though the layout guide already
  allows `closure_reason` and `closure_evidence`. `ClassifyBacklogEntry` retires
  only `declined`, so a `done` or `superseded` Backlog Entry never reaches
  `docs/history/backlog/`, and nothing reports one left in `docs/backlog/`.
- **The history guide does not name every retired family.** The layout guide
  names the backlog, findings and handoff families, not the ADR, review and Spec
  families that `spec.ArchiveDir` resolves, and the resolver has no handoff
  family at all.
- **Capture commands a commit that another owner makes.** The Secondbrain guide
  says "Commit the entry at the moment of capture", while the installed autosync
  job commits and pushes the shared checkout. A session cannot tell a local file
  from a published one.

## Project Constraints

- Identifier strategy: applicable — dated artifact basenames, Spec slugs, Review
  Artifact directory names and ADR numbers are kept; a terminal disposition
  invents no owner, and the new decision is ADR-0163. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Review retirement reads recorded
  files only; no provider read, credential or network call is added. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — This Spec adds ADR-0163, which supersedes
  ADR-0123 and absorbs the proposed ADR-0152; both move to the history root
  before the Review retirement changes. ADR-0120 puts retired documentation under
  one history root and ADR-0122 keeps a proposed decision pending, not retired.
  ADR-0092 puts triage intent in a typed backlog. ADR-0081 and ADR-0149 govern the
  sanctioned regeneration. ADR-0093 checks Spec consistency by citation,
  ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a path
  governed once bounded, ADR-0155 makes the `qa` Task declare the matrix and
  ADR-0156 makes a declared promise name a consuming Task. This Spec's gate is
  bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization on
  2026-09-24, recorded in [_authorization.md](_authorization.md); bounded files:
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/secondbrain.json`,
  `internal/spec/archive.go`, `internal/spec/archive_test.go`,
  `internal/speccheck/backlog.go`, `internal/speccheck/backlog_test.go`,
  `internal/docscontract/publicdocs_test.go`, `docs/agents/docs-layout.md`,
  `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`. Sanctioned
  regeneration: `make baseline-digests`, `make skills-sync`. The bounded set is
  the intersection of this Spec's changed paths with `GovernedPath`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Goals

- A Review Artifact's retirement is a property of the artifact, not of the
  machine that reads it.
- Completed evidence and intent reach history truthfully without a new Spec.
- Every retired family has one documented home.
- A capture names who publishes it and what was actually observed.

## Core Features

1. **Review retirement reads recorded evidence.** An orphan Review Artifact
   retires on its recorded outcome, and a recorded squash merge with its merge
   commit is a valid receipt. Without a recorded outcome the answer is an
   explicit unknown, and the artifact is retained. No local Git object,
   ref or ancestry changes the answer. A legacy `docs/specs/_reviews/` Review
   Artifact relocates to `docs/specs/reviews/` whatever its liveness.
2. **A terminal record closes without a Spec.** An archived Finding with no
   absorber is licensed by a non-empty `closure_reason` and `closure_evidence`;
   an invalid existing `absorbed_by` stays an error. Every terminal Backlog
   status retires to history when the entry names its Spec or records its
   reason, and a terminal entry left in `docs/backlog/` is reported. An
   unresolved or unknown status is never read as terminal.
3. **Every retired family has one documented home.** The resolver names the
   handoff family beside the other five, and the layout guide names every
   history directory the resolver returns, which ADRs and Review Artifacts
   retire, and how each family records its disposition.
4. **Capture names its publication owner.** The Secondbrain guide names the
   publication owner — an installed scheduler, or the capturing session when none
   is installed — and three evidence levels: captured locally, committed
   locally, confirmed on the remote. A configured scheduler is not a push.

## Non-Goals / Out of Scope

- Writing a recorded outcome from a provider read; no Roundfix flow merges an
  orphan Pull Request today, so the outcome is recorded by whoever observed it.
- The companion Secondbrain contract (`AGENTS.md`, `inbox/README.md` in
  `marcioaltoe/secondbrain`), which needs its own repository and authorization.
- Spec 0120 Core Features 1, 5, 6 and 7, cut in the 2026-09-24 restructuring.
- Rewriting any original observation or archived Spec.

## Success Metrics

1. The same orphan Review Artifact classifies identically whether its recorded
   head is present, absent, or reachable only from fetched refs, and a recorded
   squash receipt retires it with the merge commit absent locally.
2. An archived Finding with closure fields and no absorber passes the archive
   check, one with an invalid `absorbed_by` still fails, and a terminal Backlog
   Entry left in `docs/backlog/` is reported.
3. The rendered layout guide names every history directory the resolver returns,
   and a legacy `docs/specs/_reviews/` artifact relocates whatever its liveness.
4. The rendered Secondbrain guide names the publication owner and the three
   evidence levels, and no longer commands a commit at capture.

## Declared breaks

- An orphan Review Artifact with no recorded outcome, which today can retire on
  ancestry or on an unreachable head whose branch is gone, is retained as
  unknown. `TestClassifyReviewLocalGit` pins the old behaviour and changes with
  it.
- `SC-BACKLOG-UNMOVED` also reports a terminal Backlog Entry in `docs/backlog/`.

## Recorded limits

- The archived-Finding closure check accepts a Finding without an absorber only
  when both `closure_reason` and `closure_evidence` are non-empty, but no test
  pins the one-field and blank-field cases. Found by the pre-PR review of
  2026-09-24; carried as a test-coverage gap.

## Decisions

- **Record, don't rediscover.** Liveness read from the object store depends on a
  machine's fetch and garbage-collection history; a recorded outcome does not.
- **Supersede, don't edit.** ADR-0123 moves to history as superseded by ADR-0163,
  following the repository's consolidation convention.
- **Reuse the existing detector codes.** Closure is judged by
  `SC-ARCHIVE-LICENSE` and `SC-BACKLOG-UNMOVED`, so the detector registry and the
  corpus golden need no new code.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: closure fields that hid an
invalid absorber, or a Review that retired because an object happened to be
present, would pass every happy path.

## Research basis

The Review defects are measured in
`docs/history/findings/2026-08-14-a-review-retires-on-whatever-the-object-store-happens-to-hold.md`.
The closure, history-family and capture defects are the Backlog Entries adopted
by Spec 0120 on 2026-09-08, including the observed autosync job (33 runs, last
exit zero, interval 7,200 seconds), which established a configured owner but no
remote durability. The Core Feature selection is the 2026-09-24 restructuring
recorded in Spec 0120's PRD.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
