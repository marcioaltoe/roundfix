---
spec: 0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written
prd: _prd.md
created: 2026-11-09
---

# The strict check reads Unreachable Acceptance and related ADRs as written — Technical Spec

## Executive Summary

The declaration reader in `internal/spec` becomes one walk of the PRD's
`## Unreachable Acceptance` section that returns, beside the declarations it
returns today, every non-empty line it does not read. A new detector in
`internal/speccheck`, `SC-UNREACHABLE-UNREAD`, reports those lines and a
declaration without one of its fields, from the PRD stage on, as an error.
`detectADRConsistency` seeds the related-ADR candidates with the listed ADRs
minus those whose TechSpec `## ADR Alternatives` line begins with
`not applicable`, as ADR-0377 decides, and the guides, the glossary, the two
amended ADRs and three owned skills say so. The trade-off this design accepts
is that the QA gate and `roundfix archive` keep reading a section that holds
unread text without refusing it, so a Spec already inside a Run is never
stopped by the new rule; the strict check, which runs before every Run and
at queue revalidation, is where the text is named.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the new finding code joins the existing `SC-*` vocabulary.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local files and Git history
  only. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0377 (new) records the rule,
  ADR-0377: "Only an ADR the Spec applies or changes seeds the related-ADR
  candidates." It amends ADR-0093, ADR-0093: "an ADR that cites an ADR the
  Spec already lists is a candidate the Spec should account for", and
  ADR-0168, ADR-0168: "reports an accepted ADR that cites an ADR a Spec
  lists"; task_02 writes both notices. ADR-0265 defines the ADR Alternatives
  line, ADR-0265: "The check reads structure, never meaning." ADR-0094 keeps
  a detector quiet when its artifact is missing, ADR-0094: "is skipped, not
  failed, when one of them is absent". ADR-0343 matches each declared-blocked
  row with its declaration, ADR-0343: "a declared row matches the criterion
  it names". ADR-0208 lets one Spec adopt both sources, ADR-0208: "One Spec
  may adopt several Inbox Entries, Backlog Entries and Findings whose context
  is similar or complementary". Not applicable and marked so below: ADR-0210
  and ADR-0211, cited by the Fiscus reference in Fiscus's numbering, and
  ADR-0156, ADR-0176, ADR-0183, ADR-0184, ADR-0209, ADR-0220, ADR-0274,
  ADR-0276, ADR-0280, ADR-0295, ADR-0319, ADR-0335 and ADR-0353.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BA answers of 2026-11-09 quoted in the authorization record; no
  Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0367-the-strict-check-reads-unreachable-acceptance-and-related-adrs-as-written/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-techspec/references/techspec-template.md`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`,
  `skills/write-techspec/SKILL.md`,
  `skills/write-techspec/references/techspec-template.md`.

## System Architecture

- `internal/spec/spec.go` (task_01): `parseUnreachableDeclarations` and the
  new exported `UnreachableUnread` share one unexported walk of the section.
  `Unreachable` keeps its signature and returns the same declarations and
  the same `UnreachableDeclarationError` as at `2e59aebe` for every input.
- `internal/speccheck/unreachable.go` (new, task_01): `CodeUnreachableUnread`
  and `detectUnreachableDeclarations`, called by `Check` (`constraints.go`)
  right after the PRD is found present and by `checkAuthoringStage`
  (`coherence.go`) before the contract horizon; both add a skip naming the
  PRD path when the PRD is absent. `stagedDetectors` lists the code at
  `StagePRD`.
- `internal/speccheck/citations.go` and `authoring_proof.go` (task_02):
  `detectADRConsistency` builds the seed set with a new
  `appliedADRNumbers(listed, techSpecPath)`, which reads the TechSpec's
  `## ADR Alternatives` lines through the existing `adrAlternativePattern`
  and drops each number whose relation matches the new
  `notApplicableRelationPattern`, `(?i)^not applicable\b`. `listed` itself is
  unchanged, so a marked ADR is still listed for `SC-ADR-UNLISTED`, is never
  reported by `SC-ADR-RELATED`, and still needs its line under
  `SC-ADR-ALTERNATIVES-UNRELATED`.
