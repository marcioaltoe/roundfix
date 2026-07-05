# One Agent Session per Run

A Run drives all of its Work Items through one persistent Agent Session, named by the Run and closed when the Run reaches a terminal outcome; sessions are never reused across Runs. Consecutive Batches and Tasks get a warm runtime and in-Run context continuity, crash recovery rides acpx's resume-and-respawn, and cross-Run reproducibility is preserved because every Run still starts from a clean session. The acpx default of reusing sessions across invocations in the same working directory was rejected: Runs would stop being independent.
