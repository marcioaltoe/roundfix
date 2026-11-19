---
spec: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
status: active
created: 2026-11-19
surfaces: [cli, docs]
type: feat
---

# A session runs the flow's own commands without a permission prompt

A supervised session of Claude Code or Codex asks the person before it runs a shell command that no rule allows, and the flow is made of `roundfix` commands, so the flow stops whenever the person is away. Today the choice is to answer every prompt or to switch every prompt off. After this Spec, Roundfix ships one closed list of the commands its flow asks a session to run, and one command shows that list in the rule syntax of a harness and, only when the user confirms what it showed, writes it into the project's own permission file for that harness. The harness then stops asking about those commands, in the session that drives the flow and in the Author Subagents it starts, and it still asks about everything else.

## Sources

1. [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md), the maintainer's request of 2026-11-19 in his own words, the answers of the Idea stage, the approval of the Idea and the authorization of this Spec, recorded by the operator. The request exists only in the session, so this record is the source; no Backlog Entry, Finding or ai-memory message is adopted and nobody asked for a notice.
2. [`references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md`](references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md), recorded by the author of the Idea: the count of the commands the owned skills name, what Roundfix writes and answers today, the passages of each harness's documentation and the published hazards.
3. [`references/2026-11-19-harness-rule-measurements.md`](references/2026-11-19-harness-rule-measurements.md), recorded by the author of this PRD: what Claude Code and Codex do with a listed command, an unlisted one, a listed read that carries an option that writes, a subagent and a session that skips prompts.

The approved [`_idea.md`](_idea.md) is the exploration this PRD starts from: its features F1 to F10 are this Spec. This project's ai-memory inbox was listed on 2026-11-19 and held no pending message, and `docs/findings/` holds no Finding. Six Backlog Entries are open at the authoring base. One reports three unreproduced symptoms from another team, among them a command that waits for an authorization in silence; it is related and not adopted, for the reason the Idea gives: its place is unknown, nobody reproduced it, and it waits on its reproduction steps. One names dead code in the file that holds the Doctor Command's check names, which this Spec also changes, and shares no cause with it. The other four are about the context an Authoring Phase starts with, the name of a documentation folder, a limit of the Spec Consistency Check's Verification reader and a deadline of the Baseline's analysis. None shares this Spec's context.

## Prerequisites

No behavior of this Spec depends on another Spec's code, and no Verification of it reads another Spec's artifact. It was authored against `abe2ebe1`, where Spec 0411 is delivered and Specs 0409, 0410 and 0414 are planned, and it is delivered last, after Specs 0414, 0409, 0410 and 0412, in the order the maintainer set. These files are shared, and each Task regenerates the shared records on the tree it runs on, with no version number written by hand:

- Spec 0410 adds two checks to the Doctor Command, in the same list and the same files this Spec adds one to, and changes the tests that pin that list, the Doctor guide, the Roundfix skill's `setup` and `modes` references, the command dispatch file, the usage guide and `CONTEXT.md`. This Spec's check comes after the checks Spec 0410 adds, and its Task moves every test that pins the list on the default branch when it runs.
- Spec 0410 also gives every ADR two header fields. The three ADRs of this plan carry the Knowledge Mark from the start, and none is amended.
- Spec 0414 changes the command dispatch file, the usage guide, the Skill Coverage Map and `CONTEXT.md`.
- Spec 0409 changes the Roundfix skill's `modes` reference and `CONTEXT.md`.
- Spec 0412 changes the Roundfix skill's `modes` reference, the configuration guide and `CONTEXT.md`. It adds a Project Config key and touches no permission file of a harness; this Spec adds no Project Config key.
- All four raise the Roundfix skill's version and change the owned-skill version record and the Behavior Surface Record.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted for a product entity; an entry of the list is identified by its own text, and a preview by a digest of what it shows. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call is added; the command reads and writes files of the repository, and the list holds no credential. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0444 (new) records when Roundfix writes a harness's file (ADR-0444: "a person ran the command that shows the closed Flow Command List for that harness, and then confirmed that exact preview"), ADR-0445 (new) records what the list holds (ADR-0445: "A subcommand is on the list or off it."), and ADR-0446 (new, written with the TechSpec) records how a confirmation and an entry are recognized (ADR-0446: "An entry is known by its text."). The earlier refusal to merge into a harness's project file stands and is not taken: it refused a write the Baseline makes by itself (ADR-0371: "which is tooling configuration the core clauses forbid changing without express authorization"), and here a person asks, sees and confirms. The user's own runtime files stay untouched (ADR-0011: "not by editing the user's runtime configuration files"; ADR-0293: "Editing `~/.codex/config.toml` was rejected again"). The governed set grows and never shrinks (ADR-0130: "The universal clause defines the governed class by kind"). The one question the skill puts before applying follows the question contract (ADR-0320: "Up to three, independent, in one pending call."). The skill and the guides describe the command in the same delivery (ADR-0256: "A release waits for the skills that describe a changed surface"), the authorization record carries the Release Ceiling (ADR-0349: "A Spec's Release Ceiling answers the release plan's approval"), and each new ADR carries its mark (ADR-0439: "The author marks a new ADR"). Five ADRs are listed to record why they do not apply. The two reports copy the form of an earlier report and read no user-level skill copy (ADR-0386: "None of them writes to those directories."; ADR-0405: "name each user-level copy they report with the command that removes it"). The Author Subagent is not changed: this Spec relies only on a subagent taking the session's rules (ADR-0401: "Author Subagent cannot ask the user, so the session gathers the decisions first"). The Research Subagent has no shell for a rule to allow (ADR-0402: "from a set that holds no shell and no file-writing tool"). ADR-0166, ADR-0179, ADR-0314, ADR-0338, ADR-0375, ADR-0412, ADR-0416 and ADR-0438 are listed only to record why they do not apply; the TechSpec relates each one. No accepted ADR is superseded or amended. Source: `docs/agents/domain.md`, `docs/agents/docs-layout.md`.
- Tooling authority: applicable — express maintainer authorization: his answer of 2026-11-19, "Mesmas do BE para as duas (Recommended)", which grants Specs 0412 and 0413 the cycle's authorization "Mesmas do BE (Recommended)" of 2026-11-18 and which the authorization record quotes with its provenance. It keeps the grant of every owned skill and of the Governed Paths each Spec declares. It grants no Makefile, workflow, Go module, lint or formatter change except for a correction that requires one, and the design needs none; `.roundfixrc.yml` is not granted and no Task changes it. It is not a confirmation of the Flow Command List for this repository: no Task writes a harness permission file here. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `internal/speccheck/governed.go`, `internal/speccheck/governed_repocontract_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- A supervised session, and the Author Subagents it starts, run the commands of the flow that are on the list without a permission prompt, once a project applied the list and the person trusted the folder.
- Nothing is written into a harness's file that the user was not shown first, and nothing but the entries of the list is written.
- The list stays equal to what the skills ask for: a command a skill names is on it or is left out by name, and a project whose applied list fell behind is told so.
- A project can leave: one command removes what Roundfix wrote and nothing else.
- The list never becomes a way around a prompt for a shell, an interpreter, a push or a merge, and every surface that shows it says it is not a security boundary.

