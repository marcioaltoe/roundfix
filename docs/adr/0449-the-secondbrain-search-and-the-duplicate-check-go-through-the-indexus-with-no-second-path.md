---
status: accepted
knowledge: true
created_at: 2026-11-20T15:34:29Z
updated_at: 2026-11-20T15:34:29Z
deprecated_at: null
superseded_by: null
---

# The Secondbrain search and the duplicate check go through the Indexus, with no second path

## Context

On 2026-11-20 the maintainer decided: "Agora, os testes devem ser com o
objetivo inicial de ser superiores ao que o qmd/secondbrain entrega no que os
projetos/jarvis usam. Se for, já podemos usar somente o indexus e eliminar o
qmd no mini e macbook." He authorized preparing the change in a Pull Request
while the test runs. The earlier command required a local installation and
index, could be slow, and the consultation attempted during authoring stopped
at a locked database. Keeping it in the guide after its removal would direct
Agents to a command the machines no longer have.

## Decision

The Secondbrain search after reading `wiki/index.md` and the advisory duplicate
check use the `pesquisar_acervos` tool of the `indexus` MCP server. No Baseline
clause, generated guide or owned skill names a second search tool. The search
opens a returned page with `abrir_conteudo`, using its returned `caminho`; the
duplicate check verifies that each returned page opens and names a strong
verified match so the knowledge owner extends what exists. A relevance label
or a rank alone never decides.

When the tool is absent, does not answer, or answers without a wiki page, the
Agent records that limitation and continues with the evidence it has. It
substitutes no other search tool and claims no consultation that did not
happen.

Both clauses keep their identities and mandatory force. The duplicate-check
clause keeps its replacement of the earlier research-capture clause, so an
adopter's update accounts every clause it held. Roundfix's code never calls
the Indexus and adds neither a Baseline decision nor a Doctor check for it:
the session provides the tool and its reach.

This decision takes effect with the release that carries it. On 2026-11-20
the Indexus project's side-by-side test found the Indexus superior in every
use the maintainer named, and the maintainer decided to remove the earlier
command at once. The project that owns the Secondbrain removes it from the
machines after the release, not as part of this repository's change.

## Considered Options

- Indexus first with the earlier command as a reserve: rejected by the
  maintainer's choice of one path; a reserve would preserve the installation
  and an instruction that its removal breaks.
- A Baseline decision declaring whether the Indexus is available: rejected
  because availability belongs to the session, not the repository.
- New clause identities: rejected because rewording keeps identity and force;
  new identities would require new Source Baseline accounting.
- Describing the search capability without naming the tool: rejected because
  an Agent needs the tool's name to perform the consultation.
- Naming one harness's joined tool identifier: rejected because that
  identifier differs by harness; server and tool names are stable separately.

This decision amends the consultation rule of
[ADR-0330](0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md).
The Secondbrain remains read-only, and ai-memory continues to carry handoffs
and cross-project messages.
