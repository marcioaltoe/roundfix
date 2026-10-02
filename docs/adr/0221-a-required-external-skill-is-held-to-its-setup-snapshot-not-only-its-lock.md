---
status: accepted
created_at: 2026-10-02T00:00:00Z
updated_at: 2026-10-02T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A required external skill is held to its Setup Snapshot, not only its lock

A required upstream skill was ready when its installed bytes matched the
`computedHash` of its `skills-lock.json` entry. A lock written by another
skills tool records `ref: main`, so a skill installed months ago keeps
matching its own lock while the Setup Snapshot the binary carries pins a newer
tree. On 2026-10-02 eleven of the twenty-nine required external skills in this
repository matched their lock entries and differed from the `treeDigest` their
snapshot pins at `b3c45a4`, and Doctor reported `skills: ok`.

Doctor now compares each required external skill that is installed and
matches its lock with the `treeDigest` of the snapshot entry the restore
command would use for the repository's Baseline Profile, through the same
portable tree digest. A skill that differs is reported as trailing the
snapshot, with the code `DR-SKILL-TRAILS-SNAPSHOT`. Trailing alone makes the
`skills` line `warn`, not `failed`, and the managed refresh
(`roundfix baseline update`) restores a trailing skill to its snapshot commit
through the existing restore path.

## Consequences

The lock stays the authority for whether a skill is installed and intact, and
Repository Skill Set readiness keeps its meaning, so no repository fails
Doctor on upgrade for this alone. A repository that edits an upstream skill by
hand now sees it as trailing, which is true. The comparison reads only local
files and the snapshot embedded in the binary; restoring needs the source
commit, from the network or from `--source-dir`, as before.
