---
status: accepted
created_at: 2026-10-20T01:00:00Z
updated_at: 2026-10-20T01:00:00Z
deprecated_at: null
superseded_by: null
---

# The Linux Lane reads the repository history read-only

On 2026-10-19 the on-demand Linux Lane run on `main` failed two tests that
read historical commits by object name, because the lane records the
candidate in a fresh repository that holds no history, while CI checks out
with `fetch-depth: 0`. The lane now mounts the repository's object store
read-only and registers it as an alternate of the container repository, so a
historical object is reachable as in CI; a store it cannot mount whole skips
the lane with its reason.

## Decision

**The object store is mounted read-only.** The lane resolves, on the host,
the object store Git reads in the lane's working directory, `git rev-parse
--path-format=absolute --git-path objects`. For a Run Worktree or an item
worktree this is the common directory's `objects`. `docker run` gains `-v
<objects>:/repo-objects:ro` after the module cache mount. Only the object
store is mounted: never the working tree, never refs, and never writable.

**The container repository borrows from it.** The lane script checks that
`/repo-objects` is a readable directory, exiting `3` with `cannot read the
repository history at /repo-objects` otherwise, and writes `/repo-objects`
into `.git/objects/info/alternates` right after `git init`. The candidate is
still recorded as one fresh commit; objects that commit adds are written into
the container's own repository.

**A store the lane cannot mount whole is a skip.** Before Docker is called:

- a store path holding a colon, outside a Windows volume name, cannot be
  expressed in `-v`, and skips with `the object store <objects> cannot be
  mounted: its path holds a colon`;
- a store path that is not a directory skips with `the object store <objects>
  cannot be mounted: it is not a directory`, so Docker never creates a host
  directory for a missing source;
- a store whose `info/alternates` holds an entry (a line that is neither blank
  nor a `#` comment) skips with `the object store <objects> borrows objects
  through info/alternates, which the Linux Lane does not mount`.

Refs are not made available. The 2026-10-19 run failed only the two tests
that read objects by name, and every other package passed without history.

## Considered Options

- Mounting each alternate store as well was rejected. Git inside the container
  reads the mounted store's own `info/alternates`, whose host paths do not
  exist there: a probe on 2026-10-19 printed `error: unable to normalize
  alternate object path` from every Git call and still failed the test.
  Mounting each store at its own host path made the probe pass, but a host
  path then becomes a container path that can shadow `/tmp`, `/go` or the
  image's tools, a relative entry resolves against `/repo-objects` instead of
  its host store, and a Windows path is not a container path. Roundfix's own
  worktrees never create alternates; a clone made with `--shared` or
  `--reference` does, and `git repack -a` followed by removing
  `objects/info/alternates` makes it self-contained.
- Running without history when the store cannot be mounted was rejected: the
  history tests would fail as a Lane Change Failure, which is the defect this
  record fixes.
- Mounting the checkout or its `.git` was rejected for the reasons ADR-0288
  gives: a worktree's `.git` is a file pointing outside the mount, and refs,
  the index and hooks are not the candidate's to read or change.
- Fetching or bundling the history into the container was rejected: the lane
  keeps `--network none`, and copying the store costs time and disk on every
  run for what a read-only mount gives at once.
- Copying the default branch and tags as refs was rejected until a test needs
  one: no test in this repository failed for a missing ref.

## Consequences

- A test reads any object the repository holds, including other branches',
  inside a container without network that is removed when it ends.
- An owner whose checkout borrows objects through `info/alternates` sees the
  lane skip with that reason until the store is made self-contained.
- The lane log's command line shows the object store's host path.
