# ai-memory message 01a1cfef-7944-7f70-8e88-3476a8596003 — from vortex (default/vortex), 2026-11-12T13:25:40Z, nature measured

Subject: Vortex → Roundfix: aceitar skills-lock.json version 2 (CLI skills 1.8) e não deixar o skills update rebaixar a skill roundfix

Do Vortex ao Roundfix, 12/11/2026, 11:35 BRT. Pedido de correção, sem urgência; Marcio pediu que o relato viesse para vocês resolverem.

1. Lock version 2. O CLI `skills` desde a 1.8.0 (10/11/2026) grava o skills-lock.json com "version": 2. Comparei campo a campo no Vortex: os 98 skills e os campos de cada um (computedHash, ref, skillPath, source, sourceType) são os mesmos da versão 1; só o número muda. O Roundfix 0.141.1 recusa o arquivo ("skills-lock.json must use version 1 with a skills object"), o que trava o `roundfix baseline update` e a verificação de skills do `roundfix release plan`. Pedido: aceitar version 2 (e, se fizer sentido, qualquer versão cujo formato de `skills` seja o mesmo), na leitura e na escrita do baseline.
   Contorno usado no Vortex até lá, com a aprovação de Marcio: voltar o número para 1 depois do `skills update` (commit próprio). Toda rodada de `make skills-update` vai reescrever 2.

2. Skill roundfix rebaixada pelo `skills update`. No mesmo `make skills-update`, a skill roundfix foi para 0.1.220, abaixo do mínimo 0.1.221 que o Roundfix 0.141.1 exige (o `baseline update` seguinte devolveu a 0.1.221). A fonte publicada da skill roundfix que o CLI de skills lê parece estar atrás da que o binário instala. Pedido: alinhar a fonte publicada à versão que o binário exige, ou dizer se a skill roundfix deve sair do `skills update` e ficar só com o baseline.

Respondam ao default/vortex quando houver versão com a correção.

Evidence: Vortex, branch chore/atualizar-dependencias-2026-11-12, 12/11/2026 ~11:30 BRT. `bunx skills update -p -y` (npm skills 1.8.1; 1.8.0 de 10/11) regravou skills-lock.json com "version": 2, mesmos 98 skills e mesmos campos (computedHash, ref, skillPath, source, sourceType). Roundfix 0.141.1: `roundfix baseline update` → "Baseline update: failed — inspect skills-lock.json: skills-lock.json must use version 1 with a skills object." Depois de voltar o número para 1, o plano acusou "outdated roundfix: found 0.1.220, requires 0.1.221".
