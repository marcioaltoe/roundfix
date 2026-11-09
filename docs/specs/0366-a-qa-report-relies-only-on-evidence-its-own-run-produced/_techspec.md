---
spec: 0366-a-qa-report-relies-only-on-evidence-its-own-run-produced
prd: _prd.md
created: 2026-11-09
---

# A QA Report relies only on evidence its own Run produced — Technical Spec

## Executive Summary

`internal/spec` gains one check, `FindBorrowedEvidence`, that reads a QA
Report and names its Borrowed Evidence by four reasons: `reused` (a linked
file under the Spec's `qa/evidence/` whose bytes the audited head already
holds at that path), `imported` (the same, against the QA Report Import's
copy), `outside the worktree` (a `user_flow_executables` entry outside the
report's worktree) and `names another Run Worktree <path>`. `qa-report accept`
and the Settle Command apply it after the evidence-link rule, as ADR-0351
placed that rule; the Daemon applies it after the QA turn with the import and
Run Branch it knows, gives one repair turn and refuses at the QA Report
commit boundary. The qa-gate skill writes each pass's evidence under its own
directory and lists the executables it ran. The trade-off accepted: bytes
and paths decide, never a commit parsed out of an adopter's evidence format,
so a regenerated file that still pins an old commit inside its own text is
caught only when it also names another Run Worktree or the skill is followed.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the optional `user_flow_executables` key is read only when
  present. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the check reads local files and the local Git repository.
  Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0376 (new) records the rule.
  ADR-0351: "A QA Report settles only when its local links resolve in the
  committed tree", whose three paths and repair turn this rule joins, after
  that check. The carry-forward and import decisions are not listed: this
  rule reads their results (a Carried Row's status, the files the QA Report
  Import copied) and changes neither. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BA rules answer of 2026-11-09 quoted in the authorization record; no
  Makefile, workflow, Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0366-a-qa-report-relies-only-on-evidence-its-own-run-produced/_authorization.md`;
  bounded files: `.agents/skills/qa-gate/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/settle.md`, `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read at `2e59aebe`:

1. `runQAReportCommand` (`internal/cli/qa_report.go`) applies
   `spec.QAReportEligibility`, then `spec.RequireResolvedEvidenceLinks`, and
   exits `0`; nothing compares evidence with the audited head.
2. `requireSettleQAReportEligibility` (`internal/cli/settle.go`) applies the
   same two checks and wraps a refusal as `<error> (report <path>)`.
3. `runQAGate` (`internal/daemon/task_engine.go`) imports a prior pass
   (`importPriorQAPass`, which copies its report and evidence into the
   worktree), seeds the report, appends `Seeded QA Report: <path>` to the
   prompt, runs the QA turn, gives one repair turn for unresolved links
   (`qaEvidenceLinkRepairPrompt`), and at the commit boundary refuses a
   Secret Match, a non-image binary or an unresolved link before
   `commitPreparedQAReport`. The prompt names the Run Worktree branch
   (`writeQACheckoutFacts`, `internal/agent/spec_prompt.go`) but not its path
   or the audited head.
4. The qa-gate skill stores evidence under
   `qa/evidence/<YYYY-MM-DD>-<run-slug>/`: two passes on one date with one
   slug share a directory, as Quaestor's `2026-11-06-gate` did.
5. Characterized on 2026-11-09: an uncommitted `pass` report linking two
   files an earlier commit added under `qa/evidence/2026-11-06-gate/`, one
   naming a registered `roundfix/run-old` worktree's `.venv/bin/tool`, makes
   `bin/roundfix qa-report accept` exit `0`.

## System Architecture

- `internal/spec/qa_evidence_binding.go` (new): `FindBorrowedEvidence`,
  `BorrowedEvidenceError`, `RequireOwnRunEvidence`, the reason constants and
  `BorrowedEvidenceOptions`. It reads links through `RelativeLinks`, Git
  through the package's existing reader `archiveGitRead`, and the Results
  table the way `qaReportResultsCounts` reads it. It imports no new package.
- `internal/cli/qa_report.go` `runQAReportCommand` and
  `internal/cli/settle.go` `requireSettleQAReportEligibility`: call
  `spec.RequireOwnRunEvidence` with `RunBranchPrefix: store.RunBranchPrefix`
  after the evidence-link check accepts.
- `internal/daemon/task_engine.go` `runQAGate`: the prompt facts, the shared
  repair turn and the commit-boundary refusal; new helpers
  `borrowedQAEvidence` and `qaBorrowedEvidenceRepairPrompt`.
- Unchanged by design: `spec.QAReportEligibility`, `readQAReport`,
  `internal/spec/archive.go`, `internal/speccheck`, `internal/agent`,
  `importPriorQAPass`, `copyPriorQAPass`, `settleQAVerdict`,
  `mechanicalQAReportContent` and `recordQAEvidenceSnapshots`.
- Docs: the qa-gate skill; `docs/user-guide/commands/qa-report.md`,
  `settle.md` and `implement.md`; the Roundfix `settle` and `implement`
  references; `CONTEXT.md`.

## Implementation Design

### Interfaces

```go
// internal/spec/qa_evidence_binding.go
const (
	BorrowedEvidenceReused          = "reused"
	BorrowedEvidenceImported        = "imported"
	BorrowedEvidenceOutsideWorktree = "outside the worktree"
	BorrowedEvidenceMalformed       = "not a list of paths"
	BorrowedEvidenceForeignWorktree = "names another Run Worktree" // rendered with " <path>"
)
type BorrowedEvidence struct{ Line int; Subject, Reason string }
type BorrowedEvidenceOptions struct {
	RunBranchPrefix string   // required: the Run Branch prefix, roundfix/run-
	RunBranch       string   // "" derives the report worktree's own branch
	AuditedHead     string   // "" derives it (API Contract 1)
	ImportCommit    string   // the QA Report Import's commit, or ""
	ImportedPaths   []string // repository-relative paths it copied
}
func FindBorrowedEvidence(ctx context.Context, reportPath string, options BorrowedEvidenceOptions) ([]BorrowedEvidence, error)
type BorrowedEvidenceError struct{ Evidence []BorrowedEvidence }
func (err BorrowedEvidenceError) Error() string
func RequireOwnRunEvidence(ctx context.Context, reportPath string, options BorrowedEvidenceOptions) error
```

```go
// internal/daemon/task_engine.go
func borrowedQAEvidence(ctx context.Context, plan TaskPlan, prior priorQAPass, auditedHead string) (string, []spec.BorrowedEvidence, error)
func qaBorrowedEvidenceRepairPrompt(reportPath, workDir, auditedHead, evidenceDir string, evidence []spec.BorrowedEvidence) string
```

### Data Models

No schema change. A QA Report may carry the optional frontmatter key
`user_flow_executables`, a YAML list of paths; `readQAReport` and the
`QAReport` struct do not read it.

### API Contracts

1. API Contract: `FindBorrowedEvidence` scope. The worktree root is
   `git rev-parse --show-toplevel` from the report's directory; outside Git
   it returns no item and no error. `RunBranchPrefix` empty is an error. The
   audited head is `options.AuditedHead`, else `HEAD` when
   `git status --porcelain -- <report>` prints anything (untracked or
   modified), else the first parent of
   `git log --diff-filter=A --format=%H -1 -- <report>`; an unborn `HEAD` or
   a root commit leaves no audited head, and then `reused` is not checked.
   Evidence links are the report's `RelativeLinks` whose target lies under
   `<spec dir>/qa/evidence/` and exists; a directory target is checked file
   by file, each with subject `<destination>/<path below it>`. A link is
   exempt when its line is the Results line of a row whose status starts
   with `carried (` or lies in that row's `### <row-id>` block (until the
   next heading of depth 3 or less). Items are returned in report-line
   order; a read or Git failure is an error, never an item.
