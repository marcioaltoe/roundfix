## Release planning

When the user asks to cut, prepare, or validate a release, start with the
read-only Release Plan Command:

```bash
roundfix release plan
```

Run it before changelog edits, version-file edits, tags, pushes, package
publication, asset uploads, or GitHub Release creation. The command creates no
Run, reads no Roundfix configuration, contacts no external service, and
mutates no repository or release state.

A range plan also reports the skills and baseline checks read-only after its
`Next action:` line. The `skills:` line reports the comparison that
`roundfix doctor` makes, and the `baseline:` line reports what
`roundfix baseline update --no-skills` would report. The JSON output carries
both checks under `checks`. When either line is not `ok` or `current`, it ends
with the next action `complete the skills and guides check in the release
runbook before the release Pull Request`. The skills and baseline checks remain advisory and preserve the
decision state and proposed version.

The `skill-coverage` check reads the Skill Coverage Map and Behavior Surface
Record at the base and target commits, including a past `--to` revision. A
Lagging Surface is a changed, added or removed surface whose covering skills
are unchanged and which has no applicable changed Coverage Review. Update a
covering skill or record a Coverage Review in
`docs/references/skill-coverage.json`, then rerun the plan. A removed surface
requires a covering skill edit. Uncovered surfaces never lag. A `behind` or
`failed` check blocks any range except `no_release`: the plan exits 3, prints
`Release blocked: skill-coverage`, and names the blocking check in its next
action. Its decision state, proposed version and approval question remain
unchanged. No target map (`not_declared`) or no base map (`introduced`) never
blocks; an unreadable map or missing record reports `failed`. JSON carries
`checks.skillCoverage` with its status, detail, optional next action,
`blocking`, and optional `lagging` items naming the surface, change and skills.
Reset mode is unchanged.

Stable tags may be written as `MAJOR.MINOR.PATCH` or
`vMAJOR.MINOR.PATCH`; the planner accepts both spellings. If the highest
reachable version exists under both spellings, preflight refuses as ambiguous,
naming both refs and the `--from` selector that resolves the ambiguity. When a
version is proposed, it keeps the spelling of the tag it was selected from.

A generic release request authorizes only a conclusive patch plan: state
`ready` with a patch proposed version. State `approval_required` for a minor,
major, or version-zero breaking proposal requires explicit human approval of
the printed approval question before any release mutation. State
`manual_classification_required` requires a rerun with
`--impact <none|patch|minor|major> --reason <text>`; that classification
records the impact and reason, but it does not approve a resulting minor,
major, or version-zero breaking version. State `no_release` means no release
is required for the committed range.

For the exceptional `0.0.1` release-history reset, require a clean committed
target and run:

```bash
roundfix release plan --reset-to v0.0.1
```

Reset mode is mutually exclusive with `--from`, `--to`, `--impact`, and
`--reason`. It inventories every local and remote stable tag and every GitHub
Release through complete pagination, sorts the inventory deterministically,
and binds the reset target, target revision, tag identities, and Release
identities to `planDigest`. Use `--format json` for the
`roundfix.release-plan/0.0.1` result.

A complete reset plan always returns `approval_required` with exit `3`. This is
the read-only planning boundary: the command creates no Run, changes no files
or refs, and exposes no tag or GitHub Release deletion action. Missing or
partial inventory, a dirty tree, invalid flags, or an invalid target exits `2`
without a partial plan.

Review the complete inventory and stop. After implementation and QA pass,
rerun the same command for a fresh inventory and digest. Any tag or GitHub
Release deletion is a separate destructive release operation and requires
explicit human approval for that fresh plan. Approval of setup,
implementation, QA, an earlier plan, or the printed approval question does not
grant deletion authority.

For ordinary range planning, after the plan's required decision is satisfied,
follow the repository release runbook. Preserve the existing tag-triggered
workflow: validate the tag, keep artifact versions in agreement, publish npm
packages through the release workflow, upload GitHub Release assets, and leave
the Upgrade Command asset contract unchanged.
