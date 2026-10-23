---
schema: roundfix/archive-record/v1
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
title: A subsumed test selector is flagged and a QA refusal names its commit
status: archived
created: "2026-10-22"
archived: "2026-10-23"
disposition: partial
source: docs/specs/0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
source_revision: 46de47be57b4bad5375991af8799c6e5898d88f0
qa_task: task_05
qa_report: qa-report-2026-10-23.md
qa_verdict: partial
unproven:
  - Quaestor's next plan check after the operator's release message to that project.
adrs:
  - ADR-0315
sources:
  - 2026-10-22-a-pytest-node-id-beside-its-own-file-passes-before-the-test-exists.md
  - 2026-10-22-a-qa-precondition-refusal-report-records-no-audited-commit.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "645"
---

# A subsumed test selector is flagged and a QA refusal names its commit

The Spec now reports a pending pytest node id that its invocation also collects through a containing argument, names that argument and the separate-command fix, and records the audited commit in final Daemon QA reports while preserving the pending seed and report readers; the declared matrix passed every runnable row, including the public CLI repair flow, all named tests, thirteen isolated sabotages and the scope audit.
