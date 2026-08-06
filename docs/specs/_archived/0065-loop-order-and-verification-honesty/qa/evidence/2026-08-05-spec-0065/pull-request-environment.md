# Pull Request environment and equivalent evidence

Status: `blocked (environment: no open Pull Request)`.

Cause: the QA prompt states that no Pull Request is open for the Spec target
branch and that Pull Request journeys are environment-blocked. The per-Run
branch is never pushed and has no Pull Request of its own.

Exact unblocking action: open a Pull Request from
`ma/0065-loop-order-and-verification-honesty`, then rerun this row read-only to
observe checks, approval, unresolved review state, and Merge-Ready evidence.
This QA gate must not create or mutate that Pull Request.

Equivalent observed/supervised evidence:

- the current built CLI and authored gate run successfully before a Pull
  Request exists;
- accepted ADR-0080 defines the typed-equivalent-evidence verdict rule;
- archived Spec 0078 report
  `qa-report-2026-08-05.md` closed `verdict: pass` with 18 terminal rows, 7
  passes, and 11 environment-blocked rows backed by equivalent evidence;
- its `equivalent-pr-evidence.md` records the public/supervised request,
  idempotency, Final Push ordering, refusal, round-cap, and Event Stream
  evidence and the same exact hosted unblocking action.

The Spec 0065 carriers incorrectly call that count nine rather than eleven;
F-002 records the documentation defect separately. The count error does not
erase the archived report's observed proof that the typed path can pass.
