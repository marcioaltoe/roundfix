# Provenance and absence checks

`golang-dependency-management`, `golang-safety`, and
`golang-structs-interfaces` exist under `.agents/skills/`; the Go Baseline
module and generated dispatch route to them. Fresh searches found none of the
three names in `skills-lock.json` or any setup snapshot, and `roundfix skills
list` omitted all three. The PRD's declared cause is exact: the operative
authorization explicitly excludes third-party Skills and does not grant a
`skills-lock.json` change. The declared unblocking action is a maintainer grant
naming that file and the three Skills, followed by a bounded provenance Task.

Other fresh absence and documentation evidence:

- `make skills-version-check`, `make skills-sync-check`, and the Roundfix
  canonical/mirror recursive comparison exited 0.
- Injecting unversioned `coding-guidelines` into `OWNED_SKILLS` exited 2 and
  named the missing top-level declaration.
- The eight implementation-range Go test files contain no `0.0.2` or 64-hex
  assertion.
- No `contentDigest` appears in setup snapshots, catalog diagnostics, or the
  plan-characterization corpus.
- The Roundfix Skill names the minimum comparison, all three states, the
  four-fact failure, the third-party boundary, and the no-regeneration rule.
- Generated-guide strictness tests passed; archived Specs and Spec 0067's
  derived-ownership record are byte-identical to the pre-Spec range.
