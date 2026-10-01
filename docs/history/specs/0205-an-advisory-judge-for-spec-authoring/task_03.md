---
task: task_03
spec: 0205-an-advisory-judge-for-spec-authoring
status: completed
type: backend
complexity: medium
---

# Task 03: `roundfix spec judge` prints what the judge raises

## Overview

task_02 can ask a Spec's judgments and nothing reaches it from the command
line. This Task adds `roundfix spec judge <slug> [--stage <prd|techspec>]
[--format <text|json>]`, which resolves the active Spec through the Spec
Root, reads the Jev keys from the command's own environment, runs the judge, and
prints the raised judgments and the summary line exactly as the Surface
Transcripts state. It then documents the command in the command reference and
the Roundfix Skill. It is verifiable on its own: with a fake transport and a
temporary Roundfix Home, the command reproduces Surface Transcripts 1 to 6.

## Requirements

1. MUST add `judge` to `runSpecCommand` and implement API Contract 1: the flags `--stage` (`prd` or `techspec`, default both) and `--format` (`text` or `json`, default `text`), the active Spec resolved through the configured Spec Root, the exit codes that contract names, and the stderr form of Surface Transcript 6.
2. MUST read `ROUNDFIX_OPENROUTER_API_KEY` and `ROUNDFIX_TYPESAFE_API_KEY` only from the command environment's `environ`, never from the process, MUST NOT read `OPENROUTER_API_KEY`, MUST pass them to `judge.Run` as `Request.Keys`, and MUST pass Roundfix Home, a transport and a clock to `judge.Run` through two new fields of `commandDependencies`, `judgeTransport` and `judgeNow`, which default to `http.DefaultTransport` and `time.Now`.
3. MUST print the text form of `_techspec.md` → Surface Transcripts: one `advisory` line per raised judgment, one `skipped` line per skipped judgment or artifact, nothing for a clear judgment, and the summary line with its `stopped:` suffix or its run-level `skipped:` form, with costs to four decimals and confidence and probability to two. Surface Transcripts 1 to 6 MUST be reproduced byte for byte.
4. MUST print the JSON document of API Contract 2 with `--format json`, including `transport`, listing clear judgments too.
5. MUST add the help of API Contract 5 to `specUsage`, the top-level usage in `internal/cli/cli.go` and a new `specJudgeUsage`, keeping every string `TestRunSpecAuditHelpAppearsInUsageAndCommandList` and `TestRunCommandHelp` require.
6. MUST put the tests in `internal/cli/spec_judge_test.go`, each with a temporary repository, a temporary Roundfix Home, a fake `judgeTransport`, a fixed `judgeNow` and the key each Surface Transcript names set in the command's own environment (`ROUNDFIX_OPENROUTER_API_KEY` for Transcripts 1 and 4, `ROUNDFIX_TYPESAFE_API_KEY` for Transcripts 3 and 5, and only the generic `OPENROUTER_API_KEY` for Transcript 2), so no test reaches the network whatever the developer's shell exports.
7. MUST describe the command under a new section `### spec judge` in `docs/user-guide/commands/spec.md`: its synopsis, both stages, the transport order (`ROUNDFIX_OPENROUTER_API_KEY` for OpenRouter first, `ROUNDFIX_TYPESAFE_API_KEY` for TypeSafe directly second), that the generic `OPENROUTER_API_KEY` is not read so Roundfix's Jev cost stays on its own key, the Judge Log path, the monthly ceiling, the skip and stop reasons, and that it exits `0` whenever it ran.
8. MUST describe the command under a new heading `### Advisory judge` in `.agents/skills/roundfix/references/spec.md`: the synopsis, that it never gates and never changes another command's exit code, that each `advisory` line is answered by correcting the artifact or by stating why the text stands, and that a `skipped` result is neither a failure nor a clean result. It MUST add no text inside the `### QA settlement` section of any skill.
9. MUST raise the Roundfix Skill's version in both front-matter fields of `.agents/skills/roundfix/SKILL.md` by one patch step from the value on this Task's starting commit, following the version rule in force on that commit, then run `make skills-sync`, record the version with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, run `make baseline-digests`, and name in the Result every file those commands rewrote.
10. MUST prove each new gate can fail. The Result MUST record one sabotage of the exit code (for example exiting `1` on a stopped run) and one of the key source (for example reading the process environment, or reading `OPENROUTER_API_KEY`), each with the test that failed, and that the code was restored.

