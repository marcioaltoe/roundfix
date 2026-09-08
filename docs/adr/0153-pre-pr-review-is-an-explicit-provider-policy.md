---
status: accepted
created_at: 2026-09-08T19:18:10Z
updated_at: 2026-09-08T19:18:10Z
deprecated_at: null
superseded_by: null
---

# Pre-PR review is an explicit provider policy

The maintainer chose optional pre-PR review through Codex, Claude, CodeRabbit,
or explicit `none` to make the cost and assurance trade-off visible. This
supersedes the requirements to remove CodeRabbit entirely and always perform
independent review while preserving ADR-0151's Codex default and configured
precedence. User Config overrides built-in defaults, Project Config overrides
User Config, and Agent Selection Profiles supply runtime/model configuration
only when an agent provider is selected.

## Consequences

With `codex`, `claude`, or `coderabbit`, review examines the current candidate
before the PR and its result must satisfy the enabled review policy. Explicit
`none` performs no review/provider call and allows otherwise authorized delivery
with QA and required checks, recording that review was disabled rather than
passed. A failed, unavailable, incomplete, stale or missing enabled review is
an error, never implicit permission to select none. External branch protection
and required checks retain their authority.

The local CodeRabbit CLI is the proposed pre-PR integration surface; its PR
feedback service is a separate surface. Configuration schema, adapter support
and evidence persistence remain implementation work. Declaring these choices
does not install a provider, select none for this repository, or approve
additional billing. The trade-off of explicit none is delivery without an
independent review result; its deliberate omission remains visible.

The decision is the maintainer's. Exa located and read the official
[CodeRabbit CLI reference](https://docs.coderabbit.ai/cli/reference), which
supports local Git review and structured agent output; this establishes an
integration surface, not implemented Roundfix support. The earlier Secondbrain
review captures remain historical evidence of the prior mandatory policy.
