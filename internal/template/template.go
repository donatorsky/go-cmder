package template

import (
	"io"
	"strings"
	"text/template"
)

const (
	commandTemplate = `package {{ .PackageName }}
{{ if gt (.Imports | len) 0 }}
import ({{ range .Imports }}
	{{ if .Alias }}{{ .Alias }} {{ end }}"{{ .Path }}"{{ end }}
)
{{ end }}
type {{ .CommandName }} struct {{ print "{" }}{{ range .Fields }}
	v{{ .Name | Title }}   {{ .Pointer }}{{ .Type }}
	has{{ .Name | Title }} bool
{{ end }}}
{{range .Constructors }}
{{ . }}
{{ end -}}
{{range .Methods }}
{{ . }}
{{ end }}`

	constructorTemplate = `func New{{ .Name | Title }}({{ if gt (.Fields | len) 0 }}{{ range .Fields }}
	v{{ .Name | Title }} {{ .Pointer }}{{ .Type }},{{ end }}
{{ end }}) {{ if .Mutable }}*{{ end }}{{ .CommandName }} {
	return {{ if .Mutable }}&{{ end }}{{ .CommandName }}{{ print "{" }}{{ if gt (.Fields | len) 0 }}{{ range .Fields }}
		v{{ .Name | Title }}: v{{ .Name | Title }},
		has{{ .Name | Title }}: true,{{ end }}
	}{{ else }}{{ print "}" }}{{ end }}
}`

	getterTemplate = `func (cmd {{ if .Mutable }}*{{ end }}{{ .CommandName }}) {{ .Name | Title }}() {{ .Pointer }}{{ .Type }} {
	return cmd.v{{ .Name | Title }}
}`

	setterTemplate = `func (cmd {{ if .Mutable }}*{{ end }}{{ .CommandName }}) Set{{ .Name | Title }}(v {{ .Pointer }}{{ .Type }}) {{ if .Mutable }}*{{ end }}{{ .CommandName }} {
	cmd.has{{ .Name | Title }} = true
	cmd.v{{ .Name | Title }} = v

	return cmd
}`

	haserTemplate = `func (cmd {{ if .Mutable }}*{{ end }}{{ .CommandName }}) Has{{ .Name | Title }}() bool {
	return cmd.has{{ .Name | Title }}
}`
)

var templateFuncs = template.FuncMap{
	"Title": func(s string) string {
		if len(s) == 0 {
			return ""
		}

		return strings.ToUpper(s[:1]) + s[1:]
	},
}

func NewTemplate() (*Template, error) {
	commandTemplate, err := template.New("command").Funcs(templateFuncs).Parse(commandTemplate)
	if err != nil {
		return nil, err
	}

	constructorTemplate, err := template.New("constructor").Funcs(templateFuncs).Parse(constructorTemplate)
	if err != nil {
		return nil, err
	}

	getterTemplate, err := template.New("getter").Funcs(templateFuncs).Parse(getterTemplate)
	if err != nil {
		return nil, err
	}

	setterTemplate, err := template.New("setter").Funcs(templateFuncs).Parse(setterTemplate)
	if err != nil {
		return nil, err
	}

	haserTemplate, err := template.New("haser").Funcs(templateFuncs).Parse(haserTemplate)
	if err != nil {
		return nil, err
	}

	return &Template{
		commandTemplate:     commandTemplate,
		constructorTemplate: constructorTemplate,
		getterTemplate:      getterTemplate,
		setterTemplate:      setterTemplate,
		haserTemplate:       haserTemplate,
	}, nil
}

type Template struct {
	commandTemplate     *template.Template
	constructorTemplate *template.Template
	getterTemplate      *template.Template
	setterTemplate      *template.Template
	haserTemplate       *template.Template
}

func (t *Template) ExecuteCommandTemplate(writer io.Writer, data *CommandData) error {
	return t.commandTemplate.Execute(writer, data)
}

func (t *Template) ExecuteConstructorTemplate(writer io.Writer, data *ConstructorData) error {
	return t.constructorTemplate.Execute(writer, data)
}

func (t *Template) ExecuteGetterTemplate(writer io.Writer, data *FieldData) error {
	return t.getterTemplate.Execute(writer, data)
}

func (t *Template) ExecuteSetterTemplate(writer io.Writer, data *FieldData) error {
	return t.setterTemplate.Execute(writer, data)
}

func (t *Template) ExecuteHaserTemplate(writer io.Writer, data *FieldData) error {
	return t.haserTemplate.Execute(writer, data)
}
