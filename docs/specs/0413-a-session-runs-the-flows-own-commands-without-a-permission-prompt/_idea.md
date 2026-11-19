---
spec: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
created: 2026-11-19
status: promoted
leverage: strategic-bet
---

# A session runs the flow's own commands without a permission prompt

## Overview

Roundfix publishes one closed list of the shell commands its flow tells a session to run. A command shows that list in the rule syntax of the harness the session uses and, only when the user says yes, writes it into the project's own permission file for that harness. From then on the harness stops asking before each of those commands, in the session that drives the flow and in the Author Subagents it starts.

It is for the person who drives the flow from a supervised session of Claude Code, Codex or OpenCode, and for the Agent of that session. A session that skips every prompt gains nothing from it.

It is valuable because the flow is made of Roundfix commands, and a harness that asks about each one stops the flow while the person is away. The skills name 35 `roundfix` subcommands in 447 places, against 29 mentions of `git` and 5 of `gh` ([measurement](references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md)). Today the choice is to answer every prompt or to switch every prompt off. This Idea adds a bounded middle that a project can read before it accepts it.

V1 covers the session and its authors, not the Agents of a Run. It writes for Claude Code and Codex, reports for OpenCode, adds one line to the Doctor Command and removes what it wrote on request. It is a bet: its worth depends on people choosing a short reviewed list over no prompts at all, and on three harness formats that change.

## Problem

The maintainer's request of 2026-11-19, from [his answers](references/2026-11-19-maintainer-answers.md): "Outra coisa que normalmente trava o fluxo é o uso de comandos bash, se houver algo que possa ser feito para resolver isso tb, uma inclusão de autorização permanente seria muito bom. Mas depende do harness."

It does depend on the harness, and on which session runs the command. Line numbers are of this repository at `145d6c7610b77c4345e8a4eb0922bbe3857b4623`.

- In a Run nobody is asked. The Daemon starts each Agent Session with the answer already given: `--approve-all` for a session that may write, and reads approved with every other request denied for a session that may only read (`internal/agent/acpx_runner.go:1907`, `:1909`). What limits a Run is the sandbox, which has its own switch and its own decisions (ADR-0011, ADR-0293).
- In the session that drives the flow, the harness asks. Claude Code asks before a command that no rule allows and that it does not class as a read. Codex runs a command inside its sandbox and asks before one leaves it. A `roundfix` command writes Roundfix Home, which is outside the repository, and reaches GitHub, so the author expects it to be such a command; no Codex session was measured for this Idea. OpenCode's base policy allows every action, so it asks only where a project's configuration says so.
- An Author Subagent "inherits the parent conversation's permission rules" in Claude Code, so it meets the same prompts as the session that started it.

Roundfix does nothing about it today. It writes the Author Subagent files and the skill links into a harness's directory and no settings file, no Doctor check reads the harness of the person's session, and no page of the user guide mentions the prompt ([measurement](references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md), section 2).

The workarounds are three, and each fails in its own way. Answering each prompt needs the person there, which is the stop the maintainer describes. Telling the harness to remember an answer builds a list one command at a time, in a file of the person's own, that no teammate and no fresh clone receives. Switching every prompt off removes the stop and every other look with it.

One relayed report may describe the same stop from the other side. The Conexus team said that "A command sometimes stops waiting for an authorization and the calling agent never gets an answer, so it looks hung" (`docs/backlog/2026-11-11-the-conexus-team-saw-a-false-status-a-silent-authorization-wait-and-a-black-screen.md:18`). Nobody reproduced it, and the entry does not say which authorization.

### Market Data

- People approve almost every prompt. "Over time that leads to approval fatigue, where people stop paying close attention to what they're approving", and "in practice users accept 93% of them anyway" (Anthropic, 2026-03-25). The figure is the vendor's own. It says a prompt for a routine command protects little, and it does not say a standing rule protects more.
- A list of allowed commands is an attack surface. "Command-name allowlists don't work for commands with complex option semantics. git, curl, tar, rsync—any tool with flags that change behavior from read-only to write-capable can't be safely allowlisted by name alone" (Pillar Security, 2026-07-20, on a `git show --output` that led to code execution in Codex). In Claude Code, "I discovered 8 ways to execute arbitrary commands in Claude Code without user approval" (GMO Flatt Security, 2026-01-12, CVE-2025-66032).
- The harness says so itself: "Bash permission patterns that try to constrain command arguments are fragile" (Claude Code documentation).
- Each harness applies a project's rules only to a folder the person trusted. In Claude Code they apply "only after you accept the workspace trust dialog for that folder", and in Codex "Project-local rules under <repo>/.codex/rules/ load only when the project .codex/ layer is trusted".
- In Codex an allow rule is wider than a missing prompt: it means "Run the command outside the sandbox without prompting."
- One other command-line tool proposes allow entries for its own commands at setup and asks first: "Add apm commands to Claude allow list?". Its list also adds `Bash(python3 *)` in a Python project, which is the kind of entry this Idea forbids.

