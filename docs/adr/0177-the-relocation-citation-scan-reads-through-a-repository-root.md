---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The Relocation Citation scan reads through a repository root

The Relocation Citation scan (ADR-0173) checked each tracked path's parent
directories with `os.Lstat`, cached what it saw, and then opened the file by its
joined path with no-follow on the last component only. If a concurrent process
replaced an already-checked directory with a symbolic link, the open resolved
through that link and could read a file outside the repository. The same file
also satisfied the later `os.SameFile` check. This is the classic
time-of-check to time-of-use race (CWE-367): any check made on a path before it
is opened is stale by the time it is opened.

The scan now opens every tracked file through an `os.Root` opened once on the
repository root. That is the standard library's traversal-resistant API, added
in Go 1.24. On Unix it walks each component with `openat` and `O_NOFOLLOW`. On
Windows it opens components relative to held handles and never follows a
reparse point out of the root. A read can therefore never leave the
repository, whatever changes during the scan. The existing no-follow checks
stay: the per-component `Lstat`, the final `os.SameFile` comparison, and the
non-blocking open on Unix. They keep the promise that the scan never follows a
symbolic link inside the repository. No dependency is added.

## Consequences

`os.Root` follows a symbolic link that stays inside the root. A link planted
inside the repository mid-scan could still redirect a read to another
repository file, until the post-open checks see it. That read stays inside the
repository, and its content only produces citation warnings. The residual is
recorded as a limit, not hidden.
