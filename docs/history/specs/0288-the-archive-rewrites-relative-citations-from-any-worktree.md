---
schema: roundfix/archive-record/v1
spec: 0288-the-archive-rewrites-relative-citations-from-any-worktree
title: The archive rewrites relative citations from any worktree
status: archived
created: "2026-10-20"
archived: "2026-10-20"
disposition: partial
source: docs/specs/0288-the-archive-rewrites-relative-citations-from-any-worktree
source_revision: ebe05f3b4f4ea9e0866e10478ce487bd69148f3a
qa_task: task_03
qa_report: qa-report-2026-10-20.md
qa_verdict: partial
unproven:
  - the operator's next `qa-environment-partial` recovery, or any archive run from an item worktree, after this Spec's release.
adrs:
  - ADR-0268
  - ADR-0272
sources:
  - 2026-10-20-o-archive-com-qa-override-nao-reescreve-link-relativo-no-item-worktree.md
regeneration:
  - command: go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
pull_request: "607"
---

# The archive rewrites relative citations from any worktree

The archive and cancel now print and persist identical results from the required resolved and symbolic-link checkout spellings, including relative-link rebasing and Live Document relinking. Regression tests, refusal checks, Windows vet, glossary and coverage contracts, authorization and Delivery Base invariance checks passed; both sabotage defects were detected, and the Delivery Base independently reproduced the reported archive refusal. The verdict is a qualifying declared partial because the real post-release item-worktree acceptance remains unproven, while external network sources are prohibited by the Spec and no Pull Request exists.
