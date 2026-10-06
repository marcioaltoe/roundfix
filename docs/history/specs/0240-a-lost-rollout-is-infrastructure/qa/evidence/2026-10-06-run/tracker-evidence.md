# Upstream evidence observed 2026-10-06

Direct curl-cffi requests to api.github.com failed with CONNECT 403. The web tool reached the primary issue pages, so this source row is reachable through the alternate read-only boundary.

- [Codex issue 16872](https://github.com/openai/codex/issues/16872), opened 2026-04-05: the reporter describes completed turns with absent rollout persistence, successful lightweight thread reads, and sibling resume failure. Supports recovery for a loss after Agent output; does not establish this machine's writer-flush cause.
- [Codex issue 42099](https://github.com/openai/codex/issues/42099), opened 2026-09-01: the reporter contrasts 0.150.1 with 0.151.0; zero-turn threads are indexed without rollout persistence and cannot resume. Supports distinguishing empty-session recovery from worked-session loss.
- [Codex issue 28496](https://github.com/openai/codex/issues/28496), opened 2026-06-16: the reporter describes -32600 used for both stale threads and malformed config. Supports matching the diagnostic phrase as well as the code; code-only fallback would hide unrelated config defects.

All three pages showed Open. This is external issue-report evidence, not independent reproduction of Codex itself. The assertions above are concise paraphrases of the issue bodies; no comments or unobserved writer-flush claims are credited.
