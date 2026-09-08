---
status: proposed
granted: null
action: align knowledge capture, terminal dispositions, and history guidance
consuming: 0120-knowledge-lifecycle-and-durable-capture
paths:
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/secondbrain.json
  - internal/spec/archive.go
  - internal/spec/archive_test.go
  - internal/speccheck/backlog.go
  - internal/speccheck/backlog_test.go
  - internal/docscontract/publicdocs_test.go
  - .agents/skills/archive-spec/SKILL.md
  - skills/archive-spec/SKILL.md
  - docs/agents/docs-layout.md
  - docs/agents/secondbrain.md
  - docs/agents/setup-context.json
---

# Proposed authority for Spec 0120

This record grants no implementation action. The maintainer authorized
triage, preservation of completed documentation, research capture, and
authoring/commit/push of the plan. The exact protected changes above remain
proposed until a maintainer decision is recorded.

## Bounded proposal

Complete the canonical lifecycle/capture clauses, their governed consistency
checks, and the Roundfix-owned archive skill. Use the existing History Root
and archive resolver. `make baseline-digests` produces sanctioned pins;
`make skills-sync` may change only the named skill output; the public Baseline
update regenerates only the named managed guides and manifest.

The companion Secondbrain contract changes are proposed separately at exact
repository-relative paths `AGENTS.md` and `inbox/README.md` in repository
`marcioaltoe/secondbrain`. They are not Roundfix paths and must not be resolved
against the Roundfix working directory. The proposal changes guidance about
the already installed autosync owner; it does not authorize changing jobs,
scripts, credentials, mirrors, or immutable sources.

## Limits

- Approval of one repository's files does not imply approval in the other.
- Keep original observations and legacy Spec/authorization bytes intact;
  record new dispositions in dated addenda and preserve all archive licenses.
- Do not promise remote durability from the existence of an autosync job.
- A capture fallback must not create a second schedule or publish unrelated
  concurrent changes.
- Commit the final granted record separately before consuming tooling changes.
- No paid API use, release, deployment, or destructive cleanup is included.

## Approval evidence

Pending. Preserve `status: proposed` and `granted: null` until the maintainer
approves the concrete files and unresolved lifecycle/ownership decisions.
