### setup

```bash
roundfix setup [--yes] [--no-input]
```

Verifies Node.js, the minimum supported acpx version, the effective adapters,
generated Agent Selection Profiles, acpx local adapter overrides, User Config,
and Project Config. Adapter Readiness requires official
`@agentclientprotocol/codex-acp` lineage at version `2.0.1` or newer and
official `@agentclientprotocol/claude-agent-acp` lineage at version `0.84.0`
or newer. The deterministic install actions are
`npm install -g @agentclientprotocol/codex-acp@2.0.1` and
`npm install -g @agentclientprotocol/claude-agent-acp@0.84.0`.

After `acpx` and before adapter work, Setup prints Doctor's five readiness
lines in order: `gh`, `git`, `remote`, `toolchain`, and `environment`. They
also appear when acpx is unavailable. Each uses `<name>: <status> (<detail>)`;
`failed` and `warn` findings include a stable `DR-` code and `; next: <action>`.
Setup offers no install or change for these lines, including with `--yes`.
Only `failed` makes Setup exit `1` at the end; `warn` alone does not.

- `gh` checks GitHub CLI version, login for the repository's forge, and write
  permission. `DR-GH-UNAUTHENTICATED` points to `gh auth login --hostname <host>`.
- `git` checks Git version and repository `user.name` and `user.email`.
- `remote` checks the delivery remote (`watch.push_remote`, otherwise `origin`)
  and whether it is a reachable GitHub forge remote.
- `toolchain` finds executables named by configured Verification, bootstrap,
  regeneration commands, and `verification.tools`, without running them.
  `DR-TOOL-MISSING` names the missing executable and the command that needs it.
- `environment` names missing `NODE_OPTIONS` preload files and reports optional
  `ROUNDFIX_` key variables as set or not set, never their values.

Forge reads use the user's `gh` and Git credentials, no standard input, disabled
Git terminal prompts, and a ten-second cancellation limit per read. A timeout
or unreachable forge reports `warn`; Roundfix does not log in or change Git or
shell configuration. Run `roundfix doctor` again after taking the next action.

A stale or bare Codex override that fails official lineage proof produces one
migration offer to `npx -y @agentclientprotocol/codex-acp@2.0.1`, and a Claude
override that fails the same proof produces one migration offer to
`npx -y @agentclientprotocol/claude-agent-acp@0.84.0`. The offer follows from
the failed proof, so it covers a differently named or differently scoped
package without naming any superseded one. Setup proves each proposal before
asking; declining preserves the acpx configuration bytes.

Setup builds every proposed file in memory and runs exact Agent Selection proof
before writing. It never changes explicit Sol/high to model-managed reasoning
when proof fails. Each check prints one deterministic report line such as
`node: ok`, `adapter: migration proposed`, `profile readiness: passed`, or
`User Config: skipped`. `--yes` accepts every offered install or file change;
`--no-input` performs diagnosis and skips offers without writing.
When acpx is missing or older than `0.12.0`, Setup offers
`npm install -g acpx@0.12.0`. It accepts `0.12.0` and newer versions without
offering a downgrade.