Spec 0194 created `.agents/skills/roundfix/references/spec.md` and its mirror, but no user-guide command guide for `spec` (`_prd.md` → Prerequisites). This Task creates `docs/user-guide/commands/spec.md` as the `spec` command guide and adds its row ``| `spec` | [command guide](commands/spec.md) |`` to the command index in `docs/user-guide/commands.md`, between the `skills` and `spec-audit` rows, so `TestCommandIndexNamesEveryCommandFile` keeps passing.

## Subtasks

- [ ] Parse the subcommand, resolve the Spec and read the key.
- [ ] Render the text and JSON forms.
- [ ] Add the help text.
- [ ] Document the command, then describe it in the skill and record its version.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] Surface Transcripts 1 to 6 are reproduced byte for byte.
- [ ] `--format json` prints schema `roundfix/spec-judge/v1` with every field of API Contract 2, including clear judgments.
- [ ] A stopped, skipped or advisory run exits `0`; `--stage qa`, `--format yaml` and an unknown Spec exit `2`.
- [ ] With `ROUNDFIX_OPENROUTER_API_KEY` and `ROUNDFIX_TYPESAFE_API_KEY` present in the process but absent from the command's environment, which holds only `OPENROUTER_API_KEY`, the command prints Surface Transcript 2 and the fake transport receives nothing.
- [ ] Surface Transcript 1 runs on the OpenRouter transport and Surface Transcript 5 on the TypeSafe transport, each summary naming its transport.
- [ ] The command reference and the skill reference carry their new sections, and the skill's version changed and is recorded.

## Context

- creates: `internal/cli/spec_judge.go`
- creates: `internal/cli/spec_judge_test.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/cli/cli.go`
- creates: `docs/user-guide/commands/spec.md`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- creates: `.agents/skills/roundfix/references/spec.md`
- creates: `skills/roundfix/references/spec.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/cli/spec_check_test.go`
- instruction: `internal/cli/cli_test.go`

## Verification

- `grep -qE '^[|] .spec. [|] [[]command guide[]][(]commands/spec[.]md[)] [|]$' docs/user-guide/commands.md && test -f docs/user-guide/commands/spec.md && go test -tags docscontract -count=1 -run '^(TestCommandIndexNamesEveryCommandFile|TestCommandReferenceLinksResolve)$' ./internal/docscontract`
- `out="$(go test -count=1 -v -run "^(TestSpecJudgeReportsAdvisoryJudgments|TestSpecJudgeSkipsWithoutAKey|TestSpecJudgeReadsTheKeyFromItsOwnEnvironment|TestSpecJudgeSkipsAtTheMonthlyCeiling|TestSpecJudgeSkipsANonEnglishSpec|TestSpecJudgeStopsWhenTheServiceFails|TestSpecJudgeRefusesAnUnknownSpec|TestSpecJudgeRefusesAnUnknownStageOrFormat|TestSpecJudgePrintsJSON|TestSpecJudgeHelp|TestRunSpecAuditHelpAppearsInUsageAndCommandList|TestRunCommandHelp|TestEveryOwnedSkillVersionIsRecorded)$" ./internal/cli ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSpecJudgeReportsAdvisoryJudgments TestSpecJudgeSkipsWithoutAKey TestSpecJudgeReadsTheKeyFromItsOwnEnvironment TestSpecJudgeSkipsAtTheMonthlyCeiling TestSpecJudgeSkipsANonEnglishSpec TestSpecJudgeStopsWhenTheServiceFails TestSpecJudgeRefusesAnUnknownSpec TestSpecJudgeRefusesAnUnknownStageOrFormat TestSpecJudgePrintsJSON TestSpecJudgeHelp TestRunSpecAuditHelpAppearsInUsageAndCommandList TestRunCommandHelp TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task ten of the thirteen named tests do not exist, so the command fails.
- `for pair in "docs/user-guide/commands/spec.md|### spec judge" "docs/user-guide/commands/spec.md|roundfix spec judge <slug>" "docs/user-guide/commands/spec.md|ROUNDFIX_OPENROUTER_API_KEY" "docs/user-guide/commands/spec.md|ROUNDFIX_TYPESAFE_API_KEY" "docs/user-guide/commands/spec.md|monthly ceiling" ".agents/skills/roundfix/references/spec.md|### Advisory judge" ".agents/skills/roundfix/references/spec.md|roundfix spec judge <slug> --stage prd" ".agents/skills/roundfix/references/spec.md|never gates"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; add="$(git log -1 --format=%H --diff-filter=A -- internal/cli/spec_judge_test.go)" || exit 1; base="${add:+$add^}"; base="${base:-HEAD}"; tmp="$(mktemp -d)" || exit 1; git show "$base:.agents/skills/roundfix/SKILL.md" > "$tmp/before.md" || exit 1; awk 'substr($0,1,8)=="version:" {print; exit}' "$tmp/before.md" > "$tmp/version-old"; awk 'substr($0,1,8)=="version:" {print; exit}' .agents/skills/roundfix/SKILL.md > "$tmp/version-new"; if cmp -s "$tmp/version-old" "$tmp/version-new"; then printf 'the skill version did not change\n' >&2; exit 1; fi; diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task no guide or skill file names `roundfix spec judge`, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1 and 2; User Stories 1, 2 and 3; Core Features 1, 6, 8 and 11; Success Metrics 1 and 2; Declared breaks
- [_techspec.md](_techspec.md) — Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; API Contract 1; API Contract 2; API Contract 5; Testing Approach 5; Testing Approach 6; Build Order 3
- ADR-0200; ADR-0201; ADR-0089; ADR-0184; ADR-0187; ADR-0189

