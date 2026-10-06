# Why the rollout is not written — investigation record

Recorded on 2026-10-06 while authoring Spec 0240. Everything here was read
without writing: `~/.codex`, `~/.acpx/sessions`, the installed adapter and acpx
bundles, and the Roundfix source. Measurements used a fake ACP adapter under a
temporary `HOME` and never started a Codex session.

## Installed versions

- `codex-cli 0.159.3` (`~/.local/bin/codex`), reached by `codex-acp` through
  the `CODEX_PATH` that Roundfix sets on macOS.
- `@agentclientprotocol/codex-acp` 2.1.1 (`/opt/homebrew/lib/node_modules`).
- acpx 0.19.4.

## Findings

1. **Roundfix sets no `CODEX_HOME`.** `internal/agent/codex_spawn.go` adds only
   `CODEX_PATH` to the acpx environment, and nothing in `internal/` sets
   `CODEX_HOME`. Every Run, Task and disposable proof shares `~/.codex`, so a
   per-Run Codex home is not the cause.
2. **Codex writes a thread's rollout lazily.** The adapter's own comment in
   `dist/index.js` (`CodexAcpClient.resumeThread`) says: "Codex writes a
   thread's rollout file on its first user message, so `thread/resume` fails
   with "no rollout found" for a session that was created but never prompted."
   The adapter recovers from that only while the original app-server is alive:
   it falls back to `thread/read` on the live app-server. A new `acpx prompt`
   process starts a new app-server, which has never seen the thread, so the
   original error comes back.
3. **The adapter turns the app-server error into `-32603`.** `errorToResult`
   maps any error that is not a `RequestError` to
   `RequestError.internalError({details})`. The line on the wire is
   `{"code":-32603,"message":"Internal error","data":{"details":"no rollout found for thread id <uuid>"}}`.
   The phrase is in `data.details`, not in `message`.
4. **acpx hides the empty case and exposes the worked case.** In acpx 0.19.4,
   `recoverRuntimeSessionLoadFailure` creates a fresh session after a failed
   `session/resume` when `shouldFallbackToNewSession` holds. For `-32603` that
   is only while the session record holds no Agent message
   (`!sessionHasAgentMessages(record) && isFallbackSafeEmptySessionError`).
   A session that has produced Agent output fails the prompt with exit `1`.
5. **Census of this machine.** Of the 993 acpx stream files in
   `~/.acpx/sessions` (2026-09-22 to 2026-10-05), 934 carry the error, all with
   code `-32603`, message `Internal error`, all answering `session/resume`, and
   all followed by acpx's fresh `session/new`. That pattern fits sessions
   created by `acpx sessions ensure` and first prompted from a later process,
   which never prompted them in the app-server that created them. Every one was
   recovered silently. The three thread ids in the Oraculum report have no acpx record and
   no rollout on this machine; the records were gone by 2026-10-06.
6. **Measured with a fake adapter.** A fake adapter answered with the exact
   `codex-acp` error shape after a first successful turn and an expired queue
   owner (`--ttl 1`). During `session/resume`, acpx 0.19.4 exited `1` with the
   error line on stdout and only its session banner on stderr. During
   `session/prompt`, after an `agent_message_chunk`, it did the same. Replayed
   through today's ACPX Runner with `go test -overlay`, the first became
   `Agent Selection failed for runtime "codex": agent/protocol error at session setup: Internal error (JSON-RPC -32603)`
   and the second
   `Agent Batch failed after acpx exited with code 1: agent/protocol error at session/prompt: Internal error (JSON-RPC -32603)`.
   The phrase was gone, and `session/resume` was not a tracked step, so
   resuming was reported as `session setup`.

## What remains open

How a thread that produced Agent output ends up with no rollout file. The
Codex tracker has the same shape on other clients:

- openai/codex#16872: a turn completes, but the rollout never materializes and a
  sibling `thread/resume` fails with `no rollout found`; the app-server logged
  `failed to record rollout items: failed to queue rollout items: channel closed`,
  and a commenter noted that turn completion does not flush the rollout while
  shutdown does.
- openai/codex#42099: since 0.151.0 `thread/start` no longer persists a
  zero-turn thread, so resuming it fails with the same message.
- openai/codex#28496: the app-server reports `no rollout found` with the
  generic `-32600`, and clients must match the message text.

