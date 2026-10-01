# Published outside evidence — observed 2026-10-01

Sources read through web open/find. curl-cffi failed with CONNECT tunnel 403; the independent web read succeeded for both sources.

[ACP session setup](https://agentclientprotocol.com/protocol/v1/session-setup), sections Resuming Sessions / Checking Support / Resuming a Session (lines 169–232 of observed text): resume support must be advertised in sessionCapabilities.resume before the client calls session/resume. Resume restores context without replay. Loading Sessions / Checking Support (lines 67–89) separately requires loadSession support. This supports capability-based continuation and the Spec’s requirement to observe matching ACP session IDs rather than infer continuation from a reused name.

[GitHub review comments reference](https://docs.github.com/en/rest/pulls/comments), About pull request review comments (line 26) and Create a review comment parameters (lines 357–364): comments attach to the unified diff; the line/path and side fields locate the comment in that diff, with context and additions on the right. This independently supports checking candidate-diff anchors rather than accepting arbitrary repository lines.

Comparison: published contracts agree with the adopted resume-capability and anchor rules. They establish protocol semantics; current implementation behavior is measured separately by the hermetic tests and built-CLI transcripts. No provider session or GitHub mutation was performed.
