---
task: task_04
spec: 0233-a-light-tier-on-open-models
status: pending
type: backend
complexity: medium
---

# Task 04: The advisory judge suggests a model tier for each Task at the tasks stage

## Overview

The light tier is decided from the authored `complexity`, and the measurement
found the Jev judge's tier no better than it. The maintainer still wants the
judge's opinion visible while the Tasks are written, so a Task author can
reconsider a `complexity` before the planning Pull Request merges, without the
suggestion ever reaching dispatch. This Task adds a `tasks` stage to
`roundfix spec judge` that asks one `model-tier` Choice per non-QA Task and
prints it as a suggestion. It is verifiable on its own through the judge's
planner and the command with the existing fake transport.

## Requirements

1. MUST add a `model-tier` judgment to the judge's question catalog: a Choice
   asking for the cheapest model tier that would most likely pass the Task's
   Verification on the first attempt, with criteria `light`, `standard` and
   `heavy`, each defined in one sentence, and a bounded Task text length.
2. MUST accept `--stage tasks`, which requires `_tasks.md`, plans one
   `model-tier` judgment per Task file the graph lists whose `type` is not
   `qa`, sends the Task file text without its `status`, `complexity`,
   `## Result` and Daemon-written sections, and judges neither the PRD nor the
   TechSpec; without `--stage` the judge MUST still judge the PRD and the
   TechSpec only, byte-identical in output to today.
3. MUST report each answer with outcome `suggested` and print
   `suggested model-tier <task file>: <answer> at confidence <c>` in text
   output (API Contract 7), carry it in JSON as kind `model-tier`, append the
   usual Judge Log line per request with judgment `model-tier`, honor the
   existing key, ceiling and skip rules, and write no Spec file.
4. MUST change the unknown-stage message to
   `unsupported --stage "<value>"; use prd, techspec or tasks` (Surface
   Transcript 1) and every usage line to `[--stage <prd|techspec|tasks>]`.
5. MUST add tests: the planner plans one judgment per non-QA Task and none for
   the QA gate, strips the named sections, and plans nothing else at the
   `tasks` stage; through the CLI with the existing fake transport, `--stage
   tasks` prints one suggestion per non-QA Task, records one Judge Log line
   each in a disposable Home and leaves the Spec directory byte-identical, and
   with no key prints Surface Transcript 2; and MUST update the existing
   unknown-stage and help assertions to the new texts.
6. MUST NOT change dispatch, the configuration keys, any guide or skill, the
   Judge Log schema or the judge's transports; no test may reach the network
   or the real `~/.roundfix`.

## Subtasks

- [ ] Add the `model-tier` question.
- [ ] Plan and run the `tasks` stage.
- [ ] Render the suggestion and update the stage message and usage.
- [ ] Add the planner and command tests and update the existing assertions.

## Acceptance Criteria

- [ ] `--stage tasks` prints one `suggested model-tier` line per non-QA Task.
- [ ] Without `--stage` the judge's output is unchanged.
- [ ] An unknown stage names the three stages.
- [ ] No Spec file changes and every request leaves a Judge Log line.

## Context

- instruction: `docs/adr/0200-spec-authoring-gets-an-advisory-judge-that-never-gates.md`
- interface: `internal/judge/questions.json`
- interface: `internal/judge/questions.go`
- interface: `internal/judge/pairs.go`
- interface: `internal/judge/judge.go`
- interface: `internal/judge/questions_test.go`
- interface: `internal/cli/spec_judge.go`
- interface: `internal/cli/spec_judge_test.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/cli/cli.go`
- creates: `internal/judge/model_tier_test.go`
- creates: `internal/cli/spec_judge_tier_test.go`

## Verification

- `grep -q '"model-tier"' internal/judge/questions.json && grep -q 'use prd, techspec or tasks' internal/cli/spec_judge.go && go build -buildvcs=false ./... && go vet ./internal/judge ./internal/cli` — expected: exit 0; before this Task the catalog has no model-tier question and the stage message names two stages, so the command fails.
- `jd="$(go test -count=1 -v -run '^TestModelTierPlansNonQATasks$' ./internal/judge 2>&1)" || { printf '%s\n' "$jd"; exit 1; }; cl="$(go test -count=1 -v -run '^(TestSpecJudgeSuggestsModelTier|TestSpecJudgeRefusesAnUnknownStageOrFormat)$' ./internal/cli 2>&1)" || { printf '%s\n' "$cl"; exit 1; }; for name in TestModelTierPlansNonQATasks TestSpecJudgeSuggestsModelTier TestSpecJudgeRefusesAnUnknownStageOrFormat; do printf '%s\n%s\n' "$jd" "$cl" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new tests do not exist, so the command fails.

## References

- `_prd.md` → Goals; User Story 6; Core Feature 7; Success Metric 6
- `_techspec.md` → API Contract 7; Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 4
- ADR-0238; ADR-0200; ADR-0201
