---
spec: 0340-qa-evidence-links-resolve-in-the-committed-tree
prd: _prd.md
created: 2026-11-01
---

# QA evidence links resolve in the committed tree — Technical Spec

## Executive Summary

No QA settlement path reads a QA Report's links today: `runQAGate`,
`roundfix qa-report accept` and the Settle Command judge the report through
`spec.QAReportEligibility` alone, and the QA Report commit refuses only a
Secret Match or a non-image binary. The fix moves the Markdown link extraction
of `MissingRelativeLinks` into an exported `spec.RelativeLinks`, adds
`spec.UnresolvedEvidenceLinks`, which classifies each local link of a report
as resolved, `ignored` or `missing` by asking Git what the commit would hold,
and makes all three paths refuse an Unresolved Evidence Link with one refusal
text. The Daemon first returns the QA Agent Session once with the links, then
checks again at the QA Report commit boundary and, on a remaining link,
settles the QA Task `failed` without a QA Report commit. A new `qa-gate`
section, the guides, the Roundfix references, `CONTEXT.md` and ADR-0351
state the rule.

The trade-off accepted: the Daemon's check asks Git about the Run Worktree
rather than about the exact path list `prepareQAReportCommit` stages, so one
rule serves all three paths; `_techspec.md` → Risks names the gap.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0351, new with this plan, records
  the resolution rule, the three paths, the repair turn and the commit
  refusal, ADR-0351: "A local link resolves when the committed tree would
  hold its target."; the design implements it as recorded. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AU answers of 2026-11-01 quoted in the authorization record; the
  Governed Paths are the `qa-gate` and Roundfix skill files the Tasks change;
  no Makefile, workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0340-qa-evidence-links-resolve-in-the-committed-tree/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/settle.md`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-01 at `dd3b2b62`; `_prd.md` → Validation records the run.

1. `runQAGate` (`internal/daemon/task_engine.go`) runs the QA turn, then
   `settleQAVerdict`, which reads the newest report and applies
   `spec.QAReportEligibility` and `qaSettlementReportEligibility`. It then
   prepares the QA Report commit with `prepareQAReportCommit`, runs the
   Secret Scan and the non-image binary check over the prepared paths, settles
   the QA Task, and commits only when neither refuses. No step reads the
   report's links, and the QA step has no repair turn.
2. `snapshotQAPaths` sees untracked files that Git does not ignore, so an
   ignored evidence file never enters the prepared paths and never reaches
   the QA Report commit, silently.
3. `runQAReportCommand` (`internal/cli/qa_report.go`) reads the report named
   on the command line and exits 1 only when `ReadQAReportFile` or
   `QAReportEligibility` refuses; its stderr line is
   `roundfix: QA Report eligibility refused: <reason>`.
4. `requireSettleQAReportEligibility` (`internal/cli/settle.go`) applies the
   same decision to the newest report of a `qa` Task's surface and wraps the
   reason as `<reason> (report <path>)`; `runSettleCommand` prints
   `roundfix: settle QA Report is ineligible: ` before it.
5. `MissingRelativeLinks` (`internal/spec/record_links.go`) scans Markdown
   with `scanArchiveMarkdown`, which masks code and skips URLs, resolves each
   destination with `archiveLinkTarget` against the file's directory, keeps
   targets inside the repository (`historyRepositoryPath`) and reports those
   absent on disk or pointing at the file itself. `detectHistoryRetained`
   (`internal/speccheck/history_retained.go`), `PlanLiveLinks` and
   `PlanRecordLinks` read it.
6. The `qa-gate` skill's evidence sections cover runnable evidence, text or
   image content and secrets; none mentions ignore rules or links.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/spec/record_links.go`: `RelativeLinks`, and `MissingRelativeLinks`
  rewritten over it.
- `internal/spec/qa_evidence_links.go` (new): the classification and the
  refusal error.
- `internal/cli/qa_report.go` and `internal/cli/settle.go`: the CLI paths.
- `internal/daemon/task_engine.go`: the repair turn and the commit-boundary
  refusal.
- Tests in `internal/spec`, `internal/cli` and `internal/daemon`; the
  `qa-gate` skill, the guides, the Roundfix references, `CONTEXT.md`.

## Implementation Design

### Interfaces

