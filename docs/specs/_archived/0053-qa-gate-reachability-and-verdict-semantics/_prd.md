---
spec: 0053-qa-gate-reachability-and-verdict-semantics
status: archived
created: 2026-07-28
surfaces: [backend, cli, docs]
archived: "2026-07-30"
source_slug: 0053-qa-gate-reachability-and-verdict-semantics
---


# QA gate reachability and verdict semantics

The QA gate runs in a Run Worktree with no commit or push authority, so any
acceptance row whose journey needs a live Pull Request or Final Push is
structurally unreachable: it is recorded `blocked`, the verdict caps at
`partial`, and the Spec can never archive through the normal loop — Spec 0039
needed a maintainer `qa_override` for exactly this. Meanwhile every failed QA
attempt strands a Run Branch holding a superseded report that Branch Integrity
Preflight then counts against the next review Run, and the Daemon's QA prompt
and the qa-gate Skill still disagree about the report filename. Evidence:
[QA gate cannot reach Pull Request journeys](../../findings/2026-07-28-qa-gate-cannot-reach-pull-request-journeys.md),
[failed QA Runs strand branches](../../findings/2026-07-28-failed-qa-runs-strand-branches-that-block-review-runs.md),
and the naming defect recorded in
[same-day QA reruns were ignored](../../findings/2026-07-28-same-day-qa-reruns-are-ignored-by-the-verdict-selector.md).

## Project Constraints

- Identifier strategy: not applicable — this Spec creates no project-owned
  Internal Identifier; Run IDs, Git heads, QA report paths, and Review
  Source-native identities are reused. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — read-only Pull Request observation
  must continue through the repository's existing `gh` and Review Source
  boundaries; this Spec adds no authentication provider, credential policy,
  or HTTP route. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0015 keeps the QA gate a phase
  inside the Spec Run; ADR-0051 keeps QA in its own Agent Session; ADR-0052
  protects terminal completion; ADR-0053 requires reconciliation to stay
  explicit and proof-based, which the new `superseded` classification must
  honor; ADR-0057 keeps Task status Daemon-owned; ADR-0036 preserves the
  separate review-artifact commit the inheritance journeys observe; ADR-0080
  defines the environment-blocked verdict rule. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.agents/skills/qa-gate/SKILL.md`,
  `skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, and
  `skills/roundfix/SKILL.md`, plus the deterministic Skill-digest fallout in
  exactly `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json` (a
  non-roundfix Skill edit cascades into all three setup snapshots; extension
  granted 2026-07-28). No other protected tooling mutation is authorized.
  Source: `docs/agents/agent-instructions.md`.

## Goals

- A Spec whose acceptance matrix includes live Pull Request journeys can
  reach `pass` and archive through the normal loop, without granting Agents
  commit or push authority.
- QA reports distinguish rows blocked by the gate's environment from rows
  blocked by a finding, and the verdict rule treats them differently.
- Consecutive failed QA attempts leave no stranded branch that blocks the
  next review Run or invites an integration that would overwrite a newer
  artifact.
- The Daemon QA prompt and the qa-gate Skill state one report-naming
  contract.

## User Stories

1. As a maintainer archiving a Spec whose only non-pass rows are
   environment-blocked with recorded equivalent evidence, I want the verdict
   to reach `pass`, so that the Spec archives without a `qa_override`.
2. As a QA Agent, I want a read-only Pull Request journey mode that observes
   approval, Merge-Ready acceptance, and artifact-inheritance evidence
   against the real Pull Request resolved from the Run's target branch, so
   that the matrix's most valuable rows execute without push authority.
3. As a user starting a review Run after a Spec needed several QA attempts,
   I want superseded QA-report branches classified and released by the
   Reconcile Command, so that Branch Integrity Preflight does not block the
   Run on debris.
4. As an operator reading reconcile output, I want no suggested next action
   whose integration would overwrite a newer artifact on the target branch,
   so that following the printed remedy is always safe.
5. As a reader comparing QA reports, I want blocked-by-environment and
   blocked-by-finding visibly distinct, so that `partial` caused by the
   environment is not mistaken for incomplete diligence.
6. As a QA Agent rerunning the gate on the same day, I want one filename
   contract from both the Daemon prompt and the qa-gate Skill, so that both
   surfaces produce the same suffixed sibling artifact.

