package cli

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecommendationCheckIsReachedOnlyFromItsCommands(t *testing.T) {
	t.Parallel()
	root := filepath.Join("..", "..")
	allowed := map[string]bool{
		"internal/config/recommendation_check.go": true,
		"internal/cli/profiles_check.go":          true,
		"internal/cli/profiles_check_apply.go":    true,
		"internal/cli/profiles.go":                true,
		"internal/cli/doctor.go":                  true,
		"internal/cli/upgrade.go":                 true,
	}
	references := 0
	for _, dir := range []string{"internal", "cmd"} {
		err := filepath.WalkDir(filepath.Join(root, dir), func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, path, nil, 0)
			if err != nil {
				return err
			}
			relative, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			relative = filepath.ToSlash(relative)
			ast.Inspect(file, func(node ast.Node) bool {
				ident, ok := node.(*ast.Ident)
				if !ok || ident.Name != "CheckRecommendations" {
					return true
				}
				// A declaration alone is not evidence that a command reaches it.
				if ident.Obj != nil {
					if declaration, ok := ident.Obj.Decl.(*ast.FuncDecl); ok && declaration.Name == ident {
						return true
					}
				}
				references++
				if !allowed[relative] {
					t.Errorf("CheckRecommendations referenced outside its commands: %s", fset.Position(ident.Pos()))
				}
				return true
			})
			return nil
		})
		if err != nil {
			t.Fatalf("parse %s: %v", dir, err)
		}
	}
	if references == 0 {
		t.Fatal("no CheckRecommendations references found")
	}
}
