# Daemon owns batch commits and final push

Roundfix creates one commit per successful Batch in the Daemon after the Agent finishes editing, updating assigned Review Issues, and passing verification. A successful Batch may be committed even when other Unresolved Review Issues remain. The Daemon never pushes per Batch or Round. It performs one Final Push only after no Unresolved Review Issues remain for the pull request.

Final Push is enabled by default and sends the complete local PR Head Branch state, including commits that were already unpushed when the Run started. Users disable Final Push explicitly through configuration.

This keeps git policy, Review Source resolution, and push safety in one place instead of letting each Agent runtime decide when or how to commit or push.
