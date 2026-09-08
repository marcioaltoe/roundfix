---
status: proposed
granted: null
created: 2026-09-08
action: A Spec exposes its executable promises and recovery path
consuming: 0129-spec-authoring-and-gate-recovery
paths:
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-tasks/references/task-template.md
  - .agents/skills/implement-task/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-techspec/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-tasks/references/task-template.md
  - skills/implement-task/SKILL.md
  - internal/speccheck/coherence.go
  - internal/baseline/assets/modules/spec-workflow.json
  - docs/agents/spec-routing.md
  - docs/agents/setup-context.json
---

# Proposed bounded implementation authority

The maintainer selected this intent for the queue. This record is a proposal,
not authority to mutate the named governed paths. The complete technical
candidate and its exact path changes must be reviewed before the grant becomes
operative; the grant lands separately before consuming changes.

Ordinary product sources belong in the TechSpec implementation map. Source
skills remain authoritative; `make skills-sync` and sanctioned
`make baseline-digests` describe regeneration after source approval. Generated
guides use a reviewed public Baseline plan. No mirror, credential, release,
workflow, dependency or unspecified file mutation is granted here.

Implementation, correction, native review and publication retain the agreed
through-merge gate. Runtime access and resource limits are recorded separately;
missing values and an unanswered question are not unlimited authority.
