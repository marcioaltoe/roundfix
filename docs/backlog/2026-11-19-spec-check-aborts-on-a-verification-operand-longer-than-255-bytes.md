---
type: fix
status: open
created: 2026-11-19
spec: null
reason: null
---

# `roundfix spec check` aborts when a Verification operand is longer than 255 bytes

## Symptom

A Task whose Verification holds one shell word longer than 255 bytes with no `/`, such as a long `grep` phrase, makes `roundfix spec check` stop with exit 2 instead of reporting findings. Reproduced on 2026-11-19 at `a6a4a15b` with `bin/roundfix` built from it, on the active Spec 0411, by adding this bullet to a Task's `## Verification`, where the phrase is 256 letters `a`:

```text
- `grep -q -- 'aaaa…a' README.md` — expected: exit 0.
```

`bin/roundfix spec check <slug> --stage tasks` then prints:

```text
roundfix: spec check failed: read Task "<slug>/task_01.md" Verification operands: inspect Verification path "aaaa…a": lstat <repo>/aaaa…a: file name too long
Run 'roundfix spec check --help' for usage.
```

and exits 2. With 255 letters it reports `No findings.` The planning of Spec 0411 hit it with an ordinary 261-byte English sentence matched by `grep -qE`, and had to split the sentence in two.

## Where

`TaskVerificationFiles` in `internal/spec/collision.go:168-198` treats every shell word of a Verification command as a candidate path and calls `repositoryFile` (`:219`), whose `lstat` fails with `ENAMETOOLONG` for a path component above the 255-byte `NAME_MAX`. Only a missing file is treated as "not a path"; any other error is returned as `inspect Verification path …` and aborts the check (`internal/speccheck/undeclared.go:50`). A word holding one of `$`*?[]{}` is skipped earlier, which is why phrases with an escaped bracket pass.

## Expected

An operand that cannot name a file, because a component is longer than the file system allows, is treated as not a path, like a missing one, and the check goes on to report its findings.

## Evidence

The reproduction above, run by the author of Spec 0411. No finding exists.

## Decisions

- 2026-11-19, about 00:05 BRT (AskUserQuestion during the planning of Spec 0411): the maintainer chose to record it — "(c) bug do spec check".
