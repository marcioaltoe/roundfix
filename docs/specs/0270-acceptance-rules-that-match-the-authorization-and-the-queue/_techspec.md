---
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
prd: _prd.md
created: 2026-10-13
---

# Acceptance rules that match the authorization and the queue — Technical Spec

## Executive Summary

The shared QA policy in `internal/spec/qa.go` gains a third exempt row kind,
the Network-Prohibited Row, counted from the Results table the way the
Network-Denied Row already is. Every caller reads the policy through
`QAReportEligibility` and `EnvironmentRowsNeedingOverride`, so settlement,
archive, `qa-report accept` and the Delivery Queue change together. A new
structural detector in `internal/speccheck` reports `SC-QA-BYTES-PINNED` for a
pending QA requirement that pins a byte invariant to a commit. The Baseline
clauses, the owned skills, the user guide and the glossary state both rules.

The trade-off accepted is ADR-0276's trust in the text the gate writes: the
policy does not read `_authorization.md` to confirm the prohibition, as
ADR-0240 does not confirm a denied host. In exchange no authorization record
changes shape, and the claim stays checkable afterwards. An authoring
prototype at `9963ffd3` measured the cost: the policy change alone leaves
every existing test green, the detector leaves the speccheck, cli, daemon,
delivery, judge and skills suites green, and `make verify-docs` and
`make verify-contracts` pass with every text and regeneration change applied.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes; the code
  follows the existing `SC-` naming. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added, and the detector executes no authored command. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0276 governs every decision below.
  ADR-0240: "The policy has one rule, in `internal/spec`, which every caller
  reads". ADR-0104: "a row whose external evidence cannot be obtained is
  recorded as blocked with its reason rather than dropped". ADR-0167: "That
  row no longer counts against a qualifying `partial`". ADR-0093: "it never
  judges whether a decision is correct". ADR-0094: "is skipped, not failed,
  when one of them is absent". ADR-0265: "Every active Spec is held, with no
  horizon". ADR-0148: "a command that exits zero before the work exists
  proves nothing about the work". ADR-0259: "A test runs in parallel unless it
  names why it cannot". ADR-0271: "A test bounds by real time only what it
  asserts". ADR-0244: "A Spec now carries a Glossary Declaration". ADR-0256:
  "Spec authoring requires a skills Task for a change a skill describes".
  ADR-0156: "with the reason none applies". ADR-0168 and ADR-0183 chose
  horizons this Spec does not take. ADR-0176: "The citation walk now reads
  only the authored projection". ADR-0229 parks a partial with a non-exempt
  environment row, as ADR-0240 narrowed it. ADR-0237 does not apply; no
  Delivery Retry rule changes. ADR-0236: "Every attempt,
  including a re-run, tests the default branch as it stands when the job
  starts".
  ADR-0184 asks for Surface Transcripts, declared `None.` below. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable. Express maintainer authorization: "considere
  autorizado a ajustar todas as skills se necessário", "Autorizar os dois",
  the 2026-10-11 grant "Concedo" and the 2026-10-13 cycle AG approval. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.
  Spec-contained authorization record:
  `docs/specs/0270-acceptance-rules-that-match-the-authorization-and-the-queue/_authorization.md`.
  Bounded files: `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/settle.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-tasks/SKILL.md`, `docs/agents/autonomous-work.md`,
  `docs/agents/setup-context.json`, `docs/agents/spec-routing.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`,
  `internal/baseline/assets/modules/autonomous-work.json`,
  `internal/baseline/assets/modules/spec-workflow.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `skills/archive-spec/SKILL.md`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-tasks/SKILL.md`.

## Current behavior

- `qaReportResultsCounts` counts two exempt kinds: the pre-PR Pull Request
  row and a row whose status starts `blocked (environment: network denied: `
  with a host and whose provenance names `outside-evidence row`.
  `EnvironmentRowsNeedingOverride` subtracts both from
  `rows_blocked_environment`, and `deliver_workflow.go` parks
  `qa-environment-partial` when the remainder is positive.
