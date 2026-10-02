# Published-source observations — 2026-10-02 QA rerun

Fresh web.open/web.find reads after curl-cffi failed with CONNECT tunnel 403. These independent pages remain outside the Spec’s authored premise.

- [GitHub CLI 2.81.0 release](https://github.com/cli/cli/releases/tag/v2.81.0), release section at lines 172–181: introduced JSON support for auth status. Supports the version floor; current public fake-forge behavior and named tests confirm integration.
- [Versioned auth status source](https://raw.githubusercontent.com/cli/cli/v2.81.0/pkg/cmd/auth/status/status.go), lines 21–45: explicit success/timeout/error states and initialized hosts map; lines 185–201: unauthenticated JSON uses an empty hosts map with successful return; lines 260–275: JSON omits token unless explicitly requested and returns success after exporting states. Inference: Roundfix must diagnose JSON state rather than process exit alone. Fresh unauthenticated/offline public transcripts independently exercise this distinction.
- [Git 2.23.0 release notes](https://raw.githubusercontent.com/git/git/v2.23.0/Documentation/RelNotes/2.23.0.txt), lines 54–58: introduces switch and restore, supporting the stated Git floor. The PRD’s master .adoc link has moved; this immutable versioned source supplies the same release notes. Fresh version-floor tests confirm Roundfix’s integration.

No live gh credential or repository permission query was made. These observations are paraphrases of the published records.
