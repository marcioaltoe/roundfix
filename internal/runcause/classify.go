package runcause

// Evidence carries independent texts. Empty text means absent evidence.
type Evidence struct{ Diagnostic, Command, Task string }

// Classify returns the first matching signature, or unclassified without a
// signature or source. An any signature tests each text independently.
func (t Table) Classify(e Evidence) (class, signature, source string) {
	texts := []struct{ source, text string }{{"diagnostic", e.Diagnostic}, {"command", e.Command}, {"task", e.Task}}
	for _, s := range t.Signatures {
		for _, text := range texts {
			if (s.Source == "any" || s.Source == text.source) && text.text != "" && s.re.MatchString(text.text) {
				return s.Class, s.ID, text.source
			}
		}
	}
	return "unclassified", "", ""
}

// Trigger names the first trigger in table order, else unknown.
func (t Table) Trigger(taskText string) string {
	for _, trigger := range t.Triggers {
		if trigger.re.MatchString(taskText) {
			return trigger.ID
		}
	}
	return "unknown"
}
