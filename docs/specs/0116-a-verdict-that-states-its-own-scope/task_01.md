---
status: pending
type: backend
---

# Task: The auditing binary and its staleness

A QA Report records the commit it audited and nothing about the Roundfix that
produced the verdict, so a finding raised by a stale auditor is
indistinguishable from a finding about the code. This Task builds the identity
and the comparison; later Tasks record them.

## Work

- Assemble the running binary's identity from the existing build stamp beside
  the code that already composes `--version`: the version, the build commit,
  and the build time. Release builds leave the commit and time empty by design,
  so the identity must be complete and readable in that case too.
- Answer whether the binary predates the tree it audits, from two signals
  because neither is always available: commit ancestry when the build is
  stamped, and the version the audited tree declares when it is not.
- Return three states, never two. `unknown` is what a reader is shown when
  neither signal answers, and it must be impossible to reach `current` by
  absence of evidence.
- Return the reason alongside the state, so a later report can say which signal
  answered rather than leaving a reader to guess.
- Cover all three states table-driven, including both ways `unknown` arises: no
  build stamp, and no declared tree version.

## References

- `_prd.md` → Goal 3, Core Features 3 and 4
- `_techspec.md` → Build Order 1; Interfaces: `AuditingBinary`, `Staleness`,
  `CompareToTree`
- ADR-0135 makes an absent answer a reported state rather than an empty one

## Verification
- `grep -q "AuditingBinary" internal/app/version.go && grep -q "StalenessUnknown" internal/app/version.go && grep -q "TestCompareToTree" internal/app/version_test.go && go test -count=1 ./internal/app`
