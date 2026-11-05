---
task: task_04
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
status: completed
type: docs
complexity: medium
---

# Task 04: The skill, the guides, the glossary and ADR-0361 state the new behavior

## Overview

The Roundfix skill, the implement and runs guides and the glossary
still say the Task commit is scanned before the Task settles with no repair,
that `runs causes` reports only Verification attempts and corrective Tasks,
and nothing about the detail of a `run-unresolved` park. After this Task they
state what task_01, task_02 and task_03 make true, and ADR-0361 records the
decision. This answers the report reported by fiscus on 2026-11-04
([references](references/2026-11-04-fiscus-a-task-commit-secret-refusal-is-invisible.md)).

## Requirements

1. MUST edit the canonical `.agents/skills/roundfix/SKILL.md` section
   `### Settlement Checks` (never `### QA settlement`): replace everything
   from "Settlement Checks apply to every non-QA Task" up to, not including,
   the paragraph that starts "The checks inspect one Task's tree." with
   exactly this text, so the entry file stays within the 20000-byte budget
   `TestRoundfixSkillEntryFileStaysWithinItsBudget` enforces (it measured
   19943 bytes with this text applied):

   ```markdown
   Settlement Checks apply to every non-QA Task in a Task Graph, with or without a
   QA gate: the repository Verification and the secret scan always, the two
   in-process checks only with an authored QA gate Task. A Run without a QA gate
   also runs the Integrated Repository Verification under API Contract 5 and ends
   `Unresolved` with the API Contract 6 failure line when it fails. Each attempt
   runs them after its commands, whatever their result, in this order:

   1. The repository Verification, as the attempt's last command, when
      `verification.repository_at_settlement` is enabled.
   2. `settlement check: spec consistency`: new Spec Consistency findings.
   3. `settlement check: authorization`: the prospective Task commit against the
      frozen authorization record.
   4. `settlement check: secret scan`: the lines that commit would add.

   A failed check returns its diagnostics as Verification Feedback for the one
   repair turn; a Secret Match names path, line and kind, never the value. If the
   final attempt still fails, the Task settles `failed`, a final Secret Match with
   `Task commit carries a secret: <path>:<line>: <kind>`; a scan error halts the
   Run. The `verification.repository_at_settlement` switch turns off only the
   repository Verification at settlement.
   ```

2. MUST replace, in `.agents/skills/roundfix/references/implement.md` and in
   `docs/user-guide/commands/implement.md`, the paragraph that starts "Before
   each Task settles, the Secret Scan reads" with, verbatim apart from line
   wrapping and the reference's list indentation:
   "In every Verification attempt, after the Settlement Checks,
   `settlement check: secret scan` reads the lines the Task commit would add
   relative to `HEAD`. A Secret Match fails the attempt with `Task commit
   carries a secret: <path>:<line>: <kind>` and a `Fix:` line, so a match on
   the first attempt returns as Verification Feedback for the one repair turn
   and the second attempt scans again. A match on the final attempt settles
   the Task `failed` with that reason and creates no Task commit. The
   diagnostic, the prompt, the reason and the `secret_scan` event, which names
   its attempt, carry only path, line and kind; the Task's changes stay in the
   Run Worktree. A Git or read error halts the Run as an infrastructure
   failure."
3. MUST, in `docs/user-guide/commands/implement.md` section
   `#### Settlement Checks`, add the same fourth list item as requirement 1
   and, after the sentence that gives the reason `Settlement check failed:
   <label>: <first diagnostic line>; diagnostics: <path>`, add: "The secret
   scan runs in every graph, last, and its final reason is its first
   diagnostic line, `Task commit carries a secret: <path>:<line>: <kind>`."
4. MUST add after the paragraph on signature classes, in
   `.agents/skills/roundfix/references/runs.md` section `### Why Verification
   failed` and in `docs/user-guide/commands/runs.md` section `#### Why
   Verification failed`, verbatim apart from line wrapping: "A secret refusal
   of a Task commit is its own item of kind `secret-scan`, one per
   `secret_scan` event, with its Spec, Task, Run and attempt (`null` for an
   event an earlier release wrote) and its locations as the check; JSON adds
   `matches` with each path, line and kind, and never the value. Its class is
   `repository-convention`, set by its kind rather than by a signature, so its
   signature and evidence are null, and the attempt's failed
   `settlement check: secret scan` is not listed again as a `verification`
   item. Every other item has `matches: null`."
