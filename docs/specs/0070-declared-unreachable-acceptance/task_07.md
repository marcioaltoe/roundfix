---
task: task_07
spec: 0070-declared-unreachable-acceptance
status: completed
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

## Result

Updated the Archive Command's operator contract in the user guide and the
canonical Roundfix Skill. Both now teach the two archive-eligible QA shapes:
`verdict: pass`, or a `partial` verdict whose only unmet rows are declared
unreachable and fully covered by the Spec's declarations. The declared-only
case records the declarations' `satisfied-by` actions under `unproven` so the
archive preserves what was never verified.

Both sources also preserve the refusal boundary explicitly: finding-blocked
rows, environment-blocked rows, an uncovered declared count, `verdict: fail`,
missing QA, and incomplete Tasks still refuse as applicable. `qa_override`
keeps its existing meaning for explicitly authorized archival of genuinely
failed or missing evidence. No Archive Command behavior or Go source changed.

Acceptance evidence:

- Pass-only sweep: `rtk grep -n "verdict: pass"
  docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md` found four
  operator-facing occurrences. Every occurrence names `partial` and the
  declared-unreachable case on the same line; the surrounding contract names
  full declaration coverage.
- Unproven and refusal contract: `rtk grep -n
  "unproven\\|finding-blocked\\|environment-blocked\\|qa_override"
  docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md` found the
  `unproven` record and unchanged refusal/override wording in both sources.
- Mirror identity: after `rtk make skills-sync`, `rtk cmp -s
  .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited 0.
- Repository Verification: `make verify` remains Daemon-owned and was not run
  in this Agent turn, so the full-gate criterion has no Agent-side terminal
  evidence.
- Behavior-source boundary: `rtk git diff --name-only HEAD` showed no Go file
  and no path under `internal/` outside `DERIVED_DIGEST_PATHS`. The only
  `internal/` changes are command-generated files under
  `internal/baseline/assets/setups/` and `internal/baseline/testdata/`.

Focused checks and required regeneration:

- `rtk make skills-sync` exited 0 and regenerated the authorized embedded
  Roundfix Skill mirror.
- `rtk make baseline-digests` exited 0, ran all six reported Go test steps,
  and regenerated five derived digest artifacts.
- The two exact corpus re-records from task_04 initially reached the managed
  sandbox's inaccessible configured Go build cache. With approved cache
  access, `rtk go test ./internal/baseline -count=1 -run
  TestBaselinePlanCharacterization -update-baseline-plan-characterization`
  passed 7 tests and `rtk go test ./internal/baseline -count=1 -run
  TestCatalogDiagnosticCharacterization -update-catalog-diagnostics` passed 2
  tests.
- `rtk git diff --check` exited 0 after the documentation, Skill, mirror, and
  generated-artifact changes.

No follow-up work was discovered inside this Task's slice. No command from
this Task's `## Verification` section was run; the Daemon owns terminal
Verification, Task status, and settlement.
