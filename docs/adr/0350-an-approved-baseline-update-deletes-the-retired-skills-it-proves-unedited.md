---
status: accepted
created_at: 2026-11-01T22:30:00Z
updated_at: 2026-11-01T22:30:00Z
deprecated_at: null
superseded_by: null
---

# An approved Baseline update deletes the Retired Skills it proves unedited

Supersedes ADR-0246 (2026-11-01); it is retired, and Git keeps its text.
Amends ADR-0191 (2026-11-01) for Retired Skills only.

ADR-0246 kept the Retired Skills out of the Baseline but left the installed
copies to the adopter: `roundfix baseline update` listed them under
`Skills retired` with the paths to delete and deleted nothing. The v0.127.0
release note said the update removes them, adopters trusted it, and Argus
removed seven copies by hand on 2026-11-01. The maintainer decided the same
day: "Apply apaga, valide se as skills context7 e exa estão fora da baseline
e sendo removidas, deveriam ficar."

The Retired Skills rules of ADR-0246 that stay:

- A Retired Skill is a skill the Baseline stopped requiring while the
  upstream catalog may still list it. Roundfix keeps the set in code.
- No module requires or dispatches a Retired Skill, no activation bundle
  names one, and the asset sync drops one from every Setup Snapshot it
  writes, even when the upstream list or the recorded snapshot names it.
- A Retired Skill that was Roundfix-owned leaves the binary's bundle, its
  canonical copy, its mirror and the owned-version record.
- `--no-skills` skips the Retired Skills with the rest of the skills stage,
  and Doctor ignores unrequired installed skills.

What changes:

- `roundfix baseline update` gives each Retired Skill the repository holds
  an action. It is `delete` when every copy it finds is provably the one the
  repository installed: each folder at `.agents/skills/<skill>` or
  `.claude/skills/<skill>` hashes to the `computedHash` of the skill's
  `skills-lock.json` entry, and a `.claude/skills/<skill>` link points at
  `../../.agents/skills/<skill>`. A lock entry or a link with no folder
  behind it is `delete` too, because no content can be lost.
- Otherwise the action is `keep`, with one reason: `edited` when a folder
  differs from the recorded hash, `unrecorded` when a folder exists and the
  lock records no hash, `unsafe` when an entry is neither a folder nor such a
  link. A kept copy is listed with the paths to delete by hand and never
  changes the update's state or exit code.
- The preview lists each deletion, and the Plan Digest binds the deletion
  set. When nothing else is pending, a pending deletion makes the preview
  `plan_ready` with exit 3 instead of `current`. A preview with no deletion
  keeps its digest bytes.
- An approved apply (`--confirm-plan` with that digest, or `--yes`) deletes,
  after the skills stage, exactly the planned `delete` entries: the
  `.claude/skills/<skill>` entry, the `.agents/skills/<skill>` folder and the
  lock entry. It re-inspects first and deletes nothing when any planned entry
  changed. It reports each deletion under `Skills deleted` (`skills.deleted`
  in JSON). The repository commits the deletions with the rest of the
  update.
- Nothing outside the Retired Skills set is ever deleted. ADR-0191's rule
  that Roundfix never deletes an installed skill tree holds for every other
  skill: a renamed skill's old folder, `context7-cli` and `exa-web-search`,
  which every Setup Snapshot requires, and any skill the repository installed
  on its own.

Two alternatives were rejected. Deleting every listed copy, edited or not,
would destroy a repository's own changes with no copy outside Git, the case
FreeBSD's `pkg delete` refuses for a file that "fails original SHA256
checksum, not removing". Refusing the whole update when one retired copy was
edited would block every Baseline change on a skill the repository chose to
keep; a kept copy is reported and the rest proceeds.

## Consequences

An adopter's first approved update after this change deletes the unedited
retired copies its plan lists and keeps the rest with a reason. A repository
that wants a Retired Skill must keep it under another name or edit it, and
the next retirement only adds a name to the set.
