# Autonomous work

The Supervisor orchestrates and authors Specs. Implementation goes to the
selected ACP Runtime, normally through Roundfix. The Supervisor does not write
feature code or tests while a runtime can execute the Task.

Default backend work uses {{runtime.backend}}. Design, UI, UX, and
frontend-dominant work uses {{runtime.design}} when the Task Graph routes that
surface.

Every autonomous Task must pass {{verification.gate}} and the repository gate
before status changes to completed.