2. API Contract: the four reasons, checked per file in this order, at most
   one of the first two per file:
   - `reused`: the audited head holds a blob at the file's repository path
     with the same bytes as the file;
   - `imported`: `ImportCommit` is set, the path is in `ImportedPaths`, and
     `<ImportCommit>:<path>` holds the same bytes;
   - `names another Run Worktree <path>`: for every file not named by the
     first two, and for the report itself, the content holds `<root>/` for a
     registered worktree (`git worktree list --porcelain`) whose branch
     starts with `RunBranchPrefix`, is not the own Run Branch nor starts with
     the own Run Branch followed by `-` (its Task branches), and whose path
     differs from the report's root; `<root>` is the registered path and, if
     different, its symbolic-link resolution. The own Run Branch is
     `RunBranch`, else the report worktree's branch when it starts with
     `RunBranchPrefix`. Subject: the link destination, or the report's file
     name with the line of its first match. One item per file and worktree;
   - `outside the worktree`: each `user_flow_executables` entry, cleaned and
     joined to the root when relative, that is inside neither the root nor
     its symbolic-link resolution, compared lexically without resolving the
     entry; subject: the entry as written, line: the entry's report line.
     A key that is not a list of strings gives one item with subject
     `user_flow_executables`, reason `not a list of paths`, and the key's
     line. An absent key gives none.
