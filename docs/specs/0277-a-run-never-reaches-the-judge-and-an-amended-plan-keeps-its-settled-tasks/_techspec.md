---
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
prd: _prd.md
created: 2026-10-17
---

# A Run never reaches the judge, and an amended plan keeps its settled Tasks — Technical Spec

## Executive Summary

Four slices extend existing seams; no package is new. The Daemon gives every
Agent Session and Verification command of a Run the Run Marker
`ROUNDFIX_RUN_ID`, the judge in `internal/judge` refuses to send while it is
set, and `TYPESAFE_API_KEY` joins the variables a Run never receives
(ADR-0283). The bootstrap serialization test gives each sibling its own
release FIFO, and every Unix-only test moves into a `unix`-constrained file,
proved by a new Repository Contract Test that compiles every test for Windows
(ADR-0285). Task Carry-Forward in `internal/cli` compares a settled Task's
Binding Text instead of whole planning files (ADR-0284). The history cleanup in
`internal/spec` reads `superseded_by` and `supersedes` in every form adopters
write.

The trade-off accepted is a narrow guarantee over a broad one. Roundfix cannot
stop an Agent runtime from re-reading the user's shell startup, so it
guarantees that its own judge sends nothing from a Run rather than that no
command of a Run ever sees a key. Carry-forward binds a Task to whole sections
its References name, not to numbered items, so an amendment inside a named
section still re-runs the Task.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or schema is
  minted; `ROUNDFIX_RUN_ID` already names a Run for the notification command,
  and the carry-forward candidate gains one optional list. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the judge's credential lookup and its
  one HTTP request path gain a refusal before either runs; no request,
  credential or endpoint is added, and every test uses an injected or local
  transport. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0283 governs Build Order 1: "The
  judge sends nothing while the Run Marker is set." ADR-0285 governs Build
  Order 2: "A test that needs Unix lives in a file that says so." ADR-0284
  governs Build Order 3: "A settled Task binds only the planning text it
  names." ADR-0273, amended: "The judge's keys never enter a Run." ADR-0271:
  "A test bounds by real time only the budget it asserts". ADR-0266: "A Proven
  Signal is the only signal Roundfix code or tests send". ADR-0250: "The
  coverage record lists tests for every release platform." ADR-0253: "Every
  Repository Contract Test runs on main and before a release". ADR-0252: "The
  selective gate runs the repository contracts a change makes relevant".
  ADR-0170: "A Task already completed on its target is nothing to carry".
  ADR-0158: "A budget-expired Run settles BudgetExceeded and stays
  recoverable". ADR-0282: "The cleanup writes the lineage a successor lacks."
  ADR-0278: "A retirement deletes." ADR-0254: "Refused Units never block".
  ADR-0189: "A version names one content." ADR-0280: "A guide Task records its
  own fingerprint." ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0276:
  "A row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "Each claim carries its proof as
  a section." Set no obligation here, because this Spec changes none of their rules: ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0261 and ADR-0277 (history families, Archive Records, the History Full Tag and legacy folder conversion); ADR-0258 and ADR-0260 (reduced history and live citations); ADR-0164 (Run budget renewal); ADR-0233 (version regeneration at merge); ADR-0240 (the QA partial policy); ADR-0281 (rebuilt executables in tests); ADR-0251, ADR-0262 and ADR-0269 (lenient legacy reading, cancellation and archive refusals); ADR-0165 and ADR-0272 (review parks and records); ADR-0167 and ADR-0229 (QA rows and the operator archive); ADR-0169, ADR-0196 and ADR-0279 (the review's merge base, finding validation and dismissals); ADR-0237 (the Delivery Retry, whose carry-forward rule ADR-0284 refines without changing the retry). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH grant of 2026-10-16 ("Governed Paths declarados") and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/reconcile.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `docs/references/coverage-record.json`,
  `internal/baseline/repository_test.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-17 with `bin/roundfix` built at `23ab0904`.

- Judge key, Roundfix side. `judge.WithoutKeyVariables`
  (`internal/judge/environment.go:10`) removes `Questions.KeyVariables()`,
  which is `ROUNDFIX_OPENROUTER_JUDGE_API_KEY`, `ROUNDFIX_OPENROUTER_API_KEY`
  and `ROUNDFIX_TYPESAFE_API_KEY` (`internal/judge/questions.go:145-160` and
  the two transports of `questions.json`). `acpxCommandEnv` applies it to
  every acpx invocation (`internal/agent/codex_spawn.go:139`) and
  `ExecVerifier.Verify` to every Verification command
  (`internal/daemon/daemon.go:131`);
  `TestNoAgentSessionReceivesAJudgeKeyVariable` and
  `TestVerificationReceivesNoJudgeKeyVariable` prove both. `TYPESAFE_API_KEY`
  is not removed. The Queue Owner passes its whole environment, with `HOME`
  replaced, to the item binary (`deliveryCommandEnvironment`,
  `internal/cli/deliver_workflow.go:1675` and `:2161`), and the Daemon needs
  it: a Light Tier session gets the shared key's value under the
  implementation Stage Key name (`internal/agent/light_session.go:35`). The
  owner itself makes no judge call; the only callers are
  `internal/cli/spec_judge.go:134`, `internal/cli/archive.go:478` (the plan
  without `--offline`) and `internal/cli/history.go:498`, each reading the
  keys from its own environment. The Worktree Bootstrap (`sh -c` at
  `internal/worktree/worktree.go:3096`), Git hooks and the QA format command
  inherit the Daemon's environment, as ADR-0273 decided.
- Judge key, runtime side. The leak happens after Roundfix's boundary. Codex
  `0.163.0`, the QA Agents' runtime in both incidents, creates a shell
  snapshot per session: its zsh script sources `$ZDOTDIR/.zshrc` or
  `~/.zshrc`, captures `export -p` with only `PWD` and `OLDPWD` excluded, and
  sources the snapshot before each command (Acceptance evidence). On the
  maintainer's machine `~/.zshrc` sources a secrets file, and the one snapshot
  under `~/.codex/shell_snapshots/` holds one line naming
  `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` (counted by name on 2026-10-17; no value
  was read). So the QA Agent's `roundfix archive --plan` found the key and
  sent two archive-value requests in each incident. `ROUNDFIX_RUN_ID` is set
  today only for the notification command (`internal/notify/notify.go:179`).
- Bootstrap test. `bootstrapSerializationRound`
  (`internal/worktree/worktree_test.go:471`) creates one release FIFO per
  round, and every sibling's command runs `read release < <that FIFO>`
  (`:486`). `tryWriteBootstrapRelease` opens it write-only and non-blocking,
  writes `release`, then closes it (`:667-671`). The production permit
  (`internal/worktree/worktree.go:3088-3091`) admits the next sibling as soon
  as the released one exits. When the released sibling reads, exits and
  returns its permit before the helper's `Close`, the next sibling prints
  `bootstrap-started` and opens the FIFO while the old writer is still open;
  the open succeeds, and when the writer closes its `read` gets end-of-file,
  `test` fails, and the sibling ends before its own release, which is the
  "watched work ended" of run 39044083565. Nothing but scheduling decides the
  order.
- Windows. `GOOS=windows go vet ./...` exits 1 in about 4 s warm: one error per
  package in `internal/agent/acpx_integration_test.go:375` (`syscall.Kill`),
  `internal/baseline/repository_test.go:359` (`syscall.Mkfifo`) and
  `internal/cli/deliver_item_probe_report_test.go:50` (`syscall.Stat_t`).
  Compiling the three packages' tests with every error listed finds 18 in
  seven files: also `deliver_item_probe_report_test.go:126-129`
  (`syscall.Getpgrp`), `detach_fixture_group_other_test.go:14` (its `!darwin &&
  !linux` constraint admits Windows), `detach_test.go:491,517`,
  `implement_detach_teardown_test.go:62-265` and `implement_test.go:1599,1897,1938`.
  The same check on the tree of `e5da154a`, before Spec 0275, fails in the
  same three packages. The coverage record lists those files' tests for
  Windows.
- Carry-forward. `inspectCarryForwardCandidate`
  (`internal/cli/carryforward.go:538`) compares every input of
  `spec.CarryForwardInputs` (`internal/spec/spec.go:980`: the Task file,
  `_prd.md`, `_techspec.md`, `_tasks.md`, `AGENTS.md`, `CONTEXT.md`, the
  `implement-task` skill and the Task's Context paths) between the settlement
  parent and the staged target with `git diff --quiet`; any difference is a
  moved input, and `inspectCarryForwardsForRun` refuses the whole set with
  `Task <id> declared input(s) moved: <paths>`. The reconcile text lists
  `inputs-moved` and `moved-inputs` (`internal/cli/reconcile.go:1323-1324`).
- ADR reference forms. `historyADRSuccessor`
  (`internal/spec/history_entries.go:250-264`) accepts only `superseded_by:
  ADR-NNNN`, so `0064-application-routes-follow-resource-semantics` yields no
  successor and the ADR is refused. The "already names it" test
  (`:133-134`) looks for `ADR-NNNN` or the retired file's basename with `.md`,
  so `supersedes: 0002-better-auth-protocol-route-exception` would not count
  and a lineage line would be appended.

## System Architecture

- Judge (`internal/judge`, `internal/agent`, `internal/daemon`,
  `internal/cli`). `internal/judge` gains the marker name and the refusal in
  `Run` and `AdviseArchive`; `WithoutKeyVariables` also removes
  `TYPESAFE_API_KEY`. The acpx runner appends the marker to the environment of
  each invocation for a request with a Run ID; `VerifyRequest` carries the Run
  ID to `ExecVerifier`. The three CLI callers pass the marker from their own
  environment.
- Tests (`internal/worktree`, `internal/agent`, `internal/baseline`,
  `internal/cli`, `internal/verifyselect`). Test-only changes, plus the
  coverage record and a new `repocontract` test in `internal/verifyselect`,
  a guarded spawning package.
- Carry-forward (`internal/spec`, `internal/cli`). A new
  `internal/spec/carry_forward_binding.go` extracts Binding Text;
  `inspectCarryForwardCandidate` uses it for the three planning files.
- History (`internal/spec`). `historyADRSuccessor` and the naming test of
  `PlanHistoryKind` read the forms.

## Implementation Design

### Interfaces

```go
// internal/judge
const RunMarkerVariable = "ROUNDFIX_RUN_ID"
// Request and ArchiveAdviceRequest gain: RunID string
func WithoutKeyVariables(environ []string) []string // also drops TYPESAFE_API_KEY

