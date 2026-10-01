---
status: accepted
created_at: 2026-10-01T00:00:00Z
updated_at: 2026-10-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The agent environment drops a Node preload whose file is missing

Roundfix passes its own environment to `acpx` and the ACP adapters, and every
one of them is a Node program. When `NODE_OPTIONS` asks Node to preload a file
with `--require` or `--import` and that file no longer exists, Node exits with
`Cannot find module` before it runs anything. On 2026-10-01 a preload placed in
the system's temporary directory by another tool was removed by temporary-file
cleanup, and the Implement preflight of a queued Spec refused with a
profile-proof failure that named the adapter, not the preload.

Roundfix now removes, from the environment it gives an agent process, each
`NODE_OPTIONS` preload whose value is a file path that does not exist, keeps
every other option in its original order, and prints one notice naming the
removed path. A preload given as a package name, or as a path that exists, is
kept. The user's shell environment and every other variable stay as they
were, and no other Node option is interpreted.

## Consequences

An agent process starts in an environment that differs from the user's only
by preloads that could not have loaded. A missing preload is no longer a
silent cause of a profile-proof or adapter failure; the notice names it. The
maintainer chose this over refusing with a reason that names `NODE_OPTIONS`,
because refusing would still stop an unattended queue on a file nothing in the
Run needs. A `NODE_OPTIONS` value Roundfix cannot split, such as one with an
unbalanced quote, is passed through unchanged.
