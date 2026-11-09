---
schema: roundfix/archive-record/v1
spec: 0366-a-qa-report-relies-only-on-evidence-its-own-run-produced
title: A QA Report relies only on evidence its own Run produced
status: archived
created: "2026-11-09"
archived: "2026-11-09"
disposition: pass
source: docs/specs/0366-a-qa-report-relies-only-on-evidence-its-own-run-produced
source_revision: 6133c99b0deb03c75625b432d621f4e625fc5e0b
qa_task: task_04
qa_report: qa-report-2026-11-09.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0351
  - ADR-0376
sources:
  - 2026-11-07-a-later-runs-qa-gate-reused-an-earlier-runs-evidence.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/qa-gate/SKILL.md
      - .agents/skills/roundfix/SKILL.md
      - skills/qa-gate/SKILL.md
      - skills/roundfix/SKILL.md
      - skills/testdata/owned-skill-versions.json
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "795"
type: fix
---

# A QA Report relies only on evidence its own Run produced

The gate observed that borrowed QA evidence is refused by the public CLI and shared settlement paths, with one Daemon repair turn and no refused QA Report commit; regenerated evidence and the documented compatibility cases passed. All nine independent sabotages were detected. The offline archive plan exited 0 and reported `judge: skipped (offline)`, with only the adopted Quaestor measurement listed as a non-core candidate. Its historical incident remains recoverable with the Spec at the Archive Record source revision; the reusable rule is already in CONTEXT.md, ADR-0376 and the maintained guides and skills. No network lookup, provider, judge, review, archive mutation, commit, push or Pull Request creation was performed by this gate.
