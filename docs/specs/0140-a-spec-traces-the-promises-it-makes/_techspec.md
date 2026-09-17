---
spec: 0140-a-spec-traces-the-promises-it-makes
prd: _prd.md
created: 2026-09-17
---

# A Spec traces the promises it makes — Technical Spec

## Executive Summary

The Spec Consistency Check already carries a coverage mechanism: it parses
numbered units from the PRD, requires each in the TechSpec Coverage Map, and
requires each to be named in some Task's References. This design adds two unit
kinds to that mechanism — Success Metric, read from the PRD, and API Contract,
read from the TechSpec — and states the declaration form in the two authoring
skills that write those sections.

The trade-off this design accepts is severity. A section that is absent or
written as prose becomes a `gap`, not an `error`: the check cannot settle
whether a Spec has nothing to promise or forgot to say so, and the ten Specs in
authoring today would otherwise turn the repository's own documentation gate
red on delivery. A promise that is declared and then named by nobody is an
`error`, because both sides are located.

## Project Constraints

- Identifier strategy: applicable — the `SC-*` codes are stable and never renumbered once shipped; this design adds two codes, extends two existing codes to new unit kinds, and changes no existing code's number, stage or severity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the checker reads Spec artifacts from the local filesystem and opens no network connection, credential store or HTTP surface. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the detectors, their placement and their evidence are governed by accepted decisions this design preserves. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0093 applies: each finding names both compared sides and no coverage is inferred from resemblance.
  ADR-0117 applies: the declaration is checked at the stage that writes the artifact, and coverage in References at the first stage where a Task exists.
  ADR-0155 applies: the gate covers what the `qa` Task declares, which is why an untraced promise is now uncovered downstream.
  ADR-0104 applies: acceptance replays the ten Specs in authoring, which this Spec did not author.
  ADR-0156 applies: a declared promise names a consuming Task, and an explicit `None.` with a reason is a declaration.
- Tooling authority: applicable — the authoring rules live in Roundfix-owned Skills, so this Spec carries an operative grant. Express maintainer authorization: "Aprovar como proposto" on 2026-09-17, widened the same day to the Task authoring skill, and on the same date to the three paths an earlier authorization had already bounded, recorded in [_authorization.md](_authorization.md); bounded files: `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-techspec/references/techspec-template.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-techspec/SKILL.md`, `skills/write-techspec/references/techspec-template.md`, `skills/write-tasks/SKILL.md`, `internal/speccheck/coherence.go`, `internal/speccheck/constraints_characterization_test.go`, `internal/docscontract/testdata/corpus-golden.json`. The checker is ordinary source that no authorization has bounded. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Coverage unit set | `internal/speccheck` citation stage | Carry Success Metric and API Contract beside User Story and Core Feature. |
| Declaration reader | `internal/speccheck` citation stage | Read each promise section and classify it as declared units, an explicit none, or no declaration. |
| Reference vocabulary | `internal/speccheck` citation stage | Recognize the two new kinds where Task References and the Coverage Map already name units. |
| Stage table | `internal/speccheck` coherence stage table | Place each new code at the stage that can establish it. |
| Authoring rules | `write-prd`, `write-techspec` and `write-tasks` Skills and their templates | State the declaration form, its `None.` variant, the report gate, and the References obligation where Tasks are written. |
| Corpus expectation | `internal/docscontract` corpus golden | Characterize the codes the sweep may see and record the counts the new detectors produce. |

No new package, file or directory is proposed. Every seam above exists and is
already exercised by tests.

## Implementation Design

### Declaring a promise

A promise section is `## Success Metrics` in the PRD and `API Contracts` in the
TechSpec, at heading level two or three. Its content is read until the next
heading of the same or higher level. The reader classifies it as one of three
states:

- **Declared units** — one or more numbered items, read exactly as the existing
  parser reads user stories and Core Features: a leading integer, its number and
  its line.