// internal/daemon
// VerifyRequest gains: RunID string

// internal/spec
// CarryForward gains: AmendedInputs []string `json:"amendedInputs"`
func PlanningBinding(file string, planning []byte, taskID string, task []byte) (binding string, whole bool, err error)
```

### Data Models

No new file format. The carry-forward JSON object gains `amendedInputs`, an
array that is empty when nothing was tolerated.

### API Contracts

1. API Contract: the Run Marker. Every acpx invocation the runner makes for a
   prepare, ensure or prompt request whose Run ID is not empty receives
   `ROUNDFIX_RUN_ID=<Run ID>` after the inherited environment, replacing any
   inherited value; every `ExecVerifier.Verify` call whose request carries a
   Run ID runs its command with `ROUNDFIX_RUN_ID=<Run ID>` the same way. The
   engine passes the Run ID for a Task's Verification, the probe on the
   unchanged tree and the repository Verification at settlement. Neither
   environment holds a Judge Key Variable or `TYPESAFE_API_KEY`, and every
   other inherited variable is kept, including the implementation Stage Key a
   light session receives.
2. API Contract: the silent judge. `judge.Run` and `judge.AdviseArchive`
   check `RunID` first: when it is not empty they select no transport, read no
   key, read and write no Judge Log, send no request, and set `Skipped` to
   `inside Run <run-id>, named by ROUNDFIX_RUN_ID`. `roundfix spec judge`
   prints `Judge: skipped: inside Run <run-id>, named by ROUNDFIX_RUN_ID;
   <n> judgment(s) not asked`; `roundfix archive <slug> --plan` and the advice
   of `roundfix history sanitize` print each candidate as `no advice (inside
   Run <run-id>, named by ROUNDFIX_RUN_ID)`. Each caller reads
   `ROUNDFIX_RUN_ID` from the environment its command was given, so a test
   that injects an environment is unaffected by its own process's marker. An
   unset or empty marker leaves every behavior of today; `--offline` keeps its
   own lines.
3. API Contract: release FIFOs. `bootstrapSerializationRound` creates one
   release FIFO per sibling and each sibling's command reads only its own;
   the release helper writes only to the FIFO of the sibling whose start it
   observed. The lock collision check, the six rounds, the four siblings and
   every failure message keep their text.
4. API Contract: Unix-only tests. Every top-level test and helper that uses a
   symbol absent from Windows moves, unchanged and under its name, into a file
   whose first constraint line is `//go:build unix` (or narrower), and
   `detach_fixture_group_other_test.go` gains `unix`. No Windows stub that
   skips or fakes a Unix behavior is added. The coverage record is re-recorded
   with its sanctioned command and lists each moved test for `darwin` and
   `linux` only.
