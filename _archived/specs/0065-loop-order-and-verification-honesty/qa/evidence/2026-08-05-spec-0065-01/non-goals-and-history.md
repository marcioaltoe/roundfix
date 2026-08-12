# Non-Goals and history preservation

Build range:
`4d796ed22f3e2a81cc00a368b8d9af1747d5fbbe^..d603031e808e3c64a539c4875f00d62ed778f522`.

Fresh name-only range inspection shows product code changes only under
`internal/spec/` and `internal/speccheck/`, plus the declared documentation,
Skill, Baseline product, generated characterization, and active Spec paths.
No `internal/cli/`, Daemon lifecycle, QA-gate Skill, command flag, transport,
authentication, or status-ownership implementation path changed.

Fresh `git diff --quiet` exited 0 for ADR-0080 and ADR-0091. QA verdict typing,
environment-blocked equivalent-evidence semantics, the authored terminal QA
node, and Daemon-owned Task status therefore remain unchanged. The feature
does not require a bespoke harness for every Task; the public false-positive
canaries prove repository-wide gates remain accepted when an effect assertion
is also declared.

Fresh `git diff --name-only` for the complete range under
`docs/specs/_archived/` produced no output, and current archived-tree status is
clean. No archived Spec artifact was rewritten. The passing active/archived
corpus characterization independently confirms historical completed artifacts
retain their recorded finding counts; the new authoring refusals apply to
non-completed Tasks rather than retroactively invalidating settled evidence.
