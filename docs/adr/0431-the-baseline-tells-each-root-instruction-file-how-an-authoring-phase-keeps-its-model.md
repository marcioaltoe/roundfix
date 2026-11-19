---
status: accepted
knowledge: true
created_at: 2026-11-18T05:40:00Z
updated_at: 2026-11-18T05:40:00Z
deprecated_at: null
superseded_by: null
---

# The Baseline tells each root instruction file how an Authoring Phase keeps its model

The Baseline's root block `root.spec-workflow`, which every built-in profile
writes into an adopting repository's `AGENTS.md`, gains a short mandatory
statement scoped to Claude Code: in Claude Code, a skill's `model` and
`effort` header applies only when a person types the skill; a skill an Agent activates runs on the session's model; to
hold an Authoring Phase to its declared model, the Agent hands the phase to
its Author Subagent, and when that model's usage limit ends the subagent it
follows the fallback `roundfix phase fallback` prints. The statement says it
does not apply to Codex or OpenCode, so `AGENTS.md` stays the one instruction
file and `CLAUDE.md` its link. No user-level file is
needed for the Baseline to work.

## Context

ADR-0401 measured on 2026-11-14 that Claude Code applies a skill header's
`model` and `effort` when the user types the skill and ignores them when the
Agent activates it (anthropics/claude-code#98898), and made the authoring
skills hand a feature's phase to its Author Subagent. That rule lives in the
skills, which an Agent reads only after it has activated one on the session's
model. On 2026-11-18 the maintainer asked that nothing depend on the user's
`~/.claude/CLAUDE.md`: "Tudo tem que funcionar por configurações no
repositório. [...] se for necessário a informação de que o cabeçalho só vale
com /skill isso deve ser adicionado ao agents.md e claude.md do repositório.
E quem deve fazer essa alteração/ajuste no baseline é o roundfix."

The Baseline writes its root blocks into `AGENTS.md`, and the adopters
measured on 2026-11-18 keep `CLAUDE.md` as a symbolic link to it. Asked
about the root files the same day, about 06:30 BRT, the maintainer set the
standard: "Temos que ter sempre um padrão igual para todos os projetos que
usam o baseline do roundfix. Então se for necessário instruções diferentes no
AGENTS.md e CLAUDE.md, vamos ter sempre dois arquivos e sem symlink já sem nem
dar opção ao usuário, é regra de uso do baseline. Se não for necessário e
formos manter arquivo AGENTS.md e o CLAUDE.md com symlink para o AGENTS.md,
vai ser sempre assim." He added: "o trecho que fala sobre o uso da skill afeta
de alguma forma … o uso das skills no Codex ou Open Code? Se não afeta, a
gente vai ter dessa forma: um AGENTS.md, que é o arquivo, e o CLAUDE.md
sempre o symlink". Codex and OpenCode read the same `AGENTS.md`, and their
models come from the Agent Selection Profiles, so a statement scoped to
Claude Code does not affect them, and the standard is one file with
`CLAUDE.md` its link. Spec 0406 makes that link the rule.

## Decision

**The block.** `root.spec-workflow` keeps its id and gains one bullet after
its pointer. The bullet opens with "In Claude Code," and ends by saying it
does not apply to Codex or OpenCode. It names no model, no ADR and no Spec: it states the
header rule with its upstream issue and the date it was measured, the path of
the Author Subagent, the command that generates it from
`skills.authoring.<harness>`, and `roundfix phase fallback`.

**Versions.** The template's version, the root block's version and the
`spec-workflow` module's version rise by one, the module's through its record
command. The derived files change only through `make baseline-digests`, and
this repository's `AGENTS.md` and Setup Manifest through the Managed Refresh.

**Adopters.** An adopting repository receives the bullet with its next
`roundfix baseline update`, which Roundfix performs; no adopter edits a root
file by hand, and no user-level file is read or written.

## Considered Options

- **A rule in the user's `~/.claude/CLAUDE.md`.** Rejected by the maintainer.
- **A new root block.** Rejected: the statement belongs to the Spec workflow,
  and extending the existing block keeps every profile's set of blocks.
- **A clause in a Baseline guide.** Rejected: the maintainer asked for the
  root files, and an Agent reads a guide only when its pointer leads there.
- **A separate `CLAUDE.md` that carries the statement alone.** Rejected: the
  statement does not differ per harness once it is scoped to Claude Code,
  and the maintainer's standard keeps one file while no instruction differs.
- **An unscoped statement.** Rejected: a Codex or OpenCode agent could read
  it as governing its own model, which the Agent Selection Profiles set.

## Consequences

- Every adopting repository's root instruction file carries the statement
  after one Managed Refresh, and this repository's after the Spec's.
- A Codex or OpenCode agent reads the statement as not applying to it.
- A repository whose `CLAUDE.md` is not yet a link to `AGENTS.md` is Spec
  0406's case.
- The statement names `roundfix phase fallback`, which ships in the same
  release.
