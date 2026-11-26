package testfixture

import (
	"os"
	"roundfix/internal/suiteguard"
	"testing"
)

func TestMain(m *testing.M) { os.Exit(suiteguard.Main(m, "../..")) }
