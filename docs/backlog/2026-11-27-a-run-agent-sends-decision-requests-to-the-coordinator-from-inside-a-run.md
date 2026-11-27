---
type: fix
status: open
created: 2026-11-27
spec: null
reason: null
---

# A Run Agent sends decision requests to the coordinator from inside a Run

## Symptom

Twice, the Agent of a Delivery Queue Run sent a cross-project decision request (`PEDIDO DE DECISÃO`) through ai-memory to the coordinating session, from inside its Task, instead of stopping with the question for its Supervisor:

- Spec 0421, task_03, 2026-11-26 00:45 BRT, request `01a215f1`: whether to accept a `roundfix baseline update` version bump in `docs/agents/setup-context.json`.
- Spec 0417, task_01, 2026-11-27 ~14:15 BRT, requests `01a21e01-8ad4` (touch `internal/baseline/project_decision_render.go`, a file that is not a Governed Path and needed no grant) and `01a21e07-fb65` (keep an inherited line in `docs/references/coverage-record.json`).

The answers came back to the Supervisor's inbox, not to the Run, which kept going without them. The first request asked for authority the plan did not need; the second changed a QA criterion, which the Supervisor could apply only at a Run boundary (Spec 0417, correction C1). The coordinator had to relay and the Supervisor had to reconcile, while the Agent's own session never received the answer.

## Expected

A Run's Agent never sends a cross-project message: its Task settles, fails or stops with the question recorded for the Supervisor, as the Run contract already says for a person's decision. The Baseline's guidance for a Run session, or the Run's own prompt, says so in one sentence, and the Agent can tell whether a path is governed (`speccheck.GovernedPath`) before asking for a grant.

## Source

Observed by the Supervisor of Spec 0417 on 2026-11-27, with the coordinator's relays `01a21e04-2eef` and `01a21e0b-4631`, and in Spec 0421's cycle on 2026-11-26. Nature: measured.