- The archived QA Report of Spec 0258 at `078869bf` records five environment
  rows: four `blocked (environment: network prohibited by Spec authorization)`
  outside-evidence rows and the Pull Request row. The `bin/roundfix` built at
  `9963ffd3` refuses it: `rows_blocked_environment is 5, 4 outside the pre-PR
  Pull Request row; expected 0 outside it`, exit 1.
- The qa-gate skill tells the gate to write the network-denied form only when
  the Run sandbox denied the lookup; it names no form for a lookup the
  authorization forbids. The QA settlement table in three skills, the roundfix
  archive and settle references, four user guide pages and the two Baseline
  clauses list only the two exempt kinds.
- No detector reads QA requirements for a commit. The QA Task of Spec 0264 at
  the parent of `4af8655e` held a file "byte-identical to its bytes at" a
  named commit, and `spec check --strict` reported nothing.

## System Architecture

The policy change stays inside `internal/spec/qa.go`, the one place ADR-0240
names. The detector is a new file beside the other per-Task detectors in
`internal/speccheck` and runs from the Task loop of
`detectTaskCoverageAndContextReferences`, which already skips completed
Tasks. `stagedDetectors` registers the code at the Tasks stage, so the PRD and
TechSpec stages list it as skipped, and a Spec without a Task Graph lists it
as skipped. No new package, command, flag or output format.

## Implementation Design

### Interfaces

```go
// internal/spec
const QANetworkProhibitedStatus = "blocked (environment: network prohibited by Spec authorization)"

type QAReport struct {
	// ... existing fields
	RowsBlockedNetworkProhibited int
}

// internal/speccheck
const CodeQABytesPinned = "SC-QA-BYTES-PINNED"

// PinnedByteInvariants reports each requirement of a pending QA Task that
// names a byte-identity phrase and a commit-shaped token together.
func PinnedByteInvariants(task spec.Task, isQA bool) []Finding
```

### Data Models

- Network-Prohibited Row: a Results-table row whose Status cell, trimmed, is
  exactly `QANetworkProhibitedStatus` and whose Provenance cell names
  `QAOutsideEvidenceRowSource` under `qaReportProvenanceNames`. Status and
  source are case-sensitive, and no suffix is accepted. The frontmatter
  counts keep their meaning; the new count is derived from the rows.
- Pinned Byte Invariant: a requirement item of `## Requirements`, its
  continuation lines folded as `spec.Task.Requirements` holds them, of a Task
  that is the manifest's `qa:` Task or has type `qa` and is not completed. It
  matches a byte-identity phrase, case-insensitive: `byte-identical`,
  `byte identical`, `byte identity`, `byte-for-byte`, `byte-equal`,
  `identical bytes`, `same bytes` or `unchanged bytes`. It also holds a
  commit-shaped token: split the text at every character outside
  `[0-9A-Za-z_./-]`, trim trailing dots, split at `..`, and keep a part of 7 to
  40 characters of `[0-9a-f]` with at least one digit and one letter. A path
  such as `docs/abc1234`, `ADR-0276`, a date such as `20261013` and a 64-hex
  digest are not commits.

### API Contracts

1. API Contract: a Network-Prohibited Row counts in
   `RowsBlockedNetworkProhibited`. `EnvironmentRowsNeedingOverride` subtracts
   it with the other two exempt kinds, capped by `rows_blocked_environment`.
   `QAReportEligibility` treats it as an unmet exempt row, so a partial whose
   only unmet rows are such rows, Pull Request rows, Network-Denied Rows and
   covered declarations qualifies. When a non-exempt environment row stands
   beside at least one Network-Prohibited Row, the refusal reads
   `rows_blocked_environment is <n>, <m> outside the pre-PR Pull Request row and the network-denied and network-prohibited outside-evidence rows; expected 0 outside them`.
   Every existing refusal text is unchanged when the report has no such row.
