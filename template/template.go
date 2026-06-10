package template

import (
	"fmt"
	"strings"
	"text/template"
)

func RenderTemplate(templateStr string, data any) (string, error) {
	tmpl, err := template.New("").Funcs(template.FuncMap{
		"default": func(def, val any) any {
			if val == nil {
				return def
			}
			return val
		},
	}).Parse(templateStr)

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
