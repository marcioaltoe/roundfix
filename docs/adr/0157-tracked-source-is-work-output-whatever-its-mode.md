---
status: accepted
created_at: 2026-09-17T00:00:00Z
updated_at: 2026-09-17T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Tracked source is work output, whatever its mode

The Daemon decides what a Task's settlement commit carries by filtering the
changed paths. One of its filters refuses every executable regular file, on the
premise that an executable file is a build artifact. That premise is false for a
repository that tracks executable source: this one tracks a commit-message hook,
a pre-commit hook and a debugging script, and a Task that edits any of them has
its change dropped from the commit while the Task still settles completed.

A path Git already tracks is therefore work output, whatever its mode. An
executable file that is not tracked stays refused, which keeps the original
protection against a build artifact riding into the tree.

A drop also stops being an advisory. When the Daemon cannot stage a path the
work produced, the Task does not settle completed: the outcome names the paths
and the cause, and the Run does not reach Clean on a console warning nobody
reads afterwards.

Two alternatives were rejected. Reading the file's content to guess whether it
is a build artifact replaces one heuristic with a weaker one. Staging every
executable file removes the protection that stopped compiled output from being
committed.

The accepted cost is that a Task which deliberately produces an untracked
executable, such as a compiled helper, must add it to the repository before the
Daemon will carry it.
