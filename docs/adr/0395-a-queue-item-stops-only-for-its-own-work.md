---
status: accepted
created_at: 2026-11-12T21:00:00Z
updated_at: 2026-11-12T21:00:00Z
deprecated_at: null
superseded_by: null
---

# A queue item stops only for its own work

The Queue Owner writes every workflow message to its own log and never to a
nil writer, and an owner that stops while an item can still advance keeps its
owner record, so `roundfix deliver status` reports it gone. A review finding's
anchor may end a sentence. The review validator's sealed session works in a
private directory that holds nothing but a copy of the reviewed repository's
`.ai-memory.toml`, so tools that name a project from the working directory
record the session under the reviewed repository.

## Context

On 2026-11-12 a Delivery Queue item of this repository stood at `gating` for
about 8.5 hours, and quaestor reported the same stop the same morning. The
repository gate waited for the Machine Verification Capacity that another
Run held (ADR-0384), and its wait callback wrote the wait line to the
workflow's log, which the production constructor never set. The Go
specification says that calling a method of a nil interface value "causes a
run-time panic", so the owner panicked. The same specification says that a
panicking function's deferred functions "are then executed as usual"; the
owner's deferred release cleared its owner record on the way out. Status
therefore showed the item at `gating` and no `Owner:` line, because ADR-0279
reports only a recorded owner proven gone, and nothing told the operator that
the owner had died. An engine error that ends the owner leaves the same
silence today.

On 2026-11-11 the pre-PR review of Spec 0376 answered two real findings whose
anchors ended a sentence, `docs/user-guide/commands/skills.md:34.`. The anchor
grammar of ADR-0196 accepted a space, `:`, `,` or `)` after the line, so
neither anchor was read and the review blocked the item with `findings name no
file and line`.

On 2026-11-12 pantheon reported eight ai-memory projects named
`roundfix-review-validator-<n>`. The review validator opens its sealed session
in a new private directory under the system temporary directory, outside
every checkout, so the session has no repository file to read and loads no
repository instruction, the property ADR-0069 states for Baseline analysis.
ai-memory's published marker reference says that "a non-git directory outside
`$HOME` needs a marker in its exact cwd", and otherwise names the project
after the directory.

## Decision

**The owner's workflow writes to the owner's log.** The Delivery Engine and
every workflow step of one Queue Owner share one writer, the owner's standard
error, which is its console log. A workflow built without a writer writes to
standard error; no workflow message is ever written to a nil writer.

**An owner that stops with an item left keeps its record.** When the Queue
Owner process ends, by an engine error or a panic, it releases its owner
record only if no item can still advance, the condition its normal exit
already uses. Otherwise the record stays, and `deliver status` reports the
owner with the `Owner:` line of ADR-0279; `deliver resume`, `retry` and
`start` already reclaim an owner proven gone, and `deliver stop` releases the
record after the owner it stopped has exited.

**An anchor may end a sentence.** A `path:line` or `path:start-end` anchor,
backticked or not, followed by `.` or `;` reads as one followed by a space
when a space or the end of the finding comes next. `path:2.5` and `path:2.md`
remain unreadable, so punctuation never shortens a line number or a file name.

**A sealed session carries its repository's marker.** When the reviewed
checkout's root holds a regular `.ai-memory.toml` of at most 64 KiB, Roundfix
copies it into the validator's private working directory before the session
starts. A sealed working directory may hold that one regular file and nothing
else; any other entry, a directory or a symbolic link of that name, is still
refused. The session keeps its denied tools and reads no other repository
file.

## Considered Options

- **Report a queue with no recorded owner in status.** Rejected for the same
  reason ADR-0279 rejected it: it would also report the window inside `deliver
  start` and `deliver retry` before the new owner claims the queue, and the
  state after `deliver stop`. Keeping the record of an owner that died names
  exactly the defect.
- **Recover the panic in the owner and park the item.** Rejected: a panic is a
  defect of Roundfix, not of the item, and parking would ask the item's author
  to act; the record kept makes `deliver resume` the next step instead.
- **Run the validator in the reviewed checkout or a Run Worktree.** Rejected:
  a checkout gives the session the repository's files and loads its agent
  instructions, which the sealed session must not see.
- **Create the directory under the Roundfix artifact area.** Rejected: that
  area is under the home directory, where ai-memory walks to the user's own
  `~/.ai-memory.toml` and still names the project after the directory.
- **Give the private directory a Git identity copied from the repository's
  remotes.** Rejected: it fabricates a repository that agents treat as a
  project root, and a repository whose marker names its project differently
  from its remote would still be recorded elsewhere.
- **Write a marker with a workspace and project Roundfix derives.** Rejected:
  Roundfix would encode ai-memory's naming rules. Copying the repository's own
  marker keeps that choice with the repository. ai-memory is named here
  because the Baseline already names it as the cross-project continuity tool,
  and there is no marker convention shared across tools; a repository without
  the file gets an empty directory, as before.

## Consequences

- A gate that waits for the Machine Verification Capacity prints its wait line
  in the owner's console log and continues when the capacity frees.
- After an owner dies with an item left, status prints `Owner: Delivery Queue
  owner PID <pid> is not running; run roundfix deliver resume` until it is
  resumed.
- A finding anchored `path:line.` or `path:line;` is validated against the
  candidate diff like any other.
- A repository that commits `.ai-memory.toml` at its root has its review
  validator sessions recorded under its own project; Baseline analysis still
  works in an empty directory (ADR-0069).
- This decision refines ADR-0279, ADR-0384 and ADR-0196 and supersedes none.
