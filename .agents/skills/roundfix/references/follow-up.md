## Follow-up Command

```bash
roundfix follow-up [--status] [--all] [--messages <path>] [--format <text|json>]
```

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

`repository-inbox` reads every file in `docs/_inbox/` except a file or folder
whose name starts with a dot, such as `.gitkeep`, and everything inside such
a folder.

When the Release Plan needs a manual classification, the status `release.state`
is `manual_classification_required` and `release.proposed_version` is empty.
`release.reason` says how many commits since the base tag need a manual release
classification, and the `release_classification` item of `maintainer` names the
`roundfix release plan --from <base> --to HEAD --impact <none|patch|minor|major> --reason <text>`
command that classifies them.

ai-memory messages are handed in by the caller with `--messages <path>`, never
read through an MCP. The file uses
`roundfix/follow-up-messages/v1`. Each message needs a unique `id`, an RFC
3339 or `YYYY-MM-DD` `created_at`, and `text`; `from`, `subject`, `type` and
`blocked_work` are optional. The file is at most 1 MiB.

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
Each source carries these JSON fields: `name`, `state`, `count`, `detail`,
`waiting`.
The judge carries these JSON fields: `state`, `reason`, `model`, `calls`,
`cost_usd`.
Each item carries these JSON fields: `rank`, `priority`, `priority_name`,
`kind`, `ref`, `type`, `date`, `title`, `requested_by`, `blocks`, `blocked_work`,
`basis`, `answers`, `reason`, `waiting_on`, `since`, `waiting_error`.

## Declared blocked work

A message entry may carry `blocked_work`: one line, in the sender's own words,
naming the implementation, Spec or other work the sender says cannot start or
continue until this project answers. The session that writes the messages file
fills it from what the message says; the command never reads a message body
to find it. The value is trimmed, and an empty or blank value is the same as
no field. A value with a line break is refused with
`message "<id>": blocked_work must be one line`, and a value on an entry
without `from` with
`message "<id>": blocked_work needs from, the project whose work is blocked`;
both exit `2`.

An entry with `blocked_work` ranks in `unblocking` without the judge and is
never sent to it: priority `1`, `basis` `declared`, `blocks` its `from`,
`blocked_work` its value, and null `answers` and `reason`. Declared entries come
before every other item, oldest `created_at` first, then reference in byte
order. Its text line appends `: <blocked_work>` after ` — blocks <from>`.
When no item goes to the judge and an item declares blocked work, the judge is
`not_asked` with reason `only declared or waiting items`.

An item that declares no blocked work omits `blocked_work`, so a file without
the field gives the same report as before. An older `roundfix` that does not
know the field ignores it, as it ignores any key it does not read, and orders
that entry by the judge or by its type.

## Waiting Items

An open Backlog Entry or a pending or partial Finding may carry a Waiting
Declaration in its frontmatter:

```yaml
waiting_on: ERP vendor support answer
since: 2026-10-25
```

`waiting_on` names who or what the item waits on: one non-empty line, trimmed
when read. `since` is the calendar date the wait began, in `YYYY-MM-DD` form.
An absent key or YAML null is absent. The declaration applies only to Backlog
Entries and Findings, never to messages, notes or Specs.

A valid declaration makes the item a Waiting Item in the fifth class
`waiting`, outside the four ranked classes. The default list holds Waiting
Items back. `--all` lists them after every ranked item, oldest `since` first,
then reference in byte order, with priority `5`, `priority_name` and `basis`
`waiting`, and null `answers`, `blocks` and `reason`.
A Waiting Item is never sent to the judge, with or without --all.

Each source's `count` includes every item it read. The `backlog` and
`findings` sources add a non-zero `waiting` count; their text line appends
`: <n> waiting, held back without --all` to the count, or `: <n> waiting`
with `--all`. A listed Waiting Item appends
` — waiting on <waiting_on> since <since>` to its text line.

A malformed declaration keeps its item in its ordinary class, adds
`waiting_error` in JSON and ` — invalid waiting declaration: <waiting_error>`
in text, and never makes a source unavailable. Validation returns the first
applicable reason, in this order:

1. `waiting_on must be one line of non-empty text` — `waiting_on` is
   non-scalar, blank after trimming, or contains a line break.
2. `since is missing` — `waiting_on` is valid and `since` is absent.
3. `since must be a YYYY-MM-DD date` — `waiting_on` is valid and `since`
   is non-scalar or is not a calendar date in that form.
4. `since needs waiting_on` — `since` is present and `waiting_on` is absent.

The new JSON keys `waiting`, `waiting_on`, `since` and `waiting_error` are
omitted when zero or empty. With no ranked items the judge is `not_asked`,
with reason `no item` if the list is empty or `only waiting items` when every
listed item waits. Judge answer totals count only the items whose basis
is `judge` or `type`.

`--all` is refused with `--status`, with exit `2` and the diagnostic
`--all applies only to the follow-up list; drop it with --status`.
The status view reads no Backlog Entry or Finding and prints no waiting count.

The four ranked classes are `unblocking`, `health`, `feature` and `other`. Without
an answer, an item's declared type orders it, and only an entry that declares
blocked work is unblocking:
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