5. API Contract: the Windows compile contract.
   `TestEveryTestCompilesForWindows`, in a new `repocontract` file of
   `internal/verifyselect` with the default Contract Relevance, runs `go vet
   -tags 'repocontract docscontract' ./...` from the repository root with
   `GOOS=windows` added to its environment, and fails with the command's
   combined output when it exits non-zero.
6. API Contract: Binding Text. `PlanningBinding` returns, for `_tasks.md`,
   the Task's graph node (`id`, `file` and `needs`, in order) and its
   projection table row, and `whole` when either is missing. For `_prd.md` and
   `_techspec.md` it reads the Task file's `## References` bullets that begin
   with the file name in backticks and an arrow; each `;`-separated unit, cut
   at its first `:` and stripped of a trailing run of numbers, ranges, commas
   and `and`, names every heading of level 2 to 4 outside fenced blocks whose
   text equals it, case-insensitively, or its plural (`s`, or `y` to `ies`).
   The binding is the text of each named section, from its heading to the
   next heading of the same or a higher level, in file order. It is `whole`
   when the Task has no bullet for the file or any unit names no heading.
7. API Contract: carry-forward over amendments. For a moved input that is the
   `_prd.md`, `_techspec.md` or `_tasks.md` of the Task's Spec,
   `inspectCarryForwardCandidate` computes `PlanningBinding` at the settlement
   parent and on the staged target, both from the Task file at the parent;
   when neither is `whole` and the two are equal, the path leaves
   `MovedInputs` and joins `AmendedInputs`. Every other input, and a binding
   that differs or is `whole`, stays moved with today's refusal text and
   recovery. The reconcile text prints `  amended-inputs: [<paths>]` after
   `  moved-inputs:`, and JSON carries `amendedInputs`. The rule serves the
   Reconcile Command, the Implement Command's Preflight and a Delivery Retry.
