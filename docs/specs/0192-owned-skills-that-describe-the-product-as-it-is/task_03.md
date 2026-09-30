---
task: task_03
spec: 0192-owned-skills-that-describe-the-product-as-it-is
status: pending
type: docs
complexity: high
---

# Task 03: The authoring skills teach a Spec the Delivery Queue accepts

## Overview

`write-prd`, `write-techspec` and `write-tasks` are what an Agent follows to author a Spec. They treat the authorization record as a tooling-only concern, although `roundfix deliver start` refuses any Spec whose record does not grant the five delivery operations. They also omit the `## Unreachable Acceptance` section, the coverage units the checker enforces, the delivery order, and three lessons the last waves paid for. This Task corrects the three skills and their templates. The readers are authoring Agents in this repository and in every repository that adopts the Baseline, so each sentence describes the Roundfix contract, not this repository.

## Requirements

1. MUST state in `.agents/skills/write-prd/SKILL.md` and `.agents/skills/write-techspec/SKILL.md`, next to the existing authorization paragraph, that every Spec a Delivery Queue delivers carries an approved `_authorization.md`, with or without protected tooling. The text names the record's `operations` list with `implement`, `commit`, `push`, `pull_request` and `merge`, and says a Spec that changes no Governed Path records `paths: []`. It MUST contain the words `with or without protected tooling`.
2. MUST keep the literal `docs/agents/backend.md` in both skills and add a sentence that starts `A repository without that guide` and says such a repository cites the guide that owns the policy for its surfaces.
3. MUST align the question example under "Clarify" in `write-prd` with the structured question form in `_techspec.md`: one question, its three options with the recommended one first and labelled `(Recommended)`, the consequence of each option, and no `D) Other — describe` line. The paragraph below it says to ask through the structured question tool the session exposes, and that without one the same single question is shown in chat and a custom answer is accepted.
4. MUST edit `.agents/skills/write-prd/references/prd-template.md`:
   - add an optional `## Unreachable Acceptance` section after `## Success Metrics`, whose comment shows one declaration with the three fields `criterion`, `reason` and `satisfied-by`, and says to omit the section when every criterion is reachable;
   - extend the Tooling authority comment to say that a Spec a Delivery Queue delivers also records its delivery operations in that record and that, with no Governed Path to bound, the record carries `paths: []`. The comment MUST contain the words `with no Governed Path to bound`.
5. MUST edit `.agents/skills/write-techspec/references/techspec-template.md`:
   - the Coverage Map comment asks for one line per PRD goal, user story, Core Feature and Success Metric;
   - the Tooling authority comment gains the same sentence as the PRD template.
6. MUST keep, in both templates, exactly four `<applicable | not applicable>` placeholders and exactly four ``Source: `docs/agents/`` placeholders, and the phrases `express maintainer authorization`, `bounded files` and `no protected tooling mutation`. No template may gain authorization front matter.
7. MUST edit `.agents/skills/write-tasks/SKILL.md`:
   - replace the sentence that starts `Follow one order per Spec:` with the order: implement the graph including its authored gate, run the configured pre-PR review, archive on the branch, pass the repository gate, open the Pull Request, verify current-head checks, and merge. The sentence MUST contain the words `run the configured pre-PR review, archive on the branch`;
   - replace the sentence that says the checker does not yet verify what a cited ADR says with one that names `SC-CITATION-UNSUPPORTED` as the check that reports a claim the ADR's text does not support;
   - add to the Corrective-Task text that a corrective Task added after the gate settled `completed` becomes a dependency of the gate, and that the author must reopen the settled gate first with `roundfix reopen --spec <slug>`, never by editing the QA Task file;
   - add to the Project Constraint preflight that the record's `operations` list must grant every delivery operation the Spec will use, that `roundfix deliver start` refuses a record lacking `implement`, `commit`, `push`, `pull_request` or `merge`, and that a Spec with no Governed Path records `paths: []`;
   - add to the Verification rule that code for another operating system is verified by building its non-test code for that system, with an example such as `GOOS=windows go build -buildvcs=false ./<package>`, and never with `go vet`, which also compiles tests written for the host.
8. MUST change the `status` comment in `.agents/skills/write-tasks/references/task-template.md` to say that the Daemon writes this during a Run and that `implement-task` writes it only in standalone execution.
9. MUST move both version fields of `write-prd` and `write-techspec` to `0.0.3` and of `write-tasks` to `0.0.4`.
10. MUST regenerate the mirrors with `make skills-sync`, then run `make baseline-digests`, and name in the Result every path either command rewrote.
11. MUST write each wording this Task quotes literally, because the Verification checks it on whitespace-folded text; line wrapping is free, other rewording is not.
12. MUST keep every phrase the existing contract tests require from these skills and templates, and MUST NOT edit any test or any other skill. The tests this Task's Verification names in `./skills` carry the required phrases.

## Subtasks

- [ ] Teach the full authorization record in `write-prd` and `write-techspec` and their templates.
- [ ] Add the backend-guide fallback sentence to both skills.
- [ ] Align the `write-prd` question example with the structured question form.
- [ ] Add `## Unreachable Acceptance` to the PRD template and correct the Coverage Map comment.
- [ ] Correct `write-tasks`: order, citation check, reopen, operations, cross-OS verification.
- [ ] Correct the status comment in the task template.
- [ ] Raise the three versions and regenerate the mirrors and digests.

## Acceptance Criteria