3. API Contract: `BorrowedEvidenceError.Error()` is
   `QA Report evidence does not belong to this Run: <subject> (line <n>, <reason>)`
   for each item, joined by `; `. `RequireOwnRunEvidence` returns it when
   items exist, the error when the check failed, and nil otherwise.
4. API Contract: `qa-report accept`. After `RequireResolvedEvidenceLinks`
   accepts, `RequireOwnRunEvidence(ctx, reportPath,
   spec.BorrowedEvidenceOptions{RunBranchPrefix: store.RunBranchPrefix})`;
   a refusal prints `roundfix: QA Report eligibility refused: <error>` and
   exits `1` with empty stdout. The usage text stays.
5. API Contract: Settle Command. `requireSettleQAReportEligibility` applies
   the same call after the link check accepts and wraps a refusal as
   `<error> (report <path>)`.
6. API Contract: Daemon repair turn. In `runQAGate`, where the evidence-link
   repair runs today, `borrowedQAEvidence` reads the newest report with
   `RunBranchPrefix: store.RunBranchPrefix`, `RunBranch:
   plan.RunWorktree.Branch`, `AuditedHead: auditedHead` and, when a pass was
   imported, its commit and files. When links or Borrowed Evidence exist,
   exactly one repair turn runs. Its prompt is the existing link prompt
   unchanged when links exist, followed by `qaBorrowedEvidenceRepairPrompt`
   when Borrowed Evidence exists, whose lines are
   `QA evidence does not belong to this Run.`, `Report: <path>`,
   `Run Worktree: <work dir>`, `Audited head: <sha>`,
   `Evidence directory: <dir>/`, one `- <subject> (line <n>, <reason>)` per
   item, and
   `Run each affected row again at the audited head with executables built in this Run Worktree, write its evidence under the evidence directory, update its links and user_flow_executables, and never copy, replay or cite evidence an earlier Run left; never change a verdict or a row the new evidence does not support.`
   It publishes the `evidence_binding` QA Run Event (`phase`
   `evidence_binding`, `outcome` `repair`, `items` of `subject`, `line`,
   `reason`) and prints
   `QA evidence does not belong to this Run: <n> item(s); returning the QA turn once.`
7. API Contract: Daemon commit boundary. After the link check finds
   nothing, remaining Borrowed Evidence settles the QA Task `failed` with
   `spec.BorrowedEvidenceError{Evidence: items}.Error()`, publishes
   `evidence_binding` with outcome `refused`, prints
   `QA Report refused: <n> borrowed evidence item(s)`, and returns before
   `commitPreparedQAReport`. A Git or read error is an infrastructure error.
8. API Contract: QA prompt facts. After the `Seeded QA Report:` lines,
   `runQAGate` appends `Run Worktree: <work dir>`, `Audited head: <sha>`,
   `Evidence directory: <dir>/` and
   `Build and run every executable a row exercises inside this Run Worktree, list their paths in user_flow_executables, write this pass's evidence under the evidence directory, and never copy, replay or cite evidence or an executable an earlier Run left.`
   `<dir>` is the seeded report's directory joined with `evidence` and the
   report's file name without `.md`.
9. API Contract: qa-gate skill. A section `### Evidence belongs to this Run`
   after `### Evidence links resolve in the committed tree`; step 4 stores
   evidence under `qa/evidence/<report name>/`; the report template lists
   `user_flow_executables: []` after `user_flow_binary`.

