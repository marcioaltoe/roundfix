---
task: task_01
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
status: completed
type: backend
complexity: high
---

# Task 01: No Agent or Verification of a Run can make the judge send, and the TypeSafe key stays out of Runs

## Overview

Roundfix already removes the Judge Key Variables from every Agent Session and
Verification command, yet on 2026-10-16 and 2026-10-17 two QA Agents sent Spec
files to the judge, because their runtime re-read the user's shell startup and
exported the key again (`_techspec.md` → Current behavior). This Task gives
every Agent Session and Verification command of a Run the Run Marker
`ROUNDFIX_RUN_ID`, makes the judge send nothing while it is set, and adds
`TYPESAFE_API_KEY` to the withheld variables. It is verifiable on its own
through judge, agent, daemon and command tests that send nothing to a
provider.

## Requirements

1. MUST answer the operator observations of 2026-10-16 (Spec 0273, QA-001) and
   2026-10-17 (Spec 0275, F2) cited in `_prd.md` → Sources 1 by implementing
   API Contracts 1 and 2 of `_techspec.md` (ADR-0283).
2. MUST define `judge.RunMarkerVariable` as `ROUNDFIX_RUN_ID`, add `RunID` to
   `judge.Request` and `judge.ArchiveAdviceRequest`, and make `judge.Run` and
   `judge.AdviseArchive` check it before selecting a transport: with a
   non-empty `RunID` they read no key, read and write no Judge Log, send no
   request and report the skip reason
   `inside Run <run-id>, named by ROUNDFIX_RUN_ID`. The data source is the
   environment the calling command was given; the destination that must not
   be reached is the judge provider, whatever keys that environment holds.
3. MUST make `roundfix spec judge`, `roundfix archive <slug> --plan` and the
   advice of `roundfix history sanitize` pass the marker read from their
   command environment, so `spec judge` prints
   `Judge: skipped: inside Run <run-id>, named by ROUNDFIX_RUN_ID; <n> judgment(s) not asked`
   and the other two print each candidate with
   `no advice (inside Run <run-id>, named by ROUNDFIX_RUN_ID)`, matching
   Surface Transcript 1 for the archive plan. `--offline` keeps its lines.
4. MUST make every acpx invocation the runner makes for an ensure, prepare or
   prompt request with a non-empty Run ID receive
   `ROUNDFIX_RUN_ID=<Run ID>`, replacing any inherited value, and make
   `ExecVerifier.Verify` run a command whose `VerifyRequest.RunID` is not empty
   with the same variable; the engine MUST pass the Run ID for a Task's
   Verification, the probe on the unchanged tree and the repository
   Verification at settlement. MUST NOT change the environment of Git hooks,
   the QA format command or the Worktree Bootstrap (ADR-0273), or of the
   Daemon and the Queue Owner themselves, whose Light Tier sessions still get
   the shared key under the implementation Stage Key name.
5. MUST make `judge.WithoutKeyVariables` also remove `TYPESAFE_API_KEY`, so no
   Agent Session or Verification command of a Run receives it, and keep every
   other variable.
6. MUST add the tests named in Verification, each calling `t.Parallel()`
   first unless it names the process-wide state it changes: in the new
   `internal/judge/run_marker_test.go`, `TestJudgeSendsNothingInsideARun` and
   `TestArchiveAdviceSendsNothingInsideARun`, each with keys present and a
   local `httptest` server reached through the default transport that counts
   zero requests, no Judge Log file, and the exact skip reason, and the same
   fixture with an empty `RunID` unchanged; in the new
   `internal/cli/judge_run_marker_test.go`, `TestSpecJudgeSkipsInsideARun`,
   `TestArchivePlanInsideARunMatchesItsTranscript` (Surface Transcript 1 as
   one test: every stdout and stderr line in order and the exit code) and
   `TestHistorySanitizeAdviceSkipsInsideARun`, each with an injected
   environment holding the marker and a judge key and a counting transport
   that sees no request. MUST extend
   `TestWithoutKeyVariablesRemovesEveryJudgeVariable`,
   `TestNoAgentSessionReceivesAJudgeKeyVariable` (each prepare and prompt
   invocation records `ROUNDFIX_RUN_ID` with its Run ID, and none records
   `TYPESAFE_API_KEY`) and `TestVerificationReceivesNoJudgeKeyVariable` (the
   command sees `TYPESAFE_API_KEY` unset and the marker equal to the request's
   Run ID), keeping their names. Every new test compiles for Windows: no
   Unix-only system call outside a `unix`-constrained file.
