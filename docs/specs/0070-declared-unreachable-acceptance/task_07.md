---
task: task_07
spec: 0070-declared-unreachable-acceptance
status: pending
type: docs
complexity: medium
---

# Task 07: Teach the Archive Command's real contract

## Overview

QA finding F-001: the Archive Command's operator-facing contract still teaches
the pass-only rule this Spec widened. `docs/user-guide/commands.md` and the
canonical Roundfix Skill describe a command that refuses what the shipped
binary now accepts. The built help and the real archive journey are already
correct — only the documentation lies, which is the worst shape for a
contradiction to take, because a reader trusts it.

The skill half is not optional polish. The repository's HARD RULE makes it a
precondition: a Pull Request that changes CLI behaviour ships the skill update
too, and this Spec changes CLI behaviour.

## Requirements

1. MUST state, everywhere the Archive Command's contract appears, that archive
   accepts either `verdict: pass` or a `partial` verdict whose only unmet rows
   are declared unreachable and fully covered by the Spec's declarations.
2. MUST state that the archive record carries those declarations' satisfying
   actions under `unproven`, so a reader of the archive learns what was never
   verified.
3. MUST state that every other refusal is unchanged — a finding-blocked row, an
   environment-blocked row, an uncovered declared count, and `fail` all still
   refuse — and that `qa_override` keeps its existing meaning.
4. MUST update `.agents/skills/roundfix/SKILL.md` at every place the pass-only
   rule appears, and regenerate its `skills/roundfix/**` mirror with
   `make skills-sync`.
5. MUST run `make baseline-digests` after the skill edit and re-record the two
   characterization corpora that command does not reach, exactly as task_04
   documents them, because a skill edit moves the same digests.
6. MUST change only `docs/user-guide/commands.md`, `.agents/skills/roundfix/**`,
   `skills/roundfix/**`, this Task file, and the ADR-0081 digest fallout under
   `DERIVED_DIGEST_PATHS`. The skill paths are authorized by the 2026-08-04
   close addendum that names Spec 0070; any other path is out of scope — stop
   rather than widen it.
7. MUST NOT change the Archive Command's behaviour. This Task corrects
   documentation to match the binary, never the reverse.

## Subtasks

- [ ] Update the Archive Command contract in the user guide.
- [ ] Update every pass-only statement in the canonical Roundfix Skill.
- [ ] Run `make skills-sync`, then `make baseline-digests` and both corpus
      re-records.
- [ ] Confirm no Go source changed.

## Acceptance Criteria

- [ ] No operator-facing source states that archive requires `verdict: pass`
      without naming the declared-unreachable case, proven by a search across
      the user guide and the skill rather than by inspecting known lines.
- [ ] Both sources describe the `unproven` record and state that every other
      refusal is unchanged.
- [ ] `skills/roundfix/` is byte-identical to `.agents/skills/roundfix/`.
- [ ] `make verify` exits 0 after the regeneration chain.
- [ ] No file under `internal/` outside `DERIVED_DIGEST_PATHS` changed, proving
      the binary's behaviour was not touched.

## Context

- instruction: `.agents/skills/roundfix/SKILL.md`
- instruction: `docs/user-guide/commands.md`
- instruction: `docs/agents/skill-dispatch.md`

## Verification

- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `go run -buildvcs=false ./cmd/roundfix skills check` — expected: exit 0.
- `if grep -rn "verdict: pass" docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md | grep -iv "declared\|partial\|unproven" | grep -q .; then exit 1; fi`
  — expected: exit 0; no pass-only statement survives without the declared case
  beside it.
- `grep -q "unproven" .agents/skills/roundfix/SKILL.md` — expected: exit 0.
- `grep -q "unproven" docs/user-guide/commands.md` — expected: exit 0.
- `make verify` — expected: exit 0.
- `git diff --name-only HEAD | grep "^internal/" | grep -vE "^internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no behaviour source changed.
- `git diff --name-only HEAD | grep -vE "^(docs/user-guide/commands\.md$|\.agents/skills/roundfix/|skills/roundfix/|docs/specs/0070-declared-unreachable-acceptance/task_07\.md$|internal/baseline/(assets/(setups|source-baselines|formatter-fixtures|profiles)|testdata)/)" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; only the bounded paths changed.

## References

- `_prd.md` → Core Features 3, 4 and 5; Decisions.
- `_techspec.md` → API Contracts.
- `qa/qa-report-2026-08-04.md` → F-001 and its required repair.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md` → Spec 0070
  close addendum.
- ADR-0080, ADR-0081.
