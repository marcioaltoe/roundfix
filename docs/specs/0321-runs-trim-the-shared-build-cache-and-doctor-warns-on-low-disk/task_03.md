---
task: task_03
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
status: pending
type: docs
complexity: medium
---

# Task 03: Document the automatic trims and the disk-space warning

## Overview

The guides, the Roundfix skill references and the glossary say that only the
GC Command trims a Shared Build Cache, and nothing documents a free-space
warning. This Task writes the texts of `_techspec.md` → Documentation into
the GC, configuration, implement, deliver and doctor guides and the skill's
`storage`, `implement`, `deliver` and `setup` references, revises the
**Shared Build Cache** glossary entry, then records the skill's version and
the Behavior Surface fingerprint.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-26-the-shared-build-cache-fills-the-disk-between-manual-gc-runs.md`](references/2026-10-26-the-shared-build-cache-fills-the-disk-between-manual-gc-runs.md)
   of 2026-10-26 (`_prd.md` → Sources 1) by writing the documentation of
   `_techspec.md` → Documentation, describing the behavior of
   `_techspec.md` → Interfaces and ADR-0334.
2. MUST add the `### Automatic trims` subsection (paragraph 1, the
   console-line block of Interfaces in a `text` fence, paragraph 2) as the
   last part of the `## Shared Build Caches` section of
   `docs/user-guide/commands/gc.md` and of the canonical
   `.agents/skills/roundfix/references/storage.md`, and change the GC
   guide's Retention policy row as Documentation states.
3. MUST change the `worktree.shared_caches` key row of
   `docs/user-guide/configuration.md` and add its paragraph to that guide's
   `## Shared Build Caches` section as Documentation states.
4. MUST add the implement paragraph to `docs/user-guide/commands/implement.md`
   and the canonical `.agents/skills/roundfix/references/implement.md`, and
   the deliver paragraph to `docs/user-guide/commands/deliver.md` and the
   canonical `.agents/skills/roundfix/references/deliver.md`, each at the
   place Documentation names.
5. MUST add the `disk-space:` bullet and the `DR-DISK-LOW` and
   `DR-DISK-UNREAD` table rows to `docs/user-guide/commands/doctor.md` and the
   canonical `.agents/skills/roundfix/references/setup.md`, and name
   `disk-space` among the conditional checks in that reference's sentence.
6. MUST revise the **Shared Build Cache** entry of `CONTEXT.md` through the
   `domain-modeling` skill to the definition of Documentation, keeping its
   `_Avoid_` line.
7. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   (it raises both version fields of the `roundfix` skill and records them),
   then `make skills-sync` again so the mirrors match.
8. MUST, as the Task's last step, after every guide and skill edit,
   re-record the Behavior Surface fingerprint ("re-record after the last
   guide edit") with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST NOT change any Go source, any other skill, any ADR (ADR-0334 is
   written by the planning change), `docs/agents/`, `.roundfixrc.yml`, any
   Baseline module or derived Baseline file, the Makefile or the CI
   workflows, or the `### QA settlement` section of any skill.

## Subtasks

- [ ] Write the GC, configuration and storage texts.
- [ ] Write the implement, deliver, doctor and setup texts.
- [ ] Revise the glossary entry, sync and record the skill, then re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] Every guide, canonical reference and mirror named above carries its text of Documentation.
- [ ] `CONTEXT.md` carries the revised **Shared Build Cache** entry with its `_Avoid_` line.
- [ ] The `roundfix` skill's version fields are raised and recorded, and the Behavior Surface record is current.

## Consumers

None. This Task changes documents, skill text and records only; no Go
package reads them as a contract except the documentation contracts it runs.

## Context

- instruction: `docs/adr/0334-every-implement-run-and-the-queue-gate-trim-a-declared-shared-build-cache.md`
- instruction: `docs/adr/0327-gc-trims-a-declared-shared-build-cache.md`
- instruction: `skills/skills_test.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- interface: `docs/user-guide/commands/gc.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/implement.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/storage.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/storage.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'the GC Command, the start of every Implement Run and the Delivery Queue' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Shared Build Cache**:' && for doc in docs/user-guide/commands/gc.md .agents/skills/roundfix/references/storage.md skills/roundfix/references/storage.md; do grep -q '^### Automatic trims$' "$doc" && tr -s '[:space:]' ' ' < "$doc" | grep -qF -- 'Roundfix also trims every declared Shared Build Cache by itself' || { printf 'missing automatic trims: %s\n' "$doc" >&2; exit 1; }; done && tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'apply the same trim by themselves under a thirty-second budget' && for doc in docs/user-guide/commands/implement.md .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md; do tr -s '[:space:]' ' ' < "$doc" | grep -qF -- 'the Run trims it right after the Run is recorded' || { printf 'missing implement paragraph: %s\n' "$doc" >&2; exit 1; }; done && for doc in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$doc" | grep -qF -- 'queue trims each declared Shared Build Cache' || { printf 'missing deliver paragraph: %s\n' "$doc" >&2; exit 1; }; done && for doc in docs/user-guide/commands/doctor.md .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do tr -s '[:space:]' ' ' < "$doc" | grep -qF -- 'only when the filesystem holding the worktree location' && grep -q 'DR-DISK-LOW' "$doc" && grep -q 'DR-DISK-UNREAD' "$doc" || { printf 'missing disk-space text: %s\n' "$doc" >&2; exit 1; }; done && go test -timeout 10m -count=1 ./skills && go test -timeout 15m -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract` — expected: exit 0. Before this Task the glossary phrase is absent, so the first check fails; afterwards every text is present in each guide, canonical reference and mirror, the owned-skill version is recorded and the Behavior Surface record is current.

## References

- `_prd.md` → Sources; Core Feature 4; Success Metric 6; Glossary
- `_techspec.md` → Current behavior 7; Documentation; ADR Alternatives; Build Order 3