## User Stories

1. As the person who drives the flow from a supervised session, I want the harness to stop asking before each `roundfix` command of the flow, so that the flow goes on while I am away.
2. As that person, I want to read the whole list, what it leaves out and what it does not protect against before anything is written, so that my yes is to something I saw.
3. As the maintainer of a repository, I want the entries in the project's own shared file, so that a teammate and a fresh clone receive them and review them in a Pull Request.
4. As a person who keeps prompts on for safety, I want a push, a merge, a shell and an interpreter to stay asked about, so that the list removes only the questions about the flow's own commands.
5. As the maintainer who chose every Git and GitHub read the skills name, I want the options of those reads that write a file, run a program or change GitHub to stay asked about wherever the harness can do that, so that my choice costs as little as it can.
6. As a user of Codex, I want to be told that an allowed command leaves the sandbox and that no guard can be written there, so that I confirm with that in view.
7. As a user of OpenCode, I want to be told whether my project's rules ask about shell commands at all, so that I know there is nothing to apply.
8. As a maintainer after an upgrade, I want the Doctor Command and the update report to say that my applied list is behind and to show only what changed, so that a new yes is a small read.
9. As a maintainer who changed my mind, I want one command that removes exactly what Roundfix wrote, so that my own rules and the rest of the file stay as they were.
10. As a maintainer who plans work through Specs, I want a harness's permission files to be Governed Paths, so that no Task changes what a session may run without a grant that names the file.
11. As a session that follows the Roundfix skill, I want one rule for when I may apply the list, so that I never write it on my own in any permission mode.

## Core Features

