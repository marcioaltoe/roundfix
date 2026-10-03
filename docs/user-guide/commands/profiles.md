### profiles

```bash
roundfix profiles check [--json]
roundfix profiles check --apply --scope user|project [--dry-run] [--yes] [--json]
roundfix profiles show [--category <category>] [--json]
roundfix profiles configure --scope user|project [--file <path>] [--remove <category>] [--dry-run] [--yes] [--json]
roundfix profiles validate [--category <category>] [--json]
```

`roundfix profiles check` compares every configured category with the shipped
Recommended Profile, including the full fallback order. It is read-only and
offline: it opens no Agent Session, reaches no network, and writes no file.
Undefined optional categories are omitted. Each configured category is:

- `current` when its profile equals the recommendation, even with a deviation.
- `differs` when it differs without a deviation for the shipped snapshot. An
  older deviation stays visible with the date it was declared against.
- `pinned` when it differs and its deviation names the shipped snapshot.

Text names each difference with the configured profile, its source, and the
recommended profile, then counts current, differing, and pinned categories.
When everything is current, it prints only the summary. `--json` uses schema
`roundfix/profiles-check/v1`, with `snapshot`, `current`, `differ`, `pinned`,
and `categories`. Each row has `category`, `status`, `source`, `configured`
and `recommended` (each with `preferred` and `fallbacks`), plus `deviation`
(`from`, `reason`) when declared. Rows follow Agent Work Category order.
The command exits `0` after a comparison, including one with differences,
and `2` for an unknown flag, an extra argument, or a configuration load error.

`roundfix profiles check --apply --scope user|project [--dry-run] [--yes]
[--json]` adopts only differing categories, replacing each complete profile
with the Recommended Profile and removing its old deviation. Current and
pinned categories stay unchanged. `--scope` is required; `--scope`, `--dry-run`
and `--yes` require `--apply`. Invalid flag combinations exit `2` without writing.
With `--scope user`, categories whose effective profile comes from Project
Config are skipped and named on standard error with advice to use
`--scope project`.

Adoption opens disposable Agent Sessions to prove every exact selection tuple
before confirmation or writing, using the same preview, confirmation, output
and exit codes as `profiles configure`. `--dry-run` proves and previews without
writing; `--yes` skips confirmation. Failed proof and declined confirmation
leave configuration bytes unchanged. `--json` uses
`roundfix/profiles-configure/v1`. With nothing to adopt, nothing is prepared,
proved or written: text prints `Profile configuration unchanged: nothing to
adopt`, JSON has `changed: false` with empty `profiles` and `changes`, and the
command exits `0`.

`profiles show` prints `Recommendation status: <status>` before each
Recommended profile block, with `inherited` for an undefined optional category.
A declared deviation follows that block as `Deviation: from <date> — <reason>`.
Its JSON adds `recommendation_status` and, when present, `deviation`; its schema
remains `roundfix/profiles/v2`, and its flags and exit codes are unchanged.

`profiles show` renders the effective Preferred Selection, Fallback Chain, and
dated advisory Recommended Profile: a Preferred Selection followed by its
Fallback Chain, with roles, source date, and rationales. Each of the ten
categories has its own profile. The Recommended Profile never selects, routes,
or writes configuration; interactive configure shows the same advisory rows.
Official model identifiers and advisory rank
do not prove that a tuple works in the current environment.

`profiles configure` writes a Profile Deviation its fragment carries, and
replacing a profile with a fragment without one removes the old deviation.
It merges a fragment by Agent Work Category. Every category
named in the fragment replaces that complete profile atomically; every other
configured category is preserved. Omission does not delete a category.
`--remove <category>` is the only way to remove a category and may be repeated.
Naming the same category in the fragment and with `--remove` is a validation
failure.

Before a write, the command shows one summary line for each affected category.
The three classifications are `added`, `replaced`, and `removed`; untouched
categories do not appear. The command runs Exact Agent Selection Proof for
every distinct preferred and fallback selection in the categories it adds or
replaces—the categories the operation writes—before confirmation. It does not
re-prove untouched categories. `--dry-run` performs the same proof and shows
the same summary without writing. Proof failure, refusal, or output failure
preserves the target bytes.

| `profiles configure` result | Exit code |
| --- | --- |
| Applied write | `0` |
| Already-satisfied no-op | `0` |
| Dry run, with no write | `0` |
| Refusal: declined confirmation or non-interactive use without `--yes` | `1` |
| Validation failure, including invalid flags, a fragment/removal conflict, or proof failure | `2` |

`profiles validate` is read-only, deduplicates exact tuples across category
references, proves them through disposable ACP Runtime Sessions, sends no
Agent prompt, and closes every Session on success or error. JSON schemas are
`roundfix/profiles/v2`, `roundfix/profiles-configure/v1`, and
`roundfix/profiles-validate/v1`.

`profiles validate` and `profiles configure` refuse a Cursor tuple without
the maintainer's login as `cursor_login_required`, before any session opens.
Run `cursor-agent login` in a terminal yourself, then retry the command.

A proof whose setup times out is retried once. A second timeout is classified
`temporary`; rerun the command when load drops because the configured profile
was not shown to be wrong.

