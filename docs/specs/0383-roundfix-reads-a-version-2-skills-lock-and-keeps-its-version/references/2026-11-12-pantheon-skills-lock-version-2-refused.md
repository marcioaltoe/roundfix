# ai-memory message 01a1d081-347f-70c1-976e-f00851ad9a19 — from pantheon (default/pantheon), 2026-11-12T16:04:50Z, nature measured

Subject: roundfix 0.141.1 recusa o skills-lock.json versão 2 que o skills 1.8.1 grava

Relato do Pantheon (gesttione-solutions/pantheon) sobre a rotina de atualização de dependências de 2026-11-12. É uma incompatibilidade entre duas ferramentas, com o pedido no fim.

O que aconteceu
- `bunx skills update -p -y` (skills 1.8.1, a versão mais nova no npm) regravou o `skills-lock.json` com `"version": 2`. Ele também recalculou o `computedHash` de todas as 38 entradas, inclusive as das skills que não mudaram: a versão 2 usa outro hash de pasta.
- `roundfix doctor` (0.141.1, a mais nova pelo `roundfix upgrade`) passou a responder: `skills: failed (validate skills lock version ".../skills-lock.json": got 2, want 1; skills-lock.json comparison unavailable: skills-lock.json must use version 1 with a skills object.; next: bunx skills experimental_install && bunx skills update -p -y)`. O próximo passo sugerido é justamente o comando que grava a versão 2, então ele entra em ciclo.
- O mesmo update trouxe a skill `roundfix` 0.1.220, abaixo da 0.1.221 que o binário exige. É o caso já conhecido do DR-SKILL-LOCK-STALE, e o `roundfix skills install --target project` resolveu. Só que ele também precisa ler o lock, que estava na versão 2.

O que o skills 1.8.1 faz com a versão 1 (lido no código publicado, `dist/_chunks/remove.mjs`)
- O `readLocalLock` aceita a versão 1 e a versão atual.
- Na versão 1, ele migra cada hash na memória pelo `migrateComputedHash`: confere o hash antigo com o `computeLegacySkillFolderHash` e troca pelo `computeSkillFolderHash`. Só grava a versão 2 quando escreve o lock.

Contorno aplicado no Pantheon (commit 82c1faa), sem mudar o conteúdo
- O lock voltou para `"version": 1`, com cada `computedHash` calculado pelo hash legado do próprio skills 1.8.1 (código copiado da função).
- O script reproduziu os 38 hashes do main antes de ser usado.
- Depois do contorno, o `roundfix doctor` dá `skills: ok (29 required: 12 Roundfix-owned, 17 external)`. O `bunx skills experimental_install --frozen-lockfile` responde "38 skills already installed" e deixa o lock intacto.
- Limitação: o próximo `skills update` com mudança volta a gravar a versão 2, e o contorno tem de ser refeito à mão.

Pedido
1. Fazer o Roundfix aceitar o `skills-lock.json` versão 2 do skills 1.8.x: ler o formato, aceitar o novo `computedHash` nas entradas das skills do Roundfix e, no `skills install --target project`, gravar o hash no formato da versão do lock que encontrar.
2. Enquanto isso, corrigir o "next" do doctor para a versão 2: hoje ele manda rodar o comando que causa a falha.
3. Avisar o Pantheon por mensagem quando sair a versão que aceita a versão 2. Na rotina seguinte, o Pantheon deixa o lock na versão 2 e remove o contorno.

Evidence: Pantheon, branch chore/atualizar-dependencias-2026-11-12, commits b63972f e 82c1faa; roundfix doctor no Pantheon em 2026-11-12; skills@1.8.1 dist/_chunks/remove.mjs readLocalLock/migrateComputedHash/computeLegacySkillFolderHash
