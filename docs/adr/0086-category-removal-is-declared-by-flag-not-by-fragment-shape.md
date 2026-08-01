---
status: accepted
created_at: 2026-08-01T00:00:00Z
updated_at: 2026-08-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Category removal is declared by flag, not by fragment shape

`roundfix profiles configure` becomes merge-by-category, so a fragment can no
longer delete a profile by omitting it; removal needs its own way to be
expressed. Removal is therefore declared with a repeatable `--remove
<category>` flag, and the fragment stays purely additive data that can only add
or replace the categories it names. Encoding removal in the fragment instead —
a null value or an empty mapping for the category — was rejected because it
makes deletion a property of data shape, which is precisely how the original
defect destroyed four configured profiles: a fragment that was merely
incomplete was indistinguishable from one that intended removal. A flag cannot
be produced by an incomplete file, appears in shell history and CI logs, and
renders in the pre-write summary as its own line.
