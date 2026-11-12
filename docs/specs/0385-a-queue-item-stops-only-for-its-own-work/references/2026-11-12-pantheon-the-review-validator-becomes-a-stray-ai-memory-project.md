# ai-memory message 01a1d01a-fa8c-7211-bd88-5023f39563e1 — from pantheon (default/pantheon), 2026-11-12T14:13:11Z, nature measured

Subject: Pantheon: .ai-memory.toml na raiz e o validador de revisão que vira projeto solto no ai-memory

(Item 1 of the message — commit `.ai-memory.toml` at the repository root — is done by the operator in the planning PR; this Spec owns only item 2.)

2. FINDING, sobre o comportamento do Roundfix: o validador de revisão roda o Codex num diretório temporário fora da home e fora de qualquer checkout (/var/folders/<x>/T/roundfix-review-validator-<n>). O ai-memory não acha marcador nem remote ali e cria um projeto com o nome da pasta: havia 8 roundfix-review-validator-<n>. Os Runs em worktrees também criaram projetos run_* e roundfix-deliver-* (sem sessões), provavelmente pelo sandbox do Codex sem acesso ao .git do repositório principal. O marcador do item 1, versionado em cada repositório, cobre as worktrees.
- Mitigação já aplicada: o `~/.ai-memory.toml` do MacBook e do mini tem uma rota de /var/folders para default/marcioaltoe-roundfix.
- Correção sugerida: o validador grava um `.ai-memory.toml` no diretório temporário, com o workspace e o projeto do repositório revisado, ou roda dentro de uma worktree desse repositório. Assim as sessões dele caem no projeto certo.
- Limpeza: os projetos soltos já foram removidos.

Evidence: ai-memory sessions (sqlite readonly no mini, 12/11/2026): 8 projetos roundfix-review-validator-<n> com cwd /var/folders/<x>/T/roundfix-review-validator-<n> (agent codex); projetos run_*, roundfix-deliver-* com 0 sessões; pantheon docs/references/2026-11-12-projetos-do-ai-memory.md (PR #200)
