---
type: fix
status: promoted
created: 2026-10-26
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
reason: null
---

# The unreaped-group test fails under the full suite on macOS

## Symptom

`TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` (`internal/cli`)
fails now and then when the whole package runs on macOS, and passes when it
runs alone:

```text
--- FAIL: TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded (0.00s)
    implement_detach_teardown_test.go:256: probe unreaped fixture group 6892: got <nil>, want EPERM
FAIL roundfix/internal/cli 219.605s
```

A QA gate that runs the package fails on a test its Spec never touched, and
the Run is retried.

## Where

`internal/cli/implement_detach_teardown_test.go`
`TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded`: after closing the
fixture `cat`'s input it polls until `detachFixtureGroupLiveMembers` reports
no live member, then expects `kill(-pgid, 0)` to fail with `EPERM` on Darwin.
`internal/cli/detach_fixture_group_darwin_test.go`
`detachFixtureGroupLiveMembers` leaves out both a zombie (`SZOMB`) and a
process that is still exiting (`P_WEXIT` set, not yet `SZOMB`).

## Expected

The test asserts `EPERM` only once its member is the exited, unreaped member
it names, and passes under any load; the assertion itself, and the other
three assertions of the test, stay as they are.

## Evidence

- Spec 0317 QA, Run `run_20261025T143420Z_c1664469a7b4b38b`, 2026-10-25 at
  about 12:10 BRT, audited head `f463cdfec7b7e2d1d9fd51c1385fb34047ad4d75`:
  `go test -timeout 20m -count=1 ./internal/baseline ./skills ./internal/config ./internal/cli`
  exited 1 with the failure above; the same test run alone on the audited
  head and on the Delivery Base exited 0 (QA finding F-01 of that Run,
  "Root cause is unresolved").
- Operator queue record of 2026-10-25 12:10 BRT and triage item 12 of the
  end of cycle AP: "sob a suíte completa no macOS o probe de sinal do grupo
  devolve nil em vez de EPERM (0317 QA 25/10). Candidato à auditoria de
  testes."
