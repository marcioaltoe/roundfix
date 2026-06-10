# YAML user and project config

Roundfix uses YAML for User Config at `~/.roundfix/config.yml` and Project Config at `<repo>/.roundfixrc.yml`. YAML is a better fit for nested Review Source, ACP Runtime, and watch settings than a flat CLI-only model, and it matches the surrounding agent and skill configuration ecosystem.

Roundfix provides `roundfix init` to create either config file. When `--scope` is omitted, the CLI asks for the scope and defaults to Project Config so first-time setup stays repository-local unless the user chooses User Config explicitly.