2. API Contract: `SC-QA-BYTES-PINNED` is a gap from the Tasks stage, reported
   once per matching requirement:
   `<task path> line <N> pins a byte invariant to commit <token>`, located at
   the requirement's first line. Its fix reads: `Compare with the Delivery
   Base, the merge base of the audited head with the default branch, and keep
   the commit where the plan measured out of the requirement, following the
   write-tasks rule "Byte invariants name the delivery base".` It holds every
   active Spec with no horizon, and `--strict` turns it into a failing result.
3. API Contract: the code joins `corpusFindingCodes`, the corpus golden and the
   characterization that pins it at 0. Both update texts gain the identical
   sentence `SC-QA-BYTES-PINNED joins the characterized code set at 0, with
   every earlier count unchanged.` The text and JSON formats, exit codes and
   flags do not change.
4. API Contract: the skill, guide, Baseline and glossary texts of the next
   section appear verbatim, after whitespace normalization, in the files
   named, and the `### QA settlement` section stays byte-identical in all
   three skills at the delivery base.

### Surface Transcripts

None. No command, flag, output line format or exit code changes; a new
refusal reason and the new code render through the existing lines, and the
tests of task_01 and task_02 assert them.

### Exact texts

`qa-gate`, `archive-spec` and `roundfix` gain this section right after their
QA settlement table, before the next heading:

```markdown
### Rows the Spec's authorization keeps offline

The qualifying declared `partial` above also accepts an outside-evidence row
the gate did not look up because the Spec's own authorization forbids network
access, recorded as `blocked (environment: network prohibited by Spec authorization)`
with the outside-evidence row named in its provenance. Like the network-denied
row, it never decides a qualifying partial and needs no Unreachable Acceptance
declaration, and settlement, `roundfix qa-report accept` and `roundfix archive`
read it through the same policy (ADR-0276).
```

In `qa-gate`'s outside-evidence paragraph, after "never decides a qualifying
partial.", insert: "When the Spec's own `_authorization.md` forbids network
access, do not attempt the lookup: record the row exactly as
`blocked (environment: network prohibited by Spec authorization)`, keep an
`outside-evidence row` item in its provenance, and name in the report the
authorization text that forbids it. It never decides a qualifying partial
either (ADR-0276). A source the row can check offline, such as a passage the
PRD quotes or a document under `docs/references/`, is checked, never
blocked." The edge-case bullet on repositories this environment does not hold
gains: "When the Spec's own authorization forbids network access, record
`blocked (environment: network prohibited by Spec authorization)` with the
same provenance item instead; it never decides a qualifying partial either."

The roundfix `archive.md` and `settle.md` references gain, after "need no
Unreachable Acceptance declaration.": "So does an outside-evidence row the
gate did not look up because the Spec's own authorization forbids network
access, recorded as `blocked (environment: network prohibited by Spec authorization)`
with the outside-evidence row named in its provenance (ADR-0276)." The refusal
paragraph of `archive.md` names "a network-denied outside-evidence row or a
network-prohibited outside-evidence row". `spec.md` gains the bullet:
"`SC-QA-BYTES-PINNED` — a requirement of a pending QA Task names a byte
invariant and a commit-shaped hexadecimal token together; compare with the
Delivery Base instead. A gap that runs from the Tasks stage."

`write-tasks` gains this section before `## Process`:

```markdown
## Byte invariants name the delivery base

A QA requirement that holds a file byte-identical compares it with the
Delivery Base, the merge base of the audited head with the default branch,
never with a named commit. A plan waits in the queue behind other Specs, and
each one that merges first may change the file, so the commit where the plan
measured stops being the right comparison. Keep that commit out of the
requirement; when it matters, record it in the TechSpec. `roundfix spec check`
reports `SC-QA-BYTES-PINNED` for a requirement of a pending QA Task that names
a byte invariant and a commit-shaped hexadecimal token together.
```

