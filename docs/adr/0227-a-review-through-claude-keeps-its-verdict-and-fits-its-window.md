---
status: accepted
created_at: 2026-10-04T00:00:00Z
updated_at: 2026-10-04T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A review through Claude keeps its verdict and fits its window

On 2026-10-02 two repositories reviewing through `claude / opus / high` saw
`roundfix review` block with `review transport anomaly: acpx exited with exit
code 5 after parsed session/prompt result` while the answer file held a
verdict, and one saw a candidate fail as `agent/protocol error` while the
answer file said `Prompt is too long`. A reproduction on 2026-10-04 showed the
mechanism: the reviewer ran `go vet` through its terminal tool, the read-only
session answered the `session/request_permission` with `reject`, the turn
continued and ended with `stopReason: end_turn` and one finding, and acpx
exited `5` because every permission request of the turn had been denied.
Read-only commands such as `sed` and `grep` never asked, so the defect appears
only when the reviewer tries a command Claude Code does not classify as
read-only.

A read-only Agent Session whose turn ended with `end_turn` and whose acpx exit
is `5` therefore delivered its answer: the runner keeps the result, marks that
the session refused a permission, and publishes the permission-denied status;
`roundfix review` classifies the answer by its verdict, so findings get their
ids and can be disposed, and the record says the refusal happened. Every other
exit after a parsed result stays a transport anomaly, and a read-write session
is unchanged.

Denying the tools up front was preferred and is not reachable: acpx 0.19.4
forwards only `allowedTools`, `maxTurns`, the model and the system prompt to
the Claude adapter, and the Agent SDK treats `allowedTools` as pre-approval,
not restriction; the adapter's `tools` and `disallowedTools` options and the
`dontAsk` mode are not offered through acpx. A prompt instruction is not
enough either: the review prompt already says not to run a shell, and both
reproductions used the terminal anyway.

The Claude review prompt is bounded in tokens, not bytes: the whole prompt the
command sends, diff, instructions and Spec context, may use at most half of
the 1,000,000-token context window the adapter reports, estimated at two bytes
per token. The other half covers what the session adds, measured at about
25,000 tokens of system prompt and tools, about 57,000 tokens of files that
`@path` mentions in a diff attach, the reviewer's own reads and up to 128,000
output tokens. The Codex review keeps the 917,504-byte diff bound measured on
Codex. When a runtime still answers `Prompt is too long`, the command's reason
quotes that line instead of the protocol error.

## Consequences

A Claude reviewer that tries a permission-requiring command no longer loses
its verdict, and its findings follow the normal disposition path. The Claude
bound admits a prompt of about one million bytes, a little more than the Codex
diff bound, because it now counts what is sent against the provider's own
window. If acpx later forwards `disallowedTools` to the Claude adapter, the
review session can deny the terminal up front and this classification remains
the fallback.
