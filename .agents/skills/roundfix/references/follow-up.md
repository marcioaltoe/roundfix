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
it never calls the judge. When the repository's Spec Root is absent, the
`specs` source is `checked` with count `0` (ADR-0348).

When the Release Plan needs a manual classification, the status `release.state`
is `manual_classification_required` and `release.proposed_version` is empty.
`release.reason` says how many commits since the base tag need a manual release
classification, and the `release_classification` item of `maintainer` names the
`roundfix release plan --from <base> --to HEAD --impact <none|patch|minor|major> --reason <text>`
command that classifies them.

ai-memory messages are handed in by the caller with `--messages <path>`, never
read through an MCP. The file uses
`roundfix/follow-up-messages/v1`. Each message needs a unique `id`, an RFC
3339 or `YYYY-MM-DD` `created_at`, and `text`; `from`, `subject` and `type`
are optional. The file is at most 1 MiB.

A message type is empty or one of `feat`, `feature`, `finding`, `fix`,
`information`, `other`, `perf`, `refactor`. An empty or omitted type is stored
as `message`; a listed type keeps its given name. A refused type's error
lists the accepted types and says that `type` can be omitted. Invalid messages
exit `2`, with empty standard output and the error followed by
`Run 'roundfix follow-up --help' for usage.` on standard error.

```json
{"schema":"roundfix/follow-up-messages/v1","messages":[{"id":"m-17","from":"billing","created_at":"2026-10-21","subject":"Fix export","text":"The export blocks billing.","type":"fix"}]}
```

Read the JSON `state`, `sources`, `judge` and `items` fields of the list
report. `state` is `work_found` or `nothing_to_do`; the latter means there are
no items and no work in progress. `sources` says what was checked or missed.

The list report carries these JSON fields: `schema`, `view`, `state`,
`repository`, `sources`, `judge`, `items`, `in_progress`.
The status report carries these JSON fields: `schema`, `view`, `state`,
`repository`, `sources`, `in_progress`, `pending_specs`, `release`, `blockers`,
`maintainer`.
Each source carries these JSON fields: `name`, `state`, `count`, `detail`.
The judge carries these JSON fields: `state`, `reason`, `model`, `calls`,
`cost_usd`.
Each item carries these JSON fields: `rank`, `priority`, `priority_name`,
`kind`, `ref`, `type`, `date`, `title`, `requested_by`, `blocks`, `basis`,
`answers`, `reason`.

The four classes are `unblocking`, `health`, `feature` and `other`. Without
an answer, an item's declared type orders it, and nothing is unblocking:
`finding`, `fix` and `perf`, and any item whose kind is `finding`, are
`health`; `feat` and `feature` are `feature`; `information`, `other` and
`refactor` are `other`. An omitted type (`message`) and any unlisted type are
also `other`.

The judge sends nothing inside a Run, without a key, at the monthly ceiling,
with an unreadable Judge Log, or for a secret-pattern match. The command never
reads or writes the Secondbrain, calls ai-memory, or writes to the repository.

The Secondbrain inbox source was retired. The keys `secondbrain.root` and
`secondbrain.inbox_namespace` are deprecated and ignored with a warning. The
project name comes from the `origin` remote or the checkout directory.
