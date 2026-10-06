// Package lighttier classifies Tasks and records monthly light model spend.
package lighttier

import (
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

type Tier string

const (
	Light    Tier = "light"
	Standard Tier = "standard"
)

// TierFor excludes QA and every Task declaring a Governed Path to write.
func TierFor(task spec.Task) Tier {
	if task.Complexity != "low" || task.Type == spec.TaskTypeQA {
		return Standard
	}
	for _, ref := range task.Context {
		switch ref.Kind {
		case spec.ContextKindCreates, spec.ContextKindInterface, spec.ContextKindDeletes:
			if speccheck.GovernedPath(ref.Path) {
				return Standard
			}
		}
	}
	return Light
}

// Plan carries Run wiring; a one-Run override leaves the Plan zero or Models empty.
type Plan struct {
	Models      []string
	CeilingUSD  float64
	KeyVariable string
	KeyPresent  bool
	HomeDir     string
	Repository  string
}

// Enabled keeps availability checks separate so dispatch can explain skips.
func (plan Plan) Enabled() bool {
	return len(plan.Models) > 0
}
