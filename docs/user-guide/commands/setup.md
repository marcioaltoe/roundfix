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

