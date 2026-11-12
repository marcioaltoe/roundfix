---
schema: roundfix/archive-record/v1
spec: 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
title: An older Roundfix refuses a repository a newer Roundfix wrote
status: archived
created: "2026-11-11"
archived: "2026-11-12"
disposition: pass
source: docs/specs/0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
source_revision: 6999de103cb7f437f66adfe968047d433837aa5b
qa_task: task_04
qa_report: qa-report-2026-11-12.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0392
sources:
  - 2026-11-11-conexus-an-older-roundfix-ran-against-a-newer-repository.md
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
pull_request: "840"
type: feat
---

# An older Roundfix refuses a repository a newer Roundfix wrote

The initial plan invocation used a temporary HOME under the Run Worktree and correctly refused worktree.location before any side effect; initial refusal (`docs/specs/0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote/qa/evidence/qa-report-2026-11-12/archive-plan-initial.txt` (not reachable from any branch or tag)) is retained. The corrected offline archive plan was executed with the built binary, isolated HOME/TMPDIR, preserved Go cache locations, fake gh and all Judge Key Variables plus TYPESAFE_API_KEY unset; it exited 0. Archive plan and cut (`docs/specs/0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote/qa/evidence/qa-report-2026-11-12/archive-plan.txt` (not reachable from any branch or tag)) record the runtime result.