8. API Contract: ADR reference forms. `superseded_by` resolves when it is
   `ADR-<n>` or `<n>`, to the accepted `docs/adr/<n>-*.md`, or `<n>-<slug>`,
   `<n>-<slug>.md` or a path whose last element is `<n>-<slug>.md`, to the
   accepted `docs/adr/<n>-<slug>.md`. A successor already names the retired
   ADR when its body does so today or when its front matter `supersedes`, a
   string or a list of strings, holds `ADR-<n>`, `<n>`, `<stem>`,
   `<stem>.md` or a path ending in `/<stem>.md`, where `<n>` and `<stem>` are
   the retired file's number and name without `.md`. The lineage line, the
   refusal text and every other unit are unchanged.
9. API Contract: everything else stays. The offline archive plan, the Jev
   Ceiling, every judge answer, the bootstrap permit, the environment of Git
   hooks, the QA format command and the Worktree Bootstrap, every carry-forward
   refusal text, and every existing test name keep their behavior.

### Surface Transcripts

1. Surface Transcript: the archive plan inside a Run, for a Spec whose only
   candidate is one file under `references/` that no Live Document cites, with
   no Judge Key Variable set.

   ```transcript
   $ ROUNDFIX_RUN_ID=<run-id> roundfix archive <slug> --plan
   stdout:
   archive plan for <slug>: removes <count> file(s) (<bytes> bytes) and writes docs/history/specs/<slug>.md
   core <count> file(s): <files>
   evidence 0 file(s) (0 bytes) under qa/evidence/
   candidate references/<file> <bytes> bytes: no advice (inside Run <run-id>, named by ROUNDFIX_RUN_ID)
   promote with: roundfix archive <slug> --promote <path>
   stderr:
   exit: 0
   ```

