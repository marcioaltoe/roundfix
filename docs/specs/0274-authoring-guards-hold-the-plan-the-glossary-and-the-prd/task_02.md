---
task: task_02
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
status: completed
type: docs
complexity: low
---

# Task 02: The QA secret refusal is documented as the Daemon prints it

## Overview

The Daemon settles a QA Task that carries a secret with
`QA evidence carries a secret: <path>:<line>: <kind>`, while the Roundfix
`implement` reference, its mirror and the `implement` user guide write
`<path>:<line> (<kind>)` (`_techspec.md` → Current behavior). This Task
corrects the three documents and adds a Repository Contract Test that derives
the form from the Daemon's source and the Secret Match format, so the two can
no longer drift apart. It is verifiable on its own by that test and the
mirror, version and Behavior Surface Record checks.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-14,
   [the implement reference states the QA secret refusal in another form](references/2026-10-14-the-implement-reference-states-the-qa-secret-refusal-in-another-form.md),
   by writing `QA evidence carries a secret: <path>:<line>: <kind>` in
   `.agents/skills/roundfix/references/implement.md` and in
   `docs/user-guide/commands/implement.md`, keeping the rest of each sentence,
   and running `make skills-sync` for the mirror. No text inside any
   `### QA settlement` section may change.
2. MUST add, to the existing `internal/docscontract/command_documentation_test.go`
   (`docscontract` build tag, `//verify:always`), and to no new file,
   implementing API Contract 2 of `_techspec.md` in
   `TestTheQASecretRefusalIsDocumentedAsTheDaemonPrintsIt`, and
   `TestTheQASecretRefusalContractRejectsTheParenthesizedForm`, which runs the
   same comparison on the old sentence and requires a failure naming it. Both
   call `t.Parallel()` first. The documented form MUST be computed from the
   Daemon literal and `secretscan.Match.String()`, never copied as a literal,
   and no production symbol may be added for the test.
3. MUST NOT change the Daemon's refusal text, `secretscan.Match.String()`,
   `TestQAReportWithASecretMatchSettlesFailedWithoutACommit` or
   `CHANGELOG.md`.
4. MUST raise and record the Roundfix skill version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`,
   because this Task edits a user guide page. No record or mirror is
   hand-edited.

## Subtasks

- [ ] Correct the reference and the user guide, and sync the mirror.
- [ ] Add the contract test and its negative case.
- [ ] Record the skill version and the Behavior Surface Record.

## Acceptance Criteria

- [ ] The three documents state the colon form and none states the parenthesized one.
- [ ] The contract test passes on the delivered tree and its negative case fails on the old sentence.
- [ ] The mirror, the skill version and the Behavior Surface Record are current.

## Consumers

None. This Task changes documentation and adds a contract test; no code consumes either, and the Daemon's message is unchanged.

## Context

- instruction: `docs/adr/0273-no-secret-and-no-spec-content-leaves-the-machine-from-a-run.md`
- instruction: `internal/daemon/task_engine.go`
- instruction: `internal/secretscan/scan.go`
- interface: `internal/docscontract/command_documentation_test.go`
- instruction: `internal/verifyselect/contracts_repository_test.go`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `docs/user-guide/commands/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run '^(TestTheQASecretRefusalIsDocumentedAsTheDaemonPrintsIt|TestTheQASecretRefusalContractRejectsTheParenthesizedForm)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheQASecretRefusalIsDocumentedAsTheDaemonPrintsIt TestTheQASecretRefusalContractRejectsTheParenthesizedForm; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md docs/user-guide/commands/implement.md; do if tr -s '[:space:]' ' ' < "$f" | grep -qF -- '<line> (<kind>)'; then printf 'parenthesized refusal still in %s\n' "$f" >&2; exit 1; fi; done; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1; go test -count=1 -run '^TestQAReportWithASecretMatchSettlesFailedWithoutACommit$' ./internal/daemon || exit 1; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract` — expected: exit 0. Before this Task the contract tests do not exist and the documents carry the parenthesized form, so the command fails at its first pass check.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 2; Core Feature 5; Success Metric 2
- `_techspec.md` → API Contract 2; API Contract 5; Invariant 4; Current behavior; Build Order 2
- ADR-0273; ADR-0280

## Result

Implemented the assigned documentation and contract-test slice.

- The canonical Roundfix `implement` reference and the user guide now state `QA evidence carries a secret: <path>:<line>: <kind>`. `make skills-sync` refreshed the mirror, and canonical/mirror comparisons matched.
- `internal/docscontract/command_documentation_test.go` now derives the documented form from the Daemon string literal and `secretscan.Match.String()`, checks every relevant Markdown file, and includes the negative parenthesized-form case. Both tests call `t.Parallel()` first.
- `go test -count=1 -tags docscontract ./internal/docscontract -run 'TestTheQASecretRefusal'`: passed after the mirror sync; this focused run exercised both new contract tests.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed and recorded Roundfix skill version `0.1.84`.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed and refreshed the `implement` guide Behavior Surface fingerprint.
- The Daemon refusal source, `secretscan.Match.String()`, the existing QA secret settlement test, `CHANGELOG.md`, and all `### QA settlement` sections were left unchanged.

The Daemon-owned declared Verification and Task status were not changed or run in this handoff.

## Carry-forward provenance

- Source Run: `run_20261017T095043Z_d5b2f5791da3cba6`
- Source commit: `26be1a00e117a1580bbcb07df2b8d44e527780a2`
