---
status: accepted
knowledge: false
created_at: 2026-11-19T18:30:00Z
updated_at: 2026-11-19T18:30:00Z
deprecated_at: null
superseded_by: null
---

# A write to a harness permission file takes a preview digest, and Roundfix knows its entries by their text

`roundfix permissions` writes or removes only when it is given the digest of
a preview, and the digest covers the harness, the action, the target path,
the bytes of the target file and the entries to add and to remove. In a file
that takes no comment, an entry is Roundfix's when its text equals an entry
the Flow Command List holds or once held. The list carries no version
number.

## Context

ADR-0444 decided that Roundfix writes a harness permission file only as a
list the user saw and confirmed. Three technical questions remained: what
makes a confirmation belong to one preview, how Roundfix finds its own
entries later, and how a project learns that its applied list is behind.

Claude Code's settings file is strict JSON, "a `//` comment or a trailing
comma is a syntax error", so an entry cannot carry a mark. A Codex rules file
is a small program, and Codex reads every rules file of the project's rules
folder. The binary has no library that edits a JSON file and keeps its
bytes; the one ordered editor it has rewrites the indentation of the file it
owns.

## Decision

**The write takes the digest.** A preview prints the first twelve hexadecimal
characters of a SHA-256 over the harness, the action, the repository-relative
target path, the SHA-256 of the target file's bytes or the word `absent`, and
the sorted entries to add, Guard Entries to add and entries to remove. The
write computes the same plan again and proceeds only when the digests are
equal. Any change of the binary's list or of the file between the two runs
changes the digest, and the command then writes nothing and exits `3`.

**An entry is known by its text.** In the Claude Code file, an entry of the
allow list or of the ask list is Roundfix's when it equals a current entry or
Guard Entry, or a retired one. The binary keeps every entry a released list
ever held as retired once the list drops it. A project applied the list when
its file holds at least one `roundfix` entry, current or retired; a Git entry
alone proves nothing, because people write those themselves.

**Behind is a set difference.** A project is behind when its file lacks a
current entry or holds a retired one. No list version is written anywhere.

**The Claude Code file is edited in place.** The command finds the two lists
by their byte positions and inserts or deletes only string elements and
their separators. It never decodes and encodes the whole file.

**Codex gets a file of Roundfix's own.** `.codex/rules/roundfix-flow.rules`
is written whole and deleted whole. Its entries are read back from the fixed
form Roundfix writes; a file in another form counts as behind and is
replaced on the next confirmed write.

## Considered Options

- **A flag that applies without a digest.** Rejected. A session could then
  write a list nobody was shown, and a list that changed after the preview.
- **A marker entry or a side file that lists what Roundfix wrote.** Rejected.
  A marker is a rule the harness would try to read, and a side file drifts
  from the file it describes when a person edits either one.
- **A list version in Project Config.** Rejected with ADR-0444: no key.
- **Decode, change and encode the settings file.** Rejected. It reorders
  nothing in Go, yet it rewrites indentation, escapes and the final newline,
  and the product promise is that no other byte changes.
- **Merge into a rules file the project already has in Codex.** Rejected.
  Reading a project's own rules means evaluating a program.

## Consequences

- A removal also removes an entry with the same text that the user wrote,
  and the removal preview says so and names each entry.
- An entry that a release drops must stay in the retired set for good, or a
  project that applied the older list would keep it unseen. A test pins that
  the retired set never shrinks once a release has shipped it.
- A hand-edited Codex rules file of Roundfix's own is replaced, never
  merged.
- The digest is short on purpose: it guards against a stale preview, not
  against an attacker, who could compute it.
