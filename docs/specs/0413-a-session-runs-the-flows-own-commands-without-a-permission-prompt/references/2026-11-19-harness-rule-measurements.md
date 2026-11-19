# What the harness rules do with a listed command

Measured by the author of the PRD on 2026-11-19, on the Mac mini, with Claude Code 2.1.328 and codex-cli 0.170.0. It answers what the Idea left unmeasured as far as this machine allows without an interactive session: how each harness decides a listed command, an unlisted one and a listed read that carries an option that writes, and what a subagent takes from the session.

Every run used a throwaway Git repository in a temporary directory outside this repository, removed afterwards. It held one file and a stand-in script named `flowprobe` that only prints its arguments. The real `roundfix` was never run by a probe, so that no probe could start an upgrade. No file or text of this repository was sent anywhere: the Claude Code runs sent the probe's own prompt, and the Codex measurements ran offline.

## 1. Claude Code, headless

Method: `claude -p <prompt> --model haiku --output-format stream-json --verbose --setting-sources project,local --no-session-persistence`, default permission mode, reading each tool call and its result from the stream. The prompt asked for the commands one by one and forbade retries. A headless session cannot show a prompt, so a command that would be asked about is refused, and the refusal is the evidence that a prompt would have appeared.

The rules under test, as `permissions.allow` and `permissions.ask` of a settings file:

```json
{
  "permissions": {
    "allow": ["Bash(./flowprobe spec *)", "Bash(git show *)"],
    "ask": ["Bash(git show *--output*)"]
  }
}
```

| Run | Settings file | `./flowprobe spec check demo` | `./flowprobe upgrade` | `git show HEAD:README.md` | `git show --output=probe-out.txt HEAD` |
| --- | --- | --- | --- | --- | --- |
| A | untracked `.claude/settings.local.json`, both lists | ran | refused: "This command requires approval" | ran | refused |
| B | the same file without the `ask` list | ran | refused | ran | ran, and wrote `probe-out.txt` |
| C | `.claude/settings.json` in a folder nobody trusted, both lists | refused | refused | ran | refused |
| D | as A, with the four commands handed to a subagent defined in `.claude/agents/` | ran | refused | ran | refused |
| E | as A, with `--permission-mode bypassPermissions` | ran | ran | ran | refused |

What the rows show:

- An allow entry of the form `Bash(<command> <subcommand> *)` lets the listed subcommand run and leaves every other subcommand of the same command asked about (A).
- An allow entry for a Git read also allows the same read with an option that writes a file (B). This is the hazard the published reports describe, reproduced here.
- An ask entry for that option wins over the allow entry (A against B).
- A subagent takes the session's allow and ask entries (D).
- An ask entry still stops the command in a session that skips every prompt (E). There the unlisted command ran and the guarded one did not.
- In a folder nobody trusted, the shared file's allow entries are ignored (C). Standard error said: "Ignoring 2 permissions.allow entries from .claude/settings.json: this workspace has not been trusted. Run Claude Code interactively here once and accept the trust dialog". The plain read still ran, because Claude Code runs it by itself.

A sixth run, with no settings file at all, asked for the Git reads the owned skills name, in the forms they name:

| Ran with no rule | Refused with no rule |
| --- | --- |
| `git status --porcelain`, `git diff --name-only HEAD`, `git grep -l hi HEAD`, `git ls-files --cached --others --exclude-standard`, `git rev-list --first-parent --reverse HEAD`, `git cat-file -e HEAD` | `git rev-parse --path-format=absolute --git-path objects`, `git merge-tree --write-tree HEAD HEAD`, `git archive --format=tar HEAD README.md` |

So in Claude Code most of the named Git reads already run without a rule. An entry changes the outcome for the forms Claude Code does not class as a read by itself, and for a read that carries an option.

## 2. Codex, offline

### The rule decision

Method: `codex execpolicy check --rules <file> -- <command>`, which prints the decision a rules file gives a command and starts no session. The file held an allow rule for `["roundfix", ["spec", "phase", "follow-up"]]`, an allow rule for `["git", "show"]`, an allow rule for `["gh", "api"]` and a prompt rule for `["git", "show", "--output"]`.

