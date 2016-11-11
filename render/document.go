package render

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"strings"
	"text/template"
)

// A Document is an unrendered template
type Document struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

// Render renders the document from the context values based on the given target
func (d Document) Render(target Target, values ContextValues) (string, error) {

	fMap := make(template.FuncMap)
	fMap["b64encode"] = func(input string) string { return base64.StdEncoding.EncodeToString([]byte(input)) }
	fMap["quote"] = func(input string) string { return fmt.Sprintf(`"%s"`, strings.Replace(input, `"`, `\"`, -1)) }

	tmpl, err := template.New(
		d.Name,
	).Funcs(
		fMap,
	).Option(
		"missingkey=error",
	).Parse(
		d.Content,
	)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	if err := tmpl.Execute(buf, target.GetMatchingValues(values).GetTemplateContext(target.Selector)); err != nil {
		return "", err
	}

	return buf.String(), nil
}