### Surface Transcripts

None. The new stderr lines are asserted byte for byte in
`internal/cli/qa_report_test.go` over temporary repositories; the Daemon text
is asserted in `internal/daemon` with the scripted QA Agent.

### Invariants

1. A report with no evidence link under `qa/evidence/`, no
   `user_flow_executables` key and no other Run's worktree path gets no item.
2. A link the existing rule calls unresolved is never also an item: the new
   check runs only after that rule accepts, and missing targets are skipped.
3. A Carried Row's evidence is never an item.
4. `roundfix archive` and every archived report are untouched.
5. The link repair prompt keeps its bytes when no Borrowed Evidence exists.

### Declared breaks

- A report that links evidence an earlier commit holds unchanged, or names
  another Run's Run Worktree, was accepted and is now refused on all three
  paths. No existing test or fixture pins such a report (searched:
  `grep -rln '](evidence/' internal --include='*_test.go'` finds only the
  evidence-link tests, whose evidence is untracked).
- The qa-gate skill's evidence directory moves from
  `qa/evidence/<YYYY-MM-DD>-<run-slug>/` to `qa/evidence/<report name>/`;
  no test pins the old text.

## Coverage Map

- Goal 1, Story 1, Core Features 1-2, Success Metrics 1-5 → API Contracts
  1-5, Build Order 1.
- Goal 2, Story 3, Core Feature 3, Success Metrics 6-7 → API Contracts 6-8,
  Build Order 2.
- Goal 3, Core Feature 4 → API Contract 9, Build Order 3; guides and
  references → Build Orders 1-2.
- Goal 4, Story 2, Success Metric 8 → Invariants 1-5, API Contract 2
  (lexical executables), Testing Approach.

Deviations from the PRD: None.

## Absence Evidence

1. `FindBorrowedEvidence` (the shared check) — command: `grep -rn 'BorrowedEvidence\|user_flow_executables\|worktree list' internal/spec internal/cli/qa_report.go internal/cli/settle.go`; result: no match; the only evidence checks are `RequireResolvedEvidenceLinks` and, in the Daemon, the secret and binary scans.
2. The Daemon repair and refusal for Borrowed Evidence — command: `grep -n 'evidence_binding\|Audited head\|Evidence directory' internal/daemon/task_engine.go internal/agent/spec_prompt.go`; result: no match; the prompt names the Run Worktree branch only.
3. The qa-gate section — command: `grep -n 'Evidence belongs to this Run\|user_flow_executables' .agents/skills/qa-gate/SKILL.md`; result: no match.

## ADR Alternatives

- ADR-0376: the design takes none of its rejected options: no prose-only
  rule, no commit parsed from evidence, no fixed directory required, no flag
  on the user checkout, no resolved executable links.
- ADR-0351: rejected warning and leaving the check to `spec check`; this
  rule likewise refuses before the commit and keeps one repair turn, and it
  leaves archive unchanged as that ADR did.

## Integration Points

The local Git repository only: `rev-parse`, `status --porcelain`, `log`,
`cat-file`/`ls-tree`, `worktree list --porcelain`, through
`archiveGitRead`. No network.

## Testing Approach

