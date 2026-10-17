---
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
prd: _prd.md
created: 2026-10-16
---

# Authoring guards hold the plan, the glossary and the PRD to their rules — Technical Spec

## Executive Summary

Two guards are new Spec Consistency Check detectors in `internal/speccheck`,
beside the skills and glossary detectors they resemble: `SC-SURFACE-RECORD-UNTASKED`
reads a pending guide Task's declarations and its Spec's `_authorization.md`
Limits, and `SC-DECISIONS-UNRECORDED` reads a PRD's `## Decisions` section
under a horizon. Two guards are new Repository Contract Tests in
`internal/docscontract`: one derives the QA secret refusal from the Daemon's
source and the Secret Match format and compares every document that states
it; the other holds `CONTEXT.md` to five structural rules. The skills,
guides and glossary that describe each guard change in the Task that adds it.

The trade-off accepted is precision over reach. The surface record detector
reads only user guide pages, whose fingerprint is their bytes, and leaves a
command-help change to the record's own test; the phrase rule reads runs of
capitalized words and needs an allowlist; and the Clarification Record holds
only PRDs committed after its guide. Each choice keeps the guard green on
honest work, so a failure means a real gap (ADR-0280).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is added; two finding codes join the code set. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call; the checks read local files and Git history. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0280 governs every decision below:
  "A guide Task records its own fingerprint.", "A PRD records who decided."
  and "The glossary is held to its own rule.". ADR-0256: "A release waits for
  the skills that describe a changed surface". ADR-0244: "A Spec declares the
  domain terms it introduces". ADR-0273: "A Secret Match reports the
  repository-relative path, the line and the kind".
  ADR-0091: "The QA gate is a Task node of its own type". ADR-0217: "A Cursor
  selection names the advertised model value". ADR-0182: "A Task settles on
  the facts its gate will check". ADR-0117: "A defect is checked by the stage
  that can produce it". ADR-0104: "A Spec accepts on evidence it did not
  author". ADR-0096: "The QA gate proves machine facts before it spends an
  agent turn". ADR-0167: "The pre-PR Pull Request row never decides a
  qualifying partial". ADR-0240: "One QA partial policy, and rows a Run
  sandbox cannot reach". ADR-0276: "A row the authorization keeps offline
  qualifies, and a byte invariant compares with the Delivery Base". ADR-0229:
  "An operator archive resumes any park from the Run start". ADR-0237: "A
  Delivery Retry records a merge made outside the queue". Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH grant of 2026-10-16 ("Governed Paths declarados") and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-prd/references/clarification-record.md`,
  `.agents/skills/write-prd/references/prd-template.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`,
  `skills/write-prd/references/prd-template.md`,
  `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Current behavior

Characterized on 2026-10-16 with `bin/roundfix` 0.82.0 built at `311c11bc`,
on the authoring branch at `9e723b82`:

- Surface record. A tree exported from `e1e376dd` (the 0268 plan), committed
  into a scratch repository, gives
  `roundfix spec check 0268-a-turn-that-did-not-finish-is-not-accepted --strict`
  exit 0 with `No findings.`, although its task_01 declares
  `docs/user-guide/configuration.md`, `docs/user-guide/usage.md` and
  `docs/user-guide/commands/deliver.md` under `interface:`, does not declare
  `docs/references/behavior-surfaces.json`, names no `-record-skill-coverage`
  step, and its `_authorization.md` Limits list
  `docs/references/behavior-surfaces.json`. The fix commit `cc0ac876`, no
  longer reachable from any ref, added the record to task_01's Context, the
  step to its requirements, and moved the record out of Limits. `internal/speccheck` names
  `skillcoverage.MapPath` only; nothing reads `skillcoverage.RecordPath`.
  `TestTheSkillCoverageMapIsCurrent` is `//verify:always`, so the repository
  Verification at every Task's settlement runs it.
- QA secret refusal. `internal/daemon/task_engine.go` builds
  `"QA evidence carries a secret: " + strings.Join(details, "; ")` from
  `secretscan.Match.String()`, which formats `<path>:<line>: <kind>`;
  `TestQAReportWithASecretMatchSettlesFailedWithoutACommit` asserts
  `path + ":2: openrouter-key"`. `.agents/skills/roundfix/references/implement.md`,
  its mirror and `docs/user-guide/commands/implement.md` write
  `QA evidence carries a secret: <path>:<line> (<kind>)`. `CHANGELOG.md`
  already writes the colon form. No test compares them.
