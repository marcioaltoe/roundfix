---
status: accepted
created_at: 2026-11-06T15:00:00Z
updated_at: 2026-11-06T15:00:00Z
deprecated_at: null
superseded_by: null
---

# A session starts with an Update Check, and the Roundfix skill applies the update

Each Roundfix release reached the adopting repositories as a note, "Roundfix
vX: o que muda neste repositório", sent by the operator to every project's
ai-memory inbox. On 2026-11-06 thirty of them were waiting, unread, beside the
real requests between projects (Argus 11, Quaestor 8, Fiscus 6, Vortex 5). The
maintainer asked, through his assistant, to stop the notes and let each
repository see a new release at session start, apply it with `/roundfix
--update`, and choose per repository between applying automatically and only
being told. He chose "Spec agora, padrão avisar (Recommended)" (option text:
"Paro as notas já. A Baseline ganha uma checagem no início da sessão: compara
a versão do Roundfix e da Baseline com a última publicada e mostra o que muda.
/roundfix --update aplica. Uma chave por repositório liga o modo automático; o
padrão é só avisar. Quando sair, cancelo as notas pendentes nas caixas.").

## Decision

**`roundfix update check` is the Update Check.** It is read-only and always
exits `0` for a report. It compares the installed Roundfix with the latest
published release and lists, newest first, the `### ` headings of every
GitHub Release body between them: the `CHANGELOG.md` section titles the
release workflow publishes. It then runs, in process, the preview that
`roundfix baseline update --repo <root> --format json` prints, and reports the
repository's Baseline as current or as a plan with its file, history and skill
counts. The release listing is one `gh api` call, cached for 24 hours in
Roundfix Home (`update-check.json`) and shared by every repository on the
machine, so the network is read at most once a day; `--refresh` reads it now. A
lookup that fails or times out becomes one line, `latest release not reached
(<reason>); the session continues`, and the cache keeps that failure for the
same 24 hours. A development build compares nothing.

**What changes is the published release notes plus this repository's own
Baseline plan.** Before an upgrade only the release headings and the installed
binary's plan are known; after `roundfix upgrade` the new binary's
`roundfix baseline update` preview names every file the update changes in this
repository, and its Plan Digest binds what is applied. Nothing is written by
a model or by hand.

**`update.mode` in Project Config selects `notify` (default) or `auto`.** It is
read like every Project Config key, over User Config, and validated at load.
`notify` reports and applies nothing. `auto` asks the session to run the
Roundfix skill's `--update` mode before its first task.

**The Roundfix skill's `--update` Skill Mode verifies, shows and applies.** It
runs the Update Check with `--refresh`, runs `roundfix upgrade` when a release
is newer, previews `roundfix baseline update` with the binary it then has, and
applies that plan with `--confirm-plan <digest>`. It asks first when the plan
deletes files (a history unit or a Retired Skill) or when the worktree has
uncommitted changes; in `auto` at session start it applies no Baseline plan
over uncommitted changes. It never commits, pushes or opens a pull request.

**The Baseline runs the check through an instruction, not a harness hook.** A
mandatory clause of the core module's cross-project continuity rule tells
every session to run the Update Check once before its first task, to relay it,
and in a Daemon-assigned turn to skip it. The clause reads the same in Claude
Code, Codex and OpenCode.

## Considered Options

- A SessionStart hook in each harness's project configuration. Rejected as the
  Baseline's mechanism. The Baseline manages Markdown guides and root blocks
  only, and a hook means merging repository-owned JSON (`.claude/settings.json`,
  `.codex/hooks.json`) and an OpenCode plugin, which is tooling configuration the
  core clauses forbid changing without express authorization. Codex loads a
  project's hooks only when the project is trusted and runs each one only after
  its exact definition is reviewed: an untrusted hook is skipped without a
  word, and a SessionStart hook trusted later never runs for that session.
  The command guide shows the hook a repository may add itself, and the check
  then runs once because the clause skips a session whose context already
  holds the report.
- The key in the Setup Manifest as a Baseline decision. Rejected: a decision
  renders guide text and changes only through a Baseline plan and its Plan
  Digest, while `update.mode` changes no guide byte and is read by the
  command at run time, as Project Config is.
- Showing the full release notes. Rejected: a session start carries headings,
  and `CHANGELOG.md` keeps the detail.
- Applying the update from the command itself in `auto`. Rejected: the
  command stays read-only, and the skill mode is where the session asks
  before a step that deletes files.
- Checking on every session without a cache. Rejected: it would add a network
  call to every session start in every repository, and update notifiers check
  once a day by default.
- Keeping the inbox notes. Rejected by the maintainer.

## Consequences

- The operator stops sending release notes to the inboxes; after the release
  that carries this decision, the operator cancels the notes still pending.
- Every adopter receives the clause at its next managed refresh, and its next
  session runs the check.
- A repository whose `roundfix` is older than this decision shows nothing until
  it runs `roundfix upgrade` once; the operational commands' daily freshness
  line still names the newer release.
- `auto` runs `roundfix upgrade`, which replaces the machine's binary for every
  repository, not only the one whose session started.
