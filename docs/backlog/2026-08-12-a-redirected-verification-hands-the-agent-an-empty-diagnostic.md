---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-12
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A redirected Verification hands the Agent an empty diagnostic

## Symptom

The Verification command pattern adopted widely in this repository

```sh
go test ./pkg -run 'X' -v > /tmp/log 2>&1 && grep -q '^--- PASS: .*X' /tmp/log
```

was written to fix a real defect: `go test … | grep -q` exits with the status of
`grep` and hides the test's failure. The redirection fixes that and creates
another.

Run `run_20260808T153649Z_78746d4b80d08fc7`, task_03 of Spec 0084: the
Verification failed and the Daemon recorded the diagnostic in
`verification/batch-005-attempt-2.log` — **0 bytes**. Attempt 1 as well: 0 bytes.
All output went to the redirected file, and the `grep -q` that failed prints
nothing, so the Daemon captured empty stdout and stderr. The Verification
Feedback returned to the Agent Session carried no cause at all — only the command
and the exit status. The real cause sat in the redirected file, out of the
Daemon's reach:

```text
testing: warning: no tests to run
PASS
ok  	roundfix/internal/cli	(cached) [no tests to run]
```

The Daemon returns **only** the Verification's diagnostics to the Agent Session,
and only on failure. With none, the retry happens blind: in task_03 the Agent
spent its second attempt rewriting its own task file with a diagnostic it had
deduced, which is reasonable behaviour given zero information.

The pattern was present in at least 15 Verification commands written on
2026-08-07 across Specs 0080, 0081, 0082 and 0083, and in the ten Tasks of 0084.

## Where

`internal/daemon` — Verification execution and the Verification Feedback prompt —
and the authoring guidance that recommends the redirection form.

## Expected

A Verification that fails with empty stdout and stderr is reported as exactly
that, so the Agent knows to look for a log rather than deduce a cause. The
portable authoring form that keeps the status honest and still emits is worth
recording beside the pattern:

```sh
go test ./pkg -run 'X' -v > /tmp/log 2>&1; s=$?; grep -q '^--- PASS: .*X' /tmp/log || { cat /tmp/log; exit 1; }; exit $s
```

Worth settling in the same work: a cheap detector that refuses a Verification
command whose failure path cannot emit anything.

## Evidence

Minted from the Inbox Entry
`inbox/roundfix/2026-08-08-verification-que-redireciona-esconde-o-diagnostico-do-agente.md`
in the Secondbrain. Related:
`docs/backlog/2026-08-09-a-verification-command-passes-only-by-exiting-zero.md`,
which reaches the same empty-diagnostic behaviour by a different root cause.
