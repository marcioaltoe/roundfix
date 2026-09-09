---
spec: 0130-documentation-cleanup-compatibility
prd: _prd.md
created: 2026-09-09
---

# Documentation cleanup compatibility — Technical Spec

## Project Constraints

- Identifier strategy: applicable — preserve Spec slugs, repository-relative paths,
  immutable Git object identities and the existing ownership tokens; introduce no
  new domain term. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — filesystem and local Git reads only;
  no credential, HTTP interface, external mutation or paid call is introduced.
  Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0149 states that the grant names the regeneration command and the tree names its outputs; the suite guard and audit read the same ownership declaration. This repair preserves that command/ownership boundary. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: "Autorizo os reparos", 2026-09-09, recorded in [the approved repair grant](_authorization.md); bounded files: `internal/baseline/derived_ownership_test.go`, `internal/speccheck/mechanical_test.go`, `internal/speccheck/governed_repocontract_test.go`, `internal/suiteguardcontract/regeneration.go`, `internal/suiteguardcontract/regeneration_test.go`. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Design

### Operative grant discovery

Keep discovery in the existing suiteguardcontract reader. Read the legacy
`docs/workflow/authorizations` directory when present and Spec-contained records
under `docs/specs/<slug>/`, including `references/`. Spec grants need valid YAML
frontmatter, `status: approved`, a valid granted date, a non-empty action, exact
bounded paths and a matching consuming Spec. An unrelated Markdown document,
proposed/null/malformed grant or mismatched consumer grants nothing; malformed
frontmatter never falls back to legacy prose interpretation. Preserve existing
legacy unfrontmattered declaration support. Do not follow grant symlinks or read
history as operative authority. Missing optional roots mean no records; actual
read errors remain errors. Results are deterministic.

### Ownership-based regeneration

A grant sanctions a command; the tree owns its outputs under ADR-0149. For a
command-only declaration, read `DERIVED_DIGEST_PATHS` and `_ownership.yml`/`.yaml`
with the existing sanctioned/dedicated/frozen, nearest-directory, sidecar and
exception semantics. Return only sorted regular output files, excluding ownership
records and frozen or differently owned files. Preserve explicit legacy outputs.
Reject unsafe paths, symlink escapes, malformed/conflicting ownership and incomplete
resolution; never return a partially authorized set on a read error.

Do not import baseline from suiteguardcontract: Baseline's test guard would form
an import cycle. Both consumers read the same declarations; the authorized
Baseline test file compares their results against each other and against concrete
expected sets for real and isolated regular-file trees. This is parity of readers,
not a second inventory. Filesystem-only tests belong in regeneration_test.go;
that package must not gain subprocesses or an unapproved suiteguard inventory edit.

### Historical evidence and fixtures

Use the permanent main ancestor
`81a6afb48f4a3683d0e5fad52f3919cf1bdfbbf4` to recover the 42 historical grant
records, including the proof-cost enumeration. Do not pin an unmerged PR commit:
squash and branch cleanup must leave the source reachable. Existing historical
Task SHAs may be unavailable; mark those cases explicitly skipped rather than
logging and returning a successful assertion. Separately replay a recovered real
grant against controlled temporary Git commits and assert both permitted and
unpermitted governed changes, so missing old Task objects cannot empty coverage.
Require non-empty historical records and bounded paths in the repocontract test;
retain its synthetic unmatched-path negative case. Do not expand governed-path
classification or treat all modern grants as members of that historical corpus.

Regeneration fixtures copy current Spec grants without assuming the legacy
folder exists. Update the dedicated-command fixture to a valid Spec-contained
grant. Keep existing deep equality, frozen-corpus and command-ownership assertions.
Temporary historical fixtures may materialize historical paths inside disposable
test repositories; the delivered checkout must retain the requested deletion.

## Coverage Map

- Core Feature 1 → Operative grant discovery → task_01, task_02.
- Core Feature 2 → Ownership-based regeneration → task_01, task_02.
- Core Feature 3 → Historical evidence and fixtures → task_01, task_02.
- Core Feature 4 → Historical evidence and fixtures → task_01, task_02.
- Goals 1–3 and OE-1 → all three design sections and the terminal QA Task.

## Vocabulary Contract

No glossary term is added, changed or retired. Spec, Task, Verification and
Baseline keep their existing meanings; the terminal QA Task confirms this.

## Build Order

Implement the five-file reader/fixture repair as one vertical slice, then run
independent terminal QA. Execution starts from main before the cleanup deletion,
where repository verification is green. Delivery subsequently applies the repair
to PR #180 and runs sanctioned regeneration and both repository gates there.
The executable dependency graph belongs only in `_tasks.md`.

## Research and decision evidence

Secondbrain was consulted on 2026-09-09 through its wiki index and qmd. The mirrored
Spec 0062 TechSpec documents strict validation after regeneration, and its historical
regeneration decision explains why guard and audit share declaration authority; current local
ADR-0149 supersedes it and governs this repair. These sources rule out bypassing
the guard or copying output inventories into grants.

Exa fetched the primary [Git show documentation](https://git-scm.com/docs/git-show)
and [Go io/fs documentation](https://pkg.go.dev/io/fs#WalkDir) on 2026-09-09.
Git documents direct blob-content reads; Go documents deterministic traversal and
explicit filesystem errors. These support immutable historical reads and bounded
filesystem discovery. Choosing the retained main ancestor and refusing unapproved
records are project decisions, not claims that the external sources validated
the implementation. Exa content was bounded; no repository code was sent externally.

## Existing QA execution boundary

The terminal Task renders the exact Verification derived by Roundfix; custom QA
Verification is rejected by the current loader/checker. That shell command tests
the report's terminal marker, not the truth of its rows. The qa-gate Agent must
independently perform and record the matrix, and the Supervisor must inspect that
evidence before delivery. Strengthening the general QA parser remains Spec 0122;
this repair neither changes it nor treats a pass marker alone as acceptance.

Secondbrain sources read: `projects/roundfix/mirror/docs/history/specs/0062-baseline-digest-regeneration-bootstrap/_techspec.md`
and `projects/roundfix/mirror/docs/history/adr/0128-the-guard-and-the-audit-read-one-regeneration-declaration.md`.
The advisory qmd query also ran on 2026-09-09; its broad results did not replace
these directly read historical sources or current local ADR-0149.
