---
status: accepted
knowledge: true
created_at: 2026-11-14T11:00:00Z
updated_at: 2026-11-14T11:00:00Z
deprecated_at: null
superseded_by: null
---

# The Research Subagent reads and never runs a command or writes a file

The generated Research Subagent reads text nobody vouches for: web pages,
knowledge-base notes and files of the repository. Its tool list is therefore
closed by Roundfix: three tools that read the checkout, plus only the tools a
repository names from a set that holds no shell and no file-writing tool.

## Context

indexus wrote its research subagent by hand on 2026-11-13, and its pre-PR
review found a shell in the tool list beside web reading and memory reading:
an instruction injected into a page it read could run a command. Saying
"read-only" in the subagent's prompt does not remove a tool. Measured on
Claude Code 2.1.325 on 2026-11-14: a subagent whose file lists `Read, Grep,
Glob` could not create a file when asked to, while the same subagent without
a `tools` list wrote it through the shell. OWASP's LLM06:2025 Excessive Agency
asks to "Limit the extensions that LLM agents are allowed to call to only the
minimum necessary" and to "Avoid the use of open-ended extensions where
possible (e.g., run a shell command, fetch a URL, etc.)". The maintainer
decided that nothing assumes his machine: no search provider and no memory
server is built in.

## Decision

- **The list.** The generated file's `tools` holds `Read`, `Grep` and `Glob`,
  then each name of `skills.authoring.<harness>.research.tools` in the
  declared order.
- **What a repository may add.** `research.tools` accepts `WebSearch`,
  `WebFetch` and names that begin with `mcp__`. Any other name, `Bash`,
  `Write`, `Edit` and `Agent` included, is refused when the configuration is
  loaded, with the name and the accepted set.
- **No further subagent.** The list never holds `Agent`, so the Research
  Subagent starts nothing.
- **Sweeps that need a command stay with the author.** The authoring skill
  runs the command and sends the paths or the output for reading.
- **A listed MCP tool is the repository's judgment.** Roundfix cannot tell
  whether such a tool writes; the documentation says to list only tools that
  read.

## Considered Options

- **A read-only instruction in the prompt.** Rejected: the measurement shows
  the list is the boundary and the prompt is not.
- **Web tools in the fixed list.** Rejected: a tool that fetches a URL can
  carry out what the subagent read, so reaching the network is each
  repository's explicit choice.
- **A free tool list.** Rejected: it lets a repository give the subagent a
  shell, which is the defect this decision answers.
- **A list of forbidden tools.** Rejected: harnesses add tools, and a
  forbidden list is open on the side that matters.
- **A search provider or a memory server built in.** Rejected by the
  maintainer's audience decision.

## Consequences

- A repository that lists a web tool or an MCP tool gives the subagent a way
  to send out what it can read; the shell and the write tools stay out
  whatever it lists.
- A knowledge-base search that runs as a shell command is run by the author,
  and the Research Subagent reads the files it returns.
