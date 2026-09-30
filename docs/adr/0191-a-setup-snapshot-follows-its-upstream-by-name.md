---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A setup snapshot follows its upstream by name

The Baseline's setup snapshots pinned the upstream skills repository at
`14fdf46` (2026-07-25) while upstream moved to `a4e18e4`. Upstream renamed three
skills (`context7` to `context7-cli`, `feature-systems-pattern` to
`app-renderer-systems`, `rust` to `rust-expert`) and moved the terminal-interface
skills out of `go-cli` into a new `go-tui` list. The asset sync could not take
the refresh: it refused to start while the catalog named the new skills, and
the refreshed snapshots were refused while the catalog named the old ones. The
catalog also let a module dispatch a skill its profile's setup did not list,
because only required skills were checked.

The Baseline now follows the upstream list by name:

- A setup snapshot mirrors one upstream list, and the asset sync is its only
  writer. The sync validates the catalog it produces, not the one it replaces,
  so a module edit that names a renamed skill and the refresh that brings the
  skill land together.
- A built-in profile takes the upstream setup that lists every skill its
  selected modules name. The Go CLI/TUI profile therefore takes `go-tui`. The
  catalog refuses a profile whose modules dispatch a skill, or whose owned
  activation bundles name a skill, that its setup does not list, and a module
  requires every skill it dispatches.
- A renamed skill is followed by its new name everywhere the catalog names it.
  Roundfix never deletes an installed skill tree. The next managed refresh
  installs the new name; `baseline skills reconcile` removes the old lock
  entry; the old directory is the repository's to delete.
- A capability whose evidence is an installed skill names the current skill
  and may also accept the names earlier snapshots installed for it, so an
  adopter holding the old name is not refused before the refresh that installs
  the new one.

## Consequences

Upstream membership changes arrive whole: every skill an upstream list adds
joins the snapshot, and only the skills a module requires are installed. A
future rename is one Spec with one sync run. An adopter's first update after a
rename keeps the old directory until the repository removes it, and a second
update reports `current`. The Baseline's frozen parity evidence keeps the old
Context7 path, which the prior-name rule still satisfies.
