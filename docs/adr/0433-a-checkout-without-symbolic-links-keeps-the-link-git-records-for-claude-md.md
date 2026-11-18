---
status: accepted
created_at: 2026-11-18T09:55:00Z
updated_at: 2026-11-18T09:55:00Z
deprecated_at: null
superseded_by: null
---

# A checkout without symbolic links keeps the link Git records for CLAUDE.md

When Git checks out the `CLAUDE.md` link on a host that does not create
symbolic links, the regular file it writes is the link, not a deviation to
convert. A root `CLAUDE.md` that is a regular file whose bytes are exactly
`AGENTS.md`, with no newline, and that `git ls-files -s` records with mode
`120000` is the link in an unlinked checkout.

## Context

ADR-0432 makes `CLAUDE.md` always a link to `AGENTS.md`. Roundfix ships for
Windows, where Git without Developer Mode runs with `core.symlinks=false` and
writes a link as a regular file holding its target. Converting that file would
write a real link Git does not expect, or fail on a host that cannot create
one, and the Misplaced Repository Rule detector would report the bytes
`AGENTS.md` as a passage to move. The Baseline already avoids writing links
where a repository keeps copies (ADR-0341), but the root link has no copy
alternative.

## Decision

A managed refresh and an adoption plan leave an unlinked checkout of
`CLAUDE.md` byte-identical and add the warning
`baseline.root.claude-link.unlinked-checkout`, whose message tells the person
to let Git create symbolic links (`core.symlinks`) and check `CLAUDE.md` out
again. `roundfix doctor` reports it as `warn` with the same next action. Git
is asked only when the file's bytes are exactly `AGENTS.md`, through
`git ls-files -s -- CLAUDE.md`. A file with those bytes that Git records as a
regular file, or does not track, holds no passage of its own and is converted.

## Considered Options

- **Convert it like any regular file.** Rejected. Git still records the link,
  so the checkout would show a change no one made, and the conversion fails
  where links cannot be created.
- **Keep a copy of `AGENTS.md` on such hosts.** Rejected. ADR-0432 rejects a
  copy because it drifts, and a copy would also change the file Git tracks.
- **Ignore the file silently.** Rejected. Claude Code on that checkout reads
  the word `AGENTS.md` instead of the instructions, so the person must be told.

## Consequences

- The Plan Digest of an unlinked checkout does not move because of
  `CLAUDE.md`, and no Windows test is needed: a test reproduces the checkout
  with `core.symlinks=false` on any host.
- Roundfix reads Git's index for this one file, read only, without a lock.
