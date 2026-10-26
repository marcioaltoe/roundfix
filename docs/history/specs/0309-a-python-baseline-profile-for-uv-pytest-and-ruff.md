---
schema: roundfix/archive-record/v1
spec: 0309-a-python-baseline-profile-for-uv-pytest-and-ruff
title: A Python Baseline Profile for uv, pytest and ruff
status: archived
created: "2026-10-23"
archived: "2026-10-26"
disposition: qa-override
source: docs/specs/0309-a-python-baseline-profile-for-uv-pytest-and-ruff
source_revision: 06a9b7ce56e18021516c90f14966558445336387
qa_task: task_04
qa_report: qa-report-2026-10-26.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'Marcio, AskUserQuestion 2026-10-26 ~07:55 BRT: "Autorizo o override (Recomendado)"'
qa_override_reason: 'QA partial with no failed, skipped or finding-blocked row; R07 (Quaestor adoption report) was recorded as ''another repository prohibited by Spec authorization'', outside the accepted wording; the operator verified it read-only through the GitHub API: marcioaltoe/quaestor at d17c98a has 47 skills in skills-lock.json and 47 directories in .agents/skills, handoff included, and docs/agents/specific-repository.md. R06/R08 are the accepted network-prohibited rows, R12 the Pull Request row, R09 the declared Unreachable Acceptance.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 06a9b7ce56e18021516c90f14966558445336387
adrs:
  - ADR-0322
sources:
  - 2026-10-22-a-baseline-profile-for-python-repositories.md
regeneration:
  - command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
    outputs:
      - internal/baseline/assets/setups/python.json
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/core.json
      - internal/baseline/assets/modules/python.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted:
  - docs/references/2026-10-22-a-baseline-profile-for-python-repositories.md
relinked:
  - docs/specs/0310-the-python-profile-writes-its-starter-files/_prd.md
  - docs/specs/0310-the-python-profile-writes-its-starter-files/task_01.md
  - docs/specs/0310-the-python-profile-writes-its-starter-files/task_02.md
  - docs/specs/0310-the-python-profile-writes-its-starter-files/task_03.md
---

# A Python Baseline Profile for uv, pytest and ruff

The implemented Python profile plans, applies and re-plans without drift in a disposable adopter, all nine sabotage cases detect their intended defects, and the protected-file and reproducible-generation audits find no issue. Real adoption remains unproven under the single PRD declaration, while the upstream file, coordinator report, published documentation and pre-PR acceptance are not freshly observed. The offline archive plan lists the Spec cut, its adopted Backlog Entry as a candidate promotion, and the citations in Spec 0310 that an eventual archive would relink; ADR-0322 and the user guide already hold the durable decisions. QA performs no promotion or archive mutation, and this partial report leaves the missing operator evidence explicit for the daemon.
