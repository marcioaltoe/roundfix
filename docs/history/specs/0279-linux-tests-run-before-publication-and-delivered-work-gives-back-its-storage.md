---
schema: roundfix/archive-record/v1
spec: 0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
title: Linux tests run before publication, and delivered work gives back its storage
status: archived
created: "2026-10-18"
archived: "2026-10-18"
disposition: partial
source: docs/specs/0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage
source_revision: 19b177e8656248295a4b61347c2681bc4353821a
qa_task: task_04
qa_report: qa-report-2026-10-18.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0288
  - ADR-0289
  - ADR-0290
sources: []
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - skills/testdata/owned-skill-versions.json
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
---

# Linux tests run before publication, and delivered work gives back its storage

The offline archive preview exited 0 and printed `judge: skipped (offline)`, naming the cut and `docs/specs/0279-linux-tests-run-before-publication-and-delivered-work-gives-back-its-storage@b29c54309d36daf812c1cacbd33c114f3cb1d985` destination. Its file/byte counts describe the report as it existed at preview time; later evidence additions do not change the preview's scope decision. Core planning and QA evidence are Spec-owned; durable vocabulary and policy already reside in CONTEXT.md, accepted ADRs and user guides, so no promotion candidate was identified. The Archive Record should carry partial, with the three prohibited publication lookups and the absent Pull Request explicitly unverified. This matches the qualifying offline partial policy; the Daemon owns final QA settlement and any archive. Evidence: `evidence/2026-10-18-gate/archive-plan.txt`.
