---
task: task_01
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
status: completed
type: docs
complexity: low
---

# Task 01: The glossary, the guides and the Roundfix Skill describe the Provider System Error, the refused background tasks and the Unfinished Turn

## Overview

Describes the behavior that task_02, task_03 and task_04 build, so the skill
ships with the CLI behavior. The places are `CONTEXT.md`, the fallback
passages of the configuration and usage guides, the `deliver` command guide,
and the Roundfix Skill's `profiles` and `deliver` references with their
mirrors. It answers the two Backlog Entries of 2026-10-13,
[a model-at-capacity reply is accepted as a finished turn](references/2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md)
and
[a Claude QA turn ends while its work runs in the background](references/2026-10-13-a-claude-qa-turn-ends-while-its-work-runs-in-the-background.md).
It is verifiable on its own by the phrases, the mirrors and the recorded skill
version.

## Requirements

1. MUST add to `CONTEXT.md`, through the `domain-modeling` skill and in the
   glossary's `**Term**:` / definition / `_Avoid_:` shape, the three terms of
   `_techspec.md` → Glossary:
   - **Provider System Error**: the runtime infrastructure failure
     (`provider_system_error`) of a turn whose stream reports that the model
     provider failed it — the turn's last `threadStatus` is `systemError` —
     whatever its stop reason. It is recovered like a Lost Rollout, without a
     retry on the same selection (ADR-0274).
   - **Agent Background Task**: a command an Agent starts through its ACP
     Runtime's background facility, such as a Claude Code background Bash
     command or Monitor, which keeps running after its tool call returns.
     Roundfix refuses them in Run sessions and records each one a turn starts.
     Its `_Avoid_` line names "Background job", which the glossary keeps for
     the Detached Run.
   - **Unfinished Turn**: the outcome (`turn_unfinished`) of a QA turn that
     ends with its seeded report at `verdict: pending` after starting an Agent
     Background Task. The QA Task settles `failed` with a reason that names
     each task.
2. MUST amend the **Fallback Chain**, **Fallback Selection** and **Agent Work
   Started** entries so that each states that a Lost Rollout or a Provider
   System Error before the First Handoff, or in a QA Task before its report,
   lets the next Fallback Selection take the Task. Each amended entry must
   contain the exact phrase "Provider System Error".
3. MUST extend, in `docs/user-guide/configuration.md`,
   `docs/user-guide/usage.md` and
   `.agents/skills/roundfix/references/profiles.md`, the Lost Rollout passage
   with `_techspec.md` → Invariants 7, 8 and 11: a Provider System Error takes
   the next Fallback Selection under the same condition, has no same-selection
   retry, and otherwise settles the Task as runtime infrastructure. Each of the
   three files must contain the exact phrase "Provider System Error", and no
   guide, skill copy or glossary entry may still call the Lost Rollout "the
   one exception".
4. MUST describe in `docs/user-guide/usage.md` the `daemon.task` phases
   `provider_system_error` and `turn_unfinished` with their payload keys
   (API Contracts 2 and 4), the QA report line `- provider_system_error:`
   (Invariant 10), and that Roundfix starts every `claude` Agent Session with
   `CLAUDE_CODE_DISABLE_BACKGROUND_TASKS=1` and asks QA turns to run commands
   in the foreground (API Contract 5).
5. MUST state in `docs/user-guide/commands/deliver.md` and in
   `.agents/skills/roundfix/references/deliver.md` that an unrecovered
   Provider System Error also parks the item as `runtime-infrastructure`, with
   the blocker detail of API Contract 3, and keep the sentence that a retry
   from that blocker is not counted. Each file must contain the exact phrase
   "hit a provider system error".
6. MUST raise the Roundfix Skill's version in both front-matter fields of
   `.agents/skills/roundfix/SKILL.md` to the next version, run
   `make skills-sync`, and record it with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
   These are implementation steps, never part of Verification, and no version
   is written by hand.
7. MUST, after the guide and skill edits, record the Behavior Surface
   fingerprints with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   which rewrites only `docs/references/behavior-surfaces.json`, the way
   Spec 0266 recorded its guide changes. This is an implementation step,
   never part of Verification, and the record is never edited by hand.
8. MUST NOT change any other command reference, any other glossary entry, the
   `### QA settlement` section of any skill, the `qa-gate` skill, or
   `docs/references/skill-coverage.json`.

## Subtasks

- [ ] Add the three glossary terms and amend the three fallback entries.
- [ ] Extend the Lost Rollout passage in the three fallback descriptions.
- [ ] Describe the two Run Event phases, the report line and the `claude`
      refusal in the usage guide.
- [ ] Add the Provider System Error park to both `deliver` texts.
- [ ] Raise and record the skill version and sync the mirrors.
- [ ] Record the Behavior Surface fingerprints.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines the three terms, and the three fallback entries
      name the Provider System Error.
- [ ] The configuration and usage guides and the `profiles` reference describe
      the recovery; the usage guide describes both Run Event phases.
