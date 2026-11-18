# The maintainer's decision on the root instruction files (2026-11-18)

The maintainer answered in the Supervisor session of 2026-11-18, about 06:30 BRT, while the plan of Spec 0405 was being amended. The operator recorded the answers in Portuguese, and they are copied here verbatim.

## The standard

> "Temos que ter sempre um padrão igual para todos os projetos que usam o baseline do roundfix. Então se for necessário instruções diferentes no AGENTS.md e CLAUDE.md, vamos ter sempre dois arquivos e sem symlink já sem nem dar opção ao usuário, é regra de uso do baseline. Se não for necessário e formos manter arquivo AGENTS.md e o CLAUDE.md com symlink para o AGENTS.md, vai ser sempre assim. Sem opção ao usuário de mudar o baseline no projeto. O único caminho para personalização do baseline é o arquivo que contem regras especificas do repositorio, hoje chamado de specific-repository.md"

> "eu gosto da ideia de usar o CLAUDE.md como symlink para o AGENTS.md, porém… o trecho que fala sobre o uso da skill afeta de alguma forma … o uso das skills no Codex ou Open Code? Se não afeta, a gente vai ter dessa forma: um AGENTS.md, que é o arquivo, e o CLAUDE.md sempre o symlink… Mas se houver necessidade de separar… essa instrução, que é específica para o Claude, ficaria no CLAUDE.md."

The Supervisor answered that the Claude-only statement Spec 0405 adds to the root block opens with "In Claude Code," and says it does not apply to Codex or OpenCode. No instruction differs per harness, so the standard is one `AGENTS.md` and a `CLAUDE.md` that is always a symbolic link to it.

## Route and authorization

- Where the rule goes: "Spec 0406 depois da 0405 (Recommended)".
- Authorization: "Mesmas do BE (Recommended)": implement, commit, push, pull_request, merge and release, in the single release at the end of the queue; up to five corrections; the Makefile, CI, `go.mod` and lint only if a correction requires one.

## The conversion and product questions (about 07:00 BRT)

Recorded in the operator's answers file, verbatim: "0406 conversão → \"Só depois de mover (Recommended)\"" and "perguntas de produto da 0406 → \"deve parar e perguntar.\"". The coordinator relayed the chosen option: when the regular `CLAUDE.md` holds passages of its own, the update plan lists each one to move to `docs/agents/specific-repository.md` and does not convert; with no passage of its own it converts it to the symbolic link, with the Baseline's backup; no content is ever lost. Any other product question stops the authoring and goes back to the maintainer; technical choices stay with the author and are recorded as ADRs.

## The operator's brief

The operator's brief for this Spec, of 2026-11-18, adds the scope the answers imply:

- Adoption always creates `CLAUDE.md` as a relative symbolic link to `AGENTS.md`.
- `roundfix baseline update` and `roundfix doctor` treat a regular `CLAUDE.md`, or a link elsewhere, as a deviation; the update plan converts it into the link, and no content is lost. A regular `CLAUDE.md` whose own passages are Misplaced Repository Rules is reported with each passage to move first. The author decides, taking the recommended option, whether the plan refuses the conversion until the file holds no passage of its own or converts it and keeps a backup.
- The guides, skills and documents that describe the root files state the rule, and the Baseline's handling of symbolic links on Windows stays consistent.
- The adopters' checkouts on the host are inventoried read only.
