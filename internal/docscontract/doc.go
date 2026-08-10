// Package docscontract is the invalidation domain for repository-markdown
// validation. Every test here reads prose the repository publishes — user
// guides, agent guidance, the Spec corpus — so any docs commit invalidates
// exactly this package and nothing else. make verify excludes it; make
// verify-docs runs it, and a pull request must not open before verify-docs
// passes. Nothing under a docs/**/_archived tree is asserted here.
package docscontract