1. **The list is closed and ships with the binary.** The **Flow Command List** holds three kinds of entry: a `roundfix` subcommand, a Git read and a GitHub read. Each entry names a command down to its subcommand and never constrains its arguments. It names the installed `roundfix` by that name; a command run through a path, through a wrapper or with a variable set before it is not covered and stays asked about.
2. **Which `roundfix` subcommands are on it.** Every top-level subcommand that the binary registers and an owned skill names, one entry each, unless it is a **Named Exclusion**. A Named Exclusion is a command the skills name that the list leaves out on purpose, with its reason. Two reasons exclude a `roundfix` subcommand: it is the command that applies the list, or the Roundfix skill tells a session to run it, or one of its forms, only after it asked the user or the user asked for it. At the authoring base the second reason excludes `baseline`, `skills`, `profiles` and `release`. A subcommand is on the list or off it; no entry tells one flag from another, and a dangerous flag keeps the guard it has inside Roundfix.
3. **Which Git and GitHub reads are on it.** Every Git and GitHub command an owned skill names that, in the form the skill names it, changes no ref, index, working-tree file or configuration of the repository and changes nothing on GitHub. At the authoring base these are the Git subcommands `status`, `diff`, `diff-tree`, `show`, `grep`, `ls-files`, `check-ignore`, `rev-list`, `rev-parse`, `cat-file`, `merge-tree` and `archive`, and the GitHub commands `gh pr view` and `gh api`. Every other Git and GitHub command the skills name is a Named Exclusion: at the authoring base the Git subcommands `config`, `init`, `merge`, `mv`, `push`, `remote`, `repack`, `switch` and `worktree`, and `gh auth`. A word that follows `git` or `gh` in a code span and is no command, as in a quoted diagnostic, is declared as such and is on neither side.
4. **What the list never holds.** A shell, an interpreter, a package manager or a build tool, with or without arguments; a wildcard in the place of a command or of a subcommand; a command that writes through Git or GitHub, such as a commit, a push, a merge or a Pull Request command, whether a skill names it or not; and a Named Exclusion. A push and a merge of a delivery stay where they are, behind the Spec's authorization record.
5. **A read's dangerous options stay asked about where the harness can do it.** A **Guard Entry** keeps one option of a listed read asked about. The list holds one for each option the author knows to write a file, run a program or change GitHub through a listed read: writing the output to a file and running an external diff program for the Git reads that take those options, opening matches in a program for `grep`, writing to a file or reaching another host for `archive`, and choosing a method, adding a field or sending a body for `gh api`. Claude Code can express a Guard Entry, and there it wins over the entry it guards, in a subagent and in a session that skips every prompt as well. Codex cannot express one, so its rendering holds none and every surface that shows it says so. A Guard Entry adds a question and never removes one.
6. **Show.** `roundfix permissions` prints, for one harness, the file it would write, the entries and Guard Entries it would add, the entries it would remove because the list no longer holds them, the Named Exclusions with their reasons, what the list does not protect against, that the rules apply only after the person trusts the folder in that harness, and the digest of this preview with the command that applies it. With no harness named it prints the same for every harness whose files the repository holds. It prints text or JSON and writes nothing.
7. **Write only what a preview showed.** The same command writes only when it is given the digest a preview printed. When that digest no longer describes the list, the file or the harness, because the binary, the file or the list changed in between, it writes nothing, says that the preview changed and exits with the code that means a confirmation is required. There is no flag that writes without a digest.
8. **What is written, and what is kept.** For Claude Code the entries go into the project's shared settings file, under its allow list, and the Guard Entries under its ask list; the command creates the file or a missing list when needed, adds nothing that is already there, sets no default mode and keeps every other byte of the file. A file that is not strict JSON is reported as not readable and left alone. For Codex the entries go into one rules file of Roundfix's own in the project's rules folder, which the command writes whole, and no other file of that folder is read or changed; it sets no approval policy and no sandbox mode. The command never writes outside the repository and refuses a target that is, or passes through, a link that leaves it. Each such file is a **Harness Permission File**.
9. **OpenCode is reported, not written.** For OpenCode the command says that its base policy already allows every shell command, or that the project's own configuration holds rules that ask about or deny shell commands, or that the configuration could not be read. It prints the list in plain form for comparison and writes nothing.
10. **Remove.** With its removal option the command shows, and with that preview's digest removes, exactly the entries and Guard Entries Roundfix wrote. In the Claude Code file an entry is Roundfix's when it equals an entry the list holds or once held; the preview names each one, and the rest of the file stays byte for byte. The Codex rules file of Roundfix's own is deleted. No key, no other file and no folder is removed.
11. **The yes is the user's.** The command is a Named Exclusion, so a harness that asks will ask the person before it runs. The Roundfix skill tells a session to run the preview, show it, put one question to the user and apply only on that answer, with the digest, in every permission mode; a recommendation, silence or a missing answer is never that yes. No other part of Roundfix applies the list: not `roundfix baseline update`, not `roundfix skills install`, not `roundfix setup`, not the Doctor Command, not the Update Check and not an Agent of a Run.
12. **The Doctor Command says where a project stands.** One more check reports, for each harness whose files the repository holds, one state: current; behind, with the count of entries to add and to remove; not applied; not readable; and, for OpenCode, not needed when nothing asks and own rules when the project's configuration asks or denies. Each state that calls for an action names the command. Not applied is information and no warning; behind and not readable are warnings. The check never fails the Doctor Command, reads no file outside the repository and writes nothing.
13. **After an upgrade the difference is named.** The Update Check carries the same states, as text and as JSON, for each harness where the list was applied, and the `update` Skill Mode names the command in its report. The mode never applies the list. A preview in a project that is behind shows only what would be added and removed, so the new yes is to the change.
14. **The list stays true.** A contract of the repository fails, naming the command, when an owned skill names a top-level `roundfix` subcommand, or a Git or GitHub command, that is neither on the list nor a Named Exclusion. A second contract fails when any rendering of the list holds an entry Core Feature 4 forbids.
15. **The permission files are tooling configuration.** The project permission files of the three harnesses become Governed Paths: for Claude Code the shared and the personal project settings files, for Codex every file of the project's rules folder and the project configuration file, for OpenCode the project configuration file. A Task that changes one needs a grant that names it, from the release on, in every adopter. The tooling-authority clause keeps its text: it already covers repository-tooling configuration.
16. **Skill, guide, glossary and decisions.** The Roundfix skill names the command and its `setup` reference says how a session shows the list, when it may apply it and what to tell the user when a listed command is still asked about; its `modes` reference names the report line of Core Feature 13. The user guide gains a page for the command and its part of the Doctor, update and configuration pages. `CONTEXT.md` adds Flow Command List, Named Exclusion, Guard Entry and Harness Permission File, and revises Governed Path. ADR-0444 and ADR-0445 record the decisions.

