package main

import (
	"github.com/yuya-takeyama/db2yaml/model"
	"gopkg.in/yaml.v2"
	"io"
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
{{range $index, $element := .Columns}}{{$element.Name}}|{{$element.Comment}}|{{$element.Type}}|{{if $element.Length}}{{$element.Length}}{{end}}|{{if $element.Default}}{{$element.Default}}{{end}}|{{if $element.Nullable}}✓{{end}}|{{if $element.AutoIncrement}}✓{{end}}|
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

func (mdWriter *MdWriter) writeMarkdownFromTable(file io.Writer, table *model.Table) error {
	frontMatter := mdWriter.frontMatter
	frontMatter["table"] = table.Name

	frontMatterYaml, err := yaml.Marshal(mdWriter.frontMatterWithTableName(table))
	if err != nil {
		return err
	}

	tmpl, err := template.New(table.Name).Parse("---\n" + string(frontMatterYaml) + mdTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, table)
	if err != nil {
		return err
	}

	return nil
}

func (mdWriter *MdWriter) frontMatterWithTableName(table *model.Table) FrontMatter {
	frontMatter := mdWriter.frontMatter
	frontMatter["table"] = table.Name

	return frontMatter
}

func (mdWriter *MdWriter) frontMatterWithTableNameBytes(table *model.Table) ([]byte, error) {
	return yaml.Marshal(mdWriter.frontMatterWithTableName(table))
}
