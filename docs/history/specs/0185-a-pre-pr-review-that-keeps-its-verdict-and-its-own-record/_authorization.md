---
status: approved
granted: 2026-09-29
action: keep Agent message boundaries so the pre-PR review classifies the reviewer's final message and a sealed prompt parses its final message, and keep one pre-PR review record and answer per checkout under the Artifact Directory
consuming: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
paths:
  - .agents/skills/roundfix/SKILL.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0185

On 2026-09-29 the maintainer approved adding this Spec to the second queue.
They answered a structured question, which named both frictions and their
Backlog Entries, with "Nova Spec na Onda 5". This Spec joins Specs 0183 and
0184 in the Onda 5 `roundfix deliver start`, which the maintainer approved the
same day with "Duas filas". The skill files ride the standing grant of
2026-09-18 for keeping the shipped skills true to the CLI. The set was measured
with `GovernedPath` on the authoring branch at `b13c96dc`, through a
`go test -overlay` probe that wrote nothing to the repository.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md` — the
  Roundfix skill describes how `roundfix review` classifies the reviewer's
  answer and keeps it (task_01), and where the record and the answer live
  (task_02). `.agents/skills/` is canonical and `skills/` its mirror.

## What is not governed

The following are ordinary:

- `internal/agent/stream.go`, `internal/agent/acp_stream.go`,
  `internal/agent/agent.go`, `internal/agent/acpx_runner.go`,
  `internal/agent/sealed.go` and the new `internal/agent/agent_messages.go`;
- `internal/cli/review.go`;
- the existing tests `internal/cli/review_test.go`,
  `internal/cli/review_archived_spec_test.go` and
  `internal/cli/review_disposition_test.go`;
- the new test files `internal/agent/agent_messages_test.go`,
  `internal/cli/review_final_message_test.go` and
  `internal/cli/review_record_checkout_test.go`;
- `docs/user-guide/commands.md` and `CONTEXT.md`.

`internal/cli/cli_test.go` and `docs/references/coverage-record.json` are
governed and are not touched. No help text changes, and no existing top-level
test is renamed or removed.

## Sanctioned regeneration

```yaml
command: make skills-sync
```

## Limits

- No change to the reviewer prompt, the Pre-PR Review Policy, the review
  profiles, Spec 0182's merge-base contract, the Delivery Queue's verdict path,
  the disposition ledger or archived Specs.
- The live Artifact Directories under `~/.roundfix` are never written by a Task
  or the gate. The old shared record is neither migrated nor deleted.
- No reviewer, provider or network call from a test, and no paid API use,
  release, tag, deployment or branch-policy exception.
- Verification stays Daemon-owned, and Task status stays Daemon-written.
