---
schema: roundfix/archive-record/v1
spec: 0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
title: The Codex full-access mode is the one the adapter offers
status: archived
created: "2026-10-19"
archived: "2026-10-19"
disposition: partial
source: docs/specs/0281-the-codex-full-access-mode-is-the-one-the-adapter-offers
source_revision: d28c122732eba1aa162833c68eb67be00804d7a6
qa_task: task_03
qa_report: qa-report-2026-10-19-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0293
sources:
  - 2026-10-18-agent-full-access-pede-um-modo-que-o-codex-acp-recusa.md
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

# The Codex full-access mode is the one the adapter offers

The Archive Record can carry a qualifying partial: eight passed rows, three outside-evidence rows prohibited by the Spec's network authorization, and one pre-PR row without an open Pull Request. Prior F1 is resolved and no standing finding remains. The offline plan confirms one adopted incident-reference candidate; the reusable selection rule already lives in ADR-0293, CONTEXT.md and the configuration/runtime guides, so this QA pass promotes nothing. Daemon settlement, archive and review remain outside this child's actions.
