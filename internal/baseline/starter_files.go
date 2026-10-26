package baseline

import (
	"crypto/sha256"
	"fmt"
)

// StarterFile is a module's repository-owned tooling starter, resolved verbatim.
type StarterFile struct {
	ID              string
	Module          string
	Path            string
	Content         []byte
	ContentIdentity string
}

// StarterFiles returns immutable template bytes in module and declaration order.
func (c *Catalog) StarterFiles(modules []string) ([]StarterFile, error) {
	files := make([]StarterFile, 0)
	for _, moduleID := range modules {
		module, ok := c.modules[moduleID]
		if !ok {
			return nil, fmt.Errorf("resolve starter files: unknown module %q", moduleID)
		}
		for _, declaration := range objectsOrEmpty(module["starterFiles"]) {
			id, _ := stringValue(declaration, "id")
			target, _ := stringValue(declaration, "path")
			templateID, _ := stringValue(declaration, "template")
			asset, ok := c.TemplateContent(templateID)
			if !ok {
				return nil, fmt.Errorf("resolve starter file %q: unknown template %q", id, templateID)
			}
			files = append(files, StarterFile{
				ID: id, Module: moduleID, Path: target, Content: asset.Data,
				ContentIdentity: fmt.Sprintf("sha256:%x", sha256.Sum256(asset.Data)),
			})
		}
	}
	return files, nil
}
