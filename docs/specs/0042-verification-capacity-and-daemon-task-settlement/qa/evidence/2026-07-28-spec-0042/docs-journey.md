# Documentation journey

The repository maintainer journey followed the built CLI and these pages:

- `README.md`;
- `docs/user-guide/configuration.md`;
- `docs/user-guide/commands.md`;
- `docs/user-guide/usage.md`;
- `docs/agents/autonomous-work.md`;
- `CONTEXT.md`.

The generated Config, recommended 2:1 flow, public Implement command, event
filter, Attach replay, exit-75 ownership, deterministic repair, per-Task Agent
Session, and Daemon settlement wording matched observed behavior and current
integration evidence.

`Test(CommandUsage|DocumentationContract)` passed.

The journey failed one advertised CLI-help form:
`roundfix attach <run-id> --no-input` exits 2 although root help and Attach
help place `--no-input` after the optional Run ID. See `cli-contract.md`.
