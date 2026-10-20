---
task: task_02
spec: 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
status: pending
type: docs
complexity: low
---

# Task 02: The runs reference and the glossary say where diagnostics are read

## Overview

The Roundfix `runs` reference does not say where `roundfix runs causes` reads
a failed attempt's diagnostic, and the glossary entry for **Artifact Root**
names three callers of the one resolver. This Task describes what task_01
built in the reference, raises the Roundfix skill version and names the run
cause report in **Artifact Root** (`_prd.md` → Core Feature 3). The `runs`
guide is task_01's.

## Requirements

1. MUST add to `.agents/skills/roundfix/references/runs.md`, in its "Why
   Verification failed" section, the sentence "Diagnostics are read under the
   Artifact Root the Run wrote to, through the resolver that `storage report`,
   `gc sanitize` and Run Retention share;" completed with "a diagnostic
   outside the Run's `runs/<run-id>` directory or behind a symbolic link gives
   no evidence." No text inside any `### QA settlement` section may change.
2. MUST revise **Artifact Root** in `CONTEXT.md` through `domain-modeling`:
   its sentence naming the resolver's callers becomes "One resolver serves the
   storage report, GC sanitation, Run Retention and the run cause report",
   keeping its ADR-0301 and ADR-0255 links and adding a link to ADR-0214; the
   rest of the entry and its `_Avoid_` line stay.
3. MUST run these commands, because this Task edits the Roundfix skill. Before
   them, raise both version fields of the Roundfix skill's canonical
   `SKILL.md` (`metadata.version` and the top-level `version`) by one patch
   over the value the tree holds when this Task starts, because the record
   command refuses changed content under an unchanged version. The Agent runs
   them itself before handing back; no mirror or record is hand-edited.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version.
4. MUST NOT change any Go source, guide page or ADR: the code and the `runs` guide belong to task_01.

## Subtasks

- [ ] Add the diagnostic sentence to the Roundfix `runs` reference.
- [ ] Revise **Artifact Root** in the glossary.
- [ ] Raise the skill version, sync the mirror and record the version.

## Acceptance Criteria

- [ ] The reference states the shared resolver and the guards, in the canonical copy and its mirror.
- [ ] `CONTEXT.md` names the run cause report among the resolver's callers.
- [ ] The Roundfix skill version is raised and recorded, and the mirror is in sync.

## Consumers

None. This Task changes documentation, the skill mirror and the skill version record only; no code consumes them.

## Context

- instruction: `internal/runcause/report.go`
- instruction: `docs/user-guide/commands/runs.md`
- instruction: `docs/adr/0301-gc-sanitation-classifies-the-artifact-root-each-run-wrote-to.md`
- instruction: `docs/adr/0214-why-a-verification-failed-is-classified-by-deterministic-signatures.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`

## Verification

- `for f in .agents/skills/roundfix/references/runs.md skills/roundfix/references/runs.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Diagnostics are read under the Artifact Root the Run wrote to' || { printf 'missing reference sentence in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Artifact Root**:' || { echo 'missing glossary term' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'GC sanitation, Run Retention and the run cause report' || { echo 'glossary does not name the run cause report' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither the reference sentence nor the revised glossary sentence exists, so the command fails at its first phrase check. Afterwards both copies of the reference and the glossary state where diagnostics are read, the mirror is in sync, the skill version is recorded, and the documentation contracts pass.

## References

- `_prd.md` → Core Feature 3; Glossary; Goals
- `_techspec.md` → API Contract 1; API Contract 3; Build Order 2
- ADR-0301; ADR-0214; ADR-0189; ADR-0256
