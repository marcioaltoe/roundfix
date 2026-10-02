# Published source observations — 2026-10-02

Read using web.open/web.find; curl-cffi returned CONNECT tunnel 403, so its empty downloads are not evidence of source content.

- https://github.com/cli/cli/releases/tag/v2.81.0 — release heading confirms 2.81.0; the authentication section introduces JSON output for auth status. Supports the minimum version.
- https://raw.githubusercontent.com/cli/cli/v2.81.0/pkg/cmd/auth/status/status.go — authEntryState enumerates success, timeout and error; authStatus contains a hosts map. Export mode returns nil for an empty host list and after emitting statuses; absent --show-token clears Token before JSON export. Authentication failures therefore require parsing state rather than relying on process exit. Supports coded findings and no-token probes.
- https://raw.githubusercontent.com/git/git/v2.23.0/Documentation/RelNotes/2.23.0.txt — the user-interface notes introduce switch and restore. Supports the Git floor used by Roundfix.

These published records predate this Spec and independently support its integration assumptions. No real authentication was performed.
