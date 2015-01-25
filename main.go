package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/yuya-takeyama/db2yaml/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var context Context

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
			Usage: "Directory to output markdown files",
		},
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Config file",
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

		var config *Config

		if c.IsSet("config") {
			file, err := os.Open(c.String("config"))
			panicIf(err)

			config, err = LoadConfig(file)
			panicIf(err)
		} else {
			config = NewEmptyConfig()
		}

		context = NewAppContext(c, config)

		err = generateMarkdownFiles(&tables)
		panicIf(err)
	}
	app.Run(os.Args)
}

func generateMarkdownFiles(tables *map[string]*model.Table) error {
	out := context.OutDirectory()

	for k, table := range *tables {
		fmt.Fprintf(os.Stderr, "Generating %s.md ...\n", k)

		file, err := os.OpenFile(out+"/"+k+".md", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		defer file.Close()

		mdWriter := &MdWriter{context.FrontMatter()}

		err = mdWriter.writeMarkdown(file, table)
	}

	fmt.Fprintln(os.Stderr, "Generated all files successfully")

	return nil
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