## Core Features

The closed list is called the Flow Command List in this Idea. The name is a candidate for the glossary, which the PRD settles.

| #   | Feature | Priority | Description |
| --- | ------- | -------- | ----------- |
| F1  | The Flow Command List | Critical | Roundfix ships one list, versioned with the binary, of the commands its owned skills tell a session or an author to run: the `roundfix` subcommands, and the Git and GitHub reads the flow needs. Every entry names a command down to its subcommand and never tries to constrain its arguments. |
| F2  | What the list never holds | Critical | A shell, an interpreter, a package manager or a build tool; a wildcard in the place of the command; a command that writes through Git or GitHub, such as a commit, a push, a merge or a Pull Request; `gh api`; the commands the Roundfix skill already tells a session to ask the user about first; and the command that applies the list. |
| F3  | Show | Critical | A command prints the list in the rule syntax of the harness it runs under or of the one the user names, the file it would write, what F2 leaves out and why, and what the entries do not protect against. It writes nothing. |
| F4  | Write on the user's yes | Critical | With an explicit confirmation the same command merges the entries into the project's permission file of that harness. It keeps every other byte of the file, adds nothing twice and sets no permission mode, approval policy or sandbox mode. It never writes a file outside the repository. |
| F5  | The yes is the user's | Critical | Because F2 keeps the applying command off the list, a harness that asks will ask the person before it runs. The Roundfix skill tells a session to show the list and put one question to the user before applying, in every permission mode. |
| F6  | Each harness, as it is | High | Claude Code: allow entries in the project's settings file, which also reach the Author Subagents. Codex: a rules file in the project's `.codex` layer. OpenCode: a report that its base policy already allows shell, or that the project's own rules ask or deny, and no write in V1. The output says that a rule applies only after the person trusts the folder in that harness. |
| F7  | The Doctor Command says where a project stands | High | For each harness whose files the repository holds: applied and current, applied from an older list, not applied, not needed, or not readable. Each state names the next command. Nothing is written. |
| F8  | The list stays true | High | A contract of the repository fails when an owned skill names a `roundfix` subcommand that is neither on the list nor among the named exclusions. After an upgrade that changes the list, the Doctor Command and the `update` Skill Mode report the difference, and applying it takes a new yes that shows only what changed. |
| F9  | Remove | Medium | The command removes exactly the entries Roundfix wrote and leaves the rest of the file as it was. |
| F10 | The permission files are tooling configuration | Medium | The files this Spec writes become Governed Paths, so no Task changes one without a grant that names it. The user guide gains a page that says what the list is, what it is not, and how each harness reads it. |

## What the user sees

An illustration of the output of F3 in a repository opened in Claude Code. The PRD and the TechSpec settle the command's name, its flags and the exact entries.

```text
Flow Command List for Claude Code — 31 entries, from Roundfix <version>

Would add to .claude/settings.json, under permissions.allow:
  Bash(roundfix follow-up *)
  Bash(roundfix phase *)
  Bash(roundfix spec *)
  ...

Left out, and still asked about:
  the command that applies this list
  roundfix upgrade, roundfix gc, ...   the skill asks you before it runs them
  git commit, git push, gh pr merge    a push and a merge stay behind the Spec's authorization

What this does not do:
  It is not a security boundary. A roundfix command runs what this repository
  declares: its Verification, its bootstrap and the commands of its Task files.
  The rules apply after you trust this folder in Claude Code.

Nothing was written. To write these entries, run the same command with its
apply flag.
```

The number 31 and the entries are invented. The three blocks are not: what would be written and where, what stays asked about, and what the list does not protect against.

## KPIs