- Documentation: ADR-0377 (in the tree with this plan), the notices in
  ADR-0093 and ADR-0168 and the two `CONTEXT.md` entries (task_01 and
  task_02), the two user guides and three owned skills (task_03).

## Implementation Design

### Interfaces

```go
// internal/spec
type UnreachableUnreadLine struct {
	Line int    // 1-based line number in the PRD
	Text string // the line, surrounding whitespace trimmed
}

// UnreachableUnread reads the PRD content's ## Unreachable Acceptance
// section as Unreachable does and returns, in file order, every non-empty
// line no declaration reads, with the error Unreachable returns for a
// declaration without one of its fields.
func UnreachableUnread(content []byte, prdPath string) ([]UnreachableUnreadLine, error)

// internal/speccheck
const CodeUnreachableUnread = "SC-UNREACHABLE-UNREAD"
func detectUnreachableDeclarations(result *Result, repoRoot, prdPath string) error
func appliedADRNumbers(listed map[string]bool, techSpecPath string) (map[string]bool, error)
```

The reading rule, which `UnreachableUnread` implements and the guides state.
Inside the section a line is read when it is blank; inside an HTML comment,
from a line whose trimmed text starts with `<!--` through the line that
holds `-->`; a top-level `- ` line, which opens a declaration; a
`criterion:`, `reason:` or `satisfied-by:` line inside an open declaration;
or a line indented by a space or a tab that continues the field last read.
Before the first declaration, the paragraph whose first line is the
section's first non-blank, non-comment line and begins with `None.` is read
when the section holds no declaration. Every other non-empty line is unread.
The section ends at the next `# ` or `## ` heading, as today.

### Data Models

No changes. No file, record or schema is written by either check.

### API Contracts

1. API Contract: `SC-UNREACHABLE-UNREAD`, severity error, from the PRD
   stage. One finding per run of consecutive unread lines, located at the
   run's first line, with the summary
   `<prd> line <N> under ## Unreachable Acceptance is not part of a fielded declaration, so the QA gate does not read it: <quoted first line>`,
   the line quoted as Go's `strconv.Quote` writes it. Then, when a
   declaration lacks a field, one finding at its line with the summary
   `<prd> line <N> starts an Unreachable Acceptance declaration without its <field>: field`.
   Both carry the fix
   `Write each entry in the fielded form: a top-level "- criterion: <text the QA row repeats>" line, then indented "reason: <why no Verification or QA row can reach it>" and "satisfied-by: <evidence or follow-up that satisfies it>" lines; start the section with None. when no criterion is unreachable.`
   A missing PRD skips the code with the PRD path; an absent section reports
   and skips nothing.
2. API Contract: `SC-ADR-RELATED` keeps its summary, locations and gap
   severity. Its fix becomes
   `Account for ADR-<N> in the Active ADR obligations row in <prd> by listing it or recording why it does not apply; when it does not apply, begin its ADR Alternatives line in the TechSpec with "not applicable", so the ADRs that cite only it are not drawn in.`
   An accepted unlisted ADR is reported only when it cites a listed ADR whose
   TechSpec line does not begin with `not applicable`; the cited ADR the
   summary names is the lowest such number.

### Documentation

`CONTEXT.md` appends one sentence to the definition of **Unreachable
Acceptance** (task_01) and one to **Spec Consistency Check** (task_02);
ADR-0093 and ADR-0168 gain the amendment notice naming ADR-0377 (task_02).
task_03 adds the `SC-UNREACHABLE-UNREAD` row and the marker sentence to the
`spec` command guide, a paragraph to the context-driven development guide,
the same two facts to the Roundfix skill's `spec` reference, the marker to
the `write-techspec` skill's Rejected alternatives rule and template comment,
and the marker to the `write-tasks` skill's gap sentence; it raises the three
skills' versions by the record command and re-records the Behavior Surface
Record last.

### Surface Transcripts

None. `roundfix spec check` keeps its report layout; the new finding and the
changed fix are lines of the existing layout, asserted by the tests below and
by the QA gate on the built binary.

## Coverage Map