`docs/user-guide/commands/qa-report.md`, `settle.md` and
`docs/user-guide/context-driven-development.md` gain, before "Any other
blocked outside-evidence row": "So does an outside-evidence row the gate did
not look up because the Spec's own authorization forbids network access,
recorded as `blocked (environment: network prohibited by Spec authorization)`
with the outside-evidence row named in its provenance." The refusal paragraph
of `docs/user-guide/commands/archive.md` names
"`blocked (environment: network denied: <host>)` or
`blocked (environment: network prohibited by Spec authorization)`".
`docs/user-guide/commands/spec.md` gains the table row
"`SC-QA-BYTES-PINNED` | A requirement of a pending QA Task names a byte
invariant and a commit-shaped hexadecimal token together; compare with the
Delivery Base instead. A gap that runs from the Tasks stage."

The guidance of `clause.spec.project-constraints-06-outside-evidence` gains,
after its last sentence: "So is an outside-evidence row the QA gate did not
look up because the Spec's own authorization forbids network access, recorded
as `blocked (environment: network prohibited by Spec authorization)`." The
guidance of `clause.autonomous.loop-01-qa-once` gains, after "recorded as
`blocked (environment: network denied: <host>)`.": "Nor does an
outside-evidence row the QA gate did not look up because the Spec's own
authorization forbids network access, recorded as
`blocked (environment: network prohibited by Spec authorization)`." Every
existing sentence stays byte-identical, and no clause is removed or renamed,
so no Source Baseline retention disposition applies.

`CONTEXT.md` gains, after the Network-Denied Row entry:

```markdown
**Network-Prohibited Row**:
An outside-evidence QA row the gate did not look up because the Spec's own authorization forbids network access, recorded as `blocked (environment: network prohibited by Spec authorization)` with the outside-evidence row named in its provenance. Like the Network-Denied Row, it never decides a qualifying partial; any other wording still needs a declaration or the override.
_Avoid_: Network-denied row, authorization override, skipped source

**Pinned Byte Invariant**:
A QA requirement that holds a file byte-identical to its bytes at a named commit instead of at the Delivery Base. It fails once a Spec that merges first changes the file, and `roundfix spec check` reports it as `SC-QA-BYTES-PINNED`.
_Avoid_: Byte invariant, golden pin, frozen file check
```

The Delivery Base entry gains the sentence: "A QA requirement that holds a
file byte-identical compares it with the Delivery Base, never with a named
commit."

### Invariants

```text
1. The QA policy stays one rule in internal/spec; no caller re-derives it.
2. A Network-Prohibited Row needs the exact status and the outside-evidence
   provenance; anything else is an ordinary environment row.
3. Existing refusal texts are unchanged when a report has no such row.
4. The detector reads requirement text only, never executes a command, and
   never reads a completed Task or a non-QA Task.
5. SC-QA-BYTES-PINNED is a gap that holds every active Spec with no horizon.
6. The QA settlement section of every skill is byte-identical at the
   delivery base.
7. Archived Specs, Archive Records and reports stay byte-identical.
```

## Coverage Map

Deviations from the PRD: None.

- User Story 1 → Network-Prohibited Row in `internal/spec/qa.go`, the park decision through `EnvironmentRowsNeedingOverride`
- User Story 2 → qa-gate wording, the settlement extension heading, the exact-status rule
- User Story 3 → `SC-QA-BYTES-PINNED` detector, write-tasks "Byte invariants name the delivery base"
- User Story 4 → Baseline clauses, user guide pages, glossary
- Core Feature 1 → API Contract 1
- Core Feature 2 → Exact texts for qa-gate, archive-spec, roundfix and its references
- Core Feature 3 → API Contract 2, API Contract 3, write-tasks rule
- Core Feature 4 → API Contract 4, Baseline regeneration, glossary entries
- Success Metric 1 → task_01 policy and the QA replay of the 0258 report
- Success Metric 2 → task_01 shape tests across the commands and the Daemon
- Success Metric 3 → task_02 detector tests and the QA replay of the 0264 Task
- Success Metric 4 → corpus golden at 0 and the QA corpus run
- Success Metric 5 → task_03 clause test and regeneration, task_04 text tests

