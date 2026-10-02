package agent

// MissingNodePreloads returns the missing absolute preload paths using the same
// NODE_OPTIONS rules as the agent environment, without changing that environment.
func MissingNodePreloads(value string, exists func(string) bool) []string {
	_, missing, _ := agentNodeOptions(value, exists)
	return missing
}
