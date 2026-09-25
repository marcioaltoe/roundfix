---
status: done
created_at: 2026-09-16
updated_at: 2026-09-25
spec: 0172-a-qa-gate-that-tells-the-truth
---

# Run Events — An unobserved Verification is published without its classification (2026-09-16)

This came up while authoring Spec 0139, which has the Daemon run the QA gate's
repository Verification. The Spec wanted to promise that an outcome the runner
could not observe reaches the Run Event Stream as an unknown result. Reading the
publisher and the projection showed it cannot, for any Verification, so the Spec
promises nothing about event shape and the gap is recorded here.

## 1. The publisher omits the classification the projection requires

- Symptom / evidence:
  - `publishUnknownFailure` in `internal/daemon/engine.go` writes `error` and an
    optional `diagnostic_path` into the payload, then publishes the verdict with
    empty failure metadata.
  - `failureMetadata` returns an empty classification unless the request carries
    one or the failure is temporary, so the event has no `classification` field.
  - `internal/runevent/stream.go` projects an event as unknown only when
    `classification` is `verification_unknown`, and then requires `command`,
    `reason` and `diagnostic_path`.
  - A Supervisor following the stream therefore reads an unobserved Verification
    as an ordinary failure: "the work is wrong" instead of "we did not find out".
- Root cause: the unknown outcome carries its cause in prose fields rather than
  in the typed metadata the projection reads. ADR-0111 requires the distinction
  to reach the record; the Task path's pre-work probe sets the classification by
  hand, and the shared attempt publisher does not.
- Action / suggestion:
  - Route to a Spec that makes the shared publisher emit
    `classification: verification_unknown` with the required `command`, `reason`
    and `diagnostic_path` for an unobserved outcome, and projects it.
  - Keep a command verdict unclassified, so an ordinary failure stays an ordinary
    failure.
  - Spec 0139 does not change the publisher, and its QA step therefore claims no
    projected event shape.

## Addendum — 2026-09-25 — Revalidated in triage

Revalidated against main 7a9b6ec6: still holds: publishUnknownFailure publishes with empty verificationFailureMetadata (internal/daemon/engine.go). Ranked in the 2026-09-25 triage priority list.

## Addendum — 2026-09-25 — Implementation owner

[0172-a-qa-gate-that-tells-the-truth](../_prd.md) adopts this Finding: its
shared attempt publisher will emit `classification: verification_unknown` with
`command`, `reason` and `diagnostic_path` for an unobserved outcome, and a
command verdict stays unclassified. The source moves once into that Spec's
reference index. Its lifecycle status records adoption, not implementation or
QA completion.
