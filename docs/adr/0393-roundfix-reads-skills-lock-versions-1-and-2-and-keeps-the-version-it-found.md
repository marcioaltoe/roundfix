---
status: accepted
created_at: 2026-11-12T00:00:00Z
updated_at: 2026-11-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Roundfix reads skills-lock.json versions 1 and 2 and keeps the version it found

Roundfix reads a `skills-lock.json` of version 1 or 2. Every `computedHash`
it writes into the lock or compares with an installed folder follows the hash
algorithm of the version the file declares, and Roundfix never changes that
version. Any other version is refused by name, and the next action Roundfix
names for it is `roundfix upgrade`, never a skills command.

## Context

The `skills` CLI 1.8.0 (2026-11-10) writes `skills-lock.json` with
`"version": 2` and a new folder hash. Roundfix 0.141.1 refused every lock but
version 1 in `roundfix doctor`, `roundfix baseline update`, the release plan's
`skills` check and `roundfix skills install --target project`, and Doctor's
next action, `bunx skills experimental_install && bunx skills update -p -y`,
rewrote version 2 again (reported by vortex and by pantheon, 2026-11-12).

The published CLI 1.8.1 reads versions 1 and 2. When it reads version 1 it
recomputes each entry's legacy hash over the installed folder and, when that
matches, replaces it with the version 2 hash; it writes version 2 only when it
writes the lock. Its README states the version 2 hash: leave out
`metadata.json` files and the `.git`, `__pycache__` and `__pypackages__`
folders; turn CRLF into LF in a file without a NUL byte; sort paths by their
UTF-8 bytes; and feed each path, a NUL byte, the decimal content length, a NUL
byte and the content into one SHA-256. The version 1 hash Roundfix computes
(`internal/skillhash`, paths in en-US collation order, path then content, no
separators) was measured equal to the CLI's own legacy hash on all 33 entries
of this repository's lock and on a folder holding CRLF text, `metadata.json`,
`__pycache__`, `node_modules` and mixed-case names, so a version 1 hash
Roundfix writes migrates cleanly.

## Decision

**Versions read.** The two lock readers, the Repository Skill Set readiness
and the Baseline lock loader, accept the numbers 1 and 2. A missing version or
one that is not a number is a malformed lock: `skills-lock.json must hold a
version number and a skills object`. Any other number is refused with
`skills-lock.json version <N> is not supported; this Roundfix reads versions
1 and 2`.

**Hashes written and compared.** The bundle install's lock update, the Stale
Owned Lock Entry scan, the Retired Skill proof and the external skill restore
compute a folder's `computedHash` in the algorithm of the lock's version:
version 1 keeps today's digest, version 2 uses the CLI 1.8 digest above. A
lock Roundfix creates where none exists stays version 1, as before.

**The version is kept.** Roundfix never migrates a lock between versions. A
version 1 lock, this repository's included, stays version 1 and keeps its
bytes until the skills CLI rewrites it.

**Next actions.** Doctor's `skills:` line names `roundfix upgrade` as its only
next action when the lock version is unsupported, and `roundfix baseline
update` names `run roundfix upgrade, then rerun roundfix baseline update`. A
next action never names a command that would write the refused version again.

This decision amends no earlier ADR: ADR-0359, ADR-0350 and ADR-0189 speak of
"the installed tree's `computedHash`", which is now computed in the lock
version's algorithm, and ADR-0364 still checks a required external skill only
for presence.

## Considered Options

- Accept any version whose `skills` shape matches (vortex's suggestion).
  Rejected: version 2 changed the meaning of `computedHash` while keeping the
  shape, so an unknown version can change it again and every hash Roundfix
  wrote would be wrong without a sign.
- Rewrite a version 2 lock as version 1, or migrate every lock to version 2.
  Rejected: Roundfix would fight the tool that owns the file, and every
  `skills update` or `baseline update` would rewrite all entries.
- Write a version 1 hash into a version 2 lock. Rejected: the CLI then reports
  each such skill as modified and `--frozen-lockfile` fails.
- Keep Doctor's skills command as the next action for every lock failure.
  Rejected: for a version Roundfix cannot read, that command recreates the
  failure.

## Consequences

Adopters whose `skills update` wrote version 2 pass `doctor`, `baseline
update` and the release plan again, and `roundfix skills install --target
project` repairs a Stale Owned Lock Entry in a version 2 lock. A future lock
version is refused until a Roundfix release reads it. A symbolic link inside
a skill folder stays refused, although the CLI 1.8 hash follows one that
stays inside the folder.
