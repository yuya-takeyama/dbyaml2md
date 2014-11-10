package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/yuya-takeyama/db2yaml/model"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

var mdTemplate = `# {{.Name}}

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

func main() {
	app := cli.NewApp()
	app.Name = "dbyaml2md"
	app.Usage = "Generate markdown files from YAML generated by db2yaml"
	app.HideHelp = true

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} [options] [arguments...]

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}
`

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "out, o",
			Value: "./dbyaml2md_out",
			Usage: "Directory to output markdown files",
		},
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "Show usage",
		},
	}

	app.Action = func(c *cli.Context) {
		tables := make(map[string]*model.Table)
		buf, err := ioutil.ReadAll(os.Stdin)
		panicIf(err)

		err = yaml.Unmarshal(buf, tables)
		panicIf(err)

		err = generateMarkdownFiles(c, &tables)
		panicIf(err)
	}
	app.Run(os.Args)
}

func generateMarkdownFiles(c *cli.Context, tables *map[string]*model.Table) error {
	out := c.String("out")

	for k, table := range *tables {
		fmt.Fprintf(os.Stderr, "Generating %s.md ...\n", k)

		file, err := os.OpenFile(out+"/"+k+".md", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		defer file.Close()

		err = WriteMarkdownFromTable(file, table)
	}

	fmt.Fprintln(os.Stderr, "Generated all files successfully")

	return nil
}

func WriteMarkdownFromTable(file io.Writer, table *model.Table) error {
	tmpl, err := template.New(table.Name).Parse(mdTemplate)
	if err != nil {
		return err
	}

	err = tmpl.Execute(file, table)
	if err != nil {
		return err
	}

	return nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
