---
status: pending
created_at: 2026-08-14
updated_at: 2026-08-14
kind: finding
---

# Preflight starves when the machine is busy (2026-08-14)

`profile proof failed … adapter error: context deadline exceeded` was already on
record from 2026-08-08 as tuple unavailability. This adds the missing variable:
it is not random intermittence. It appears when other Roundfix Runs are active on
the same machine and disappears when they finish. Minted from the Inbox Entry
`inbox/roundfix/2026-08-12-preflight-e-inanicao-por-runs-concorrentes.md` in the
Secondbrain, observed in `conexus` on 2026-08-11 and 2026-08-12.

## 1. The failure tracks load, and widens as load rises

- Symptom / evidence: with five Runs active across `fiscus`, `fluxus`,
  `roundfix` and two worktrees, load average at 23, `roundfix implement` failed
  in Preflight without creating a Run. The failure first affected only the `data`
  and `docs` categories; the next attempt affected `general`, `backend`,
  `frontend`, `qa` and `review`. With the machine idle at load average 4.9 and no
  concurrent Run, `roundfix profiles validate` passed all three configured tuples
  on the first attempt with nothing else changed. A third occasion, at
  intermediate load, hit only `qa`.
- Root cause: not established. The correlation is measured across three
  occasions; the mechanism inside the disposable Agent Session is not. A
  disposable session per tuple, each paying process startup while the host is
  already saturated, is the leading hypothesis and needs its own reproduction
  before anyone acts on it.
- Action / suggestion: two directions, both non-binding. A deadline that scales
  with observed host load rather than a fixed budget would stop reading
  saturation as unavailability. Proving the preferred Selection eagerly and every
  other tuple lazily would shrink the work Preflight does at the moment the host
  is least able to do it — which is what the preflight Spec in the current set
  already proposes for a different reason.

## 2. There is no fallback, because the Run does not exist yet

- Symptom / evidence: the failure happens in Preflight, before Run creation, so
  the Fallback Chain cannot activate — it activates only after a Run exists
  (ADR-0050). A machine under load therefore cannot start work at all, rather
  than starting it on a slower selection.
- Root cause: the chain's activation boundary is Run creation, and Preflight sits
  before it by design.
- Action / suggestion: worth settling whether a starved proof should be
  distinguishable from a refused one. A tuple that timed out is unknown, not
  unavailable, and ADR-0111 already draws that line for Verification.

## What worked — keep

The failure is loud and creates nothing: no Run, no Agent Session, no partial
state. An operator who reads the message and waits for the machine to clear gets
a clean first attempt.
