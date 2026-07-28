---
status: accepted
created_at: 2026-07-28T20:02:08Z
updated_at: 2026-07-28T20:02:08Z
deprecated_at: null
superseded_by: null
---

# Independent reasoning controls make advertised model identifiers opaque

When an ACP adapter advertises an independent reasoning-effort configuration
option, Roundfix treats every advertised Agent Model identifier as an opaque
selectable value: a bracketed suffix such as `opus[1m]` is an adapter
annotation (a context window), not Roundfix's `canonical[effort]` variant
encoding, so parsing it as a reasoning effort made advertised identifiers
unselectable and silently accepted a context window as an effort. The variant
encoding remains only for adapters that advertise no independent reasoning
control. This supersedes ADR-0040's premise that the Claude adapter lacks the
reasoning configuration option — official
`@agentclientprotocol/claude-agent-acp` advertises both controls — while
preserving ADR-0040's explicit-empty semantics, and refines ADR-0055's
capability-driven assignment.
