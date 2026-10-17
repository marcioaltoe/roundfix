---
task: task_01
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
status: pending
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
