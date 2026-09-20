---
task: task_05
spec: 0151-a-supersession-the-archive-can-see
status: pending
type: backend
complexity: medium
---

# Task 05: A deliverer check that checks what it claims

## Overview

Pre-PR review found two gaps, both verified before this Task was written.

**The existence check does not check existence as documented.**
`knownSpecDirectory` stats `_prd.md` and returns true when it is not a
directory. Nothing reads its status. So a `--by` slug naming a draft, an
already-archived-in-place, or a malformed Spec in the active root is accepted,
and the command records a deliverer that is neither active nor archived —
exactly the condition its own refusal claims to catch. A record naming a Spec
that does not really exist is worse than no record, because it reads as an
answer.

**The glossary was never checked.** `docs/agents/domain.md` makes the check
mandatory at the close of work that introduces a term, and this Spec introduces
supersession as a durable lifecycle concept with a public command. `CONTEXT.md`
has no mention of it. The contract has to outlive this Spec's folder.

## Requirements

1. MUST accept a `--by` slug only when the named Spec is genuinely active in the
   configured Spec Root or present in the archive, judged by reading its
   declared status rather than by the presence of a file.
2. MUST refuse a `--by` Spec whose status is draft, archived-in-place, absent or
   malformed, naming that condition, and MUST write nothing.
3. MUST define supersession in `CONTEXT.md`, including its relationship to
   archive: that it is the proof archive accepts for a Spec whose content
   another Spec delivered.
4. MUST keep every behavior Tasks 01 through 03 delivered.

## Subtasks

- [ ] Read the declared status in the deliverer check.
- [ ] Add a test per rejected status.
- [ ] Add the glossary term and its archive relationship.

## Acceptance Criteria

- [ ] A `--by` Spec whose `_prd.md` declares a non-active status is refused.
- [ ] A `--by` Spec whose `_prd.md` is malformed is refused, not accepted.
- [ ] An active `--by` Spec and an archived one are both still accepted.
- [ ] `CONTEXT.md` defines supersession and names its archive relationship.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/supersede.go`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run "^TestSupersedeRejectsANonActiveDeliverer" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestSupersedeRejectsANonActiveDeliverer"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `grep -qi "supersession" CONTEXT.md` — expected: exit 0; the glossary carries the term. Before this Task it does not, so the command fails.

## References

- [_techspec.md](_techspec.md) — The command
