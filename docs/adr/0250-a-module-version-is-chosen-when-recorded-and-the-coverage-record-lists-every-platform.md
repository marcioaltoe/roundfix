---
status: accepted
created_at: 2026-10-07T00:00:00Z
updated_at: 2026-10-07T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A module version is chosen when it is recorded, and the coverage record lists every platform

Two checked-in records were wrong for the head CI tested on 2026-10-06 and
2026-10-07, and the operator rewrote each by hand. Specs 0239 and 0241 were
authored in parallel and both wrote `context-workflow` version 22 into their
Tasks; nothing checks a Baseline module's `version`, so only a review noticed,
and the operator raised 0241 to 23. Specs 0242 and 0243 then changed the same
module's content without raising it at all. Separately, the coverage record
written by `go test -list` on macOS listed a test that only macOS builds, and
Linux CI reported it as a regression.

Two rules settle it.

- **A Baseline module's version names one content, and the record step
  chooses it.** `internal/baseline/module-versions.json` holds, for every
  catalog module, each version it has had and the digest of its content with
  the top-level `version` removed. A test refuses a module whose content
  differs from the digest recorded under its version, or whose version is not
  recorded. Its record flag keeps a version that is higher than every recorded
  one, and otherwise writes the next integer above the highest recorded
  version into the module's top-level `version` line, then records it. It
  never replaces or reorders a recorded entry. Every module holds that line
  on its own, so the line-scoped derived merge of ADR-0233 resolves a
  conflict on it. The record step runs before `make baseline-digests` in
  this repository's derived-path declaration for the Baseline catalog, so two
  items that change one module merge with the next free version.
- **The coverage record lists tests for every release platform.** The record
  is collected without running a test binary: `go list` reports each
  package's test files under every `GOOS` in `dist/npm/platforms.json`, and
  the Go parser finds each top-level `TestXxx(*testing.T)` in them. A test
  built on every platform is listed as before; a test built on only some is
  listed with those platforms. The comparison uses the same collection on any
  host, so the record is the same bytes on macOS and Linux and a removed
  macOS-only test fails Linux CI too. A separate test proves that the
  collection for the host's `GOOS` equals what `go test -list` reports.
  Because any host now writes the same record, it is also declared a derived
  path, re-recorded when two items conflict on it.

Comparing only the host's subset of a host-written record was rejected: it
still needs each test's platforms to know what to skip, and it cannot see a
regression on another platform. Listing with `GOOS=<os> go test -list` was
rejected because `go test -list` runs the test binary, and a binary built for
another system fails with `exec format error`. Recording only unconstrained
tests was rejected because it drops the 48 platform-limited tests from
coverage. Making `make baseline-digests` run the record step was rejected
because the Makefile is outside this Spec's authority; the derived-path
declaration and the repository rule run it instead.

This decision follows
[ADR-0189](0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md)
and [ADR-0233](0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md)
for owned skills and applies them to Baseline modules, and uses
[ADR-0192](0192-a-conflict-confined-to-declared-derived-paths-is-resolved-by-regeneration.md)
unchanged. It supersedes none of them.

## Consequences

A module edit that does not run the record step fails the suite with the
command to run. Rule, guide and clause versions inside a module stay
hand-maintained. A version number may be skipped when a Task records twice,
as ADR-0233 accepts for skills. Adding a release platform changes the
coverage record, which is then re-recorded deliberately.
