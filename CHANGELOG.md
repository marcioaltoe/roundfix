# Changelog

All notable changes to Roundfix are documented here. The format follows
[Keep a Changelog](https://keepachangelog.com/), versions follow
[Semantic Versioning](https://semver.org/), and commit subjects use
[Conventional Commits](https://www.conventionalcommits.org/) (validated by
Cocogitto). The tag drives the release: pushing `vX.Y.Z` runs the release
workflow, which publishes the npm packages and creates the GitHub Release from
this file's matching section.

## [0.1.0] - 2026-07-06

First public release. Roundfix drives two kinds of work to a clean, verified,
integrated outcome with local coding agents.

### Added

- **Review resolution** — `fetch`, `resolve`, and `watch` pull unresolved
  CodeRabbit review threads into local Round and Review Issue artifacts, drive a
  selected ACP Runtime (`codex`, `claude`, `opencode`) over bounded Batches,
  verify, commit per Batch, resolve source threads, and push only when nothing
  unresolved remains (`watch --until-clean`).
- **Spec execution** — `implement` runs a Spec's Task Graph in isolated Run and
  Task Worktrees with a ready-set wave scheduler (`worktree.concurrency`), one
  commit per verified Task, an optional `--qa` gate, and porcelain-only
  integration onto the user's branch.
- **Detached Runs** — `--detach` on `resolve`, `watch`, and `implement`
  re-executes Roundfix as a session leader that survives its caller; followed
  with `attach`, controlled with `stop`.
- **`doctor`** — read-only machine-readiness diagnosis (Node.js, pinned acpx,
  the configured Agent probe, and codex runtime hygiene) that mutates nothing.
- **Support commands** — `setup`, `init`, `upgrade`, `settle`, `archive`, `gc`,
  `stop`, and `skills` (`check`/`install`).
- **Worktree bootstrap** — `worktree.bootstrap` prepares each new worktree
  (install dependencies, migrate and seed databases, warm caches) after
  `worktree.copy` and before Agent work, so stateful monorepos verify cleanly.
- **Run store retention** — `store.journal_retention` plus `roundfix gc` bound
  the Run Event Journal and reclaim run-artifact storage; Active Runs, `runs`
  rows, and locks are never pruned.
- **Review artifact placement** — resolved into the repository's spec tree by
  default: an explicit Artifact Directory wins, else a spec-associated PR writes
  under `docs/specs/<slug>/reviews/`, else `docs/specs/_reviews/pr-<n>/`.
- **npm distribution** — `npx roundfix`, `bunx roundfix`, and global installs
  from one channel via a launcher package plus per-platform binary packages,
  built and published by a tag-triggered release workflow that also uploads
  GitHub Release assets for `roundfix upgrade`.

### Notes

- Requires Node.js 22.13+, acpx 0.12.0, and the GitHub CLI authenticated for the
  target repository.
- codex runtime hygiene gates on the `com.apple.quarantine` attribute (the real
  XProtect trigger) and code-signature validity — not `spctl --assess`, which
  rejects any signed CLI that is not a notarized app.

[0.1.0]: https://github.com/marcioaltoe/roundfix/releases/tag/v0.1.0