- **Explicit none** — the section's only content is a first non-empty line that
  begins with `None.` and carries at least one further sentence on that line,
  which is the reason. A section that mixes `None.` with numbered items is not an
  explicit none: its numbered items are the declaration, and each carries the
  obligations below.
- **No declaration** — the section is absent, empty, or carries only unnumbered
  prose, bullets or a table.

The third state is the one the corpus is in today: of the ten Specs in
authoring, seven carry no Success Metrics section and three carry one written as
bullets or a table, and all ten carry an API Contracts section written as prose.

### Severity, and why it splits

- **No declaration** is a `gap`. The checker locates the artifact and the
  missing or unaddressable section, but cannot settle whether the Spec promises
  nothing or omitted the promise. `gap` is the severity the repository already
  reserves for a candidate the check cannot settle, and `--strict` still
  promotes it for a session that wants the stronger reading.
- **A declared unit nobody names** is an `error`. The declaration and the
  artifact that should have named it are both located, which is the existing
  `error` condition.

A finding carries the location of the artifact that declared the unit. Coverage
units are parsed from the PRD today and their findings are rendered against the
PRD, so the API Contract kind, which is declared in the TechSpec, must carry its
own declaring artifact through to the finding. A contract reported at a PRD line
would send the author to the wrong file.

### Codes and their stages

Two codes are added, one per stage, because the stage table binds one stage to
one code:

- `SC-METRIC-UNDECLARED`, at the PRD stage, for the PRD's own section.
- `SC-CONTRACT-UNDECLARED`, at the TechSpec stage, for the TechSpec's own
  section.

Two existing codes gain the new kinds and keep their number, stage, severity and
message shape:

- `SC-COVERAGE-UNMAPPED` (TechSpec stage) additionally requires each declared
  Success Metric in the Coverage Map. A declared API Contract is not required
  there: it already lives in the TechSpec, and mapping a section to itself
  proves nothing. A Spec with no TechSpec keeps the detector's existing
  behavior, which is to skip for a missing artifact; its metrics are still
  required in Task References, and no contract rule applies to it.
- `SC-COVERAGE-UNTASKED` (Task Graph stage) additionally requires each declared
  Success Metric and each declared API Contract to be named in some Task's
  References.

### Interfaces

The public surface is the existing checker API and its `SC-*` result codes. No
exported function signature changes; the unit kind enumeration and the reference
patterns grow.

```go
// Existing kinds, extended. Values are the words a reader writes.
const (
    coverageFeature  coverageKind = "Core Feature"
    coverageStory    coverageKind = "User Story"
    coverageMetric   coverageKind = "Success Metric"
    coverageContract coverageKind = "API Contract"
)
```

References are matched the way they already are, by a case-insensitive phrase
followed by a number sequence, so `Success Metrics 1-2` and `API Contract 3`
resolve without a new grammar.

### Data Models

No entity, schema or stored record changes, and the checker stays read-only over
Spec artifacts.

The corpus golden keeps its current shape: one aggregate map from finding code
to count over the whole active corpus. The codes the sweep accepts come from a
separate list compiled into the corpus contract, not from the golden's keys: an
emitted code absent from that list fails the sweep as uncharacterized before any
count is compared. Adding a code therefore means adding it in both places — the
accepted-code list and the golden map.

### API Contracts

1. `SC-METRIC-UNDECLARED` — emitted at the PRD stage with severity `gap`. It
   names the PRD and the absent or unaddressable Success Metrics section, and
   its fix line states the two accepted forms.
2. `SC-CONTRACT-UNDECLARED` — emitted at the TechSpec stage with severity `gap`,
   naming the TechSpec and its API Contracts section under the same rule.
3. `SC-COVERAGE-UNMAPPED` and `SC-COVERAGE-UNTASKED` keep their codes, stages,
   severities and message shapes, and report the new unit kinds by the same
   words a reader writes: `Success Metric 2`, `API Contract 1`.

