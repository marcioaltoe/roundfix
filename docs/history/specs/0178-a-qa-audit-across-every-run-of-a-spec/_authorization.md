---
status: approved
granted: 2026-09-28
action: audit every Task commit a Spec accumulated across its Runs against the grant at the Delivery Base, name each audited commit in the QA Report, measure the Daemon's staleness against the Delivery Base and refuse a stale gate, and record and check the QA Agent's user-flow binary
consuming: 0178-a-qa-audit-across-every-run-of-a-spec
paths:
  - .agents/skills/qa-gate/SKILL.md
  - skills/qa-gate/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0178

The maintainer authorized continuing with the next wave after the v0.18.0
release in chat on 2026-09-28 ("Após o release, pode continuar com as
implementações da onda seguinte"), and this Spec is part of that wave. The
skill files ride the standing grant of 2026-09-18 for keeping the shipped
skills true to the CLI; Go sources, were any of them governed, would ride the
standing grant of 2026-09-21 for governed source. The set was measured with
`GovernedPath` on 2026-09-28.

## Why each governed path is unavoidable

- `.agents/skills/qa-gate/SKILL.md`, `skills/qa-gate/SKILL.md` — the qa-gate
  skill states what the mechanical stage now audits and when the gate still
  audits Task commits by command (task_02), and that the auditor fields are
  Daemon-owned and the QA Agent records `user_flow_binary` (task_04);
  `.agents/skills/` is canonical and `skills/` its mirror.

## What is not governed

Measured with `GovernedPath`, these files are ordinary:

- `internal/daemon/task_engine.go`, `internal/daemon/engine.go`,
  `internal/daemon/task_engine_test.go` and the new
  `internal/daemon/qa_delivery_base.go`,
  `internal/daemon/qa_every_run_audit_test.go`,
  `internal/daemon/qa_auditor_staleness_test.go` and
  `internal/daemon/qa_user_flow_binary_test.go`;
- `internal/speccheck/mechanical.go`, `internal/speccheck/report.go` and the new
  `internal/speccheck/mechanical_commit_column_test.go`;
- `internal/spec/auditor_evidence.go`, `internal/spec/auditor_evidence_test.go`,
  `internal/spec/qa.go` and the new
  `internal/spec/auditor_delivery_base_test.go` and
  `internal/spec/qa_user_flow_binary_test.go`;
- `internal/app/version.go` and `internal/app/version_test.go`;
- `CONTEXT.md`.

`internal/speccheck/mechanical_test.go` is governed and no Task changes it; the
new speccheck tests live in their own file.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to `QAReportEligibility`, to archived Specs or to archived QA
  Reports, and no change to `docs/references/coverage-record.json`: no Task
  renames or replaces a top-level test.
- No paid API use, release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
