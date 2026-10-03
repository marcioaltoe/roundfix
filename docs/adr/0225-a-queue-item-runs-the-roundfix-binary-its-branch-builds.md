---
status: accepted
created_at: 2026-10-03T00:00:00Z
updated_at: 2026-10-03T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A queue item runs the Roundfix binary its branch builds

On 2026-10-02 Spec 0215 added the Project Config key `verification.tools` and
wrote it into the repository's own `.roundfixrc.yml`. The queue owner had been
built from `main`, and the `roundfix implement` it started in the item
worktree read the item's Project Config and refused the key, so the item
parked `delivery-error`. The operator built `bin/roundfix` from the item branch
and retried with it, three times before the item merged. Every child the owner
starts in an item worktree, `implement`, `archive` and `review`, is the
owner's own executable reading the item's configuration.

Now a repository may declare, in Project Config, how to build the Roundfix
binary from its own tree: `delivery.item_binary` with a `build` command and the
repository-relative `path` it writes. The owner reads the declaration it loaded
at start, never the item's. Before each of those three children it runs the
build in the item worktree through the Verification executor and asks the
built binary `migrate --check` whether it would change the Run Database
schema. When it would not, the child runs with the built binary. When it
would, the child runs with the owner's own binary as before, and the console
log says why. A build that fails, or a binary that cannot start, parks the item
`delivery-error` with the build log path. A repository without the declaration
runs the owner's binary exactly as before; adopters install Roundfix from npm
and never declare it.

Two alternatives were rejected. Letting an older binary accept a configuration
key the item introduces would read the key and ignore its meaning, so the
item's own Runs would exercise the behavior it replaces; it would also weaken
the unknown-key refusal for every command. Re-executing the whole owner as the
item binary would let one item's unmerged queue code drive the merges of every
other item.

## Consequences

The item's code already runs as its Tasks' Verification, so running its binary
adds no new class of trust. It does mean the item's own Daemon settles its
Tasks and its own `review` assembles its review; the repository gate the owner
runs and the Pull Request's required checks still decide the merge. The schema
check exists because the Run Database is shared: an item binary that migrated
it would leave the installed binary unable to open it if the item never
merged. The Spec that introduces the key cannot use it: Roundfix's own Project
Config gains the declaration only after that Spec merges, in a separate change
made once the owner runs a binary that reads it, because an older owner would
refuse the key in the item worktree the same way. A build costs about a second on a
warm Go build cache, three times per item.
