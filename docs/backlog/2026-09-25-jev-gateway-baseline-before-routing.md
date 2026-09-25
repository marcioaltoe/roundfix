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