- Glossary. `CONTEXT.md` holds 320 entries, each a `**Term**:` line followed
  by its definition and one `_Avoid_:` line (one entry, Review Source, has a
  blank line before it). Measured with the rules of API Contract 4 by a
  scratch script: rule 1 finds 27 uncovered runs (Pull Request ×11, Pre-PR
  Review ×3, Better Auth ×3, QA Task ×2, Claude Code ×2, Run ID ×2, Roundfix
  Runs ×2, and once each Active ADR, Coverage Map, Conventional Commits, Agent
  Client Protocol, Agent Selection Roundfix, Agent Model Roundfix, Agent
  Models Roundfix, Agent Selection Profile Roundfix, Implement Task, Task
  Verification, Passing Verification, Spec Consistency, Spec Run Worktrees,
  Specs Root, GitHub Release, Success Metrics, API Contracts, Normal QA,
  Architecture Decision Record, Roundfix OpenRouter); rule 2 finds every one
  of the 61 linked ADRs accepted and no other link; rule 3 finds the four
  mapped lists equal to their constants once `QA-PRECONDITION` is named as a
  row label; rule 4 finds one item with prose, in Corrective Spec
  (`inherited authorization; a Corrective Task fixes an active Spec, not an
  archived one`); rule 5 finds no repeated sentence once link destinations are
  dropped. 77 `_Avoid_` items name another defined term. No test reads
  `CONTEXT.md`.
- Clarification. `write-prd`'s template ends with `## Decisions` ("Product
  decisions made during clarification, one line each") and says nothing about
  who answered; the minimal PRD of `write-techspec` has no such section.
  `roundfix spec check --stage prd` reports nothing about it. Spec 0272's PRD
  already opens its section with dated answers.

## System Architecture

- `internal/speccheck/surface_record.go` (new) holds
  `detectSurfaceRecord`, called from `Check` in `constraints.go` next to
  `detectSkillsDeclaration` with the loaded graph. `coherence.go` registers
  the code at `StageTasks`.
- `internal/speccheck/clarification_record.go` (new) holds
  `detectClarificationRecord` and its horizon, called from `Check` in
  `constraints.go` and from the authoring-stage sweep in `coherence.go`, which
  registers the code at `StagePRD`.
- `internal/docscontract/command_documentation_test.go` (existing, `//verify:always`) gains the QA secret refusal
  tests, which parse the Daemon's
  non-test sources with `go/parser` and imports `internal/secretscan`.
- `internal/docscontract/user_guide_contract_test.go` (existing, `//verify:always`) gains the glossary contract
  tests, which read `CONTEXT.md`
  and `docs/adr/`, and imports `internal/spec`, `internal/speccheck`,
  `internal/worktree` and `internal/agent` for the mapped constants. Its
  allowlist and mapping are two package-level tables at the top of the file.
- Corpus pins: `corpusFindingCodes` in `internal/docscontract/corpus_test.go`,
  `internal/docscontract/testdata/corpus-golden.json` and
  `internal/spec/archive_layout_characterization_test.go` gain each code.

## Implementation Design

### Interfaces

```go
// internal/speccheck/surface_record.go
const CodeSurfaceRecordUntasked = "SC-SURFACE-RECORD-UNTASKED"

func detectSurfaceRecord(result *Result, repoRoot, specDir string, graph *spec.Graph) error

// internal/speccheck/clarification_record.go
const (
	CodeDecisionsUnrecorded = "SC-DECISIONS-UNRECORDED"
	// ClarificationGuidePath is the guide whose adding commit starts the horizon.
	ClarificationGuidePath = ".agents/skills/write-prd/references/clarification-record.md"
)

func detectClarificationRecord(result *Result, repoRoot, prdPath string) error
```

### Data Models

No new file format. The Clarification Record is Markdown inside a PRD's
`## Decisions` section:

```text
- YYYY-MM-DD: <question a person answered> — <answer>
Autonomous route: no person present; decisions taken from <source>.
```

### API Contracts

