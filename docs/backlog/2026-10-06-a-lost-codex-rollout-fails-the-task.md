---
type: fix
status: open
created: 2026-10-06
spec: null
---

# A lost Codex rollout fails the Task and spends a retry

## Problem

In Oraculum's Delivery Queue (Roundfix 0.44–0.46, `@agentclientprotocol/codex-acp`
2.1.1), between 2026-10-05 and 2026-10-06, `implement` Runs failed in the middle
of the agent's turn with:

```
internal error -32603: no rollout found for thread id <uuid>
```

- **Lost rollouts.** Three distinct threads failed
  (`01a10ec8-2071-7912-b2ed-645c09812e6a`,
  `01a110db-9b73-7073-865d-7f41e10c32da`,
  `01a110e3-5f35-70c2-bdd4-9a833b8d49a7`). None of them has a file in
  `~/.codex/sessions/`, while 47 other sessions of the same day do. The session
  is created, the rollout is never persisted, and the adapter's resume of the
  thread fails.
- **The QA case costs the most.** The QA agent dies before writing its
  report, which stays `verdict: pending`, and the item parks `run-unresolved`.
  Every attempt spends one of the queue's retries.
- **Manual workaround.** The same QA re-run by hand with
  `roundfix implement --spec <spec> --agent claude --model opus
  --reasoning-effort high --detach` reached a verdict (Oraculum 0059, 0062,
  0065).

## Expected

1. **Classify it as infrastructure.** `-32603 no rollout found` (and
   adapter-internal errors of the same class) is a runtime infrastructure
   failure. It takes the infrastructure path: no repair, no retry spent, and
   the Task resumes in a new session.
2. **Fall back when the runtime breaks.** When the profile's runtime fails
   that way before the Task's first handoff, or for a QA Task before its
   report, the profile's Fallback Chain takes the Task. ADR-0114 today allows
   the fallback only before Agent work began; decide whether a lost rollout
   counts as "no work" for QA.
3. **Find why the rollout is not written.** Candidates are a `CODEX_HOME` per
   Run, concurrent sessions, or disk. Measure with the Run's ACP stream and
   `~/.codex/sessions`, and report it upstream if it is an adapter defect.

## Sources

Secondbrain inbox, triaged 2026-10-06:
`inbox/roundfix/_triaged/2026-10-06-codex-acp-perde-o-rollout-e-a-task-falha-com-32603.md`
(Oraculum).
