---
status: done
absorbed_by: 2026-08-06-rollup-agent-selection-and-execution-environments.md
created_at: 2026-08-07
updated_at: 2026-08-26
kind: finding
---

# A sixty-four-value bound locks out the opencode runtime (2026-08-07)

`opencode` is one of Roundfix's three supported runtimes and no profile can use
it. Its adapter advertises 431 models; Roundfix's capability validator caps a
select option at 64 values and refuses the evidence.

## Reproduction

```
roundfix profiles configure --scope project --file <opencode profile> --dry-run
  → exit 2
  → classification: capability_evidence_invalid
  → adapter error: capability evidence invalid: contradictory_response,
    missing_model_state, too_many_option_values
```

Measured against the bounds in `internal/agent/selection_capabilities.go`:

| Evidence | Observed | Bound | |
| --- | ---: | ---: | --- |
| `session/new` result size | 37,275 bytes | 65,536 | within |
| `model` option values | **431** | **64** | exceeds |
| option count | 2 | 32 | within |

Only the value count is exceeded, and it is exceeded by nearly seven times. The
response is well under the byte bound, so truncation is not the cause.

## A second, smaller mismatch

OpenCode 1.18.15 advertises exactly two config options, `mode` and `model`. It
advertises no `effort` option, unlike both first-party adapters. Roundfix decides
whether to parse model identifiers as opaque wholes by testing whether a
reasoning option exists, and a profile for this runtime still has to carry a
`reasoning_effort` value that the adapter cannot honour. That is the likely
source of the `missing_model_state` and `contradictory_response` companions to
the value-count refusal, though only the value count is directly measured here.

## Why it matters

The runtime is accepted by configuration validation
(`internal/config/profiles.go` lists `codex, claude, opencode`) and rejected by
capability proof, so the support is nominal. What it locks out is
substantial: OpenCode reaches DeepSeek, MiMo, GLM, Kimi, Grok, and Gemini, and
its OpenCode Zen provider serves several of them at no cost. Two were verified
working through the adapter on 2026-08-07:

```
opencode/deepseek-v4-flash-free  → replied "OK"
opencode/mimo-v2.5-free          → replied "OK"
```

So the models answer, the runtime is listed as supported, and the profile cannot
be written.

The bounds themselves are sound in intent — they keep a hostile or broken
adapter from flooding the validator. They were sized for adapters advertising a
handful of models: codex advertises 7, claude 5.

## Route

Not fixed here. The bound cannot simply be raised to 431 without deciding what
it is protecting against, and a fixed larger number only moves the wall. Worth
weighing: proving only that the *requested* model appears in the advertised set,
rather than admitting and validating the whole list, which makes the check
independent of catalogue size. The absent `effort` option needs its own answer —
either a profile for such a runtime omits reasoning effort, or Roundfix records
that the adapter cannot honour it.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
