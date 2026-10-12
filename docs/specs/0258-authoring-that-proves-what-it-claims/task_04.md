---
task: task_04
spec: 0258-authoring-that-proves-what-it-claims
status: pending
type: docs
complexity: low
---

# Task 04: The command reference, the user guide and the glossary name the four proofs

## Overview

The four codes task_01 and task_02 add reach adopters through the roundfix
command reference and the user guide, with the paragraph of
`_techspec.md` → Skill and document texts that tells an adopter what an
active Spec completes before its next Run. The Behavior Surface Record is
re-recorded for the changed guide, and the glossary gains the six terms of
`_prd.md` → Glossary. It is verifiable on its own through phrase checks, the
mirrors, the version record and the surface record test.

This is an authorized tooling Task. It may change only the files in its
Context and this Task file.

## Requirements

1. MUST list `SC-ABSENCE-UNPROVEN`, `SC-ADR-ALTERNATIVES-UNRELATED`,
   `SC-COVERAGE-UNREAD` and `SC-CONSUMERS-UNPROVEN` with their stage in the
   roundfix command reference for `spec check` and in the user guide for the
   `spec` command, each with the paragraph of `_techspec.md` → Skill and
   document texts, so the adopter notes of `_prd.md` are true of the shipped
   documentation (ADR-0265, ADR-0256).
2. MUST write, through `domain-modeling`, **Absence Evidence**,
   **ADR Alternatives**, **Coverage Cross-Read**, **Consumer Declaration**,
   **Shared Contract** and **Integration Checkpoint** into `CONTEXT.md`, each
   as the bolded term, a colon, a one-sentence definition that links ADR-0265,
   and an `_Avoid_` line (ADR-0244).
3. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   then `make skills-sync` again, so that the roundfix skill records a raised
   version and its mirrors match.
4. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   after the user guide changes.
5. MUST NOT change any `### QA settlement` section, any other skill, the
   Skill Coverage Map or any production Go file.

## Subtasks

- [ ] Document the four codes and the adopter paragraph in both places.
- [ ] Write the six glossary terms.
- [ ] Record the roundfix skill version and sync the mirrors.
- [ ] Re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The command reference and the user guide list the four codes and the
      adopter paragraph.
- [ ] `CONTEXT.md` defines the six terms.
- [ ] The roundfix skill records a raised version, and the mirrors match.
- [ ] The Behavior Surface Record is current.

## Consumers

- The spec command reference and user guide — consumers: `./internal/docscontract`, `./skills`

## Context

- instruction: `docs/adr/0265-a-spec-proves-what-it-claims-and-every-active-spec-is-held-to-it.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `for code in SC-ABSENCE-UNPROVEN SC-ADR-ALTERNATIVES-UNRELATED SC-COVERAGE-UNREAD SC-CONSUMERS-UNPROVEN; do grep -qF -- "$code" .agents/skills/roundfix/references/spec.md || { printf 'the command reference lacks %s\n' "$code" >&2; exit 1; }; grep -qF -- "$code" docs/user-guide/commands/spec.md || { printf 'the user guide lacks %s\n' "$code" >&2; exit 1; }; done; for file in .agents/skills/roundfix/references/spec.md docs/user-guide/commands/spec.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "These four gaps hold every active Spec, whatever its age." || { printf '%s lacks the adopter paragraph\n' "$file" >&2; exit 1; }; done; for file in references/spec.md SKILL.md; do cmp -s ".agents/skills/roundfix/$file" "skills/roundfix/$file" || { printf 'mirror differs: roundfix/%s\n' "$file" >&2; exit 1; }; done; for term in "**Absence Evidence**:" "**ADR Alternatives**:" "**Coverage Cross-Read**:" "**Consumer Declaration**:" "**Shared Contract**:" "**Integration Checkpoint**:"; do tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "$term" || { printf 'CONTEXT.md lacks %s\n' "$term" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract` — expected: exit 0; before this Task the command reference does not name SC-ABSENCE-UNPROVEN, so the command fails at its first check.

## References

- [_prd.md](_prd.md) — Goal 5; User Story 5; Core Features 5 and 6; Success Metric 3; Adopter notes; Glossary
- [_techspec.md](_techspec.md) — API Contract 5; Skill and document texts; Testing Approach; Build Order 4
- ADR-0265; ADR-0244; ADR-0256
