### profiles

```bash
roundfix profiles show [--category <category>] [--json]
roundfix profiles configure --scope user|project [--file <path>] [--remove <category>] [--dry-run] [--yes] [--json]
roundfix profiles validate [--category <category>] [--json]
```

`profiles show` renders the effective Preferred Selection, Fallback Chain, and
dated advisory recommendations. Official model identifiers and advisory rank
do not prove that a tuple works in the current environment.

`profiles configure` merges a fragment by Agent Work Category. Every category
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
`roundfix/profiles/v1`, `roundfix/profiles-configure/v1`, and
`roundfix/profiles-validate/v1`.
A proof whose setup times out is retried once. A second timeout is classified
`temporary`; rerun the command when load drops because the configured profile
was not shown to be wrong.

