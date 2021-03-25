package main

import (
	"fmt"
	"github.com/clok/cdocs"
	"github.com/clok/hev-cli/commands"
	"github.com/clok/kemba"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
)

var (
	version string
	k       = kemba.New("hev")
)

func main() {
	k.Println("executing")

	im, err := cdocs.InstallManpageCommand(&cdocs.InstallManpageCommandInput{
		AppName: "hev",
	})
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "hev"
	app.Version = version
	app.Usage = "scan H-E-B vaccine availability and open a browser when there is one available within a radius of miles"
	app.Commands = []*cli.Command{
		commands.CommandWatch,
		commands.CommandPublish,
		im,
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Print version info",
			Action: func(c *cli.Context) error {
				fmt.Printf("%s %s (%s/%s)\n", "hev", version, runtime.GOOS, runtime.GOARCH)
				return nil
			},
		},
	}

	if os.Getenv("DOCS_MD") != "" {
		docs, err := cdocs.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := cdocs.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
