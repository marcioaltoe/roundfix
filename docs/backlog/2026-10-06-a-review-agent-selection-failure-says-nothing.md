---
type: fix
status: open
created: 2026-10-06
spec: null
---

# A pre-PR review agent selection failure says nothing actionable

## Problem

In Pantheon on 2026-10-02 (Roundfix 0.26), `roundfix review --base main` with
the default `codex` provider (`codex / gpt-5.6-luna / max`) blocked twice in a
row at head `e1c3af16`:

```
reason: review runtime failure: Agent Selection failed for runtime "codex": agent/protocol error
```

The answer file was empty. In the same minute `roundfix doctor` reported the
codex adapter ok and `codex exec` answered. On 2026-10-05 (0.43.0) the same
review ran and returned findings with no configuration change. The error does
not say which protocol step failed or whether a retry is worthwhile, and the
maintainer skipped the review for Pantheon PR #45.

## Expected

The blocked review names the adapter's own error message and the protocol
step (initialize, session/new, set model, prompt), and a failure during agent
selection — before any prompt — is retried once automatically, as the Run's
own selection already does for transport failures.

## Sources

Secondbrain inbox, triaged 2026-10-06:
`inbox/roundfix/_triaged/2026-10-05-revisao-pre-pr-com-codex-falha-por-protocolo-de-forma-intermitente.md`
(Pantheon).
