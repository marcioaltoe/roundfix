# Published source confirmation — 2026-10-01

Both sources were fetched through the web reader; neither is a Spec-authored source.

- [Build Systems à la Carte](https://www.microsoft.com/en-us/research/wp-content/uploads/2018/03/build-systems-a-la-carte.pdf), Mokhov, Mitchell and Peyton Jones, ICFP 2018, section 4.2.2, PDF pages 14–15. The source reader resolved a 29-page PDF. It describes recording dependency values or hashes and checking them on a later build to decide whether the key needs rebuilding. This supports digest-based reuse; it does not prove Roundfix's implementation.
- [Go test command documentation](https://pkg.go.dev/cmd/go/internal/test), Package documentation, source-reader lines 242–267. Successful results can be cached; tests reading module files or environment variables match later runs only when those inputs stay unchanged. The documented way to disable caching is `-count=1`, used in this QA's test commands. This supports success-only reuse bound to input identity; it does not prove Roundfix's implementation.

Independent confirmation: each fetched document's title and named subsection agree with the PRD citation. Local archived reports and Git comparisons supply separate operational evidence.
