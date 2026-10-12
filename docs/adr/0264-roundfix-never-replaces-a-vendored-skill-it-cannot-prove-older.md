---
status: accepted
created_at: 2026-10-12T00:00:00Z
updated_at: 2026-10-12T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Roundfix never replaces a vendored skill it cannot prove older, and a divergence is corrected in the Baseline first

ADR-0221 made Doctor compare each required external skill that matches its
lock with the `treeDigest` of its Setup Snapshot and call a difference
"trailing", and made the managed refresh restore it. The comparison is an
equality, not an order. On 2026-10-09 an adopter updated its external skills
from the upstream one hour after a Roundfix release; Doctor then reported five
newer skills as `DR-SKILL-TRAILS-SNAPSHOT`, and `baseline update --yes` would
have restored the older snapshot over them. The binary carries only the
snapshot's tree digest and commit, the lock records a branch such as `main` or
some commit, and the restore fetches one commit at depth one, so no offline
evidence orders two upstream revisions.

**A difference is not trailing.** A required external skill that is installed,
matches its lock and differs from its Setup Snapshot is a Setup Snapshot
Difference: newer or older is unknown. Doctor names it with
`DR-SKILL-DIFFERS-FROM-SNAPSHOT` without raising the `skills` line to `warn`
and without a next action. The managed refresh lists it under
`skills.differing`, keeps it, and never restores it unasked; its reason names
`roundfix baseline skills restore --profile <id> --skill <name>`, which keeps
its own preview and confirmation. `DR-SKILL-TRAILS-SNAPSHOT` is retired,
because nothing the binary holds can prove a lock-matching copy behind. A
missing skill, or one that no longer matches its lock, is restored as before.
This repository still holds its own required skills to the snapshot
(ADR-0224).

**A divergence is corrected in this order.** A vendored skill is never edited
locally, because the next sync replaces its folder. When it diverges from what
a repository needs: first, a rule that governs over the skill's default, in
the Baseline or as a Repository-Specific Normative Rule, which answers a
convention difference such as the package manager, the test runner or a
library choice; second, a repository-authored skill that supplies or overrides
only what is missing beside the vendored one; third, a repository-authored
replacement only when the vendored skill is useless or almost useless. An
objective upstream defect goes to the upstream as a pull request and never
justifies the third choice. The maintainer of the skills repository sent this
order on 2026-10-09.

## Considered options

Recording the tree digests of earlier snapshots, so that a copy equal to one
of them is provably behind, was considered and left out: it changes the Setup
Snapshot schema and the asset sync, and today no earlier digest is recorded,
so it would prove nothing until the next sync. Keeping `warn` for a
difference was rejected because a warning whose only remedy may downgrade the
copy is the defect. Comparing upstream ancestry during the update was
rejected because the restore's depth-one fetch holds no history.

## Consequences

An adopter that updates its skills with its own tool sees an informational
difference instead of a warning and a downgrade. An adopter whose copies are
genuinely stale is no longer prompted to restore them; it restores a skill
deliberately or updates through its own tool. ADR-0221 keeps its lock rule and
its comparison; its trailing classification and automatic restore are
replaced by this record.
