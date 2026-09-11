---
spec: 0132-a-grant-read-exactly-where-it-lives
status: active
created: 2026-09-10
surfaces: [backend, cli]
---

# A grant read exactly where it lives

Spec 0119 made the authorization record the thing that decides whether governed
work may proceed. Its pre-Pull-Request review found that the reader accepts a
record it should refuse, and refuses records it should accept. A malformed
frontmatter delimiter still yields a grant, and a Spec Root configured outside
the code repository is unreadable because the record path is assembled from a
constant instead of the resolved root. Two of 0119's own tests also pin the
record to its active path, so archiving the Spec breaks them.

The defects were found on 2026-09-10 by review of the 0119 candidate, after its
terminal QA gate passed. They cannot be repaired inside 0119: its gate is
settled, and adding work under a passing gate would make that pass certify work
the gate never saw. They are repaired here, on the same branch, before either
Spec reaches the target.

## Project Constraints

- Identifier strategy: applicable — preserve Spec slugs, Task IDs, refusal codes and record field names; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched; the record is read from the local filesystem and Git objects as it is today. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0130 keeps the Governed Path set monotonic, so no path leaves it; ADR-0057 keeps the Daemon the exclusive writer of Task status; ADR-0117 places each check at the stage that can establish it, and record shape is established at parse time; ADR-0096 requires mechanical facts before the QA Agent turn and is preserved unchanged. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: "Aprovar os quatro caminhos", 2026-09-10, extended by "Aprovar o teste também", 2026-09-11, both recorded in `docs/specs/0132-a-grant-read-exactly-where-it-lives/_authorization.md`; bounded files: `internal/speccheck/constraints.go`, `internal/speccheck/constraints_characterization_test.go`, `internal/speccheck/governed_repocontract_test.go`, `internal/suiteguardcontract/regeneration.go`, `internal/suiteguardcontract/regeneration_test.go`. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A record grants only when every part a reader relies on is exactly what the
  format requires, so a malformed record refuses instead of authorizing.
- A record is read where the repository says Specs live, so a configured Spec
  Root outside the code repository resolves rather than refusing valid work.
- A test that reads a Spec's record keeps working when that Spec is archived.

## User Stories

1. As the maintainer, I want a record whose frontmatter is malformed to grant
   nothing, so that a near-miss on the format cannot become authority.
2. As a maintainer whose Specs live in a separate knowledge repository, I want
   Implement, Settle and Spec Check to resolve my grants, so that a configured
   Spec Root is usable without duplicating records into the code repository.
3. As the Supervisor, I want archiving a Spec to leave the suite green, so that
   the last step before a Pull Request does not break the gate that just passed.

## Core Features

1. Frontmatter is delimited by complete marker lines. A purported closing line
   carrying anything besides the marker is malformed, and a malformed record
   resolves to a refusal naming the offending line rather than to a grant.
2. The consuming record's path and revision derive from the Spec repository,
   and the bounded paths it declares are judged against the project repository.
   The two are named separately, so an external root is read where it lives
   while the audit still judges the tree whose changes it governs.
3. A Spec-relative citation resolves against the artifact that carries it, so a
   record beside its PRD resolves whether that PRD sits inside the code
   repository or outside it.
4. Tests that read a Spec's authorization record discover it wherever it
   currently lives, and corpus invariants hold over the whole corpus rather than
   requiring at least one active Spec.
5. Record discovery recognizes the naming already in use for authorization
   records, so an existing valid record is not silently ignored.

## User Experience

The maintainer sees no new prompt or flag. A malformed record produces a
refusal that names the line that made it malformed. A configured external Spec
Root behaves like the default one. Archiving a Spec changes no test outcome.

## Non-Goals / Out of Scope

- The execution trust boundary promoted out of Spec 0119. Its symlink identity,
  revision pinning and per-command cost belong to that promoted Spec.
- The residual archive-traversal cost in the suite guard, which is repository
  performance work rather than a correctness defect.
- Any change to what a grant permits, to the operation vocabulary, or to the
  Governed Path set.

## Success Metrics

- A record whose closing delimiter carries extra characters refuses, and the
  refusal names that line.
- A Spec Root outside the code repository resolves its grant through Spec
  Check, Implement dispatch and Settle.
- Archiving a Spec leaves every test that reads its record passing, proven by
  running the suite against both the active and the archived layout.

## Decisions

- Repaired here rather than in Spec 0119, because 0119's QA gate is settled and
  work added under a passing gate would make that pass certify what it never
  saw. The refusal that blocks it is correct and stays.
- Delivered on the same branch as Spec 0119 rather than after it, because these
  defects are in code that has not reached the target, and merging first would
  place a record that grants on malformed input into the default branch.

### Declared intentional breaks

1. A record whose frontmatter closing line carries anything besides the marker
   stops resolving to a grant and starts refusing.
2. A test or invariant that assumed the authorization record sits at its active
   path now reads it wherever it lives; the recorded expectation moves with it.

Everything else keeps behaving as it does: a well-formed record grants exactly
what it granted before, the Governed Path set does not change, refusal codes
keep their tokens, and Daemon ownership of Task status is untouched.

### Regression locks

- Today's parse and resolution answers are captured as characterization before
  the change, including the answers this Spec intends to move.
- Every new refusal is fail-closed on evidence: it fires on a positive
  observation of a malformed delimiter or an unresolvable root, never on an
  unreadable or merely unfamiliar input, which stays unresolved as it is today.

## Open Questions

None.
