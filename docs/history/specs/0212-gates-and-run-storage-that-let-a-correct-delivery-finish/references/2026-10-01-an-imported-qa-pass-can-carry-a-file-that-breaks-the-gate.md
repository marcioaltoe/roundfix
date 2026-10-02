---
type: fix
status: promoted
created: 2026-10-01
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
reason: null
---

# An imported QA pass can carry a file that breaks the repository gate

## Opportunity

Spec 0202 imports a failed QA pass's files into the next pass. In Spec 0203's delivery, the first QA Agent wrote evidence as Go source (`qa/evidence/ledger_replay_test.go`, `package cli`, unformatted, and depending on `internal/cli` symbols). The next pass imported it byte for byte. The QA precondition (`make verify-changed` → `fmt-check`) then failed on that imported file before any row ran, and every retry would import it again. The operator broke the loop by committing a compiling placeholder at that path, which made the import refuse (`path differs`), and kept the original as `.go.txt`.

## Value

A QA Agent's evidence must never be able to fail the repository gate of later passes or loop a delivery.

## Shape

The import refuses or skips evidence files that the repository gate compiles or formats (Go sources, at minimum), and records the reason. The qa-gate skill also tells the Agent to keep runnable evidence out of compiled paths or under a non-compiled extension.