```go
// internal/spec/record_links.go

// RelativeLink is a Markdown link outside code whose destination is a
// relative path; Target is resolved against the linking file's directory and
// may leave the repository.
type RelativeLink struct {
	Line                int
	Destination, Target string
}

// RelativeLinks lists every relative link of a Markdown file at path, in
// document order, through the scanner MissingRelativeLinks uses.
func RelativeLinks(path string, content []byte) []RelativeLink

// internal/spec/qa_evidence_links.go

const (
	EvidenceLinkIgnored = "ignored"
	EvidenceLinkMissing = "missing"
)

// UnresolvedEvidenceLink is a local link of a QA Report the committed tree
// would not hold.
type UnresolvedEvidenceLink struct {
	Line        int
	Destination string
	Target      string
	Reason      string // EvidenceLinkIgnored or EvidenceLinkMissing
}

// UnresolvedEvidenceLinks classifies every local link of the QA Report at
// reportPath; an empty result means every link resolves.
func UnresolvedEvidenceLinks(ctx context.Context, reportPath string) ([]UnresolvedEvidenceLink, error)

// EvidenceLinksError is the shared refusal of an Unresolved Evidence Link.
type EvidenceLinksError struct{ Links []UnresolvedEvidenceLink }

func (err EvidenceLinksError) Error() string

// RequireResolvedEvidenceLinks returns nil, an EvidenceLinksError, or an
// infrastructure error from reading the report or Git.
func RequireResolvedEvidenceLinks(ctx context.Context, reportPath string) error
```

### Data Models

No schema changes. The Daemon publishes one new QA Run Event phase,
`evidence_links`, with payload
`{"phase": "evidence_links", "outcome": "repair" | "refused", "links": [{"link": <destination>, "line": <n>, "reason": "ignored" | "missing"}]}`.

### API Contracts

1. API Contract: the extraction. `RelativeLinks(path, content)` returns, for
   each destination `scanArchiveMarkdown` finds, the 1-based line, the
   destination as written and `filepath.ToSlash(target)` where
   `archiveLinkTarget(destination, filepath.Dir(path))` reports it relative;
   other destinations are skipped. `MissingRelativeLinks` keeps its signature
   and results: it filters `RelativeLinks` by `historyRepositoryPath` and by
   self or absent target, exactly as today.
2. API Contract: the classification. `UnresolvedEvidenceLinks` reads the
   report and resolves the repository with `git rev-parse --show-toplevel`
   from the report's directory.
   - In a Git repository it computes each link's target relative to the
     repository root. A target outside the root is `missing`. A target that
     does not exist on disk is `missing`. An existing target resolves when
     one `git ls-files -z --cached --others --exclude-standard` over the
     targets (with literal pathspecs) lists the path itself or, for a
     directory, a path under it; otherwise it is `ignored`.
   - Outside a Git repository a target resolves when it exists, and is
     `missing` otherwise.
   - A link whose target is the report itself resolves. Results keep
     document order.
3. API Contract: the refusal text. `EvidenceLinksError.Error()` is
   `QA Report evidence links do not resolve in the committed tree: ` followed
   by `<destination> (line <n>, <reason>)` per link, joined by `; `.
4. API Contract: `qa-report accept`. After `QAReportEligibility` accepts,
   `runQAReportCommand` calls `RequireResolvedEvidenceLinks(ctx, reportPath)`;
   a refusal prints `roundfix: QA Report eligibility refused: <error>` and
   exits 1. The usage text keeps its bytes: an unresolved link is one more way
   a report is "unacceptable".
5. API Contract: the Settle Command. After `QAReportEligibility` accepts,
   `requireSettleQAReportEligibility` calls `RequireResolvedEvidenceLinks` on
   the same report and wraps a refusal as `<error> (report <path>)`, so the
   stderr line is
   `roundfix: settle QA Report is ineligible: QA Report evidence links do not resolve in the committed tree: … (report <path>)`
   and the Task file stays unchanged.
6. API Contract: the Daemon repair turn. In `runQAGate`, after the QA turn
   and its background-task check, inside the branch that ran the Agent, the
   Daemon calls `UnresolvedEvidenceLinks` on the newest report. When it
   returns links, the Daemon publishes the `evidence_links` event with outcome
   `repair`, prints `QA evidence links unresolved: <n>; returning the QA turn
   once.` to progress, and runs one turn in the same Agent Session through
   `runAgentSession` with the prompt `qaEvidenceLinkRepairPrompt` builds. The
   prompt names the report path and, per link, the destination, line and
   reason, and says: save evidence under a name the repository does not
   ignore, such as a `.txt` extension, and update the link; restore a missing
   file or stop citing a file you did not keep; never change a verdict or a
   row to hide lost evidence. Stop, deadline and Agent errors of that turn are
   handled as the QA turn's are. No second repair turn runs.
