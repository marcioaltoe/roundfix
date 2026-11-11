### follow-up

```bash
roundfix follow-up [--status] [--all] [--messages <path>] [--format <text|json>]
```

Lists this repository's pending work in Follow-up Priority order. Use
`--status` to show work in progress, pending Specs, the next release and what
waits on the maintainer. The default output is text; `--format json` prints
one machine-readable report. `--messages <path>` adds messages supplied by the
caller, and is valid only for the list view.

It is read-only: it never writes a file, commits, pushes or prompts. It reads
the repository and its local Run Database, but never opens that database for
writing. It never reads or writes the Secondbrain, calls ai-memory, or reads
another checkout.

Exit codes:

- `0` — a report was produced, including when there is nothing to do or the
  judge did not answer.
- `1` — the Run Database exists but cannot be read.
- `2` — invalid usage or configuration, an unreadable or invalid messages
  file, or a working directory outside a Git repository.

Errors go to standard error as `roundfix: follow-up failed: <reason>` followed
by `Run 'roundfix follow-up --help' for usage.` Warnings also go to standard
error; the report goes to standard output.

## Sources

The list view reports these sources in this order. Each source is `checked`
with a count, `missing` when it was not configured or provided, or
`unavailable` with a detail when reading failed. An optional missing or
unavailable source does not fail the command.

- `ai-memory-messages` reads the optional `--messages` file. Without the flag
  it is `missing` with `no --messages file`.
- `backlog` reads open Backlog Entries.
- `findings` reads pending or partial Findings.
- `repository-inbox` reads every file in `docs/_inbox/` except a file or folder
  whose name starts with a dot, such as `.gitkeep`, and everything inside such
  a folder.
- `specs` reads active Specs. When the repository's Spec Root is absent,
  this source is `checked` with count `0` (ADR-0348).
- `runs` reads active Runs for this Git root; no Run Database is `checked`
  with count zero.
- `delivery-queue` reads unfinished Delivery Queue items.
- `pull-requests` reads open Pull Requests through `gh`; a failed `gh` call is
  `unavailable` with its first error line.
- `git` reads the branch, head, upstream, changed paths and unpushed commits;
  inspection failure is `unavailable`.

Each source carries these JSON fields: `name`, `state`, `count`, `detail`,
`waiting`.
A source-reading error can make the affected sources `unavailable`; the
report still lists the other sources.

## Messages file

The file passed to `--messages` is at most 1 MiB and must have this schema:
`roundfix/follow-up-messages/v1`. Each message needs a unique `id`, an RFC
3339 or `YYYY-MM-DD` `created_at`, and `text`; `from`, `subject`, `type` and
`blocked_work` are also supported.

A message type is empty or one of `feat`, `feature`, `finding`, `fix`,
`information`, `other`, `perf`, `refactor`. An empty or omitted type is stored
as `message`; a listed type keeps its given name. A refused type's error
lists the accepted types and says that `type` can be omitted.

```json
{"schema":"roundfix/follow-up-messages/v1","messages":[{"id":"m-17","from":"billing","created_at":"2026-10-21","subject":"Fix export","text":"The export blocks billing.","type":"fix"}]}
```

Messages become `message:<id>` items. Their title is `subject`, or the first
line of `text` cut to 80 characters, and `requested_by` is `from`.

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

## Follow-up Priority

The command asks four yes/no questions per item it sends to the judge and
computes the order in code: `is_fix`, `blocks_other_project`, `is_small` and
`affects_production`.
The classes are `unblocking`, `health`, `feature` and `other`. An item is
`unblocking` when it declares blocked work, or when it blocks another project
and is a fix or small; otherwise a fix is `health`, production work is
`feature`, and the remainder is `other`.
Small items come first inside a class, then the oldest dated item, then the
reference in byte order.

Without an answer, an item's declared type orders it, and only an entry that
declares blocked work is unblocking.
Types `finding`, `fix` and `perf`, and any item whose kind is `finding`, are
`health`; `feat` and `feature` are `feature`; `information`, `other` and
`refactor` are `other`. An omitted type (`message`) and any unlisted type are
also `other`. The report identifies whether its basis was `declared`,
`judge` or `type` and gives the reason when the fallback was used.

The judge sends nothing inside a Run, without a key, at the monthly ceiling,
when the Judge Log cannot be read, or for an item that matches a secret
pattern. A service failure stops the remaining requests. These cases still
produce a report and exit `0`.

## Text and JSON reports

For an empty list with no work in progress, including when all items are held back, the list view prints:

```text
Nothing to do: no pending item and nothing in progress.
```

Its JSON `state` is `nothing_to_do`; otherwise it is `work_found`. The JSON
report has schema `roundfix/follow-up/v1` and `view: "follow-up"`.

The list report carries these JSON fields: `schema`, `view`, `state`,
`repository`, `sources`, `judge`, `items`, `in_progress`.
The judge carries these JSON fields: `state`, `reason`, `model`, `calls`,
`cost_usd`.
Each item carries these JSON fields: `rank`, `priority`, `priority_name`,
`kind`, `ref`, `type`, `date`, `title`, `requested_by`, `blocks`, `blocked_work`,
`basis`, `answers`, `reason`, `waiting_on`, `since`, `waiting_error`.

For example:

```json
{"schema":"roundfix/follow-up/v1","view":"follow-up","state":"work_found","repository":{"root":"/work/app","namespace":"app"},"sources":[{"name":"backlog","state":"checked","count":1}],"judge":{"state":"unavailable","reason":"ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set","model":null,"calls":0,"cost_usd":0},"items":[{"rank":1,"priority":2,"priority_name":"health","kind":"backlog","ref":"docs/backlog/fix.md","type":"fix","date":"2026-10-21","title":"Fix export","requested_by":null,"blocks":null,"basis":"type","answers":null,"reason":"ROUNDFIX_OPENROUTER_JUDGE_API_KEY is not set"}],"in_progress":{"runs":[],"delivery_queue":[],"pull_requests":[],"git":{}}}
```

The status view reads only `runs`, `delivery-queue`, `pull-requests`, `git`,
`specs` and the Release Plan; it never calls the judge. Its text sections are
`In progress`, `Pending Specs`, `Next release`, `Blockers`, `Waiting on the
maintainer` and `Sources`. When there is no release change it includes the
line:

```text
Next release: no unreleased changes
```

When the Release Plan needs a manual classification, the status proposes no
version and says why instead of leaving the version blank. Its `Next release`
line reads
`Next release: not proposed (manual_classification_required): <n> commit(s) since <base> need a manual release classification`,
and `Waiting on the maintainer` lists
`release classification: classify them with roundfix release plan --from <base> --to HEAD --impact <none|patch|minor|major> --reason <text>`.
The JSON carries the first sentence in `release.reason` and the second in the
`detail` of the `release_classification` item of `maintainer`;
`release.proposed_version` stays empty, as in every state that proposes no
version.

Its JSON has `view: "status"`.
The status report carries these JSON fields: `schema`, `view`, `state`,
`repository`, `sources`, `in_progress`, `pending_specs`, `release`, `blockers`,
`maintainer`.

For example:

```json
{"schema":"roundfix/follow-up/v1","view":"status","state":"nothing_to_do","repository":{"root":"/work/app","namespace":"app"},"sources":[{"name":"runs","state":"checked","count":0}],"in_progress":{"runs":[],"delivery_queue":[],"pull_requests":[],"git":{}},"pending_specs":[],"release":{"state":"no_unreleased_changes","base":"","proposed_version":"","changes":[],"reason":null},"blockers":[],"maintainer":[]}
```