7. MUST document the Run Marker, the skip reason and the withheld TypeSafe key
   in `docs/user-guide/commands/spec.md`, `docs/user-guide/commands/archive.md`,
   `docs/user-guide/commands/history.md`,
   `.agents/skills/roundfix/references/spec.md` and
   `.agents/skills/roundfix/references/archive.md`, each naming
   `ROUNDFIX_RUN_ID`, and search the guides and skill copies for text that says
   the judge keys are the only protection inside a Run. No text inside any
   `### QA settlement` section may change.
8. MUST add **Run Marker** to `CONTEXT.md` and revise **Judge Key Variables**
   through domain-modeling, with ADR-0283's link.
9. MUST run `make skills-sync`, raise and record the Roundfix skill version
   with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
   because this Task edits user guide pages and command sources. No record or
   mirror is hand-edited.

## Subtasks

- [ ] Add the marker refusal to the judge and the three callers.
- [ ] Pass the marker to acpx invocations and Verification commands of a Run.
- [ ] Withhold `TYPESAFE_API_KEY` from Runs.
- [ ] Add and extend the tests, the transcript test included.
- [ ] Document the marker in the guides, the references and the glossary.
- [ ] Sync the mirrors, record the version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] With `ROUNDFIX_RUN_ID` set and a judge key present, neither judge entry point sends a request or writes a Judge Log, and each reports the Run.
- [ ] Every prepare and prompt acpx invocation and every Verification command of a Run carries the marker with its Run ID and no Judge Key Variable or `TYPESAFE_API_KEY`.
- [ ] The archive plan inside a Run matches Surface Transcript 1.
- [ ] With the marker unset, every judge behavior is unchanged.

## Consumers

- The judge request contract and its skip reasons — consumers: `./internal/judge`, `./internal/cli`
- The Agent Session and Verification environment — consumers: `./internal/agent`, `./internal/daemon`

## Context