| KPI | Target | How to measure |
| --- | ------ | -------------- |
| `roundfix` subcommands the owned skills name that are on the list or among its named exclusions | 100%, from no list today and 35 subcommands named | The contract of F8, which reads the skills and the list in the repository's Verification. |
| Entries F2 forbids | 0 in every harness's syntax | A test that renders the list for each harness and checks every entry against F2. |
| Bytes changed outside Roundfix's own entries | 0 across apply, a second apply and removal | Tests over permission files that already hold a user's rules, other keys and an unusual layout. |
| Prompts for a listed command in a supervised session | 0 in a walk through the status view, a phase start and a Spec Consistency Check, in Claude Code and in Codex, after the list is applied and the folder trusted | A walk a person records in each harness when the Spec is delivered, kept as a reference document. A Run cannot perform it. No count of today's prompts exists. |
| Projects behind the installed list that are told so | 100% | A test that applies an older list and reads the Doctor Command's line. |
| Reports of a flow stopped by a prompt for a listed command | 0 in the 30 days after a repository applies the list | At each cycle close, count the Backlog Entries, Findings and messages that name that cause. Today's evidence is the maintainer's statement and one relayed report. |

## Feature Assessment

| Criteria        | Question                                            | Score |
| --------------- | --------------------------------------------------- | ----- |
| Impact          | How much more valuable does this make the product?  | Strong — it removes the prompt from the commands that make up most of the flow, and it gives a project a third choice between answering every prompt and switching all of them off. |
| Reach           | What % of users would this affect?                  | Maybe — only a session that keeps prompts on. Nobody has counted how many of the sessions that drive the flow do. The maintainer's statement and the Conexus report are the evidence that some do. |
| Frequency       | How often would users encounter this value?         | Must do — a flow runs `roundfix` commands from its first step to its last: the skills name them in 447 places. |
| Differentiation | Does this set us apart or just match competitors?   | Maybe — every harness can remember one answer, and one other tool proposes entries for itself. A closed list with stated exclusions, checked against the skills and reported when it falls behind, is the part the sweep did not find. |
| Defensibility   | Is this easy to copy or does it compound over time? | Maybe — the list is text. The contract that keeps it equal to what the skills ask for is what a copy would lack. |
| Feasibility     | Can we actually build this?                         | Strong for Claude Code and Codex, whose formats and trust rules are documented. Maybe for OpenCode, whose documentation shows two shapes of the same setting, and for proof: the absence of a prompt needs a real harness and a person. |

The leverage type is strategic bet. The verdict is to proceed with V1 as bounded. Two results would change it: no project applies the list because its sessions already run with every prompt off, which makes the Doctor line the only part worth keeping, or a harness changes its rule format faster than releases follow, which makes Roundfix the author of rules that no longer mean what it shows.

## Trade-off Insights

- **Recommended approach:** the maintainer's two answers. The list is closed and shown before anything is written. The write is the user's decision, lands in the project's file and is the project's to review. A push and a merge stay where they are, behind the Spec's authorization.
- **Key trade-offs:**
  - A missing prompt against a missing look. Each entry removes a moment in which a person could have seen a command. The vendor's own figure says people approve 93% of those moments without change, so the look is weak. It is still the only one a supervised session has.
  - One entry per subcommand against one entry for all of `roundfix`. A single entry never falls behind and would also cover the commands the skill asks about first and the applying command itself. One entry per subcommand needs F8 and can exclude by name.
  - Exclusion by subcommand against exclusion by flag. A rule that tries to tell one flag from another is what the harness documentation calls fragile. So a subcommand is on the list or off it, and a dangerous flag keeps the guard it has inside Roundfix.
  - The project's shared file against the person's own. The shared file is reviewed in a Pull Request and reaches teammates and fresh clones. One person's yes then becomes the team's, once each teammate trusts the folder.
  - Git reads by name against the harness's own judgment of a read. The maintainer's answer includes Git and GitHub reads. Published attacks went through exactly such names, and each harness already runs some reads without asking.
  - Writing a harness's file against printing a snippet. ADR-0371 refused a merge into the same files as something the Baseline does on its own at every update. Here a person asks, sees and confirms, and the maintainer chose the write.
