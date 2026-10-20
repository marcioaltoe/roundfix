---
schema: roundfix/archive-record/v1
spec: 0289-gc-sanitize-reaches-the-default-artifact-roots
title: GC sanitize reaches the default Artifact Roots
status: archived
created: "2026-10-20"
archived: "2026-10-20"
disposition: partial
source: docs/specs/0289-gc-sanitize-reaches-the-default-artifact-roots
source_revision: fa922b19828a2912d832c9806786f4a8acc95e85
qa_task: task_04
qa_report: qa-report-2026-10-20.md
qa_verdict: partial
unproven:
  - the maintainer runs `roundfix gc sanitize`, and then `roundfix gc sanitize --apply` only if the dry run is clean, after the release that carries this Spec.
adrs:
  - ADR-0301
sources:
  - 2026-10-20-gc-sanitize-preserves-every-artifact-root-a-run-recorded-empty.md
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

# GC sanitize reaches the default Artifact Roots

The Archive Record can carry this qualifying partial and the single real-home declaration as unproven. The offline archive plan exited 0 and skipped the judge; its adopted defect report needs no promotion because the durable rules already live in the guide, skill reference, glossary and ADR-0301. No archive, Task status change, commit, push or Pull Request was performed by this gate.
