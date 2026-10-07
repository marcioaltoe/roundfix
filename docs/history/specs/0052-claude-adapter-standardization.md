---
schema: roundfix/archive-record/v1
spec: 0052-claude-adapter-standardization
title: Official Claude adapter and opaque model identifiers
status: archived
created: "2026-07-28"
archived: "2026-07-28"
disposition: pass
source: docs/history/specs/0052-claude-adapter-standardization
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-28.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0079
sources: []
regeneration: []
promoted: []
pull_request: "42"
delivery_commit: 85aad72727420f3d509109845be17f3477d50d55
---

# Official Claude adapter and opaque model identifiers

Roundfix still resolves the claude ACP Runtime to the deprecated `@zed-industries/claude-code-acp`, which npm renamed to `@agentclientprotocol/claude-agent-acp` — the same rename Roundfix already followed for Codex. The official adapter advertises exactly the model and reasoning controls Roundfix requires, but Roundfix misreads its advertised model identifiers (`opus[1m]`) as an embedded reasoning-effort encoding, so the identifiers Roundfix itself prints are unselectable, a context window is silently accepted as a reasoning effort, and Adapter Readiness proves lineage for Codex only. The maintainer directs Roundfix to support only the official Claude adapter, including in the Doctor and Setup Commands. Evidence and verification live in the [standardization finding](../../findings/2026-07-27-claude-adapter-standardization.md) and its [predecessor](../../findings/2026-07-26-claude-adapter-configoptions-migration.md).
