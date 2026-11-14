---
schema: roundfix/archive-record/v1
spec: 0390-a-nested-carrier-git-does-not-list-raises-no-conflict
title: A nested carrier Git does not list raises no conflict
status: archived
created: "2026-11-13"
archived: "2026-11-13"
disposition: partial
source: docs/specs/0390-a-nested-carrier-git-does-not-list-raises-no-conflict
source_revision: 9cb7119c33c76a082cbcea43861ef189244253a1
qa_task: task_03
qa_report: qa-report-2026-11-13.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0399
sources:
  - 2026-11-13-quaestor-a-nested-carrier-conflict-inside-a-git-ignored-worktree.md
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
pull_request: "867"
type: fix
---

# A nested carrier Git does not list raises no conflict

The Baseline now omits ignored and foreign nested instruction carriers while retaining tracked and untracked nonignored carriers, leaving its digest stable when an ignored worktree is added. The offline archive preview exited 0 and skipped the judge; it identified the adopted quaestor message as the only non-core promotion candidate. ADR-0399, the ADR-0070 amendment and the published guide/reference already retain the reusable rule, so no additional promotion is needed. The QA result is a qualifying partial solely for the pre-PR control row, with 20 passes and no product findings; the Daemon owns Task settlement and any later archive.