7. API Contract: the commit-boundary refusal. After the Secret Scan and the
   non-image binary check find nothing, the Daemon calls
   `UnresolvedEvidenceLinks` on the newest report. A remaining link sets the
   QA status `failed`, `accepted` false and the reason to
   `EvidenceLinksError{Links}.Error()`, publishes the `evidence_links` event
   with outcome `refused`, prints `QA Report refused: <n> unresolved evidence
   link(s)` to progress, and returns before `commitPreparedQAReport`, so no QA
   Report commit is created and the report and evidence stay in the Run
   Worktree. A mechanically closed report, which ran no Agent turn, gets this
   check without a repair turn. A read or Git error halts the Run as an
   infrastructure failure, as the Secret Scan's does.
8. API Contract: the rules in writing.
   - `.agents/skills/qa-gate/SKILL.md`: a new section
     `### Evidence links resolve in the committed tree`, after
     `### Evidence carries no secret` and before
     `### Declared rows name their criterion`, that says: save evidence under
     names the repository does not ignore (a `.txt` extension for logs and
     command output); never cite a file you did not keep; every local link of
     the report must resolve in the committed tree, and an ignored or missing
     target settles the QA Task `failed` after one repair turn. The sentence
     "Never cite a file you did not keep." appears verbatim.
   - `docs/user-guide/commands/qa-report.md` and
     `.agents/skills/roundfix/references/settle.md` (QA Report acceptance):
     the rule, both reasons and the refusal text. The sentence "Every local
     link of the report must resolve in the committed tree" appears verbatim.
   - `docs/user-guide/commands/settle.md`: the same rule for a `qa` Task with
     its stderr line.
   - `docs/user-guide/commands/implement.md` and
     `.agents/skills/roundfix/references/implement.md`: after the non-image
     binary paragraph, the repair turn, the `evidence_links` event and the
     refusal without a QA Report commit. The sentence "the Daemon returns the
     QA turn once" appears verbatim.
   - `CONTEXT.md`: adds **Unresolved Evidence Link** and adds to **QA
     Report** that it settles only when its local links resolve in the
     committed tree.

### Surface Transcripts

None. The changed lines are asserted by the named tests; no whole-command
transcript is declared.

### Invariants

```text
1. No QA settlement path accepts a report with a local link whose target the committed tree would not hold.
2. A tracked target resolves whatever the ignore patterns say; an existing untracked ignored target is ignored; an absent target or one outside the repository is missing.
3. MissingRelativeLinks, and through it SC-HISTORY-RETAINED, returns the same links as before.
4. The Daemon runs at most one repair turn per QA step, and a refused report is never committed.
5. An eligible report with no unresolved link is accepted exactly as before on every path.
```

### Declared breaks

- `roundfix qa-report accept` and the Settle Command now exit 1 on an
  eligible report with an Unresolved Evidence Link; the Daemon settles such a
  QA Task `failed`. No existing test report carries a relative link
  (`grep -rn '](' internal/cli/*qa*_test.go internal/daemon/*_test.go` finds
  none), so no existing assertion changes.

## Coverage Map

- Goal 1, Core Features 2 and 3, Success Metrics 1, 3, 4 and 5 →
  `runQAReportCommand`, `requireSettleQAReportEligibility`, `runQAGate`
  (API Contracts 4 to 7), Invariants 1 and 4.
- Goal 2, Core Feature 1, Success Metrics 1, 2, 3 and 5 →
  `UnresolvedEvidenceLinks` and `EvidenceLinksError` (API Contracts 2 and 3),
  Invariant 2.
- Goal 3, Core Feature 3, Success Metric 4 → the repair turn (API Contract
  6), Invariant 4.
- Goal 4, Core Feature 1, Success Metric 6 → `RelativeLinks` and
  `MissingRelativeLinks` (API Contract 1), Invariant 3.
- Goal 5, Core Feature 4, Success Metric 7 → the `qa-gate` section (API
  Contract 8).
- Goal 6 → Invariants 3 and 5, Declared breaks.

Deviations from the PRD: None.

## Absence Evidence

