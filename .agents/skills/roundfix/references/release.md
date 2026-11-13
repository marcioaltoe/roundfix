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

A Release Ceiling recorded in the approved authorization record of every
Spec the range delivers answers `approval_required` within every ceiling when
each record grants `release` and every change that sets the impact is such a
delivery. The decision becomes `ready` with the same proposed version.
A Release Ceiling never answers a breaking plan or a manual classification.
A missing grant, an impact above a ceiling, or an impact-setting change outside
those deliveries keeps the approval question.

Text prints `Release authorization: answered` or
`Release authorization: not answered — <reason>`, followed by each Spec's
ceiling, record and revision, and release paths. JSON adds
`releaseAuthorization` with `answered`, optional `reason`, and `grants`.
When no delivered record grants a ceiling, both outputs keep their bytes and
omit this authorization evidence. The `skill-coverage` check can still block
an answered plan.

Backlog Entries under `docs/backlog/` are planning evidence alongside
`docs/specs/`, `docs/adr/`, `docs/findings/` and `docs/handoffs/`; commits
changing only those paths are maintenance-only. A `refactor:` delivery that
changes released code is a commit to classify by hand.

The record carries `release_ceiling: patch | minor | major` and
`release_paths`, a list of exact repository-relative version-only file paths,
using the same path rules as `executables`. A ceiling requires `release` in
`operations` and non-empty release paths; paths require a ceiling. The plan
reads the record at the delivering commit's first parent and reports those
paths; the release runbook governs which files the later release changes.

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

## Completion notices

A Notice Request, a Cross-project Message that asks to be told when its work
is done, receives one completion notice under the cross-project continuity
rule of `docs/agents/agent-instructions.md`. Until then the request lives in
its record: the Backlog Entry or Finding that records the message, or the
`## Sources` entry of the PRD that adopts it. A PRD that adopts a Backlog
Entry or Finding holding a Notice Request carries the request into its
`## Sources` entry, so the Spec's own record keeps it after the Backlog Entry
or Finding is deleted.

Each notice falls due at one event, and the session that follows that event
finds the requests due with Git alone:

- **A release is published.** List the Archive Records the release adds with
  `git diff --name-only --diff-filter=A <previous tag> <new tag> -- <archive root>`,
  read each record's Spec PRD at the record's `source_revision` with
  `git show <source_revision>:<source>/_prd.md`, and send one notice to each
  sending project whose message its `## Sources` records as a Notice Request:
  what changed, the released version, and anything else the request asked the
  notice to carry. A record whose disposition is `cancelled`, written by
  `roundfix cancel`, delivered nothing, and its notice says that the work was
  closed without being delivered.
- **A Spec merges and releases nothing.** When the release plan for the range
  that holds the merge answers `no_release`, or the maintainer decides that no
  release carries it, read the Archive Records the merge adds with
  `git diff --name-only --diff-filter=A <merge commit>^ <merge commit> -- <archive root>`
  the same way and send each notice naming the merge commit instead of a
  version.
- **A record is closed without a Spec.** List the Backlog Entries and Findings
  the range deletes with
  `git log --diff-filter=D --name-only --format=%H <previous tag>..<new tag> -- docs/backlog docs/findings`,
  read each one before its deleting commit with
  `git show <deleting commit>^:<path>`, and send one notice for each Notice
  Request it holds: what was delivered and where when its work shipped without
  a Spec, or, when the work is closed without being delivered, such as a
  declined Backlog Entry, the notice says so instead. A record that a PRD's
  `## Sources` names was adopted by that Spec, and its notice follows the Spec
  instead.

When the work of one Notice Request spans several Specs, its notice follows
the last of them.

Send each notice once, from the session that follows the event and never
from a Daemon-assigned Run, and only to the project that asked: a completion
notice is not a release note. A tag, a merge commit and a deleting commit
each occur once, so each request falls due at exactly one event, and a
session sends notices only for the event it follows. Before sending, list the
outbox with `memory_message_list` and `box: "outbox"`: a pending notice to the
same project naming the same version, merge or closure was already sent, and
is not sent again. A session that cannot tell whether a notice for that event
already went, because an earlier session for the same event may have sent it
and the recipient may have consumed it, asks the user instead of sending a
second one. When the session exposes no ai-memory tools, put each notice's
text in the report for the user to deliver.
