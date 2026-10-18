---
status: done
created_at: 2026-10-18
updated_at: 2026-10-18
absorbed_by: 0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born
---

# Delivery Queue — `secret-found` asks for a history rewrite that the retry of an archived item refuses (2026-10-18)

Triaged from the Secondbrain Inbox Entry
`inbox/roundfix/2026-10-18-secret-found-manda-reescrever-mas-o-retry-recusa-o-head-arquivado.md`
(origin `conexus`, captured manually on 2026-10-18) and checked against this
repository at `2ef4f080` (v0.89.0). The observation comes from the Conexus
repository (`useconexus/conexus`), Roundfix 0.87.0, Spec
`0001-usuario-verificado-com-segundo-fator`, whose Delivery Queue item was
archived on its branch and parked at the final push as `secret-found`.

## 1. The printed next action leads to a head the retry refuses

- Symptom / evidence: the push found five `credential-assignment` matches,
  all one test fixture: an assignment to `BETTER_AUTH_SECRET` of a value of 20
  or more characters with letters and digits. The printed next action was
  "remove the secret from every commit the blocker names on the item branch
  ..., then run roundfix deliver retry". Removing a value from earlier commits
  rewrites history and moves the head. After a `git filter-branch` over
  `6bcddf3..HEAD` in the item worktree, the retry refused with `archived item
  head "<new>" differs from candidate head "<old>"` and no side effects. The
  workaround that worked: `git reset --hard` to the QA Report commit before
  `docs: archive ...`, so the Spec was active again with every Task completed;
  a second rewrite to a letters-only value; and a retry, which moved
  `secret-found -> reviewing` but restarted the pre-PR review from round 1,
  because the Reviewer Lineage did not survive the rewrite.
- Root cause: proven in this repository. `ClassifyPark` prints the next action
  (`internal/delivery/park_class.go:65-67`). The retry of an archived item
  (`internal/delivery/engine.go:547-582`) accepts a moved head only when it
  descends from the newest candidate, or from the Run start when no candidate
  is recorded, and otherwise refuses with "retry Delivery Queue item %q:
  archived item head %q differs from candidate head %q" (`:571`). A history
  rewrite never descends from the commits it replaced, so the action the park
  prints cannot be followed by a retry that succeeds. ADR-0273 states the
  operator "removes the value from every commit the blocker names on the item
  branch before the retry" and does not say how the retry recognizes that
  work.
- Action / suggestion: for `secret-found` on an archived item, accept a
  rewritten head that keeps the final tree except at the reported lines, or
  print the path that works (return to the commit before the archive) in the
  next action.

## 2. The Secret Scan finds a fixture only at the final push

- Symptom / evidence: the Daemon committed the fixture in each Task without
  refusing it, and the problem appeared only after QA, review and archive. The
  first rewrite's replacement value, with letters and digits and more than 20
  characters, also matched the rule; the value pattern that Vortex uses in its
  fixtures today would be refused if delivered again.
- Root cause: proven in this repository. The Secret Scan runs at the QA Report
  commit (`scanQAReportPaths`, `internal/daemon/task_engine.go:3058`), at the
  Daemon's Final Push (`internal/daemon/daemon.go:343-350`) and at the
  Delivery Queue's push (`internal/delivery/github.go:179-188`). A Task's
  settlement (`executeTask`, `internal/daemon/task_engine.go:1145-1227`) and
  its commit (`commitTask`) run no scan. The `credential-assignment` pattern
  (`internal/secretscan/patterns.go:23`) matches a name ending in `API_KEY`,
  `TOKEN`, `SECRET` or `PASSWORD`, an assignment operator and a value of at
  least 20 characters, and `matchesKind` (`:34-43`) reports it only when the
  value holds a letter and a digit; nothing exempts a test file.
- Action / suggestion: run the Secret Scan at each Task's commit, as it runs
  at the QA Report commit, so a fixture is refused when it is born, and
  document for test authors which fixture values the rule does not match, or
  consider exempting evident fixture values such as ones containing `test` or
  `fixture`.

## What worked — keep

- The push refused before anything left the machine, and the park named every
  location without its value.
- The retry refusal changed nothing, so the item could still be recovered by
  hand.

## Addendum — 2026-10-18 — Adopted by Spec 0280

[0280-qa-evidence-keeps-its-images-and-a-secret-is-refused-where-it-is-born](../_prd.md) adopts this Finding. The source moves once into that Spec's
reference index; its lifecycle status records adoption, not implementation or
QA completion.