No command, flag, exit code or output format changes. A `gap` still leaves the
command's exit status unchanged unless `--strict` is passed.

## Coverage Map

- Goal 1 → Declaration reader, Coverage unit set, Reference vocabulary.
- Goal 2 → Declaration reader (explicit none), Authoring rules.
- Goal 3 → Severity split; Declaration reader.
- Goal 4 → Corpus expectation (Testing Approach 4).
- User Story 1 → Declaration reader, Stage table, Authoring rules.
- User Story 2 → Reference vocabulary, Coverage unit set.
- User Story 3 → Declaration reader (explicit none).
- User Story 4 → Codes and their stages; API Contracts 1-3.
- Core Feature 1 → Declaration reader; Authoring rules (`write-prd`,
  `write-techspec`).
- Core Feature 2 → Coverage unit set; Reference vocabulary; Authoring rules
  (`write-tasks`).
- Core Feature 3 → Severity split.
- Core Feature 4 → Stage table.
- Core Feature 5 → Coverage unit set (existing kinds untouched); Testing
  Approach 2.
- Success Metric 1 → Corpus expectation (Testing Approach 3 and 4).
- Success Metric 2 → Authoring rules; Declaration reader.
- Success Metric 3 → Testing Approach 2.
- API Contracts 1-3 → Codes and their stages; Vocabulary Contract. The rule
  does not require a contract in the Coverage Map; this line records where the
  three land, for the reader.

## Integration Points

- **Documentation gate.** `make verify-docs` runs the checker over the active
  corpus without `--strict`, so a new `gap` does not fail it. The
  corpus golden, one aggregate map from finding code to count over the whole
  active corpus, must be updated deliberately in the same delivery, with the two
  new codes added as keys.
- **Skill distribution.** The canonical skills regenerate their mirrors through
  the repository's skill-sync target; the mirror parity test and the owned-skill
  version check keep the copies identical. No skill's version moves, so no
  setup minimum moves.
- **Spec 0129.** This Spec is the first slice carved from it. The slices that
  remain — inherited ADR obligations, property-shaped acceptance, amendment,
  stale-QA recovery, temporal prerequisites — keep their place there and are not
  touched here.
- **Archived Specs.** They are never swept and never edited. The replay reads
  only the active corpus.

## Testing Approach

1. **Declaration reader, at the existing fixture seam.** Unit tests in the
   checker's citation corpus cover: a numbered declaration that some Task names
   (silent); a numbered declaration nobody names (`error`); an absent section
   (`gap`); `None.` with a reason (silent); `None.` with no reason (`gap`);
   bullets, a table and unnumbered prose (`gap` each). Each new case fails on
   the tree as it stands today.
2. **Existing kinds unchanged.** The stage table test asserts every existing
   code keeps its stage, and the existing coverage tests keep passing unedited,
   which is what proves user stories and Core Features kept their behavior.
3. **Corpus expectation.** The active-corpus sweep produces exactly the
   predicted counts: ten `SC-METRIC-UNDECLARED`, ten `SC-CONTRACT-UNDECLARED`,
   and no new `error`. The golden is updated with those counts in the same Task
   that makes them true.
4. **Outside evidence.** The ten Specs in authoring were written before this one
   and by a different intent. The replay states the prediction before running
   it — seven PRDs with no Success Metrics section, three with an unaddressable
   one, ten TechSpecs with prose API Contracts — and the sweep either matches it
   or the difference is reported. Where a Spec cannot be swept, the row records
   that reason and does not block.
5. **Skill text.** Task Verification asserts, in each canonical copy and its
   mirror, that the new clauses are present — the declaration form in the PRD
   and TechSpec skills, the mapping and References obligation in the Task skill,
   and the References example in the Task template — and that the mirror parity
   test passes.
