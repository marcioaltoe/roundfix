---
spec: 0267-no-secret-and-no-spec-content-leaves-the-machine
prd: _prd.md
created: 2026-10-13
---

# No secret and no Spec content leaves the machine — Technical Spec

## Executive Summary

A new leaf package, `internal/secretscan`, holds the one reviewed table of
Secret Patterns and the pure functions that turn file content or a pushed
history into Secret Matches of path, line and kind. Three call sites use it:
the QA Report commit, the Delivery Queue push and the Final Push. A small
helper in `internal/judge` removes the Judge Key Variables from an
environment, and the two builders of a Run's command environments, the acpx
environment of every Agent Session and `ExecVerifier`, call it. The trade-off
accepted: refusing is coarser than redacting. A pattern that is too broad
stops an honest QA Report or push until the table is narrowed; in exchange
nothing is rewritten after QA observed it, and the refusal points at the
command that printed the value (ADR-0273).

## Project Constraints

- Identifier strategy: not applicable. No identifier, schema or stored format
  gains a key; `secret-found` is a value of the existing Delivery Queue
  blocker string. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable. The change decides which processes
  receive the judge's credentials and what a push may carry; it adds no
  endpoint, credential or transport, and credentials stay in the environment
  (`docs/agents/agent-instructions.md`). Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0273 is this design. ADR-0015 says
  the Daemon commits the QA Report on any verdict (ADR-0015: "commits the QA
  Report in its own commit either way"); ADR-0273 refines it, and task_01
  keeps every clean report's commit unchanged. ADR-0201 reads each
  judge key from the environment only (ADR-0201: "Each is read from its
  variable in the environment and nowhere else"), so
  removing the variables from a process removes the judge from it. ADR-0239:
  "One list in Roundfix names each stage's variables in that order, and every
  reader takes it from there"; the helper reads `Questions.KeyVariables`, never
  a copy. ADR-0238: the light session references its key variable by name;
  task_02 moves the reference to the implementation Stage Key. ADR-0138 owns
  the Final Push, which task_03 guards.
  ADR-0256: task_04 changes the `deliver`, `implement` and configuration
  user-guide surfaces with their covering skill files. ADR-0189 and ADR-0233
  govern the owned skills' version raise. ADR-0244 governs the Glossary
  Declaration. ADR-0259 and ADR-0271 govern every new test. ADR-0265 governs
  the three proofs below. ADR-0200, ADR-0208, ADR-0209 and ADR-0231 do not
  apply: the judge's advisory role, source grouping, grouping suggestions and
  ceiling are unchanged; only which processes hold the judge's keys changes.
  ADR-0235 does not apply: no model selection or subscription rule changes.
  ADR-0250 does not apply: no Baseline module version or coverage record
  changes. ADR-0257 does not apply: the Daemon stays the only full gate inside
  a Run, and only the environment of its Verification commands changes.
  ADR-0258: the QA Task's rows declare their inputs, the scope row a
  `commit_range`. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. No lint, formatter, test-runner, Makefile or
  CI configuration changes. The Governed Paths, measured with
  `speccheck.GovernedPath` through a `go test -overlay` probe that wrote
  nothing, are owned skill files granted by the cycle AG approval of
  2026-10-13. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. Spec-contained authorization record:
  `docs/specs/0267-no-secret-and-no-spec-content-leaves-the-machine/_authorization.md`.
  Bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/implement.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `skills/qa-gate/SKILL.md`, `skills/roundfix/SKILL.md`.

## System Architecture

New:

- `internal/secretscan` (task_01) — `patterns.go` holds the Secret Pattern
  table, the one reviewed place; `scan.go` holds the scan of file content, the
  git arguments for a pushed history and the parser of that history, and the
  `MatchesError` both push sites return. It imports only the standard library.
- `internal/judge/environment.go` (task_02) — `WithoutKeyVariables`.

Changed:

- `internal/daemon/task_engine.go` (task_01) — the QA step computes the QA
  Report commit's path set once, scans it before `settleTask`, and on a
  Secret Match settles `failed`, publishes one `daemon.qa` Run Event with
  phase `secret_scan`, and skips `commitQAReport`.
- `internal/agent/codex_spawn.go` and `internal/agent/light_session.go`
  (task_02) — `acpxCommandEnv` drops the Judge Key Variables from its base;
  `lightSessionEnvironment` references
  `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY` and, when the selected variable is
  another one, adds that name with the selected variable's value.
- `internal/daemon/daemon.go` (task_02, then task_03) — `ExecVerifier` sets
  `cmd.Env` from `judge.WithoutKeyVariables(os.Environ())`; `GitPusher.Push`
  scans the pushed history before `git push`.
- `internal/delivery/github.go`, `engine.go`, `park_class.go` (task_03) —
  `GitHubCLI.PushBranch` scans before `git push` through its `CommandRunner`;
  `Engine.push` parks `secret-found` on a `*secretscan.MatchesError`;
  `ClassifyPark` gives it class `finding` and its next action.
- Skill, prompt, guide and glossary text (task_04).

Characterization at `9963ffd3`, through a `go test -overlay` probe in
`internal/agent` and `internal/daemon` that wrote nothing to the repository:
`acpxCommandEnv` over a base holding the three Judge Key Variables returned
all three; an `ExecVerifier.Verify` command testing that
`ROUNDFIX_TYPESAFE_API_KEY` is unset failed when the test process set it.
`commitQAReport`, `GitHubCLI.PushBranch` and `GitPusher.Push` read no file
content. Existing tests that pin today's behavior and change:
`TestOpenModelImplementationFallsBackToTheSharedKey` and
`TestLightSessionReadsTheImplementKeyByName` in
`internal/agent/light_session_test.go` expect the shared variable in the
light session's reference and environment (task_02 moves them to API
Contract 5); `TestPushPublishesTheRecordedHeadNotTheCheckout` and
`TestPullRequestBoundaryRejectsAPushWithoutAnObservedRemoteHead` in
`internal/delivery/github_test.go` script only `git push` and `ls-remote`
(task_03 adds the two scan steps to their scripts); `TestParkClassesAndNextActions`
in `internal/delivery/park_class_test.go` gains the new row.

## Implementation Design

### Interfaces

```go
package secretscan

// Match is one Secret Match. It never holds the matched value.
type Match struct{ Path string; Line int; Kind, Commit string }
func (m Match) String() string // "<path>:<line>: <kind>" plus " (commit <12 hex>)" when Commit is set
type MatchesError struct{ Matches []Match } // Error lists every Match, one per line
func Kinds() []string                          // table order
func ScanContent(path string, content []byte) []Match
func PushedHistoryArgs(head, remote string) []string // git arguments, see API Contract 3
func ScanPushedHistory(log []byte) []Match
```

```go
package judge

// WithoutKeyVariables returns a copy of environ without any entry whose name
// is one of Questions.KeyVariables of the embedded questions.
func WithoutKeyVariables(environ []string) []string
```

Invariants:

1. No Match, error, Run Event, progress line, Task file reason, park blocker
   or log contains any byte sequence a Secret Pattern matched.
2. Every Secret Pattern sample in a test is built at runtime by
   concatenation; no repository file holds a string that matches a pattern.
3. A clean scan changes nothing: the QA Report commit, its events and both
   pushes behave as at `9963ffd3`.
4. A line yields at most one Match, of the first kind in table order that
   matches it.
5. A Light Tier session's environment holds no Judge Key Variable and holds
   its OpenRouter key under `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY`.

### Data Models

The Secret Pattern table is a Go slice of `{kind, expression}` in
`internal/secretscan/patterns.go`, in this order, every expression compiled
once with `regexp.MustCompile`:

| kind | shape |
| --- | --- |
| `private-key-block` | a `-----BEGIN ` header ending in `PRIVATE KEY-----` |
| `github-token` | `ghp_`, `gho_`, `ghu_`, `ghs_` or `ghr_` and 36 letters or digits; or `github_pat_` and 82 letters, digits or underscores |
| `openrouter-key` | `sk-or-v1-` and 64 lowercase hex digits |
| `anthropic-key` | `sk-ant-` and at least 32 letters, digits, `-` or `_` |
| `openai-key` | `sk-`, an optional `proj-`, and at least 32 letters, digits, `-` or `_` |
| `aws-access-key-id` | `AKIA` or `ASIA` and 16 uppercase letters or digits |
| `slack-token` | `xox` and one of `a b p r s`, `-`, and at least 10 letters, digits or `-` |
| `google-api-key` | `AIza` and 35 letters, digits, `-` or `_` |
| `stripe-secret-key` | `sk_live_` or `rk_live_` and at least 24 letters or digits |
| `credential-assignment` | an uppercase name ending in `API_KEY`, `TOKEN`, `SECRET` or `PASSWORD`, then `=`, `:=`, `?=` or `:` with optional spaces and an optional quote, then a value of at least 20 characters from letters, digits and `_ + / . = -` that holds at least one letter and one digit and does not start with `=` |

Token shapes are bounded by non-word characters on both sides. The last row
catches what `env`, `printenv` and `make -p` print for any credential
variable, whatever the value's shape. Its value rule was measured on
2026-10-13 at `9963ffd3` over every file `git ls-files` lists: with a
16-character minimum and no digit requirement it matched 12 test fixture
lines such as a sentinel word assigned to a key variable and one placeholder
in an upstream skill; with the rule above, the whole table matched 0 lines,
and a `make -p` style line holding an OpenRouter-shaped value still matched.
A placeholder such as `$VAR`, `<value>`, `${{ secrets.X }}` or `[redacted]`
never matches.

The Delivery Queue stores the refusal in the existing blocker string:
`secret-found: <match>, <match>, …`, every Match's `String()` without the
leading `secret-found: `.

### API Contracts

1. API Contract 1: `ScanContent` treats content holding a NUL byte in its
   first 8000 bytes as binary and returns no Match. Otherwise it reports, per
   line numbered from 1, at most one Match whose `Path` is the given path.
2. API Contract 2: the QA step's Secret Scan reads, from the Run Worktree,
   exactly the paths the QA Report commit would stage (the step's snapshot
   difference, the report, the imported prior-pass files and the QA Task file,
   after the Project Config and stageability filters). On a Match it settles
   the QA Task `failed` with the reason `QA evidence carries a secret: ` and
   the Matches joined by `; `, publishes one `daemon.qa` event whose message
   is `QA Report refused: N Secret Match(es)` and whose data carries
   `phase: secret_scan` and `matches` as a list of `{path, line, kind}`, and
   creates no QA Report commit. The report, its evidence and the Task file
   stay in the Run Worktree.
3. API Contract 3: `PushedHistoryArgs(head, remote)` returns
   `log -p -U0 --reverse --no-color --no-ext-diff --no-renames --format=commit %H <head> --not --remotes=<remote>`;
   the push scans every commit not reachable from the remote's tracking refs,
   merge commits contributing no lines. `ScanPushedHistory` reads that output
   and reports, for each added line of each commit, at most one Match with the
   commit, the new path and the line number in that commit's version of the
   file; a deleted file and a binary file contribute nothing.
4. API Contract 4: `GitHubCLI.PushBranch` and `GitPusher.Push` run the
   history read through their own runner before `git push`; on one or more
   Matches they return a `*secretscan.MatchesError` (wrapped with
   `push branch: refuse push:` and `run Final Push: refuse push:`) and run no
   `git push`. A failed history read refuses the push too. `Engine.push`
   parks the item with blocker `secret-found: …`, whose Park Class is
   `finding` and whose next action is `remove the secret from every commit
   the blocker names on the item branch in <worktree>, then run roundfix
   deliver retry <slug>`.
5. API Contract 5: `acpxCommandEnv` drops every Judge Key Variable from its
   base before it applies overrides. `lightSessionEnvironment` writes the
   OpenCode option `apiKey` as `{env:ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY}`;
   when `RuntimeSpec.OpenRouterKeyVariable` names another variable, it adds
   `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY=<that variable's value>` to the
   overrides. It keeps removing the generic `OPENROUTER_API_KEY`. The Light
   Spend Log still records the selected variable's name.
6. API Contract 6: `ExecVerifier.Verify` runs its command with
   `judge.WithoutKeyVariables(os.Environ())`. Every other variable, including
   `GOCACHE`, `HOME` and `PATH`, passes unchanged. This covers Task and QA
   Verification, the repository Verification, the Delivery Queue's gate and
   `spec check --run-verification`, which all run through it.

### Surface Transcripts

None. No command's help, usage, exit code or output format changes: the
refusals reuse the existing QA settlement reason, Run Event and `Park:` line
shapes, and their text is API Contracts 2 and 4.

## Coverage Map

Deviations from the PRD: None.

Cross-read from `_prd.md`:

- Goal 1, User Story 5, Core Feature 1 → the `qa-gate` skill heading and the
  `qaGateContract` line (task_04).
- Goal 2, User Stories 1-2, Core Features 2-3, Success Metric 2 →
  `internal/secretscan` (`ScanContent`, the table) and the QA step's scan;
  API Contracts 1-2; Invariants 1-4.
- Goal 3, User Story 3, Core Feature 4, Success Metric 3 →
  `PushedHistoryArgs`, `ScanPushedHistory`, `GitHubCLI.PushBranch`,
  `GitPusher.Push`, `Engine.push`, `ClassifyPark`; API Contracts 3-4.
- Goal 4, User Story 4, Core Feature 5, Success Metric 1 →
  `judge.WithoutKeyVariables`, `acpxCommandEnv`, `lightSessionEnvironment`,
  `ExecVerifier`; API Contracts 5-6; Invariant 5.
- Success Metric 4 → the table's precision, measured by the QA gate over
  `git ls-files`.
- Core Feature 6 → the `roundfix` skill references, the three user guides and
  `CONTEXT.md` (task_04).

Every PRD goal, story, feature and metric has a row; no component is
unmapped.

## Absence Evidence

1. `internal/secretscan` — command: `ls internal/secretscan`; result: no such directory at `9963ffd3`.
2. A secret check before a commit or push — command: `grep -rn 'secretscan\|SecretMatch\|secret-found\|PRIVATE KEY' internal/daemon internal/delivery internal/cli`; result: no match at `9963ffd3`; `commitQAReport`, `PushBranch` and `GitPusher.Push` read no file content.
3. `judge.WithoutKeyVariables` and any removal of the judge's variables — command: `grep -rn 'WithoutKeyVariables\|KeyVariables()' internal/agent internal/daemon`; result: no match at `9963ffd3`; the overlay probe shows all three variables reach `acpxCommandEnv`'s result and `ExecVerifier`'s command.
4. The new tests — command: `grep -rn 'TestSecretPatterns\|TestQAReportWithASecretMatch\|TestNoAgentSessionReceivesAJudgeKeyVariable\|TestVerificationReceivesNoJudgeKeyVariable\|TestPushRefusesAHistoryThatAddsASecret\|TestFinalPushRefusesAHistoryThatAddsASecret' internal`; result: no match at `9963ffd3`.

## ADR Alternatives

- ADR-0273: the design is this ADR; its rejected options (Run detection in `archive --plan`, a Spec authorization for the keys, QA-only withholding, evidence-only scanning, redaction, an external scanner) are all left out.
- ADR-0015: no rejected alternative is taken; a clean QA Report keeps its own commit, and only a report with a Secret Match is not committed, as ADR-0273 records.
- ADR-0201: no rejected alternative is taken; the judge still reads its keys only from its environment, and its recipients and Judge Log are unchanged.
- ADR-0239: its rejected single-key and no-fallback options stay rejected; the helper reads the existing stage list, and the light session still selects the implementation key first and the shared key second.
- ADR-0238: no rejected alternative is taken; the light session still receives a variable reference, never a value in its configuration.
- ADR-0138: no rejected alternative is taken; the Final Push still pushes only at Clean, after the scan.
- ADR-0256: its rejected model-judging alternative stays rejected; task_04 is the skills Task for the changed user-guide surfaces.
- ADR-0189: no rejected alternative is taken; each owned skill's version rises in both `SKILL.md` fields.
- ADR-0233: no rejected alternative is taken; versions are recorded with the record command, never by hand.
- ADR-0244: no rejected alternative is taken; the four terms are declared and task_04 writes them.
- ADR-0259: no rejected alternative is taken; each new test calls `t.Parallel()` first or names the process-wide state it changes.
- ADR-0271: no rejected alternative is taken; no new test bounds by real time what it does not assert.
- ADR-0265: no rejected alternative is taken; this TechSpec carries the three proofs.
- ADR-0200: does not apply; the judge stays advisory and never gates, and nothing here changes how it answers.
- ADR-0208: does not apply; how sources share a Spec is unchanged.
- ADR-0209: does not apply; grouping suggestions are unchanged.
- ADR-0231: does not apply; the Jev Ceiling is unchanged.
- ADR-0235: does not apply; no model selection or subscription rule changes.
- ADR-0250: does not apply; no Baseline module version or coverage record changes, so none of its alternatives is in play.
- ADR-0257: does not apply; the Daemon stays the only full gate inside a Run, and only its Verification environment changes.
- ADR-0258: no rejected alternative is taken; the QA rows declare their inputs.

## Integration Points

- Git, through each push site's existing runner. The push scan needs the
  remote's tracking refs; when they are stale it scans more commits, never
  fewer.
- No network. The judge's providers lose access from inside a Run; nothing
  else external is touched.

## Testing Approach

- `internal/secretscan/scan_test.go` (task_01): `TestSecretPatternsMatchTheirKindAndNeverEchoTheValue`
  builds one sample per kind at runtime and requires its kind, and that
  `Match.String()` and `MatchesError.Error()` hold no sample byte run of 8 or
  more; `TestScanContentReportsPathLineAndKind`; `TestScanSkipsBinaryAndPlaceholders`
  (`$VAR`, `<value>`, `[redacted]`, a short value); `TestScanPushedHistoryReportsCommitPathAndLine`
  over a real temporary repository whose second commit removes the sample the
  first added. All parallel.
- `internal/daemon/qa_secret_scan_test.go` (task_01): `TestQAReportWithASecretMatchSettlesFailedWithoutACommit`
  and `TestQAReportWithoutASecretMatchKeepsItsCommit`, through the task
  engine's existing QA fixtures with a fake Agent that writes the report and
  one evidence file.
- `internal/judge/environment_test.go`, `internal/agent/judge_keys_test.go`
  and `internal/daemon/verification_environment_test.go` (task_02):
  `TestWithoutKeyVariablesRemovesEveryJudgeVariable`,
  `TestNoAgentSessionReceivesAJudgeKeyVariable` (every acpx command a
  standard and a light session run, through the existing fake acpx, records
  the names it received; a control variable must be present), and
  `TestVerificationReceivesNoJudgeKeyVariable` (sequential: it sets the
  process environment).
- `internal/delivery/push_secret_scan_test.go` and
  `internal/daemon/push_secret_scan_test.go` (task_03):
  `TestPushRefusesAHistoryThatAddsASecret`, `TestEnginePushParksSecretFound`
  and `TestFinalPushRefusesAHistoryThatAddsASecret`, over a temporary
  repository with a bare remote; the remote branch must be absent after a
  refusal and present after a clean push.

## Build Order

1. The Secret Scan package and the QA Report refusal (task_01).
2. The withheld Judge Key Variables (task_02). Independent of step 1.
3. The push scans and the `secret-found` park (task_03) (depends on: 1, 2).
   It depends on step 2 because both edit `internal/daemon/daemon.go`.
4. The skills, the QA prompt, the user guides, the Behavior Surface Record
   and the glossary (task_04) (depends on: 1, 2, 3).
5. Final QA gate (task_05) (depends on: 4).

Spec 0266 merges before this Spec and changes no file this Build Order
declares; the plan is rebased on it before its Run.

## Risks & Considerations

- A too-broad pattern refuses honest text. Success Metric 4 measures the
  repository's own tree, and the table is the one place to narrow.
- An operator's custom OpenCode configuration that reads the shared Roundfix
  variable outside the Light Tier loses it inside a Run; Roundfix's own
  contract never reads it there (ADR-0239).
- A QA Agent that reloads its shell profile (`zsh -ic`) could read keys the
  profile exports. The QA contract forbids it, and no Roundfix instruction
  asks for it inside a Run.
- This Spec's own QA must never write a sample to its evidence, or its own
  scan refuses its report; task_05 requires recording only the Matches.

## Glossary

- adds: **Secret Pattern**
- adds: **Secret Match**
- adds: **Secret Scan**
- adds: **Judge Key Variables**

## Decisions

- One leaf package for the table and the parsing, called from each site's own
  runner, rather than a scanner that runs git itself: the delivery tests
  script every command, and the daemon owns its git invocations.
- Scan before the QA Task settles, so a refused report never reads
  `completed` in the Run Worktree. See ADR-0273.
- Withhold at the two environment builders rather than mutate the Run
  owner's process environment, which in-process Runs and parallel tests share.
  See ADR-0273.
- Research: the Secondbrain and Exa consultations are recorded in
  `_prd.md` → Research record.