## User Experience

The blocks below show the content of each surface. The TechSpec fixes the flags and the exact wording.

The preview for Claude Code, in a repository that never applied the list:

```text
Flow Command List for Claude Code, from Roundfix <version>

Would add to .claude/settings.json, under permissions.allow (45 entries):
  Bash(roundfix follow-up *)
  Bash(roundfix phase *)
  Bash(roundfix spec *)
  ...
  Bash(git show *)
  Bash(gh api *)

Would add under permissions.ask, to keep these options asked about (19 entries):
  Bash(git show *--output*)
  Bash(gh api *--method*)
  ...
  These ask in every permission mode, a session that skips prompts included.

Left out, and still asked about:
  roundfix permissions                 the command that applies this list
  roundfix baseline, skills, profiles, release
                                       the skill asks you before it runs them
  git config, git merge, git push, ... they change the repository or its configuration
  gh auth                              it changes your credentials
  A commit, a push, a merge and a Pull Request stay asked about, and a
  delivery's push and merge stay behind the Spec's authorization.

What this does not do:
  It is not a security boundary. A roundfix command runs what this repository
  declares: its Verification, its bootstrap and the commands of its Task files.
  A listed read accepts any option; the entries above keep the known ones asked about.
  The rules apply after you trust this folder in Claude Code.

Nothing was written. To write exactly this:
  roundfix permissions --harness claude --apply 3f2a9c1e84d0
```

The counts are illustrations. The preview for Codex differs in three places: the file is a rules file of Roundfix's own, there is no block of asked options, and the last block says more:

```text
What this does not do:
  It is not a security boundary.
  In Codex an allowed command runs outside the sandbox with nobody asked.
  Codex rules match the start of a command, so no option of a listed read can
  be kept asked about: git show --output=<file> would write that file.
  The rules apply after you trust this project in Codex.
```

Applying with a digest that no longer matches:

```text
The preview changed since digest 3f2a9c1e84d0 was printed. Nothing was written.
Run the preview again: roundfix permissions --harness claude
```

The Doctor Command's line, in a project that applied the list for Claude Code before an upgrade and holds a Codex folder it never applied it to:

```text
flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: claude: behind, 2 to add, 1 to remove; codex: not applied; next: roundfix permissions --harness claude)
```

The removal preview names each entry it would remove and says that an entry with the same text that the user wrote is removed with them.

## Non-Goals / Out of Scope

- The Agents of a Run: the Daemon already answers their permission requests, and the sandbox of a Run keeps its own switch and decisions.
- The user's own harness files: nothing under the user's home directory is read for its rules or written, and Roundfix reads no trust state.
- A permission mode: the command never sets a default mode, an approval policy or a sandbox mode, and never recommends switching prompts off.
- A commit, a push, a Pull Request and a merge by the session: they stay asked about.
- Shells, interpreters, package managers and build tools: the repository's Verification is reached only through the `roundfix` command that runs it.
- The ai-memory tools and every other tool that is not a shell command.
- A hook that approves commands, and any write by the Baseline or by installation.
- A security boundary: the list removes prompts, and isolation stays the sandbox's work and the user's choice.
- Writing for OpenCode, and harnesses other than the three.
- A guard in Codex, and a narrower list for one harness: the maintainer chose every named read, and the Codex rendering holds the same entries.
- A Project Config key: the harness's file is the record that a project applied the list.
- Applying the list to this repository: no Task writes a harness permission file here, and whether Roundfix applies its own list is the maintainer's act after the release.
- The personal project settings file of Claude Code as a target, and a choice of target file.
- The three symptoms of the related Backlog Entry, which wait on their reproduction steps.
- A `roundfix` command that publishes a plan so that the session needs no Git write of its own, a sandbox profile for each harness, and entries for the user's own files: later work, if the first version proves its worth.

## Success Metrics

