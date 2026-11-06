---
schema: roundfix/archive-record/v1
spec: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
title: The build cache stays within a bound, and a full disk is named
status: archived
created: "2026-11-06"
archived: "2026-11-06"
disposition: partial
source: docs/specs/0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
source_revision: d2112e31283e04ae4fa71c5f29cc85c0390899b4
qa_task: task_04
qa_report: qa-report-2026-11-06.md
qa_verdict: partial
unproven:
  - the operator's queue records and `du` of the cache after the release that carries this Spec, recorded session evidence.
adrs:
  - ADR-0370
  - ADR-0327
  - ADR-0334
sources:
  - 2026-11-05-the-roundfix-build-cache-fills-the-mac-mini-internal-disk.md
  - 2026-11-06-the-repository-build-cache-stays-within-a-size-bound.md
regeneration:
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "774"
type: fix
---

# The build cache stays within a bound, and a full disk is named

The isolated gate confirms that the Makefile trims only idle cache entries, Doctor and the queue report low space without stopping work, and a full filesystem during the queue's build or gate is named disk-full with an uncounted retry. The report closes partial with 40 passed rows, one declared post-release production-day measurement still unproven and one Pull Request row blocked because none is open; these are the only unmet rows.