6. **Self-application.** This Spec declares numbered Success Metrics and
   numbered API Contracts, maps the metrics in its own Coverage Map, and names
   all of them in its own Task References. Its own check must be clean under the
   rule it ships.

## Build Order

1. Declaration reader, the two new codes and their stage placement, with the
   unit tests above (depends on: none).
2. Coverage extension: the two new kinds in the Coverage Map requirement and in
   Task References, with unit tests (depends on: 1).
3. Glossary owners for the two coined codes, the two codes added to the corpus
   contract's accepted-code list, and the golden map updated to the counts steps
   1 and 2 produce, with the prediction recorded (depends on: 1, 2).
4. Authoring rules in the three skills and their templates, canonical then
   mirror: the declaration form in `write-prd` and `write-techspec`, and the
   mapping and References obligation in `write-tasks` (depends on: 1, 2, 3).
5. Terminal QA (depends on: 1, 2, 3, 4).

Step 4 follows the behavior steps deliberately: a documentation Task written
against a rule that is still moving is a Task written against a draft.

## Risks & Considerations

- **Turning a rule on over an existing corpus.** Mitigated by the severity
  split and by updating the golden in the same delivery. If the predicted counts
  are wrong, step 3 fails loudly rather than silently accepting a new number.
- **A gap blocks decomposition.** The Task authoring skill treats a `[gap]` as
  blocking, so each of the ten Specs in authoring must declare its promises
  before it can be decomposed. The maintainer chose that cost over converting
  twenty artifacts the coming slices will replace; this Spec converts none of
  them and leaves the twenty gaps as its measured evidence.
- **`None.` as an escape hatch.** A Spec can declare `None.` with a thin reason
  and satisfy the rule. Mechanically requiring honesty is not possible here;
  independent review, which Spec 0126 makes configurable, is where a thin reason
  is challenged. This limit is accepted, not hidden.
- **The templates carry a label the checker rejects.** Both templates tell the
  author to write `bounded paths:` in the Tooling authority row, while the
  checker accepts only `bounded files:`, `bounded proposed files:` or
  `bounded repository-relative`. This Spec does not repair that, to keep the
  slice to one outcome; it is captured as a pending Inbox Entry for the
  repository.
- **Reference phrases could collide.** `API Contract` appears in prose as well
  as in References. The existing matcher requires a number sequence after the
  phrase, which is what keeps a sentence about API contracts from registering as
  a reference.

## Decisions

- **Two codes, one per stage.** The stage table binds one stage per code, and a
  PRD section and a TechSpec section are established at different stages.
- **Extend the existing coverage codes rather than mint parallel ones.** The
  obligation is identical to the one user stories already carry, so a new code
  would split one rule across two vocabularies.
- **A numbered list is the declaration form.** A unit must be addressable as
  `Success Metric 2` for a Task to name it; bullets and tables are not. The
  corpus conversion cost lands on Specs still in authoring.
- **An explicit `None.` with a reason is a declaration.** See ADR-0156.
- **Severity splits between absence and contradiction.** See the Severity
  section; this is what keeps the documentation gate green on delivery.
- **The API Contract unit is not required in the Coverage Map.** It is already a
  TechSpec section, and requiring the TechSpec to map itself proves nothing.

## Vocabulary Contract

- emits: `internal/speccheck/citations.go`
  pattern: `SC-METRIC-UNDECLARED`
  documented-in: `CONTEXT.md`
- emits: `internal/speccheck/citations.go`
  pattern: `SC-CONTRACT-UNDECLARED`
  documented-in: `CONTEXT.md`

Both codes are coined here, so each needs a glossary owner before it ships, and
the declaration is what makes `SC-VOCABULARY-UNDOCUMENTED` run over them instead
of skipping. The patterns do not match the glossary today; Build Order step 3
makes them match.

No other token is coined. "Promise" is used descriptively for the two declared
unit kinds and names no new artifact, field, state or command.