5. MUST add after the `secret-found` next-action paragraph of
   `.agents/skills/roundfix/references/deliver.md` the paragraph task_03
   wrote into `docs/user-guide/commands/deliver.md`, verbatim apart from line
   wrapping:
   "A `run-unresolved` park names each Task the Run settled `failed`, with its
   reason cut to 300 characters, and the Tasks the Run skipped:
   `run-unresolved: <task> failed: <reason>; skipped: <task>, <task>`. A Run
   with neither keeps the bare `run-unresolved` and its answer. When a failed
   Task's commit carried a secret, the next action is `if the failed Task's
   text prescribes the value the blocker locates, amend it in <worktree> to
   one no Secret Pattern matches (letters only, or fewer than 20 characters);
   then run roundfix deliver retry <slug>, whose rerun returns any Secret
   Match to the Agent as Verification Feedback`."
6. MUST, through `domain-modeling`, revise three `CONTEXT.md` entries, keeping
   each `_Avoid_:` line:
   - **Settlement Check**: its definition becomes "A gate fact the Daemon
     checks for a non-QA Task before it settles in a Task Graph, with or
     without a QA gate. The repository Verification at settlement and the
     Task commit's Secret Scan run with or without a QA gate; the two
     in-process checks run only with one. The checks cover the Spec
     Consistency Check, the repository Verification at settlement,
     authorization of the prospective Task commit and its Secret Scan, which
     runs last in every attempt", followed by a link to ADR-0361 in the form the
     entry's other links use.
   - **Secret Scan**: append "The Task commit's scan runs as the last
     Settlement Check of every Verification attempt, so a first-attempt
     Secret Match returns as Verification Feedback", followed by the same link
     to ADR-0361.
   - **Cause Class**: "It is assigned by signature, never by a model
     judgment" becomes "It is assigned by signature, or by kind for a
     `secret-scan` item, never by a model judgment", and the entry also links
     ADR-0361.
7. MUST create
   `docs/adr/0361-a-task-commit-secret-match-returns-as-verification-feedback-and-names-its-cause.md`
   from `_techspec.md` → ADR-0361, with the lifecycle overlay of
   `docs/agents/docs-layout.md`.
8. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md` and its mirror with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`.
   Never hand-edit a mirror or a version field.
9. MUST re-record the Behavior Surface Record as the LAST step, after the
   last edit of any guide or reference, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   and commit `docs/references/behavior-surfaces.json`.
10. MUST search `docs/user-guide/`, `docs/agents/`, `CONTEXT.md` and every
    copy of the Roundfix skill for text saying the Task commit's scan runs
    after Verification without a repair turn, or that `runs causes` reports
    only Verification attempts and corrective Tasks, and record in the Result
    each hit and whether it changed.
11. MUST NOT touch any Go file, any `### QA settlement` section, any other
    skill or reference, any Baseline asset, `docs/agents/`, the Makefile, CI
    workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Write the Settlement Checks and secret scan text into the skill and the implement guide.
- [ ] Write the Causes item and the park detail into the runs and deliver references and guides.
- [ ] Revise the glossary entries and create ADR-0361.
- [ ] Sync, record the skill version, re-record the Behavior Surface Record last.

## Acceptance Criteria