- Goal 1, Story 1, Core Feature 1, Success Metric 1 → Interfaces
  (`UnreachableUnread`, `detectUnreachableDeclarations`); API Contract 1;
  Build Order 1.
- Goal 2, Success Metric 2 → Interfaces (the reading rule); API Contract 1;
  Build Order 1.
- Goal 3, Story 2, Core Feature 2, Success Metric 3 → Interfaces
  (`appliedADRNumbers`); API Contract 2; Build Order 2.
- Goal 4, Story 3, Core Feature 3 → Documentation; Build Orders 1-3.
- Success Metric 4 → Build Order 4 (the QA gate runs the delivered check on
  this Spec).

Deviations from the PRD: None.

## Absence Evidence

1. The unread-line reader — command: `grep -rn 'Unread' internal/spec/spec.go internal/speccheck/*.go`; result: only `CodeCoverageUnread` and its detector, which read the TechSpec Coverage Map; nothing reads the Unreachable Acceptance section outside `parseUnreachableDeclarations`.
2. A strict-check reading of Unreachable Acceptance — command: `grep -rn 'Unreachable' internal/speccheck --include='*.go'`; result: no match (exit 1).
3. A marker that stops an ADR from seeding the related candidates — command: `grep -rn -i 'not applicable' internal/speccheck --include='*.go'`; result: only `notApplicable applicability = "not applicable"` in `constraints.go`, the Project Constraints row token; `detectADRConsistency` passes every listed number to `firstListedCitation`.

## ADR Alternatives

- ADR-0377: the decision this Spec records; it rejects a marker inside the
  PRD row, a fifth Project Constraints row, a limit on rounds and removing
  `SC-ADR-RELATED`, and the design takes none of them.
- ADR-0093: amended; its candidate rule keeps one citation edge and reads only
  written text, and it rejected judging relevance, which the fixed token does
  not do.
- ADR-0168: amended; its horizon is unchanged, and the design takes none of
  its rejected options.
- ADR-0265: kept; it rejected judging whether a relation is true, and the
  check reads only the line's leading token.
- ADR-0094: kept; it rejected demanding every artifact, and the new detector
  skips a missing PRD while the seeds fall back to every listed ADR without
  a TechSpec.
- ADR-0343: kept; it rejected matching by count alone, and the new detector
  only makes sure each declaration it matches is read.
- ADR-0208: kept; both adopted sources change the same detector family, and
  the design takes none of its rejected options.
- ADR-0209: not applicable; the judge's grouping suggestion is unchanged.
- ADR-0210: not applicable; the adopted Fiscus reference names a Fiscus ADR
  of that number, and this repository's evidence snapshot is unchanged.
- ADR-0211: not applicable; the adopted Fiscus reference names a Fiscus ADR
  of that number, and the agent environment is unchanged.
- ADR-0220: not applicable; Doctor is unchanged.
- ADR-0319: not applicable; the follow-up order is unchanged.
- ADR-0353: not applicable; the follow-up list is unchanged.
- ADR-0156: not applicable; it concerns Success Metric and API Contract
  declarations, which this Spec does not change.
- ADR-0176: not applicable; the check still reads only authored artifacts,
  and the TechSpec it now reads for the marker is one of them.
- ADR-0183: not applicable; Claim Receipts are unchanged.
- ADR-0184: not applicable; Surface Transcripts are unchanged, and this
  Spec declares none.
- ADR-0274: not applicable; it concerns when an Agent turn counts as
  finished, which this Spec does not change.
- ADR-0276: not applicable; the Network-Prohibited Row and the Pinned Byte
  Invariant are unchanged.
- ADR-0280: not applicable; the guide fingerprint, the decision record and
  the glossary's own rule are unchanged, and task_03 follows the first.
- ADR-0295: not applicable; it concerns how reachability reads the default
  branch, which this Spec does not read.
- ADR-0335: not applicable; it concerns how a Run counts a stop, which this
  Spec does not change.

## Integration Points

None outside the repository. The QA gate (`qa-report accept`), `settle` and
`roundfix archive` keep calling `spec.Unreachable`; the Delivery Queue's
revalidation and the QA gate's precondition run `speccheck.Check`, so they
report the new error and the narrower `SC-ADR-RELATED` without a change of
their own.

