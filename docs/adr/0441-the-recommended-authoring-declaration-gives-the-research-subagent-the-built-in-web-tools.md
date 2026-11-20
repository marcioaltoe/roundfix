---
status: accepted
knowledge: true
created_at: 2026-11-19T03:30:00Z
updated_at: 2026-11-19T03:30:00Z
deprecated_at: null
superseded_by: null
---

# The Recommended Authoring Declaration gives the Research Subagent the built-in web tools

The Recommended Authoring Declaration lists `tools: [WebSearch, WebFetch]`
under `research`, and Doctor and the Update Check name those tools when a
repository declares a research entry without a `tools` key. The generated
subagent's fixed list stays `Read, Grep, Glob`. This amends ADR-0415, whose
block held no `research.tools`.

## Context

ADR-0415 kept `research.tools` out of the block because this repository lists
MCP tools of its own machine and ADR-0402 leaves the network to each
repository. Its consequence was that an adopter pasting the block got "a
Research Subagent with no web tool until it lists one". On 2026-11-18 conexus
reported, in message 01a1f18f, that the generated `roundfix-researcher` holds
`Read, Grep, Glob` while its text says it reads the web and returns URLs, and
that its pre-PR review flagged the gap. Nothing Doctor or the Update Check
prints names the `tools` key, so conexus concluded that no key existed. The
Claude Code tools reference says of a subagent definition: "`tools` only: the
subagent gets only the listed tools." The maintainer chose on 2026-11-18 to fix
it in Spec 0411 ("Juntar à Spec 0411 (Recommended)").

## Decision

**The block.** The research entry of the Recommended Authoring Declaration
lists `WebSearch` and `WebFetch`, rendered as `tools: [WebSearch, WebFetch]`
after its `reasoning_effort`. It names no MCP tool. Pasting the block is the
repository's explicit choice of the network that ADR-0402 asks for, made where
the repository reads it.

**The generator.** `skills/authoring.go` keeps writing `Read`, `Grep` and
`Glob` followed by the declared tools. The accepted names stay `WebSearch`,
`WebFetch` and names that begin with `mcp__`; a shell or a write tool is still
refused when the configuration loads.

**The report.** When the declared `skills.authoring.claude.research` entry has
no `tools` key and the inspection of ADR-0401 finds no problem, Doctor prints
`authoring-models: warn` with the detail `research declares no tools, so the
Research Subagent cannot read the web` and the next action `add tools:
[WebSearch, WebFetch] under skills.authoring.claude.research in .roundfixrc.yml,
then run roundfix skills install --target project`. The Update Check prints the
same in one text line and adds `recommendedResearchTools` to `authoringModels`,
whose `state` stays `applied`. An empty list, `tools: []`, and any declared
list are the repository's choice and are not compared. Nothing is written, and
no exit code or top-level `next` changes. A problem the inspection finds keeps
its own report instead.

**The contract.** The repository contract test of ADR-0415 also requires this
repository's research tools to include every recommended tool; further tools
pass.

## Considered Options

- **Web tools in the generator's fixed list.** Rejected, as in ADR-0402: a
  repository could no longer keep the subagent off the network, and every
  adopter's generated file would change at once.
- **MCP tool names in the block.** Rejected, as in ADR-0415: they name tools
  of one machine.
- **Report every declared list that lacks a recommended tool.** Rejected: a
  list is a deliberate choice, and ADR-0181 holds that Roundfix does not argue
  with one on every check.
- **Rewrite the generated description when no web tool is declared.**
  Rejected: it changes the generated text of every adopter for a sentence the
  subagent's rule 3 already covers, since it reports what it could not read.
- **A snapshot date for the authoring recommendation.** Rejected: the report
  names the missing key directly, and the declaration's models are not
  compared (ADR-0415).
- **Write the tools into `.roundfixrc.yml`.** Rejected, as in ADR-0415: the
  checks never write.

## Consequences

- An adopter that pastes the block after the release gets a Research Subagent
  that can search and fetch the web, subject to the harness's own permission
  rules; on a provider without the server-side web search, `WebSearch` is
  simply unavailable.
- An adopter that pasted the earlier block sees a Doctor warning and an Update
  Check line until it adds the tools or writes `tools: []`.
- The Research Subagent can now send out what it reads through a URL it
  fetches, as ADR-0402 already says of any listed web tool.
