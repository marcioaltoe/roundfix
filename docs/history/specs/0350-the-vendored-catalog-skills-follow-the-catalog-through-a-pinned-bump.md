---
schema: roundfix/archive-record/v1
spec: 0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
title: The vendored catalog skills follow the catalog's main through a pinned bump
status: archived
created: "2026-11-04"
archived: "2026-11-04"
disposition: partial
source: docs/specs/0350-the-vendored-catalog-skills-follow-the-catalog-through-a-pinned-bump
source_revision: 5361eed2d51cedd352b55b2f40b158853cde8a06
qa_task: task_04
qa_report: qa-report-2026-11-04-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0224
  - ADR-0360
sources:
  - 2026-11-04-the-vendored-catalog-skills-trail-the-catalog-and-name-branch-refs.md
regeneration:
  - command: make skills-third-party-bump SKILLS_CATALOG=/Users/marcio/.roundfix-operator/scratch/spec-0350/skills SKILLS_REF=c8b2e8511a7ff89dbe9e87a56cc3c29a6acbcb84
    outputs:
      - .agents/skills/systematic-debugging/CREATION-LOG.md
      - .agents/skills/systematic-debugging/test-academic.md
      - .agents/skills/systematic-debugging/test-pressure-1.md
      - .agents/skills/systematic-debugging/test-pressure-2.md
      - .agents/skills/systematic-debugging/test-pressure-3.md
      - docs/agents/setup-context.json
      - internal/baseline/assets/setups/go-cli-typescript-bun.json
      - internal/baseline/assets/setups/go.json
      - internal/baseline/assets/setups/python.json
      - internal/baseline/assets/setups/rust.json
      - internal/baseline/assets/setups/typescript.json
      - skills-lock.json
  - command: roundfix baseline assets sync --source-dir <fetched>/setups --format text
    outputs:
      - internal/baseline/assets/setups/go-cli-typescript-bun.json
      - internal/baseline/assets/setups/go.json
      - internal/baseline/assets/setups/python.json
      - internal/baseline/assets/setups/rust.json
      - internal/baseline/assets/setups/typescript.json
  - command: roundfix baseline skills restore --profile go-cli-tui --skill <each catalog entry> --source-dir <fetched> --repo . --confirm-plan <digest> --format text
    outputs:
      - .agents/skills/systematic-debugging/CREATION-LOG.md
      - .agents/skills/systematic-debugging/test-academic.md
      - .agents/skills/systematic-debugging/test-pressure-1.md
      - .agents/skills/systematic-debugging/test-pressure-2.md
      - .agents/skills/systematic-debugging/test-pressure-3.md
      - skills-lock.json
  - command: make baseline-digests
  - command: go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text
    outputs:
      - docs/agents/setup-context.json
promoted: []
type: fix
---

# The vendored catalog skills follow the catalog's main through a pinned bump

The 33 vendored catalog skills, their lock entries and the five upstream Setup Snapshots agree on the pinned Catalog Revision, the obsolete systematic-debugging files are absent, and the 12 owned skills remain intact. The real repeat bump succeeds without changes, the suites and nine required negative probes meet their expectations, and the history audit finds no scope violation. This is a qualifying partial solely because no Pull Request is open and current-head review and Merge-Ready acceptance do not yet exist. The offline archive preview skipped the judge and listed the adopted catalog message as its sole promotion candidate; no archive or promotion was performed.
