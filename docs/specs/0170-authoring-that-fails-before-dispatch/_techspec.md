---
spec: 0170-authoring-that-fails-before-dispatch
status: active
created: 2026-09-25
surfaces: [backend, cli, docs]
---

# Authoring that fails before dispatch

## Executive Summary

Two new Task-stage detectors in `internal/speccheck` read what a pending Task
declares: one refuses an undeclared Governed Path, the other reports a CLI
change without its guide. Three existing detectors close honesty gaps, and the
templates and authoring skills stop producing the mistakes.

## Project Constraints

- Identifier strategy: applicable — `SC-TOOLING-UNDECLARED` and
  `SC-CLI-UNDOCUMENTED` are new stable codes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files only; no credential
  and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0083, ADR-0091, ADR-0092,
  ADR-0093, ADR-0094, ADR-0096, ADR-0097, ADR-0104, ADR-0116, ADR-0117,
  ADR-0130, ADR-0131, ADR-0149, ADR-0155 and ADR-0156 hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/speccheck/constraints.go`, `internal/speccheck/coherence.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `skills/write-techspec/references/techspec-template.md`,
  `.agents/skills/write-tasks/references/task-template.md`,
  `.agents/skills/write-tasks/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `skills/baseline_skill_contract_test.go`. Sanctioned regeneration:
  `make skills-sync`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Undeclared Governed Paths

`Check` in `internal/speccheck/constraints.go` already parses the Tooling
authority row of the PRD and TechSpec and selects the cited authorization
record. The row parser is extended so a row carries the backticked paths of its
`bounded files:` list (up to `Sanctioned regeneration:` or `Source:`), not only
the boolean `recordsBoundedFiles` returns today. A new detector in
`internal/speccheck/undeclared.go`, called from `Check` when a Task Graph is
present, collects for every Task that is neither completed nor `type: qa`:

- each Context entry of kind `interface` or `creates` (`instruction` names
  guidance and is not read);
- each Verification operand that resolves to an existing repository file, read
  through one function exported from `internal/spec/collision.go` that
  `declaredTaskTouches` also uses, plus any operand the Task declares under
  `creates:`.

Each collected path for which `GovernedPath` holds must appear in the
authorization record's `paths:` or among its sanctioned regeneration outputs
(resolved with `mechanicalRegenerationOutputs`, the audit's own predicate), and
in the `bounded files:` list of every present Tooling authority row unless it
is a regeneration output. A row path for which `GovernedPath` holds and that
the record does not grant is refused the same way. A Spec with no granted
record refuses every collected Governed Path. `CodeToolingUndeclared` is
declared in `constraints.go` beside the other tooling codes and registered in
`stagedDetectors` in `internal/speccheck/coherence.go` as a Task-stage
detector; a Spec without `_tasks.md` lists it as skipped.

## CLI changes without their guide

`CodeCLIUndocumented` is declared in `coherence.go` and registered as a
Task-stage detector. A new detector in `internal/speccheck/surface.go`, called
from `Check` beside the first, treats as a CLI surface a Context entry of kind
`interface` or `creates` naming a non-test `.go` file under `internal/cli/` or
`cmd/roundfix/`, or `internal/cli/cli_test.go`. A guide is an entry of the same
kinds under `.agents/skills/`, `skills/`, `docs/user-guide/` or `docs/agents/`.
A pending non-QA Task naming a surface, without a guide in itself or in any
Task in its transitive `needs`, is reported as a gap.

## Honesty gaps

- `internal/speccheck/verification.go`: `invertedExitVerificationCommand`
  gains the form of a command substitution holding a pipeline whose last member
  is `grep` and whose first member is not `printf` or `echo`, when the
  substitution, or a variable assigned from it, is tested with `test -z` or
  `[ -z`, and no `set -o pipefail` precedes it.
- `internal/speccheck/citations.go`: `parsePromiseSection` accepts one
  paragraph that begins with `None.` and carries a reason on the same or a
  following line; a blank-line-separated second block, a list or table line, or
  `None.` without a reason stays refused.
- `internal/speccheck/citations.go`: `detectReferenceIndex` reads each row's
  `source` cell and refuses with `SC-REF-UNRESOLVED` when that source is under
  `docs/findings/` or `docs/backlog/` and is still a file in the repository.
- `internal/speccheck/citations_test.go`: reason-only, evidence-only and
  whitespace-only-reason archived Findings each expect `SC-ARCHIVE-LICENSE`.

## Templates and skills

The Tooling authority rows of the PRD and TechSpec templates say
`bounded files:`. The Task template's Verification comment and the write-tasks
skill state the status-preserving form, that every edited path is declared
under `interface:` or `creates:` and never `instruction:`, that a declared
Governed Path must be authorized, and that a CLI change names its guide in the
same Task or a dependency. `skills/baseline_skill_contract_test.go` pins the
label.

## API Contracts

1. `SC-TOOLING-UNDECLARED` is an error naming the Task file and line, the
   Governed Path, and each declaration that omits it.
2. `SC-CLI-UNDOCUMENTED` is a gap naming the Task file and the CLI surface
   path; `--strict` promotes it to an error.
3. `SC-VERIFY-INVERTED-EXIT`, `SC-METRIC-UNDECLARED`, `SC-CONTRACT-UNDECLARED`,
   `SC-REF-UNRESOLVED` and `SC-ARCHIVE-LICENSE` keep their codes and severities;
   only the forms they accept or refuse change as described above.

## Vocabulary Contract

- emits: `internal/speccheck/constraints.go`
  pattern: `SC-TOOLING-UNDECLARED`
  documented-in: `.agents/skills/roundfix/SKILL.md`
- emits: `internal/speccheck/constraints.go`
  pattern: `SC-TOOLING-UNDECLARED`
  documented-in: `.agents/skills/qa-gate/SKILL.md`
- emits: `internal/speccheck/coherence.go`
  pattern: `SC-CLI-UNDOCUMENTED`
  documented-in: `.agents/skills/roundfix/SKILL.md`
- emits: `internal/speccheck/coherence.go`
  pattern: `SC-CLI-UNDOCUMENTED`
  documented-in: `.agents/skills/qa-gate/SKILL.md`

The glossary gains **Undeclared Governed Path** and **Undocumented CLI
Surface** in `CONTEXT.md`.

## Coverage Map

- Goal 1 → Undeclared Governed Paths; API Contract 1.
- Goal 2 → CLI changes without their guide; API Contract 2.
- Goal 3 → Templates and skills.
- Goal 4 → Honesty gaps; API Contract 3.
- Core Feature 1 → Undeclared Governed Paths.
- Core Feature 2 → CLI changes without their guide.
- Core Feature 3 → Honesty gaps.
- Core Feature 4 → Templates and skills.
- Success Metrics 1-2, 4 → Testing Approach 1, 5.
- Success Metric 3 → Testing Approach 2.
- Success Metric 5 → Testing Approach 3.
- Success Metric 6 → Testing Approach 4.
- API Contracts 1-3 → Undeclared Governed Paths, CLI changes without their
  guide, Honesty gaps.

## Integration Points

- **Corpus golden.** `internal/docscontract/corpus_test.go` characterizes
  every counted code; both new codes join `corpusFindingCodes` with count `0`
  in `testdata/corpus-golden.json` and in the pin in
  `internal/spec/archive_layout_characterization_test.go`.
- **Wave collisions.** `spec.Collisions` keeps its result; it reads
  Verification operands through the same exported function.

## Testing Approach

1. **Undeclared Governed Paths.** Fixture Specs through the public `Check`: an
   undeclared `interface:`, `creates:` and Verification path are each refused;
   the fully declared set passes; a path missing from one row, a row path the
   record lacks, and a Spec without authority are each refused; `instruction:`,
   completed Tasks and ordinary paths are not; a regeneration output counts as
   declared; a missing Task Graph lists the skip.
2. **CLI surface.** A surface without a guide is reported; a guide in the Task
   or a dependency passes; a guide only in a dependent Task does not;
   `internal/cli/cli_test.go` counts and an ordinary CLI test does not.
3. **Honesty gaps.** Both pipe-to-`grep` forms are refused and the
   status-preserving, skill-recommended and `pipefail` forms pass; a wrapped
   `None.` passes in both sections while reasonless and mixed sections fail;
   a left-behind Backlog Entry and Finding are refused and a single move
   passes; the three closure cases are refused.
4. **Templates.** A contract test pins `bounded files:` in both templates and
   the absence of `bounded paths:`.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification and `make verify-docs` results as facts.

## Build Order

1. Undeclared Governed Paths (depends on: none).
2. CLI changes without their guide (depends on: 1).
3. Honesty gaps (depends on: none).
4. Templates and skills (depends on: none).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A false refusal of a read-only path.** A Verification that only reads a
  Governed Path must now declare it; the write-tasks skill states this, and the
  rule is the one the Daemon's audit already enforces for an edit.
- **Two detectors editing one file.** Build Order 2 depends on 1 because both
  wire into `Check` and both change the corpus golden and the same two skills.