- `internal/spec/qa_evidence_binding_test.go` (new), each test calling
  `t.Parallel()` first, over temporary repositories built with `gittest`:
  - `TestBorrowedEvidenceRefusesAnEarlierRunsEvidenceDirectory`: the
    Quaestor layout (an earlier commit holds
    `qa/evidence/2026-11-06-gate/head-exclusive.execution.txt` and
    `replay.py.txt`; a later commit; an uncommitted report links both) →
    two `reused` items; a new untracked file and a tracked file with changed
    bytes in the same report → no item.
  - `TestBorrowedEvidenceReadsACommittedReportAgainstItsOwnCommit`: report
    and evidence committed together → none; report committed later linking
    evidence an earlier commit added → `reused`.
  - `TestBorrowedEvidenceExemptsACarriedRow`: the link on a `carried (...)`
    Results line and in its `### <row-id>` block → none; the same link on a
    `pass` row → `reused`.
  - `TestBorrowedEvidenceNamesAnExecutableOutsideTheWorktree`: entries
    `./bin/tool`, `.venv/bin/tool` and `.venv/bin/python` (a symbolic link to
    a file outside the worktree) → none; `/elsewhere/bin/tool` and
    `../other/bin/tool` → `outside the worktree`; a scalar → `not a list of
    paths`; an absent key → none.
  - `TestBorrowedEvidenceNamesAnotherRunWorktree`: worktrees on
    `roundfix/run-old`, `roundfix/run-own-task_01` and a non-Run branch;
    the report in the `roundfix/run-own` worktree; a fresh evidence file
    naming `<run-old>/.venv/bin/tool` → one item; paths in the own worktree,
    the own Task worktree and the non-Run worktree → none; the path in the
    report body → an item on its line.
  - `TestBorrowedEvidenceNamesImportedEvidence`: `ImportCommit` and
    `ImportedPaths` with an unchanged copy → `imported`; a changed copy →
    none.
  - `TestBorrowedEvidenceErrorNamesEachItem`: the exact error text.
  - `TestBorrowedEvidenceOutsideGitChecksNothing`.
- `internal/cli/qa_report_test.go`:
  `TestQAReportAcceptRefusesBorrowedEvidence` and
  `TestSettleRefusesAQAReportWithBorrowedEvidence`: exact stderr, then
  accepted after the evidence is regenerated under
  `qa/evidence/qa-report-2026-11-09/`.
- `internal/daemon/qa_evidence_binding_test.go` (new), with a wrapper runner
  over `taskFakeRunner` that recognises the repair prompt by its first line,
  as `qa_evidence_links_test.go` does:
  `TestQABorrowedEvidenceRepairTurnRegeneratesTheEvidence`,
  `TestQABorrowedEvidenceAfterRepairSettlesFailed`,
  `TestQAPromptNamesTheRunWorktreeAuditedHeadAndEvidenceDirectory`.
- Existing seams that must keep passing unchanged:
  `TestUnresolvedEvidenceLinksNameTheReason`,
  `TestQAReportAcceptRefusesAnUnresolvedEvidenceLink`,
  `TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink`,
  `TestQAEvidenceLinkRepairTurnKeepsTheEvidence`,
  `TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed`,
  `TestQAReportWithResolvedLinksRunsNoRepairTurn`,
  `TestQAGateCarriesARowFromAnUnintegratedFailedPass` and
  `TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte`.

## Build Order

1. The shared check, `qa-report accept` and the Settle Command, with their
   tests, the qa-report and settle guides, the Roundfix settle reference,
   `CONTEXT.md`, the Roundfix skill version and the Behavior Surface Record.
2. The Daemon's prompt facts, repair turn and commit-boundary refusal, with
   their tests, the implement guide and the Roundfix implement reference
   (depends on: 1).
3. The qa-gate skill section, evidence directory and template key (depends
   on: 2).
4. Final QA gate (depends on: 3).

## Risks & Considerations

- A regenerated file that pins an old commit in its own text and names no
  other Run Worktree passes the check; the skill's rule and the prompt's
  audited head address it, and ADR-0376 records why no format is parsed.
- An adopter whose qa-gate copy predates this Spec may still write into an
  earlier pass's dated directory; a byte-identical rewrite is then `reused`
  and the repair prompt names the new directory. A changed rewrite passes.
- Rebase order: 0365 merges first and may raise the Roundfix skill version
  and re-record the Behavior Surface Record; 0367 follows. Every shared
  record is regenerated by its command on the Task's own base.
- The new tests use `git worktree add`, which works on every platform; no
  file carries a build constraint.

## Glossary

- adds: **Borrowed Evidence**
- changes: **QA Report**

## Decisions

- One shared check on the three settlement paths, after the link rule. See
  ADR-0376.
- Bytes against the audited head decide `reused`; no evidence format is
  parsed.
- `user_flow_executables` optional and lexical; only other Runs' Run
  Worktrees named; Carried Rows exempt; archive unchanged.
- Three implementation Tasks and the QA gate, in series, mirroring the
  evidence-link Spec that this rule follows.
