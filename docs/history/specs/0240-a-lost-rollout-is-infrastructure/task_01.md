---
task: task_01
spec: 0240-a-lost-rollout-is-infrastructure
status: completed
type: docs
complexity: low
---

# Task 01: The glossary, the guides and the Roundfix Skill describe the Lost Rollout, its recovery and its park

## Overview

Describes the behavior that task_03 builds. The places are `CONTEXT.md`, the
fallback passages of the configuration and usage guides, the `deliver` guide,
and the Roundfix Skill's `profiles` and `deliver` references with their
mirrors. Then the skill ships with the CLI behavior. This Task answers the
Backlog Entry "A lost Codex rollout fails the Task and spends a retry" of
2026-10-06. It is verifiable on its own by the phrases, the mirrors and the
recorded skill version.

## Requirements

1. MUST add to `CONTEXT.md`, following the `domain-modeling` skill, the terms
   **Lost Rollout** and **First Handoff**, in the glossary's
   `**Term**:` / definition / `_Avoid_:` shape. **Lost Rollout** is the runtime
   infrastructure failure (`rollout_lost`) in which the ACP Runtime cannot find
   the persisted record of the Agent Session it is asked to continue. The Daemon
   recovers it in a new Agent Session without repair or a counted queue retry
   (ADR-0245). **First Handoff** is the first Agent turn of a Task that returns
   to the Daemon (ADR-0057's "Agent handoff"); before it, a Lost Rollout lets
   the Fallback Chain take the Task.
2. MUST amend the **Fallback Chain**, **Fallback Selection** and **Agent Work
   Started** entries so that each states the one exception. A Lost Rollout
   before the First Handoff, or in a QA Task before its report, lets the next
   Fallback Selection take the Task. The amended text must contain the exact
   phrase "before the First Handoff".
3. MUST replace, in `docs/user-guide/configuration.md`,
   `docs/user-guide/usage.md` and `.agents/skills/roundfix/references/profiles.md`,
   the sentences saying that no session-loss failure can start a replacement
   session. The new text describes `_techspec.md` → API Contract 2: a Lost
   Rollout gets no repair; the Task continues in a new Agent Session; the next
   Fallback Selection takes it before the First Handoff or while a QA report is
   pending, and the same selection continues otherwise; at most two recoveries
   happen per Agent Session owner; the third settles the Task as runtime
   infrastructure. Each of the three files must contain the exact phrase
   "Lost Rollout".
4. MUST describe in `docs/user-guide/usage.md` the Run Event of Invariant 9:
   the `daemon.task` phase `rollout_lost`, its `recovery` values and
   `retry_spent: false`, and the QA report section `## Agent runtime fallback`
   of Invariant 11.
5. MUST add `runtime-infrastructure` to the `environment` row of the Park
   Class table in `docs/user-guide/commands/deliver.md` and in
   `.agents/skills/roundfix/references/deliver.md`. Each file must also state,
   with the exact phrase "is not counted", that a Delivery Retry from that
   blocker is not counted against the retry limit and re-enters as a
   `run-unresolved` retry does (API Contract 4).
6. MUST raise the Roundfix Skill's version in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md` from the value it holds when this Task
   starts, run `make skills-sync`, and record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   Running those commands is an implementation step, never part of
   Verification, and no version is written by hand.
7. MUST NOT change any other command reference, the skill's other text, the
   `### QA settlement` section of any skill, or another glossary entry.

## Subtasks

- [x] Add the two glossary terms and amend the three fallback entries.
- [x] Rewrite the replacement-session sentences in the three fallback passages.
- [x] Describe the Run Event and the QA report section in the usage guide.
- [x] Add the park to both Park Class tables.
- [x] Raise and record the skill version and sync the mirrors.

## Acceptance Criteria

- [x] `CONTEXT.md` defines **Lost Rollout** and **First Handoff**, and the
      fallback entries name the exception.
- [x] No guide or skill reference still says that a session-loss failure never
      starts a replacement session.
- [x] Both Park Class tables list `runtime-infrastructure` with its uncounted
      retry.
- [x] Every skill mirror equals its canonical file, and the raised version is
      recorded.

## Result

Implemented the Lost Rollout and First Handoff glossary terms, the fallback
exception, recovery behavior, Run Event and QA report documentation, and the
runtime-infrastructure park and uncounted retry contract. Raised the Roundfix
Skill from `0.1.42` to `0.1.43` through the owned-skill recorder and synced its
canonical files to the embedded mirrors.

Focused implementation evidence:

- `make skills-sync`: passed before and after recording the raised version.
- `GOCACHE=/tmp/roundfix-go-cache-recorder go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed; the recorder raised both front-matter fields and updated `skills/testdata/owned-skill-versions.json`.
- Focused phrase checks over the changed documentation: passed for both glossary terms, `before the First Handoff`, `Lost Rollout`, `rollout_lost`, `retry_spent: false`, `## Agent runtime fallback`, `runtime-infrastructure`, `is not counted`, and `run-unresolved`.
- Canonical/embedded skill mirror comparison: passed for `SKILL.md`, `references/profiles.md`, and `references/deliver.md`.
- `git diff --check`: passed.

The Task's authored Verification commands remain for the Daemon to run.

## Context

- interface: `CONTEXT.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/profiles.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/profiles.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0245-a-lost-rollout-is-runtime-infrastructure.md`
- instruction: `docs/adr/0114-opening-an-agent-session-is-not-agent-work.md`

## Verification

- `f() { tr -s '[:space:]' ' ' < "$1" | grep -qF -- "$2" || { printf 'missing in %s: %s\n' "$1" "$2" >&2; exit 1; }; }; f CONTEXT.md '**Lost Rollout**:'; f CONTEXT.md '**First Handoff**:'; f CONTEXT.md 'before the First Handoff'; f docs/user-guide/configuration.md 'Lost Rollout'; f docs/user-guide/usage.md 'Lost Rollout'; f docs/user-guide/usage.md 'rollout_lost'; f docs/user-guide/usage.md 'Agent runtime fallback'; f docs/user-guide/commands/deliver.md 'runtime-infrastructure'; f docs/user-guide/commands/deliver.md 'is not counted'; f .agents/skills/roundfix/references/profiles.md 'Lost Rollout'; f .agents/skills/roundfix/references/deliver.md 'runtime-infrastructure'; f .agents/skills/roundfix/references/deliver.md 'is not counted'; cmp .agents/skills/roundfix/references/profiles.md skills/roundfix/references/profiles.md || exit 1; cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md || exit 1; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0; before this Task no file carries the new terms or the park, so the command fails at its first phrase check, and after it the raised version must be recorded for the skill test to pass.
- `g() { tr -s '[:space:]' ' ' < CONTEXT.md | awk -v t="$1" '{ i = index($0, t); if (!i) exit 1; s = substr($0, i); j = index(s, "_Avoid_"); if (!j || !index(substr(s, 1, j), "First Handoff")) exit 1 }' || { printf 'entry lacks the exception: %s\n' "$1" >&2; exit 1; }; }; g '**Fallback Chain**:'; g '**Fallback Selection**:'; g '**Agent Work Started**:'` — expected: exit 0; before this Task none of the three entries names the First Handoff, so the command fails at its first entry, and after it each amended entry states the exception.

## References

- `_prd.md` → Core Feature 5; User Stories 1-3; Goals
- `_techspec.md` → Build Order 1; API Contract 2; API Contract 4; Invariant 9; Invariant 11; Vocabulary Contract
- ADR-0245
- ADR-0189
