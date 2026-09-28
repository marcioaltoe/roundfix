---
type: feat
status: open
created: 2026-09-25
spec: null
reason: null
---

# Measure agent token use through jev-gateway before turning routing on

## Opportunity

Roundfix cannot say how many requests and tokens a Task costs per runtime. jev-gateway sits between the ACP agent and its provider; with `--routing off` it only counts, and with routing on it sends tool-choice decisions to Jev. A community benchmark reports large savings on GPT bug fixes and regressions on Claude models, so the gain is runtime-dependent.

## Value

A measured baseline per runtime and Task kind, then an evidence-based decision on routing. Routing sends tool schemas and context to TypeSafe, so it also needs a data-sharing decision.

## Shape

Run the codex runtime through the gateway in baseline mode without changing the maintainer's personal agent configuration, record counts per Run, and compare with routing on only after the deterministic-test work (D3) removes load flakes. Show the comparison to the maintainer before adopting routing.

Evidence: secondbrain `raw/roundfix/2026-09-25-sequencia-de-eficiencia.md` and `raw/web/2026-09-25-exa-roundfix-eficiencia-e-custo.json`; https://github.com/vinilana/jev-gateway.

## Result — 2026-09-28

Baseline measured over Specs 0170–0172 (1,303 gpt-5.6-sol requests, 0 failures, 97% of input cached, ~540 output tokens and 16.5 s per request). Routing on (`jev-1.13.0`, pinned) during Onda 2: before the outage 8% of requests failed and requests took 46.9 s on average; from 12:30 every request through the gateway failed with `ERR_HTTP2_STREAM_ERROR` (502), including passthrough ones, and Agent sessions died mid-Task (Spec 0173 lost two Runs). Direct codex worked at the same time. Routing was switched off and is not adopted; baseline mode stays available for measurement. Data: secondbrain `raw/roundfix/2026-09-25-jev-gateway-baseline-onda-1.md` and `~/.roundfix/jev-gateway/codex-routed.jsonl`.
