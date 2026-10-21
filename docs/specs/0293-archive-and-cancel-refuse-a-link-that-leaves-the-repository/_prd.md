---
spec: 0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository
status: active
created: 2026-10-20
surfaces: [backend, cli, docs]
---

# Archive and cancel refuse a link that leaves the repository

Minimal PRD for a bug fix, minted by `write-techspec`. `roundfix archive` and
`roundfix cancel` accept a path lexically inside the repository before
anything resolves it, so a symbolic link inside the repository whose target
lies outside it is followed: the Archive Record, a promoted copy or a
returned Backlog Entry is written outside the repository while the command
exits 0 and removes the Spec folder. After this Spec both commands refuse,
before writing anything, every path whose resolved target lies outside the
repository, with the message they already print for a path outside the
repository, naming the path. This is a declared break: a repository whose
`docs/history`, `docs/references`, `docs/backlog`, `docs/findings` or Spec
folder links outside can no longer archive or cancel until the link points
inside, while a link that stays inside, a checkout reached through a link and
an external Spec Root keep working with today's bytes. The TechSpec holds the
design.

## Sources

One Backlog Entry, written and adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-10-20-archive-and-cancel-accept-a-link-that-leaves-the-repository.md`](references/2026-10-20-archive-and-cancel-accept-a-link-that-leaves-the-repository.md),
   type `fix`. Spec 0288's pre-PR review finding F1 of 2026-10-20, dismissed
   then as the pre-existing lexical behavior of `RepositoryPath`, and the
   characterization at `70c78d15` that shows the record, a promoted copy and
   a returned Backlog Entry written outside the repository.

No open Backlog Entry or unresolved Finding shares this context. The only open
Backlog Entry,
`docs/backlog/2026-10-20-runs-causes-reads-no-diagnostic-of-a-run-that-recorded-no-artifact-root.md`,
is about where `roundfix runs causes` reads a Run's diagnostic; Spec 0292
owns it. `docs/findings/` holds no entry.

## Prerequisites

None. This Spec is authored against `70c78d15` (v0.100.0). Specs 0290, 0291
and 0292 are authored beside it and deliver before it; none of them declares
`internal/spec/repository_path.go`, `internal/spec/archive.go`,
`internal/spec/cancel.go`, the `archive` or `cancel` guide or the Roundfix
skill's `archive` reference, so the plan is rebased and checked again before
its Run.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the Archive Record keeps `roundfix/archive-record/v1` and its
  fields. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; both commands read and write only the
  checkout they run in, and every test uses temporary repositories and
  homes. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the fix records its decision and
  keeps every other archive and cancel decision. ADR-0306: "Archive and
  cancel refuse a path whose resolved target leaves the repository".
  ADR-0269: "Archive and cancel refuse ignored Spec files and finish a failed
  removal forward". ADR-0268: "A recorded worktree path is resolved within
  its Worktree Root". ADR-0272: "A review record survives a fix, and a
  correction inside the archived Spec's declared paths returns to review".
ADR-0247: "An archive leaves an Archive
  Record, and the Spec folder stays in Git". ADR-0260: "An archive relinks
  live citations, and a Baseline update lists the citations it breaks".
  ADR-0278: "A retirement deletes, and the History Root holds only Archive
  Records". ADR-0262: "A cancelled Spec leaves an Archive Record and returns
  its sources". ADR-0154: "A QA archive override records user authority, not
  a pass". ADR-0299: "The history conversion rewrites every citation it
  lists, and the check reports a removed history file cited without a Git
  Locator". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0271: "A test bounds by real time only what it asserts".
  ADR-0285: "Every test file compiles for Windows, and a repository contract
  proves it". ADR-0104: "A Spec accepts on evidence it did not author".
  ADR-0240: "One QA partial policy, and rows a Run sandbox cannot reach".
  ADR-0276: "A row the authorization keeps offline qualifies, and a byte
  invariant compares with the Delivery Base". ADR-0265: "A Spec proves what
  it claims, and every active Spec is held to it". ADR-0183: "An attributed
  claim carries a receipt the check proves". ADR-0184: "A TechSpec states a
  command surface as a transcript". ADR-0280: "A guide Task records its own
  surface fingerprint, a PRD records who decided, and the glossary is held to
  its own rule". ADR-0256: "A release waits for the skills that describe a
  changed surface". Set no obligation here, because this Spec changes none of
  their rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120
  (retired documentation lives under one history root); ADR-0167 (the pre-PR
  Pull Request row never decides a qualifying partial); ADR-0177 (the
  relocation citation scan reads through a repository root); ADR-0229 (an
  operator archive resumes any park from the Run start); ADR-0230 (an
  archived Spec keeps its relative links); ADR-0248 (existing history is
  sanitized in batches after a history-full tag); ADR-0250 (a module version
  is chosen when it is recorded, and the coverage record lists every
  platform); ADR-0251 (a Legacy Archive Folder is read leniently); ADR-0254
  (baseline update sanitizes pending history in the change it plans);
  ADR-0258 (a Spec is authored against the QA rerun classes); ADR-0261 (an
  Archive Record names its pull request); ADR-0277 (a Legacy Archive Folder
  conversion keeps the summary, links and dates its folder held); ADR-0279
  (status reports a Queue Owner proven gone); ADR-0281 (a test replaces an
  executable it rebuilds); ADR-0282 (the history cleanup names a superseded
  ADR in its successor); ADR-0286 (a fixed round-2 finding gets one
  Correction Review); ADR-0295 (reachability reads the fetched default
  branch); ADR-0298 (a timed-out profile proof gets a third attempt);
  ADR-0165 (a blocking review after archive parks publication for a
  corrective Spec); ADR-0197 (a pre-PR reviewer lineage spans at most two
  rounds); ADR-0237 (a Delivery Retry records a merge made outside the
  queue); ADR-0290 (merged items give back their Item Binaries); ADR-0169
  (the pre-PR review diffs the candidate from its merge base); ADR-0233 (a
  skill version raise is regenerated at merge); ADR-0196 (a pre-PR review
  finding parks only after validation); ADR-0301 (GC sanitation classifies
  the Artifact Root each Run wrote to). Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AM approval of 2026-10-20 through AskUserQuestion, authorization "the same
  rules as cycle AL", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).") and the 2026-09-30
  skills grant "considere autorizado a ajustar todas as skills se
  necessário". Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0293-archive-and-cancel-refuse-a-link-that-leaves-the-repository/_authorization.md`;
  bounded files: `internal/spec/archive.go`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`. No
  Makefile, workflow, Go module, lint, formatter, Baseline module or Project
  Config file change.

## Goals

- `roundfix archive`, in every form (normal, `--qa-override`, `--promote`,
  `--pull-request`, `--plan`), refuses before writing anything when a path it
  reads or a destination it would write resolves outside the repository, and
  names that path.
- `roundfix cancel` does the same, for its record, its returned sources, its
  promotions and its citation advice.
- A path whose target stays inside the repository keeps today's relative
  path, record bytes, relinked bytes and output, including a checkout
  reached through a symbolic link and an external Spec Root.
- The break is declared to adopters in the guides, the Roundfix skill and
  the decision record.

## Core Features

1. The resolver accepts a path only when its resolved target lies inside the
   resolved repository root. It still returns the lexical relative path when
   the path is lexically inside, and the resolved one otherwise.
2. The archive's promotion destinations and the cancel's promotion and
   returned-source destinations, which are joined to the repository root and
   never pass through the resolver today, pass through it before anything is
   written.
3. Every refusal uses the text the step already prints for a path outside
   the repository, naming the path, and exits `2` through Preflight
   Validation with no write.
4. Documentation. The `archive` and `cancel` user guides and the Roundfix
   skill's `archive` reference state the refusal and what keeps working;
   ADR-0269 records that ADR-0306 narrows its symbolic-link bullet.

## Declared break

Today, in the cases the Source lists, both commands exit 0 and write the
record, a promoted copy or a returned source outside the repository. After
this Spec they exit `2`, print the refusal naming the path, and write
nothing. The release notes carry it under the version's section as: "A path
whose resolved target lies outside the repository, such as a `docs/history`
symbolic link to another directory, is now refused by `roundfix archive` and
`roundfix cancel` with exit `2`, where it used to be written through; a link
that stays inside, a checkout reached through a link and an external Spec
Root are unchanged (ADR-0306)." No existing test, golden or fixture changes.

## Non-Goals / Out of Scope

- Any caller of the resolver outside the two commands: it has none.
- Refusing every symbolic link on a path, or changing how Roundfix resolves
  worktree paths (ADR-0268).
- The pre-existing refusal of a knowledge workspace whose `docs/specs` is a
  link into another repository: Git refuses the linked spelling in the Spec
  folder status read before the resolver runs; that is unchanged.
- Writing through `os.Root`, or any change to how the commands write once a
  path is accepted.
- Rewriting any Archive Record already written.

## Success Metrics

1. Success Metric: in a temporary repository whose `docs/history` is a
   committed symbolic link to a directory outside it, `roundfix archive
   <slug> --qa-override --approval <a> --reason <r>` exits 2 with the refusal
   naming `<repo>/docs/history/specs`, and the outside directory, the Spec
   folder and `git status` are unchanged.
2. Success Metric: `roundfix cancel <slug> --reason <r>` in a repository whose
   `docs/backlog` links outside exits 2 with `cancellation path
   "<repo>/docs/backlog/<file>" is outside the repository`, writing nothing;
   the same holds for the archive's promotion into a linked
   `docs/references`, a Spec folder that is itself such a link, and `archive
   --plan --offline` with a linked `docs/history`.
3. Success Metric: a `docs/history` link to a directory inside the
   repository, a checkout reached through a symbolic link, and an external
   Spec Root configured through a symbolic-link spelling archive with exit 0
   and the same output, record and relinked bytes as the control.
4. Success Metric: `RepositoryPath` given a root and a path both spelled
   through a link to the checkout returns the lexical relative path, and
   `spec.Archive` and `spec.Cancel` refuse every outside destination before
   writing.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-20.

- The Go standard library documentation of `os.Root`, readable offline with
  `go doc os.Root`: "Methods on Root will follow symbolic links, but symbolic
  links may not reference a location outside the root."
- The Go blog post "Traversal-resistant file APIs" by Damien Neil, 2025-03-12
  (https://go.dev/blog/osroot): "If the attacker controls part of the local
  filesystem, they may be able to use symbolic links to cause a program to
  access the wrong file".
- The Go standard library documentation of `filepath.Rel`, readable offline
  with `go doc path/filepath.Rel`: "Rel returns a relative path that is
  lexically equivalent to targPath when joined to basePath".

## Unreachable Acceptance

None. Every criterion runs in temporary repositories; the blog post is
recorded under its own network-prohibited row and never decides a qualifying
`partial`.

## Validation

The decisions were checked on 2026-10-20 against these sources.

- Characterization at `70c78d15` with `bin/roundfix` built from that commit,
  in temporary repositories under a temporary `HOME` and
  `GIT_CONFIG_GLOBAL`. With `docs/history` linked outside, archive and cancel
  each exited 0 and wrote the record outside; with `docs/references` linked
  outside, `archive --promote` wrote the promoted copy outside; with
  `docs/backlog` linked outside, cancel wrote the returned Backlog Entry
  outside; with the Spec folder itself linked outside, archive exited 0. A
  `docs/history` link inside the repository, a checkout reached through a
  link, the macOS `/var` spelling and an external `specs.root` spelled
  through `/var` all archived. A tracked Live Document whose directory was
  replaced by an outside link was already refused (`citing file … is not a
  regular repository file`), and cancel already refused a Spec folder link
  (`unknown active Spec slug`). The TechSpec records each.
- A prototype of Core Features 1-2 on a scratch copy of `70c78d15` refused
  every outside case with the texts the TechSpec's transcripts give, kept the
  inside cases at exit 0 with the same output, and left every existing test
  in `internal/spec` and the archive, cancel and folder-status tests in
  `internal/cli` passing; the one `internal/spec` failure,
  `TestAuthorizationReaderResolvesPreservedHistoricalRecords`, read a
  historical commit that the scratch copy did not hold. The prototype was
  discarded.
- Secondbrain, read only. Read `wiki/index.md` and
  `wiki/concepts/caminhos-de-worktree-e-links-simbolicos.md`, which records
  that a tool comparing a stored spelling with Git's must resolve both sides
  and that refusing links outright discarded whole Runs; ran `qmd query
  "archive cancel symbolic link outside repository containment canonical
  path" --all --files --min-score 0.3`, which returned the mirrors of
  ADR-0262, ADR-0177, the external Spec Root decision, ADR-0269, the changelog, the archive record
  of Spec 0018 and the cancel guide, and `qmd query "link simbólico que sai
  do repositório escrita fora da raiz" --all --files --min-score 0.3`, which
  returned ADR-0268, ADR-0177, a fluxus reference on broken links and the
  2026-10-12 worktree digest. They confirmed that both sides must be resolved
  and that the external Spec Root is a supported layout; no Inbox Entry was
  written, because the briefing forbids writing there.
