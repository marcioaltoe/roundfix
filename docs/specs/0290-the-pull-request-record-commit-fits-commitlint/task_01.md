---
task: task_01
spec: 0290-the-pull-request-record-commit-fits-commitlint
status: completed
type: docs
complexity: low
---

# Task 01: The deliver guide, the skill reference and the glossary give the record subject and the hand commit

## Overview

The `deliver` user guide and the Roundfix skill's `deliver` reference say the
Pull Request Record Commit changes only `pull_request`, but not what its
subject is or what to do when a commit hook refuses a commit the queue makes.
The glossary's Hook Refusal entry covers only a Task commit. This Task states
the subject and the recovery in both documents, raises the Roundfix skill
version, re-records the Behavior Surface Record and revises **Hook Refusal**
(`_prd.md` → Core Feature 4).

## Requirements

1. MUST add to `docs/user-guide/commands/deliver.md`, right after the sentence
   "The commit changes the record's `pull_request` field and the queue pushes
   it before judging the pull request checks.", this paragraph:
   "Its subject is `docs: record pull request #<number> in the <spec number>
   Archive Record`, where the Spec number is the run of digits that starts the
   slug, such as `0088`; a slug that does not start with digits followed by
   `-` gives `docs: record pull request #<number> in the Archive Record`. The
   subject stays within 100 characters for every slug. The queue accepts an
   existing record commit by its change, whatever its subject."
   Then this paragraph: "When a commit hook refuses a commit the queue makes
   (the archive commit, the Pull Request Record Commit or the derived merge
   commit), the item parks as `delivery-error` and the reason names the hook,
   the refused subject, the hook's exit status and output, and the recovery.
   The queue never bypasses the hook. To recover, commit the staged change in
   the item worktree with a subject the hook accepts, then run `roundfix
   deliver retry <slug>`. A refused derived merge is aborted, so merge the
   default branch into the item branch by hand as for a conflict park, commit
   it with a message the hook accepts, then retry. A hook that does not name
   itself is named from the installed `pre-commit` and `commit-msg` hooks."
2. MUST add the same two paragraphs to the canonical Roundfix reference
   `.agents/skills/roundfix/references/deliver.md`, right after its sentence
   "The commit changes the record's `pull_request` field and the queue pushes
   it before judging the pull request checks.". No text inside any `###
   QA settlement` section may change.
3. MUST revise **Hook Refusal** in `CONTEXT.md` through `domain-modeling`,
   keeping its existing sentences and `_Avoid_` line, by appending: "A
   Delivery Queue commit — the archive commit, the Pull Request Record Commit
   or the derived merge commit — that a hook refuses is a Hook Refusal too:
   the item parks `delivery-error` with a reason naming the hook, the refused
   subject and the hand commit that recovers it, and the queue never bypasses
   the hook (ADR-0302)." The ADR-0302 reference takes the glossary's link form, to
   `docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md`.
4. MUST run these commands, because this Task edits the Roundfix skill and a
   guide. Before them, raise both version fields of the Roundfix skill's
   canonical `SKILL.md` (`metadata.version` and the top-level `version`) by
   one patch over the value the tree holds when this Task starts, because the
   record command refuses changed content under an unchanged version. The
   Agent runs them itself before handing back, because the Daemon never
   regenerates a record; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     `deliver` guide's fingerprint changes.
5. MUST search the guides, the agent guides and the skill copies for wording
   that gives the record subject with the slug or says a refused queue commit
   parks with a Git error only, and report what it found.
6. MUST NOT change any Go source or ADR, `docs/references/skill-coverage.json`
   or any other guide page.

## Subtasks

- [ ] Add the subject and recovery paragraphs to the `deliver` guide.
- [ ] Add them to the Roundfix `deliver` reference.
- [ ] Revise **Hook Refusal** in the glossary.
- [ ] Raise the skill version, sync the mirror, record the version and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The guide and both copies of the reference give the record subject and the hand-commit recovery.
- [ ] `CONTEXT.md` says a refused Delivery Queue commit is a Hook Refusal and links ADR-0302.
- [ ] The Roundfix skill version is raised and recorded, the mirror is in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirror, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'docs: record pull request #<number> in the <spec number> Archive Record' || { printf 'missing subject in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'commit the staged change in the item worktree with a subject the hook accepts' || { printf 'missing recovery in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Hook Refusal**:' || { echo 'missing glossary term' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'that a hook refuses is a Hook Refusal too' || { echo 'missing glossary revision' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document gives the record subject and the glossary entry covers only Task commits, so the command fails at its first phrase check. Afterwards the guide and both reference copies give the subject and the recovery, the glossary revises **Hook Refusal**, the mirror is in sync, the skill version is recorded, and the documentation contracts, the Skill Coverage Map among them, pass.

## References

- `_prd.md` → Core Feature 4; Goals; Glossary
- `_techspec.md` → API Contract 7; Build Order 1
- ADR-0302

## Result

Implemented the documentation slice while leaving the Daemon-owned status
`in_progress` unchanged:

- Added the numeric Spec subject rule and the named commit-hook refusal and
  hand-commit recovery to `docs/user-guide/commands/deliver.md` and the
  canonical `.agents/skills/roundfix/references/deliver.md`.
- Revised the `Hook Refusal` glossary entry in `CONTEXT.md` with Delivery Queue
  commits and the ADR-0302 link, preserving its existing definition and
  `_Avoid_` line.
- Raised both canonical Roundfix skill version fields from `0.1.117` to
  `0.1.118`, then ran `make skills-sync`; the canonical and mirrored skill and
  reference files are byte-identical.
- Ran the required recorders successfully: `go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` and `go
  test -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
  They updated `skills/testdata/owned-skill-versions.json` and
  `docs/references/behavior-surfaces.json` through their owning generators.
- Focused checks passed: `git diff --check`; canonical/mirror `cmp`; no diff
  inside `### QA settlement`; and the changed-path audit found only the
  Task-declared documentation, skill, mirror and derived-record paths.
- The required wording search found no old slug-based record-subject or
  Git-only queue-refusal wording in `docs/user-guide`, `.agents/skills` or
  `skills`. Historical slug wording remains in the Spec's PRD/TechSpec,
  adopted reference and ADR, which are outside this Task's allowed changes.

Acceptance evidence:

- Guide and both reference copies contain the record subject and hand-commit
  recovery paragraphs, with the reference's `### QA settlement` unchanged.
- `CONTEXT.md` identifies a refused Delivery Queue commit as a Hook Refusal
  and links ADR-0302.
- Roundfix is recorded at `0.1.118`, the mirror is synchronized, and the
  Behavior Surface Record was re-recorded successfully.
