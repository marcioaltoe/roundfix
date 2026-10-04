---
type: fix
status: promoted
created: 2026-10-03
spec: 0222-a-review-through-claude-that-keeps-its-verdict
reason: null
---

# The review through Claude blocks on a refused permission and loses its findings

## Symptom

Reported by fiscus and oraculum on 2026-10-02 with Roundfix 0.26.0, acpx 0.19.4, `pre_pr_review.provider: claude` and the `review` profile `claude / opus / high`:

```text
roundfix: review blocked: review transport anomaly: acpx exited with exit code 5 after parsed session/prompt result
```

- **The failure is a permission refusal, not a transport fault.** acpx exit 5 is `PERMISSION_DENIED`. The Claude reviewer opened a `Terminal` call. The read-only session answered the `request_permission` with `reject`, and the turn still ended with `stopReason: end_turn` and a verdict. `internal/agent/acpx_runner.go` turns any exit other than 130 after a parsed result into `TransportAnomaly`, and `classifyReviewCommandResult` blocks before it reads the answer.
- **The block hides real verdicts.** It blocked a clean `No findings.` on fiscus, and four rounds with real findings on oraculum (Spec 0057), which therefore got no finding ids and could not be disposed.
- **The prompt can overflow.** On oraculum, a ~10,000-line candidate (Spec 0058) produced a ~1.29 M-token request against a 1 M limit. The command printed only `agent/protocol error`; the real reason, `Prompt is too long`, was only in `pre-pr-review-answer.txt`. Spec 0212's byte bound (917,504 bytes of diff) does not account for the Spec context and whole files the Claude session attaches.

Codex does not show this, because its read-only sandbox never asks for permission.

## Expected

- A turn that delivered its verdict after the read-only session refused a permission is classified by its verdict, with finding ids. The block reason names the refusal.
- The review session for Claude denies tools that need permission up front, or the prompt forbids them.
- The review bounds what it sends by the provider's real context limit, and when a runtime reports `Prompt is too long`, the command prints that reason instead of a protocol error.

## Sources

Secondbrain `inbox/roundfix/2026-10-02-review-claude-bloqueia-quando-o-revisor-pede-terminal.md` (fiscus) and `inbox/roundfix/2026-10-02-review-claude-perde-achados-e-estoura-o-prompt-no-oraculum.md` (oraculum).
