---
schema: roundfix/archive-record/v1
spec: 0400-planning-references-and-release-preparation-need-no-manual-classification
title: Planning references and release preparation need no manual classification
status: archived
created: "2026-11-16"
archived: "2026-11-17"
disposition: qa-override
source: docs/specs/0400-planning-references-and-release-preparation-need-no-manual-classification
source_revision: f21929b04cdd6fdc56c4150a4aa87f0190b92e85
qa_task: task_04
qa_report: ""
qa_verdict: ""
unproven: []
qa_override: true
qa_override_approval: Marcio, standing authorization 2026-08-09 (qa_override por ambiente, rows_blocked_finding=0) and 2026-11-17 unattended grant
qa_override_reason: 'QA partial only by environment (rows_blocked_finding 0): isolated HOME had no rustup toolchain for TestGuidanceCompositionJourney/rust-cli (same failure on the unchanged Delivery Base); the full ./internal/cli and ./internal/releaseplan packages pass with the host toolchain on f21929b0, and make verify-changed passes; the PR and Daemon verification rows are satisfied by the delivery checks.'
qa_override_qa_outcome: missing
qa_override_qa_task_status: pending
qa_override_revision: f21929b04cdd6fdc56c4150a4aa87f0190b92e85
adrs:
  - ADR-0412
sources:
  - 2026-11-16-a-planning-change-with-references-makes-every-release-a-manual-classification.md
  - 2026-11-16-the-release-guide-assumes-tag-triggered-releases.md
  - 2026-11-16-maintainer-answers.md
  - 2026-11-16-release-plan-classification-measurement.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/roundfix/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted:
  - docs/references/2026-11-16-the-release-guide-assumes-tag-triggered-releases.md
relinked:
  - docs/backlog/2026-11-16-the-delivery-queue-and-the-release-guide-assume-loose-branch-rules-and-tag-releases.md
pull_request: "897"
type: fix
---

# Planning references and release preparation need no manual classification

The Release Plan Command asks for a classification by hand for two kinds of commit that ship nothing. The first is a Spec's planning commit that adds outside evidence under `docs/references/`, as `docs/agents/spec-routing.md` requires. The second is the release-preparation commit that writes the changelog and the version files. A Release Ceiling never answers a classification by hand, so in adopters that record ceilings every release stops and asks again. The Oraculum met it on 2026-11-16, and Roundfix met it at `v0.143.4`. The release guide also tells every adopter to preserve a tag-triggered workflow, while the Quaestor publishes its release when the version change merges. After this Spec, the plan classifies both commits itself, and the release guide describes both publication triggers.