## Result

Implemented the Task 03 command slice for Daemon Verification. Task status,
the Task Graph, and other Task files were not edited. The starting worktree
already held the Daemon's `pending` → `in_progress` status change in this
file; no implementation files for this Task existed and `runSpecCommand`
had no `judge` case.

### Implementation and acceptance evidence

- Surface Transcripts 1–6: the six transcript tests assert stdout, stderr,
  and exit code byte for byte through `runWithContext`, using real temporary
  repositories, temporary Roundfix Homes, fixed UTC clocks, and injected
  fake transports. `TestSpecJudgeReportsAdvisoryJudgments`,
  `TestSpecJudgeSkipsWithoutAKey`, `TestSpecJudgeSkipsAtTheMonthlyCeiling`,
  `TestSpecJudgeSkipsANonEnglishSpec`, `TestSpecJudgeStopsWhenTheServiceFails`,
  and `TestSpecJudgeRefusesAnUnknownSpec` passed in the focused check.
- JSON: `TestSpecJudgePrintsJSON` passed, asserting schema
  `roundfix/spec-judge/v1`, every document field including `transport`, every
  judgment field, five judgments including clear results, and applicable
  nulls. `TestSpecJudgeJSONSkipHasNullTransport` additionally passed for the
  no-key JSON document, null transport/stopped fields, and an empty artifact
  skip array. JSON encodes `judge.Report` directly rather than duplicating
  its schema.
- Exit codes and validation: advisory, run-level skips, and service stops
  exit `0` in the transcript tests. `TestSpecJudgeRefusesAnUnknownStageOrFormat`
  passed for `--stage qa`, `--format yaml`, unknown flags, missing slug,
  missing flag value, and extra slug (exit `2`). Unknown Spec exits `2` with
  Transcript 6's diagnostic. `TestSpecJudgeRequiresStageArtifacts` passed
  for a missing PRD and a missing TechSpec at `--stage techspec`.
- Environment isolation: `TestSpecJudgeReadsTheKeyFromItsOwnEnvironment`
  passed with both Roundfix keys set to fake values in the process and only
  `OPENROUTER_API_KEY` in the command environment. It asserts Transcript 2
  and its fake transport fails on any request. Production builds
  `Request.Keys` solely from the two Roundfix names in `environment.environ`.
  New dependencies `judgeTransport` and `judgeNow` default to
  `http.DefaultTransport` and `time.Now`, and the command passes them and its
  Roundfix Home to `judge.Run`.
- Transport and cost: Transcript 1's fake validates the OpenRouter host and
  command key, returns `typesafe/jev-1.13-20260917` with reported cost, and
  asserts five log lines under the temporary Home with that model ID.
  Transcript 5's fake validates the TypeSafe host and command key, returns
  `jev-1.13.0`, and asserts exactly two calls before stopping. Their summaries
  name `openrouter` and `typesafe` respectively. Costs have four decimal
  places; confidence and delivery probability have two.
