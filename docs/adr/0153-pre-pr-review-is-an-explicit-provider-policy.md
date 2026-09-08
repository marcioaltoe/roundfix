---
status: accepted
created_at: 2026-09-08T19:18:10Z
updated_at: 2026-09-08T19:18:10Z
deprecated_at: null
superseded_by: null
---

# Pre-PR review is an explicit provider policy

The maintainer clarified that pre-PR review may use Codex, Claude, CodeRabbit,
or none. This replaces the earlier requirements to remove CodeRabbit entirely
and always perform independent review. CodeRabbit remains optional. ADR-0151's
Codex default and explicit project-selection precedence remain; provider choice
also includes the service provider and the explicit absence of a reviewer.
Preserve User Config over built-in defaults and Project Config over User Config.
An Agent Selection Profile describes model/runtime configuration when an agent
provider is selected; it does not represent `none` as an Agent Runtime.

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