1. A link check in QA settlement — command: `grep -rn 'RelativeLinks\|UnresolvedEvidence\|check-ignore' internal/daemon/task_engine.go internal/cli/qa_report.go internal/cli/settle.go`; result: no match; the three paths call only `QAReportEligibility`.
2. An exported link extraction without the existence filter — command: `grep -rn 'func RelativeLinks\|type RelativeLink' internal/spec`; result: no match; `MissingRelativeLinks` scans and filters in one loop.
3. A QA repair turn — command: `grep -n 'evidence_links\|repairQA' internal/daemon/task_engine.go`; result: no match; `repairTaskVerification` serves non-QA Tasks only.
4. Ignore guidance in the gate skill — command: `grep -n -i 'gitignore\|ignored\|did not keep' .agents/skills/qa-gate/SKILL.md`; result: no match.

## ADR Alternatives

- ADR-0351: implemented as recorded, including its rejection of forcing
  ignored evidence with `git add -f`, of a Daemon rename, of a warning, of a
  check left to `spec check` and of a refusal without a repair turn.

## Integration Points

None external. The check reads local files and the local Git index.

## Testing Approach

All tests are hermetic: temporary Git repositories, the existing CLI
workspace helpers, the Daemon's `newTaskCycleFixture` with its fake QA Agent,
and no network.

- `internal/spec/qa_evidence_links_test.go` (new):
  - `TestUnresolvedEvidenceLinksNameTheReason`, table subtests over a
    temporary Git repository whose committed `.gitignore` holds `*.log`, with
    a report under `docs/specs/x/qa/`: a present ignored `evidence/a.log` is
    `ignored`; an absent `evidence/b.txt` is `missing`; a committed
    `evidence/c.log` added with `git add -f` resolves; an untracked
    `evidence/d.txt` resolves; a directory link `evidence/` holding `d.txt`
    resolves; `../../../../../outside.txt`, which leaves the repository, is `missing`; a link inside a code
    span, an `https://` URL and a `#anchor` are absent from the result; the
    result keeps document order and lines.
  - `TestUnresolvedEvidenceLinksOutsideGitReadTheDisk`: under a plain
    temporary directory a present target resolves and an absent one is
    `missing`.
  - `TestEvidenceLinksErrorNamesEachLink`: the exact text of API Contract 3
    for two links.
  - `TestMissingRelativeLinksReadsTheSharedExtraction`: for one Markdown
    input holding a present link, an absent link, a self link, a link leaving
    the repository and a code-span link, `MissingRelativeLinks` equals the
    filter of `RelativeLinks` that API Contract 1 states and lists exactly
    the absent and self links.
- `internal/cli/qa_report_test.go`:
  - `TestQAReportAcceptRefusesAnUnresolvedEvidenceLink`: a temporary Git
    repository with `.gitignore` `*.log`, an eligible `pass` report with one
    Results row linking `evidence/consumer-tests.log` (present) and
    `evidence/journeys.txt` (absent): exit 1, empty stdout, the exact stderr
    line of Success Metric 1; after the evidence is saved as `.txt` and the
    links are updated, exit 0 with empty output.
  - `TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink`: through
    `newImplementWorkspace` with a failed `qa` Task, a committed `.gitignore`
    holding `*.log` and an eligible report linking a present ignored `.log`:
    exit 1, the stderr line of Success Metric 3, the Task file unchanged.
- `internal/daemon/qa_evidence_links_test.go` (new), with a wrapper runner
  over `taskFakeRunner` that recognises the repair prompt by its first line
  and runs a scripted `onRepair`:
  - `TestQAEvidenceLinkRepairTurnKeepsTheEvidence`: the fixture commits a
    `.gitignore` holding `*.log`; the QA turn writes a `pass` report with a
    Results row linking `evidence/consumer-tests.log` and writes the file;
    `onRepair` moves it to `evidence/consumer-tests.txt` and rewrites the
    link. The QA Task is accepted, one repair prompt ran, it names
    `evidence/consumer-tests.log`, its line and `ignored`, one
    `evidence_links` event with outcome `repair` was published, and
    `HEAD:<qa>/evidence/consumer-tests.txt` exists.
  - `TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed`: `onRepair`
    changes nothing. The QA Task is `failed` with the exact reason of Success
    Metric 5; exactly one repair prompt ran; one `evidence_links` event with
    outcome `refused` names the link; no `docs: qa report for` commit exists;
    the report and the `.log` remain in the worktree.
  - `TestQAReportWithResolvedLinksRunsNoRepairTurn`: a report whose link
    points at a `.txt` it wrote is accepted with no repair prompt and no
    `evidence_links` event.
- Existing suites: `./internal/spec`, `./internal/speccheck` (the
  `SC-HISTORY-RETAINED` tests), `./internal/daemon`, the `qa-report`,
  settle and partial-policy tests of `./internal/cli`, `./skills`, and
  `-tags docscontract ./internal/docscontract`.