## Core Features

1. Read-only Pull Request journeys: the QA surface resolves the Pull Request
   from the Run's recorded target branch and repository — never from the Run
   Worktree's checked-out branch — and observes approval state, check and
   status evidence, unresolved threads, and Daemon review-artifact
   descendants without any commit, push, or Review Source mutation.
2. Blocked rows carry a typed cause: `blocked-by-environment` for journeys
   the gate's boundary cannot execute, `blocked-by-finding` for rows a
   defect prevents. The report and its coverage summary render the
   distinction.
3. The verdict rule per ADR-0080: finding-blocked or failed rows cap the
   verdict exactly as today; a matrix whose only non-pass rows are
   environment-blocked reaches `pass` when the report records the equivalent
   observed or supervised evidence for each such row, and stays `partial`
   when it does not.
4. The Reconcile Command gains a `superseded` classification: a terminal Run
   whose only unintegrated commits touch paths a later Run already wrote on
   the target branch. `--apply` releases superseded Runs with the reason
   recorded; genuinely unintegrated work keeps its preserve-by-default
   behavior.
5. No printed next action proposes a fast-forward integration that would
   overwrite a newer artifact on the target branch; the
   conflict-by-supersession is detected and named instead.
6. A Run Branch holding only a superseded QA report is never silently
   integrated into the user's head branch: Branch Integrity Preflight's
   automatic fast-forward integration recognizes QA-report-only branches and
   routes them to supersession handling instead of merging a failing report
   onto the Pull Request.
7. Branch Integrity Preflight names which pending Run Branches reconcile
   classifies as superseded, so the operator is never left inspecting each
   branch to answer a question Roundfix can compute.
8. One report-naming contract — dated report with a collision-safe suffix for
   same-day reruns — stated identically by the Daemon QA prompt and the
   qa-gate Skill, and ranked correctly by the newest-report selector for
   every suffix form the contract allows.

## User Experience

- A QA report's blocked rows read `blocked (environment: no push authority)`
  or `blocked (finding: <id>)`, and the coverage section counts the two
  separately.
- `roundfix reconcile` lists superseded Runs with the newer Run and target
  evidence that proves supersession; `--apply` reports each release reason.
- Branch Integrity Preflight's refusal distinguishes branches needing human
  judgment from branches reconcile can release, and names the command.
- Nothing in the QA flow prints a zero-issue or pass claim for a journey it
  did not execute or observe.

## Non-Goals / Out of Scope

- Granting QA Agents commit, push, or Pull Request mutation authority.
- Weakening reconcile's preserve-by-default for work that is not proven
  superseded, or any automatic release without `--apply`.
- Changing `roundfix archive`'s requirement that the newest report's verdict
  be `pass`; `qa_override` remains a human decision.
- Crediting unexecuted journeys from Task evidence or Result prose — the
  gate's refusal to guess is preserved.
- Re-implementing the report recency ordering fixed on 2026-07-28.

## Success Metrics

- A Spec whose matrix includes Pull Request journeys reaches `pass` when
  those journeys succeed against a real Pull Request observed read-only.
- The QA Agent resolves the Pull Request without depending on the Run
  Worktree's branch name.
- Three consecutive failing QA gates followed by a pass leave no branch that
  blocks a review Run; `reconcile --apply` releases the superseded branches
  and preserves genuinely unintegrated work.
- No reconcile or preflight output suggests an integration that would
  overwrite a newer artifact.
- The Daemon prompt and the qa-gate Skill state the same report filename
  contract, and a same-day rerun creates a suffixed sibling on both
  surfaces.

## Decisions

- Environment-blocked rows do not cap the verdict when their equivalent
  evidence is recorded; finding-blocked rows always do. See
  [ADR-0080](../../adr/0080-qa-verdicts-distinguish-environment-blocked-rows.md).
- The `superseded` classification stays proof-based per ADR-0053: release
  requires proving every unintegrated commit's paths were rewritten by a
  later Run on the target branch.
- The naming contract is reconciled by aligning the Daemon prompt and the
  qa-gate Skill to the suffixed-sibling scheme the recency fix already
  understands.

## Open Questions

None.