The `spec judge` line and the reconcile line are fixed by API Contracts 2 and
7 and asserted by the named tests.

### Invariants

```text
1. No Agent Session and no Verification command of a Run receives a Judge Key Variable or TYPESAFE_API_KEY.
2. Every Agent Session and Verification command of a Run receives ROUNDFIX_RUN_ID with its Run ID.
3. While ROUNDFIX_RUN_ID is set and not empty, the judge sends no request, whatever keys are set.
4. A bootstrap release reaches only the sibling it was written for.
5. Every test of the repository compiles for Windows.
6. A Task whose Binding Text changed, or whose non-planning input moved, refuses the whole set.
7. A successor that names its predecessor in any accepted form gets no lineage line.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Story 1, Core Features 1-3, Success Metric 1 → the Run Marker in the acpx runner and `ExecVerifier`, the refusal in `judge.Run` and `judge.AdviseArchive`, `WithoutKeyVariables` (API Contracts 1, 2; Surface Transcript 1)
- Goal 2, User Story 2, Core Feature 4, Success Metric 2 → `bootstrapSerializationRound` (API Contract 3)
- Goal 3, User Story 3, Core Features 5 and 6, Success Metric 3 → the moved tests, the coverage record, `TestEveryTestCompilesForWindows` (API Contracts 4, 5)
- Goal 4, User Story 4, Core Feature 7, Success Metric 4 → `PlanningBinding`, `inspectCarryForwardCandidate` (API Contracts 6, 7)
- Goal 5, User Story 5, Core Feature 8, Success Metric 5 → `historyADRSuccessor` and the naming test (API Contract 8)
- Core Feature 9 → the five user guides, the Roundfix references, `CONTEXT.md`
- Success Metrics 1-5, API Contract 9 and Invariants 1-7 → the named tests and the scope audit

## Absence Evidence

1. The Run Marker for Agents and Verification — command: `grep -rn 'ROUNDFIX_RUN_ID' internal cmd --include='*.go'`; result: only `internal/notify/notify.go:179` and its test at `23ab0904`; no judge, agent or daemon match.
2. `TYPESAFE_API_KEY` withheld — command: `grep -rn 'TYPESAFE_API_KEY' internal/judge/environment.go internal/judge/questions.json`; result: no match at `23ab0904`.
3. `PlanningBinding` and `AmendedInputs` — command: `grep -rn 'AmendedInputs\|amendedInputs\|PlanningBinding' internal --include='*.go'`; result: no match at `23ab0904`.
4. The Windows compile contract — command: `grep -rln 'GOOS=windows' internal --include='*_test.go'`; result: no match at `23ab0904`.
5. ADR reference forms — command: `grep -n 'supersedes' internal/spec/history_entries.go`; result: no match at `23ab0904`, and `historyADRSuccessor` reads only the `ADR-` prefix.
6. Per-sibling release FIFOs — command: `grep -n 'releasePath := ' internal/worktree/worktree_test.go`; result: one path per round at `:471`, shared by every sibling's command.

## ADR Alternatives

- ADR-0283: this Spec's decision; it rejects a runtime setting, an empty `ZDOTDIR` or `HOME`, detection from the working directory and a refusal only when keys are present.
- ADR-0285: this Spec's decision; it rejects accepting the gap, a static selector scan, an `always` contract and a Makefile or workflow change.
- ADR-0284: this Spec's decision; it rejects item-level binding, ignoring planning files and carrying part of a refused set.
- ADR-0273: takes, in a new form, the rejected "Skipping Archive Advice when archive --plan runs inside a Run": the Run is detected from a variable the Daemon sets, not from Roundfix Home or the worktree, and every judge caller honours it; ADR-0283 amends ADR-0273 with a dated note, and the keys' removal stays.
- ADR-0271: no rejected alternative is taken; the bootstrap fix adds no sleep, retry or timing allowance.
- ADR-0266: no rejected alternative is taken; the moved tests keep their Proven Signals.
- ADR-0250: no rejected alternative is taken; the record is still collected without running a test binary, and it now lists the moved tests for two platforms.
- ADR-0253: no rejected alternative is taken; the Full Contract Run discovers the new contract with no hand-kept list.
- ADR-0252: no rejected alternative is taken; the new contract keeps the default relevance.
- ADR-0170: no rejected alternative is taken; a refused Task still refuses the whole remaining set, and completed Tasks are still nothing to carry.
- ADR-0158: no rejected alternative is taken; the accepted outcomes are unchanged.
- ADR-0282: no rejected alternative is taken; the lineage line is still written only where an accepted successor does not name its predecessor.
- ADR-0278: no rejected alternative is taken; a superseded ADR with a resolved successor still leaves the tree.
- ADR-0254: no rejected alternative is taken; a kept ADR is still a Refused Unit that never blocks.
- ADR-0189: no rejected alternative is taken; each skill edit raises its recorded version.
- ADR-0280: no rejected alternative is taken; each guide Task records its own fingerprint.
- ADR-0256: no rejected alternative is taken; each changed surface's skill reference changes with it.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0104: no rejected alternative is taken; the outside-evidence row quotes published sources.
- ADR-0276: no rejected alternative is taken; network-free rows qualify and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0092: no rejected alternative is taken; the backlog and retirement rules are unchanged.
- ADR-0120: no rejected alternative is taken; the History Root keeps its one job.
- ADR-0247: no rejected alternative is taken; an Archive Record and its folder in Git are unchanged.
- ADR-0248: no rejected alternative is taken; the History Full Tag still precedes every removal.
- ADR-0261: no rejected alternative is taken; an Archive Record still names its pull request.
- ADR-0277: no rejected alternative is taken; legacy folder conversion is unchanged.
- ADR-0258: no rejected alternative is taken; no Reduced History Entry is written.
- ADR-0260: no rejected alternative is taken; the update still lists each citation it writes.
- ADR-0164: no rejected alternative is taken; the Run budget still renews at each settlement.
- ADR-0233: no rejected alternative is taken; a version raise is still recorded by its command.
- ADR-0240: no rejected alternative is taken; the QA partial policy is unchanged.
- ADR-0281: no rejected alternative is taken; no test rebuilds an executable in place.
- ADR-0251: no rejected alternative is taken; lenient legacy reading is unchanged.
- ADR-0262: no rejected alternative is taken; cancellation is unchanged.
- ADR-0269: no rejected alternative is taken; archive and cancel refusals are unchanged.
- ADR-0165: no rejected alternative is taken; a blocking review still parks publication.
- ADR-0272: no rejected alternative is taken; a review record still survives a fix at its head.
- ADR-0167: no rejected alternative is taken; the Pull Request row stays blocked before a Pull Request exists.
- ADR-0229: no rejected alternative is taken; the operator archive is unchanged.
- ADR-0169: no rejected alternative is taken; the review still diffs the candidate from its merge base.
- ADR-0279: no rejected alternative is taken; a dismissal still survives a fix that leaves its anchor.
- ADR-0196: no rejected alternative is taken; a pre-PR review finding still parks only after validation.
- ADR-0237: no rejected alternative is taken; the Delivery Retry still carries forward before it retries, under the rule ADR-0284 states.

## Integration Points

None outside the repository. The Windows contract runs the local Go
toolchain; nothing reaches a provider, a runtime's state or a shell startup
file.

## Testing Approach

Every new top-level test calls `t.Parallel()` first and builds temporary
repositories with the fixtures already used beside it.

- `internal/judge/run_marker_test.go`: with `RunID` set and keys present, a
  local `httptest` server reached through the default transport sees no
  request from `Run` or `AdviseArchive`, no Judge Log appears, and the skip
  reason is API Contract 2's; with `RunID` empty the same fixture is
  unchanged. `TestWithoutKeyVariablesRemovesEveryJudgeVariable` also asserts
  `TYPESAFE_API_KEY`.
- `internal/agent/judge_keys_test.go`: `TestNoAgentSessionReceivesAJudgeKeyVariable`
  also sets `TYPESAFE_API_KEY` and asserts that each prepare and prompt
  invocation records `ROUNDFIX_RUN_ID` with its Run ID.
- `internal/daemon/verification_environment_test.go`:
  `TestVerificationReceivesNoJudgeKeyVariable` also checks `TYPESAFE_API_KEY`
  unset and the marker equal to the request's Run ID.
- `internal/cli/judge_run_marker_test.go`: `spec judge`, `archive --plan`
  (Surface Transcript 1 as one test) and `history sanitize` advice with an
  injected environment holding the marker and a counting transport see no
  request and print the reasons.
- `internal/worktree/worktree_test.go`: the round uses per-sibling FIFOs, and
  `TestBootstrapReleaseWriterOutlivingItsSiblingReachesNoOtherSibling` holds
  each release writer open until the next sibling has ended and asserts every
  sibling succeeds.
- The moved tests keep their names; `internal/verifyselect`'s new contract
  runs only under `repocontract`.
- `internal/spec/carry_forward_binding_test.go`: node and row extraction,
  section cuts, plural matching, fenced headings, unresolved units and a
  missing bullet.
- `internal/cli/carryforward_amendment_test.go`: an amendment adding a
  Tooling authority row and a corrective Task node carries the settled Task
  with `amendedInputs`; an amendment of a named section refuses with today's
  text; `TestRetryRefusalNamesTheAmendingCommitAndTheRecovery` keeps passing.
- `internal/spec/history_successor_forms_test.go`: each `superseded_by` form
  resolves; `supersedes` in each form suppresses the lineage line; the fluxus
  shape removes the ADR with no line.

## Build Order

1. The Run Marker, the silent judge and the withheld TypeSafe key, with their
   tests, the `spec`, `archive` and `history` guides, the Roundfix `spec` and
   `archive` references, the glossary, the skill version and the Behavior
   Surface Record (no dependency).
2. The race-free bootstrap test, the Unix-only test files, the coverage record
   and the Windows compile contract (no dependency).
3. Carry-forward over amendments, with its tests, the `reconcile` guide and
   reference, the glossary, the skill version and the record (depends on: 1,
   which raises the same Roundfix version and record and edits the glossary).
4. ADR reference forms, with their tests, the `history` and `baseline` guides,
   the Roundfix `archive` and `baseline` references, the skill version and the
   record (depends on: 3, the same version and record, and on 1, which edits
   the `history` guide and the `archive` reference).
5. QA gate (depends on: 2, 4).

No other Spec is a prerequisite.

## Risks & Considerations

- On a host whose shell startup exports a judge key, an Agent's commands still
  see it; ADR-0283 records that Roundfix guarantees only its own judge.
- Steps 1, 3 and 4 each raise the Roundfix skill version and re-record two
  shared files, so they run in series; step 2 shares no file with them.
- The Windows list of the coverage record shrinks; ADR-0250's comparison
  accepts that only through the sanctioned re-record.
- A test added by steps 3 or 4 that uses a Unix-only call fails the Windows
  contract in the gate; their Tasks say so.

## Glossary

- adds: **Run Marker**
- adds: **Binding Text**
- changes: **Judge Key Variables**
- changes: **Task Carry-Forward**

## Decisions

- The marker reuses the notification command's `ROUNDFIX_RUN_ID`. See ADR-0283.
- The refusal lives in the judge's request path and reads the caller's environment. See ADR-0283.
- Repository-configured commands keep the Run's environment, as ADR-0273 decided.
- The bootstrap fix is test-only: per-sibling release FIFOs.
- The Windows contract has default relevance and vets with both contract tags. See ADR-0285.
- Binding Text is section-level, from the References bullets. See ADR-0284.
- The ADR reference forms are read in `superseded_by` and `supersedes`.

## Amendment of 2026-10-18

`TestEveryTestCompilesForWindows` is also listed in the Makefile's
`REPO_CONTRACT_TESTS`, which the repository gate audit requires of every
repocontract test; the maintainer authorized that one line on 2026-10-18, and
the audit keeps its delivery-base behavior.
