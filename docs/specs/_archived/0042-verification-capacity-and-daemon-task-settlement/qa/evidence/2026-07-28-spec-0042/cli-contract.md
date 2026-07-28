# CLI contract evidence

The built `roundfix implement --help` exposes no Task/Verification capacity
flag and retains the existing exit/output surface. The built successful flow
kept Task results and the final Clean summary on stdout; Run identity,
waiting, Verification progress, diagnostics, and notification warnings stayed
outside requested stdout.

`roundfix events` emitted only `roundfix-events/v1` JSONL records for the
requested categories and replayed a terminal Run with exit 0.

Finding: root help and `roundfix attach --help` advertise:

```text
roundfix attach [<run-id>] [--no-input]
```

The exact advertised command:

```text
roundfix attach run_20260728T140033Z_867f7808dae52080 --no-input
```

exited 2 with:

```text
roundfix attach failed: unexpected argument "--no-input"
```

Both `roundfix attach <run-id>` and
`roundfix attach --no-input <run-id>` replayed the same Run successfully.
The failure is therefore the advertised stdlib-flag ordering, not Run data.
