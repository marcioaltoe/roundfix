---
spec: 0099-a-tool-that-meets-the-repositorys-conventions
status: active
created: 2026-08-12
surfaces: [cli, backend]
---

# A tool that meets the repository's conventions

Twice in one session Roundfix assumed a convention the target repository does not
follow, failed closed, and offered no way to accommodate the real one. Archiving
refuses a Spec whose graph declined the QA gate — a shape the authoring contract
documents as valid — and refuses a declared-only partial verdict its own help text
promises to accept. Release planning demands a `v`-prefixed tag in a repository
whose release workflow uses the absence of that prefix to decide a release is
stable, so following the tool's instruction would publish the release as a
prerelease. The cost is not the command that fails. It is the operator concluding
they must change their repository to please the tool.

## Project Constraints

- Identifier strategy: applicable — Archive Command, QA Report, and the verdict
  vocabulary are glossary terms whose meaning this Spec widens, and a release tag
  form is an identifier the repository's own workflow reads. The closing node
  checks whether the accepted shapes changed a term the glossary should carry.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. Release planning is read-only over local Git refs,
  and archiving is a filesystem and frontmatter operation. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0080 makes QA verdicts distinguish
  environment-blocked rows, which is the vocabulary the declared-only partial path
  is written in; ADR-0088 authors the QA gate into the graph rather than
  requesting it per run, which is what makes a declined gate a graph-level fact
  the Archive Command must read. The decisions built on ADR-0080 are accounted:
  ADR-0091 makes the gate a Task node of its own type, which is what carries the
  declined shape this Spec teaches the Archive Command; ADR-0097 lets a QA row
  carry forward only on declared, unmoved evidence, and applies directly, since the
  declared-only partial this Spec accepts is exactly evidence of that kind;
  ADR-0093 checks Spec consistency by citation rather than inference, and ADR-0096
  with ADR-0117 place the gate's mechanical stage and its checks — this Spec
  changes none of those three. ADR-0104 makes a Spec accept on evidence it did not
  author, and applies: this Spec's acceptance rests on a repository it did not
  build, whose Spec declined the gate and could not be archived. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work is production Go in the archive, spec, and release
  planning packages plus their tests, and the CLI help text they own. Source:
  `docs/agents/agent-instructions.md`.

## Goals

1. A Spec that legitimately declined the QA gate can be archived by the tool.
2. A partial verdict whose blocked rows are all declared reaches archive, as the
   command's own help already promises.
3. A repository's existing release tag form is discovered rather than prescribed.
4. A refusal that Roundfix cannot accommodate names what the repository would have
   to change and why, instead of implying the repository is wrong.

## User Stories

1. As a maintainer of a Spec whose graph declined the QA gate with a recorded
   reason, I want to archive it with the tool, so that I do not archive by hand
   around the preflight that exists to protect archiving.
2. As a maintainer of a Spec whose QA verdict is partial with only declared
   blocked rows, I want archive to accept it, so that the documented path is
   reachable rather than only described.
3. As a maintainer of a repository whose release tags carry no prefix, I want the
   release plan to read my tags, so that following it does not publish my release
   as a prerelease.

## Core Features

1. **A declined gate is an equivalent verdict.** Archive preflight accepts a graph
   that declined the QA gate with a non-empty reason, recording that reason in the
   archival artifact instead of requiring a QA Report. The contract already
   obliges the reason to exist and to be specific; it is the evidence.
2. **A declared-only partial settles its gate Task.** The QA Task settles as
   completed when the partial verdict's blocked rows are covered only by declared
   unreachable acceptance, so the two conditions archive states can both hold.
3. **The release tag form is discovered.** Release planning detects the stable tag
   pattern from the repository's existing tags rather than requiring one. With no
   tags at all, it asks for a base without prescribing a specific form.
4. **Consistency across the commands that read these shapes.** Whatever accepts a
   declined gate at archive accepts it wherever else the shape is read, so the
   contract does not hold in one command and not another.

## User Experience

Archiving a Spec that declined the gate succeeds and prints the recorded reason
where a QA verdict would otherwise be named, so the archive record shows why no
gate ran. A release plan in a repository with unprefixed tags names the tag it
found and the next version in the same form. A repository with no tags at all
gets a request for a base version that does not name a format.

## Non-Goals / Out of Scope

- Changing what a QA gate proves, when it may be declined, or the reason contract
  behind a decline.
- Changing verdict semantics or the typed blocked-cause counts.
- Supporting tag forms beyond the presence or absence of the `v` prefix.
- Changing release execution, publication, or the changelog.
- Auditing whether any particular Spec's decline was justified.

## Success Metrics

- A Spec that declined the gate, in a repository this Spec did not build, archives
  with the tool and records its reason. Source: a fleet repository whose Run
  reached Clean with every Task completed and which could not be archived on
  2026-08-10.
- A partial verdict with only declared blocked rows reaches archive without a
  maintainer override.
- A release plan run in a repository whose eight release tags carry no prefix
  produces a next version in that same form.
- No command that reads a declined gate refuses it while another accepts it.

## Decisions

- The decline reason is the evidence recorded at archive, rather than a synthetic
  QA Report standing in for one. The contract already requires the reason to be
  specific.

## Open Questions

- Whether the tag form is inferred per run from the newest stable tag or recorded
  once in configuration. Inferring is the default until answered, because it
  requires nothing of a repository adopting the tool.
- Whether a repository with a mixed tag history — some prefixed, some not — is a
  supported state or a refusal that names the ambiguity. The default is to refuse
  and name it, because guessing picks a form that silently breaks a consumer.
