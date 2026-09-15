---
task: task_01
spec: 0138-a-matrix-the-spec-declares
status: pending
type: docs
complexity: medium
---

# Task 01: Rewrite the qa-gate matrix rules and regenerate the mirror

## Overview

The qa-gate skill has the gate build its own matrix from every promise and
exclusion. It also stops flow rows after any audit problem, and it lets a row's
inputs be written after the row ran. This slice rewrites those rules in the
canonical skill:

- a `qa` Task declares its coverage sources through Requirements that start with
  `MUST verify` or `MUST run`;
- every row names the sources it covers, every source appears in some row, and
  no row covers anything else;
- once the matrix exists, a finding blocks only the rows that depend on it.

The slice then regenerates the distributed mirror. The change can be checked on
its own: the anchor clauses are present, the removed sentences are absent, and
the skill's existing contract tests still pass.

## Requirements

1. MUST rewrite the canonical qa-gate skill as the TechSpec's "Where each Core
   Feature lands in the skill" describes. The rewrite covers sections 1 and 2,
   the Row input declaration, and section 3.
2. MUST carry each of these anchor clauses verbatim in the canonical skill:

   ```text
   Each numbered `qa` Task Requirement that starts with `MUST verify` or `MUST run` is a declared verification Requirement
   declaring a matrix does not make the run partial
   Every coverage source appears in at least one row's provenance
   no row's provenance names a source outside the coverage sources
   Once the matrix exists, a finding blocks only the rows that depend on it
   fixed when the row is planned as `pending`
   is not a proved environmental cause
   identical on the delivery target
   ```

3. MUST remove these two sentences from the canonical skill:

   ```text
   Any problem blocks flow QA after the complete command audit has been reported.
   The scope is complete when every promise and explicit exclusion in the spec maps to a planned row
   ```

4. MUST state the declaration and coverage rules exactly as the TechSpec gives
   them:
   - Every other Requirement constrains how the gate runs.
   - A declared Requirement that verifies every Acceptance Criterion of named
     Tasks makes each of those criteria its own source.
   - A `qa` Task without a declared verification Requirement gets the bounded
     default.
   - Rows may group several sources that one observation settles.
   - No row's result is computed from other rows.
   - No row re-checks a Mechanical Refusal Code.
5. MUST keep these coverage sources in every matrix, declared or default:
   - the outside-evidence row;
   - the Pull Request row, on its equivalent-evidence path;
   - the repository Verification;
   - each declaration under the PRD's `## Unreachable Acceptance` section, with
     the existing paragraph that reads those declarations kept;
   - the frontend sweep, when `frontend` is declared.
6. MUST keep the behavior of the tooling-audit bullet: the gate still audits
   Task commits by command. Only the first sentence in Requirement 3 leaves that
   bullet.
7. MUST keep verbatim every clause the skill's contract tests pin:
   - `roundfix spec check`
   - `Authoring rule removed from the QA matrix`
   - `did not run. It is not an equivalent`
   - `actual paths from the Daemon-owned Task commit`
   - `git diff-tree --no-commit-id --name-only -r`
   - `missing, late, or untraceable authorization`
   - `outside the exact bounded list`
   - `Task status or Task Graph dependencies`
8. MUST NOT change any of these, and the skill MUST keep `version: 0.0.2`:
   - the verdict rules in the closing section;
   - the three typed blocked causes;
   - QA Report keys;
   - the Results table columns;
   - report naming;
   - the Pull Request row's equivalent-evidence path;
   - the outside-evidence obligation;
   - the frontend sweep.
9. MUST regenerate the distributed mirror only through `make skills-sync`. MUST
   NOT hand-edit the mirror.
10. MUST change only the approved grant's bounded files and this Task file:
    - `.agents/skills/qa-gate/SKILL.md`
    - `skills/qa-gate/SKILL.md`
    - this Task file

    If a derived pin moves because of this edit, rewrite it only with
    `make baseline-digests`, which the grant sanctions. If any other test or
    file outside those paths fails, stop and record the failure in the Result
    instead of editing it.
11. MUST record a timeout or intermittent failure of the repository
    Verification as a failure, never an environment block. Contention, a
    non-reproducible run and the same failure on the unchanged delivery target
    MUST NOT count as proved environmental causes.
12. SHOULD ask provenance to use the `qa` Task's own identifiers, such as
    "Requirement 3", "Task 01 criterion 2" and "repository Verification".
13. SHOULD state that the mechanical stage's withholding of the Agent Session
    before a matrix exists is unchanged.

## Subtasks

- [ ] Rewrite section 1: the scope sentence, the partial clarification, and the
      removal of the blocking sentence.
- [ ] Rewrite section 2: the declaration rule, criterion expansion, bounded
      default, non-waivable sources and coverage rule.
- [ ] Rewrite the Row input declaration so inputs are fixed at planning.
- [ ] Rewrite section 3: timing failures and analyzer scope.
- [ ] Regenerate the mirror and confirm parity.

## Acceptance Criteria

- [ ] Every anchor clause appears in both the canonical skill and the mirror;
      today none does.
- [ ] Neither removed sentence appears in either copy; today both do.
- [ ] The canonical skill and the mirror are byte-identical, and the skill still
      declares `version: 0.0.2`.
- [ ] The skill's workflow contract tests still pass with their pinned clauses
      unchanged. These include the QA gate constraint and the
      tooling-authorization journey.
- [ ] No path outside the two bounded files, this Task file and any sanctioned
      regeneration output changed.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/baseline_skill_contract_test.go`

## Verification

- `for f in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md; do n="$(tr '\n' ' ' < "$f" | tr -s ' ')" || exit 1; for a in 'or MUST run is a declared verification Requirement' 'declaring a matrix does not make the run partial' 'Every coverage source appears in at least one row' 'names a source outside the coverage sources' 'the matrix exists, a finding blocks only the rows that depend on it' 'fixed when the row is planned as' 'is not a proved environmental cause' 'identical on the delivery target'; do printf '%s' "$n" | tr -d '\140' | grep -qF -- "$a" || { echo "$f is missing: $a"; exit 1; }; done; done` — every anchor clause is present in both copies; this fails today.
- `for f in .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md; do test -f "$f" || exit 1; n="$(tr '\n' ' ' < "$f" | tr -s ' ')" || exit 1; for a in 'Any problem blocks flow QA after the complete command audit has been reported.' 'The scope is complete when every promise and explicit exclusion in the spec maps to a planned row'; do if printf '%s' "$n" | grep -qF -- "$a"; then echo "$f still carries: $a"; exit 1; fi; done; done` — both removed sentences are gone from both copies; this fails today.
- `grep -qF 'is a declared verification Requirement' .agents/skills/qa-gate/SKILL.md || exit 1; grep -qx 'version: 0.0.2' .agents/skills/qa-gate/SKILL.md || exit 1; go test -count=1 ./skills -run '^(TestAuthorialSkillSync|TestProjectConstraintQAGate|TestToolingAuthorizationJourney|TestLegacySpecConstraintExemption)$'` — the mirror matches, the version stays and the pinned contract clauses survive.

## References

- `_prd.md` → User Stories 1-5; Core Features 1-8; Declared intentional breaks;
  Regression locks.
- `_techspec.md` → Implementation Design: Where each Core Feature lands in the
  skill, Anchor clauses; Integration Points; Build Order 1.
- ADR-0155.
