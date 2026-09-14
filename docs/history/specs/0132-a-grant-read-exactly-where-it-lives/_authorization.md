---
status: approved
granted: 2026-09-10
action: repair authorization record parsing, Spec Root resolution, and archive-stable test discovery
consuming: 0132-a-grant-read-exactly-where-it-lives
paths:
  - internal/speccheck/constraints.go
  - internal/speccheck/constraints_characterization_test.go
  - internal/speccheck/governed_repocontract_test.go
  - internal/suiteguardcontract/regeneration.go
  - internal/suiteguardcontract/regeneration_test.go
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0132

The maintainer approved the exact governed mutations above on 2026-09-10. The
frontmatter is the operative statement; this prose explains it and never widens
it.

## Why these four paths

The repairs are mostly in ungoverned source — `internal/authorization`,
`internal/spec` and `internal/speccheck/mechanical.go` — which needs no grant.
Four governed paths are unavoidable:

- `internal/speccheck/constraints.go` resolves a Spec-relative citation and must
  resolve it against the artifact that carries it, so a record beside a PRD in
  an external Spec Root stops being refused.
- `internal/speccheck/constraints_characterization_test.go` records the
  resolution answers this Spec intentionally moves, and holds the rest.
- `internal/speccheck/governed_repocontract_test.go` pins Spec 0119's record to
  its active path, so archiving that Spec breaks the contract; it must discover
  the record instead.
- `internal/suiteguardcontract/regeneration.go` filters discovery on filenames
  containing `authorization` and therefore ignores the existing valid record
  named `2026-09-08-authorized-qa-archive-override.md`.

Every path is already in the Governed Path set, so ADR-0130 keeps them governed
regardless of this Spec; none is added to or removed from that set here.

## Amendment — 2026-09-11 — the discovery filter's own test

`internal/suiteguardcontract/regeneration_test.go` joins the bounded set. The
maintainer answered the explicit request with "Aprovar o teste também",
raising the effective set to five paths.

The pre-Pull-Request review found the omission. The file is bounded by Spec
0119's grant, and Spec 0119's own widening of the Governed Path set covers every
path a discovered record bounds, so it is governed on the implementation branch.
The discovery regression this Spec requires — proving the record named
`2026-09-08-authorized-qa-archive-override.md` is found — belongs beside the
code in that file, and without the path the Task would fail the changed-file
audit.

The alternative was a second test file in the same package, split by
authorization rather than by cohesion. The maintainer chose the boundary over
the split.

## Limits

- No action, operation or path beyond the five above: the four approved on
  2026-09-10 and the test file approved on 2026-09-11.
- No change to what any grant permits, to the operation vocabulary, or to the
  Governed Path set.
- No paid API use, release, tag, deployment, or branch-policy exception.
- Verification remains Daemon-owned; Task status remains Daemon-written.
- The execution trust boundary promoted out of Spec 0119 stays out of scope.

## Approval evidence

On 2026-09-10 the maintainer answered the explicit request for these four exact
paths with "Aprovar os quatro caminhos", approving the bounded scope above.

The decision was requested rather than taken under the standing purpose-bounded
extension authority, because Spec 0119's Core Feature 2 requires a newly
recorded maintainer decision for a grant. A standing delegation is authority to
ask narrowly, not a substitute for the record.

The request named what each path is needed for: the citation resolver that
refuses a valid record in an external Spec Root, the characterization that holds
the resolution answers this Spec moves, the governed contract that pins Spec
0119's record to its active path and breaks on archive, and the discovery filter
that ignores the existing record named with `authorized`.
