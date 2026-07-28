# Config and race evidence

The built public `roundfix init --scope project` command generated:

```yaml
worktree:
  concurrency: 2

verification:
  concurrency: 1
```

The generated comments state that Verification Capacity is independent from
Task Capacity. The existing review command remained
`defaults.verification: make verify`.

The built public Implement Command rejected both
`verification.concurrency: 0` and `-1` with:

```text
verification.concurrency must be greater than 0
```

Both failures stated that no Run or Agent side effect occurred. A public
machine-wide `runs list` from the same Roundfix Home still contained only the
one earlier Clean Run.

Capacity 2 overlap, default/precedence, non-integer and unknown-key validation,
waiting-before-started, repeated fairness, and permit restoration passed the
current focused and full suites. The exact full
`go test -race ./... -count=1` passed with process-inspection access.
