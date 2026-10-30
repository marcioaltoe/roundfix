---
spec: 0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty
status: active
created: 2026-10-30
surfaces: [backend]
---

# The follow-up command takes every documented message type, an absent Spec Root counts as empty, and the follow-up guides name the JSON fields

Quaestor's run of `roundfix follow-up` on v0.126.0 hit three frictions. The
`--messages` file was refused with exit 2 for the types `other` and `chore`,
although the reference says the rest of the types are `other` and lists no
accepted value. A repository whose last Spec was archived reported its `specs`
source `unavailable`, because the archive removes the Spec folder and Git
keeps no empty directory, so the Spec Root is gone from every fresh checkout.
And the reference named the item fields `class`, `reference` and `requester`,
while the JSON carries `priority_name`, `ref` and `requested_by`. After this
Spec the accepted message types are one list that the parser, its error, the
ordering fallback and both follow-up guides share; an absent Spec Root inside
the repository reads exactly as an empty one in every command that enumerates
Specs, as ADR-0348 decides; and a documentation contract test holds the guides
to the JSON field names the command emits.

## Sources

One Backlog Entry, minted from an ai-memory message at adoption and adopted
by this Spec's planning change (see `references/_index.md`):

1. [`references/2026-10-30-the-follow-up-command-rejects-documented-message-types-and-an-empty-spec-root.md`](references/2026-10-30-the-follow-up-command-rejects-documented-message-types-and-an-empty-spec-root.md),
   `type: fix`: Quaestor's report (ai-memory message
   `01a18c1a-90b6-7ca0-b720-b5c8672092b5` from
   `default/gesttione-solutions-quaestor`, 2026-10-30 06:18 BRT), whose body
   is kept verbatim. Only its three frictions belong to this Spec; the proof
   results in the same message are context.

At the authoring base `docs/backlog/` and `docs/findings/` hold no open entry,
so no other source shares this context.

## Prerequisites

None in `requires`. Authored against `35de7f74` (v0.126.0). Specs 0332, 0333
and 0334 are authored beside this one and deliver before it, serially. Spec
0333 also changes the Roundfix owned skill, so this Spec's Tasks re-record the
shared generated files (`.agents/skills/roundfix/SKILL.md` and its mirror,
`skills/testdata/owned-skill-versions.json`,
`docs/references/behavior-surfaces.json`, `CONTEXT.md` entries) with their
record commands after the rebase; they never hand-edit them. No Verification
of this Spec depends on another Spec's code.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme is minted or
  changed; message ids and item references keep their current form. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; the follow-up command's judge and `gh` calls are
  untouched, and no Verification or QA row reaches the network. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0348, new, records the Spec Root
  decision: ADR-0348: "An absent Spec Root holds no Spec." It keeps the
  archive's removal of exactly the Spec folder, so the exact-archive proof
  and every archive reader are untouched, and it keeps a missing external
  root refused. ADR-0319 owns the Follow-up Priority, whose fallback by
  declared type now reads one table: ADR-0319: "order is computed in code,
  never asked". Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AS decisions of 2026-10-30 quoted in the authorization record, and
  the Governed Paths this Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0335-follow-up-takes-every-documented-type-and-an-absent-spec-root-is-empty/_authorization.md`;
  bounded files:
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/follow-up.md`,
  `.agents/skills/roundfix/references/modes.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`.

## Goals

- A messages file whose types are the ones the guides document is accepted,
  and a refused type is answered with the accepted list.
- A healthy repository that archived every Spec is reported healthy by every
  command that enumerates Specs.
- An agent that reads the follow-up guides reads the JSON field names the
  command emits.

## Core Features

1. One message type table. The `--messages` parser accepts an empty `type`
   or one of `feat`, `feature`, `finding`, `fix`, `information`, `other`,
   `perf` and `refactor`: the Backlog types, the message type hints the
   Baseline names, and `other`. The same table gives an item's class when the
   judge did not answer, and a refused type's error lists the table.
