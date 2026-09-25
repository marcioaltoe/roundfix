---
status: approved
granted: 2026-09-08
action: make user-authorized QA override an explicit canonical archive policy
consuming: 0122-verified-content-and-terminal-settlement
paths:
  - internal/baseline/assets/modules/spec-workflow.json
  - internal/baseline/assets/modules/context-workflow.json
  - internal/baseline/assets/modules/autonomous-work.json
  - docs/agents/docs-layout.md
  - docs/agents/skill-dispatch.md
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
---

# Canonical QA archive override authorization

The maintainer directed on 2026-09-08: "Outra regra que deve ser possível é o
arquivamente de specs com override de qa quando solicitado ou autorizado pelo
usuário". Earlier instructions require every canonical rule to live in the
Baseline modules and its generated guidance.

This grants the named source-policy and derived-guidance edits. An explicit
user request or applicable prior authorization may waive the QA prerequisite
for archiving its covered Spec. Record the request/approval, scope, date and
actual QA outcome or absence, and stamp `qa_override: true`. Preserve original
QA reports, Task states and implementation evidence. The archive may not be
reported as QA passed or the Run as Clean because an override was consumed.
Only the QA gate's archive prerequisite is waived; non-QA Task completion and
self-contained source references remain required. Publication, merge, release
and external checks retain their separate authority.

This record authorizes the policy capability, not an override for any particular
existing Spec. The broader 0122 implementation remains proposed. No Go, test,
skill source, configuration parser, runtime flag or account mutation is granted
here. The seven named files and sanctioned digest regeneration form this narrow
canonical scope. The command/skill integration must be implemented through the
consuming Spec before it is represented as a supported runtime operation.

Commit this record separately before its consuming tooling changes. Known
regeneration/verification compatibility failures remain unresolved; this grant
does not suppress them or claim a passing repository gate.

## Sanctioned regeneration

The repository-owned command resolves its generated outputs. This declaration
records the digest regeneration already granted above; it adds no source paths.

```yaml
command: make baseline-digests
```
