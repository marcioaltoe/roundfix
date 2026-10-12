---
status: approved
granted: 2026-10-11
action: let roundfix archive promote a cited file of the archived Spec to docs/references/ and rewrite the Live Document citations into the archived folder, let the Delivery Queue stage and prove those rewrites, let roundfix baseline update list the History Citations its Pending History sanitization breaks, bind them into its Plan Digest, promote only on --promote and report Stale History Citations, and describe all of it in the glossary, the command references, the Roundfix Skill and the authoring and archive skills
consuming: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
paths:
  - .agents/skills/archive-spec/SKILL.md
  - .agents/skills/roundfix/SKILL.md
  - .agents/skills/roundfix/references/archive.md
  - .agents/skills/roundfix/references/baseline.md
  - .agents/skills/roundfix/references/deliver.md
  - .agents/skills/write-idea/SKILL.md
  - .agents/skills/write-prd/SKILL.md
  - .agents/skills/write-tasks/SKILL.md
  - .agents/skills/write-techspec/SKILL.md
  - internal/spec/archive.go
  - skills/archive-spec/SKILL.md
  - skills/roundfix/SKILL.md
  - skills/write-idea/SKILL.md
  - skills/write-prd/SKILL.md
  - skills/write-tasks/SKILL.md
  - skills/write-techspec/SKILL.md
operations:
  - implement
  - commit
  - push
  - pull_request
  - merge
---

# Approved authority for Spec 0255

On 2026-09-30 the maintainer asked for unattended work through every release
of the program and said of the skills "considere autorizado a ajustar todas as
skills se necessário". Of Baseline sources, guides and `.roundfixrc.yml` the
maintainer said "Autorizar os dois". On 2026-10-11, answering "Grants dos
caminhos governados" for cycle AE, the maintainer said "Concedo" for the
Governed Paths each Spec declares.

The maintainer decided the scope of this Spec:

- 2026-10-11, for the adopted Backlog Entry
  `docs/backlog/2026-10-12-an-archive-breaks-links-other-specs-make-into-its-references.md`:
  "Archive promove e reescreve". When the owner Spec archives, a file a live
  document still cites moves to `docs/references/` and the inbound links are
  rewritten; the source keeps one owner.
- 2026-10-12, for the adopted Backlog Entry
  `docs/backlog/2026-10-12-baseline-update-sanitizes-history-without-reporting-the-citations-it-breaks.md`:
  `baseline update` lists every citation the Pending History sanitization
  breaks, in text and JSON and covered by the Plan Digest, and promotes the
  cited files only when the operator passes `--promote`, as
  `history sanitize` does. No automatic promotion in the update.

This Spec changes owned skills and `internal/spec/archive.go`, and no Baseline
module, Baseline guide, Makefile, `.roundfixrc.yml`, CI workflow, lint or
formatter configuration, or `go.mod`. It therefore bounds only the paths
above and declares no module regeneration.

No live provider call is authorized by this record. Authoring, tests,
Verification and QA use temporary repositories, temporary homes and synthetic
Specs and history. No adopter repository is read or written by any Task or QA
row.

The governed set was measured with `GovernedPath` on the authoring branch at
`7a8b7f70`. The probe was a `go test -overlay` test in `internal/speccheck`
that wrote nothing to the repository, run against every file the Tasks
declare. Sixteen of them are governed.

## Why each governed path is unavoidable

- `internal/spec/archive.go` holds the archive itself, which must find,
  promote and rewrite the citations before it removes the folder.
- `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md` and
  `.agents/skills/roundfix/references/deliver.md` describe the three changed
  command surfaces, and the skill-sync rule requires them to ship with the
  behavior.
- `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md` and
  `.agents/skills/write-tasks/SKILL.md` tell a secondary Spec to link the
  owner's copy; each must say that the archive relinks that copy.
- `.agents/skills/archive-spec/SKILL.md` describes what an archive promotes.
- `.agents/skills/roundfix/SKILL.md` and the `SKILL.md` mirrors under
  `skills/` are included because an owned skill's content changes only with
  its version, and the record command raises both front-matter fields.

## What is not governed

These paths are ordinary:

- The other Go sources and tests under `internal/spec` and `internal/cli`.
- `CONTEXT.md` and `docs/user-guide/commands/{archive,baseline,deliver}.md`.
- The skill reference mirrors under `skills/roundfix/references/`.
- `skills/testdata/owned-skill-versions.json`.
- `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
  and `docs/references/2026-10-12-live-citations-characterization.md`.

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
  pushes, or reads or writes an adopter repository.
- No change to archived Specs, existing Archive Records, ADR-0083, ADR-0247,
  ADR-0254, `CHANGELOG.md` or the `### QA settlement` section of any skill.
- No release, tag or deployment of this repository. Verification stays
  Daemon-owned, and Task status stays Daemon-written.
