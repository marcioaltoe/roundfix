---
type: feat
status: done
created: 2026-09-25
spec: null
reason: "Measured: the baseline and the routed comparison ran, routing was not adopted, and the maintainer closed the Jev front on 2026-09-29."
---

# Measure agent token use through jev-gateway before turning routing on

Roundfix cannot say how many requests and tokens a Task costs per runtime. jev-gateway sits between the ACP agent and its provider; with `--routing off` it only counts, and with routing on it sends tool-choice decisions to Jev. A community benchmark reports large savings on GPT bug fixes and regressions on Claude models, so the gain is runtime-dependent.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-09-25-jev-gateway-baseline-before-routing.md`.
