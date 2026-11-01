---
task: task_02
spec: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
status: pending
type: backend
complexity: medium
---

# Task 02: The authored-QA-Verification finding prints the expected block, and write-tasks teaches copying it

## Overview

`SC-QA-VERIFICATION-AUTHORED` gains the exact `## Verification` section the
QA Task must carry: the text report prints it after the `fix:` line as a
fenced `markdown` block, and JSON carries it in the finding's new optional
`expected` field. Copying it clears the finding. The `write-tasks` template
and `SKILL.md`, the Roundfix skill's spec reference and the spec command guide
teach the step. This answers the adopted reference of 2026-11-01, "No
authoring path shows the QA gate's derived Verification"
([references](references/2026-11-01-no-authoring-path-shows-the-qa-gates-derived-verification.md)),
the Quaestor report of 2026-11-01.

## Requirements

1. MUST add the field `Expected string` with the JSON tag
   `json:"expected,omitempty"` to `speccheck.Finding`
   (`internal/speccheck/report.go`), with a comment that it holds the exact
   text the fix asks for when the checker can derive it, and make
   `RenderText` write, right after the `  fix: ` line of a finding whose
   `Expected` is non-empty, `  expected block:`, a line with three backticks
   and `markdown`, the block, and a line with three backticks
   (`_techspec.md` → API Contract 5). A finding with an empty `Expected`
   renders the same bytes as before in text and JSON.
2. MUST add `ExpectedQAVerificationBlock(derived []string) string` to
   `internal/speccheck/coherence.go`, returning `## Verification`, a blank
   line, and one bullet per command holding it in one code span, each ending
   in a newline (API Contract 4). MUST make `AuthoredQAVerification` set
   `Expected` to the block of `task.Verification` and its fix to "Replace the
   qa Task's ## Verification section with the expected block this finding
   prints, unchanged."; its code, severity, summary and location keep their
   bytes, and it still fires for every differing QA Verification.
3. MUST add `TestAuthoredQAVerificationPrintsTheExpectedBlock` to
   `internal/speccheck/coherence_test.go` as `_techspec.md` → Testing
   Approach states: a QA Task with an authored command gives one finding whose
   `Expected` equals the block of `spec.DerivedQAVerification(slug)`, and the
   same Task file with its `## Verification` section replaced by that text
   gives no `SC-QA-VERIFICATION-AUTHORED`.
4. MUST add `TestAFindingRendersItsExpectedBlock` to
   `internal/speccheck/report_test.go`: `RenderText` of a finding with
   `Expected` followed by one without equals the exact expected text, and
   `RenderJSON` carries `expected` on the first and no `expected` key on the
   second. Keep every existing test of both files unchanged.
5. MUST document the finding in `docs/user-guide/commands/spec.md`, as a row
   of the code table placed before `SC-WAVE-COLLISION`, and in the canonical
   `.agents/skills/roundfix/references/spec.md`, as a bullet before the
   paragraph that begins "These four gaps hold every active Spec". Each says:
   a pending QA Task's `## Verification` renders a command other than the one
   Roundfix derives; an error from the Tasks stage; after its `fix:` line the
   text report prints `expected block:` and the exact section as a fenced
   `markdown` block; JSON carries it in the finding's `expected` field,
   present only on a finding that has one; replace the whole section with it,
   unchanged. Both contain the phrase "the finding's `expected` field".
