# Autonomous work

Agent Selection Profiles in Project Config choose each Agent Session's ACP
Runtime, model, and reasoning effort (`roundfix profiles show`). The
repository prefers {{runtime.backend}} for backend work and {{runtime.design}}
for design, UI, UX, and frontend-dominant work. A `complexity: low` Task that
is not `qa` and changes no Governed Path runs on the Light Tier when it is
available.

{{artifact.rules}}