2. An absent Spec Root is an empty one. A command that enumerates Specs
   reports an absent Spec Root inside the repository exactly as an empty
   Spec Root: the follow-up `specs` source is `checked` with count `0`, and
   `spec check`, `deliver plan` without a slug, `deliver status`, `stop`,
   `resume`, `history sanitize` and the history plan of `baseline update`
   give the output and exit code they give for an empty root.
3. Named Specs keep their refusal. `spec check <slug>`, `deliver plan
   <slug>`, `deliver start`, `deliver retry` and every other reader of a
   named Spec keep today's refusal text, and a broken link, a file or a
   missing external root stays refused.
4. Guides name the JSON. The follow-up user guide and the Roundfix skill
   reference name the JSON fields of the list report, its sources, its judge
   and its items, the status report's fields, and the accepted message types;
   a documentation contract test compares them with the command's types.
5. The `--follow-up` mode writes only documented types. The Roundfix skill's
   modes reference tells the agent to give a message entry a `type` only from
   the documented list, and to omit it otherwise.

## User Experience

An agent that writes `"type":"other"` gets the same report as for an omitted
type, with the item in the `other` class. An agent that writes
`"type":"chore"` sees `roundfix: follow-up failed: messages file "<path>":
invalid type "chore"; use one of feat, feature, finding, fix, information,
other, perf, refactor, or omit type` and the usage line, exit 2. In a
repository that archived everything, the follow-up `Sources` show
`specs: checked 0`, and `roundfix spec check` exits 0 as it does with an empty
`docs/specs/`.

Declared behavior changes: `feature`, `finding`, `information` and `other`
become accepted message types; a refused type's error text gains the accepted
list; `roundfix spec check`, `roundfix deliver plan` without a slug,
`roundfix deliver status`, `stop` and `resume`, `roundfix history sanitize`
and `roundfix baseline update` no longer refuse in a repository whose Spec
Root is absent; the follow-up guides' field prose changes.

## Non-Goals / Out of Scope

- Making `roundfix archive` or `roundfix cancel` write a placeholder file in
  the Spec Root (ADR-0348, Considered Options).
- Accepting every Conventional Commits type, such as `chore`, `docs` or `ci`,
  as a message type; `other` covers them.
- Changing the refusal of commands given a named Spec, the interactive
  `implement` picker (which refuses with or without a Spec Root), or the
  history notice of `roundfix upgrade`, which already skips silently.
- Changing the Follow-up Priority's questions, thresholds or classes, the
  judge, or the JSON schema version and field names.

## Success Metrics

1. Success Metric: a messages file with each accepted type is listed with
   exit 0, and `other`, `information` and `refactor` items fall in the
   `other` class, `feature` and `feat` in `feature`, `finding`, `fix` and
   `perf` in `health` when the judge does not answer.
2. Success Metric: a messages file with type `chore` exits 2 with the error
   that lists the accepted types from the same table the parser reads.
3. Success Metric: in one temporary repository, each enumerating command
   prints the same standard output and standard error and exits with the same
   code before and after an empty Spec Root is created, and no output names
   `which does not exist`.
4. Success Metric: a named Spec in an absent Spec Root, a dangling Spec Root
   link and a missing external Spec Root are still refused with today's text.
5. Success Metric: the documentation contract test fails when a follow-up
   guide omits or adds a JSON field or a message type, and passes on the
   delivered guides.
6. Success Metric: after release, Quaestor reruns `roundfix follow-up
   --messages` with type `other` in its repository that archived every Spec
   and reports both accepted and the `specs` source `checked`.

## Acceptance evidence

The outside-evidence row rests on text this Spec did not write. First,
Quaestor's report kept verbatim in the adopted Backlog Entry under
`references/` (ai-memory message `01a18c1a`, 2026-10-30). Second, the Git
SCM Wiki's GitFaq, section "Can I add empty directories?" (archived at
`https://web.archive.org/web/20120723224559/https:/git.wiki.kernel.org/index.php/GitFaq`,
retrieved 2026-10-30), whose passages the QA row quotes: "Directories are
added automatically when adding files inside them. That is, directories never
have to be added to the repository, and are not tracked on their own." and
"If you really need a directory to exist in checkouts you should create a
file in it." The QA compares, offline, Git's own behavior in a temporary
repository (a commit that removes the last file under `docs/specs/` leaves no
`docs/specs` in a fresh clone) with that passage, and then runs the delivered
binary there.