1. Success Metric: the contract of Core Feature 14 passes on the delivered tree, and it fails, naming the command, when one entry is taken off the list and when a skill names a subcommand that is on neither side.
2. Success Metric: no rendering of the list for any harness holds an entry Core Feature 4 forbids, and the test that proves it fails when such an entry is added.
3. Success Metric: over Claude Code settings files that already hold a user's rules, other keys, an unusual layout and no final newline, an apply changes no byte outside the entries it adds, a second apply changes nothing, and a removal after an apply returns the file to its first bytes when both lists existed and held no entry of the list before; in every other case the removal changes no byte outside the two lists. A file that is not strict JSON is left unchanged and reported.
4. Success Metric: a preview writes nothing; an apply without a digest is refused; an apply with a digest that no longer matches writes nothing and exits with the confirmation-required code; an apply with the current digest writes exactly what the preview listed.
5. Success Metric: the Claude Code rendering holds one Guard Entry for each known option of Core Feature 5 and the Codex rendering holds none and says why; the Codex rules file allows every entry of the list and holds no rule for a Named Exclusion.
6. Success Metric: the Doctor Command prints each state of Core Feature 12 with its next command and never fails on any of them, and the Update Check prints the same states as text and JSON; a project behind the installed list is always told so.
7. Success Metric: neither `roundfix baseline update`, `roundfix skills install`, the Doctor Command nor the Update Check creates or changes a Harness Permission File in a temporary repository, and no target outside the repository is ever written.
8. Success Metric: the paths of Core Feature 15 are Governed Paths, every path governed at the delivery base stays governed, and an ordinary source file stays ordinary.
9. Success Metric: both copies of the Roundfix skill name the command and carry the texts of Core Feature 16, `CONTEXT.md` defines the four new terms and the revised one, the guide pages describe their parts, and the documentation, skill and glossary contracts pass.
10. Success Metric: every test of the command, Doctor, update, governed-path and skills packages that exists at the delivery base passes, the tests that pin the Doctor Command's list and the usage text having moved to the new output.
11. Success Metric: in a supervised session of a trusted folder where the list was applied, a walk through the status view, a phase start and a Spec Consistency Check shows no prompt for a listed command, in Claude Code and in Codex, and a Named Exclusion is still asked about.
12. Success Metric: in the 30 days after a repository applies the list, no Backlog Entry, Finding or message reports a flow stopped by a prompt for a listed command.

## Acceptance evidence

The outside-evidence rows rest on published sources this Spec did not write. The Research Subagent read them for the Idea on 2026-11-19 and again for this PRD the same day; the author read again each passage quoted here except where the measurement record says otherwise.

