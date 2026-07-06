---
spec: 0011-storage-lifecycle
status: archived
created: 2026-07-06
surfaces: [cli, infra, docs]
archived: "2026-07-06"
source_slug: 0011-storage-lifecycle
---


# Review Artifacts, Run Logs, and Spec Archiving

Roundfix scatters durable output into the user-scoped Roundfix Home: Round and
Review Issue artifacts land in a loose `reviews/pr-<n>/` root, and every Batch
writes an agent log file that duplicates what the Run Event Journal already
holds — tens of megabytes a day during dogfooding. Neither survives review or
travels with the feature it belongs to, and a completed Spec has no way to
retire itself out of the active set. This Spec puts each kind of durable output
where it belongs: review artifacts in the repository's spec tree, per-Batch
agent logs behind an opt-in switch, and completed Specs into an archive.

## Goals

- Review artifacts default into the repository's spec tree, never loose in
  Roundfix Home, resolved by an explicit-config-first hierarchy. See ADR-0029.
- Per-Batch agent log files stop accumulating by default and turn back on
  through one config key for development. See ADR-0030.
- A completed Spec can be archived — its folder moves under an archive root
  with archive metadata — through a first-class Roundfix surface.
- Review Issue titles and status-poll output read cleanly: no raw markup in
  titles, no repeated identical poll lines.

## User Stories

1. As a developer resolving a PR tied to a Spec, I want its Round and Review
   Issue artifacts written under that Spec's folder, so that the review record
   lives and versions with the feature instead of loose in Roundfix Home.
2. As a developer resolving a PR with no Spec, I want its artifacts written to
   an in-repo review root, so that nothing durable lands outside the
   repository by default.
3. As a developer who set an Artifact Directory explicitly, I want that
   location to keep winning, so that upgrading never overrides my chosen
   layout. See ADR-0029.
4. As a developer running production Runs, I want no per-Batch agent log files
   written unless I opt in, so that my disk stops filling with a redundant
   copy of the Run Event Journal. See ADR-0030.
5. As a developer debugging an Agent, I want a config key that turns per-Batch
   agent logs back on, so that I can inspect raw payloads on disk when I need
   to.
6. As a developer who finished a Spec, I want to archive it once every Task is
   completed and QA passed, so that the active spec set shows only live work.
7. As a developer reading a Review Issue list, I want issue titles free of
   table markup and emoji, so that the Work Queue and issue files are legible.
8. As a developer watching a Run, I want the status-poll line to print only
   when it changes, so that the stderr stream is not flooded with identical
   lines every interval.

## Core Features

1. **Review artifact location hierarchy.** Round and Review Issue artifacts
   resolve their storage root by: an explicitly configured Artifact Directory
   wins; else a PR associated with a Spec (newest `Roundfix-Spec` commit
   trailer on the PR head, or an explicit spec selector) stores under
   `docs/specs/<slug>/reviews/`; else `docs/specs/_reviews/pr-<n>/`. Roundfix
   still never commits or ignores these artifacts — versioning stays the
   repository owner's choice. See ADR-0029. Supersedes ADR-0003's default.
2. **Opt-in agent logs.** Per-Batch agent log files are not written by
   default; a config key (User or Project) enables them. The Run Event Journal
   remains the durable record either way. A Detached Run's console log is
   unaffected — it stays unconditional. See ADR-0030.
3. **Spec archiving.** A Roundfix surface archives a completed Spec: it
   verifies every Task is completed and QA passed, stamps archive metadata,
   and moves `docs/specs/<slug>/` to `docs/specs/_archived/<slug>/`. It refuses
   to archive a Spec with incomplete Tasks or no passing QA verdict.
4. **Review Issue title hygiene.** The Review Issue title derivation strips
   Review Source markup and emoji from CodeRabbit table-fragment titles, so
   titles are plain text.
5. **Status-poll dedup.** The watch status-poll stderr line prints on change
   only, collapsing runs of identical poll lines.
6. **Merge-readiness docs note.** The merge-readiness `missing` path names the
   documentation expectation that explains the state, so the output points at
   the next useful action.

## User Experience

Review artifacts appear inside the repository tree by default instead of
Roundfix Home; a developer who set `artifact_dir` sees no change. Production
Runs write no `*.log` files under the Run's artifact area unless the agent-log
config key is set. Archiving is one command (or one loop step) that either
moves the folder and reports the new path or refuses with the unmet condition
named. Review Issue titles and the watch poll line read cleanly; every other
output stays byte-stable.

## Non-Goals / Out of Scope

- Auto-committing or git-ignoring review artifacts — placement only; the
  repository owner decides versioning. See ADR-0029.
- Rewriting or migrating artifacts already written to Roundfix Home by earlier
  Runs — the new default applies to new Runs.
- Changing the Run Event Journal format or storage (ADR-0008 stands).
- A general log-rotation or retention policy for the opt-in agent logs.
- Un-archiving a Spec or any archive-browsing surface.

## Success Metrics

- A spec-associated PR resolve writes its Round artifacts under
  `docs/specs/<slug>/reviews/` and a spec-less PR under
  `docs/specs/_reviews/pr-<n>/`; an explicit `artifact_dir` still wins — all
  asserted in tests.
- A production Run produces zero per-Batch agent log files; with the config key
  set, the same Run writes them.
- Archiving a fully-completed, QA-passed Spec moves its folder and stamps
  metadata; archiving an incomplete Spec is refused with the reason named.
- A CodeRabbit table-fragment issue yields a plain-text title, and a stable
  watch poll prints one line, not one per interval.

## Decisions

- Review artifact location resolves explicit-config → spec-associated →
  in-repo default; Roundfix never commits or ignores them. See ADR-0029.
- Per-Batch agent logs are opt-in; the Detached Run console log stays
  unconditional. See ADR-0030.
- Spec archiving is exposed as a Roundfix surface that enforces the
  all-completed-and-QA-passed precondition; the `_archived` spelling stands.
- Title hygiene and poll dedup are output-shaping fixes with no config surface.

## Open Questions

None.
