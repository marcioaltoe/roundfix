---
task: task_03
spec: 0289-gc-sanitize-reaches-the-default-artifact-roots
status: completed
type: docs
complexity: low
---

# Task 03: The storage reference and the glossary name the resolved Artifact Root

## Overview

The Roundfix `storage` reference still says that `roundfix gc sanitize`
discovers every recorded Artifact Root and preserves ambiguous paths, and the
glossary has no entry for the root a Run writes to. This Task describes the
shared resolver and sanitation over resolved roots in the reference, raises the
Roundfix skill version, and adds **Artifact Root** to `CONTEXT.md`
(`_prd.md` → Core Feature 4).

## Requirements

1. MUST describe in `.agents/skills/roundfix/references/storage.md` what
   task_01 and task_02 built, replacing the paragraph that begins
   "`roundfix gc sanitize` discovers every recorded Artifact Root", with these
   sentences: "`roundfix storage report`, `roundfix gc sanitize` and Run
   Retention resolve a recorded Artifact Directory with one resolver." and
   "`roundfix gc sanitize` groups Runs by the Artifact Root each wrote to and
   classifies each root with its existing proofs.", plus the recheck before
   removal, the `unsafe` group of an unresolved Run, and the `Recorded as` and
   per-root lines. The storage report paragraph that Spec
   `reports-that-count-what-they-say` added stays as it is. No text inside any
   `### QA settlement` section may change.
2. MUST add **Artifact Root** to `CONTEXT.md` through `domain-modeling`,
   beside **Artifact Directory**. It is the directory under which a Run writes
   `runs/<run-id>`: the recorded Artifact Directory resolved the way the Run
   resolved it, a blank one to the repository's default root (or, for older
   Runs, the checkout's default when only it holds the Run's directory), a
   relative one under the Run's Git root and an absolute one as recorded. One
   resolver serves the storage report, GC sanitation and Run Retention. The
   entry links ADR-0301 and ADR-0255 and gives an `_Avoid_` line.
3. MUST run these commands, because this Task edits the Roundfix skill. Before
   them, raise both version fields of the Roundfix skill's canonical
   `SKILL.md` (`metadata.version` and the top-level `version`) by one patch
   over the value the tree holds when this Task starts, because the record
   command refuses changed content under an unchanged version. The Agent runs
   them itself before handing back; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version.
4. MUST NOT change any Go source, guide page or ADR: the code and the `gc`
   guide belong to task_01 and task_02.

## Subtasks

- [x] Rewrite the sanitation paragraph of the `storage` reference.
- [x] Add **Artifact Root** to the glossary.
- [x] Raise the skill version, sync the mirror and record the version.

## Acceptance Criteria

- [x] The reference states the one resolver and sanitation over resolved roots, in the canonical copy and its mirror.
- [x] `CONTEXT.md` defines **Artifact Root** with its resolution rules and ADR links.
- [x] The Roundfix skill version is raised and recorded, and the mirror is in sync.

## Consumers

None. This Task changes documentation, the skill mirror and the skill version record only; no code consumes them.

## Context

- instruction: `docs/user-guide/commands/gc.md`
- instruction: `docs/adr/0301-gc-sanitation-classifies-the-artifact-root-each-run-wrote-to.md`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`

## Verification

- `for f in .agents/skills/roundfix/references/storage.md skills/roundfix/references/storage.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'resolve a recorded Artifact Directory with one resolver' || { printf 'missing resolver sentence in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'groups Runs by the Artifact Root each wrote to' || { printf 'missing sanitation sentence in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Artifact Root**:' || { echo 'missing glossary term' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither sentence nor the glossary term exists, so the command fails at its first phrase check. Afterwards both copies of the reference state the resolver and sanitation over resolved roots, the glossary defines the term, the mirror is in sync, the skill version is recorded, and the documentation contracts pass.

## References

- `_prd.md` → Core Feature 4; Glossary; Goals
- `_techspec.md` → API Contract 1; API Contract 5; API Contract 6; API Contract 7; Build Order 3
- ADR-0301; ADR-0255; ADR-0189; ADR-0256

## Result

- Replaced the storage reference's sanitation description in the canonical copy
  and mirror with the shared resolver, resolved-root grouping, existing-proof
  classification, unresolved `unsafe` groups, pre-removal recheck, `Recorded
  as:` lines and per-root reclaimable output.
- Added **Artifact Root** beside **Artifact Directory** in `CONTEXT.md`, with
  blank, older-Run fallback, relative and absolute resolution rules, the shared
  resolver consumers, ADR-0301 and ADR-0255 links, and an `_Avoid_` line.
- Raised the canonical Roundfix skill from `0.1.116` to `0.1.117`; `make
  skills-sync` completed successfully, and
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions` reported `1 passed in 1 packages` and recorded the
  new digest in `skills/testdata/owned-skill-versions.json`.
- Focused checks: canonical and mirrored Roundfix skill files and storage
  references compare identical; `git diff --check` passed; the final diff
  contains only the authorized skill/reference files, `CONTEXT.md`, the skill
  version record and this Task file. The Daemon-owned declared Verification was
  not run.
