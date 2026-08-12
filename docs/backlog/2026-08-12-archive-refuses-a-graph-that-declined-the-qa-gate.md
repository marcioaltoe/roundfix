---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-12
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# `archive` refuses a graph that declined the QA gate

## Symptom

Two parts of the contract disagree. The `write-tasks` skill declares
`qa: declined` as one of the **two** valid shapes of a post-contract graph:

> **Decline the gate.** Add `qa: declined` and a non-empty `qa_reason` to
> `_tasks.md` frontmatter. Do not emit a `qa` Task when the gate is declined.

`roundfix archive` does not know that shape. Against a Spec whose graph declines
the gate, it refuses at preflight:

```text
Preflight failed

Reason:
  no passing QA verdict: no QA Report found for Spec directory
  ".../docs/specs/0023-convergence-rejection-observability";
  run the qa-gate workflow to produce one
```

Measured on 2026-08-10 with roundfix 0.4.0 (`0c7d2d3c`), in `fluxus`, against
Spec `0023-convergence-rejection-observability`, whose Run reached `Clean` with
4 of 4 Tasks completed.

The decline was not laziness: that Spec declares only a backend surface, and its
whole correctness claim — a deterministic refusal writes exactly one quarantine
row, stops being retried, and stops downgrading the Convergence Report — is
established by an integration proof against a real PostgreSQL. No user-reachable
flow exercises it; a user walk would rehearse the Spec's own premise.

The practical effect is that a legitimately declined Spec **cannot be archived by
the tool**. Either it is archived by hand, bypassing the preflight that exists to
protect archival, or a QA gate runs that the Spec already argued proves nothing —
which is exactly the waste the decline option exists to avoid.

## Where

`internal/spec/archive.go`, which reads a QA Report unconditionally before any
other check (`Archive`, lines 60–79). `internal/spec/spec.go` already parses
`qa: declined` with its required `qa_reason`; `archive` never consults it.

## Expected

Preflight accepts `qa: declined` with a non-empty `qa_reason` as an equivalent
verdict, recording the reason in the archival artifact instead of requiring a QA
Report. The contract already requires the reason to exist and to be specific; it
is the evidence.

Worth checking in the same work whether `implement-spec` and the Supervisor treat
the declined shape consistently, or whether `archive` is the only point that does
not know it. Sibling defect, same command:
`docs/findings/2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md`
finding 1, where a declared-only `partial` cannot reach archive either.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-10-archive-exige-qa-report-mesmo-com-qa-declined.md` in
the Secondbrain, captured from a `fluxus` session.
