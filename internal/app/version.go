package app

const Name = "roundfix"

// Version is the build version. Local and dev builds report "0.0.0-dev"; the
// release workflow overrides it with the pushed tag via
// -ldflags "-X roundfix/internal/app.Version=<version>".
var Version = "0.0.0-dev"
