### follow-up

```bash
roundfix follow-up [--status] [--messages <path>] [--format <text|json>]
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
- `repository-inbox` reads every file in `docs/_inbox/`.
- `specs` reads active Specs. When the repository's Spec Root is absent,
  this source is `checked` with count `0` (ADR-0348).
- `runs` reads active Runs for this Git root; no Run Database is `checked`
  with count zero.
- `delivery-queue` reads unfinished Delivery Queue items.
- `pull-requests` reads open Pull Requests through `gh`; a failed `gh` call is
  `unavailable` with its first error line.
- `git` reads the branch, head, upstream, changed paths and unpushed commits;
  inspection failure is `unavailable`.

Each source carries these JSON fields: `name`, `state`, `count`, `detail`.
A source-reading error can make the affected sources `unavailable`; the
report still lists the other sources.

## Messages file

The file passed to `--messages` is at most 1 MiB and must have this schema:
`roundfix/follow-up-messages/v1`. Each message needs a unique `id`, an RFC
3339 or `YYYY-MM-DD` `created_at`, and `text`; `from`, `subject` and `type`
are also supported.

A message type is empty or one of `feat`, `feature`, `finding`, `fix`,
`information`, `other`, `perf`, `refactor`. An empty or omitted type is stored
as `message`; a listed type keeps its given name. A refused type's error
lists the accepted types and says that `type` can be omitted.

```json
{"schema":"roundfix/follow-up-messages/v1","messages":[{"id":"m-17","from":"billing","created_at":"2026-10-21","subject":"Fix export","text":"The export blocks billing.","type":"fix"}]}
```

Messages become `message:<id>` items. Their title is `subject`, or the first
line of `text` cut to 80 characters, and `requested_by` is `from`.

## Follow-up Priority

The command asks four yes/no questions per item and computes the order in
code: `is_fix`, `blocks_other_project`, `is_small` and `affects_production`.
The classes are `unblocking`, `health`, `feature` and `other`. An item is
`unblocking` when it blocks another project and is a fix or small; otherwise a
fix is `health`, production work is `feature`, and the remainder is `other`.
Small items come first inside a class, then the oldest dated item, then the
reference in byte order.

Without an answer, an item's declared type orders it, and nothing is unblocking.
Types `finding`, `fix` and `perf`, and any item whose kind is `finding`, are
`health`; `feat` and `feature` are `feature`; `information`, `other` and
`refactor` are `other`. An omitted type (`message`) and any unlisted type are
also `other`. The report identifies whether its basis was `judge` or
`type` and gives the reason when the fallback was used.

The judge sends nothing inside a Run, without a key, at the monthly ceiling,
when the Judge Log cannot be read, or for an item that matches a secret
pattern. A service failure stops the remaining requests. These cases still
produce a report and exit `0`.

## Text and JSON reports

For an empty inventory with no work in progress, the list view prints:

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
`kind`, `ref`, `type`, `date`, `title`, `requested_by`, `blocks`, `basis`,
`answers`, `reason`.

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