- instruction: `docs/adr/0283-a-run-marker-keeps-the-judge-silent-inside-a-run.md`
- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- instruction: `internal/agent/light_session.go`
- instruction: `internal/notify/notify.go`
- interface: `internal/judge/environment.go`
- interface: `internal/judge/environment_test.go`
- interface: `internal/judge/judge.go`
- interface: `internal/judge/archive_advice.go`
- creates: `internal/judge/run_marker_test.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/codex_spawn.go`
- interface: `internal/agent/judge_keys_test.go`
- interface: `internal/daemon/daemon.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/verification_probe.go`
- interface: `internal/daemon/verification_environment_test.go`
- interface: `internal/cli/spec_judge.go`
- interface: `internal/cli/archive.go`
- interface: `internal/cli/history.go`
- creates: `internal/cli/judge_run_marker_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/history.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestJudgeSendsNothingInsideARun|TestArchiveAdviceSendsNothingInsideARun|TestWithoutKeyVariablesRemovesEveryJudgeVariable)$' ./internal/judge 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestJudgeSendsNothingInsideARun TestArchiveAdviceSendsNothingInsideARun TestWithoutKeyVariablesRemovesEveryJudgeVariable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^TestNoAgentSessionReceivesAJudgeKeyVariable$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestNoAgentSessionReceivesAJudgeKeyVariable " || { printf 'missing pass: TestNoAgentSessionReceivesAJudgeKeyVariable\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^TestVerificationReceivesNoJudgeKeyVariable$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestVerificationReceivesNoJudgeKeyVariable " || { printf 'missing pass: TestVerificationReceivesNoJudgeKeyVariable\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestSpecJudgeSkipsInsideARun|TestArchivePlanInsideARunMatchesItsTranscript|TestHistorySanitizeAdviceSkipsInsideARun)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSpecJudgeSkipsInsideARun TestArchivePlanInsideARunMatchesItsTranscript TestHistorySanitizeAdviceSkipsInsideARun; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/spec.md docs/user-guide/commands/archive.md docs/user-guide/commands/history.md .agents/skills/roundfix/references/spec.md .agents/skills/roundfix/references/archive.md; do grep -qF -- 'ROUNDFIX_RUN_ID' "$f" || { printf 'missing ROUNDFIX_RUN_ID in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Run Marker**:' || { printf 'missing glossary term Run Marker\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Judge Key Variables**:' || { printf 'missing glossary term Judge Key Variables\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 ./internal/judge || exit 1; go test -count=1 -run 'Judge|Archive|History' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the two judge tests do not exist, so the command fails at its first pass check; afterwards the judge, agent, daemon and command tests, the documentation, the glossary, the skill version and the Behavior Surface Record checks pass.

## References

- `_prd.md` → Sources; Goals; User Story 1; Core Feature 1; Core Feature 2; Core Feature 3; Core Feature 9; Success Metric 1; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 9; Surface Transcript 1; Invariants 1, 2 and 3; Current behavior; Testing Approach; Build Order 1
- ADR-0283; ADR-0273; ADR-0280; ADR-0189; ADR-0256; ADR-0259

## Result

### Implementation

- Added `judge.RunMarkerVariable` (`ROUNDFIX_RUN_ID`) and explicit `RunID`
  fields on both judge requests. A non-empty marker returns the exact Run
  skip reason before transport selection, credential lookup, the Judge Log
  clock or log access. Spec judgments retain their planned skipped rows;
  archive candidates each receive the Run reason.
- The three command callers read the marker from their injected command
  environment and collect no judge credentials while it is non-empty.
  Explicit `--offline` keeps its existing transcript.
- Session ensure/prepare and prompt invocations append a request-specific
  marker override without mutating the runner's inherited environment or
  cached runtime overrides. Existing acpx environment replacement removes an
  inherited marker. `WithoutKeyVariables` now also removes
  `TYPESAFE_API_KEY`, preserving unrelated entries and implementation Stage
  Keys.
- `ExecVerifier` replaces the inherited marker with `VerifyRequest.RunID`.
  The engine forwards the Run ID for Task and repository Verification. The
  pre-work probe uses an internal Run-aware entry point; standalone
  `ProbeCommands` keeps its interface and empty Run ID. The necessary caller
  change in `internal/daemon/task_engine.go` belongs to this Task's unchanged
  tree probe requirement.
- Updated the three command guides, canonical spec/archive skill references,
  and the Run Marker and Judge Key Variables glossary terms with ADR-0283.
  Searched the guides and skill copies for descriptions of judge protection;
  remaining key-withholding statements describe the still-active removal
  rule. No QA settlement section changed.
- Regenerated the mirrors and records using their owning commands; the
  Roundfix skill version is `0.1.90`. No record or mirror was hand-edited.

### Acceptance evidence

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| A marker with keys present prevents both judge entry points from sending or writing a log, and names the Run. | `TestJudgeSendsNothingInsideARun` and `TestArchiveAdviceSendsNothingInsideARun` use the default HTTP transport with a local counting server, assert zero requests, no Judge Log directory, the exact reason and no calls to the log clock. All three new command tests inject a marker and key and assert zero transport calls. |
| Each prepare/prompt and Verification command carries its Run ID and no withheld credential variable. | Extended `TestNoAgentSessionReceivesAJudgeKeyVariable` records the marker at acpx boundaries for standard and light sessions, checks replacement of an inherited marker, rejects `TYPESAFE_API_KEY`, retains the implementation Stage Key and verifies the runner's environment stays intact. Extended `TestVerificationReceivesNoJudgeKeyVariable` executes a child that checks all withheld variables, its request marker and a control variable, then checks that the parent's variables remain intact. `TestEnginePassesRunIDToEveryVerificationCommand` exercises and observes the engine's pre-work probe, Task Verification and repository settlement Verification. |
| The archive plan inside a Run matches Surface Transcript 1. | `TestArchivePlanInsideARunMatchesItsTranscript` asserts every stdout/stderr line and the exit code for the keyless transcript and the same fixture with a key present; it also checks Spec bytes and the existing offline transcript. |
| An empty marker preserves existing judge behavior. | Empty-marker controls in both judge tests make real local HTTP requests and write Judge Logs. Focused existing judge/CLI controls cover advisory results, the command's injected key, monthly ceilings, archive advice and history advice; the existing offline transcript also passes. |

### Focused checks and generation

- Red starting point: `go test ./internal/judge -run 'InsideARun' -count=1`
  exited 1 because the two request types lacked `RunID`, before the production
  change.
- `go test -race ./internal/judge ./internal/agent ./internal/daemon ./internal/cli -run 'InsideARun|WithoutKeyVariables|NoAgentSessionReceives|VerificationReceivesNoJudge|EnginePassesRunID' -count=1`
  exited 0 in all four packages after generation finished.
- An earlier invocation of that race check had passing assertions but exited
  1 in `internal/agent`: skill version regeneration overlapped its suiteguard
  fingerprint window. After the generator finished, the same check passed
  without concurrent repository writes or any guard changes.
- `go test ./internal/judge ./internal/cli -run '^(TestRunRaisesAtTheMeasuredThresholds|TestRunStopsAtTheMonthlyCeiling|TestRunSendsNothingWithoutAKey|TestArchiveAdviceJudgesOnlyCandidateFiles|TestSpecJudgeReportsAdvisoryJudgments|TestSpecJudgeReadsTheKeyFromItsOwnEnvironment|TestSpecJudgeSkipsAtTheMonthlyCeiling|TestHistorySanitizeAdviseWithAFakeJudge|TestArchivePlanOfflineTranscript)$' -count=1`
  exited 0 in both packages.
- `make skills-sync` exited 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and recorded `0.1.90` in both skill copies and the version record.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited 0 and re-recorded the three changed guide fingerprints.
- `make baseline-digests` exited 0 with `changed:false`; no Baseline artifact
  remains changed.
- Read-only byte comparisons confirmed that the generated Roundfix mirror
  matches its canonical folder and every existing `### QA settlement`
  section in the changed canonical skill files matches `HEAD`. The five
  requested documentation files contain both `ROUNDFIX_RUN_ID` and
  `TYPESAFE_API_KEY`. `git diff --check` exited 0.
- A focused Windows compile, `GOOS=windows GOARCH=amd64 go test -c -o <temporary-directory> ./internal/judge ./internal/agent ./internal/daemon ./internal/cli`,
  produced `judge.test.exe` and `daemon.test.exe`, but exited 1 for the
  pre-existing Unix-only tests in Agent and CLI: `syscall.Kill` in
  `internal/agent/acpx_integration_test.go`, and `syscall.Stat_t`,
  `syscall.Getpgrp`, signal helpers, `syscall.Wait4` and `Setpgid` in the CLI
  files listed in the TechSpec's Current behavior. Those repairs belong to
  task_02. New tests use no Unix-only system call or new platform constraint.

### Handoff boundary

The pre-existing status change to `in_progress` is preserved. The Task's
status, authored Verification, other Tasks and `_tasks.md` were not edited.
No authored Verification command, selected repository Verification or full
suite was run; the three recording commands above are required generation
steps. No commit, push or pull request was made. All provider behavior used
injected transports or the local test server; no request reached a provider.
Daemon Verification and settlement remain pending.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/daemon/task_engine.go`
