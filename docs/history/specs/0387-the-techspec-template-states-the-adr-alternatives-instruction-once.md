---
schema: roundfix/archive-record/v1
spec: 0387-the-techspec-template-states-the-adr-alternatives-instruction-once
title: The TechSpec template states the ADR Alternatives instruction once
status: archived
created: "2026-11-12"
archived: "2026-11-13"
disposition: partial
source: docs/specs/0387-the-techspec-template-states-the-adr-alternatives-instruction-once
source_revision: 4ba6d05842230ef9b3fb74580d1b73e8945fcb3b
qa_task: task_02
qa_report: qa-report-2026-11-13.md
qa_verdict: partial
unproven:
  - the operator's or Argus's first managed refresh after the release, recorded session evidence.
adrs: []
sources:
  - 2026-11-12-argus-the-techspec-template-repeats-begin-the-relation-with.md
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/write-techspec/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-techspec/SKILL.md
promoted: []
type: fix
---

# The TechSpec template states the ADR Alternatives instruction once

The `write-techspec` skill's TechSpec template repeats "Begin the relation with" in the comment of its `## ADR Alternatives` section: one line ends with the phrase and the next line begins with it again, so the instruction that tells an author how to mark an ADR that does not apply reads as a stutter. Spec 0379 introduced the repetition, every adopter received it through `roundfix baseline update` with `write-techspec` 0.0.19, and Argus's pre-PR review reported it. After this Spec both copies of the template state the instruction once, the skill's version rises so adopters receive the correction through their next `roundfix baseline update`, and a contract test fails whenever a line of any Roundfix-owned skill ends with a phrase of two or more words that the next line begins with.
