---
type: feat
status: open
created: 2026-10-05
spec: null
---

# A judge-assigned model tier per Task

## Opportunity

The Jev Router (Spec 0218) chose a model for every request inside a Task, and
through OpenRouter it chose `openai/gpt-6-astra`, `openai/gpt-6.1-sol` and
`anthropic/claude-opus-5.5`, which the maintainer's subscriptions already
cover. On 2026-10-05 the maintainer set the rule that OpenAI and Anthropic
models run only through the Codex and Claude subscriptions, and the router was
retired (Spec 0230, ADR-0235). The maintainer found the idea of choosing a
model dynamically worth keeping in another form, for later analysis:

- **Choose per Task, not per request.** A subscription is reachable only
  through its own agent (codex-acp, claude-agent-acp), so one agent session
  cannot mix subscription and OpenRouter requests. The choice has to be made
  before a Task is dispatched.
- **The judge assigns the tier when the Tasks are written.** `write-tasks`
  asks the Jev judge for a `model_tier` (`light`, `standard`, `heavy`) for each
  Task and records it in the Task's front matter, where the maintainer reviews
  and can change it in the planning Pull Request. Dispatch only maps the tier
  to a profile candidate, with no Jev call at run time:
  - `light`: an open model through OpenRouter (DeepSeek first, from a
    configurable allow-list);
  - `standard`: `codex / gpt-5.6-luna`;
  - `heavy`: `codex / gpt-6.1-sol` or `claude / opus`.
- **The graph is unchanged.** The tier is metadata. Dependencies and waves are
  written as today. Without Jev the tier defaults to `standard`, which is
  today's behavior. QA Tasks and Tasks that touch Governed Paths always stay on
  the subscriptions. A `light` Task that fails Verification escalates one tier
  once.
- **Separate OpenRouter keys per stage.** The judge and the implementation use
  their own keys (for example `ROUNDFIX_OPENROUTER_JUDGE_API_KEY` and
  `ROUNDFIX_OPENROUTER_IMPLEMENT_API_KEY`), each with its own monthly limit, so
  the OpenRouter activity export shows the cost of each stage. The maintainer
  would create the keys.
- **Gate kept.** Any OpenRouter model id under `openai/` or `anthropic/` stays
  refused (ADR-0235).

## Value

Offload simple Tasks from the subscription quota at very low cost, with a
decision the maintainer can see and override before anything runs.

Measured inputs (OpenRouter export `docs/_inbox/openrouter_activity_2026-10-05.csv`
and the archived Task files, 2026-10-05):

- Effective prices per million tokens, cache included:
  `deepseek-v4.1-flash` US$0.05, `typesafe/jev-1.13` US$0.04,
  `gpt-6.1-sol` US$0.55, `gpt-6-astra` US$2.17, `claude-opus-5.5` US$4.23.
- A Task in the measured Runs used 1 to 6 million tokens, so a `light` Task on
  DeepSeek would cost about US$0.05 to US$0.30. The judge's tier call costs
  about US$0.0001 per Task.
- Of 1,247 archived Tasks, about 230 (18%) are natural `light` candidates:
  docs low 91, backend low 96, test low 26, chore low 18. Because they use
  fewer tokens than heavy Tasks, the subscription quota saved would be roughly
  10 to 15%.

## Risks

- DeepSeek has never run a whole Roundfix Task on its own. A failed `light`
  Task adds repair attempts and an escalation, which costs more time than
  sending it to Codex directly. The idea pays off only if `light` Tasks pass
  Verification on the first attempt in the large majority of cases (about 80%
  or more).
- Documentation Tasks in this repository are precise (required phrases, skill
  versions, mirrors), which raises the failure risk for a weaker model.
- Routed OpenCode Runs took 6 to 10 minutes against 3 to 10 minutes on Codex.
  There is no wall-time data for DeepSeek alone.

## Shape

Non-binding. Measure before building:

1. Replay six archived `light` Tasks (two docs low, two backend low, one test
   low, one chore) directly on `opencode` with DeepSeek, no router, and compare
   first-attempt Verification, wall time and cost with their original Codex
   Runs. Estimated cost under US$2.
2. Ask the Jev judge for the tier of about 50 archived Tasks offline and
   compare with the repairs each Task actually needed. Estimated cost about
   US$0.01.
3. Proceed only if first-attempt Verification is at least 80% and wall time at
   most 1.5 times Codex. Otherwise keep only the separate keys.

Any measurement sends only this repository's code and Spec artifacts, uses
only open models and the Jev judge, and stops at US$5.
