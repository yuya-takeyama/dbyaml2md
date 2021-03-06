package main

import (
	"github.com/yuya-takeyama/db2yaml/model"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
	"text/template"
)

type MdWriter struct {
	frontMatter FrontMatter
}

var mdTemplate = `---
# {{.Name}}

{{.Comment}}

## Columns

Name|Description|Type|Length|Default|Nullable|AUTO_INCREMENT|
----|-----------|----|-----:|-------|-------:|-------------:|
{{range $index, $element := .Columns}}{{$element.Name}}|{{convertLineBreaks $element.Comment}}|{{$element.Type}}|{{if $element.Length}}{{$element.Length}}{{end}}|{{if $element.Default}}{{$element.Default}}{{end}}|{{if $element.Nullable}}✓{{end}}|{{if $element.AutoIncrement}}✓{{end}}|
{{end}}
## Indexes

<table>
  <thead>
    <tr>
      <th>Name</th>
      <th>Unique</th>
      <th>Columns</th>
    </tr>
  </thead>
  <tbody>{{range $index, $element := .Indexes}}
    <tr>
      <td>{{html $element.Name}}</td>
      <td style="text-align: right">{{if $element.Unique}}✓{{end}}</td>
      <td>
        <ul>{{range $cIndex, $cElement := $element.Columns}}
          <li>{{html $cElement.Name}}</li>{{end}}
        </ul>
      </td>
    </tr>{{end}}
  </tbody>
</table>
`

var funcMap = template.FuncMap{
	"convertLineBreaks": func(s string) string {
		s = strings.Replace(s, "\r\n", "<br>", -1)
		s = strings.Replace(s, "\r", "<br>", -1)
		s = strings.Replace(s, "\n", "<br>", -1)
		return s
	},
}

func (mdWriter *MdWriter) writeMarkdown(file io.Writer, table *model.Table) error {
	frontMatterYaml, err := yaml.Marshal(mdWriter.frontMatterWithTable(table))
	if err != nil {
		return err
	}

	tmpl, err := template.New(table.Name).Funcs(funcMap).Parse("---\n" + string(frontMatterYaml) + mdTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, table)
	if err != nil {
		return err
	}

	return nil
}

func (mdWriter *MdWriter) frontMatterWithTable(table *model.Table) FrontMatter {
	frontMatter := mdWriter.frontMatter
	if frontMatter == nil {
		frontMatter = make(FrontMatter)
	}

	frontMatter["table"] = table

	return frontMatter
}
