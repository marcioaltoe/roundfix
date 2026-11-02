---
schema: roundfix/archive-record/v1
spec: 0339-an-approved-baseline-update-deletes-the-retired-skills
title: An approved Baseline update deletes the Retired Skills
status: archived
created: "2026-11-01"
archived: "2026-11-02"
disposition: partial
source: docs/specs/0339-an-approved-baseline-update-deletes-the-retired-skills
source_revision: 495d9c30e132b134220f03e41565d9ca56832f52
qa_task: task_03
qa_report: qa-report-2026-11-02.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0246
  - ADR-0191
sources:
  - 2026-11-01-an-approved-baseline-update-deletes-the-retired-skills.md
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
promoted: []
---

# An approved Baseline update deletes the Retired Skills

Minimal PRD for a bug fix, minted by `write-techspec`. `roundfix baseline update` lists each Retired Skill an adopter still holds under `Skills retired` and deletes nothing, while the v0.127.0 release note said the update removes them; Argus removed seven copies by hand on 2026-11-01. After this Spec an approved update deletes each Retired Skill whose copies it proves unedited against `skills-lock.json`: the `.agents/skills/<skill>` folder, the `.claude/skills/<skill>` link or copy and the lock entry. The preview lists every deletion and the Plan Digest binds them, an edited or unrecorded copy is kept with its reason, and nothing outside the Retired Skills set is ever deleted, so `context7-cli`, `exa-web-search` and a repository's own skills survive. ADR-0350 supersedes ADR-0246.
