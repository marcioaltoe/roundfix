# Documentation and CLI help

- `cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited `0`.
- `./bin/roundfix skills check` exited `0` and reported the Roundfix skill set passed.
- `./bin/roundfix supersede --help` exited `0` and exposed the invocation, all three options, exit codes, one-file write, and no-Run/no-commit/no-push envelope.
- `./bin/roundfix archive --help` exited `0` and described the recorded-supersession proof only for a Spec without a Task Graph, with an unchanged move.
- Focused `rg` reads found the invocation and refusal/archive wording in the canonical skill, its mirror, and `docs/user-guide/commands.md`.

The mirror comparison was repeated after the help and skill commands and stayed
byte-identical.

