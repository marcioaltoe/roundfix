---
status: approved
granted: 2026-10-08
action: let roundfix baseline update plan and apply the pending history sanitize in the change it plans, creating the annotated history-full tag when absent, let roundfix upgrade name the pending history, accept the legacy list-of-maps unproven under Lenient Legacy Reading, print each Refused Unit reason on one line, and describe all of it in the glossary, the command references and the Roundfix Skill
consuming: 0249-a-baseline-update-that-sanitizes-pending-history
paths:
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - skills/roundfix/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0249

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". Of Baseline sources, guides and `.roundfixrc.yml` the
maintainer said "Autorizar os dois". The standing answer for the Governed Paths
each Spec declares is "Concedo".

On 2026-10-08 the maintainer decided the scope of this Spec:

- Where the sanitize runs: "No baseline update (Recommended)". After
  `roundfix upgrade`, a repository's `roundfix baseline update` detects the
  pending legacy history and includes the sanitize in its plan. It creates the
  annotated `history-full` tag when absent and applies the conversion together
  with the rest of the Baseline plan, in the same change the operator reviews.
  `roundfix upgrade` only notifies.
- The adopted Backlog Entry
  `docs/backlog/2026-10-08-legacy-unproven-maps-refuse-a-folder.md`: accept the
  legacy list-of-maps `unproven` under Lenient Legacy Reading, keep active
  Specs string-only, and print each Refused Unit reason on one line.
- Grants: the Governed Paths this Spec declares (skills, Baseline modules and
  guides under rule 18, the Makefile, `.roundfixrc.yml` and CI), and a
  `qa_override` for an environment-only `partial`. Adopter repositories are
  never touched; every fixture is synthetic.

This Spec changes the Roundfix Skill and no Baseline module, Baseline guide,
Makefile, `.roundfixrc.yml` or CI workflow. It therefore bounds only the skill
paths below and declares no module regeneration.

No live provider call is authorized by this record. Authoring, tests,
Verification and QA use temporary repositories, temporary homes and synthetic
legacy folders.

The governed set was measured with `GovernedPath` on the authoring branch at
`c73a92e0`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository, run against every file the Tasks
declare. Four of them are governed.

## Why each governed path is unavoidable

- `.agents/skills/roundfix/references/baseline.md` describes
  `roundfix baseline update`. The skill-sync rule requires the skill to ship
  with the history section, `--no-history` and the tag the update creates.
- `.agents/skills/roundfix/references/archive.md` describes the History
  Sanitize Command, which gains the list-of-maps `unproven`, one-line
  refusal reasons and its relation to the update.
- `.agents/skills/roundfix/SKILL.md` and its mirror `skills/roundfix/SKILL.md`
  are included because an owned skill's content changes only with its version,
  and the record command raises both front-matter fields.

## What is not governed

These paths are ordinary:

- The Go sources and tests under `internal/spec` and `internal/cli`.
- `CONTEXT.md` and `docs/user-guide/commands/{baseline,history,upgrade}.md`.
- The skill reference mirrors `skills/roundfix/references/baseline.md` and
  `skills/roundfix/references/archive.md`.
- `skills/testdata/owned-skill-versions.json`.
- `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`.

## Sanctioned regeneration

The repository-owned command resolves the skill mirrors. This declaration
records the regeneration that follows the approved edit and adds no source
path.

```yaml
command: make skills-sync
```

## Limits

- No new dependency in `go.mod`, and no change to the Makefile, the lint,
  formatter or test-runner configuration, `.roundfixrc.yml`, the CI workflows,
  any Baseline module or any Baseline guide under `docs/agents/`.
- No test, Verification command or QA row reaches a provider, starts a real
  Agent Session, reads a credential, reads or writes the real `~/.roundfix`,
  pushes a tag, or reads or writes an adopter repository.
- No change to archived Specs, existing Archive Records, `CHANGELOG.md` or the
  `### QA settlement` section of any skill.
- No release, tag or deployment of this repository. Verification stays
  Daemon-owned, and Task status stays Daemon-written.
