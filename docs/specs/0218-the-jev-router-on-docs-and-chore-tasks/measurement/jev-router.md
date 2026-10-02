# Jev Router measurement

Date: 2026-10-02

## Measured

No replay ran. The scoped `ROUNDFIX_OPENROUTER_API_KEY` was present and
`NODE_OPTIONS` was unset, but the required scratch pre-state could not pass
Roundfix's committed-Spec graph preflight without creating a disposable commit.
The Task explicitly prohibits commits, so no prompt was sent and no Run data,
Judge Log cost, fallback, token count, or Verification repair count exists.

| Task | Category | Selection | Wall time | Prompts | Verification repairs | Outcome | Fallback and reason | Tokens | Router cost |
| --- | --- | --- | --- | ---: | ---: | --- | --- | --- | ---: |
| none | — | — | — | 0 | 0 | preflight stopped | no replay; committed scratch graph required | — | — |

## Reading

There is no cost per settled Task on the router because no Task settled. Tool
calls through `roundfix-openrouter/typesafe/jev-router` were not exercised, so
the sample says nothing about tool-call compatibility, cost, latency, fallback
behavior, or Verification repairs. It cannot support a router proposal for
`docs` or `chore`, and no follow-up Backlog Entry was minted.

The measurement is blocked by the protocol's scratch-state construction: the
CLI reads the Spec graph from `HEAD`, while the required Task-only graph exists
only after the scratch clone is rewritten. Creating the needed ephemeral commit
would violate this Task's explicit no-commit invariant.
