---
type: feat
status: promoted
created: 2026-10-04
spec: 0229-a-jev-router-that-checks-credit-and-names-its-model
reason: null
---

# The Jev Router checks the credit it can spend and names the model it used

## Opportunity

The Jev Router measurement of 2026-10-04 (the addendum to Spec 0218's
measurement record) ended a routed Task at OpenRouter, not in Roundfix. The
gate read the month's spend (US$9.69 of the US$50 ceiling) and the key's
remaining limit (US$40.31) and sent the repair prompt. The OpenRouter
account behind the key held about US$5.11 (`total_credits` 210,
`total_usage` 204.885), and OpenRouter refused the request with HTTP 402,
"This request would exceed your available credits given your current
in-flight requests". Roundfix reported `agent/protocol error`, and no
fallback could run because work had begun. Both `router-prompt` Judge Log
lines recorded an empty `model`, so the ten-fold price spread between the
2026-10-02 and 2026-10-04 implementation prompts (US$6.81 and US$0.68 for the
same outcome) cannot be traced to a model.

The maintainer who runs routed Tasks, and any repository that opts into the
router, would gain three things: a refusal before the prompt when the
account cannot cover it, a named reason when OpenRouter refuses for credit,
and the routed model on each Judge Log line.

## Value

Hypothesis: with the account balance in the gate, a routed Task stops before
its first prompt (and the Fallback Chain takes it) instead of failing after
paying for a prompt; and with the model recorded, the next cost decision
(restricting the router's candidates or a price ceiling) rests on evidence.

## Shape

Non-binding, from the addendum's cost-reduction proposal items 1, 2 and 4:

- Read `GET /api/v1/credits` beside `GET /api/v1/key` and refuse below a
  floor, with a named reason such as `openrouter_credit_low`.
- Classify an OpenRouter credit refusal like `jev_ceiling_reached`.
- Record the model OpenRouter reports for each routed request.

Items 3 (a per-replay measurement key and an output-token cap), 5 (fixing
the 0194 task_04 replay pre-state) and 6 (keeping `docs` and `chore` on the
subscription default) are measurement practice or already the default and
stay outside this entry.
