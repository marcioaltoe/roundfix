---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-15T14:30:30Z
updated_at: 2026-08-15T14:30:30Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# A diagnostic names the literal it requires, and one cause reports once

A blocked row typed `blocked (finding: …)` is counted only when its text also
contains the exact string `" — waits on "`, and a row missing it is refused with a
message listing the three type prefixes the row already satisfies — never the
literal actually wanted. The uncounted row then makes the declared count disagree
with the table, so a second finding reports a counting defect that does not exist.
A refusal therefore names the form it requires, quoting it, and a count that
disagrees because rows failed to parse reports as that parse failure rather than
as a separate arithmetic complaint. One cause, one finding, and the fix in the
words the parser accepts.
