---
schema: roundfix/archive-record/v1
spec: 0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd
title: The Secondbrain search and the duplicate check go through the Indexus, without qmd
status: archived
created: "2026-11-20"
archived: "2026-11-20"
disposition: partial
source: docs/specs/0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd
source_revision: 3d0f4b91f33fcc5288e2c5a65e455567ed117a9a
qa_task: task_02
qa_report: qa-report-2026-11-20-01.md
qa_verdict: partial
unproven:
  - The notice of the test's result that the project that owns the Secondbrain sends to this project's inbox, which the operator records before the merge.
  - The author's live calls of 2026-11-20, recorded in the characterization, section 5.
adrs:
  - ADR-0257
  - ADR-0360
sources: []
regeneration:
  - command: make skills-sync
  - command: go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions
    outputs:
      - .agents/skills/write-idea/SKILL.md
      - skills/testdata/owned-skill-versions.json
      - skills/write-idea/SKILL.md
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/secondbrain.json
      - internal/baseline/module-versions.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - .secondbrain-export
      - docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md
      - docs/agents/secondbrain.md
      - docs/agents/setup-context.json
promoted: []
pull_request: "944"
type: feat
---

# The Secondbrain search and the duplicate check go through the Indexus, without qmd

The candidate directs Secondbrain search and duplicate checks through the Indexus with no second search tool; the isolated Delivery Base adopter receives only its guide and …