1. API Contract: surface record gap. For each Task of the graph that is not
   completed, is not the manifest's `qa:` Task and is not of type `qa`, and
   declares under `interface:`, `creates:` or `deletes:` a path that starts
   with `docs/user-guide/` and ends with `.md`, the detector reports one gap
   when the Task does not both declare `docs/references/behavior-surfaces.json`
   under `interface:` or `creates:` and name `-record-skill-coverage` in one of
   its requirements. Summary: `Task <id> declares user guide <path> without
   declaring docs/references/behavior-surfaces.json and naming its
   -record-skill-coverage step`, located at that Context line. It reports a
   second gap, located at the Limits line, when the Spec's `_authorization.md`
   has a `## Limits` section with a line containing
   `docs/references/behavior-surfaces.json`. Both fixes name the write-tasks
   rule "Skills Tasks for changed Behavior Surfaces". Another Task's
   declaration never covers the guide Task (ADR-0280). A repository without a
   regular `docs/references/behavior-surfaces.json` skips the code with that
   path; a Spec without a Task Graph skips it with `_tasks.md`; the PRD and
   TechSpec stages list it as skipped. No horizon applies.
2. API Contract: QA secret refusal contract. The test finds, in the non-test
   `.go` files of `internal/daemon`, exactly one string literal that starts
   with `QA evidence carries a secret`, and builds the documented form as that
   literal followed by `secretscan.Match{Path: "<path>", Line: 1, Kind:
   "<kind>"}.String()` with `<path>:1:` replaced by `<path>:<line>:`. Every
   `.md` file under `docs/user-guide/`, `docs/agents/`, `.agents/skills/` and
   `skills/`, and `CONTEXT.md`, is whitespace-normalized; each occurrence of
   `QA evidence carries a secret` in it must begin the documented form, and
   `.agents/skills/roundfix/references/implement.md`,
   `skills/roundfix/references/implement.md` and
   `docs/user-guide/commands/implement.md` must each hold one. A failure names
   the file and the text found. `CHANGELOG.md` is release history and is not
   read.
3. API Contract: Clarification Record gap. For a PRD that exists and whose
   adding commit descends from the oldest commit that added
   `ClarificationGuidePath`, or that is uncommitted, the detector reports one
   gap when the PRD has no `## Decisions` line, or when the lines of that
   section hold neither a line matching `^- (\d{4}-\d{2}-\d{2}): \S` with a
   valid calendar date nor a line matching
   `^Autonomous route: no person present; decisions taken from \S.*\.$`.
   Summary: `<prd> has no Clarification Record`, located at the section's
   heading or at line 1. Fix: `Under ## Decisions, write one - YYYY-MM-DD:
   <question> — <answer> line for each question a person answered, or the line
   Autonomous route: no person present; decisions taken from <source>.,
   following the write-prd rule "Clarification Record".` A PRD committed before
   the horizon, or a repository without the guide, skips the code with
   `a PRD committed at or after .agents/skills/write-prd/references/clarification-record.md`;
   a missing PRD skips it with its path. Both shapes may appear together. It
   runs from the PRD stage and in the full sweep.
4. API Contract: glossary contract. An entry is a `**Term**:` line and the
   lines after it up to the next such line or heading; its body is every line
   except its `_Avoid_:` line. Before rules 1 and 5, a Markdown link becomes
   its text and a code span becomes a break.
   - Rule 1, phrases. A run is a maximal sequence of words that start with an
     uppercase letter; a word that ends in `'s` ends its run, a hyphen followed
     by a lowercase letter ends the word before it, and a leading `A`, `An`,
     `The`, `One`, `Each`, `Every`, `Its`, `This`, `That`, `These`, `Those`,
     `What`, `Like`, `Anything`, `No`, `Any`, `Some`, `Both`, `When`, `If` or
     `Only` is dropped. A run of two or more words is covered when it, or each
     piece of a split of it into consecutive pieces of two to five words,
     equals, ignoring case, a defined term, a defined term followed by `s` or
     `es`, a term ending in `y` written with `ies`, or an allowlist phrase.
     Each uncovered run fails with its entry and line.
   - Rule 2, links. Every Markdown link destination in `CONTEXT.md` is a file
     directly under `docs/adr/` whose front matter `status` is `accepted`, or
     which has no lifecycle front matter; any other destination fails.
   - Rule 3, closed lists. For each mapped entry, the set of code spans in its
     body equals the mapped constant values plus the entry's named exceptions.
     The mapping is: Task Type to `spec.AllowedTaskTypes()`; Mechanical
     Refusal Code to the four `speccheck.CodeMechanical*` constants, with the
     exception `QA-PRECONDITION`; Run Worktree Reconciliation to the six
     `worktree.Reconciliation*` states; Selection Encoding to the five
     `agent.SelectionEncoding*` values.
   - Rule 4, avoid lines. Each entry has exactly one line that starts with
     `_Avoid_: `, and it is the entry's last non-blank line; its items, split
     on `, `, are non-blank and hold no `;`, `—`, `(`, `)` or `:` and do not
     end in `.`. An item may name a defined term.
   - Rule 5, one owner per rule. After whitespace normalization, the body is
     split into sentences after `.`, `;`, `?` or `!` followed by a space; no
     sentence longer than sixty characters occurs in two entries.
   Each allowlist entry is a phrase and a non-blank reason; the test fails on
   an allowlist entry that rule 1 no longer needs, so the list only holds live
   exceptions.