## Absence Evidence

1. Network-Prohibited Row policy — command: `grep -rn 'network prohibited' internal`; result: no match at `9963ffd3`; the only exempt outside-evidence form is `QANetworkDeniedStatusPrefix`.
2. Pinned Byte Invariant detector — command: `grep -rnE 'QA-BYTES|byte-identical|Delivery Base' internal/speccheck`; result: no match at `9963ffd3`; no detector reads QA requirement text for a commit.
3. Skill and guide texts — command: `grep -rnE 'network prohibited|Byte invariants name' .agents/skills docs/agents docs/user-guide internal/baseline/assets/modules`; result: no match at `9963ffd3`.
4. Glossary terms — command: `grep -nE '^[*][*](Network-Prohibited Row|Pinned Byte Invariant)[*][*]' CONTEXT.md`; result: no match at `9963ffd3`.

## ADR Alternatives

- ADR-0276: this Spec's own decision; it rejects reusing the network-denied form, reading the authorization, a host suffix, editing the settlement table, comparing with the plan's merge commit and exempting a requirement that says "delivery base".
- ADR-0240: no rejected alternative is taken; trusting a row only when the Daemon observed the cause stays rejected, and the new kind joins the one rule.
- ADR-0104: separating the authors of acceptance and implementation stays deferred; only the Network-Prohibited Row stops holding Pull Request preparation.
- ADR-0167: no rejected alternative is taken; the gate still never writes `pass` for evidence it did not record.
- ADR-0093: no rejected alternative is taken; the detector reads written structure and infers nothing from prose meaning.
- ADR-0094: its rejected demand for the full artifact set stays rejected; a Spec without a Task Graph lists the code as skipped.
- ADR-0265: no rejected alternative is taken; like its gaps, the new one takes no horizon and judges no content.
- ADR-0148: no rejected alternative is taken; every Verification names a test or phrase only the work creates.
- ADR-0259: no rejected alternative is taken; every new test is parallel.
- ADR-0271: no rejected alternative is taken; the new tests read files and give no child, goroutine or Git a real-time allowance.
- ADR-0244: no rejected alternative is taken; the three terms follow its declaration contract.
- ADR-0256: its rejected model-judging alternative stays rejected; task_04 updates the command references with the user guide.
- ADR-0184: no rejected alternative is taken; the Surface Transcripts section declares `None.` with its reason.
- ADR-0156: of its rejected options, inferring coverage from prose stays rejected; the Consumers sections reuse its `None.` with a reason form.
- ADR-0168: its horizon stays for `SC-ADR-RELATED`; the new gap takes none, as ADR-0265 decided for its gaps.
- ADR-0183: its horizon stays for receipts; the new gap holds every active Spec instead.
- ADR-0176: no rejected alternative is taken; a QA requirement is authored text inside the projection the detector reads.
- ADR-0229: no rejected alternative is taken; a Network-Prohibited Row joins the exempt rows that never park, and every other environment row parks as before.
- ADR-0237: no rejected alternative is taken; no Delivery Retry rule changes.
- ADR-0236: no rejected alternative is taken; like a check judged on the current default branch, a byte invariant is judged at the Delivery Base, never at a fixed commit.

## Integration Points

None outside the repository. The Delivery Plan, queue revalidation and the QA
gate already promote gaps through `spec check --strict`, so they refuse an
active Spec the detector reports. Adopters receive the two Baseline sentences
through `roundfix baseline update`.

## Testing Approach

New tests, each calling `t.Parallel()` first:

- `internal/spec/qa_partial_policy_test.go`:
  `TestANetworkProhibitedOutsideEvidenceRowNeverDecidesAQualifyingPartial`
  (prohibited only, beside a Pull Request row, beside a Network-Denied Row,
  and the 0258 shape of four prohibited rows, a Pull Request row and one
  covered declaration), `TestANetworkProhibitedRowNeedsItsExactStatusAndOutsideEvidence`
  (no outside-evidence provenance, a host suffix, a changed case, a wrong
  source boundary; each counts 0 and refuses) and
  `TestAnotherEnvironmentRowBesideANetworkProhibitedRowNamesBoth`.