- [ ] Every copy of the skill and each guide states that the Task commit's scan is the last Settlement Check of every attempt, the `secret-scan` Causes item and the `run-unresolved` park detail.
- [ ] `CONTEXT.md` and ADR-0361 state the rules, the mirrors equal their canonical copies, and the skill version and documentation contracts pass.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/implement.md`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- creates: `docs/adr/0361-a-task-commit-secret-match-returns-as-verification-feedback-and-names-its-cause.md`

## Consumers

- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; printf '%s' "$t" | grep -q -- 'settlement check: secret scan.: the lines that commit would add' && printf '%s' "$t" | grep -qF -- 'a Secret Match names path, line and kind, never the value' || { printf 'missing Settlement Check text in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md docs/user-guide/commands/implement.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'returns as Verification Feedback for the one repair turn and the second attempt scans again' || { printf 'missing secret scan text in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md docs/user-guide/commands/runs.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'set by its kind rather than by a signature' || { printf 'missing Causes item text in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md docs/user-guide/commands/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'park names each Task the Run settled' || { printf 'missing park text in %s\n' "$f" >&2; exit 1; }; done; c="$(tr -s '[:space:]' ' ' < CONTEXT.md)"; printf '%s' "$c" | grep -qF -- '**Settlement Check**:' && printf '%s' "$c" | grep -qF -- 'authorization of the prospective Task commit and its Secret Scan' && printf '%s' "$c" | grep -qF -- '**Secret Scan**:' && printf '%s' "$c" | grep -qF -- 'a first-attempt Secret Match returns as Verification Feedback' && printf '%s' "$c" | grep -qF -- '**Cause Class**:' && printf '%s' "$c" | grep -qF -- 'or by kind for a' || { echo 'missing glossary revision in CONTEXT.md' >&2; exit 1; }; adr=docs/adr/0361-a-task-commit-secret-match-returns-as-verification-feedback-and-names-its-cause.md; test -f "$adr" && grep -q '^status: accepted' "$adr" || { echo 'missing ADR-0361' >&2; exit 1; }; for f in SKILL.md references/implement.md references/runs.md references/deliver.md; do cmp ".agents/skills/roundfix/$f" "skills/roundfix/$f" || exit 1; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no copy states the new behavior and ADR-0361 does not exist, so the first check fails; afterwards every copy does, the mirrors equal their canonical copies, and the skill version and documentation contracts pass.

## References

- `_prd.md` → Goals; Core Feature 6; Success Metric 6; Glossary
- `_techspec.md` → API Contract 10; ADR-0361; Build Order 4

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, evidence-gate

Implemented the documented settlement contract: the canonical Roundfix skill
now puts the Task commit Secret Scan last in every attempt; Implement documents
the one repair turn and second scan; Runs documents `secret-scan` Causes items;
Deliver references document the detailed `run-unresolved` park; `CONTEXT.md`
defines the updated Settlement Check, Secret Scan and Cause Class; and ADR-0361
records the decision with its lifecycle overlay. `make skills-sync` copied the
canonical skill and references to the mirrors, and the sanctioned skill-version
recorder raised both Roundfix version fields and updated its record.

Focused checks:

- `make skills-sync`: passed.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions -count=1`: passed (`roundfix/skills`, 0.505s) after the final skill edit.
- `git diff --check`: passed.
- Skill entry size: 19943 bytes for both canonical and mirror, under the 20000-byte budget.
- SHA-256 checks for the four canonical/mirror skill files matched pairwise.
- Stale-wording sweep: `.agents/skills/roundfix/SKILL.md` had the old three-check settlement list and changed; `skills/roundfix/SKILL.md` received that change through sync. `.agents/skills/roundfix/references/implement.md` and `docs/user-guide/commands/implement.md` had the post-attempt Secret Scan wording and changed; the mirror reference received it through sync. `.agents/skills/roundfix/references/runs.md`, `skills/roundfix/references/runs.md` and `docs/user-guide/commands/runs.md` lacked the secret-refusal Causes item and changed. `.agents/skills/roundfix/references/deliver.md` lacked task_03's park paragraph and changed; its mirror received it through sync. `docs/user-guide/commands/deliver.md` already had that paragraph and was unchanged. The old hits in `CONTEXT.md` were the three glossary definitions and all changed; no stale hit remained in `docs/agents/` or elsewhere in the searched skill trees.
- No Go file, Makefile, workflow, `go.mod` or `go.sum` changed; the assigned Task file retains `status: in_progress`.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed (`roundfix/internal/docscontract`, 0.435s); this re-recorded `docs/references/behavior-surfaces.json` after the final guide/reference edit.

The declared Task Verification commands, including the full skills and
`docscontract` suites, were not run because this is a Daemon-assigned turn.

## Carry-forward provenance

- Source Run: `run_20261105T015847Z_709cb0723a1fbaae`
- Source commit: `8151c5827b1e2982b4479002ca1c9bf94f4ce35f`
