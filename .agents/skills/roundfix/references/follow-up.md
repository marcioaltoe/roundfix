## Follow-up Command

Use `roundfix follow-up --format json` when an agent needs the repository's
pending work in Follow-up Priority order. Use `roundfix follow-up --status
--format json` when it needs only work in progress, pending Specs, the next
release and maintainer blockers. Both commands are read-only and exit `0` for
an advisory report, including a report with no work.

The list command reads the sources `ai-memory-messages`,
`backlog`, `findings`, `repository-inbox`, `specs`, `runs`, `delivery-queue`,
`pull-requests` and `git`. Each source reports `checked`, `missing` or
`unavailable`, with a count or detail. The status command reads only
`runs`, `delivery-queue`, `pull-requests`, `git`, `specs` and the Release Plan;
it never calls the judge.

ai-memory messages are handed in by the caller with `--messages <path>`, never
read through an MCP. The file uses
`roundfix/follow-up-messages/v1`:

```json
{"schema":"roundfix/follow-up-messages/v1","messages":[{"id":"m-17","from":"billing","created_at":"2026-10-21","subject":"Fix export","text":"The export blocks billing.","type":"fix"}]}
```

Read the JSON `state`, `sources`, `judge` and `items` fields. `state` is
`work_found` or `nothing_to_do`; the latter means there are no items and no
work in progress. `sources` says what was checked or missed. `judge` names its
state, reason, model, calls and cost. Each item carries its rank, class, kind,
reference, type, date, title, requester, blocking project, basis, answers and
fallback reason. The four classes are `unblocking`, `health`, `feature` and
`other`. Without an answer, an item's declared type orders it, and nothing is
unblocking: `fix`, `perf` and Findings are `health`, `feat` is `feature`, and
the rest are `other`.

The judge sends nothing inside a Run, without a key, at the monthly ceiling,
with an unreadable Judge Log, or for a secret-pattern match. The command never
reads or writes the Secondbrain, calls ai-memory, or writes to the repository.

The Secondbrain inbox source was retired. The keys `secondbrain.root` and
`secondbrain.inbox_namespace` are deprecated and ignored with a warning. The
project name comes from the `origin` remote or the checkout directory.