- [ ] Both `deliver` texts name the Provider System Error park.
- [ ] Every skill mirror equals its canonical file, and the raised version is
      recorded.

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
- interface: `docs/references/behavior-surfaces.json`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0274-a-turn-that-did-not-finish-is-not-accepted.md`
- instruction: `docs/adr/0245-a-lost-rollout-is-runtime-infrastructure.md`

## Consumers

None. This Task changes documentation and skill text only; no code reads it,
and the owned-skill version record is checked in this Task's own Verification.

## Verification

- `f() { tr -s '[:space:]' ' ' < "$1" | grep -qF -- "$2" || { printf 'missing in %s: %s\n' "$1" "$2" >&2; exit 1; }; }; f CONTEXT.md '**Provider System Error**:'; f CONTEXT.md '**Agent Background Task**:'; f CONTEXT.md '**Unfinished Turn**:'; f docs/user-guide/configuration.md 'Provider System Error'; f docs/user-guide/usage.md 'Provider System Error'; f docs/user-guide/usage.md 'provider_system_error'; f docs/user-guide/usage.md 'turn_unfinished'; f docs/user-guide/usage.md 'CLAUDE_CODE_DISABLE_BACKGROUND_TASKS=1'; f docs/user-guide/commands/deliver.md 'hit a provider system error'; f .agents/skills/roundfix/references/profiles.md 'Provider System Error'; f .agents/skills/roundfix/references/deliver.md 'hit a provider system error'; for d in .agents/skills/roundfix/references/profiles.md docs/user-guide/configuration.md docs/user-guide/usage.md CONTEXT.md; do if tr -s '[:space:]' ' ' < "$d" | grep -qF -- 'Lost Rollout is the one exception'; then printf 'still the one exception in %s\n' "$d" >&2; exit 1; fi; done; cmp .agents/skills/roundfix/references/profiles.md skills/roundfix/references/profiles.md || exit 1; cmp .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md || exit 1; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0; before this Task no file carries the new terms, so the command fails at its first phrase check, and after it the raised version must be recorded for the skill test to pass.
- `g() { tr -s '[:space:]' ' ' < CONTEXT.md | awk -v t="$1" '{ i = index($0, t); if (!i) exit 1; s = substr($0, i); j = index(s, "_Avoid_"); if (!j || !index(substr(s, 1, j), "Provider System Error")) exit 1 }' || { printf 'entry lacks the Provider System Error: %s\n' "$1" >&2; exit 1; }; }; g '**Fallback Chain**:'; g '**Fallback Selection**:'; g '**Agent Work Started**:'` — expected: exit 0; before this Task none of the three entries names the Provider System Error, so the command fails at its first entry.

## References

- [_prd.md](_prd.md) — Core Feature 4; User Stories 1-4; Goals 1-5
- [_techspec.md](_techspec.md) — Build Order 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Invariants 7-11; Glossary
- [2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md](references/2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md)
- [2026-10-13-a-claude-qa-turn-ends-while-its-work-runs-in-the-background.md](references/2026-10-13-a-claude-qa-turn-ends-while-its-work-runs-in-the-background.md)
- ADR-0274; ADR-0245; ADR-0189; ADR-0233; ADR-0244; ADR-0256

## Result

- Added the Provider System Error, Agent Background Task and Unfinished Turn
  glossary terms to `CONTEXT.md`; updated Fallback Chain, Fallback Selection
  and Agent Work Started with the Provider System Error eligibility rule.
- Extended the configuration guide, usage guide and canonical profiles
  reference with Provider System Error recovery, no same-selection retry and
  runtime-infrastructure settlement. The usage guide now records the
  `provider_system_error` and `turn_unfinished` event payloads, QA report line,
  Claude environment variable and foreground command rule.
- Added the unrecovered Provider System Error park detail to both deliver
  references, including the uncounted retry sentence.
- Ran `make skills-sync` successfully, then
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` successfully; the recorder raised the Roundfix Skill
  from `0.1.73` to the final recorded version `0.1.75` and updated its
  history. The sync and recorder were rerun after the final wording refinement.
- Ran `go test -count=1 -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` successfully;
  it was rerun after the final wording refinement and only the Behavior Surface
  record changed from that recorder.
- Focused checks: required phrases were found with `rg`; the deprecated
  `Lost Rollout is the one exception` wording was absent; canonical and mirror
  Roundfix Skill files compared equal with `cmp`; and `git diff --check`
  reported no whitespace errors.
- The Task status remains daemon-owned. The pre-existing `status:
  in_progress` change was preserved, and the Task's declared Verification is
  left for the Daemon.
- Verification Feedback attempt 1: `make verify-changed` reached the unrelated
  `internal/store` package but its parallel process-controller test timed out
  while proving graceful termination; all other listed packages passed. The
  changed-path diff contains no `internal/store` source or test. Focused
  reproduction `go test -count=1 ./internal/store -run
  '^TestOwnerProcessControllerReprovesBeforeForceKill$'` passed with 6 tests,
  so no out-of-scope repair was made.