- Claude Code documentation, permissions (https://code.claude.com/docs/en/permissions): "Rules are evaluated in order: deny, then ask, then allow." and "Bash permission patterns that try to constrain command arguments are fragile." The first supports Core Feature 5, the second the rule of Core Features 2 and 3 that an entry never constrains arguments.
- Claude Code documentation, sub-agents (https://code.claude.com/docs/en/sub-agents): "Each inherits the parent conversation's permission rules" and "Background subagents surface every permission prompt in your main session." They support the goal for Author Subagents: an author meets the session's prompts, and the session's rules remove them.
- Codex documentation, rules (https://learn.chatgpt.com/docs/agent-configuration/rules): "Run the command outside the sandbox without prompting." and "Rules are experimental and may change." They support what the Codex preview says and the risk that the format moves.
- Anthropic, "How we built Claude Code auto mode", 2026-03-25 (https://www.anthropic.com/research/claude-code-auto-mode): "in practice users accept 93% of them anyway". The figure is the vendor's own. It supports the problem and limits what a prompt is worth.
- Pillar Security, "GitPwned: Allowlist to RCE", 2026-07-20 (https://www.pillar.security/blog/gitpwned-allowlist-to-rce): "Command-name allowlists don't work for commands with complex option semantics." It challenges Core Feature 3 and is why Core Feature 5 exists.
- The source of `apm-cli` 0.1.40 (https://docs.rs/apm-cli/0.1.40/src/apm/cmd/init.rs.html): "Add apm commands to Claude allow list?" It shows the form is known, and its entry for an interpreter shows what Core Feature 4 must refuse.

The measurement this plan commits, [`references/2026-11-19-harness-rule-measurements.md`](references/2026-11-19-harness-rule-measurements.md), is the recorded evidence for how each harness decides an entry, a Guard Entry and a subagent.

## Unreachable Acceptance

- criterion: in an interactive Claude Code session of a trusted folder where the list was applied, a walk through the status view, a phase start and a Spec Consistency Check shows no permission prompt for a listed command, in the session and in an Author Subagent it starts, and a Named Exclusion is still asked about.
  reason: it needs a person at an interactive session of the harness, which no Verification or QA row may start; a headless session cannot show a prompt.
  satisfied-by: the operator's walk in Claude Code after the release, with the count of prompts before and after applying, kept as a reference document, recorded session evidence.
- criterion: in an interactive Codex session of a trusted project where the list was applied, the same walk shows no prompt for a listed command, and the session loaded the rules file of Roundfix's own from the project's rules folder.
  reason: it needs a person at an interactive Codex session, and no test may run the real harness.
  satisfied-by: the operator's walk in Codex after the release, kept in the same reference document, recorded session evidence.
- criterion: Success Metric 12 holds in the 30 days after a repository applies the list.
  reason: it needs weeks of use after the release, in a repository whose maintainer chose to apply the list.
  satisfied-by: the operator's count, at a cycle close, of the Backlog Entries, Findings and messages that name that cause, recorded session evidence.

## Validation

The decisions were checked on 2026-11-19 against these sources.

- The Idea left three things unmeasured. What this machine allows was measured, and the record is [`references/2026-11-19-harness-rule-measurements.md`](references/2026-11-19-harness-rule-measurements.md). In headless Claude Code an allow entry of the planned form lets the listed subcommand run and leaves its siblings asked about; an allow entry for a Git read also lets an option that writes a file through; an ask entry for that option wins, in a subagent and in a session that skips prompts too; and a shared file's allow entries are ignored, with a message, in a folder nobody trusted. In Codex's offline rule checker a rule matches the first words of a command and nothing after them, and under Codex's sandbox a Git read runs, a write into `.git` is refused and a write outside the repository is refused. Whether a background subagent can show a prompt is answered by the documentation and was not run. A prompt on a screen, a real Codex session and the number of prompts one flow shows stay unmeasured and are declared under `## Unreachable Acceptance`.
- Characterization at `abe2ebe1`. Roundfix names no settings file of a harness in its code, the repository tracks five Author Subagent files and one link under `.claude/` and nothing under `.codex/` or for OpenCode, and no pattern of the governed set covers a harness permission file. None of `permissions`, `allow`, `harness`, `session` or `trust` is a registered command. The Doctor Command prints text only and fails only on a failed check. The Update Check report already carries optional parts that are left out when empty.
- The classification of Core Features 2 and 3 was made by reading every place the Roundfix skill tells a session to ask the user first, and every Git and GitHub command the twelve owned skills name inside a code span. The Roundfix skill holds such a sentence for the Baseline plan and its starters, for removing or refreshing a user-level skill copy, for applying a profile recommendation and for a release, and for no other subcommand. The `update` Skill Mode runs `roundfix upgrade` without a question, by an earlier decision.
- Secondbrain, read only, through the Idea's sweep of the same day: an earlier measurement of this repository that an allow list "só pré-aprova", the note that an approval prompt is not a safety control, and a practitioner who keeps prompts on and allows listed commands in his own file. The Secondbrain holds no earlier research on a tool that writes its own allow list.
- The web sweep of the Research Subagent on 2026-11-19 found three more tools that write allow entries into a Claude Code settings file after a question. None was read in depth and none is quoted as evidence. It also reported that OpenCode's permission page now names a file with comments and a rule array, which is why Core Feature 9 reads two shapes and writes none.

## Risks

- The list is not a boundary. A `roundfix` command runs what the repository declares, and the Spec Consistency Check can run the commands of Task files, so a session that may edit files without asking can reach any command through a listed one.
- The Git and GitHub reads are on the list by the maintainer's choice, against the author's recommendation. In Claude Code the Guard Entries keep the known options asked about, and they are patterns of the kind the harness's documentation calls fragile: one that misses leaves the read allowed with that option. In Codex nothing guards them, and an allowed read leaves the sandbox. The measurement also shows that in both harnesses most of these reads already run with no rule, so in Codex the entries for Git reads add little beyond that exposure. The decision record keeps the narrower Codex rendering as an option the maintainer can still take.
- A Guard Entry asks in every permission mode. In a project that applied the list, a session that skips prompts is asked, or refused when it cannot ask, before `gh api` with an option that sends a change. A flow that relied on such a call from the session stops there.
- The applied rule set is wide in what a subcommand can do: `upgrade`, `setup`, `gc`, `init`, `migrate`, `stop` and `cancel` are on the list, because the skill holds no sentence that tells a session to ask first about them. The preview lists every entry by name for that reason.
- An Agent in a session with every prompt off can run the preview and the apply itself. The skill forbids it and the tooling-authority clause forbids it; nothing mechanical stops it there. In a Run the settlement check refuses the changed file without a grant.
- A rule that is written and not in force. Until the person trusts the folder the harness ignores the project's allow rules, Roundfix cannot see that state, and the person sees prompts and believes the list failed. The preview, the Doctor line and the skill say it.
- Harness formats move. Codex calls its rules experimental, and OpenCode documents two shapes of its setting. A writer that lags writes rules a harness ignores or reads differently; the walks of `## Unreachable Acceptance` are the check after each release that changes a rendering.
- Making the files Governed Paths changes what a Task of any adopter may touch. A repository that edits its Claude Code settings in a Spec needs a grant for it.
- Text a session reads can steer a listed command. The skills already tell a session to treat such text as data; the list gives an injected instruction fewer prompts to pass.
- The flow still stops elsewhere: at a commit or a push the session makes itself, at the excluded subcommands, and at the ai-memory tools.
- Dissent kept from the Idea: a prompt is not a safety control and a list of allow rules is not one either, so the honest product may be to recommend a sandboxed session with prompts off. And with a commit and a push left out, a flow that plans a Spec still stops at its first Git write.

## Idea-stage judgment

- item: the request in the session, recorded in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md)
- judged-by: nobody (the user asked for the Idea stage)
- verdict: not judged
- stage: ran
- shape: changed — the standing authorization became a closed list that is shown first and written into the project's own harness file only on the user's yes, for the session and its authors and not for a Run, with stated exclusions, a line in the Doctor Command and a removal.

The judge was not run. The maintainer chose the Idea stage for the whole request on 2026-11-19, "Spec 0412 com Ideia agora", and then the split that made this Spec, "Duas Specs (Recommended)". The `shape` line is the author's and waits for the approval of this draft.

He approved the Idea, with its `shape` line, the same day: "Aprovar e seguir ao PRD (Recommended)".

## Rule challenges

None. The author read the rules that touch this goal: the tooling-authority clause and the clause that tells an Agent to stop and ask before it changes tool configuration; the decision that refused a merge into a harness's project file as the Baseline's own mechanism; the two decisions that keep Roundfix from editing the user's runtime files; and the rule that the governed set only grows. None blocks the goal. Each forbids a write nobody authorized or a write to the user's own files, and this Spec writes the project's file, after a person saw the entries and confirmed them.

## Glossary

- adds: **Flow Command List**
- adds: **Named Exclusion**
- adds: **Guard Entry**
- adds: **Harness Permission File**
- changes: **Governed Path**

## Decisions

Every dated line is the maintainer's answer, as recorded in [`references/2026-11-19-maintainer-answers.md`](references/2026-11-19-maintainer-answers.md). The text in quotes is the option he chose, verbatim.

- 2026-11-19: What does the operator do with an open question of an author while the maintainer is away, in this cycle? — "Seguir o padrão do autor (Recommended)": the operator adopts the option the author recommends and records it in the PRD as an operator decision by delegation, with the question and the options; his limits, which are data sent outside, new scope and spend, still stop the Spec.
- 2026-11-19: One Spec or two? — "Duas Specs (Recommended)": the standing authorization of shell commands becomes Spec 0413, with an Idea of its own, delivered after Spec 0412.
- 2026-11-19: What is the target of the first version? — "Sessão Supervisor e autores (Recommended)": the permission prompt in the session that drives the flow, in Claude Code, Codex and OpenCode, and in the Author Subagents.
- 2026-11-19: Which form does it take? — "Mostra a lista e grava com seu sim (Recommended)": a command shows the closed list, of `roundfix` commands and Git and GitHub reads, with no shell wildcard and no interpreter, and writes the project's permission file only on a confirmation; a push and a merge stay behind the authorization record.
- 2026-11-19: Is the Idea approved? — "Aprovar e seguir ao PRD (Recommended)": with the harness permission files as Governed Paths and no new key in the Project Config.
- 2026-11-19: Which `roundfix` subcommands are on the list? — "Os que as skills nomeiam (Recommended)": one per subcommand, except those the skill tells a session to ask about first and the command that writes the list.
- 2026-11-19: Which Git and GitHub reads are on the list? — "Todas as leituras nomeadas": every Git and GitHub read the skills name. This is not the option the author recommended, which was only the reads that are safe. The PRD records the measured risk, which is `git show --output` reaching code execution in Codex and `gh api`, and reduces it inside this choice without narrowing it.
- 2026-11-19: Where does Claude Code's list go? — ".claude/settings.json (Recommended)": the shared file, which is committed; OpenCode is only reported.
- 2026-11-19: Which authorization do Specs 0412 and 0413 have? — "Mesmas do BE para as duas (Recommended)": implement, commit, push, Pull Request, merge and a minor release after each one, up to five corrections, and the Makefile, CI, Go module and lint only if a correction requires one.

Autonomous route: no person present; decisions taken from the approved Idea and its defaults, the maintainer's answers above, the operator's brief of 2026-11-19 and the measurement and characterization recorded in `## Validation`.

- The open questions the approval left to the Idea's defaults are settled as those defaults say: Roundfix's entries in a file that takes no comments are recognized by their text, against the entries the list holds or once held, and a removal names each one first; the absence of a prompt is proved by the walks declared under `## Unreachable Acceptance`; the governed patterns reach adopters with the release; the preview and the Doctor line say that a rule waits on folder trust, and Roundfix reads no trust state.
- The Idea's question about a background subagent is answered by the harness's documentation, quoted in `## Acceptance evidence`, and by the measurement that a subagent takes the session's rules. No requirement depends on an Author Subagent being unable to show a prompt.
- The Idea's exclusion of `gh api` is overruled by the maintainer's answer about reads. The skills name it as the release listing of the Update Check, which is a read.
- The Release Ceiling is `minor`: one command, one check and one report field are added, and nothing is removed or renamed.
- ADR-0444 and ADR-0445 are written with this PRD. The TechSpec adds ADR-0446 for the digest and the way an entry is recognized.

Each entry below is an operator decision delegated by the maintainer (2026-11-19), with its question, its options and the one taken. The author took them under the operator's brief, which is how this cycle applies his answer.

- delegated: 2026-11-19, prd
  - question: does the Codex rendering hold the Git and GitHub reads, given that no guard can be written there and that an allowed command leaves the sandbox?
  - options: yes, every entry of the list, with the measured risk stated where the user confirms; no rule for the Git reads in Codex, which already run inside its sandbox with nobody asked; the reads behind a second confirmation in Codex
  - adopted: yes, every entry — the maintainer chose every named read with this risk named in his answer, and the other two options narrow that choice for one harness, which the brief forbids; the measurement that favors the second option is recorded for him to weigh
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: do Guard Entries exist, and do they cover the options of `gh api` that send a change?
  - options: one for every known option that writes a file, runs a program or changes GitHub; only for the options that write a file or run a program on the machine; none
  - adopted: one for every known option — it is the only reduction a harness offers inside the maintainer's choice, it restores the question a supervised session has today, and its cost, a question in a session that skips prompts, is stated in the preview
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: which `roundfix` subcommands does "those the skill tells a session to ask about first" exclude?
  - options: the ones for which the Roundfix skill holds a sentence that tells a session to ask the user first or to wait for the user's request, which are `baseline`, `skills`, `profiles` and `release`; those and every subcommand that changes the machine beyond the repository, such as `upgrade`, `setup`, `gc`, `init` and `migrate`
  - adopted: the ones the skill's text names — it is the test the maintainer's answer states, a second test of the author's own would narrow his choice, and the `update` Skill Mode already runs `upgrade` without a question
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: which mentions make a Git or GitHub command one "the skills name"?
  - options: every command inside a code span of an owned skill, whether the skill tells the session to run it or describes what Roundfix runs; only the commands a skill tells a session or an author to run
  - adopted: every command inside a code span — it is how the Idea's measurement counted and what his answer chose from, and a contract can check a code span while it cannot check who is told to run a command
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: what makes a named Git or GitHub command a read?
  - options: in the form the skill names it, it changes no ref, index, working-tree file or configuration and nothing on GitHub; it writes nothing at all, unreferenced objects and its own output file included
  - adopted: the first — it keeps `merge-tree`, `archive` and `gh api` on the list as his answer intends, and the options that make them write are what the Guard Entries are for
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: how is the yes given to the command?
  - options: the write takes the digest a preview printed and no other flag writes; an apply flag alone; an apply flag with a yes flag
  - adopted: the digest — a write can then never happen without a preview of the same list and the same file having been printed first, which is the visible confirmation the maintainer asked for in a form a test can prove
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: which harness does the preview show when none is named?
  - options: every harness whose files the repository holds; the harness the command runs under, read from its environment; none, the harness is always named
  - adopted: every harness whose files the repository holds — it is the same rule the Doctor Command's check uses, it gives the same output wherever the command runs, and no harness documents a variable Roundfix could rely on
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: where do the Codex entries go?
  - options: one rules file of Roundfix's own in the project's rules folder, written and deleted whole; merged into a rules file the project already has
  - adopted: a file of its own — Roundfix then never reads or rewrites a project's own rules, which are a small program, and the removal and the comparison need no parsing
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: which files become Governed Paths?
  - options: the project permission files of the three harnesses, the personal Claude Code file, every file of the Codex rules folder, the Codex project configuration and the OpenCode configuration included; only the two files the command writes
  - adopted: the project permission files of the three harnesses — each of them can grant a command or set a mode, and governing only the two written files would leave the others as the way around
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: which states of the Doctor Command's check are warnings?
  - options: behind and not readable warn, and not applied is information; every state but current warns
  - adopted: behind and not readable — a project that never asked for the list must not be warned at every Doctor run, and a project that applied it must be told when it falls behind
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: which skills does the list follow?
  - options: every skill the binary owns; the six a Supervisor session and an author read, which the Idea measured
  - adopted: every owned skill — the contract then needs no list of its own to keep, and the six are among them
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: which names does the glossary take?
  - options: four terms, for the list, what it leaves out, the entry that keeps an option asked about and the file a harness reads; one term, Flow Command List, with the rest described in prose
  - adopted: four terms — the preview, the skill, the guide and the governed set each name these in one or two words, and prose synonyms would drift between them
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, prd
  - question: does this Spec apply the list to this repository?
  - options: no, the maintainer applies it himself after the release if he wants it; yes, a Task writes the Claude Code file here
  - adopted: no — the yes is the user's own act, the sessions of this repository mostly run with prompts off, and a Guard Entry would then ask where nothing asks today
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, techspec
  - question: when has a project applied the list for Claude Code, in a file that cannot mark an entry?
  - options: when its file holds a `roundfix` rule of the list; when it holds any rule of the list, a Git read included
  - adopted: when it holds a `roundfix` rule — people write rules for Git reads themselves, and a project that never ran the command would otherwise be told at every Doctor run that it is behind
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, techspec
  - question: what does the command do with a Claude Code settings file that is not strict JSON or whose lists have another type?
  - options: report it as not readable and write nothing; read it tolerantly and write strict JSON back
  - adopted: report it and write nothing — the harness itself rejects such a file, and a tolerant rewrite would change bytes the command promises to keep
  - under: the user's answer of 2026-11-19
  - answer: none yet
- delegated: 2026-11-19, tasks
  - question: how does the QA gate measure the two things the PRD and the TechSpec left unmeasured, the short-option form of a Guard Entry and whether Codex loads the rules file under its name?
  - options: offline and without a harness, through the rendered rules under the documented reading of `*`, the written file's place and form against the Codex documentation and the measurement this plan commits, with the two walks of `## Unreachable Acceptance` as the proof in a real session; the same, and also Codex's offline rule checker over the written file when `codex` is installed; a session of each harness at the gate
  - adopted: offline and without a harness — the authorization keeps the real `claude`, `codex` and `opencode` out of every test and Verification, a gate row that depends on an installed harness would be blocked or skipped on a machine without it and fail the gate, and the checker's reading of the rule form is already recorded in the measurement
  - under: the user's answer of 2026-11-19
  - answer: none yet

## Open Questions

None that stops the plan. One delegated decision deserves the maintainer's eye before delivery, because a measurement made after his answer bears on it: the Codex rendering holds the Git reads as he chose, and the measurement shows that in Codex those reads already run inside the sandbox with nobody asked, so a rule for them adds only the exit from the sandbox. Until he says otherwise the rendering holds every entry. The walks and the count of `## Unreachable Acceptance` are the operator's, after the release.
