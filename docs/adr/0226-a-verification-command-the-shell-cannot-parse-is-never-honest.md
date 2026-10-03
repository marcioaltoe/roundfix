---
status: accepted
created_at: 2026-10-03T00:00:00Z
updated_at: 2026-10-03T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Verification command the shell cannot parse is never honest

On 2026-10-01 an amended Verification command held escaped backticks inside
its inline-code span. A Markdown code span ends at the next backtick and
ignores backslash escapes, so the Task Graph carried a truncated command that
the shell rejected as a syntax error. It failed before the work, so the
authoring probe called it honest; it failed after the work, so a complete Task
failed twice.

The one Verification prober now asks the shell to parse each command, without
executing it, before running it. A command the shell rejects is reported as
malformed: the Spec Consistency Check's probe prints the verdict `malformed`
and exits `1`, and the Daemon refuses the Task before the Agent starts, as it
refuses an unobserved command. The Spec Consistency Check also reports, without
running anything, a Verification bullet whose code span ends in a backslash or
leaves a backtick after it.

Authored Verification run outside a Run keeps no approval record and no
refusal. Authoring runs the probe on Specs that are not committed yet, by
design, so the probe names each Task file or Task Graph whose bytes are not
committed, and the person who runs it reads those commands first, as the
Baseline already says.

## Consequences

A truncated command is caught when it is written instead of after two failed
attempts. The provenance of an uncommitted Spec is visible in the probe's
output and in its JSON, and nothing that ran before is refused now.
