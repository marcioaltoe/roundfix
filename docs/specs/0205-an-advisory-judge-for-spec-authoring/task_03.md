---
task: task_03
spec: 0205-an-advisory-judge-for-spec-authoring
status: pending
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
2. MUST read `ROUNDFIX_JEV_OPENROUTER_API_KEY` and `ROUNDFIX_TYPESAFE_API_KEY` only from the command environment's `environ`, never from the process, MUST NOT read `OPENROUTER_API_KEY`, MUST pass them to `judge.Run` as `Request.Keys`, and MUST pass Roundfix Home, a transport and a clock to `judge.Run` through two new fields of `commandDependencies`, `judgeTransport` and `judgeNow`, which default to `http.DefaultTransport` and `time.Now`.
3. MUST print the text form of `_techspec.md` → Surface Transcripts: one `advisory` line per raised judgment, one `skipped` line per skipped judgment or artifact, nothing for a clear judgment, and the summary line with its `stopped:` suffix or its run-level `skipped:` form, with costs to four decimals and confidence and probability to two. Surface Transcripts 1 to 6 MUST be reproduced byte for byte.
4. MUST print the JSON document of API Contract 2 with `--format json`, including `transport`, listing clear judgments too.
5. MUST add the help of API Contract 5 to `specUsage`, the top-level usage in `internal/cli/cli.go` and a new `specJudgeUsage`, keeping every string `TestRunSpecAuditHelpAppearsInUsageAndCommandList` and `TestRunCommandHelp` require.
6. MUST put the tests in `internal/cli/spec_judge_test.go`, each with a temporary repository, a temporary Roundfix Home, a fake `judgeTransport`, a fixed `judgeNow` and the key each Surface Transcript names set in the command's own environment (`ROUNDFIX_JEV_OPENROUTER_API_KEY` for Transcripts 1 and 4, `ROUNDFIX_TYPESAFE_API_KEY` for Transcripts 3 and 5, and only the generic `OPENROUTER_API_KEY` for Transcript 2), so no test reaches the network whatever the developer's shell exports.
7. MUST describe the command under a new section `### spec judge` in `docs/user-guide/commands/spec.md`: its synopsis, both stages, the transport order (`ROUNDFIX_JEV_OPENROUTER_API_KEY` for OpenRouter first, `ROUNDFIX_TYPESAFE_API_KEY` for TypeSafe directly second), that the generic `OPENROUTER_API_KEY` is not read so Roundfix's Jev cost stays on its own key, the Judge Log path, the monthly ceiling, the skip and stop reasons, and that it exits `0` whenever it ran.
8. MUST describe the command under a new heading `### Advisory judge` in `.agents/skills/roundfix/references/spec.md`: the synopsis, that it never gates and never changes another command's exit code, that each `advisory` line is answered by correcting the artifact or by stating why the text stands, and that a `skipped` result is neither a failure nor a clean result. It MUST add no text inside the `### QA settlement` section of any skill.
9. MUST raise the Roundfix Skill's version in both front-matter fields of `.agents/skills/roundfix/SKILL.md` by one patch step from the value on this Task's starting commit, following the version rule in force on that commit, then run `make skills-sync`, record the version with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, run `make baseline-digests`, and name in the Result every file those commands rewrote.
10. MUST prove each new gate can fail. The Result MUST record one sabotage of the exit code (for example exiting `1` on a stopped run) and one of the key source (for example reading the process environment, or reading `OPENROUTER_API_KEY`), each with the test that failed, and that the code was restored.

Spec 0194 creates `.agents/skills/roundfix/references/spec.md`, its mirror and `docs/user-guide/commands/spec.md` (`_prd.md` → Prerequisites). They are declared under `creates:` because they do not exist when this Spec is authored. This Task MUST NOT create one that Spec 0194 has not created: when such a file is absent, the Task stops and reports that Spec 0194 has not landed.

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
- [ ] With `ROUNDFIX_JEV_OPENROUTER_API_KEY` and `ROUNDFIX_TYPESAFE_API_KEY` present in the process but absent from the command's environment, which holds only `OPENROUTER_API_KEY`, the command prints Surface Transcript 2 and the fake transport receives nothing.
- [ ] Surface Transcript 1 runs on the OpenRouter transport and Surface Transcript 5 on the TypeSafe transport, each summary naming its transport.
- [ ] The command reference and the skill reference carry their new sections, and the skill's version changed and is recorded.

