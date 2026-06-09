# One active run per PR head branch

Roundfix allows simultaneous Active Runs for different PR Head Branches, but rejects a second Active Run for the same Head Repository and PR Head Branch. The global Run Database coordinates this boundary so future daemon modes can handle multiple repositories without letting two loops mutate the same branch and review state at once, including pull requests opened from forks.
