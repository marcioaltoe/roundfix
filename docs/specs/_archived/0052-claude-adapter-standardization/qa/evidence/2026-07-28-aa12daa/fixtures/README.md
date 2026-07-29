# Legacy adapter fixtures

The QA run staged `legacy-agent-acp.sh` at a disposable symlink-resolved path
containing
`node_modules/@zed-industries/claude-agent-acp/bin/claude-agent-acp`.
It staged `legacy-code-acp.sh` the same way under
`node_modules/@zed-industries/claude-code-acp/bin/claude-agent-acp`.

Each disposable package path was exposed through a leading `PATH` directory as
`claude-agent-acp`. This makes the built Doctor exercise the production
resolved-path lineage classifier while keeping ignored `node_modules`
artifacts out of the QA report commit.
