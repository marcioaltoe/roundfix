---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-15T08:10:30Z
updated_at: 2026-08-15T08:10:30Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# The audit judges governed paths, and history keeps the set honest

The changed-path audit reads every path a Task commit touched as governed by the
tooling rules, so an ordinary Go file alongside an authorized asset fails a Task
whose only fault was being one commit. The universal clause defines the governed
class by kind — linters, formatters, typecheckers, test runners, build tools,
package managers, code generators, their configuration and scripts, ignore files,
plugin declarations, version pins — and nothing enumerates it, which is why the
audit could not tell the two apart. A declared set fixes that and invites the
drift a declared set always invites, so it is held to the record: every path any
authorization has ever bounded must be matched by the declared set, checked
mechanically. The set can grow deliberately and cannot silently narrow, because
narrowing it below what maintainers have actually protected fails.
