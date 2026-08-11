---
status: accepted
created_at: 2026-07-17T16:00:35Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Agent Selection encoding follows advertised ACP capabilities

Roundfix represents an Agent Selection as the exact ACP Runtime, Agent Model,
and reasoning-effort tuple requested by the user, but adapters do not expose
that tuple through one universal control shape. An adapter may advertise model
and reasoning as independent configuration options, encode reasoning in an
advertised model variant, or expose no way to select the requested reasoning
at all. Roundfix therefore discovers the controls advertised by a disposable
Agent Session, derives an adapter-specific assignment plan, applies it, and
proves that the resulting session represents the exact requested tuple before
configuration or Run mutation.

The effective adapter command, package lineage, and version are part of this
proof boundary. Roundfix Setup must not replace ACPX's built-in adapter with a
same-named executable found on `PATH` unless that executable is proven to be a
supported official adapter. A stale explicit override is incompatible even
when the surrounding ACPX and Codex CLI versions are current. Migration of an
existing override is diagnosed and proposed to the user; it is never edited
silently.

Static model-family assumptions, private runtime caches, and silent conversion
of a non-empty reasoning effort to model-managed reasoning are rejected. An
empty Default Reasoning Effort remains an explicit model-managed selection as
defined by ADR-0040, but Roundfix never infers that empty value merely because
a model belongs to the GPT-5.6 family. This supersedes only ADR-0040's
family-wide classification of GPT-5.6 reasoning behavior and preserves its
explicit-empty semantics. It also refines ADR-0039 and ADR-0049 by requiring
the same capability-driven proof boundary for setup, profile validation,
Doctor, and operational Preflight Validation.

A bare `--agent` value cannot prove a complete one-Run Agent Selection and is
therefore invalid for profile-driven commands. A deliberate one-Run override
must provide runtime, model, and reasoning effort together; omitting all three
uses the resolved Agent Selection Profile.