- Spec Root, individual skips, and help: `TestSpecJudgeUsesConfiguredSpecRoot`
  passed with an external configured root.
  `TestSpecJudgePrintsIndividualSkip` passed for a missing cited ADR.
  `TestSpecJudgeHelp`, `TestRunSpecAuditHelpAppearsInUsageAndCommandList`,
  and `TestRunCommandHelp` passed. The new synopsis is present in top-level,
  Spec, and judge help, and existing audit help strings remain intact.
- Documentation and skill: created `docs/user-guide/commands/spec.md` with
  `### spec judge` and inserted its command index row between `skills` and
  `spec-audit`. Added `### Advisory judge` in the Roundfix skill's Spec
  reference. The guide covers stages, key order and generic-key exclusion,
  the Judge Log, ceiling, skip/stop reasons, and advisory exit behavior. The
  skill tells authors to correct or defend every advisory and distinguish
  skips from clean judgments. Neither QA settlement section changed.
  Both Roundfix version fields rose one patch step from starting commit
  `0.1.7` to `0.1.8`; version recording added its content digest.

### Focused checks

- `GOCACHE=/tmp/roundfix-task03-gocache go test ./internal/cli -run '^TestSpecJudge' -count=1`
  — exit `0` after fixing the fake's goal question ID to read the embedded
  question configuration. Before that fixture correction, Transcript 1
  failed with two `unreadable answer` skips; the expected transcript was
  preserved.
- After both sabotages were restored and the additional boundary tests added:
  `GOCACHE=/tmp/roundfix-task03-gocache go test ./internal/cli -run 'TestSpecJudge|TestRunSpecAuditHelpAppearsInUsageAndCommandList|TestRunCommandHelp' -count=1`
  — exit `0` (`roundfix/internal/cli`, 1.120s).
- `GOCACHE=/tmp/roundfix-task03-gocache go test -tags docscontract ./internal/docscontract -run '^TestCommand' -count=1`
  — exit `0` (`roundfix/internal/docscontract`, 0.694s), exercising command
  index and reference contracts as focused documentation checks.
- A Python inspection asserted all required guide/reference phrases, both
  `0.1.8` version fields, and byte equality of every canonical Roundfix skill
  file against its shipped mirror — passed.
- `git -c core.fsmonitor=false diff --check` — exit `0`.

### Regeneration

- `make skills-sync` — exit `0`; rewrote `skills/roundfix/SKILL.md` and
  `skills/roundfix/references/spec.md`.
- `GOCACHE=/tmp/roundfix-task03-gocache go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  — exit `0`; rewrote `skills/testdata/owned-skill-versions.json` with
  Roundfix `0.1.8` and its digest.
- `make baseline-digests` — exit `0`; reported `changed: false`, derived
  artifacts already matched their sources, and rewrote no files.

### Sabotage evidence

1. Exit-code gate: temporarily changed `runSpecJudgeCommand` to return
   `exitRunFailed` when `report.Stopped != nil`.
   `GOCACHE=/tmp/roundfix-task03-gocache go test ./internal/cli -run '^TestSpecJudgeStopsWhenTheServiceFails$' -count=1`
   exited `1`; the test failed with `exit=1 want=0` while its stdout matched
   Transcript 5. Restored the production source immediately afterward.
2. Key-source gate: temporarily changed the key loop from
   `environment.environ` to `os.Environ()`.
   `GOCACHE=/tmp/roundfix-task03-gocache go test ./internal/cli -run '^TestSpecJudgeReadsTheKeyFromItsOwnEnvironment$' -count=1`
   exited `1`; its fake transport failed with `unexpected request`.
   Restored the production source immediately afterward. The final focused
   command suite passed with both sabotages removed.

### Scope and limitations

No Task Verification command was run; authored Verification and settlement
remain Daemon-owned. No commit, push, PR, or live inference request was made.
The TypeSafe skill's live index at `https://docs.typesafe.ai/llms.txt` was
read through the web tool after the shell HTTP fetch was blocked; its linked
`https://docs.typesafe.ai/api.md` returned an internal fetch error. This
slice uses task_02's existing request/report contract and changes no service
protocol or judgment threshold. There are no follow-up implementation changes
in this diff.

## Carry-forward provenance

- Source Run: `run_20261001T215854Z_cdd052c5b8440585`
- Source commit: `e846a775c0c1e14572c9ef57d087196776ce3e6e`