5. API Contract: everything else stays. No command's help, flag, exit code or
   output line shape changes; existing codes keep their counts; the Daemon's
   refusal text and the Secret Match format keep their bytes; no Baseline
   module changes.

### Surface Transcripts

None. No command's help, flag or output shape changes: `roundfix spec check`
prints the two new codes in its existing finding-line shape, and the contract
tests are not a command surface.

### Invariants

```text
1. A guide Task that declares the record and names the step yields no SC-SURFACE-RECORD-UNTASKED.
2. A completed Task, the QA Task and a Task declaring no user guide page are never read for it.
3. A PRD committed before the Clarification horizon is skipped, never reported.
4. The documented QA refusal is computed from the Daemon literal and Match.String, never copied.
5. The glossary contract passes on the delivered CONTEXT.md and fails on one planted violation of each rule.
6. Every existing finding count in the corpus golden stays.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Story 1, Core Feature 1, Success Metric 1 → `detectSurfaceRecord` (API Contract 1)
- Goal 2, User Story 2, Core Feature 2, Success Metric 2 → the QA secret refusal tests in `command_documentation_test.go` (API Contract 2)
- Goal 3, User Story 3, Core Feature 3, Success Metric 3 → the glossary contract tests in `user_guide_contract_test.go` (API Contract 4)
- Goal 4, User Story 4, Core Feature 4, Success Metric 4 → `detectClarificationRecord` (API Contract 3), the `write-prd` and `write-techspec` skills
- Core Feature 5 → the `spec` and `implement` user guides, the Roundfix `spec` and `implement` references, the `write-tasks` skill, `CONTEXT.md`
- Success Metric 5 → corpus pins, API Contract 5, Invariant 6

## Absence Evidence

1. `detectSurfaceRecord` — command: `grep -rn 'behavior-surfaces\|skillcoverage.RecordPath\|record-skill-coverage' internal/speccheck --include='*.go'`; result: no match in production code at `9e723b82`; the skills detector reads only `skillcoverage.MapPath`.
2. QA secret refusal contract — command: `grep -rn 'carries a secret' internal/docscontract skills --include='*_test.go'`; result: no match at `9e723b82`; only `internal/daemon/qa_secret_scan_test.go` asserts the Daemon's reason.
3. Glossary contract — command: `grep -rln 'CONTEXT.md' internal/docscontract`; result: no match at `9e723b82`; no docs contract reads the glossary.
4. `detectClarificationRecord` — command: `grep -rn '"## Decisions"\|Autonomous route' internal/speccheck internal/spec --include='*.go'`; result: no match at `9e723b82`.
5. `clarification-record.md` guide — command: `find .agents/skills/write-prd/references -name 'clarification-record.md'`; result: absent at `9e723b82`; the folder holds `glossary.md` and `prd-template.md`.

## ADR Alternatives

- ADR-0280: this Spec's decision; it rejects letting another Task record the fingerprint, inferring help changes, holding the record with no horizon, mandatory grilling, holding the glossary in `spec check`, and forbidding defined terms in `_Avoid_`.
- ADR-0256: no rejected alternative is taken; fingerprints stay in the generated record and the map stays authored.
- ADR-0244: no rejected alternative is taken; no model judges a term, and the glossary meets the same deterministic rule.
- ADR-0273: no rejected alternative is taken; the Secret Match format and the refusal are unchanged.
- ADR-0091: no rejected alternative is taken; the Task Type mapping reads its `qa` value from the code.
- ADR-0217: no rejected alternative is taken; the ACP Runtime list stays unmapped and keeps its text.
- ADR-0182: no rejected alternative is taken; the guide Task settles on the record its own settlement checks.
- ADR-0117: no rejected alternative is taken; each detector runs from the stage whose artifact it reads.
- ADR-0104: no rejected alternative is taken; the outside-evidence row quotes published sources.
- ADR-0096: no rejected alternative is taken; the gate's mechanical stage is unchanged.
- ADR-0167: no rejected alternative is taken; the Pull Request row stays blocked before a Pull Request exists.
- ADR-0240: no rejected alternative is taken; the QA partial policy is unchanged.
- ADR-0276: no rejected alternative is taken; the corpus row compares with the Delivery Base.
- ADR-0229: no rejected alternative is taken; the archive is unchanged.
- ADR-0237: no rejected alternative is taken; the Delivery Queue is unchanged.

## Integration Points

None outside the repository. The Clarification horizon reads `git log` in the
repository holding the Spec Root, as the glossary horizon does.

## Testing Approach

Every new top-level test calls `t.Parallel()` first and builds temporary
repositories with the speccheck fixtures already used by
`skills_declaration_test.go`, `glossary_test.go` and `pinned_bytes_test.go`.

- `internal/speccheck/surface_record_test.go`: the Spec 0268 shape (a guide
  Task without the record and Limits naming it) yields two gaps; a Task
  declaring the record without the step, and the record and step in another
  Task, still yield one; the fixed shape, a completed Task, a QA Task, a
  repository without the record and the PRD stage yield none, with the skips.
- `internal/speccheck/clarification_record_test.go`: missing section, a
  section with neither shape, an invalid date, each valid shape, a PRD before
  the horizon, an uncommitted PRD, and the PRD stage with `--strict`.
- the QA secret refusal tests in `internal/docscontract/command_documentation_test.go`: the live documents, and a
  helper case showing the parenthesized sentence fails.
- the glossary contract tests in `internal/docscontract/user_guide_contract_test.go`: one top-level test per
  rule on the live `CONTEXT.md`, and one that plants a violation of each rule
  in a synthetic glossary and requires each rule to fail.
- No existing test changes except the corpus pins, whose update text gains
  one sentence per code.

## Build Order

1. `SC-SURFACE-RECORD-UNTASKED`, its tests, corpus pins, the `spec` guide, the
   Roundfix `spec` reference, the `write-tasks` rule, versions and the
   Behavior Surface Record (no dependency).
2. The QA secret refusal contract, the `implement` reference and guide,
   versions and the Behavior Surface Record (depends on: 1, which writes the
   same Roundfix `SKILL.md`, version record and Behavior Surface Record).
3. The glossary contract and the `CONTEXT.md` fixes it needs (no dependency).
4. `SC-DECISIONS-UNRECORDED`, its tests, corpus pins, the `write-prd` guide,
   template and skill, the `write-techspec` minimal PRD, the `spec` guide and
   reference, versions, the Behavior Surface Record and the glossary entry
   (depends on: 2, which writes the same versions and record, and 3, whose
   contract the new entry must pass).
5. QA gate (depends on: 4).

No other Spec is a prerequisite and no Verification reads another Spec's
artifacts. Rebase on main after 0272 merges and run the strict check again.

## Risks & Considerations

- An active Spec authored beside this one may declare a guide without the
  record; the corpus golden then records its count and `--strict` holds that
  Spec until its plan is fixed, which is the intended effect.
- The phrase rule can flag a new legitimate compound in a later Spec; that
  Spec defines the term, rewords, or adds an allowlist entry with a reason.
- A new constant in a mapped closed list is missed until its mapping row is
  updated; the mapping names constants, so a reviewer sees the list beside it.
- Raising the Roundfix skill version in two Tasks is expected; each Task's
  record command writes the next free version.

## Glossary

- adds: **Clarification Record**

## Decisions

- The guide Task itself declares the record and the step. See ADR-0280.
- The surface detector reads user guide pages only. See ADR-0280.
- The Clarification Record has a horizon. See ADR-0280.
- The documented QA refusal is derived from the producer, so no production
  symbol is added for the test.
- The allowlist and mapping are tables in the glossary test, each allowlist
  entry with its reason. See ADR-0280.

## Amendment of 2026-10-17

The first Run settled task_01 and stopped at task_02 and task_03 on
`TestRepositoryContractTestsDeclareTheirRelevance`
(`internal/verifyselect/contracts_repository_test.go`): every `docscontract`
file is an always Repository Contract Test, and the suite pins nine such
files, so a new file with a `//verify:relevant` directive fails it twice over.
The two contract test sets now live in existing `//verify:always` files,
`command_documentation_test.go` and `user_guide_contract_test.go`, one per
Task, so the two Tasks share no file and the pinned count holds.
