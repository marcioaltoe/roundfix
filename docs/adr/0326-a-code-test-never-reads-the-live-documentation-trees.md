---
status: accepted
created_at: 2026-10-24T12:00:00Z
updated_at: 2026-10-24T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A code test never reads the live documentation trees

Go reuses a package's cached test result only when every file the test
process opened inside the module is unchanged. On 2026-10-24 thirteen
`internal/cli` tests read the repository's own `docs/history`, `docs/adr`,
`docs/specs`, `docs/user-guide` or `.agents/skills`: release plan checks on
this repository, repository-profile Baseline updates, the skill set this
repository holds, and the examples and statements its guides carry. With
those trees removed or unreadable, those thirteen fail and no other test of
the package does. In a worktree whose `internal/cli` result was cached, one
line appended to a command guide made the package run again, 162 s and 237 s
in two measurements; with the thirteen tests skipped, the same edit left the
result cached. `internal/cli` is the longest package of the suite, and the
repository's own rules make a guide or skill edit part of many Tasks, so
every Verification of such a Task paid the whole package.

## Decision

**A test that reads the repository's live documentation trees is a
Repository Contract Test.** A test in a code package that opens
`docs/history`, `docs/adr`, `docs/specs`, `docs/user-guide` or
`.agents/skills` of this repository lives in a file under the
`docscontract` build tag that declares `//verify:always`, so the selective
gate runs it on every change and the Full Contract Run runs it on main and
before a release, while `go test ./...` and the package's cached result no
longer depend on documents. A test that needs a document's shape reads a
fixture copy instead and stays a code test.

## Considered Options

- Fixture copies of the documents for every one of the thirteen tests were
  rejected as the default: most of them assert this repository's own
  documents, and a copy would stop checking what ships.
- Skipping the tests with `-skip` in the Makefile was rejected: it hides them
  from every gate that does not repeat the list, and a new reader would
  silently slow the package again.
- Moving them to `internal/docscontract` was rejected: they exercise
  unexported `internal/cli` code.

## Consequences

- A guide or skill edit reuses the `internal/cli` result; the thirteen tests
  run in their own `go test -tags docscontract` invocation, about 150 s of
  test time in parallel, compiled once more with the tag.
- The tests leave `make test` and the push-to-main "Verify" step, and run in
  `make verify-changed` (Runs and pull requests) and `make verify-contracts`
  (main and releases).
- The repository contract declarations test pins the class of these files,
  so a new file of the kind cannot land without `//verify:always`.
