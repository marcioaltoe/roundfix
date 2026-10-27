---
schema: roundfix/archive-record/v1
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
title: Runs are not lost to an authorization row or a provider at capacity
status: archived
created: "2026-10-26"
archived: "2026-10-27"
disposition: partial
source: docs/specs/0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
source_revision: 4c4a4a60f6eb8f6bbfef62575c17f8d76a645331
qa_task: task_05
qa_report: qa-report-2026-10-27.md
qa_verdict: partial
unproven:
  - the operator's queue log and the Run Event Journal of the next Run that meets the reply (event `provider_system_error` with recovery `fallback` in a QA scope or after a First Handoff), recorded session evidence.
adrs:
  - ADR-0335
  - ADR-0274
  - ADR-0276
sources:
  - 2026-10-26-an-outside-evidence-row-that-names-what-the-authorization-forbids-does-not-qualify.md
  - 2026-10-26-a-provider-at-capacity-parks-a-turn-the-fallback-chain-could-finish.md
regeneration:
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/autonomous-work.json
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/autonomous-work.md
      - docs/agents/setup-context.json
      - docs/agents/spec-routing.md
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/archive-spec/SKILL.md
      - .agents/skills/qa-gate/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "697"
---

# Runs are not lost to an authorization row or a provider at capacity

Qualifying declared partial: 10 rows passed, 3 rows are environment-blocked (two authorization-prohibited network lookups and the pre-PR Pull Request row), and exactly 1 row is declared unreachable. Repository Verification passed as supplied by the Daemon; no product finding remains. The offline archive plan exited 0 with the judge skipped and made no archive changes. The real-provider criterion remains unproven until the operator supplies the queue log and matching Run Event Journal after release. The Daemon owns Task settlement and evidence snapshots; no commit, push or Pull Request was performed.
