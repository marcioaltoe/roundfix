---
type: fix
status: open
created: 2026-11-11
spec: null
reason: null
waiting_on: exact reproduction steps from the maintainer or Diego Coutinho (Conexus)
since: 2026-11-11
---

# The Conexus team saw a false Run status, a silent authorization wait and a black screen

## Symptom

Reported by conexus in ai-memory message `01a1c84e-8dfb-7483-966c-061a3c293853` of 2026-11-11 (10/11 22:52 BRT), from the Conexus weekly meeting of 2026-11-10. None was reproduced by the sender:

1. `--status` / `roundfix status` answered that no Run was active while a Run was in progress (the maintainer, minute 00:02:59 of the meeting).
2. A command sometimes stops waiting for an authorization and the calling agent never gets an answer, so it looks hung (Diego Coutinho, confirmed by the maintainer, minute 00:04:19).
3. A black screen when running Roundfix; the maintainer updated the tooling afterwards and does not know whether that fixed it (minute 00:02:30).

## Where

Unknown. Context from the sender: until 2026-11-10 Diego ran Roundfix 0.96.0 against a repository configured for 0.137.7, and the Conexus card for its Spec 0023 says the configuration was incompatible and Specs 0022/0023 ran outside Roundfix. The version-mismatch part is planned as Spec 0382 (an older binary refuses a newer repository configuration); the three symptoms may or may not be consequences of it.

## Expected

The status view reports every active Run; a command that needs an authorization either asks on a channel the calling agent sees or fails with a message naming what to authorize, never waits silently; Roundfix never leaves a blank terminal. Each symptom needs its exact steps (command, Roundfix version, harness, terminal) before a Spec can characterize it.

## Evidence

none yet — relayed report, not reproduced. The sender asked to request reproduction steps from the maintainer or Diego Coutinho.