- **Risks identified:**
  - The list is not a boundary. A `roundfix` command runs what the repository declares, and the Spec Consistency Check with `--run-verification` runs the commands of Task files. A session that may edit files without asking can therefore reach any command through a listed one.
  - Text a session reads can steer a listed command. The skills already tell a session to treat such text as data, and the list gives an injected instruction fewer prompts to pass.
  - An Agent in a session with every prompt off could apply the list by itself. F5 and the tooling-authority clause forbid it, and nothing mechanical stops it there.
  - A rule that is written and not applied. Until the person trusts the folder, the harness ignores the project's allow rules, and Roundfix cannot see the trust state. The person sees prompts and believes the list failed.
  - Harness formats move. OpenCode documents two shapes of this setting, and a Codex rules file is a small program in a language of its own. A writer that lags writes rules a harness ignores or reads differently.
  - Making the files Governed Paths changes what a Task of any adopter may touch. A repository that already edits its Claude Code settings in a Spec will need a grant for it.
  - The flow still stops elsewhere: at a commit or a push the session makes itself, and at the ai-memory tools, which are not shell commands.
- **Dissenting view:** an approval prompt is not a safety control, and a list of allow rules is not one either. The honest product may be to say so and to recommend a sandboxed session with prompts off, where the boundary is the operating system's and no list is kept. The opposite dissent also stands: with a commit and a push left out, a flow that plans a Spec still stops at its first Git write, so the list removes many prompts and not the one that matters when nobody is there.
- **Stretch goal (V2+):** a Roundfix command that publishes a plan under the Spec's authorization, so the session needs no Git write of its own. Entries for the ai-memory tools the flow uses. A sandbox profile for each harness that lets a `roundfix` command reach Roundfix Home and GitHub without leaving the sandbox. The same list for the user's own files, on the user's request.

## Opportunity Scan

The original idea is bounded by the sessions that keep prompts on, and by the commands it leaves out on purpose. If it works fully, a supervised session reads status, records phases, checks plans and follows a delivery without one prompt, and still asks before a Git write.

- **More ambitious: one door.** Every step the flow needs becomes a `roundfix` command, the publication of a plan included, so the list is the whole flow and the session runs no Git or GitHub command of its own. Scale: large. Impact Must do, Reach Strong, Frequency Must do, Differentiation Strong, Defensibility Strong, Feasibility Maybe: each new command is its own design, and a command that pushes needs the authorization rules the Delivery Queue has.
- **Simpler: print and explain.** Roundfix prints the entries and the guide says where to paste them. Scale: small. Impact Maybe, Reach Maybe, Frequency Strong, Differentiation Pass, Defensibility Pass, Feasibility Must do. The maintainer saw this option and chose the write. Nothing would tell a project that its pasted list fell behind unless the Doctor line came with it.
- **Adjacent: a refusal that names itself.** When a command of the flow is refused or never starts, the session reports which permission was missing and where to grant it, in place of looking hung. It answers the Conexus report directly. Scale: small. Impact Strong, Reach Strong, Frequency Maybe, Differentiation Maybe, Defensibility Maybe, Feasibility Maybe: Roundfix does not run while a harness holds a prompt, so only the skill's text can carry it.

**Chosen direction:** original — the maintainer's answers of 2026-11-19 describe it. It takes from the adjacent idea the Doctor line, which says before the flow starts what would be asked. The single door is the stretch goal, starting with the publication of a plan.

## Integration with Existing Features

| Integration point | What it gains or how it is used |
| --- | --- |
| `roundfix skills install` | Unchanged. It keeps writing the skills and the Author Subagent files, and never the list: a permission is asked for by name, not installed along the way. |
| The Author Subagent files under `.claude/agents/` | Unchanged. In Claude Code a subagent takes the session's rules. |
| The Doctor Command (`CONTEXT.md:1205`) | One more check, a warning that never fails a Run, with its next command. It still changes nothing. |
| The Update Check and the `update` Skill Mode | After an upgrade that changes the list, the mode's report names the command. The mode never applies the list by itself. |
| The Roundfix skill | Its setup reference says how a session shows the list, that it asks the user before applying, and what to tell the user when a command of the flow is still asked about. |
| The authorization record | Unchanged. A push, a Pull Request, a merge and a release stay behind it, and none is on the list. |
| Full-Access Mode and the ACP Runtime (ADR-0011, ADR-0293) | Unchanged. A Run is out of this Spec. This Spec also keeps what those ADRs decided for the user's own files: Roundfix does not edit them. |
| ADR-0371 | Not contradicted. It refused a merge into a harness's project file as the Baseline's own mechanism, "which is tooling configuration the core clauses forbid changing without express authorization". Here the user gives that authorization, for entries he was shown. The PRD's ADR says so. |
| Tooling authority and Governed Paths (`internal/speccheck/governed.go`) | The harness files this Spec writes join the governed patterns. |
| Project Config | No key. The harness's own file is the record that a project applied the list, so an older binary reads the same repository. |
| The user guide | A page for the list, and a row in the Doctor Command's table. |
| `CONTEXT.md` | The names the PRD settles, starting from the Flow Command List. |