| Command | Decision |
| --- | --- |
| `roundfix spec check 0413-x` | allow |
| `roundfix spec` | allow |
| `roundfix upgrade` | no rule matches |
| `roundfix` | no rule matches |
| `git show HEAD:README.md` | allow |
| `git show --output=.git/config --format=x HEAD` | allow |
| `git show HEAD --output=.git/config` | allow |
| `git show --output .git/config HEAD` | prompt |
| `git -C . show HEAD` | no rule matches |
| `gh api -X DELETE repos/o/r/git/refs/heads/x` | allow |
| `/usr/local/bin/roundfix spec check x` | no rule matches |
| `env FOO=1 roundfix spec check x` | no rule matches |

A rule matches whole words from the start of the command and nothing after them. A list in one position names alternatives. A rule cannot tell `git show` from `git show --output=<file>`: the prompt rule caught the option only as a separate third word, and the same option written with an equals sign or after the revision was allowed. No guard for an option can therefore be written in a Codex rules file.

### The sandbox

Method: `codex sandbox -- <command>` from inside the probe repository, which runs one command under the sandbox Codex gives a session and starts no session.

| Command | Result inside the sandbox |
| --- | --- |
| `git show HEAD:README.md`, `git rev-parse --short HEAD`, `git ls-files`, `git grep -l hi HEAD` | ran |
| `git status --porcelain` | ran, with one error line about the file-system monitor socket of this machine's Git configuration |
| `git show --output=.git/probe-out HEAD` | refused: "could not open '.git/probe-out' for writing: Operation not permitted" |
| `touch <a directory outside the repository>/x` | refused: "Operation not permitted" |
| `touch <the repository>/inside.txt` | ran |
| `curl https://api.github.com/zen` | refused: the connection was blocked |

What the rows show:

- A Git read runs inside the sandbox, and Codex asks about a command only when it leaves the sandbox. In a session of the default kind, a named Git read therefore runs without a rule.
- Inside the sandbox the option that writes is harmless: the write into `.git` was refused. The documentation says an allow rule means "Run the command outside the sandbox without prompting." Taken with the first table, an allow rule for `git show` lets `git show --output=.git/config` write that file with nobody asked.
- A command that writes outside the repository or reaches the network fails inside the sandbox. A `roundfix` command writes Roundfix Home and reaches GitHub, and `gh` reaches GitHub, so these are the commands a Codex session asks about and the ones a rule removes the question from.

## 3. Passages this record rests on

The author read the first, the second and the fourth again on 2026-11-19.

- Claude Code documentation, sub-agents (https://code.claude.com/docs/en/sub-agents): "Each inherits the parent conversation's permission rules" and "Background subagents surface every permission prompt in your main session."
- Claude Code documentation, permissions (https://code.claude.com/docs/en/permissions): "Rules are evaluated in order: deny, then ask, then allow."; "A `*` at the end, with a space before it, also matches the bare command."; "Every subcommand must be allowed for a compound command to run without a prompt"; "Claude Code shows the trust dialog in interactive terminal sessions, but not in `claude -p` or SDK sessions."
- Claude Code documentation, settings (https://code.claude.com/docs/en/settings), read by the Research Subagent and not read again by the author: "Settings files are strict JSON: a `//` comment or a trailing comma is a syntax error".
- Codex documentation, rules (https://learn.chatgpt.com/docs/agent-configuration/rules): "Run the command outside the sandbox without prompting." and "Rules are experimental and may change."
- GitHub CLI manual (https://cli.github.com/manual/gh_api), read by the Research Subagent and not read again by the author: "The default HTTP request method is `GET` normally and `POST` if any parameters were added."

## 4. What stays unmeasured

- A prompt on the screen. Every result above is a refusal in a session that cannot ask, or a decision printed by a checker. That a person sees no prompt for a listed command in an interactive session of a trusted folder needs that session and that person.
- A background subagent. The documentation says its prompts appear in the main session. No run here started one.
- A real Codex session. The rule decision and the sandbox were measured apart. That a session loads `.codex/rules/roundfix-flow.rules` from a trusted project under that file name was not measured; the documentation names the folder and the `.rules` extension and gives `default.rules` as its example.
- How many prompts one flow shows today. No session of this repository recorded it.
- OpenCode. No session and no configuration of it was measured.
