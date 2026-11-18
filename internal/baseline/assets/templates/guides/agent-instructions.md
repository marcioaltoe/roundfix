# Agent instructions

This setup-owned guide defines the portable baseline. Repository authors write
every project-specific rule or variation, stricter or not, only in
`docs/agents/specific-repository.md`, never in this guide or in a root
instruction file.
The root instruction file is `AGENTS.md`; `CLAUDE.md` is a symbolic link to it, which `roundfix baseline` creates and keeps, never a file of its own.
The selected repository Verification is {{verification.gate}}.
The selected incremental Verification is {{verification.incremental}}.
{{branch.prefix}}

{{artifact.rules}}