## Context

- creates: `internal/cli/spec_judge.go`
- creates: `internal/cli/spec_judge_test.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/cli/cli.go`
- creates: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- creates: `.agents/skills/roundfix/references/spec.md`
- creates: `skills/roundfix/references/spec.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/cli/spec_check_test.go`
- instruction: `internal/cli/cli_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestSpecJudgeReportsAdvisoryJudgments|TestSpecJudgeSkipsWithoutAKey|TestSpecJudgeReadsTheKeyFromItsOwnEnvironment|TestSpecJudgeSkipsAtTheMonthlyCeiling|TestSpecJudgeSkipsANonEnglishSpec|TestSpecJudgeStopsWhenTheServiceFails|TestSpecJudgeRefusesAnUnknownSpec|TestSpecJudgeRefusesAnUnknownStageOrFormat|TestSpecJudgePrintsJSON|TestSpecJudgeHelp|TestRunSpecAuditHelpAppearsInUsageAndCommandList|TestRunCommandHelp|TestEveryOwnedSkillVersionIsRecorded)$" ./internal/cli ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSpecJudgeReportsAdvisoryJudgments TestSpecJudgeSkipsWithoutAKey TestSpecJudgeReadsTheKeyFromItsOwnEnvironment TestSpecJudgeSkipsAtTheMonthlyCeiling TestSpecJudgeSkipsANonEnglishSpec TestSpecJudgeStopsWhenTheServiceFails TestSpecJudgeRefusesAnUnknownSpec TestSpecJudgeRefusesAnUnknownStageOrFormat TestSpecJudgePrintsJSON TestSpecJudgeHelp TestRunSpecAuditHelpAppearsInUsageAndCommandList TestRunCommandHelp TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task ten of the thirteen named tests do not exist, so the command fails.
- `for pair in "docs/user-guide/commands/spec.md|### spec judge" "docs/user-guide/commands/spec.md|roundfix spec judge <slug>" "docs/user-guide/commands/spec.md|ROUNDFIX_JEV_OPENROUTER_API_KEY" "docs/user-guide/commands/spec.md|ROUNDFIX_TYPESAFE_API_KEY" "docs/user-guide/commands/spec.md|monthly ceiling" ".agents/skills/roundfix/references/spec.md|### Advisory judge" ".agents/skills/roundfix/references/spec.md|roundfix spec judge <slug> --stage prd" ".agents/skills/roundfix/references/spec.md|never gates"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; add="$(git log -1 --format=%H --diff-filter=A -- internal/cli/spec_judge_test.go)" || exit 1; base="${add:+$add^}"; base="${base:-HEAD}"; tmp="$(mktemp -d)" || exit 1; git show "$base:.agents/skills/roundfix/SKILL.md" > "$tmp/before.md" || exit 1; awk 'substr($0,1,8)=="version:" {print; exit}' "$tmp/before.md" > "$tmp/version-old"; awk 'substr($0,1,8)=="version:" {print; exit}' .agents/skills/roundfix/SKILL.md > "$tmp/version-new"; if cmp -s "$tmp/version-old" "$tmp/version-new"; then printf 'the skill version did not change\n' >&2; exit 1; fi; diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task no guide or skill file names `roundfix spec judge`, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1 and 2; User Stories 1, 2 and 3; Core Features 1, 6, 8 and 11; Success Metrics 1 and 2; Declared breaks
- [_techspec.md](_techspec.md) — Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; API Contract 1; API Contract 2; API Contract 5; Testing Approach 5; Testing Approach 6; Build Order 3
- ADR-0200; ADR-0201; ADR-0089; ADR-0184; ADR-0187; ADR-0189
