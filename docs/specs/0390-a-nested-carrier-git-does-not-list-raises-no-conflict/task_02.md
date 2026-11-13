---
task: task_02
spec: 0390-a-nested-carrier-git-does-not-list-raises-no-conflict
status: pending
type: docs
complexity: low
---

# Task 02: The baseline guide and reference state the rule, and the amendment notice lands

## Overview

The `baseline` guide and the Roundfix skill's `baseline` reference, which
operators and agents read, state that a nested carrier enters the inventory
only when Git lists it, and ADR-0070 carries the notice that ADR-0399 amends
it. Adopters receive the rule through the binary and this text through the
Roundfix skills.

## Requirements

1. MUST add Text A of `_techspec.md` → API Contract 5, word for word, at the
   anchor it names, to `docs/user-guide/commands/baseline.md` and to the
   canonical `.agents/skills/roundfix/references/baseline.md`. MUST NOT add
   or change text inside any `### QA settlement` section, and MUST add no
   capitalized phrase that `CONTEXT.md` does not define.
2. MUST append Text B of API Contract 6, word for word, as the last paragraph
   of `docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md`
   and set its `updated_at` to the time of the edit; MUST NOT change its
   other text.
3. MUST run `make skills-sync` to copy the reference to
   `skills/roundfix/references/baseline.md`, then raise the Roundfix skill
   version to the next version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   on the tree the Task runs on, after Spec 0389 raised it; never hand-edit a
   mirror, a version field or the record. Of the Governed Paths the Task may
   change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/baseline.md` and
   `skills/roundfix/SKILL.md`.
4. MUST, as its last step and after the last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Add Text A to the guide and the canonical reference.
- [ ] Append Text B to ADR-0070.
- [ ] Sync the mirror and record the next skill version.
- [ ] Re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The guide and both copies of the `baseline` reference carry Text A, and the two reference copies hold the same bytes.
- [ ] ADR-0070 carries Text B and stays `accepted`.
- [ ] The Roundfix skill version is raised and recorded; the skills and documentation contracts pass.

## Context

- instruction: `docs/adr/0399-the-baseline-inventories-a-nested-carrier-only-when-git-lists-it.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `baseline` reference and the skill version record — consumers: `./skills`
- The `baseline` guide and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `for f in docs/user-guide/commands/baseline.md .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'lists it: a tracked one, or an untracked one Git does not ignore. A copy in a path Git ignores, and a carrier inside a linked worktree, a nested repository or a submodule, is not a file of the repository: it raises no' || { echo "Text A missing: $f" >&2; exit 1; }; done; cmp -s .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md || { echo 'baseline reference mirror differs' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md | grep -qF -- '**Amended by [ADR-0399](0399-the-baseline-inventories-a-nested-carrier-only-when-git-lists-it.md) (2026-11-13).**' || { echo 'ADR-0070 notice missing' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0070-baseline-audits-all-carriers-but-preserves-root-instructions.md | grep -qF -- 'Root carriers, their backups and the warning-only handling of the nested carriers Git lists still hold.' || { echo 'ADR-0070 notice incomplete' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task none of the texts exists; afterwards the guide, both reference copies with the same bytes and ADR-0070 carry them, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record and the glossary phrase contract pass.

## References

- `_prd.md` → Goals; Core Feature 3; Success Metric 4
- `_techspec.md` → API Contracts 5 and 6; Build Order 2
- ADR-0399
- ADR-0070
