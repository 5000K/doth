package template

import (
	"fmt"
	"strings"
	"text/template"
)

type ConfigMap map[string]any

func RenderTemplate(templateStr string, data ConfigMap) (string, error) {
	tmpl, err := template.New("template").Parse(templateStr)
	tmpl = tmpl.Option("missingkey=zero")
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
