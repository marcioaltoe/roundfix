---
schema: roundfix/archive-record/v1
spec: 0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web
title: Adopting the recommendations works under load and gives the researcher the web
status: archived
created: "2026-11-18"
archived: "2026-11-19"
disposition: partial
source: docs/specs/0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web
source_revision: 656f630f02a871f5457d0ae9a1a9c8b2400c96fb
qa_task: task_03
qa_report: qa-report-2026-11-19.md
qa_verdict: partial
unproven:
  - the next adoption by fiscus or fluxus after the release that carries this Spec
adrs:
  - ADR-0441
  - ADR-0402
sources: []
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
type: fix
---

# Adopting the recommendations works under load and gives the researcher the web

Profile-write tests confirm that both adoption commands pass the configured setup deadline with Project Config precedence and reject invalid deadlines or broken YAML before proof or persistence. The recommendation produces a researcher with both web tools, the reports respect explicit tool choices, and guides, mirrors, glossary, versions and scope checks agree. All runnable checks and all eight sabotage probes pass; actual adoption on a loaded machine still awaits the declared post-release adopter run, and Pull Request acceptance awaits delivery. The offline archive plan reports no judge call and identifies two historical reference candidates whose durable rules already live in the accepted ADRs and documentation.
