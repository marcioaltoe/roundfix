---
type: feat
status: done
created: 2026-09-25
spec: null
reason: "Measured: the baseline and the routed comparison ran, routing was not adopted, and the maintainer closed the Jev front on 2026-09-29."
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

## Correction — 2026-09-29

A read-only recount of `~/.roundfix/jev-gateway/codex-baseline.jsonl` and `codex-routed.jsonl` corrects the 2026-09-28 result:

- With routing on (10:17–12:50Z), 194 of 598 requests failed with 502: 45 of 449 (10.0%, not 8%) before 12:30, and 149 of 149 afterwards. Jev chose `exec` in 494 of 543 decisions, with p50 498 ms and p95 2.1 s.
- After routing was switched off (12:50Z on 2026-09-28 to 05:28Z on 2026-09-29), 500 of 4,014 baseline requests still returned 5xx (12.5%), including windows where every request failed and no Jev call was made. The outage therefore belongs to the gateway or its upstream, not to routing, and the routed comparison is inconclusive rather than negative.
- The "0 failures" baseline holds only for 2026-09-25.
- The −36% figure the plan attributed to requests is, in the benchmark's primary source (secondbrain `raw/web/2026-09-18-github-vinilana-jev-gateway.md`), a reduction in time.

Disposition: `done`. Routing is not adopted, and Runs no longer go through the gateway: `deliver` and `implement` are launched without `CODEX_CONFIG`. Reopen only after the gateway has run a week in baseline with zero 5xx, and an A/B on the same Spec (at least three Runs per arm) shows at least a 20% reduction in requests per Task at the same pass rate. Sending context to TypeSafe also needs the maintainer's explicit authorization.

