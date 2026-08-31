---
status: completed
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

## Result

Implemented the structured Auditing Binary identity beside `VersionLine` and
made `VersionLine` render that shared identity. Released binaries render their
version when commit and build time are empty; stamped builds include the
trimmed commit and build time.

Added the `current`, `stale`, and `unknown` staleness states plus an ancestry
input state. `CompareToTree` uses commit ancestry for stamped builds and the
tree's declared version for unstamped builds. It returns the answering signal
as a reason and keeps missing ancestry or version evidence `unknown`.

Acceptance evidence:

- Identity assembly: `TestAuditorAssemblesBuildIdentity` covers released and
  stamped binaries, including the complete readable release identity.
- Two comparison signals: `TestCompareToTree` covers ancestry answers for
  stamped builds and declared-version answers for unstamped builds.
- Three states and absence handling: the same table covers `current`, `stale`,
  and `unknown`, including a stamped build with no ancestry answer and an
  unstamped build whose tree declares no version. The stamped unknown case
  supplies an equal tree version to prove missing ancestry cannot silently
  become `current` through the release fallback.
- Reasons: every comparison row asserts the exact reason returned with its
  state and identifies either commit ancestry, declared tree version, or the
  missing signal.

Focused checks:

- Red starting signal: `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test ./internal/app -run '^(TestAuditorAssemblesBuildIdentity|TestCompareToTree)$' -count=1`
  failed to compile before implementation with the new API symbols undefined.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test ./internal/app -run '^(TestAuditorAssemblesBuildIdentity|TestCompareToTree|TestVersionLineIncludesBuildIdentityWhenStamped)$' -count=1 -v`
  passed every named subtest after implementation.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go vet ./internal/app`
  passed.
- `rtk make verify-incremental` passed outside the filesystem/process sandbox.
  The sandboxed attempt had first reached `internal/app` successfully, then
  failed two unrelated `internal/cli` force-stop tests because process-table
  reads were denied; the permitted rerun passed those tests and the complete
  incremental target.

The Task's declared Verification remains unrun for the Daemon. Git ancestry
probing and QA Report serialization remain in the later Task slices.

## Carry-forward provenance

- Source Run: `run_20260830T161359Z_31aaee7e42ecc4e4`
- Source commit: `a05138c24a533dd16c674ce2fae0de11233b4f16`
