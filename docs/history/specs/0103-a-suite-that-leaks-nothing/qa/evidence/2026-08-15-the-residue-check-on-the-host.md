---
spec: 0103-a-suite-that-leaks-nothing
date: 2026-08-15
kind: equivalent-evidence
row: R07
---

# The residue check, built and run on the host

The gate's R07 was blocked because the Agent Session's sandbox denied the Go
build cache outside its workspace, so no fresh binary could be produced. Built
and run from the repository root on the operator's machine instead:

```text
$ go build -buildvcs=false -o /tmp/rf-r07 ./cmd/roundfix
$ /tmp/rf-r07 doctor; echo "exit=$?"
residue: ok (no process residue found)
exit=0
```

The check reports the empty case in words rather than as an empty table, and the
diagnostic's exit status is unchanged, which is what Core Feature 5 requires of
the reported-nothing path. The found, live-Run, and partial paths are proven by
`TestDoctorReportsProcessResidue` in the command's own package.
