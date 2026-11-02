---
status: accepted
created_at: 2026-11-01T00:00:00Z
updated_at: 2026-11-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Spec's Release Ceiling answers the release plan's approval

A maintainer who authorizes a release while a Spec is planned had no way to
make that decision reach the release: the Release Plan Command stopped at
`approval_required`, and a later session asked again (Argus, ai-memory message
`01a19920`, 2026-11-01). The maintainer chose "Teto de impacto (Recommended)":
the Spec's authorization record carries the highest impact its release grant
answers, and the plan accepts without asking while the fresh plan stays at or
below it; a breaking plan, or one above the ceiling, asks again. ADR-0143
keeps the plan read-only and forbids a generic request from silently
authorizing a version decision; a recorded, dated maintainer decision bound to
the Specs the range delivers is neither generic nor silent.

## Decision

**The record.** An approved, non-withdrawn `_authorization.md` that lists
`release` in `operations` may carry `release_ceiling` (`patch`, `minor` or
`major`) and `release_paths`, the exact repository-relative paths the release
may change, written as `executables` entries are (ADR-0340). A ceiling without
`release`, an unknown ceiling, a ceiling without paths, or paths without a
ceiling refuses the whole record. `release` without a ceiling grants what it
granted before and answers nothing.

**Finding the records.** The plan stays local and reads no Roundfix
configuration. A commit of the range delivers a Spec when it deletes
`<dir>/<slug>/_authorization.md` and adds an Archive Record `<slug>.md` in the
same change, as every Delivery Queue squash does. The record is read at that
commit's first parent, where the Spec folder still exists, with the reader the
Delivery Queue uses, asking for `<slug>`.

**When the records answer.** Only a plan whose decision is
`approval_required`, and only when all of these hold: the plan is not
breaking; no `--impact` was given; the range delivers at least one Spec; every
delivered Spec's record grants `release` with a ceiling at or above the plan's
impact; and every commit whose automatic impact equals the plan's impact and
crosses the maintenance boundary is one of those deliveries. Then the decision
is `ready` with the same proposed version and the plan names the records. A
plan whose records do not answer keeps `approval_required` and names the first
failing condition. When no delivered Spec's record carries a ceiling, the plan
prints and encodes the same bytes as before. A Lagging Surface still blocks
(ADR-0256).

**Backlog Entries are planning evidence.** `docs/backlog/` joins
`docs/specs/`, `docs/adr/`, `docs/findings/` and `docs/handoffs/` as
maintenance-only paths. Every plan commit since cycle AR adopted or filed a
Backlog Entry and therefore needed manual classification, which a ceiling
never answers.

## Considered Options

- Answering a manual classification from the ceiling was rejected by the
  maintainer's decision: a classification by hand is a new decision.
- Reading the record at the Archive Record's `source_revision` was rejected:
  after a squash merge that commit is only on the pull request head, which the
  local plan cannot fetch.
- Letting one Spec's ceiling answer for the whole range was rejected: a direct
  `feat:` commit or another Spec's delivery was never part of that decision.
- A new decision state, such as `approved`, was rejected: adopters and the
  follow-up mode already read `ready` and exit 0 as "proceed".

## Consequences

- A range that delivers only Specs whose records carry a sufficient ceiling
  releases without a second question; any other range asks as before.
- A plan commit that changes only planning paths and Backlog Entries no longer
  needs `--impact` and `--reason`.
- The record's `release_paths` are reported, not enforced: the release runbook
  keeps owning which files a release Pull Request changes.