## Size and Shared Files

By the sizing criterion of the Spec routing guide this is one Spec. It holds one context: what a harness asks before it runs the flow's own commands in the session that drives the flow. The list, the command that shows and writes it, the Doctor line and the contract complete one another. A list nobody can apply is a page of documentation, a write without the contract goes stale at the next subcommand, and a Doctor line without the list has nothing to compare. The work fits four implementation Tasks and the QA gate, each large enough to earn its verification cycle:

1. the list, its exclusions and the two contracts over it;
2. the command that shows, writes and removes, for Claude Code and Codex, with the report for OpenCode;
3. the Doctor Command's check, the difference after an upgrade and the governed patterns;
4. the Roundfix skill's reference, the user guide, the record of command surfaces and the names in `CONTEXT.md`.

This Spec is the second part of one request. The first part, what a session does with an open authoring question when nobody is present, is Spec 0412. The two share no component and no contract, and the maintainer chose "Duas Specs (Recommended)".

The Conexus Backlog Entry is related and is not adopted. It holds three symptoms, of which one may be this stop. Its `## Where` says "Unknown", nobody reproduced it, and it waits on reproduction steps. By the same criterion a shared component or a shared cause would have to be shown, and neither is. Adopting the entry would also bring a false Run status and a black screen into this Spec. If the steps show that the wait is a harness prompt for a `roundfix` command, the list and the Doctor line answer what that entry expects, and its Triage says so then.

This Spec is delivered after Spec 0412. These files can meet:

- The owned Roundfix skill: its version fields, `SKILL.md` and its setup reference. Specs 0409, 0410, 0411 and 0412 also change this skill.
- `internal/cli/health.go` and `docs/user-guide/commands/doctor.md`. Spec 0410 adds checks to the Doctor Command, and the open Backlog Entry about a dead health probe names the same file.
- `internal/speccheck/governed.go` and the tests that pin its patterns.
- `docs/references/behavior-surfaces.json`, which Specs 0410 and 0412 also change, and `docs/references/coverage-record.json` when a test is added or renamed.
- `skills/testdata/owned-skill-versions.json` and the `skills/` mirror.
- `CONTEXT.md`.

## Out of Scope (V1)

- **The Agents of a Run** — the Daemon already answers their permission requests, and the sandbox of a Run has its own switch and its own ADRs. The maintainer chose "Sessão Supervisor e autores (Recommended)".
- **The user's own harness files** — nothing under the user's home directory is read for its rules or written. The project's file is the only one a project can review.
- **A permission mode** — the command never sets a default mode, an approval policy or a sandbox mode, and never recommends switching prompts off.
- **A commit, a push, a Pull Request and a merge** — they stay asked about, and a delivery's push and merge stay behind the Spec's authorization, as the maintainer answered.
- **Shells, interpreters, package managers and build tools** — no entry for `bash`, `sh`, `python`, `node`, `go`, `make`, `npm` or the like, with or without arguments. The repository's Verification is reached only through the `roundfix` command that runs it.
- **The ai-memory tools and every other tool that is not a shell command** — the request is about shell commands. They are in the stretch goal.
- **A hook that approves commands** — a harness's hooks are the subject ADR-0371 already decided, and a list is readable where a hook is code.
- **A write by the Baseline or by installation** — `roundfix baseline update` and `roundfix skills install` never apply the list.
- **A security boundary** — the list removes prompts. Isolation is the sandbox's work and the user's choice.
- **The three symptoms of the Conexus Backlog Entry** — they wait on their reproduction steps.
- **Harnesses other than the three** — Cursor and others have no session role in the flow today.

## Idea-stage judgment

