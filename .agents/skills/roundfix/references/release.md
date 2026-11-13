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
or Finding moves into the Spec. When the work of one Notice Request spans
several Specs, each of their `## Sources` entries for it names the other
Specs.

Each notice falls due at one event, and the session that follows that event
finds the requests due with Git alone:

- **A release is published.** List the Archive Records the release adds with
  `git diff --name-only --diff-filter=A <previous tag> <new tag> -- <archive root>`,
  read each record's Spec PRD at the record's `source_revision` with
  `git show <source_revision>:<source>/_prd.md`, and send one notice for each
  Notice Request its `## Sources` records: what changed, the released
  version, and anything else the request asked the notice to carry.
- **A change merges and releases nothing.** When the release plan for the
  range that holds the merge answers `no_release`, or the maintainer decides
  that no release carries it, read the Archive Records the merge adds with
  `git diff --name-only --diff-filter=A <merge commit>^ <merge commit> -- <archive root>`
  the same way and send each notice naming the merge commit instead of a
  version.
- **A record is closed without a Spec.** Every terminal Backlog Entry or
  Finding is deleted in the change that records its disposition, so after
  each merge list the records it deletes with
  `git diff --name-only --diff-filter=D <merge commit>^ <merge commit> -- docs/backlog docs/findings`,
  read each one with `git show <merge commit>^:<path>`, set aside each record
  the merge adopted, as the next paragraph says, and send one notice for each
  Notice Request each other record holds, naming the merge commit: what was
  delivered when its work shipped without a Spec or, when the work is closed
  without being delivered, such as a declined Backlog Entry, the notice says
  so instead.

A record the merge adopted is not closed: its notice follows the Spec that
adopted it, and the closure event sends nothing for it. The adoption contract
of `docs/agents/docs-layout.md` moves an adopted Backlog Entry or Finding into
a Spec's `references/` folder, and a diff limited to `docs/backlog` and
`docs/findings` cannot pair that move with its destination under any rename
setting, so the listing of deleted records names every adopted record too.
The record read at `<merge commit>^` still carries the status it had before
the merge, so it cannot tell an adopted record from a closed one. List the
files the merge adds under the Spec root with
`git diff --name-only --diff-filter=A <merge commit>^ <merge commit> -- <spec root>`:
that listing, by file name and never by content, finds the copy of an adopted
record. To find a PRD that names a record, list the files there that mention
it with
`git grep -l -F -e <file name> <merge commit> -- <spec root>`. A listed record
was adopted when the merge adds a file of the same name under a Spec's
`references/` folder whose frontmatter, read with
`git show <merge commit>:<added path>`, records the adoption: `status: promoted`
with `spec` naming that Spec for a Backlog Entry, or `status: done` or
`status: partial` for a Finding. It was also adopted when a PRD that the
merge adds or changes names it in its `## Sources` entry at the merge commit;
list those PRDs with
`git diff --name-only --diff-filter=AM <merge commit>^ <merge commit> -- <spec root>`,
since every adoption writes the adopting Spec's `## Sources` in the same
change, and a PRD the merge leaves untouched never adopts a record. When the
merge adds a file of the same name whose frontmatter cannot be read or records
another status, and no such PRD's `## Sources` entry names the record, the
event sends no notice for it, and the session names the record and the added
file and asks the user.

A Notice Request's work is delivered at the event after which none of the
Specs its `## Sources` entry names is still under the Spec root: while one of
them remains there at the new tag or merge commit, the event sends nothing for
that request. Before sending, group the requests an event finds by the
ai-memory message each one records: entries in several Specs that the same
event archives and that record the same message are one request, which
receives one notice naming every Spec that delivered it. A
record whose disposition is `cancelled`, written by `roundfix cancel`, sends
no notice: the cancel returns each adopted Backlog Entry and Finding to its
folder, and the request stays there.

Send each notice once with `memory_message_send`, from the session that
follows the event and never from a Daemon-assigned Run, and only to the
project that asked: a completion notice is not a release note. Each notice
names the record that kept its request, so two requests from the same project
receive two notices. A tag and a merge commit each occur once, so each request
falls due at exactly one event, and a session sends notices only for the event
it follows. Before sending, list the outbox with `memory_message_list` and
`box: "outbox"`: a pending notice to the same project that names the same
record and the same event was already sent, and is not sent again. A session
that cannot tell whether a notice for that request already went, because an
earlier session for the same event may have sent it and the recipient may have
consumed it, asks the user instead of sending a second one. When the session
exposes no ai-memory tools, put each notice's text in the report for the user
to deliver.