6. MUST teach the step in the canonical `write-tasks` skill:
   - in `.agents/skills/write-tasks/references/task-template.md`, a comment
     right above the `- <command> — expected: ...` bullet of the task
     template's `## Verification` section, saying that a `type: qa` Task never
     authors its Verification, to write the bullet as a placeholder, run
     `roundfix spec check <slug> --stage tasks`, and replace the whole section
     with the block the finding prints under `expected block:` (JSON: the
     finding's `expected` field), unchanged;
   - in `.agents/skills/write-tasks/SKILL.md`, a new section
     `## QA Task Verification` placed right before `## Verification paths`,
     saying the same, that the next check reports no
     `SC-QA-VERIFICATION-AUTHORED`, that what the gate must prove belongs in
     its Requirements, and that `roundfix spec check --run-verification`
     reports the command `not probed` with reason `qa_gate`, so it needs no
     failure before the work.
   Both contain the phrase "expected block:". Change no other passage.
7. MUST copy both canonical skills with `make skills-sync`, then raise both
   version fields of `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/write-tasks/SKILL.md` and their mirrors with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
   -record-skill-versions -count=1`, recording
   `skills/testdata/owned-skill-versions.json`. task_01 raised the Roundfix
   version first; the record command raises above whatever is recorded. Never
   hand-edit a mirror or a version field.
8. MUST re-record the Behavior Surface Record as the LAST step, after the
   last edit of the guide, with `go test -count=1 -tags docscontract
   ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'
   -record-skill-coverage`, and commit `docs/references/behavior-surfaces.json`.
9. MUST NOT touch `spec.DerivedQAVerification`, `taskVerification`, the
   loader, `docs/references/skill-coverage.json`, `CONTEXT.md`, any
   `### QA settlement` section, any Baseline module, `docs/agents/`, the
   Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize: an authored QA Verification gives `SC-QA-VERIFICATION-AUTHORED` with no expected text today.
- [ ] Add the field, the block, the rendering and the two tests.
- [ ] Document the finding in the guide and the reference and teach the step in the `write-tasks` template and `SKILL.md`.
- [ ] Sync the mirrors, record both skill versions, re-record the fingerprints last.

## Acceptance Criteria

- [ ] `SC-QA-VERIFICATION-AUTHORED` carries the exact expected section in text and JSON, and copying it clears the finding.
- [ ] Every finding without an expected block renders the same bytes as before.
- [ ] The guide, both copies of the Roundfix spec reference and both copies of the `write-tasks` template and `SKILL.md` teach the step, the mirrors equal their canonical copies, and the owned-skill version and Behavior Surface records pass.

## Context

- instruction: `internal/spec/task.go`
- interface: `internal/speccheck/report.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/speccheck/coherence_test.go`
- interface: `internal/speccheck/report_test.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/write-tasks/references/task-template.md`
- interface: `skills/write-tasks/references/task-template.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `speccheck.Finding` JSON shape and `RenderText` — consumers: `./internal/speccheck`, `./internal/cli`
- The owned-skill version record and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAuthoredQAVerificationPrintsTheExpectedBlock|TestAFindingRendersItsExpectedBlock)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAuthoredQAVerificationPrintsTheExpectedBlock TestAFindingRendersItsExpectedBlock; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/speccheck && go test -count=1 -timeout 20m -run 'SpecCheck' ./internal/cli` — expected: exit 0. Before this Task neither test exists, so the pass lines are missing; afterwards both pass with the Spec Check suite and the CLI's Spec Check tests.
- `for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'SC-QA-VERIFICATION-AUTHORED' || { printf 'missing finding in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'expected block:' || { printf 'missing expected block in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md .agents/skills/write-tasks/references/task-template.md skills/write-tasks/references/task-template.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'expected block:' || { printf 'missing step in %s\n' "$f" >&2; exit 1; }; done; grep -qF '## QA Task Verification' .agents/skills/write-tasks/SKILL.md || { echo 'missing QA Task Verification section' >&2; exit 1; }; cmp .agents/skills/roundfix/references/spec.md skills/roundfix/references/spec.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md && cmp .agents/skills/write-tasks/references/task-template.md skills/write-tasks/references/task-template.md && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no document names the finding, so the first check fails. Afterwards the guide, both reference copies and both copies of the `write-tasks` files teach the expected block, the mirrors equal their canonical copies, and the owned-skill version and documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Features 3 and 4; Success Metric 4; Success Metric 5; Success Metric 6
- `_techspec.md` → Current behavior 5-7; API Contracts 4-6; Invariants 4 and 5; Declared breaks; Testing Approach; Build Order 2