- item: the request in the session, recorded in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md)
- judged-by: nobody (the user asked for the Idea stage)
- verdict: not judged
- stage: ran
- shape: changed — the standing authorization became a closed list that is shown first and written into the project's own harness file only on the user's yes, for the session and its authors and not for a Run, with stated exclusions, a line in the Doctor Command and a removal.

The judge was not run. The maintainer chose the Idea stage for the whole request on 2026-11-19, "Spec 0412 com Ideia agora", and then the split that made this Spec, "Duas Specs (Recommended)". The `shape` line is the author's and waits for the approval of this draft.

## Decisions

Every answer is the maintainer's, as recorded in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md). The text in quotes is the option he chose, verbatim. The text after the dash is that option's description as the operator recorded it.

- One Spec or two (2026-11-19, ~01:10 BRT): "Duas Specs (Recommended)" — 0412 fica com (A) preferência para pergunta aberta de autoria; (B) autorização permanente de comandos de shell vira a 0413, com Ideia própria partindo desta pesquisa, entregue depois da 0412.
- The target of V1 (2026-11-19, ~01:10 BRT): "Sessão Supervisor e autores (Recommended)" — o pedido de permissão na sessão que conduz o fluxo (Claude Code, Codex, OpenCode) e nos subagentes de autoria.
- The form (2026-11-19, ~01:10 BRT): "Mostra a lista e grava com seu sim (Recommended)" — comando mostra a lista fechada (comandos roundfix e leituras git/gh, sem curinga de shell nem interpretador) e grava no arquivo de permissões do projeto só com confirmação; push e merge continuam atrás do _authorization.md.

The author's own, made in this stage:

- The Idea of this Spec asked the maintainer no question of its own. Its two questions were asked with those of Spec 0412, from the same research, and the open points below keep their defaults until he approves the draft.
- The applying command is among the exclusions. It is what makes the user's yes something a harness asks for, and not only something a skill requests.
- The list excludes by subcommand and never by flag, because each harness matches a command from its start and the documentation of one calls argument patterns fragile.
- No Project Config key records the yes. The harness's file is the record, and a new key would make every older binary refuse the repository (`internal/config/config.go:1910`).
- OpenCode is reported and not written in V1. Its base policy already allows shell, and its documentation shows two shapes of the setting.
- The files this Spec writes become Governed Paths. A file that grants commands is tooling configuration, and today no pattern covers it.
- The Conexus Backlog Entry is named as related and not adopted, for the reasons in `## Size and Shared Files`. No other open Backlog Entry or unresolved Finding shares this context.
- The author read the rules that touch this Idea: ADR-0011, ADR-0293, ADR-0371 and the tooling-authority clauses. None blocks the goal: each forbids a write nobody authorized or a write to the user's own files, and this Idea does neither.
- No ADR is written at this stage. The numbers 0444, 0445 and 0446 are reserved for the PRD and the TechSpec.

## Open Questions

- Which `roundfix` subcommands are on the list? Default: every subcommand the owned skills name, one entry each, except those the Roundfix skill already tells a session to ask the user about before running, and the applying command. The alternatives are one entry for the whole `roundfix` command, or only the subcommands that change nothing.
- Which Git and GitHub reads are on the list? Default: only a read the harness would ask about, in a form that cannot write, and never `gh api`. A Git command whose options can write a file, as `git show --output` does, is left to the harness's own judgment of a read. The alternatives are every read the skills name, by name, or none.
- Which Claude Code file takes the entries? Default: the shared project file, which is committed and reaches teammates and fresh clones once each person trusts the folder. The alternative is the personal project file, which nobody else receives.
- Is OpenCode written in V1? Default: no, it is reported. The alternative is to write the shape its current documentation shows.
- What is the command called, and how is the yes given? Default: a preview unless an apply flag is passed, with the skill asking the user through the question tool before it passes the flag. The TechSpec names the command and the flags.
- How does Roundfix know which entries are its own in a file that takes no comments? Default: an entry is Roundfix's when it equals an entry of the current list or of an earlier one, and the removal says what it will remove before it does. The TechSpec decides.
- How is the absence of a prompt proved? Default: tests prove the entries, the bytes and the Doctor line. The walk in a real harness is recorded once by a person at delivery, under `docs/references/`, and the PRD declares it as acceptance a Run cannot reach.
- Do the governed patterns reach adopters at once? Default: yes, with the release. A Task that edits one of those files then needs a grant that names it.
- Can a subagent that runs in the background show a prompt at all? Not measured. Default: the PRD stage measures it in Claude Code before it writes the requirement about authors.
- Does the session tell the user that a rule waits on folder trust? Default: the command's output and the Doctor line say it, and Roundfix reads no trust state.
- Which names does the glossary take? Default: the PRD's Glossary settles them, starting from Flow Command List.

