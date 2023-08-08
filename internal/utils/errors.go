package utils

import (
	"fmt"
	validation "github.com/go-playground/validator/v10"
	"github.com/pterm/pterm"
)

func PrintValidationErrors(errors validation.ValidationErrors) {
	var tree = pterm.TreeNode{
		Text:     "Errors",
		Children: []pterm.TreeNode{},
	}

	for _, err := range errors {
		var field = fmt.Sprintf("%s.%s", err.Namespace(), err.Field())
		tree.Children = append(tree.Children, pterm.TreeNode{
			Text: field,
			Children: []pterm.TreeNode{
				{Text: err.Error()},
			},
		})
	}

	_ = pterm.DefaultTree.WithRoot(tree).Render()
}
