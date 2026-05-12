package template

import (
	"fmt"
	"strings"
	"text/template"
)

func RenderTemplate(templateStr string, data any) (string, error) {
	tmpl, err := template.New("template").Parse(templateStr)
	if err != nil {
		return "", fmt.Errorf("parsing template: %w", err)
	}

	writer := &strings.Builder{}
	err = tmpl.Execute(writer, data)
	if err != nil {
		return "", fmt.Errorf("executing template: %w", err)
	}

	return writer.String(), nil
}
