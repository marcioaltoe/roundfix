---
status: accepted
created_at: 2026-11-11T00:00:00Z
updated_at: 2026-11-11T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An older Roundfix refuses a repository whose Baseline a newer Roundfix wrote

The Setup Manifest records, as `roundfixVersion`, the release of the Roundfix
that last applied a Baseline change to the repository. A Roundfix binary older
than that release refuses, before it reads Project Config, every command that
loads the repository's configuration and every Baseline command, naming both
versions and `roundfix upgrade`; the read-only commands that a user needs to
recover warn and continue without Project Config.

## Context

Until 10/11 a developer of the Conexus project ran Roundfix 0.96.0 against a
repository configured with 0.137.7, and two Specs ran outside Roundfix because
"the configuration was incompatible" (reported by conexus). Measured at
`013af3b4` (v0.139.0): the repository recorded no Roundfix version anywhere;
Roundfix 0.96.0 failed every configuration-loading command with
`parse config ".roundfixrc.yml": worktree.shared_caches is not a supported
config key`, which names neither version nor a fix; and Roundfix 0.137.7 read
the same tree without complaint and proposed, as an ordinary "plan ready"
`baseline update`, to rewrite `docs/agents/agent-instructions.md` and the Setup
Manifest back to its own older text, removing the message rule Spec 0373 had
added. The Update Check of that older binary would have shown the plan with
`Next: apply`. A Run Database written by a newer binary is already refused
with "a newer roundfix wrote it, so use that binary or run 'roundfix upgrade'";
the repository had no such guard. The maintainer chose "Backlog + Spec de
versão (Recommended)" — option text: "Spec F agora: binário mais velho que o
Baseline do repositório recusa com mensagem clara. Os três sintomas viram
Backlog Entry esperando passos de reprodução de você ou do Diego."

## Decision

**The record.** The Setup Manifest gains one optional field, `roundfixVersion`,
a plain `MAJOR.MINOR.PATCH` release number. A Baseline plan that the CLI builds
carries the running binary's version as its writer version. A plan carries the
recorded value unchanged; only a plan with at least one file change records the
writer's version, and only when it is a release number above the recorded one
or none is recorded. A plan without a file change proposes none, so the
Managed Refresh still converges (ADR-0103) and a release that changes no
Baseline content makes no adopter commit. A development build and a pre-release
build record nothing new; nothing ever lowers the value.

**The refusal.** Configuration loading reads `roundfixVersion` before Project
Config and refuses when it is newer than the running release, with: `Setup
Manifest "<path>" records Roundfix <recorded>, newer than this binary's
<running>; a newer roundfix wrote this repository's Baseline, so use that
binary or run 'roundfix upgrade'`. Every command that loads configuration
fails with it through its usual failure path, and so do `roundfix init --scope
project` and the Baseline commands (`roundfix baseline`, `baseline update`,
`baseline plan`, `baseline apply`, `baseline skills restore`, `baseline skills
reconcile`), which would otherwise rewrite the repository with older content.
A value that is present but not a release number is refused too. An absent or
unreadable Setup Manifest, or one without the field, records nothing and
changes nothing. A development build is never refused.

**What keeps working.** `roundfix --version`, `help` and `roundfix upgrade`,
which read no configuration before installing, are unchanged. `roundfix doctor`
reports the refusal as its first line, `baseline-version: failed`, and runs its
other checks without Project Config. The Update Check, the Run Database readers
(`runs`, `attach`, `events`) and `roundfix stop` print the refusal as a stderr
warning, read no Project Config, and continue; the Update Check's Baseline line
shows the refused preview and its `next` stays as the release comparison sets
it.

**Reach.** Only binaries from the release that carries this decision on can
read the record. An older binary cannot learn the refusal; it sees a Setup
Manifest with an unknown field and its Baseline commands report the manifest as
incompatible instead of planning a rewrite. The record protects a repository
from its first Baseline change applied by such a release.

## Considered Options

- Record the writer's version on every plan, so it always names the last
  binary that ran: rejected because every release would then propose a
  one-line Setup Manifest change in every adopted repository and in this one,
  and the Managed Refresh would stop converging after each release.
- A minimum version that a developer raises by hand when a change breaks older
  binaries: rejected because it rests on discipline, and the maintainer asked
  for a refusal of any binary older than the Baseline.
- A separate marker file: rejected because the Setup Manifest is already the
  record every Baseline change rewrites and every Spec that changes the
  Baseline already regenerates, so no new derived file enters the repository's
  process.
- A key in Project Config: rejected because the Baseline commands do not write
  Project Config, and a person edits it by hand.
- Refuse every command, read-only ones included: rejected because the user must
  be able to run `roundfix doctor`, the Update Check and `roundfix upgrade` to
  recover, and to stop or inspect a Run.
- Explain an unknown Project Config key as a possible newer-version key:
  rejected for this decision; the refusal comes before Project Config is read
  when the record is newer, and a hand-written key a newer release documents
  keeps today's "is not a supported config key" refusal.

## Consequences

- A repository's Setup Manifest gains `roundfixVersion` at its next applied
  Baseline change; this repository's own record gains it the same way.
- Two people on different Roundfix releases in one repository: the one on the
  older release is refused after the other applies a Baseline change, and runs
  `roundfix upgrade`.
- Builds from source between releases carry the last release's number, so a
  repository written by such a build accepts that release.