Concurrency between Codex processes is not excluded. `~/.codex/thread-writer-locks`
was empty when read, and `codex-acp` reports a second writer with its own
`thread_active_writer` error, which none of the census lines carried. The
leading candidate is a rollout writer that never flushed before its app-server
process ended. That defect belongs to the adapter and Codex, and the report
below asks for it upstream.

## Upstream report text (not filed)

Title: `session/resume` of a thread with completed turns fails with `-32603 Internal error` (`no rollout found`) when the rollout was never flushed

> With `@agentclientprotocol/codex-acp` 2.1.1 and `codex-cli` 0.159.3,
> driven through acpx 0.19.4, a session that completed at least one turn
> sometimes cannot be resumed by a later process. `session/resume` answers
> `{"code":-32603,"message":"Internal error","data":{"details":"no rollout found for thread id <uuid>"}}`,
> and no `rollout-*-<uuid>.jsonl` exists under `~/.codex/sessions`, while the
> other sessions of the same day have one. Three threads failed this way on
> 2026-10-05 and 2026-10-06 in unattended runs.
>
> `CodexAcpClient.resumeThread` recovers only while the app-server that
> created the thread is alive (`thread/read` fallback). A new adapter process
> cannot, and the error reaches the client as a generic internal error.
>
> Requests:
>
> 1. Flush the rollout when a turn completes, or when a session closes, so a
>    thread that has done work can always be resumed (compare openai/codex#16872).
> 2. Report a thread the app-server cannot find with a typed error, for
>    example `RequestError.resourceNotFound` or `data.reason: "thread_not_found"`
>    as the adapter already does for `thread_active_writer`, instead of
>    `internalError` with the text in `data.details`, so clients can tell a
>    lost thread from a failed turn without matching English text
>    (compare openai/codex#28496).
>
> Reproduction without a model: start a session through acpx, complete one
> turn, let the queue owner exit, delete or withhold the rollout file, and
> prompt the session again.

## Appendix: recorded acpx 0.19.4 stdout

The working directory is normalized to `/tmp/w`. Both prompts exited `1`;
stderr held only the `[acpx] session w (...) · ... · agent starting` banner.

Lost during `session/resume`:

```text
{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{"fs":{"readTextFile":true,"writeTextFile":true},"terminal":true},"clientInfo":{"name":"acpx","version":"0.19.4"}}}
{"jsonrpc":"2.0","id":0,"result":{"protocolVersion":1,"agentCapabilities":{"loadSession":true,"sessionCapabilities":{"resume":{}}},"agentInfo":{"name":"fake-codex","version":"1"},"authMethods":[]}}
{"jsonrpc":"2.0","id":1,"method":"session/resume","params":{"sessionId":"fake-thread-0240","cwd":"/tmp/w","mcpServers":[]}}
{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"Internal error","data":{"details":"no rollout found for thread id 01a1fake-0000-7000-8000-000000000240"}}}
```

Lost during `session/prompt`:

```text
{"jsonrpc":"2.0","id":0,"method":"initialize","params":{"protocolVersion":1,"clientCapabilities":{"fs":{"readTextFile":true,"writeTextFile":true},"terminal":true},"clientInfo":{"name":"acpx","version":"0.19.4"}}}
{"jsonrpc":"2.0","id":0,"result":{"protocolVersion":1,"agentCapabilities":{"loadSession":true,"sessionCapabilities":{"resume":{}}},"agentInfo":{"name":"fake-codex","version":"1"},"authMethods":[]}}
{"jsonrpc":"2.0","id":1,"method":"session/resume","params":{"sessionId":"fake-thread-0240","cwd":"/tmp/w","mcpServers":[]}}
{"jsonrpc":"2.0","id":1,"result":{"sessionId":"fake-thread-0240","models":{"currentModelId":"default","availableModels":[{"modelId":"default","name":"Default"}]}}}
{"jsonrpc":"2.0","id":2,"method":"session/prompt","params":{"sessionId":"fake-thread-0240","prompt":[{"type":"text","text":"second turn"}]}}
{"jsonrpc":"2.0","method":"session/update","params":{"sessionId":"fake-thread-0240","update":{"sessionUpdate":"agent_message_chunk","content":{"type":"text","text":"working"}}}}
{"jsonrpc":"2.0","id":2,"error":{"code":-32603,"message":"Internal error","data":{"details":"no rollout found for thread id 01a1fake-0000-7000-8000-000000000240"}}}
```