- `internal/cli/qa_partial_policy_test.go` and
  `internal/daemon/qa_partial_policy_test.go`: `qaPartialShapes` gains
  `network prohibited only`,
  `Spec 0258 network prohibited with PR and declaration`,
  `network prohibited without outside evidence`,
  `network prohibited with a host` and
  `another environment beside network prohibited`, so the existing shape
  tests run each through `qa-report accept`, `archive`, `settle` and Daemon
  settlement. `TestRunSpecDoesNotParkAQualifyingPartialAsEnvironmentOnly`
  gains `network prohibited outside evidence` (no park) and
  `network prohibited without outside evidence` (park).
- `internal/speccheck/pinned_bytes_test.go`:
  `TestPinnedByteInvariantGapNamesTheRequirementAndTheCommit`,
  `TestPinnedByteInvariantSkipsUnpinnedRequirementsAndOtherTasks`,
  `TestPinnedByteInvariantRunsAtTheTasksStageAndFailsStrict` and
  `TestPinnedByteInvariantFixNamesTheDeliveryBaseAndTheSkillRule`, on Specs
  built in temporary directories with synthetic hexadecimal tokens.
- `internal/baseline/network_prohibited_row_clause_test.go`:
  `TestTheGuidesExemptANetworkProhibitedOutsideEvidenceRow`, which reads both
  sentences in the embedded clauses, both formatter goldens and both
  repository guides, as `TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow`
  does.
- `skills/acceptance_rules_test.go`:
  `TestTheSettlementSkillsNameTheNetworkProhibitedRow` and
  `TestWriteTasksComparesByteInvariantsWithTheDeliveryBase`.

Declared breaks: none measured. The prototype ran every existing test of the
changed packages unchanged. The derived Baseline files, the version records
and `docs/references/behavior-surfaces.json` change only through their
commands. `docs/references/coverage-record.json` records new tests as
additions and needs no re-record.

## Build Order

1. Network-Prohibited Row in the shared policy, with the spec, command and
   Daemon shape tests.
2. `SC-QA-BYTES-PINNED` detector, its tests and the corpus pins.
3. Baseline clauses, regeneration, the clause test and the glossary
   (depends on: 1, 2).
4. Skill texts, mirrors, versions, the skill test, the command references,
   the user guide pages and the Behavior Surface Record (depends on: 1, 2).
5. QA gate (depends on: 3, 4).

This Spec rebases after Specs 0266, 0267, 0268 and 0269 of cycle AG (see
`_prd.md` → Prerequisites). Its design depends on no other AG Spec's code.

## Risks & Considerations

- Sibling Specs of the cycle may change the owned skill version record, the
  corpus golden, the Behavior Surface Record or Baseline derived files. After
  a rebase the Tasks rerun the record and regeneration commands.
- A gate could label a reachable row as network-prohibited. The skill
  requires it to name the authorization text and to check offline sources;
  the report and the record make the claim checkable afterwards.
- The commit-shaped token is coarse: an all-letter hash is missed, and a
  hexadecimal word with a digit beside a byte phrase would match. The scan of
  235 archived Specs found no false match in a QA requirement.
- An active Spec queued before this one merges is refused under `--strict` if
  a QA requirement pins bytes to a commit; the adopter notes and the fix say
  how to rewrite it.

## Glossary

None.

## Decisions

- One exact status joins the policy in `internal/spec`; see ADR-0276.
- A structural gap on pending QA requirements, holding every active Spec; see
  ADR-0276.
- The settlement table extension lives under its own heading in the three
  skills that carry the table.
- The glossary and Baseline land in the Baseline Task; the skills and user
  guide land in one docs Task, so the two can run in parallel.