- Sabotages, each with the test whose input makes it observable:
  - `UnresolvedEvidenceLinks` never reports `ignored` (treats every existing
    target as resolved): `TestUnresolvedEvidenceLinksNameTheReason` fails on
    `evidence/a.log`, and `TestQAReportAcceptRefusesAnUnresolvedEvidenceLink`
    fails;
  - `--cached` dropped from the `ls-files` call:
    `TestUnresolvedEvidenceLinksNameTheReason` fails on the forced-add
    `evidence/c.log`;
  - `RelativeLinks` drops reference definitions or keeps code spans:
    `TestUnresolvedEvidenceLinksNameTheReason` fails on the code-span case;
  - the check removed from `requireSettleQAReportEligibility`:
    `TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink` fails;
  - the repair turn removed from `runQAGate`:
    `TestQAEvidenceLinkRepairTurnKeepsTheEvidence` fails (the Task is
    refused);
  - the commit-boundary check removed:
    `TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed` fails (the QA Task
    is accepted and committed);
  - the repair turn run in a loop until the links resolve:
    `TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed` fails on the
    repair prompt count.

## Build Order

1. `RelativeLinks`, `UnresolvedEvidenceLinks`, the refusal error, the
   `qa-report accept` and Settle Command paths, their tests, the qa-report and
   settle guides, the Roundfix settle reference, `CONTEXT.md`, with the
   Roundfix skill version and the fingerprints re-recorded last (task_01).
2. The Daemon repair turn and commit-boundary refusal, their tests, the
   implement guide and the Roundfix implement reference, with the Roundfix
   skill version and the fingerprints re-recorded last (task_02) (depends
   on: 1).
3. The `qa-gate` section with its version (task_03) (depends on: 2).
4. QA gate (task_04) (depends on: 3).

task_02 follows task_01 because it calls `UnresolvedEvidenceLinks` and both
raise the Roundfix skill version and re-record the owned-skill version record
and the Behavior Surface Record; task_03 follows task_02 because both
re-record the owned-skill version record.

Specs 0338 and 0339 of cycle AU deliver before this Spec and are authored in
parallel. Files this Spec may share with them, which a rebase reconciles
through their record commands rather than by hand: the Roundfix skill's
`SKILL.md` pair (version fields), `skills/testdata/owned-skill-versions.json`,
`docs/references/behavior-surfaces.json` and `CONTEXT.md`. This Spec's other
files, `internal/spec/record_links.go`, `internal/spec/qa_evidence_links.go`,
`internal/cli/qa_report.go`, `internal/cli/settle.go` (only
`requireSettleQAReportEligibility`), `runQAGate` in
`internal/daemon/task_engine.go`, the qa-report, settle and implement guides,
the Roundfix settle and implement references and the `qa-gate` skill, are
outside the release plan, Baseline apply and authoring-skill work those two
Specs own. No Verification of this Spec depends on their artifacts.

## Risks & Considerations

- The Daemon asks Git about the Run Worktree, not about the exact paths
  `prepareQAReportCommit` stages. An untracked, unignored file that was dirty
  before the QA step began, or one the stage filter drops as executable or
  behind a symbolic link, resolves here but is left out of the commit. The
  Run Worktree starts clean, and the stage filter already fails a lost path.
- A report that cites evidence only in a code span is not checked, as
  `SC-HISTORY-RETAINED` does not read it either; the `qa-gate` section asks
  for links.
- An adopter whose QA gate keeps `.log` evidence will see its first QA after
  the upgrade take the repair turn; the release notes name the `.txt` rule.

## Glossary

- adds: **Unresolved Evidence Link** — a QA Report link the committed tree would not hold
- changes: **QA Report** — settles only when its local links resolve

## Decisions

- One check, `spec.UnresolvedEvidenceLinks`, for the Daemon, `qa-report
  accept` and the Settle Command; `roundfix archive` unchanged.
- One extraction, `spec.RelativeLinks`, for the QA check and
  `SC-HISTORY-RETAINED`.
- Git decides resolution: a tracked or unignored existing target resolves,
  an existing ignored one is `ignored`, any other is `missing`.
- One repair turn in the same Agent Session; a remaining link refuses the QA
  Report commit.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found no page deciding the rule, and a staging digest that
  argues against `git add -f`. Exa: the Git manuals for `git-check-ignore`
  and `gitignore`, which state that tracked files are not subject to ignore
  rules.
- Three implementation Tasks in series and a QA gate.
