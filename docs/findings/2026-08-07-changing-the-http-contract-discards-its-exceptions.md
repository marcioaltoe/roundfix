---
status: deferred
created_at: 2026-08-07
updated_at: 2026-08-26
kind: finding
---

# Changing the HTTP Contract discards its exceptions (2026-08-07)

Answering "Change" on the `http.contract` prompt replaces the complete typed
decision with a mode and nothing else. Every recorded exception and the source
provenance are silently dropped.

## What was observed

`internal/cli/baseline_human.go`, the `http-contract` branch of
`promptBaselineDecision`, prompts only for the mode and returns:

```go
return map[string]any{"mode": declaration.Modes[selected]}, nil
```

The declaration's schema accepts `mode`, `exceptions`, and optional `source` —
`normalizeHTTPContract` validates all three. The prompt constructs only the
first. There is no follow-up prompt for exceptions, so a maintainer cannot add,
edit, or even retain one through the interactive command.

A repository observed on 2026-08-07 carried four exceptions and a source:

```json
{"exceptions":[
   {"scope":"/api/auth/*","methods":["GET","POST"],"owner":"Better Auth", ...},
   {"scope":"/health","methods":["GET"],"owner":"…operations", ...},
   {"scope":"/openapi.json","methods":["GET"],"owner":"…API documentation", ...},
   {"scope":"/reference","methods":["GET"],"owner":"…API documentation", ...}],
 "mode":"Post-only",
 "source":{"digest":"397bc399…","path":"packages/backend/src/infra/controllers/http/app.ts"}}
```

Answering "Change" and selecting a mode reduces that to `{"mode":"<selected>"}`.

## Why it matters

The loss is silent and the prompt gives no sign it is destructive. Worse, it is
on the only path to changing the mode: a maintainer who wants a different mode
must pay every exception to get it, including the provider-owned
`/api/auth/*` exception that the auth provider decision separately proposes and
depends on.

The "Keep" branch is safe — it returns the stored value whole — so the defect is
reachable only by the maintainer who wants to change something, which is exactly
the maintainer who has the most to lose.

## Workaround

Answer "Keep" whenever the stored mode is already the wanted one. It preserves
the complete typed value. There is currently no way to add an exception through
the interactive command.

## Route

Not fixed here. It is product behavior in `internal/cli`, not protected tooling,
so it needs no tooling authorization — but it does need a Spec, because the fix
has a real design question: whether "Change" should prompt for the full ordered
exception list, edit it entry by entry, or accept the mode change while carrying
the existing exceptions forward untouched.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
