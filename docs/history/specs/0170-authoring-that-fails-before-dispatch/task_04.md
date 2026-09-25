---
task: task_04
spec: 0170-authoring-that-fails-before-dispatch
status: completed
type: docs
complexity: low
---

# Task 04: Make the authoring templates and skills produce what the checker accepts

## Overview

A template-faithful PRD fails its first `spec check` because the PRD and
TechSpec templates ask for `bounded paths:`, which `SC-TOOLING-UNBOUNDED`
refuses. The Task template and the write-tasks skill do not tell an author which
Verification form keeps a tool's status, which Context entries the checker
audits for authority, or that a CLI change names its guide. This Task corrects
the templates and the skill and pins the label with a contract test.

## Requirements

1. MUST replace `bounded paths:` with `bounded files:` in the Tooling authority
   row of `.agents/skills/write-prd/references/prd-template.md` and
   `.agents/skills/write-techspec/references/techspec-template.md`, so neither
   template contains `bounded paths:`.
2. MUST state in the Verification comment of
   `.agents/skills/write-tasks/references/task-template.md` that a tool piped
   into `grep` inside a command substitution hides the tool's status and is
   refused with `SC-VERIFY-INVERTED-EXIT`, and give the status-preserving form
   `out="$(tool 2>&1)" || exit 1; ! printf '%s\n' "$out" | grep -q pattern`.
3. MUST state in `.agents/skills/write-tasks/SKILL.md` that every path a Task
   edits is declared under `interface:` or `creates:` and never `instruction:`,
   that each declared or Verification-read Governed Path must be in the Spec's
   `_authorization.md` `paths:` and both `bounded files:` rows or is refused with
   `SC-TOOLING-UNDECLARED`, and that a Task naming a CLI surface names its skill
   or guide itself or through a Task it depends on or is reported with
   `SC-CLI-UNDOCUMENTED`.
4. MUST add to `skills/baseline_skill_contract_test.go` a test that the PRD and
   TechSpec templates each contain `bounded files:` in their Tooling authority
   row and do not contain `bounded paths:`, a test that the Task template carries
   the status-preserving form, and a test that the write-tasks skill names
   `SC-TOOLING-UNDECLARED`, `SC-CLI-UNDOCUMENTED` and the `instruction:` rule.
5. MUST regenerate the mirrors under `skills/` with `make skills-sync`.

## Subtasks

- [ ] Correct the Tooling authority label in both templates.
- [ ] Add the Verification note to the Task template.
- [ ] Add the declared-path and CLI-guide rules to the write-tasks skill.
- [ ] Pin the three contracts with tests and run `make skills-sync`.

## Acceptance Criteria

- [ ] Neither template, canonical or mirror, contains `bounded paths:`.
- [ ] The Task template and the write-tasks skill carry the new rules.
- [ ] The three contract tests pass and the mirrors match.

## Context

- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`
- interface: `skills/write-techspec/references/techspec-template.md`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/write-tasks/references/task-template.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/baseline_skill_contract_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAuthoringTemplatesUseTheBoundedFilesLabel|TestTaskTemplateStatesTheStatusPreservingVerificationForm|TestWriteTasksSkillStatesTheDeclaredPathRules)$" ./skills 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestAuthoringTemplatesUseTheBoundedFilesLabel TestTaskTemplateStatesTheStatusPreservingVerificationForm TestWriteTasksSkillStatesTheDeclaredPathRules; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done && ! grep -q "bounded paths:" .agents/skills/write-prd/references/prd-template.md && ! grep -q "bounded paths:" .agents/skills/write-techspec/references/techspec-template.md && ! grep -q "bounded paths:" skills/write-prd/references/prd-template.md && ! grep -q "bounded paths:" skills/write-techspec/references/techspec-template.md && grep -q "SC-TOOLING-UNDECLARED" .agents/skills/write-tasks/SKILL.md && grep -q "SC-CLI-UNDOCUMENTED" .agents/skills/write-tasks/SKILL.md && diff -r .agents/skills/write-prd skills/write-prd >/dev/null && diff -r .agents/skills/write-techspec skills/write-techspec >/dev/null && diff -r .agents/skills/write-tasks skills/write-tasks >/dev/null` — expected: exit 0; before this Task the three tests do not exist and both templates still say `bounded paths:`, so the command fails.

## References

- `_prd.md` → Goal 3; Core Feature 4; Success Metric 6.
- `_techspec.md` → Templates and skills; Testing Approach 4; ADR-0131.

## Result

Implementation evidence:

- Updated the canonical PRD and TechSpec Tooling authority rows to use
  `bounded files:`; `make skills-sync` regenerated their mirrors.
- Added the status-preserving Verification guidance, including
  `SC-VERIFY-INVERTED-EXIT` and the `out="$(tool 2>&1)" || exit 1; ! printf
  '%s\n' "$out" | grep -q pattern` form, to the canonical Task template and
  regenerated its mirror.
- Added the explicit `interface:`/`creates:` declaration rule, authorization
  and `SC-TOOLING-UNDECLARED` rule, and CLI-guide/`SC-CLI-UNDOCUMENTED` rule to
  the canonical write-tasks skill and regenerated its mirror.
- Added the three named contract tests in
  `skills/baseline_skill_contract_test.go`.

Focused checks:

- `gofmt -w skills/baseline_skill_contract_test.go` — passed.
- The three named contract tests, run individually with a task-scoped
  `GOCACHE`, passed.
- `make skills-sync` — passed; Git emitted a non-blocking fsmonitor IPC
  diagnostic while reporting status.
- `go test -count=1 ./skills -run
  'TestAuthoringTemplatesUseTheBoundedFilesLabel|TestTaskTemplateStatesTheStatusPreservingVerificationForm|TestWriteTasksSkillStatesTheDeclaredPathRules|TestAuthorialSkillSync'`
  with task-scoped `GOCACHE` — passed.
- `git diff --check` — passed.

Acceptance evidence:

- Canonical and mirror template labels are covered by
  `TestAuthoringTemplatesUseTheBoundedFilesLabel`; the canonical/mirror skill
  trees are covered by `TestAuthorialSkillSync`.
- The Task-template wording and write-tasks rules are covered by
  `TestTaskTemplateStatesTheStatusPreservingVerificationForm` and
  `TestWriteTasksSkillStatesTheDeclaredPathRules`.
- All three contract tests and the post-sync mirror contract passed. The
  authored Verification command remains for the Daemon to run.
