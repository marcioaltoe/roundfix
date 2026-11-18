# Model evidence for the 2026-11-17 snapshot

Read by the author on 2026-11-17 on the planning machine, read-only, and on the vendors' public pages. It records what the subscription runtimes offer; it is not a benchmark of this repository's Runs.

## Codex model listing

`~/.codex/models_cache.json`, written by Codex 0.169.0 and fetched `2026-11-17T14:55:41Z`, lists these models with `visibility: list`, in this order, with the reasoning efforts each accepts:

| Model | Codex description | Efforts |
| --- | --- | --- |
| `gpt-6.1-sol` | Latest workhorse model for coding and everyday work. | low, medium, high, xhigh, max, ultra |
| `gpt-6-astra` | Frontier intelligence for the most demanding work. | low, medium, high, xhigh, max, ultra |
| `gpt-6.1-terra` | GPT-6.1 Terra Codex model. | low, medium, high, xhigh, max, ultra |
| `gpt-6-sol` | Previous generation workhorse model. | low, medium, high, xhigh, max, ultra |
| `gpt-6-luna` | Fast and affordable model for easier tasks. | low, medium, high, xhigh, max |
| `gpt-5.6-sol` | Older generation workhorse model. | low, medium, high, xhigh, max, ultra |
| `gpt-5.6-terra` | Older balanced model for straightforward work. | low, medium, high, xhigh, max, ultra |
| `gpt-5.6-luna` | Older fast and efficient model. | low, medium, high, xhigh, max |

Two more entries, `gpt-reserve` and `codex-auto-review`, have `visibility: hide`. `gpt-6-luna` accepts `max` and not `ultra`; `gpt-6.1-terra` accepts `high`. The installed ACP adapter is `@agentclientprotocol/codex-acp` 2.7.0.

## Claude adapter

`@agentclientprotocol/claude-agent-acp` 0.107.3 (`dist/session-model.js`) advertises the alias `fable` for the newest Fable and keeps the pinned id `claude-fable-5-2` or `claude-fable-5-1` only for a session resumed on that version ("A resumed Fable 5.1 session keeps `claude-fable-5-1` too, because Fable 5.2 is newer").

### Claude selections proved on this machine

The recommended Claude tuples were proved against the installed adapter, offline from any vendor page: `roundfix profiles validate` (Roundfix 0.147.1, `@agentclientprotocol/claude-agent-acp` 0.107.3, Claude Code subscription) on 2026-11-17 at about 07:15 BRT (10:15 UTC), in the Supervisor session of cycle BF, printed:

| Runtime | Model | Effort | Result |
| --- | --- | --- | --- |
| `claude` | `opus` | `high` | `2. claude / opus / high — passed` |
| `claude` | `sonnet` | `high` | `5. claude / sonnet / high — passed` |

The adapter resolves the aliases `opus` and `sonnet` against its model list (`dist/session-config-options.js`, "A model takes an alias too, such as \"opus\" or \"sonnet\" for a full model id") and offers the effort levels Claude Code lists from `low` to `max` (`dist/session-effort.js`). These are the two Claude tuples the Recommended Profile uses: `opus`/`high` (preferred for `frontend`, fallback for the implementation categories) and `sonnet`/`high` (fallback for `docs` and `chore`).

## Vendor releases since the 2026-09-30 snapshot

| Model | Released | Source and cited passage |
| --- | --- | --- |
| GPT-6.1 Terra (`gpt-6.1-terra`) | 2026-10-26 | [OpenAI model release notes](https://help.openai.com/en/articles/9624314-model-release-notes): "GPT-6.1 Terra is rolling out in the OpenAI API as gpt-6.1-terra and in ChatGPT Work and Codex." [Introducing GPT-6.1 Terra](https://openai.com/index/introducing-gpt-6-1-terra/): "clearer, more natural writing"; "At maximum reasoning effort, GPT‑6.1 Terra solves 50.3% of tasks" on Terminal-Bench 4.0; "In ChatGPT Work and Codex, your plan goes further with GPT-6.1 Terra … substantially more than with GPT‑6.1 Sol or GPT‑6 Astra"; "For the most complex coding tasks, GPT-6.1 Sol and GPT-6 Astra remain the better choice." |
| Claude Haiku 5.5 (`haiku`) | 2026-10-07 | [Introducing Claude Haiku 5.5](https://www.anthropic.com/claude-haiku-5-5): "the cheapest, fastest, and most capable small model we've ever released"; "our first Haiku-class model to come with an adjustable effort setting." |
| Claude Fable 5.2 (`fable`) | 2026-11-09 | [Introducing Claude Fable 5.2 and Mythos 5.2](https://www.anthropic.com/claude-fable-and-mythos-5-2): "For most coding and knowledge tasks, we recommend that users start with Claude Opus 5.5, our strongest daily driver." |

`gpt-5.5` left Codex with ChatGPT sign-in on 2026-10-14, as `docs/references/model-selection.md` already records.

## What the table follows from

- `review` and `chore` move from `gpt-5.6-luna`, which Codex now calls older, to `gpt-6-luna`, its fast model, at `max`, the highest effort it accepts.
- `docs` takes `gpt-6.1-terra` at `high`: the writing improvement and the quota OpenAI cites fit documentation Tasks, while the coding categories keep `gpt-6.1-sol`.
- Fable stays out of Run categories: the vendor itself points daily coding at Opus 5.5.
- No selection uses `ultra`, by the maintainer's decision.
- Every selection runs on the `codex` or `claude` runtime, the subscriptions (ADR-0235).
