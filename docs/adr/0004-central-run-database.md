# Central run database

Roundfix stores Run state in a global SQLite database at `~/.roundfix/roundfix.db` instead of inside each repository Artifact Directory. Centralizing Run state gives the Daemon one place to track progress across multiple repositories and future concurrent Runs, while repositories may contain only markdown review artifacts when users choose that layout.
