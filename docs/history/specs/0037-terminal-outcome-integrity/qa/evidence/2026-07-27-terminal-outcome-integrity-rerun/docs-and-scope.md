# Documentation and scope sweep

Build: `ef6eb44ad8951112b1c3641bb7fd21793b440f95`

## Operator guidance

The built `roundfix --help`, `roundfix stop --help`, `roundfix settle --help`,
`roundfix runs list`, `roundfix events`, and `roundfix attach` paths executed
successfully where applicable. The Stop help matches the live Force Stop and
graceful-stop observables recorded in `live-cli-flows.md`.

The supported guides contain the Force Stop and Stop Request contracts:

```text
rtk grep -n 'stop --force\|Force Stop\|Stop Request'
  docs/user-guide/commands.md docs/user-guide/usage.md
```

The focused help/report regression passed with a writable task-local Go cache:

```text
rtk env GOCACHE=/private/tmp/roundfix-qa-0037-focused-gocache \
  go test ./internal/cli -run 'Test.*Stop.*(Help|Usage|Report)' -count=1
```

All cited glossary, ADR, finding, and command-reference files resolve. The
finding anchor source heading
`## 4. Cleanup noise appeared before the actionable failure` resolves.

## Non-Goals

- `roundfix pause --help`, `roundfix resume --help`, and
  `roundfix checkpoint --help` each exited 2 as unknown commands. No pause,
  resume, or checkpoint surface shipped.
- The live graceful-stop flow reported that only the Run Database Stop Request
  changed and that the in-flight Work Item remained active.
- The unprovable/reused-PID focused tests fail closed; no arbitrary-process
  target surface exists in Stop help.
- The live watch Stop Request changed no tracked path, commit, HEAD, fetch,
  push, or Review Source state.
- The changed-path and public-help sweep found no Spec 0038 Run Worktree
  classification contract in this Spec.
- Review Source evidence and notification content remain outside this Spec.
  The live watch fixture exercised existing CodeRabbit status content; the
  change only stopped later access after Stop Request observation.

## Skill contract

The canonical and generated Roundfix Skill files are byte-identical.
`roundfix skills check` and the complete repository gate both passed. The Skill
states proof-before-completion, registered-active cleanup, next-poll graceful
interruption, immutable terminal outcomes, and winner-only publication in the
same terms as the built CLI and supported guides.

Verdict: pass.
