# Autonomous work

The Supervisor orchestrates and authors Specs. Implementation goes to the
selected ACP Runtime, normally through Roundfix. The Supervisor does not write
feature code or tests while a runtime can execute the Task.

Default backend work uses the configured Codex runtime. Design, UI, UX, and
frontend-dominant work uses the configured Claude runtime when the Task Graph
routes that surface.

Every autonomous Task must pass its Task Verification and repository gate
before status changes to completed.