## Unreachable Acceptance

- criterion: Success Metric 6, after release, Quaestor reruns `roundfix
  follow-up --messages` with type `other` in its repository that archived
  every Spec and reports both accepted and the `specs` source `checked`.
  reason: it needs the released binary in another repository, which no
  Verification or QA row may read.
  satisfied-by: the release notice to `default/gesttione-solutions-quaestor`
  and Quaestor's reply through ai-memory.

## Validation

Characterized on 2026-10-30 at `35de7f74`, with `bin/roundfix` built by
`make build`, in a temporary repository holding only `README.md` and
`docs/history/specs/keep.md`, a temporary `HOME`, and no `docs/specs`.

- `roundfix follow-up --format json` exited 0 with the `specs` source
  `unavailable`: `specs.root resolved to "<repo>/docs/specs", which does not
  exist; create the directory or update specs.root`. `roundfix follow-up
  --status` reported the same `specs` line.
- `roundfix spec check`, `roundfix deliver plan`, `roundfix deliver status`
  and `roundfix history sanitize` each exited 2 with that refusal as the
  reason.
- After `mkdir docs/specs` in the same repository, `follow-up` reported
  `specs: checked 0`, `spec check` and `deliver plan` exited 0 with empty
  output, `history sanitize` printed `history sanitize plan: 0 unit(s)
  pending; nothing pending` and exited 0. (`deliver status` then failed for a
  separate reason: the temporary `HOME` had no Run Database.)
- A messages file with types `other` and `chore` exited 2 with `roundfix:
  follow-up failed: messages file "<path>": invalid type "other"` and the
  usage line.
- `config.ResolveSpecsRoot` (`internal/config/config.go`) refuses any
  resolved root for which `os.Stat` reports not-exist; `spec.ListActiveDetailed`
  (`internal/spec/spec.go`) already returns no Spec for a missing root.
- The archive removes only the Spec folder; this repository keeps
  `docs/specs/.gitkeep` since 2026-07, and two queue tests add one with the
  commit message "keep Specs Root after archive".
- Secondbrain, read only: `wiki/index.md`; `qmd query "roundfix follow-up
  messages type spec root empty directory archive gitkeep" --all --files
  --min-score 0.3` returned the mirrors of this repository's follow-up guide
  and archive records and the skills repository's vendored copy of the
  Roundfix follow-up reference; nothing that decides the friction.
- Exa: the Git SCM Wiki GitFaq passages quoted under Acceptance evidence, and
  Ubuntu's git-ubuntu documentation on empty directories ("the empty
  directories disappear as they pass through the index").

## Glossary

- changes: **Spec Root**

## Decisions

Maintainer answers to the cycle AS questions (AskUserQuestion, 2026-10-30),
verbatim:

- 2026-10-30: Plan — "Aprovar como está (Recommended)", option text: "Quatro
  Specs na ordem validada pelo Jev; a 0332 entra primeiro para que a fila
  resolva sozinha conflitos de registros gerados durante a 0333." This Spec
  is the fourth.
- 2026-10-30: Authorization — "Mesmas regras do AR (Recommended)", option
  text: "Push de branches, PR dos planos, entrega e merge pela Delivery
  Queue, eu escolho o ponto de release e aviso os adotantes por ai-memory."
- 2026-10-30: Release — "Uma no fim do ciclo (Recommended)", option text:
  "v0.127.0 depois das quatro Specs; os adotantes recebem a revisão de skills
  e as correções juntas."

Autonomous route: no person present for the technical choices; decisions taken from the operator's cycle AS brief for Spec 0335 of 2026-10-30 and the characterization above.

Taken from those sources, with no further interview:

- An absent Spec Root inside the repository is read as empty by the commands
  that enumerate Specs; the archive writes no placeholder (ADR-0348).
- The accepted message types are one table: the Backlog types, the Baseline's
  message type hints and `other`; `chore` stays refused with the list.
- A documentation contract test pins both follow-up guides to the command's
  JSON field names and the type table.

## Open Questions

None.