## Sources

Adopted by this Spec (see [`references/_index.md`](references/_index.md)):

1. [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md), the maintainer's request and answers of 2026-11-19, as the operator recorded them. Spec 0412 keeps a copy of the same record.
2. [`references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md`](references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md), recorded by the author: the count of the commands the skills name, what Roundfix writes and answers today, the passages of each harness's documentation and the published hazards.

Related and not adopted: `docs/backlog/2026-11-11-the-conexus-team-saw-a-false-status-a-silent-authorization-wait-and-a-black-screen.md`, which stays open with its Waiting Declaration.

This project's ai-memory messages were listed on 2026-11-19 and none was pending. The Research Subagent did the codebase sweep and the web sweep. The author did the Secondbrain sweep, because it needs a command, counted the commands, and re-read the passages quoted in this file.

Secondbrain (read on 2026-11-19; `wiki/index.md`, then `qmd query --all --files --min-score 0.3`):

- `wiki/sources/digest-roundfix-claude-agent-sdk-permissoes-e-anexos-2026-10-04.md` — an earlier measurement of this repository: Claude Code approves by itself the commands it classes as reads and asks about the others, and an allow list "só pré-aprova". It is why the second open question leaves reads to the harness.
- `wiki/concepts/agent-harnesses.md` — "approval prompt não é controle de segurança (93% são aprovados)" and "denylist é sugestão; sandbox é parede". It is the dissent of this Idea, and the reason the output of F3 says what the list is not.
- `wiki/concepts/claude-code-hooks-e-automacao.md` — a scoped permission rule blocks a call and leaves the tool in place, and Git writes are the commands a guard blocks first. It supports leaving every Git write off the list.
- `wiki/concepts/claude-code-operacao-produtiva.md` — "manter aprovação humana para envios externos ou ações difíceis de desfazer". It supports the same exclusion.
- `raw/web/akita/2026-01-10-akita-ai-agents-garantindo-a-protecao-do-seu-sistema.md` — a practitioner who keeps prompts on, allows listing commands in his own settings file and puts the agent in an operating-system sandbox. It is the user this Idea serves, and it shows the list and the sandbox as two separate things.

External sources (the official pages were read by the Research Subagent and re-read by the author on 2026-11-19; the others were found through Exa):

- Claude Code documentation, https://code.claude.com/docs/en/permissions, https://code.claude.com/docs/en/settings and https://code.claude.com/docs/en/sub-agents — the allow list, its merge across files, the trust rule, the order of rules and what a subagent inherits. They set F4, F5 and F6.
- Codex documentation, "Rules", https://learn.chatgpt.com/docs/agent-configuration/rules — the prefix rule, the project layer and what `allow` means. It set F6 and the warning that an allowed command leaves the sandbox.
- OpenCode documentation, "Permissions", https://opencode.ai/docs/permissions/ — the base policy that allows every action. It is why OpenCode is reported and not written.
- Anthropic, "How we built Claude Code auto mode", https://www.anthropic.com/research/claude-code-auto-mode — approval fatigue and the 93% figure. It supports the problem and limits what a prompt is worth.
- Pillar Security, "GitPwned: Allowlist to RCE", https://www.pillar.security/blog/gitpwned-allowlist-to-rce — a Git read by name that wrote a file. It challenged the Git reads of the maintainer's answer and is the second open question.
- GMO Flatt Security, "Pwning Claude Code in 8 Different Ways", https://flatt.tech/research/posts/pwning-claude-code-in-8-different-ways/ — eight ways past a command check. It supports a short list with no entry that reaches a shell.
- The source of `apm-cli` 0.1.40, https://docs.rs/apm-cli/0.1.40/src/apm/cmd/init.rs.html — a tool that proposes entries for itself and asks first. It shows the form is known, and its interpreter entry shows what F2 must refuse. Nothing of it is copied.

The digest of this research for the project that owns the Secondbrain was not sent by the author, who sends nothing in this cycle. The operator has it in the author's report.