## Testing Approach

- `internal/spec/unreachable_unread_test.go` (new, task_01):
  `TestUnreachableUnreadReadsTheSectionAsUnreachableDoes` runs
  `UnreachableUnread` and `Unreachable` over the same inputs and asserts the
  unread line numbers and that the declarations `Unreachable` returns are the
  ones it returned before for the existing fixtures under
  `internal/spec/testdata/unreachable`.
- `internal/speccheck/unreachable_test.go` (new, task_01):
  `TestAnUnreachableSectionReadAsWrittenReportsNothing` (fielded, two
  criteria, `None.` paragraph, empty, absent, the `write-prd` template
  comment), `TestANumberedUnreachableEntryIsReportedAtItsFirstLine` (Spec
  0364's shape, at `StageAll` and `StagePRD`),
  `TestProseAfterADeclarationAndAMissingFieldAreReported` and
  `TestNoneBesideADeclarationAndAPlainBulletAreReported`.
- `internal/speccheck/adr_related_seed_test.go` (new, task_02):
  `TestANotApplicableADRDrawsInNoCitingADR` (lower and upper case),
  `TestAnAppliedADRStillDrawsInItsCitingADR` (applied, and a relation that
  only contains the words later), `TestWithoutATechSpecEveryListedADRSeeds`
  and `TestAnADRCitingAnAppliedADRIsStillRelated`.
- Fixtures are temporary directories written by the tests; no test reads a
  Spec folder of this repository, starts a process or reaches the network.
  Every new top-level test calls `t.Parallel()` first.
- A prototype of both changes on `2e59aebe` passed these tests and
  `go test ./internal/speccheck ./internal/spec`,
  `-tags docscontract ./internal/docscontract` and the `SpecCheck` tests of
  `./internal/cli`; each sabotage the QA Task lists failed its named test.

## Build Order

1. The unread reader, `SC-UNREACHABLE-UNREAD`, its tests and the
   **Unreachable Acceptance** glossary sentence.
2. The applied-ADR seeds, the new fix, their tests, the two amendment notices
   and the **Spec Consistency Check** glossary sentence (depends on: 1, which
   also edits `CONTEXT.md`).
3. The two user guides, the Roundfix `spec` reference, the `write-techspec`
   and `write-tasks` skills, their version record and the Behavior Surface
   Record (depends on: 2, so the guides describe delivered behavior).
4. Final QA gate (depends on: 3).

No step needs another Spec's code. Specs 0365 and 0366 merge first and may
move `CONTEXT.md`, `docs/references/behavior-surfaces.json` and
`skills/testdata/owned-skill-versions.json`; task_03 regenerates the two
records by their commands on its own base.

## Risks & Considerations

- An adopter Spec already in the queue whose section holds unread prose now
  stops at revalidation with `SC-UNREACHABLE-UNREAD` instead of failing at its
  QA gate; the fix line shows the form. The survey of 81 archived PRDs found
  none that the rule reports.
- An adopter whose TechSpec writes the relation in another language
  (`não se aplica`) keeps today's behavior until it writes `not applicable`;
  the guides and skills name the token.
- A TechSpec line that marks `not applicable` an ADR the PRD row describes as
  applied is read as not applying; the check reads the declaration, not the
  prose (ADR-0093).
- The `write-techspec` and `write-tasks` skills reach adopters with the next
  release through the Baseline's skill install; their version fields rise
  by the record command only.

## Glossary

- changes: **Unreachable Acceptance** — the Spec Consistency Check reports text in the section that no fielded declaration reads.
- changes: **Spec Consistency Check** — only an ADR the Spec applies or changes draws in the ADRs that cite it.

## Decisions

- The marker is the leading `not applicable` of the ADR's TechSpec line. See
  ADR-0377.
- One reader in `internal/spec` serves the QA gate and the check; the parser
  does not refuse unread text.
- `SC-UNREACHABLE-UNREAD` is an error from the PRD stage; consecutive unread
  lines are one finding.
- Three implementation Tasks and the QA gate, in series.