- [ ] `write-prd` and `write-techspec` say every delivered Spec carries an approved authorization record with its operations, and name `paths: []`.
- [ ] The PRD template carries an optional `## Unreachable Acceptance` section with `criterion`, `reason` and `satisfied-by`.
- [ ] The TechSpec template asks for a Coverage Map line per goal, user story, Core Feature and Success Metric.
- [ ] `write-tasks` states the order with the pre-PR review before the archive, names `SC-CITATION-UNSUPPORTED`, and says when `roundfix reopen` is needed.
- [ ] `write-tasks` says foreign-OS code is verified by building non-test code, never with `go vet`.
- [ ] The task template says the Daemon writes `status` during a Run.
- [ ] No owned authoring skill shows a `D) Other` option.
- [ ] Every contract test that reads these skills passes unedited, and each mirror is byte-identical to its canonical skill.

## Context

- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`
- interface: `skills/write-techspec/references/techspec-template.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/write-tasks/references/task-template.md`
- instruction: `docs/user-guide/commands.md`

## Verification

- `for pair in ".agents/skills/write-prd/SKILL.md|with or without protected tooling" ".agents/skills/write-techspec/SKILL.md|with or without protected tooling" ".agents/skills/write-prd/SKILL.md|A repository without that guide" ".agents/skills/write-techspec/SKILL.md|A repository without that guide" ".agents/skills/write-prd/SKILL.md|(Recommended)" ".agents/skills/write-prd/references/prd-template.md|## Unreachable Acceptance" ".agents/skills/write-prd/references/prd-template.md|satisfied-by" ".agents/skills/write-prd/references/prd-template.md|with no Governed Path to bound" ".agents/skills/write-techspec/references/techspec-template.md|with no Governed Path to bound" ".agents/skills/write-techspec/references/techspec-template.md|user story, Core Feature and Success Metric" ".agents/skills/write-tasks/SKILL.md|pre-PR review, archive on the branch" ".agents/skills/write-tasks/SKILL.md|SC-CITATION-UNSUPPORTED" ".agents/skills/write-tasks/SKILL.md|reopen the settled gate first" ".agents/skills/write-tasks/SKILL.md|must grant every delivery operation" ".agents/skills/write-tasks/SKILL.md|building its non-test code for that system" ".agents/skills/write-tasks/references/task-template.md|the Daemon writes this during a Run"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; for pair in ".agents/skills/write-prd/SKILL.md|D) Other" ".agents/skills/write-techspec/references/techspec-template.md|One line per PRD goal and user story" ".agents/skills/write-tasks/SKILL.md|archive, open the Pull Request, watch until Clean" ".agents/skills/write-tasks/SKILL.md|not yet verify that a cited ADR" ".agents/skills/write-tasks/references/task-template.md|only implement-task changes this"; do file="${pair%%|*}"; phrase="${pair#*|}"; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase"; then printf 'stale phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; fi; done; grep -q "^version: 0.0.3$" .agents/skills/write-prd/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/write-prd/SKILL.md || { printf 'version fields are not 0.0.3 in %s\n' .agents/skills/write-prd/SKILL.md >&2; exit 1; }; grep -q "^version: 0.0.3$" .agents/skills/write-techspec/SKILL.md && grep -q "^  version: 0.0.3$" .agents/skills/write-techspec/SKILL.md || { printf 'version fields are not 0.0.3 in %s\n' .agents/skills/write-techspec/SKILL.md >&2; exit 1; }; grep -q "^version: 0.0.4$" .agents/skills/write-tasks/SKILL.md && grep -q "^  version: 0.0.4$" .agents/skills/write-tasks/SKILL.md || { printf 'version fields are not 0.0.4 in %s\n' .agents/skills/write-tasks/SKILL.md >&2; exit 1; }; diff -r .agents/skills/write-prd skills/write-prd >/dev/null || { printf 'mirror differs: %s\n' skills/write-prd >&2; exit 1; }; diff -r .agents/skills/write-techspec skills/write-techspec >/dev/null || { printf 'mirror differs: %s\n' skills/write-techspec >&2; exit 1; }; diff -r .agents/skills/write-tasks skills/write-tasks >/dev/null || { printf 'mirror differs: %s\n' skills/write-tasks >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestWritePRDProjectConstraints|TestWriteTechSpecProjectConstraints|TestAuthoringTemplatesUseTheBoundedFilesLabel|TestTaskTemplateStatesTheStatusPreservingVerificationForm|TestTaskTemplateStatesTheWrapTolerantPhraseForm|TestWriteTasksSkillStatesTheDeclaredPathRules|TestSpecReferenceLifecycleSkillContracts|TestProjectConstraintTaskGate|TestProjectConstraintPRDGate|TestProjectConstraintTechSpecGate|TestAuthoringConstraintOwnership|TestAuthorialSkillSync)$" ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestWritePRDProjectConstraints TestWriteTechSpecProjectConstraints TestAuthoringTemplatesUseTheBoundedFilesLabel TestTaskTemplateStatesTheStatusPreservingVerificationForm TestTaskTemplateStatesTheWrapTolerantPhraseForm TestWriteTasksSkillStatesTheDeclaredPathRules TestSpecReferenceLifecycleSkillContracts TestProjectConstraintTaskGate TestProjectConstraintPRDGate TestProjectConstraintTechSpecGate TestAuthoringConstraintOwnership TestAuthorialSkillSync; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -tags docscontract -v -run "^(TestProjectConstraintDocumentation|TestProfilesDocumentationContractMatchesPublicGuidance)$" ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestProjectConstraintDocumentation TestProfilesDocumentationContractMatchesPublicGuidance; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task the first required phrase is absent from `write-prd`, so the command fails.

## References

- [_techspec.md](_techspec.md) — The authoring skills; The structured question form
- `_prd.md` → Goal 1; Core Feature 3; Core Feature 5; Core Feature 7; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Testing Approach 3; Testing Approach 4; Build Order 3
