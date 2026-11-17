# What the Phase Record and `roundfix runs causes` held on 2026-11-17

Recorded by the author of the Idea of Spec 0404 on the Mac mini on 2026-11-17, between about 14:10 and 15:10 BRT. Every command was read-only. The binary was the installed Roundfix 0.147.1, run from a clone of this repository at `74d0e3a9`.

## 1. The Phase Record

The directory `phase-record` under Roundfix Home held two day files, for the UTC dates 2026-11-16 and 2026-11-17.

- 15 phase entries. Their keys are `schema`, `kind`, `time`, `phase`, `spec`, `repository`, `harness`, `session` and `boundary`, and 11 of them also hold `previous_phase` and `previous_spec`. By phase: `techspec` 6, `tasks` 5, `follow-up` 2, `release` 1 and `idea` 1. They name Specs 0399 to 0404. All 15 hold the harness `claude` and the same session identifier.
- 45 prompt entries. Their keys are `schema`, `kind`, `time`, `started`, `ended`, `run_id`, `spec`, `scope_kind`, `scope_id`, `repository`, `work_dir`, `harness` and `session`. By scope: `task` 27, `pre-pr-review` 12 and `qa` 6, all with the harness `codex`.
- No entry of either kind holds a model or a reasoning effort.

`roundfix phase start idea --spec 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model`, run by the Author Subagent of this Idea, recorded the entry and printed the Boundary Notice naming the `tasks` phase of Spec 0403. The subagent's entry carries the session of the Supervisor that started it.

## 2. What a command run by the authoring session can read

The Author Subagent ran in Claude Code; its `AI_AGENT` variable names version 2.1.325. Its environment held `CLAUDE_CODE_SESSION_ID`, the variable `roundfix phase start` reads, and `CLAUDE_EFFORT` with the value `xhigh`. A search of every variable's value for `fable`, `opus` and `sonnet` found none. No other value was printed or recorded.

## 3. `roundfix runs causes`

The operator's capture of `roundfix runs causes --since 2026-11-09 --format json` on 2026-11-17, kept with the study, reports the schema `roundfix/runs-causes/v1`, 31 Runs and 10 items. All 10 are of kind `verification`: `implementation-defect` 4, `repository-convention` 4, `shared-section-contract` 1 and 1 unclassified. By Spec: 0386 2, 0393 1, 0397 1, 0399 2 and 0400 4. It lists no item of kind `corrective`.

The author ran the same command three times: 1.20 s, 1.07 s and 1.18 s.

## 4. The commits of a delivery

Pull request #892, the delivery of Spec 0399, held nine commits, in this order: three Task commits, each with a `Roundfix-Task` trailer; `f7a132d4` (`fix`); `d23cf4ba` (`docs`, a change to the QA plan); the QA Report commit; `a8b702c0` (`fix`); the archive commit; and the commit that records the pull request. The three commits named by their hash are the Supervisor's corrections, and none carries a `Roundfix-Task` trailer. The pull request reached the default branch as the single commit `67ae4b15`.

## Limits

- The record is one machine's and its first day file is of 2026-11-16, so it holds two days.
- Section 2 is one session of one harness version.
- Three runs do not make a baseline of the command's time; the PRD records one if a KPI rests on it.
