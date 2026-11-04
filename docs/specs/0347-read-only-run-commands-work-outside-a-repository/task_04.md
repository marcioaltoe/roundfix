---
task: task_04
spec: 0347-read-only-run-commands-work-outside-a-repository
status: completed
type: docs
complexity: low
---

# Task 04: The Roundfix skill, the configuration guide and the glossary state the Read-Only Run Command contract

## Overview

The Roundfix skill still tells agents that `runs list` exits `2` outside a
repository, and nothing says which commands tolerate a refused
`worktree.location` or how to find the `.git` the refusal comes from. This
Task writes the contract tasks 02 and 03 built into the Roundfix `runs`,
`events` and `review-runs` references, the configuration guide and
`CONTEXT.md`, which gains **Read-Only Run Command**. It answers the
maintainer report of 2026-11-03, adopted as
[references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md](references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md).

## Requirements

1. MUST change the canonical `.agents/skills/roundfix/references/runs.md`
   (`_techspec.md` → API Contract 5):
   - replace "Outside a Git repository, `runs list` without `--all` exits `2`
     and names `--all` as the alternative." with the outside-repository
     listing: outside a Git repository `runs list` lists every repository, as
     `--all` does, and exits `0`;
   - add a paragraph naming the Read-Only Run Commands (bare `runs`,
     `runs list`, `runs show`, `attach`, `events`), the warning line
     `roundfix: warning: <refusal>; commands that create a Run Worktree refuse
     it`, and that `runs causes` keeps requiring a repository;
   - make the closing sentence of the Run Browser paragraph say that `runs
     list` is repository-scoped inside a repository.
2. MUST add to the canonical `.agents/skills/roundfix/references/events.md`
   that `events` is a Read-Only Run Command and writes that warning to stderr,
   never to the JSONL on stdout.
3. MUST add to the canonical `.agents/skills/roundfix/references/review-runs.md`,
   at the `worktree.location` bullet, and to `docs/user-guide/configuration.md`,
   next to the `worktree` block, that a location inside the current Git root
   is refused with a message that names the root and the `.git` it found
   ("is inside the Git root"), that a stray `.git` in a parent directory such
   as the home directory is found this way, and that Read-Only Run Commands
   only warn.
4. MUST add **Read-Only Run Command** to `CONTEXT.md` through
   domain-modeling, near **Run Browser**, as a `**Read-Only Run Command**:`
   line: a command that only reads the Run Database (bare `runs`, `runs
   list`, `runs show`, `attach`, `events`) and never creates a Run Worktree,
   so it runs from any directory; outside a Git repository `runs list` lists
   every repository, and a `worktree.location` refusal becomes one stderr
   warning; `runs causes` is not one; linking ADR-0357 the way the glossary
   links ADRs. It keeps this repository's
   glossary contract (no undefined capitalized phrase, a plain `_Avoid_`
   line).
5. MUST NOT add or change text inside the `### QA settlement` section of any
   skill, and MUST NOT touch Go code, any other guide under `docs/`, any
   Baseline module, the Makefile, CI workflows, `go.mod` or `go.sum`.
6. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of the `roundfix` `SKILL.md` pair with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   recording `skills/testdata/owned-skill-versions.json`; never hand-edit a
   mirror or a version field.
7. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last edit of the configuration guide,
   and never hand-edit it.

## Subtasks

- [ ] Rewrite the runs, events and review-runs references.
- [ ] State the refusal and the warning in the configuration guide.
- [ ] Add the glossary entry.
- [ ] Sync the mirrors, record the versions, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of the three references state the contract and no longer say that `runs list` exits `2` outside a repository.
- [ ] The configuration guide names the Git root in the refusal and `CONTEXT.md` defines **Read-Only Run Command**.
- [ ] The mirrors equal their canonical copies and the owned-skill version record passes.

## Context

- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `docs/adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/references/events.md`
- interface: `skills/roundfix/references/events.md`
- interface: `.agents/skills/roundfix/references/review-runs.md`
- interface: `skills/roundfix/references/review-runs.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- The Roundfix skill texts and the owned-skill version record — consumers: `./skills`
- The configuration guide fingerprint and the glossary contract over `CONTEXT.md` — consumers: `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md .agents/skills/roundfix/references/events.md skills/roundfix/references/events.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Read-Only Run Command' || { echo "missing contract in $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md; do ! tr -s '[:space:]' ' ' < "$f" | grep -q -- 'without .--all. exits .2.' || { echo "old refusal still in $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/review-runs.md skills/roundfix/references/review-runs.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'is inside the Git root' || { echo "missing refusal in $f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md && cmp .agents/skills/roundfix/references/events.md skills/roundfix/references/events.md && cmp .agents/skills/roundfix/references/review-runs.md skills/roundfix/references/review-runs.md && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task no reference names a Read-Only Run Command, so the first check fails; afterwards both copies state the contract, the mirrors equal their canonical copies, and the owned-skill version record and the settlement guidance contracts pass.
- `grep -F -- '**Read-Only Run Command**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Read-Only Run Command\*\*(:| —( |$))' || { echo 'missing term in CONTEXT.md' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'is inside the Git root' || { echo 'missing refusal in the configuration guide' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the glossary has no such definition, so the first check fails; afterwards the glossary and the guide carry the contract and the documentation contracts, the glossary contract and the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Feature 5; Success Metric 5; Glossary
- `_techspec.md` → Current behavior 5; API Contract 5; Glossary; Testing Approach; Build Order 4

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, roundfix, evidence-gate

- Rewrote the canonical `runs`, `events`, and `review-runs` references with the Read-Only Run Command contract, the outside-repository `runs list` behavior, stderr warning boundary, and Git-root/`.git` refusal details.
- Added the same refusal and warning contract beside the `worktree` configuration block and defined **Read-Only Run Command** in `CONTEXT.md` with ADR-0357 and a plain `_Avoid_` line.
- `make skills-sync`: passed; copied the canonical Roundfix references to the maintained mirrors.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed (ok roundfix/skills 0.583s); raised both Roundfix skill version fields to 0.1.181 and recorded `skills/testdata/owned-skill-versions.json`.
- `git diff --check`: passed with no whitespace errors.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed (ok roundfix/internal/docscontract 0.653s); re-recorded `docs/references/behavior-surfaces.json` after the configuration-guide edit.

Acceptance evidence:

- Both canonical reference copies and their generated mirrors carry the contract after `make skills-sync`; the prior outside-repository exit-2 wording is removed from the canonical `runs` reference.
- The configuration guide names the Git root and .git path in the refusal, and `CONTEXT.md` defines **Read-Only Run Command**.
- The owned-skill version record and Behavior Surface Record were regenerated only through their sanctioned commands; the Behavior Surface command will be repeated as the final no-edit confirmation.
