---
task: task_01
spec: 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
status: pending
type: docs
complexity: low
---

# Task 01: The deliver guide, the skill references and the glossary give the archive subject and the title

## Overview

The `deliver` user guide and the Roundfix skill's `deliver` reference give the
Pull Request Record Commit subject but not the archive commit subject or the
pull request title, and the `archive-spec` skill says the Delivery Queue
always writes `docs: archive <slug>`. This Task states both rules in those
documents, raises both skills' versions, re-records the Behavior Surface
Record and adds **Spec Number** to the glossary (`_prd.md` → Core Feature 5).

## Requirements

1. MUST add to `docs/user-guide/commands/deliver.md`, right after the
   paragraph that ends "The queue accepts an existing record commit by its
   change, whatever its subject.", this paragraph:
   "The archive commit's subject is `docs: archive <slug>` while that is at
   most 100 characters, and otherwise `docs: archive the <spec number> Spec`,
   or `docs: archive the Spec` for a slug without a Spec number. The pull
   request title is `feat: deliver <slug>` while that is at most 91
   characters, which leaves room for the ` (#<number>)` a squash merge
   appends for any pull request number below 100000; otherwise it is `feat:
   deliver <spec number>: <Spec title>`, or `feat: deliver <Spec title>` for a
   slug without a Spec number, where the Spec title is the Archive Record's
   title cut at a word to fit. Lengths count UTF-16 code units, as commitlint
   does. The pull request body still names the full slug, and the queue
   accepts an existing archive commit by its change, whatever its subject."
2. MUST add the same paragraph to the canonical Roundfix reference
   `.agents/skills/roundfix/references/deliver.md`, right after its paragraph
   that ends "The queue accepts an existing record commit by its change,
   whatever its subject.".
3. MUST replace, in `.agents/skills/archive-spec/SKILL.md` under `## Steps`,
   the sentence "The Delivery Queue writes the commit subject `docs: archive
   <slug>` (Conventional Commits)." with "The Delivery Queue writes the
   commit subject `docs: archive <slug>` (Conventional Commits), or, when
   that would pass 100 characters, `docs: archive the <spec number> Spec`
   (`docs: archive the Spec` for a slug without a Spec number)." No text
   inside any `### QA settlement` section of any skill may change.
4. MUST add **Spec Number** to `CONTEXT.md` through `domain-modeling`, next to
   the **Pull Request Record Commit** entry, as: "**Spec Number**:" followed
   on the next line by "The one to ten ASCII digits that start a Spec slug
   when a `-` follows them, such as `0088`; a slug without them has no Spec
   Number. The Delivery Queue names the Spec by it in the Pull Request Record
   Commit subject, and in the archive commit subject and the pull request
   title when the slug form would pass 100 characters (ADR-0302, ADR-0309)."
   and the line "_Avoid_: Spec ID, spec prefix". Each ADR reference takes the
   glossary's link form, to
   `docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md`
   and
   `docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md`.
5. MUST run these commands, because this Task edits two owned skills and a
   guide. Before them, raise both version fields (`metadata.version` and the
   top-level `version`) of the canonical `SKILL.md` of the Roundfix skill and
   of the `archive-spec` skill by one patch over the value the tree holds
   when this Task starts, because the record command refuses changed content
   under an unchanged version. The Agent runs them itself before handing
   back, because the Daemon never regenerates a record; no mirror or record
   is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     which re-records `docs/references/behavior-surfaces.json` because the
     `deliver` guide's fingerprint changes.
6. MUST search the guides, the agent guides and the skill copies for wording
   that says the queue always writes `docs: archive <slug>` or `feat: deliver
   <slug>`, and report what it found;
   `docs/references/2026-10-12-squash-merged-archive-source-revisions.md`
   quotes historical squash subjects and stays as written.
7. MUST NOT change any Go source or ADR, `docs/references/skill-coverage.json`
   or any other guide page.

## Subtasks

- [ ] Add the archive subject and title paragraph to the `deliver` guide and the Roundfix `deliver` reference.
- [ ] Update the `archive-spec` skill's sentence.
- [ ] Add **Spec Number** to the glossary.
- [ ] Raise both skill versions, sync the mirrors, record the versions and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The guide and both copies of the reference give both archive subject forms and both title forms with their bounds.
- [ ] Both copies of the `archive-spec` skill give the number form of the archive subject.
- [ ] `CONTEXT.md` defines **Spec Number** and links ADR-0302 and ADR-0309.
- [ ] Both skill versions are raised and recorded, the mirrors are in sync, and the Behavior Surface Record is current.

## Consumers

None. This Task changes documentation, the skill mirrors, the skill version record and the Behavior Surface Record only; no code consumes them.

## Context

- instruction: `docs/adr/0302-a-queue-commit-names-the-spec-by-number-and-a-hook-refusal-parks-with-the-hand-commit.md`
- instruction: `docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md`
- instruction: `internal/docscontract/skill_coverage_test.go`
- instruction: `docs/references/skill-coverage.json`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'docs: archive the <spec number> Spec' || { printf 'missing archive subject in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'deliver <spec number>: <Spec title>' || { printf 'missing title in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'at most 91 characters' || { printf 'missing title bound in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/archive-spec/SKILL.md skills/archive-spec/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'archive the <spec number> Spec' || { printf 'missing archive subject in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Spec Number**:' || { echo 'missing glossary term' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document gives the number form of the archive subject or the title and the glossary lacks **Spec Number**, so the command fails at its first phrase check. Afterwards the guide and both reference copies give both rules, both `archive-spec` copies give the number form, the glossary defines **Spec Number**, the mirrors are in sync, both skill versions are recorded, and the documentation contracts, the Skill Coverage Map among them, pass.

## References

- `_prd.md` → Core Feature 5; Goals; Glossary
- `_techspec.md` → API Contract 7; Build Order 1
- ADR-0302
- ADR-0309
