---
task: task_03
spec: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
status: pending
type: docs
complexity: low
---

# Task 03: The Roundfix skill, the runs guide and ADR-0357 state that only a .git Git accepts marks a Git root

## Overview

The Roundfix `review-runs` reference still says that repository detection
finds any `.git` entry in an ancestor, as the configuration guide did before
task_01, the `runs`
guide calls a found `.git` stray, the `runs` reference says nothing about
which `.git` counts, and ADR-0357's Decision still says that any `.git` entry
marks a repository. This Task writes the rule tasks 01 and 02 built into
those documents and the amendment notice into ADR-0357, as ADR-0362 records.
It answers the maintainer report of 2026-11-05, adopted as
[references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md](references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md).

## Requirements

1. MUST replace, in the canonical
   `.agents/skills/roundfix/references/review-runs.md` (at the
   `worktree.location` bullet) the sentence "Repository detection finds
   any `.git` entry in an ancestor, so a stray `.git` in a parent such as the
   home directory is found this way." with the rule sentence of
   `_techspec.md` → API Contract 5, word for word, and add the same sentence
   to the canonical `.agents/skills/roundfix/references/runs.md` after
   "Outside a Git repository, `runs list` lists every repository as `runs
   list --all` does and exits `0`."
2. MUST change `docs/user-guide/commands/runs.md` as API Contract 5 states:
   "including a stray `.git` in a parent such as the home directory." becomes
   "including one in a parent such as the home directory." followed by the
   rule sentence, and "A detected Git root, including a stray one, keeps"
   becomes "A detected Git root, including one in a parent directory, keeps".
3. MUST write into
   `docs/adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md`
   the amendment of API Contract 5: the paragraph that opens with the bold
   `Amended by ADR-0362 (2026-11-05).`, where `ADR-0362` links
   `0362-only-a-git-entry-that-git-accepts-marks-a-git-root.md`, after the
   Decision paragraph "How the Git root is found does not change...", the
   sentence "ADR-0362 takes the first half of this option, recognizing a
   `.git` only when Git would." at the end of the Considered Option
   "Recognize a `.git` only when Git would...", and
   `updated_at: 2026-11-05T12:00:00Z`; every other byte of ADR-0357 stays.
4. MUST NOT add or change text inside the `### QA settlement` section of any
   skill, and MUST NOT touch Go code, `CONTEXT.md`, any other guide under
   `docs/`, any other ADR, any Baseline module, the Makefile, CI workflows,
   `go.mod` or `go.sum`.
5. MUST copy the canonical skill with `make skills-sync`, then raise both
   version fields of the `roundfix` `SKILL.md` pair with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   recording `skills/testdata/owned-skill-versions.json`; never hand-edit a
   mirror or a version field.
6. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last edit of the `runs` guide, and never
   hand-edit it.

## Subtasks

- [ ] Replace the detection sentence in the references and the `runs` guide.
- [ ] Write the amendment notice into ADR-0357.
- [ ] Sync the mirrors, record the versions, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Both copies of the `review-runs` and `runs` references and the `runs` guide state the rule sentence, and none still says that detection finds any `.git` entry.
- [ ] ADR-0357 carries the amendment notice naming ADR-0362 and its Considered Option says which half ADR-0362 takes.
- [ ] The mirrors equal their canonical copies, and the owned-skill version record and the documentation contracts pass.

## Context

- instruction: `docs/adr/0362-only-a-git-entry-that-git-accepts-marks-a-git-root.md`
- interface: `.agents/skills/roundfix/references/review-runs.md`
- interface: `skills/roundfix/references/review-runs.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md`

## Consumers

- The Roundfix skill texts and the owned-skill version record — consumers: `./skills`
- The `runs` guide fingerprint and the ADR and glossary contracts — consumers: `./internal/docscontract`

## Verification

- `for f in .agents/skills/roundfix/references/review-runs.md skills/roundfix/references/review-runs.md .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md docs/user-guide/commands/runs.md; do t="$(tr -s '[:space:]' ' ' < "$f")"; printf '%s' "$t" | grep -qF -- 'is skipped as Git skips it' || { echo "missing rule in $f" >&2; exit 1; }; ! printf '%s' "$t" | grep -q -- 'finds any .\.git. entry in an ancestor' || { echo "old rule still in $f" >&2; exit 1; }; done; cmp .agents/skills/roundfix/references/review-runs.md skills/roundfix/references/review-runs.md && cmp .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task no copy states the rule sentence, so the first check fails; afterwards all five do, none keeps the old sentence, the mirrors equal their canonical copies, and the owned-skill version record passes.
- `adr=docs/adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md; t="$(tr -s '[:space:]' ' ' < "$adr")"; printf '%s' "$t" | grep -qF -- '**Amended by [ADR-0362](0362-only-a-git-entry-that-git-accepts-marks-a-git-root.md) (2026-11-05).**' && printf '%s' "$t" | grep -qF -- 'ADR-0362 takes the first half of this option' && grep -q '^updated_at: 2026-11-05T12:00:00Z$' "$adr" || { echo 'missing amendment in ADR-0357' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task ADR-0357 has no amendment, so the first check fails; afterwards it carries the notice and the date, and the documentation contracts, the glossary and the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 5
- `_techspec.md` → Current behavior 4; API Contract 5; Declared breaks; ADR Alternatives; Build Order 3
