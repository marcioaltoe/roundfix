---
status: accepted
knowledge: true
created_at: 2026-11-17T17:00:00Z
updated_at: 2026-11-17T17:00:00Z
deprecated_at: null
superseded_by: null
---

# The Update Check names each configured profile that differs from the Recommended Profile

Amends [ADR-0181](0181-a-configuration-is-compared-with-the-recommendation-only-where-a-person-asked.md).

## Context

ADR-0181 compares a configured Agent Selection Profile with the Recommended
Profile only where a person asked: `roundfix profiles check`, one Doctor line,
`roundfix profiles show`, and a notice after `roundfix upgrade`. None of them
runs in the session's daily path. The Update Check of ADR-0371 does: every
session runs it before its first task, and the `--update` Skill Mode applies
what it reports. On 2026-11-17 the maintainer found a repository running an
older configuration that nothing had reported ("Já peguei configuração antiga
sem atualizar") and chose "Spec 0403 separada (Recommended)": refresh the
Recommended Profile, and make the Update Check and the `--update` Skill Mode
show each category that differs, with the command that fixes it.

The snapshot shipped since 2026-09-30 had aged. Codex 0.169.0 lists
`gpt-6.1-terra` (OpenAI, 2026-10-26) and `gpt-6-luna` beside the `gpt-5.6`
models it marks older; the Claude ACP adapter 0.107.3 advertises the `fable`
alias for Fable 5.2 (Anthropic, 2026-11-09), and `haiku` resolves to Haiku 5.5
(2026-10-07). The maintainer approved the new table on 2026-11-17: `review` and
`chore` lead with `codex / gpt-6-luna / max`, `docs` with
`codex / gpt-6.1-terra / high`, and the implementation and `frontend`
categories keep their selections.

Adopters copied the 2026-09-30 snapshot into their Project Config, and two of
them pin a Profile Deviation on `review` declared against it. Under ADR-0181 a
new snapshot ends such a pin and reports the category as `differs`, with no
word about what to do with the pin.

## Decision

**The Update Check compares too.** `roundfix update check` runs the same
read-only, offline comparison as `roundfix profiles check`, with the snapshot
of the binary that runs it, whenever it loaded the configuration in full:
User Config outside a repository, User Config and Project Config inside one.
It reports only the categories to act on, each on one line and in a
`recommendations` JSON member:

- A category that differs with no Profile Deviation names its configured and
  recommended profiles and the adopt command,
  `roundfix profiles check --apply --scope <user|project>`, for the scope that
  defines it.
- A category whose Profile Deviation names an older snapshot is reported as
  pinned to an older snapshot, with the date and reason of the pin, the
  recommendation, the key that re-pins it
  (`profiles.<category>.deviation.from`) set to the shipped snapshot, and the
  same adopt command.
- A current category, and a category pinned to the shipped snapshot, add
  nothing; with none to act on, the member is omitted and no line is printed.

The comparison changes neither the report's `next` nor its exit code, and it
writes nothing.

**The `--update` Skill Mode offers the commands, never runs them unasked.**
After an upgrade it runs the Update Check again with the installed binary,
whose snapshot is the one that matters, shows each reported category, and asks
once which to run. A re-pin is a one-key edit of a configuration file a person
decides; the session makes it only on that answer.

## Considered Options

- **A new `stale_pin` status in `roundfix/profiles-check/v1`.** Rejected: the
  schema's consumers, `--apply` and the Doctor counts treat an older deviation
  as `differs`, which is still true. The Update Check derives "pinned to an
  older snapshot" from the deviation the status already carries.
- **A command that re-pins.** Rejected: a re-pin changes only the date of a
  decision a person already wrote, and every configuration write in Roundfix
  goes through the proved, confirmed write of `profiles configure`. Naming the
  key is enough, and the decision stays with the person.
- **Comparing when a Run starts or in Preflight.** Still rejected, as in
  ADR-0181: a configuration is a decision, and Agent work never stops on it.
- **Recommending the `ultra` effort Codex 0.169.0 accepts.** Rejected by the
  maintainer: no selection uses it, and the picker keeps its efforts.

## Consequences

- A session learns at its start which of its profiles an older snapshot left
  behind, and the `--update` Skill Mode turns that into one question.
- Argus and Indexus see their `review` pin reported as pinned to 2026-09-30
  once they upgrade, with the choice to re-pin or adopt.
- The Update Check's JSON gains a member, so the release that ships it raises
  the minor version.
- A repository whose Project Config the Update Check did not read, because a
  newer Roundfix wrote its Baseline, is not compared.
