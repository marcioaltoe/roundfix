---
status: accepted
created_at: 2026-10-18T18:00:00Z
updated_at: 2026-10-18T18:00:00Z
deprecated_at: null
superseded_by: null
---

# A Task commit is scanned for secrets, and a retry accepts a Secret Rewrite

On 2026-10-18 a Delivery Queue item in the Conexus repository, archived on its
branch, parked at the final push as `secret-found`: five
`credential-assignment` matches, all one test fixture that each Task had
committed without a refusal. The printed next action asked the operator to
remove the value from every commit the blocker names. Doing so rewrites
history, and the retry of an archived item accepts only a head that descends
from its candidate, so it refused the rewritten head with "archived item head
... differs from candidate head ...". The operator recovered only by resetting
the branch to the commit before the archive. GitHub's documentation describes
the same work for a blocked push: "To resolve a blocked push, you must remove
the secret from all of the commits it appears in", which for an earlier commit
is an interactive rebase from that commit. On 2026-10-18 the maintainer chose
"Retry aceita head reescrito", "Scan em cada Task" and "Documentar fixture".

## Decision

**The Secret Scan reads each Task commit.** After a Task's Verification passes
and the Daemon has prepared the paths the Task commit would stage, and before
the Task settles, the Daemon scans the lines that commit would add relative to
`HEAD`: every line of a path `HEAD` does not hold and the added lines of a path
it does. A deleted path and binary content contribute nothing. A Secret Match
settles the Task `failed` with the reason `Task commit carries a secret:
<path>:<line>: <kind>; ...`, creates no Task commit, and leaves the Task's
changes in the Run Worktree, as every failed Task does. The scan reads added
lines, not whole files, so it refuses exactly what the Final Push would refuse
in that commit.

**A Delivery Retry accepts a Secret Rewrite.** For an item parked
`secret-found` whose Spec is archived on the item branch and whose head moved
without descending from the parked candidate, the retry accepts the head when
all of these hold:

- the item's Run start is an ancestor of the head, so the rewrite begins at or
  after the item's base;
- every path that differs between the trees of the parked candidate and the
  head is a path the blocker names, changed in place with the same mode;
- in each such path both versions have the same number of lines, and every
  line that differs holds a Secret Match at the candidate and none at the
  head.

The retry then records the head as the newest candidate and returns the item
to `reviewing`. Any other moved head is refused with today's reason, followed
by the rule that failed in parentheses.

**The next action names the path that works.** `secret-found` prints: "in
<worktree>, rewrite the item branch from the earliest commit the blocker names
so that each named line holds a value no Secret Pattern matches and nothing
else in the final tree changes, then run roundfix deliver retry <slug>".

**The detector is not loosened.** The user guide and the Roundfix and
`implement-task` skill texts describe fixture values the
`credential-assignment` pattern does not match: letters only, or fewer than 20
characters.

## Considered Options

- Comparing whole files the blocker names was rejected: it would admit any
  change in those files, and the QA gate audited the candidate's tree.
- Comparing the reported line numbers was rejected: a Secret Match names the
  line in the commit that added it, which later commits can move, so the
  number cannot be compared at the head. The lines that match at the candidate
  are what the review and QA saw.
- Comparing only the matched value's span was rejected: a Secret Match never
  retains its value or its position in the line, by ADR-0273.
- Printing only the reset to the commit before the archive was rejected: it
  undoes the archive, and the operator then repeats the rewrite.
- Exempting test files, or values that contain `test` or `fixture`, was
  rejected: a real credential in a test file leaks exactly as anywhere else.
- Scanning whole files at a Task commit was rejected: a Task that edits a file
  whose old line the default branch already holds would be refused for a
  value the push accepts.
- Giving the Agent a repair turn after a Task Secret Match was rejected for
  this change: the retry policy belongs to Verification diagnostics, and a new
  repair class is a separate decision.

## Consequences

- A fixture value that matches a Secret Pattern fails the Task that adds it,
  before QA, review and archive.
- The rewritten head is reviewed again from round 1, because a Reviewer
  Lineage requires the head to descend from the reviewed head; the Delivery
  Queue then gates, pushes and scans the rewritten history as usual.
- An operator who deletes a secret-bearing line, or changes anything else, is
  refused by the retry and named the rule that failed.
- This amends ADR-0273, whose scan covered the QA Report commit and every
  push, and whose operator removes the value from every commit the blocker
  names.
