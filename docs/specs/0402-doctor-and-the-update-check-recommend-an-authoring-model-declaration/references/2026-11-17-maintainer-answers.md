# Maintainer answers for Spec 0402

The maintainer answered through AskUserQuestion on 2026-11-17, about 10:20 BRT, in the Supervisor session that planned Specs 0399 to 0401. The operator recorded the answers in Portuguese, and they are copied here verbatim.

- 17/11 ~10:20 BRT: falta de skills.authoring.claude não é avisada → "Avisar e sugerir (Recommended)": doctor e update check dão warn quando as skills de autoria estão instaladas sem declaração, mostrando o bloco recomendado para colar no .roundfixrc.yml; não falha nem grava; Spec fix após a 0401. Autorização da Spec 0402 → "Mesmas do BE (Recommended)".

The operator's brief for the TechSpec adds:

- the block shows the models Roundfix itself declares or uses;
- a repository that publishes the authoring skills (ADR-0408) stays exempt, as today;
- the author decides where the block's values come from, with one source of truth and not a hard-coded copy that drifts, and records the choice as an ADR;
- the authorization covers implement, commit, push, pull request, merge and the release right after the Spec, up to three corrections; the Makefile, CI, `go.mod` and lint only if a correction requires one.