- Exa: the Go blog post above and the `os.Root` source
  (https://go.dev/src/os/root.go), whose package comment says "Methods on
  Root can only access files and directories beneath a root directory."
  Challenging: `os.Root` would confine the writes themselves, but it refuses
  only at the first escaping write, after earlier writes, and never covers
  the read-only plan and citation advice, so the decision refuses up front
  and records `os.Root` as a rejected alternative in ADR-0306.

## Glossary

None.

## Decisions

- 2026-10-20 (about 18:15–18:35 BRT): the cycle AM decisions through
  AskUserQuestion: scope "R → S → T (Recomendado)"; minor items "Distribuir +
  U (Recomendado)" (option text: "'runs causes' vai em T (mesmo resolvedor e
  storage). 'Esperar vaga' vai em S (gate e preflight da fila). O
  endurecimento de symlink vira a Spec U, no fim (v0.104.0)."); for this
  Spec "Recusar no archive/cancel (Recomendado)" (option text: "`archive` e
  `cancel` recusam esse caminho com a mesma mensagem da 0288, nomeando o
  caminho. Hoje ele é aceito, então a Spec declara a quebra."); the
  authorization, the same rules as cycle AL.
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AM briefing for Spec 0293 of 2026-10-20 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- Put the containment test in the one resolver the two commands share, and
  route the destinations it never saw through it, rather than a check per
  call site.
- Resolve both the path and the repository root, so a checkout or external
  Spec Root reached through a link is not mistaken for an escape.
- Keep the lexical relative path whenever it is inside, so nothing that
  works today changes its bytes.
- Record the decision as ADR-0306, because it is a behavior break that
  narrows ADR-0269.

## Open Questions

None.
