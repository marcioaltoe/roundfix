package spec

import "sort"

// QAGateClosure reads the declared QA Task and its sorted transitive dependency
// closure from manifest bytes, without reading Task files or their statuses.
func QAGateClosure(manifestPath string, content []byte) (string, []string, error) {
	nodes, _, _, qa, _, err := parseManifestNodes(manifestPath, content)
	if err != nil {
		return "", nil, err
	}
	if _, err := topologicalOrder(nodes); err != nil {
		return "", nil, err
	}
	if !qa.Present || qa.Declined {
		return "", nil, nil
	}
	byID := make(map[string]manifestNode, len(nodes))
	for _, node := range nodes {
		byID[node.ID] = node
	}
	gate, found := byID[qa.TaskID]
	if !found {
		return "", nil, QAGateError{ManifestPath: manifestPath, Reason: "qa: names a Task absent from the graph"}
	}
	seen := make(map[string]bool)
	stack := append([]string(nil), gate.Needs...)
	var closure []string
	for len(stack) > 0 {
		id := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if seen[id] {
			continue
		}
		seen[id] = true
		closure = append(closure, id)
		stack = append(stack, byID[id].Needs...)
	}
	sort.Strings(closure)
	return qa.TaskID, closure, nil
}
